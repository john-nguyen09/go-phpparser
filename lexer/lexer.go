package lexer

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"
)

func decodeAll(bs []byte) ([]rune, []uint8) {
	runes := make([]rune, 0, len(bs))
	sizes := make([]uint8, 0, len(bs))
	lenBs := len(bs)
	for i := 0; i < lenBs; {
		c := bs[i]
		var (
			r    rune
			size int
		)
		if c < utf8.RuneSelf {
			r = rune(c)
			size = 1
		} else {
			r, size = utf8.DecodeRune(bs[i:])
		}
		runes = append(runes, r)
		sizes = append(sizes, uint8(size))
		i += size
	}
	return runes, sizes
}

type Lexer struct {
	offset                   int
	nextOffset               int
	sourceBytes              []byte
	source                   []rune
	sourceSizes              []uint8
	modeStack                []LexerMode
	doubleQuoteScannedLength int
	heredocLabel             string
	r                        rune
	pool                     *Pool
}

func NewLexer(source []byte, modeStack []LexerMode, offset int) *Lexer {
	if modeStack == nil {
		modeStack = []LexerMode{ModeInitial}
	}
	runes, sizes := decodeAll(source)
	lexer := &Lexer{
		offset:                   offset,
		nextOffset:               offset,
		sourceBytes:              source,
		source:                   runes,
		sourceSizes:              sizes,
		modeStack:                modeStack,
		doubleQuoteScannedLength: -1,
		heredocLabel:             "",
		r:                        0,
		pool:                     NewPool(DefaultBlockSize),
	}
	lexer.step()
	return lexer
}

func (s *Lexer) step() {
	if s.nextOffset > len(s.source) {
		return
	}
	if s.nextOffset == len(s.source) {
		s.r = -1
		if s.nextOffset != 0 {
			s.offset++
		}
		s.nextOffset++
		return
	}
	r := s.source[s.nextOffset]
	s.r = r
	if s.nextOffset > 0 {
		s.offset += int(s.sourceSizes[s.nextOffset])
	}
	s.nextOffset++
}

func (s *Lexer) stepLoop(n int) {
	for i := 0; i < n; i++ {
		s.step()
	}
}

func (s *Lexer) peek(offset int) rune {
	if s.nextOffset+offset-1 >= len(s.source) {
		return -1
	}
	if s.nextOffset+offset-1 < 0 {
		return -1
	}
	c := s.source[s.nextOffset+offset-1]
	return c
}

func (s *Lexer) peekSpanString(offset int, n int) string {
	offset += s.nextOffset
	if offset >= len(s.source) {
		return ""
	}
	end := offset + n
	if end >= len(s.source) {
		end = len(s.source) - 1
	}
	return string(s.source[offset:end])
}

// ModeStack returns a copy of modeStack
func (s *Lexer) ModeStack() []LexerMode {
	modeStack := append(s.modeStack[:0:0], s.modeStack...)

	return modeStack
}

// Lex runs the lexing and returns a token
func (s *Lexer) Lex() *Token {
	if s.r == -1 {
		return NewToken(s.pool, EndOfFile, s.offset, 0)
	}

	var t *Token

	switch s.modeStack[len(s.modeStack)-1] {
	case ModeInitial:
		t = s.initial()
	case ModeScripting:
		t = s.scripting()
	case ModeLookingForProperty:
		t = s.lookingForProperty()
	case ModeDoubleQuotes:
		t = s.doubleQuotes()
	case ModeNowDoc:
		t = s.nowdoc()
	case ModeHereDoc:
		t = s.heredoc()
	case ModeEndHereDoc:
		t = s.endHeredoc()
	case ModeBacktick:
		t = s.backtick()
	case ModeVarOffset:
		t = s.varOffset()
	case ModeLookingForVarName:
		t = s.lookingForVarName()
	case ModeDocumentBlock:
		t = s.scriptingDocBlock()
	}

	if t == nil {
		t = s.Lex()
	}

	return t
}

