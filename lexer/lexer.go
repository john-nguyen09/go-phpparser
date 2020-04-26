package lexer

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type TokenType uint8

const (
	// Misc
	Undefined TokenType = iota
	Unknown
	EndOfFile

	Abstract
	Array
	As
	Break
	Callable
	Case
	Catch
	Class
	ClassConstant
	Clone
	Const
	Continue
	Declare
	Default
	Do
	Echo
	Else
	ElseIf
	Empty
	EndDeclare
	EndFor
	EndForeach
	EndIf
	EndSwitch
	EndWhile
	EndHeredoc
	Eval
	Exit
	Extends
	Final
	Finally
	For
	ForEach
	Function
	Global
	Goto
	HaltCompiler
	If
	Implements
	Include
	IncludeOnce
	InstanceOf
	InsteadOf
	Interface
	Isset
	List
	And
	Or
	Xor
	Namespace
	New
	Print
	Private
	Public
	Protected
	Require
	RequireOnce
	Return
	Static
	Switch
	Throw
	Trait
	Try
	Unset
	Use
	Var
	While
	Yield
	YieldFrom

	//keyword magic constants
	DirectoryConstant
	FileConstant
	LineConstant
	FunctionConstant
	MethodConstant
	NamespaceConstant
	TraitConstant

	//literals
	StringLiteral
	FloatingLiteral
	EncapsulatedAndWhitespace
	Text
	IntegerLiteral

	//Names
	Name
	VariableName

	//Operators and Punctuation
	Equals
	Tilde
	Colon
	Semicolon
	Exclamation
	Dollar
	ForwardSlash
	Percent
	Comma
	AtSymbol
	Backtick
	Question
	DoubleQuote
	SingleQuote
	LessThan
	GreaterThan
	Asterisk
	AmpersandAmpersand
	Ampersand
	AmpersandEquals
	CaretEquals
	LessThanLessThan
	LessThanLessThanEquals
	GreaterThanGreaterThan
	GreaterThanGreaterThanEquals
	BarEquals
	Plus
	PlusEquals
	AsteriskAsterisk
	AsteriskAsteriskEquals
	Arrow
	OpenBrace
	OpenBracket
	OpenParenthesis
	CloseBrace
	CloseBracket
	CloseParenthesis
	QuestionQuestion
	Bar
	BarBar
	Caret
	Dot
	DotEquals
	CurlyOpen
	MinusMinus
	ForwardslashEquals
	DollarCurlyOpen
	FatArrow
	ColonColon
	Ellipsis
	PlusPlus
	EqualsEquals
	GreaterThanEquals
	EqualsEqualsEquals
	ExclamationEquals
	ExclamationEqualsEquals
	LessThanEquals
	Spaceship
	Minus
	MinusEquals
	PercentEquals
	AsteriskEquals
	Backslash
	BooleanCast
	UnsetCast
	StringCast
	ObjectCast
	IntegerCast
	FloatCast
	StartHeredoc
	ArrayCast
	OpenTag
	OpenTagEcho
	CloseTag

	//Comments whitespace
	Comment
	DocumentComment
	Whitespace
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case Unknown:
		return "Unknown"
	case EndOfFile:
		return "EndOfFile"
	case Abstract:
		return "Abstract"
	case Array:
		return "Array"
	case As:
		return "As"
	case Break:
		return "Break"
	case Callable:
		return "Callable"
	case Case:
		return "Case"
	case Catch:
		return "Catch"
	case Class:
		return "Class"
	case ClassConstant:
		return "ClassConstant"
	case Clone:
		return "Clone"
	case Const:
		return "Const"
	case Continue:
		return "Continue"
	case Declare:
		return "Declare"
	case Default:
		return "Default"
	case Do:
		return "Do"
	case Echo:
		return "Echo"
	case Else:
		return "Else"
	case ElseIf:
		return "ElseIf"
	case Empty:
		return "Empty"
	case EndDeclare:
		return "EndDeclare"
	case EndFor:
		return "EndFor"
	case EndForeach:
		return "EndForeach"
	case EndIf:
		return "EndIf"
	case EndSwitch:
		return "EndSwitch"
	case EndWhile:
		return "EndWhile"
	case EndHeredoc:
		return "EndHeredoc"
	case Eval:
		return "Eval"
	case Exit:
		return "Exit"
	case Extends:
		return "Extends"
	case Final:
		return "Final"
	case Finally:
		return "Finally"
	case For:
		return "For"
	case ForEach:
		return "ForEach"
	case Function:
		return "Function"
	case Global:
		return "Global"
	case Goto:
		return "Goto"
	case HaltCompiler:
		return "HaltCompiler"
	case If:
		return "If"
	case Implements:
		return "Implements"
	case Include:
		return "Include"
	case IncludeOnce:
		return "IncludeOnce"
	case InstanceOf:
		return "InstanceOf"
	case InsteadOf:
		return "InsteadOf"
	case Interface:
		return "Interface"
	case Isset:
		return "Isset"
	case List:
		return "List"
	case And:
		return "And"
	case Or:
		return "Or"
	case Xor:
		return "Xor"
	case Namespace:
		return "Namespace"
	case New:
		return "New"
	case Print:
		return "Print"
	case Private:
		return "Private"
	case Public:
		return "Public"
	case Protected:
		return "Protected"
	case Require:
		return "Require"
	case RequireOnce:
		return "RequireOnce"
	case Return:
		return "Return"
	case Static:
		return "Static"
	case Switch:
		return "Switch"
	case Throw:
		return "Throw"
	case Trait:
		return "Trait"
	case Try:
		return "Try"
	case Unset:
		return "Unset"
	case Use:
		return "Use"
	case Var:
		return "Var"
	case While:
		return "While"
	case Yield:
		return "Yield"
	case YieldFrom:
		return "YieldFrom"
	case DirectoryConstant:
		return "DirectoryConstant"
	case FileConstant:
		return "FileConstant"
	case LineConstant:
		return "LineConstant"
	case FunctionConstant:
		return "FunctionConstant"
	case MethodConstant:
		return "MethodConstant"
	case NamespaceConstant:
		return "NamespaceConstant"
	case TraitConstant:
		return "TraitConstant"
	case StringLiteral:
		return "StringLiteral"
	case FloatingLiteral:
		return "FloatingLiteral"
	case EncapsulatedAndWhitespace:
		return "EncapsulatedAndWhitespace"
	case Text:
		return "Text"
	case IntegerLiteral:
		return "IntegerLiteral"
	case Name:
		return "Name"
	case VariableName:
		return "VariableName"
	case Equals:
		return "Equals"
	case Tilde:
		return "Tilde"
	case Colon:
		return "Colon"
	case Semicolon:
		return "Semicolon"
	case Exclamation:
		return "Exclamation"
	case Dollar:
		return "Dollar"
	case ForwardSlash:
		return "ForwardSlash"
	case Percent:
		return "Percent"
	case Comma:
		return "Comma"
	case AtSymbol:
		return "AtSymbol"
	case Backtick:
		return "Backtick"
	case Question:
		return "Question"
	case DoubleQuote:
		return "DoubleQuote"
	case SingleQuote:
		return "SingleQuote"
	case LessThan:
		return "LessThan"
	case GreaterThan:
		return "GreaterThan"
	case Asterisk:
		return "Asterisk"
	case AmpersandAmpersand:
		return "AmpersandAmpersand"
	case Ampersand:
		return "Ampersand"
	case AmpersandEquals:
		return "AmpersandEquals"
	case CaretEquals:
		return "CaretEquals"
	case LessThanLessThan:
		return "LessThanLessThan"
	case LessThanLessThanEquals:
		return "LessThanLessThanEquals"
	case GreaterThanGreaterThan:
		return "GreaterThanGreaterThan"
	case GreaterThanGreaterThanEquals:
		return "GreaterThanGreaterThanEquals"
	case BarEquals:
		return "BarEquals"
	case Plus:
		return "Plus"
	case PlusEquals:
		return "PlusEquals"
	case AsteriskAsterisk:
		return "AsteriskAsterisk"
	case AsteriskAsteriskEquals:
		return "AsteriskAsteriskEquals"
	case Arrow:
		return "Arrow"
	case OpenBrace:
		return "OpenBrace"
	case OpenBracket:
		return "OpenBracket"
	case OpenParenthesis:
		return "OpenParenthesis"
	case CloseBrace:
		return "CloseBrace"
	case CloseBracket:
		return "CloseBracket"
	case CloseParenthesis:
		return "CloseParenthesis"
	case QuestionQuestion:
		return "QuestionQuestion"
	case Bar:
		return "Bar"
	case BarBar:
		return "BarBar"
	case Caret:
		return "Caret"
	case Dot:
		return "Dot"
	case DotEquals:
		return "DotEquals"
	case CurlyOpen:
		return "CurlyOpen"
	case MinusMinus:
		return "MinusMinus"
	case ForwardslashEquals:
		return "ForwardslashEquals"
	case DollarCurlyOpen:
		return "DollarCurlyOpen"
	case FatArrow:
		return "FatArrow"
	case ColonColon:
		return "ColonColon"
	case Ellipsis:
		return "Ellipsis"
	case PlusPlus:
		return "PlusPlus"
	case EqualsEquals:
		return "EqualsEquals"
	case GreaterThanEquals:
		return "GreaterThanEquals"
	case EqualsEqualsEquals:
		return "EqualsEqualsEquals"
	case ExclamationEquals:
		return "ExclamationEquals"
	case ExclamationEqualsEquals:
		return "ExclamationEqualsEquals"
	case LessThanEquals:
		return "LessThanEquals"
	case Spaceship:
		return "Spaceship"
	case Minus:
		return "Minus"
	case MinusEquals:
		return "MinusEquals"
	case PercentEquals:
		return "PercentEquals"
	case AsteriskEquals:
		return "AsteriskEquals"
	case Backslash:
		return "Backslash"
	case BooleanCast:
		return "BooleanCast"
	case UnsetCast:
		return "UnsetCast"
	case StringCast:
		return "StringCast"
	case ObjectCast:
		return "ObjectCast"
	case IntegerCast:
		return "IntegerCast"
	case FloatCast:
		return "FloatCast"
	case StartHeredoc:
		return "StartHeredoc"
	case ArrayCast:
		return "ArrayCast"
	case OpenTag:
		return "OpenTag"
	case OpenTagEcho:
		return "OpenTagEcho"
	case CloseTag:
		return "CloseTag"
	case Comment:
		return "Comment"
	case DocumentComment:
		return "DocumentComment"
	case Whitespace:
		return "Whitespace"
	}

	return ""
}