// MarshalJSON marshals the token type into JSON
func (tokenType *TokenType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(tokenType.String())
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func (s *Lexer) initial() *Token {
	start := s.offset
	c := s.r
	if c == '<' && s.peek(1) == '?' {
		if isWhitespace(s.peek(2)) {
			if s.peek(2) == '\r' && s.peek(3) == '\n' {
				s.stepLoop(4)
			} else {
				s.stepLoop(3)
			}
			token := NewToken(s.pool, OpenTag, start, s.offset-start)
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
		if s.peek(2) == '=' &&
			(s.peek(3) == -1 || isWhitespace(s.peek(3))) {
			if s.peek(3) == '\r' && s.peek(4) == '\n' {
				s.stepLoop(5)
			} else {
				s.stepLoop(4)
			}

			token := NewToken(s.pool, OpenTagEcho, start, s.offset-start)
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
		if strings.ToLower(string(s.peekSpanString(0, 4))) == "?php" && isWhitespace(s.peek(5)) {
			if s.peek(5) == '\r' && s.peek(6) == '\n' {
				s.stepLoop(7)
			} else {
				s.stepLoop(6)
			}
			token := NewToken(s.pool, OpenTag, start, s.offset-start)
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
	}
	for s.step(); s.r != -1; s.step() {
		c = s.r
		if c == '<' && (s.peek(1) == '?') {
			if isWhitespace(s.peek(2)) {
				break
			}
			if s.peek(2) == '=' && isWhitespace(s.peek(3)) {
				break
			}
			if strings.ToLower(s.peekSpanString(0, 4)) == "?php" && isWhitespace(s.peek(5)) {
				break
			}
		}
	}
	return NewToken(s.pool, Text, start, s.offset-start)
}

func (s *Lexer) scripting() *Token {
	start := s.offset

	switch s.r {
	case ' ', '\t', '\n', '\r':
		for s.step(); isWhitespace(s.r); s.step() {
		}

		return NewToken(s.pool, Whitespace, start, s.offset-start)
	case '-':
		return s.scriptingMinus()
	case ':':
		s.step()
		if s.r == ':' {
			s.step()

			return NewToken(s.pool, ColonColon, start, 2)
		}

		return NewToken(s.pool, Colon, start, 1)

	case '.':
		return s.scriptingDot()
	case '=':
		return s.scriptingEquals()
	case '+':
		return s.scriptingPlus()
	case '!':
		return s.scriptingExclamation()
	case '<':
		return s.scriptingLessThan()
	case '>':
		return s.scriptingGreaterThan()
	case '*':
		return s.scriptingAsterisk()
	case '/':
		return s.scriptingForwardSlash()
	case '%':
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(s.pool, PercentEquals, start, 2)
		}

		return NewToken(s.pool, Percent, start, 1)
	case '&':
		return s.scriptingAmpersand()
	case '|':
		return s.scriptingBar()
	case '^':
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(s.pool, CaretEquals, start, 2)
		}

		return NewToken(s.pool, Caret, start, 1)
	case ';':
		s.step()

		return NewToken(s.pool, Semicolon, start, 1)
	case ',':
		s.step()

		return NewToken(s.pool, Comma, start, 1)
	case '[':
		s.step()

		return NewToken(s.pool, OpenBracket, start, 1)
	case ']':
		s.step()

		return NewToken(s.pool, CloseBracket, start, 1)
	case '(':
		return s.scriptingOpenParenthesis()
	case ')':
		s.step()

		return NewToken(s.pool, CloseParenthesis, start, 1)
	case '~':
		s.step()

		return NewToken(s.pool, Tilde, start, 1)
	case '?':
		return s.scriptingQuestion()
	case '@':
		s.step()

		return NewToken(s.pool, AtSymbol, start, 1)
	case '$':
		return s.scriptingDollar()
	case '#':
		s.step()

		return s.scriptingComment(start)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.scriptingNumeric()

	case '{':
		s.step()

		s.modeStack = append(s.modeStack, ModeScripting)

		return NewToken(s.pool, OpenBrace, start, 1)
	case '}':
		s.step()

		if len(s.modeStack) > 1 {
			s.modeStack = s.modeStack[:len(s.modeStack)-1]
		}

		return NewToken(s.pool, CloseBrace, start, 1)
	case '`':
		s.step()

		s.modeStack[len(s.modeStack)-1] = ModeBacktick

		return NewToken(s.pool, Backtick, start, 1)
	case '\\':
		return s.scriptingBackslash()
	case '\'':
		return s.scriptingSingleQuote(start)
	case '"':
		return s.scriptingDoubleQuote(start)
	}

	if isLabelStart(s.r) {
		return s.scriptingLabelStart()
	}

	s.step()

	return NewToken(s.pool, Unknown, start, 1)
}

func (s *Lexer) scriptingMinus() *Token {
	start := s.offset

	s.step()
	switch s.r {
	case '>':
		s.step()
		s.modeStack = append(s.modeStack, ModeLookingForProperty)

		return NewToken(s.pool, Arrow, start, 2)
	case '-':
		s.step()

		return NewToken(s.pool, MinusMinus, start, 2)
	case '=':
		s.step()

		return NewToken(s.pool, MinusEquals, start, 2)
	}

	return NewToken(s.pool, Minus, start, 1)
}

func (s *Lexer) scriptingDot() *Token {
	start := s.offset

	s.step()
	c := s.r
	if c == '=' {
		s.step()

		return NewToken(s.pool, DotEquals, start, 2)
	} else if c == '.' && s.peek(1) == '.' {
		s.stepLoop(2)

		return NewToken(s.pool, Ellipsis, start, 3)
	} else if c >= '0' && c <= '9' {
		// float
		return s.scriptingNumericStartingWithDotOrE(start, true)
	}

	return NewToken(s.pool, Dot, start, 1)
}

func (s *Lexer) scriptingEquals() *Token {
	start := s.offset

	s.step()
	switch s.r {
	case '=':
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(s.pool, EqualsEqualsEquals, start, 3)
		}

		return NewToken(s.pool, EqualsEquals, start, 2)
	case '>':
		s.step()

		return NewToken(s.pool, FatArrow, start, 2)
	}

	return NewToken(s.pool, Equals, start, 1)
}

func (s *Lexer) scriptingPlus() *Token {
	start := s.offset

	s.step()
	switch s.r {
	case '=':
		s.step()

		return NewToken(s.pool, PlusEquals, start, 2)
	case '+':
		s.step()

		return NewToken(s.pool, PlusPlus, start, 2)
	}

	return NewToken(s.pool, Plus, start, 1)
}

func (s *Lexer) scriptingExclamation() *Token {
	start := s.offset

	s.step()
	if s.r == '=' {
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(s.pool, ExclamationEqualsEquals, start, 3)
		}

		return NewToken(s.pool, ExclamationEquals, start, 2)
	}

	return NewToken(s.pool, Exclamation, start, 1)
}

func (s *Lexer) scriptingLessThan() *Token {
	start := s.offset

	switch s.peek(1) {
	case '>':
		s.stepLoop(2)
		return NewToken(s.pool, ExclamationEquals, start, 2)
	case '<':
		if s.peek(2) == '=' {
			s.stepLoop(3)
			return NewToken(s.pool, LessThanLessThanEquals, start, 3)
		} else if s.peek(2) == '<' {
			heredoc := s.scriptingHeredoc(start)
			if heredoc != nil {
				return heredoc
			}
			s.stepLoop(2)
		} else {
			s.stepLoop(2)
		}
		return NewToken(s.pool, LessThanLessThan, start, 2)
	case '=':
		s.stepLoop(2)
		if s.r == '>' {
			s.step()
			return NewToken(s.pool, Spaceship, start, 3)
		}
		return NewToken(s.pool, LessThanEquals, start, 2)
	default:
		s.step()
	}
	return NewToken(s.pool, LessThan, start, 1)
}

func (s *Lexer) scriptingGreaterThan() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '>':
		s.step()
		if s.r == '=' {
			s.step()
			return NewToken(s.pool, GreaterThanGreaterThanEquals, start, 3)
		}
		return NewToken(s.pool, GreaterThanGreaterThan, start, 2)
	case '=':
		s.step()
		return NewToken(s.pool, GreaterThanEquals, start, 2)
	}
	return NewToken(s.pool, GreaterThan, start, 1)
}

func (s *Lexer) scriptingAsterisk() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '*':
		s.step()
		if s.r == '=' {
			s.step()
			return NewToken(s.pool, AsteriskAsteriskEquals, start, 3)
		}
		return NewToken(s.pool, AsteriskAsterisk, start, 2)
	case '=':
		s.step()
		return NewToken(s.pool, AsteriskEquals, start, 2)
	}
	return NewToken(s.pool, Asterisk, start, 1)
}

func (s *Lexer) scriptingForwardSlash() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()
		return NewToken(s.pool, ForwardslashEquals, start, 2)
	case '*':
		s.step()
		return s.scriptingInlineCommentOrDocBlock()
	case '/':
		s.step()
		return s.scriptingComment(start)
	}
	return NewToken(s.pool, ForwardSlash, start, 1)
}

func (s *Lexer) scriptingAmpersand() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()
		return NewToken(s.pool, AmpersandEquals, start, 2)
	case '&':
		s.step()
		return NewToken(s.pool, AmpersandAmpersand, start, 2)
	}
	return NewToken(s.pool, Ampersand, start, 1)
}

func (s *Lexer) scriptingBar() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()

		return NewToken(s.pool, BarEquals, start, 2)
	case '|':
		s.step()

		return NewToken(s.pool, BarBar, start, 2)
	}
	return NewToken(s.pool, Bar, start, 1)
}