type LexerMode uint8

const (
	ModeInitial LexerMode = iota
	ModeScripting
	ModeLookingForProperty
	ModeDoubleQuotes
	ModeNowDoc
	ModeHereDoc
	ModeEndHereDoc
	ModeBacktick
	ModeVarOffset
	ModeLookingForVarName
)

func (mode LexerMode) String() string {
	switch mode {
	case ModeInitial:
		return "ModeInitial"
	case ModeScripting:
		return "ModeScripting"
	case ModeLookingForProperty:
		return "ModeLookingForProperty"
	case ModeDoubleQuotes:
		return "ModeDoubleQuotes"
	case ModeNowDoc:
		return "ModeNowDoc"
	case ModeHereDoc:
		return "ModeHereDoc"
	case ModeEndHereDoc:
		return "ModeEndHereDoc"
	case ModeBacktick:
		return "ModeBacktick"
	case ModeVarOffset:
		return "ModeVarOffset"
	case ModeLookingForVarName:
		return "ModeLookingForVarName"
	}

	return ""
}

type Token struct {
	Type      TokenType   `json:"TokenType"`
	Offset    int         `json:"Offset"`
	Length    int         `json:"Length"`
	ModeStack []LexerMode `json:"-"`
}

func NewToken(tokenType TokenType, offset int, length int, modeStack []LexerMode) *Token {
	return &Token{tokenType, offset, length, modeStack}
}