func (s *Lexer) scriptingOpenParenthesis() *Token {
	start := s.offset
	k := 0
	//check for cast tokens
	for k++; s.peek(k) == ' ' || s.peek(k) == '\t'; k++ {
	}
	keywordStart := k - 1
	for ; (s.peek(k) >= 'A' && s.peek(k) <= 'Z') || (s.peek(k) >= 'a' && s.peek(k) <= 'z'); k++ {
	}
	keywordEnd := k - 1
	for ; s.peek(k) == ' ' || s.peek(k) == '\t'; k++ {
	}

	//should have a ) here if valid cast token
	if s.peek(k) == ')' {
		keyword := strings.ToLower(s.peekSpanString(keywordStart, keywordEnd-keywordStart))
		tokenType := Unknown
		switch keyword {
		case "int", "integer":
			tokenType = IntegerCast
		case "real", "float", "double":
			tokenType = FloatCast
		case "string", "binary":
			tokenType = StringCast
		case "array":
			tokenType = ArrayCast
		case "object":
			tokenType = ObjectCast
		case "bool", "boolean":
			tokenType = BooleanCast
		case "unset":
			tokenType = UnsetCast
		}
		if tokenType > Unknown {
			s.stepLoop(k + 1)
			return NewToken(s.pool, tokenType, start, s.offset-start)
		}
	}
	s.step()
	return NewToken(s.pool, OpenParenthesis, start, 1)
}

func (s *Lexer) scriptingQuestion() *Token {
	start := s.offset
	s.step()
	if s.r == '?' {
		s.step()

		return NewToken(s.pool, QuestionQuestion, start, 2)
	} else if s.r == '>' {
		s.step()
		s.modeStack[len(s.modeStack)-1] = ModeInitial
		return NewToken(s.pool, CloseTag, start, s.offset-start)
	}
	return NewToken(s.pool, Question, start, 1)
}

func (s *Lexer) scriptingDollar() *Token {
	start := s.offset
	k := 1
	if isLabelStart(s.peek(k)) {
		for k++; isLabelChar(s.peek(k)); k++ {
		}
		s.stepLoop(k)
		return NewToken(s.pool, VariableName, start, s.offset-start)
	}
	s.step()
	return NewToken(s.pool, Dollar, start, 1)
}

func (s *Lexer) scriptingComment(start int) *Token {
	//s.position will be on first char after # or //
	//find first newline or closing tag
	for c := s.r; c != -1; {
		if c == '\n' || c == '\r' || (c == '?' && s.peek(1) == '>') {
			break
		}
		s.step()
		c = s.r
	}
	return NewToken(s.pool, Comment, start, s.offset-start)
}

func (s *Lexer) scriptingNumeric() *Token {
	start := s.offset
	if s.r == '0' {
		j := 2
		if s.peek(1) == 'b' && (s.peek(j) == '0' || s.peek(j) == '1') {
			for j++; s.peek(j) == '0' || s.peek(j) == '1'; j++ {
			}
			s.stepLoop(j)
			return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
		}
		if s.peek(1) == 'x' && isHexDigit(s.peek(j)) {
			for j++; isHexDigit(s.peek(j)); j++ {
			}
			s.stepLoop(j)
			return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
		}
	}
	for s.step(); s.r >= '0' && s.r <= '9'; s.step() {
	}
	if s.r == '.' {
		s.step()

		return s.scriptingNumericStartingWithDotOrE(start, true)
	} else if s.r == 'e' || s.r == 'E' {
		return s.scriptingNumericStartingWithDotOrE(start, false)
	}
	return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
}

func (s *Lexer) scriptingBackslash() *Token {
	//single quote, double quote and heredoc open have optional \
	start := s.offset
	s.step()
	var t *Token
	switch s.r {
	case '\'':
		return s.scriptingSingleQuote(start)
	case '"':
		return s.scriptingDoubleQuote(start)
	case '<':
		t = s.scriptingHeredoc(start)
		if t != nil {
			return t
		}
	}
	return NewToken(s.pool, Backslash, start, 1)
}

func (s *Lexer) scriptingSingleQuote(start int) *Token {
	//optional \ already consumed
	//find first unescaped '
	s.step()
	for {
		if s.r != -1 {
			if s.r == '\'' {
				s.step()
				break
			} else if s.r == '\\' {
				s.step()
				s.step()
			} else {
				s.step()
			}
			continue
		}
		return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
	}
	return NewToken(s.pool, StringLiteral, start, s.offset-start)
}

func (s *Lexer) scriptingDoubleQuote(start int) *Token {
	//optional \ consumed
	//consume until unescaped "
	//if ${LABEL_START}, ${, {$ found or no match return " and consume none
	s.step()
	n := 0
	var c rune
	for s.peek(n) != -1 {
		c = s.peek(n)
		n++
		switch c {
		case '"':
			s.stepLoop(n)
			return NewToken(s.pool, StringLiteral, start, s.offset-start)
		case '$':
			if isLabelStart(s.peek(n)) || s.peek(n) == '{' {
				break
			}
			continue
		case '{':
			if s.peek(n) == '$' {
				break
			}
			continue
		case '\\':
			if s.peek(n) != -1 {
				n++
			}
			continue
		default:
			continue
		}
		n--
		break
	}
	s.doubleQuoteScannedLength = n
	s.modeStack[len(s.modeStack)-1] = ModeDoubleQuotes
	return NewToken(s.pool, DoubleQuote, start, s.offset-start)
}

func (s *Lexer) scriptingLabelStart() *Token {
	start := s.offset
	startPosition := s.nextOffset - 1
	firstRune := s.r
	for s.step(); isLabelChar(s.r); s.step() {
	}
	i := 0
	for ; isWhitespace(s.peek(i)); i++ {
	}
	nextNonWhitespace := s.peek(i)

	text := string(s.source[startPosition : s.nextOffset-1])
	tokenType := Unknown
	if firstRune == '_' {
		switch text {
		case "__CLASS__":
			tokenType = ClassConstant
		case "__TRAIT__":
			tokenType = TraitConstant
		case "__FUNCTION__":
			tokenType = FunctionConstant
		case "__METHOD__":
			tokenType = MethodConstant
		case "__LINE__":
			tokenType = LineConstant
		case "__FILE__":
			tokenType = FileConstant
		case "__DIR__":
			tokenType = DirectoryConstant
		case "__NAMESPACE__":
			tokenType = NamespaceConstant
		}
		if tokenType > Unknown {
			return NewToken(s.pool, tokenType, start, s.offset-start)
		}
	}
	text = strings.ToLower(text)
	switch text {
	case "exit":
		tokenType = Exit
	case "die":
		tokenType = Exit
	case "function":
		tokenType = Function
	case "fn":
		tokenType = Fn
	case "const":
		tokenType = Const
	case "return":
		tokenType = Return
	case "yield":
		return s.scriptingYield(start)
	case "try":
		tokenType = Try
	case "catch":
		tokenType = Catch
	case "finally":
		tokenType = Finally
	case "throw":
		tokenType = Throw
	case "if":
		tokenType = If
	case "elseif":
		tokenType = ElseIf
	case "endif":
		tokenType = EndIf
	case "else":
		tokenType = Else
	case "while":
		tokenType = While
	case "endwhile":
		tokenType = EndWhile
	case "do":
		tokenType = Do
	case "for":
		tokenType = For
	case "endfor":
		tokenType = EndFor
	case "foreach":
		tokenType = ForEach
	case "endforeach":
		tokenType = EndForeach
	case "declare":
		tokenType = Declare
	case "enddeclare":
		tokenType = EndDeclare
	case "instanceof":
		tokenType = InstanceOf
	case "as":
		tokenType = As
	case "switch":
		tokenType = Switch
	case "endswitch":
		tokenType = EndSwitch
	case "case":
		tokenType = Case
	case "default":
		tokenType = Default
	case "break":
		tokenType = Break
	case "continue":
		tokenType = Continue
	case "goto":
		tokenType = Goto
	case "echo":
		tokenType = Echo
	case "print":
		tokenType = Print
	case "class":
		tokenType = Class
	case "interface":
		tokenType = Interface
	case "trait":
		tokenType = Trait
	case "extends":
		tokenType = Extends
	case "implements":
		tokenType = Implements
	case "new":
		tokenType = New
	case "clone":
		tokenType = Clone
	case "var":
		tokenType = Var
	case "eval":
		tokenType = Eval
	case "include_once":
		tokenType = IncludeOnce
	case "include":
		tokenType = Include
	case "require_once":
		tokenType = RequireOnce
	case "require":
		tokenType = Require
	case "namespace":
		tokenType = Namespace
	case "use":
		tokenType = Use
	case "insteadof":
		tokenType = InsteadOf
	case "global":
		tokenType = Global
	case "isset":
		tokenType = Isset
	case "empty":
		tokenType = Empty
	case "__halt_compiler":
		tokenType = HaltCompiler
	case "static":
		tokenType = Static
	case "abstract":
		tokenType = Abstract
	case "final":
		tokenType = Final
	case "private":
		tokenType = Private
	case "protected":
		tokenType = Protected
	case "public":
		tokenType = Public
	case "unset":
		tokenType = Unset
	case "list":
		tokenType = List
	case "array":
		if nextNonWhitespace == '(' {
			tokenType = Array
		}
	case "callable":
		tokenType = Callable
	case "or":
		tokenType = Or
	case "and":
		tokenType = And
	case "xor":
		tokenType = Xor
	}
	if tokenType > Unknown {
		return NewToken(s.pool, tokenType, start, s.offset-start)
	}
	return NewToken(s.pool, Name, start, s.offset-start)
}