// AstNode is a boilerplate for extending interface
func (token Token) AstNode() {
}

func (token Token) String() string {
	str := token.Type.String() + " " + strconv.Itoa(token.Offset) + " " + strconv.Itoa(token.Length)

	for _, mode := range token.ModeStack {
		str += " " + mode.String()
	}

	return str
}

type Lexer struct {
	offset                   int
	nextOffset               int
	source                   []byte
	modeStack                []LexerMode
	doubleQuoteScannedLength int
	heredocLabel             string
	r                        rune
	decodedQueue             []rune
	decodedSizeQueue         []int
}

func NewLexer(source []byte, modeStack []LexerMode, offset int) *Lexer {
	if modeStack == nil {
		modeStack = []LexerMode{ModeInitial}
	}
	lexer := &Lexer{
		offset:                   offset,
		source:                   source,
		modeStack:                modeStack,
		doubleQuoteScannedLength: -1,
		heredocLabel:             "",
		r:                        0,
		decodedQueue:             []rune{},
		decodedSizeQueue:         []int{},
	}
	lexer.step()
	return lexer
}

func Lex(source []byte) []*Token {
	lexer := NewLexer(source, nil, 0)
	tokens := []*Token{}
	t := lexer.Lex()
	for {
		tokens = append(tokens, t)
		if t.Type == EndOfFile {
			break
		}
		t = lexer.Lex()
	}
	return tokens
}

func (s *Lexer) step() {
	var (
		r    rune
		size int
	)
	if len(s.decodedQueue) != 0 {
		r, size, s.decodedQueue, s.decodedSizeQueue = s.decodedQueue[0], s.decodedSizeQueue[0],
			s.decodedQueue[1:], s.decodedSizeQueue[1:]
	} else {
		r, size = utf8.DecodeRune(s.source[s.nextOffset:])
	}
	if size == 0 {
		r = -1
	}
	s.r = r
	s.offset = s.nextOffset
	s.nextOffset += size
}

func (s *Lexer) stepLoop(n int) {
	for i := 0; i < n; i++ {
		s.step()
	}
}

// Param offset starts with 1 since 0 is just the current r
func (s *Lexer) decodeRuneOffset(offset int) (rune, int) {
	if offset <= 0 {
		return s.r, 0
	}
	if offset-1 < len(s.decodedQueue) {
		totalSize := 0
		for _, size := range s.decodedSizeQueue[:offset-1] {
			totalSize += size
		}
		return s.decodedQueue[offset-1], s.decodedSizeQueue[offset-1] + totalSize
	}
	if offset == 1 {
		r, size := utf8.DecodeRune(s.source[s.nextOffset:])
		s.decodedQueue = append(s.decodedQueue, r)
		s.decodedSizeQueue = append(s.decodedSizeQueue, size)
		return r, size
	}
	totalSize := 0
	if offset-1 >= 1 {
		_, totalSize = s.decodeRuneOffset(offset - 1)
	}
	r, size := utf8.DecodeRune(s.source[s.nextOffset+totalSize:])
	if size == 0 {
		r = -1
	}
	s.decodedQueue = append(s.decodedQueue, r)
	s.decodedSizeQueue = append(s.decodedSizeQueue, size)
	return r, totalSize + size
}

func (s *Lexer) peek(offset int) rune {
	if offset == 0 {
		return s.r
	}
	r, _ := s.decodeRuneOffset(offset)
	return r
}

func (s *Lexer) peekSpanString(offset int, n int) string {
	// Make sure the queue is filled
	s.peek(offset + n)
	return string(s.decodedQueue[offset : offset+n])
}

// ModeStack returns a copy of modeStack
func (s *Lexer) ModeStack() []LexerMode {
	modeStack := append(s.modeStack[:0:0], s.modeStack...)

	return modeStack
}