func (s *Lexer) scriptingNumericStartingWithDotOrE(start int, hasDot bool) *Token {
	for ; s.r >= '0' && s.r <= '9'; s.step() {

	}
	if s.r == 'e' || s.r == 'E' {
		k := 1
		if s.peek(k) == '+' || s.peek(k) == '-' {
			k++
		}
		if s.peek(k) >= '0' && s.peek(k) <= '9' {
			for k++; s.peek(k) >= '0' && s.peek(k) <= '9'; k++ {
			}
			s.stepLoop(k)
			return NewToken(s.pool, FloatingLiteral, start, s.offset-start)
		}
	}
	tokenType := IntegerLiteral
	if hasDot {
		tokenType = FloatingLiteral
	}
	return NewToken(s.pool, tokenType, start, s.offset-start)
}

func (s *Lexer) scriptingHeredoc(start int) *Token {
	//pos is on first <
	k := 0
	var labelStart int
	var labelEnd int
	for posPlus3 := k + 3; k < posPlus3; k++ {
		if s.peek(k) != '<' {
			return nil
		}
	}
	for ; s.peek(k) == ' ' || s.peek(k) == '\t'; k++ {
	}
	var quote rune
	if s.peek(k) == '\'' || s.peek(k) == '"' {
		quote = s.peek(k)
		k++
	}
	labelStart = k
	if isLabelStart(s.peek(k)) {
		for k++; isLabelChar(s.peek(k)); k++ {
		}
	} else {
		return nil
	}
	labelEnd = k
	if quote != 0 {
		if s.peek(k) == quote {
			k++
		} else {
			return nil
		}
	}
	if s.peek(k) == '\r' {
		k++
		if s.peek(k) == '\n' {
			k++
		}
	} else if s.peek(k) == '\n' {
		k++
	} else {
		return nil
	}
	s.heredocLabel = s.peekSpanString(labelStart-1, labelEnd-labelStart)
	s.stepLoop(k)
	t := NewToken(s.pool, StartHeredoc, start, s.offset-start)
	if quote == '\'' {
		s.modeStack[len(s.modeStack)-1] = ModeNowDoc
	} else {
		s.modeStack[len(s.modeStack)-1] = ModeHereDoc
	}
	//check for end on next line
	endHereDocLabel := string(s.source[s.nextOffset-1+len(s.heredocLabel) : s.nextOffset-1+len(s.heredocLabel)+3])
	isEndOfLine, err := regexp.MatchString("^;?(?:\r\n|\n|\r)", endHereDocLabel)
	if err == nil && string(s.source[s.nextOffset-1:s.nextOffset-1+len(s.heredocLabel)]) == s.heredocLabel && isEndOfLine {
		s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc
	}
	return t
}

func (s *Lexer) scriptingInlineCommentOrDocBlock() *Token {
	// /* already read
	tokenType := Comment
	start := s.offset - 2
	if s.r == '*' && s.peek(1) != '/' {
		s.step()
		s.modeStack = append(s.modeStack, ModeDocumentBlock)
		return NewToken(s.pool, DocumentCommentStart, start, s.offset-start)
	}
	//find comment end */
	for s.r != -1 {
		if s.r == '*' && s.peek(1) == '/' {
			s.stepLoop(2)
			break
		}
		s.step()
	}
	// TODO: WARN unterminated comment
	return NewToken(s.pool, tokenType, start, s.offset-start)
}

func (s *Lexer) scriptingYield(start int) *Token {
	//pos will be after yield keyword
	//check for from
	k := 0

	if isWhitespace(s.peek(k)) {
		for k++; isWhitespace(s.peek(k)); k++ {
		}
		if strings.ToLower(s.peekSpanString(k-1, 4)) == "from" {
			s.stepLoop(k + 4)
			return NewToken(s.pool, YieldFrom, start, s.offset-start)
		}
	}
	return NewToken(s.pool, Yield, start, s.offset-start)
}

func (s *Lexer) lookingForProperty() *Token {
	start := s.offset
	c := s.r
	if isWhitespace(c) {
		for s.step(); isWhitespace(s.r); s.step() {
		}
		return NewToken(s.pool, Whitespace, start, s.offset-start)
	}
	if isLabelStart(c) {
		for s.step(); isLabelChar(s.r); s.step() {
		}
		s.modeStack = s.modeStack[:len(s.modeStack)-1]
		return NewToken(s.pool, Name, start, s.offset-start)
	}
	if c == '-' && s.peek(1) == '>' {
		s.stepLoop(2)
		return NewToken(s.pool, Arrow, start, 2)
	}
	s.modeStack = s.modeStack[:len(s.modeStack)-1]
	return nil
}

func (s *Lexer) doubleQuotes() *Token {
	c := s.r
	start := s.offset
	var t *Token

	switch c {
	case '$':
		if t = s.encapsulatedDollar(); t != nil {
			return t
		}
	case '{':
		if s.peek(1) == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.step()
			return NewToken(s.pool, CurlyOpen, start, 1)
		}
	case '"':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.step()
		return NewToken(s.pool, DoubleQuote, start, 1)
	}
	return s.doubleQuotesAny()
}

func (s *Lexer) encapsulatedDollar() *Token {
	start := s.offset
	k := 1
	if s.peek(k) == -1 {
		return nil
	}
	if s.peek(k) == '{' {
		s.stepLoop(2)
		s.modeStack = append(s.modeStack, ModeLookingForVarName)
		return NewToken(s.pool, DollarCurlyOpen, start, 2)
	}
	if !isLabelStart(s.peek(k)) {
		return nil
	}
	for k++; isLabelChar(s.peek(k)); k++ {
	}
	if s.peek(k) == '[' {
		s.modeStack = append(s.modeStack, ModeVarOffset)
		s.stepLoop(k)
		return NewToken(s.pool, VariableName, start, s.offset-start)
	}
	if s.peek(k) == '-' {
		if n := k + 1; s.peek(n) == '>' {
			if n++; isLabelStart(s.peek(n)) {
				s.modeStack = append(s.modeStack, ModeLookingForProperty)
				s.stepLoop(k)
				return NewToken(s.pool, VariableName, start, s.offset-start)
			}
		}
	}
	s.stepLoop(k)
	return NewToken(s.pool, VariableName, start, s.offset-start)
}

func (s *Lexer) doubleQuotesAny() *Token {
	start := s.offset
	if s.doubleQuoteScannedLength > 0 {
		//already know position
		s.stepLoop(s.doubleQuoteScannedLength)
		s.doubleQuoteScannedLength = -1
	} else {
		//find new pos
		n := 1
		if s.r == '\\' && s.peek(n+1) != -1 {
			n++
		}
		var c rune
		for s.peek(n) != -1 {
			c = s.peek(n)
			n++
			switch c {
			case '"':
				break
			case '$':
				if isLabelStart(s.peek(n)) || s.peek(n) == '{' {
					break
				}
				continue
			case '{':
				if s.peek(n) == '$' {
					break
				}
				continue
			case '\\':
				if s.peek(n) != -1 {
					n++
				}
				continue
			default:
				continue
			}

			n--
			break
		}
		s.stepLoop(n)
	}
	return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
}

func (s *Lexer) nowdoc() *Token {
	//search for label
	start := s.offset
	n := 0
	var c rune

	for s.peek(n) != -1 {
		c = s.peek(n)
		n++
		switch c {
		case '\r', '\n':
			if c == '\r' && s.peek(n) == '\n' {
				n++
			}
			/* Check for ending label on the next line */
			heredocLabel := s.peekSpanString(n, len(s.heredocLabel))
			if s.peek(n) != -1 && s.heredocLabel == heredocLabel {
				k := n + len(s.heredocLabel)
				if s.peek(k) == ';' {
					k++
				}
				if s.peek(k) == '\n' || s.peek(k) == '\r' {
					//set position to whitespace before label
					nl := s.peekSpanString(n-2, 2)
					if nl == "\r\n" {
						n -= 2
					} else {
						n--
					}
					s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc
					break
				}
			}
			continue
		default:
			continue
		}

		break
	}
	s.stepLoop(n)
	return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
}

func (s *Lexer) heredoc() *Token {
	c := s.r
	start := s.offset
	var t *Token
	switch c {
	case '$':
		t = s.encapsulatedDollar()
		if t != nil {
			return t
		}
	case '{':
		if s.peek(1) == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.step()
			return NewToken(s.pool, CurlyOpen, start, 1)
		}
	}
	return s.heredocAny()
}

func (s *Lexer) heredocAny() *Token {
	start := s.offset
	n := 0
	var c rune
	for s.peek(n) != -1 {
		c = s.peek(n)
		n++
		switch c {
		case '\r', '\n':
			mark := n - 1
			if c == '\r' && s.peek(n) == '\n' {
				n++
			}
			/* Check for ending label on the next line */
			heredocLabel := s.peekSpanString(n-1, len(s.heredocLabel))
			if s.peek(n) != -1 && heredocLabel == s.heredocLabel {
				k := n + len(s.heredocLabel)
				if s.peek(k) == ';' {
					k++
				}
				s.stepLoop(mark)
				s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc
				return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
			}
			continue
		case '$':
			if isLabelStart(s.peek(n)) || s.peek(n) == '{' {
				break
			}
			continue
		case '{':
			if s.peek(n) == '$' {
				break
			}
			continue
		case '\\':
			if s.peek(n) != '\n' && s.peek(n) != '\r' {
				n++
			}
			continue
		default:
			continue
		}
		n--
		break
	}
	s.stepLoop(n)
	return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
}