// Lex runs the lexing and returns a token
func (s *Lexer) Lex() *Token {
	if s.r == -1 {
		return NewToken(EndOfFile, s.offset, 0, s.ModeStack())
	}

	var t *Token

	switch s.modeStack[len(s.modeStack)-1] {
	case ModeInitial:
		t = s.initial()
		break

	case ModeScripting:
		t = s.scripting()
		break

	case ModeLookingForProperty:
		t = s.lookingForProperty()
		break

	case ModeDoubleQuotes:
		t = s.doubleQuotes()
		break

	case ModeNowDoc:
		t = s.nowdoc()
		break

	case ModeHereDoc:
		t = s.heredoc()
		break

	case ModeEndHereDoc:
		t = s.endHeredoc()
		break

	case ModeBacktick:
		t = s.backtick()
		break

	case ModeVarOffset:
		t = s.varOffset()
		break

	case ModeLookingForVarName:
		t = s.lookingForVarName()
		break
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
			token := NewToken(OpenTag, start, s.offset-start, s.ModeStack())
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

			token := NewToken(OpenTagEcho, start, s.offset-start, s.ModeStack())
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
		if strings.ToLower(string(s.peekSpanString(0, 4))) == "?php" && isWhitespace(s.peek(5)) {
			if s.peek(5) == '\r' && s.peek(6) == '\n' {
				s.stepLoop(7)
			} else {
				s.stepLoop(6)
			}
			token := NewToken(OpenTag, start, s.offset-start, s.ModeStack())
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
	return NewToken(Text, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) scripting() *Token {
	start := s.offset
	modeStack := s.ModeStack()

	switch s.r {
	case ' ', '\t', '\n', '\r':
		for s.step(); isWhitespace(s.r); s.step() {
		}

		return NewToken(Whitespace, start, s.offset-start, modeStack)
	case '-':
		return s.scriptingMinus()
	case ':':
		s.step()
		if s.r == ':' {
			s.step()

			return NewToken(ColonColon, start, 2, modeStack)
		}

		return NewToken(Colon, start, 1, modeStack)

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

			return NewToken(PercentEquals, start, 2, modeStack)
		}

		return NewToken(Percent, start, 1, modeStack)
	case '&':
		return s.scriptingAmpersand()
	case '|':
		return s.scriptingBar()
	case '^':
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(CaretEquals, start, 2, modeStack)
		}

		return NewToken(Caret, start, 1, modeStack)
	case ';':
		s.step()

		return NewToken(Semicolon, start, 1, modeStack)
	case ',':
		s.step()

		return NewToken(Comma, start, 1, modeStack)
	case '[':
		s.step()

		return NewToken(OpenBracket, start, 1, modeStack)
	case ']':
		s.step()

		return NewToken(CloseBracket, start, 1, modeStack)
	case '(':
		return s.scriptingOpenParenthesis()
	case ')':
		s.step()

		return NewToken(CloseParenthesis, start, 1, modeStack)
	case '~':
		s.step()

		return NewToken(Tilde, start, 1, modeStack)
	case '?':
		return s.scriptingQuestion()
	case '@':
		s.step()

		return NewToken(AtSymbol, start, 1, modeStack)
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

		return NewToken(OpenBrace, start, 1, modeStack)
	case '}':
		s.step()

		if len(s.modeStack) > 1 {
			s.modeStack = s.modeStack[:len(s.modeStack)-1]
		}

		return NewToken(CloseBrace, start, 1, modeStack)
	case '`':
		s.step()

		s.modeStack[len(s.modeStack)-1] = ModeBacktick

		return NewToken(Backtick, start, 1, modeStack)
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

	return NewToken(Unknown, start, 1, modeStack)
}

func (s *Lexer) scriptingMinus() *Token {
	start := s.offset
	modeStack := s.ModeStack()

	s.step()
	switch s.r {
	case '>':
		s.step()
		s.modeStack = append(s.modeStack, ModeLookingForProperty)

		return NewToken(Arrow, start, 2, modeStack)
	case '-':
		s.step()

		return NewToken(MinusMinus, start, 2, modeStack)
	case '=':
		s.step()

		return NewToken(MinusEquals, start, 2, modeStack)
	}

	return NewToken(Minus, start, 1, modeStack)
}

func (s *Lexer) scriptingDot() *Token {
	start := s.offset

	s.step()
	c := s.r
	if c == '=' {
		s.step()

		return NewToken(DotEquals, start, 2, s.ModeStack())
	} else if c == '.' && s.peek(1) == '.' {
		s.stepLoop(2)

		return NewToken(Ellipsis, start, 3, s.ModeStack())
	} else if c >= '0' && c <= '9' {
		// float
		return s.scriptingNumericStartingWithDotOrE(start, true)
	}

	return NewToken(Dot, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingEquals() *Token {
	start := s.offset

	s.step()
	switch s.r {
	case '=':
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(EqualsEqualsEquals, start, 3, s.ModeStack())
		}

		return NewToken(EqualsEquals, start, 2, s.ModeStack())
	case '>':
		s.step()

		return NewToken(FatArrow, start, 2, s.ModeStack())
	}

	return NewToken(Equals, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingPlus() *Token {
	start := s.offset

	s.step()
	switch s.r {
	case '=':
		s.step()

		return NewToken(PlusEquals, start, 2, s.ModeStack())
	case '+':
		s.step()

		return NewToken(PlusPlus, start, 2, s.ModeStack())
	}

	return NewToken(Plus, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingExclamation() *Token {
	start := s.offset

	s.step()
	if s.r == '=' {
		s.step()
		if s.r == '=' {
			s.step()

			return NewToken(ExclamationEqualsEquals, start, 3, s.ModeStack())
		}

		return NewToken(ExclamationEquals, start, 2, s.ModeStack())
	}

	return NewToken(Exclamation, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingLessThan() *Token {
	start := s.offset

	switch s.peek(1) {
	case '>':
		s.stepLoop(2)
		return NewToken(ExclamationEquals, start, 2, s.ModeStack())
	case '<':
		if s.peek(2) == '=' {
			s.stepLoop(3)
			return NewToken(LessThanLessThanEquals, start, 3, s.ModeStack())
		} else if s.peek(2) == '<' {
			heredoc := s.scriptingHeredoc(start)
			if heredoc != nil {
				return heredoc
			}
			s.stepLoop(2)
		} else {
			s.step()
		}
		return NewToken(LessThanLessThan, start, 2, s.ModeStack())
	case '=':
		s.stepLoop(2)
		if s.r == '>' {
			s.step()
			return NewToken(Spaceship, start, 3, s.ModeStack())
		}
		return NewToken(LessThanEquals, start, 2, s.ModeStack())
	default:
		s.step()
	}
	return NewToken(LessThan, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingGreaterThan() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '>':
		s.step()
		if s.r == '=' {
			s.step()
			return NewToken(GreaterThanGreaterThanEquals, start, 3, s.ModeStack())
		}
		return NewToken(GreaterThanGreaterThan, start, 2, s.ModeStack())
	case '=':
		s.step()
		return NewToken(GreaterThanEquals, start, 2, s.ModeStack())
	}
	return NewToken(GreaterThan, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingAsterisk() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '*':
		s.step()
		if s.r == '=' {
			s.step()
			return NewToken(AsteriskAsteriskEquals, start, 3, s.ModeStack())
		}
		return NewToken(AsteriskAsterisk, start, 2, s.ModeStack())
	case '=':
		s.step()
		return NewToken(AsteriskEquals, start, 2, s.ModeStack())
	}
	return NewToken(Asterisk, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingForwardSlash() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()
		return NewToken(ForwardslashEquals, start, 2, s.ModeStack())
	case '*':
		s.step()
		return s.scriptingInlineCommentOrDocBlock()
	case '/':
		s.step()
		return s.scriptingComment(start)
	}
	return NewToken(ForwardSlash, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingAmpersand() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()
		return NewToken(AmpersandEquals, start, 2, s.ModeStack())
	case '&':
		s.step()
		return NewToken(AmpersandAmpersand, start, 2, s.ModeStack())
	}
	return NewToken(Ampersand, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingBar() *Token {
	start := s.offset
	s.step()
	switch s.r {
	case '=':
		s.step()

		return NewToken(BarEquals, start, 2, s.ModeStack())
	case '|':
		s.step()

		return NewToken(BarBar, start, 2, s.ModeStack())
	}
	return NewToken(Bar, start, 1, s.ModeStack())
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
			return NewToken(tokenType, start, s.offset-start, s.ModeStack())
		}
	}
	s.step()
	return NewToken(OpenParenthesis, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingQuestion() *Token {
	start := s.offset
	s.step()
	if s.r == '?' {
		s.step()

		return NewToken(QuestionQuestion, start, 2, s.ModeStack())
	} else if s.r == '>' {
		s.step()
		modeStack := s.ModeStack()
		s.modeStack[len(s.modeStack)-1] = ModeInitial
		return NewToken(CloseTag, start, s.offset-start, modeStack)
	}
	return NewToken(Question, start, 1, s.ModeStack())
}

func (s *Lexer) scriptingDollar() *Token {
	start := s.offset
	k := 1
	if isLabelStart(s.peek(k)) {
		for k++; isLabelChar(s.peek(k)); k++ {
		}
		s.stepLoop(k)
		return NewToken(VariableName, start, s.offset-start, s.ModeStack())
	}
	s.step()
	return NewToken(Dollar, start, 1, s.ModeStack())
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
	return NewToken(Comment, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) scriptingNumeric() *Token {
	start := s.offset
	if s.r == '0' {
		j := 2
		if s.peek(1) == 'b' && (s.peek(j) == '0' || s.peek(j) == '1') {
			for j++; s.peek(j) == '0' || s.peek(j) == '1'; j++ {
			}
			s.stepLoop(j)
			return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
		}
		if s.peek(1) == 'x' && isHexDigit(s.peek(j)) {
			for j++; isHexDigit(s.peek(j)); j++ {
			}
			s.stepLoop(j)
			return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
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
	return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
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
	return NewToken(Backslash, start, 1, s.ModeStack())
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
		return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, s.ModeStack())
	}
	return NewToken(StringLiteral, start, s.offset-start, s.ModeStack())
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
			return NewToken(StringLiteral, start, s.offset-start, s.ModeStack())
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
	modeStack := s.ModeStack()
	s.modeStack[len(s.modeStack)-1] = ModeDoubleQuotes
	return NewToken(DoubleQuote, start, s.offset-start, modeStack)
}

func (s *Lexer) scriptingLabelStart() *Token {
	start := s.offset
	firstRune := s.r
	for s.step(); isLabelChar(s.r); s.step() {
	}

	text := string(s.source[start:s.offset])
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
			return NewToken(tokenType, start, s.offset-start, s.ModeStack())
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
		tokenType = Array
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
		return NewToken(tokenType, start, s.offset-start, s.ModeStack())
	}
	return NewToken(Name, start, s.offset-start, s.ModeStack())
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
			return NewToken(FloatingLiteral, start, s.offset-start, s.ModeStack())
		}
	}
	tokenType := IntegerLiteral
	if hasDot {
		tokenType = FloatingLiteral
	}
	return NewToken(tokenType, start, s.offset-start, s.ModeStack())
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
	t := NewToken(StartHeredoc, start, s.offset-start, s.ModeStack())
	if quote == '\'' {
		s.modeStack[len(s.modeStack)-1] = ModeNowDoc
	} else {
		s.modeStack[len(s.modeStack)-1] = ModeHereDoc
	}
	//check for end on next line
	endHereDocLabel := string(s.source[s.offset+len(s.heredocLabel) : s.offset+len(s.heredocLabel)+3])
	isEndOfLine, err := regexp.MatchString("^;?(?:\r\n|\n|\r)", endHereDocLabel)
	if err == nil && string(s.source[s.offset:s.offset+len(s.heredocLabel)]) == s.heredocLabel && isEndOfLine {
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
		tokenType = DocumentComment
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
	return NewToken(tokenType, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) scriptingYield(start int) *Token {
	//pos will be after yield keyword
	//check for from
	k := 0

	if isWhitespace(s.peek(k)) {
		for k++; isWhitespace(s.peek(k)); k++ {
		}
		if strings.ToLower(string(s.source[k:k+4])) == "from" {
			s.stepLoop(k + 4)
			return NewToken(YieldFrom, start, s.offset-start, s.ModeStack())
		}
	}
	return NewToken(Yield, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) lookingForProperty() *Token {
	start := s.offset
	modeStack := s.ModeStack()
	c := s.r
	if isWhitespace(c) {
		for s.step(); isWhitespace(s.r); s.step() {
		}
		return NewToken(Whitespace, start, s.offset-start, modeStack)
	}
	if isLabelStart(c) {
		for s.step(); isLabelChar(s.r); s.step() {
		}
		s.modeStack = s.modeStack[:len(s.modeStack)-1]
		return NewToken(Name, start, s.offset-start, modeStack)
	}
	if c == '-' && s.peek(1) == '>' {
		s.stepLoop(2)
		return NewToken(Arrow, start, 2, modeStack)
	}
	s.modeStack = s.modeStack[:len(s.modeStack)-1]
	return nil
}

func (s *Lexer) doubleQuotes() *Token {
	c := s.r
	start := s.offset
	modeStack := s.ModeStack()
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
			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	case '"':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.step()
		return NewToken(DoubleQuote, start, 1, modeStack)
	}
	return s.doubleQuotesAny()
}

func (s *Lexer) encapsulatedDollar() *Token {
	start := s.offset
	k := 1
	modeStack := s.ModeStack()
	if s.peek(k) == -1 {
		return nil
	}
	if s.peek(k) == '{' {
		s.stepLoop(2)
		s.modeStack = append(s.modeStack, ModeLookingForVarName)
		return NewToken(DollarCurlyOpen, start, 2, modeStack)
	}
	if !isLabelStart(s.peek(k)) {
		return nil
	}
	for k++; isLabelChar(s.peek(k)); k++ {
	}
	if s.peek(k) == '[' {
		s.modeStack = append(s.modeStack, ModeVarOffset)
		s.stepLoop(k)
		return NewToken(VariableName, start, s.offset-start, modeStack)
	}
	if s.peek(k) == '-' {
		if n := k + 1; s.peek(n) == '>' {
			if n++; isLabelStart(s.peek(n)) {
				s.modeStack = append(s.modeStack, ModeLookingForProperty)
				s.stepLoop(k)
				return NewToken(VariableName, start, s.offset-start, modeStack)
			}
		}
	}
	s.stepLoop(k)
	return NewToken(VariableName, start, s.offset-start, modeStack)
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
	return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) nowdoc() *Token {
	//search for label
	start := s.offset
	n := 0
	var c rune
	modeStack := s.ModeStack()

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
	return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, modeStack)
}

func (s *Lexer) heredoc() *Token {
	c := s.r
	start := s.offset
	modeStack := s.ModeStack()
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
			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	}
	return s.heredocAny()
}

func (s *Lexer) heredocAny() *Token {
	start := s.offset
	n := 0
	var c rune
	modeStack := s.ModeStack()
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
				if s.peek(k) == '\n' || s.peek(k) == '\r' {
					s.stepLoop(mark)
					s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc
					return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, modeStack)
				}
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
	return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, modeStack)
}

func (s *Lexer) endHeredoc() *Token {
	start := s.offset
	//consume ws
	for s.step(); s.r == '\r' || s.r == '\n'; s.step() {
	}
	s.stepLoop(len(s.heredocLabel))
	s.heredocLabel = ""
	t := NewToken(EndHeredoc, start, s.offset-start, s.ModeStack())
	s.modeStack[len(s.modeStack)-1] = ModeScripting
	return t
}

func (s *Lexer) backtick() *Token {
	c := s.r
	start := s.offset
	modeStack := s.ModeStack()
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
			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	case '`':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.step()
		return NewToken(Backtick, start, 1, modeStack)
	}
	return s.backtickAny()
}

func (s *Lexer) backtickAny() *Token {
	n := s.offset
	var c rune
	start := 0
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
	return NewToken(EncapsulatedAndWhitespace, start, s.offset-start, s.modeStack)
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
			return NewToken(VariableName, start, s.offset-start, s.ModeStack())
		}
	case '[':
		s.step()

		return NewToken(OpenBracket, start, 1, s.ModeStack())
	case ']':
		s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
		s.step()

		return NewToken(CloseBracket, start, 1, s.ModeStack())
	case '-':
		s.step()

		return NewToken(Minus, start, 1, s.ModeStack())
	default:
		if c >= '0' && c <= '9' {
			return s.varOffsetNumeric()
		} else if isLabelStart(c) {
			for s.step(); isLabelChar(s.r); s.step() {
			}
			return NewToken(Name, start, s.offset-start, s.ModeStack())
		}
	}
	//unexpected char
	modeStack := s.ModeStack()
	s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
	s.step()
	return NewToken(Unknown, start, 1, modeStack)
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
				return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
			}
		}
		if s.peek(k) == 'x' {
			if k++; isHexDigit(s.peek(k)) {
				for k++; isHexDigit(s.peek(k)); k++ {
				}
				s.stepLoop(k)
				return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
			}
		}
	}
	for s.step(); s.r >= '0' && s.r <= '9'; s.step() {
	}
	return NewToken(IntegerLiteral, start, s.offset-start, s.ModeStack())
}

func (s *Lexer) lookingForVarName() *Token {
	start := s.offset
	modeStack := s.ModeStack()
	if isLabelStart(s.r) {
		k := 1
		for k++; isLabelChar(s.peek(k)); k++ {
		}
		if s.peek(k) == '[' || s.peek(k) == '}' {
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			s.stepLoop(k)
			return NewToken(VariableName, start, s.offset-start, modeStack)
		}
	}
	s.modeStack[len(s.modeStack)-1] = ModeScripting
	return nil
}

func isLabelStart(cp rune) bool {
	return (cp > 0x40 && cp < 0x5b) || (cp > 0x60 && cp < 0x7b) || cp == 0x5f || cp > 0x7f
}

func isLabelChar(cp rune) bool {
	return (cp > 0x2f && cp < 0x3a) ||
		(cp > 0x40 && cp < 0x5b) ||
		(cp > 0x60 && cp < 0x7b) ||
		cp == 0x5f ||
		cp > 0x7f
}

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