func (s *Lexer) endHeredoc() *Token {
	start := s.offset
	//consume ws
	for s.step(); s.r == '\r' || s.r == '\n'; s.step() {
	}
	s.stepLoop(len(s.heredocLabel))
	s.heredocLabel = ""
	t := NewToken(s.pool, EndHeredoc, start, s.offset-start)
	s.modeStack[len(s.modeStack)-1] = ModeScripting
	return t
}

func (s *Lexer) backtick() *Token {
	c := s.r
	start := s.offset
	var t *Token
	switch c {
	case '$':
		t = s.encapsulatedDollar()
		if t != nil {
			return t
		}
	case '{':
		if s.peek(1) == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.step()
			return NewToken(s.pool, CurlyOpen, start, 1)
		}
	case '`':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.step()
		return NewToken(s.pool, Backtick, start, 1)
	}
	return s.backtickAny()
}

func (s *Lexer) backtickAny() *Token {
	n := 0
	var c rune
	start := s.offset
	if s.peek(n) == '\\' && s.peek(n+1) != -1 {
		n++
	}
	for s.peek(n) != -1 {
		c = s.peek(n)
		n++
		switch c {
		case '`':
			break
		case '$':
			if isLabelStart(s.peek(n)) || s.peek(n) == '{' {
				break
			}
			continue
		case '{':
			if s.peek(n) == '$' {
				break
			}
			continue
		case '\\':
			if s.peek(n) != -1 {
				n++
			}
			continue
		default:
			continue
		}
		n--
		break
	}
	s.stepLoop(n)
	return NewToken(s.pool, EncapsulatedAndWhitespace, start, s.offset-start)
}

func (s *Lexer) varOffset() *Token {
	start := s.offset
	c := s.r
	switch s.r {
	case '$':
		if isLabelStart(s.peek(1)) {
			s.step()
			for s.step(); isLabelChar(s.r); s.step() {
			}
			return NewToken(s.pool, VariableName, start, s.offset-start)
		}
	case '[':
		s.step()

		return NewToken(s.pool, OpenBracket, start, 1)
	case ']':
		s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
		s.step()

		return NewToken(s.pool, CloseBracket, start, 1)
	case '-':
		s.step()

		return NewToken(s.pool, Minus, start, 1)
	default:
		if c >= '0' && c <= '9' {
			return s.varOffsetNumeric()
		} else if isLabelStart(c) {
			for s.step(); isLabelChar(s.r); s.step() {
			}
			return NewToken(s.pool, Name, start, s.offset-start)
		}
	}
	//unexpected char
	s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
	s.step()
	return NewToken(s.pool, Unknown, start, 1)
}

func (s *Lexer) varOffsetNumeric() *Token {
	start := s.offset
	c := s.r
	if c == '0' {
		k := 1
		if s.peek(k) == 'b' {
			if k++; s.peek(k) == '1' || s.peek(k) == '0' {
				for k++; s.peek(k) == '1' || s.peek(k) == '0'; k++ {
				}
				s.stepLoop(k)
				return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
			}
		}
		if s.peek(k) == 'x' {
			if k++; isHexDigit(s.peek(k)) {
				for k++; isHexDigit(s.peek(k)); k++ {
				}
				s.stepLoop(k)
				return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
			}
		}
	}
	for s.step(); s.r >= '0' && s.r <= '9'; s.step() {
	}
	return NewToken(s.pool, IntegerLiteral, start, s.offset-start)
}

func (s *Lexer) lookingForVarName() *Token {
	start := s.offset
	if isLabelStart(s.r) {
		k := 1
		for k++; isLabelChar(s.peek(k)); k++ {
		}
		if s.peek(k) == '[' || s.peek(k) == '}' {
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			s.stepLoop(k)
			return NewToken(s.pool, VariableName, start, s.offset-start)
		}
	}
	s.modeStack[len(s.modeStack)-1] = ModeScripting
	return nil
}

// GetTokenValue returns the []rune of the token source code
func (s *Lexer) GetTokenValue(t *Token) []rune {
	return []rune(string(s.sourceBytes[t.Offset : t.Offset+t.Length]))
}

func isLabelStart(cp rune) bool {
	return (cp >= 'A' && cp <= 'Z') || (cp >= 'a' && cp <= 'z') || cp == '_' || cp >= utf8.RuneSelf
}

func isLabelChar(cp rune) bool {
	return (cp >= '0' && cp <= '9') ||
		(cp >= 'A' && cp <= 'Z') ||
		(cp >= 'a' && cp <= 'z') ||
		cp == '_' ||
		cp >= utf8.RuneSelf
}

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c rune) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
