package lexer

import (
	"bytes"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

type LexerState struct {
	position                 int
	input                    []rune
	inputLength              int
	modeStack                []LexerMode
	doubleQuoteScannedLength int
	heredocLabel             string
}

func NewLexerState(text string, modeStack []LexerMode, position int) *LexerState {
	if modeStack == nil {
		modeStack = []LexerMode{ModeInitial}
	}

	// Consider utf-8 cases since len() returns number of bytes and utf-8 characters can
	// take up more than 1 byte, so we convert string to rune slice first
	characters := []rune(text)

	return &LexerState{position, characters, len(characters), modeStack, -1, ""}
}

func Lex(text string) []*Token {
	lexerState := NewLexerState(text, nil, 0)
	tokens := []*Token{}
	t := lexerState.Lex()
	for {
		tokens = append(tokens, t)
		if t.Type == EndOfFile {
			break
		}
		t = lexerState.Lex()
	}
	return tokens
}

// ModeStack returns a copy of modeStack
func (s *LexerState) ModeStack() []LexerMode {
	modeStack := append(s.modeStack[:0:0], s.modeStack...)

	return modeStack
}

func IsInToken(offset int, t *Token) int {
	if offset < t.Offset {
		return -1
	}
	if offset >= t.Offset+t.Length {
		return 1
	}
	return 0
}

func Sync(text string, change ChangeEvent, oldTokens []*Token) []*Token {
	startIndex := sort.Search(len(oldTokens), func(i int) bool {
		return IsInToken(change.Start, oldTokens[i]) <= 0
	})
	endIndex := sort.Search(len(oldTokens), func(i int) bool {
		return IsInToken(change.End, oldTokens[i]) <= 0
	})

	modeStack := []LexerMode{ModeInitial}
	position := 0
	changedTokens := []*Token{}
	if startIndex > 0 {
		copy(modeStack, oldTokens[startIndex-1].ModeStack)
		position = oldTokens[startIndex-1].Offset
		changedTokens = append(oldTokens[:0:0], oldTokens[:startIndex-1]...)
	}

	lexerState := NewLexerState(text, modeStack, position)
	t := lexerState.Lex()
	lastOffset := 0
	newOffsetDiff := len(change.Text) - (change.End - change.Start)
	for t.Offset < change.End+newOffsetDiff || t.Length == 0 {
		changedTokens = append(changedTokens, t)
		lastOffset = t.Offset + t.Length
		if t.Type == EndOfFile {
			break
		}
		t = lexerState.Lex()
	}
	for _, token := range oldTokens[endIndex:] {
		newToken := *token
		newToken.Offset += len(change.Text) - (change.End - change.Start)

		if newToken.Offset < lastOffset {
			continue
		}
		changedTokens = append(changedTokens, &newToken)
	}
	return changedTokens
}

func (s *LexerState) Lex() *Token {
	if s.position >= s.inputLength {
		return NewToken(EndOfFile, s.position, 0, s.ModeStack())
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

func (tokenType *TokenType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(tokenType.String())
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func (s *LexerState) initial() *Token {
	c := s.input[s.position]
	start := s.position

	if c == '<' && (s.position+1 < s.inputLength && s.input[s.position+1] == '?') {
		if s.position+2 >= s.inputLength || isWhitespace(s.input[s.position+2]) {
			if s.input[s.position+2] == '\r' && s.position+3 < s.inputLength && s.input[s.position+3] == '\n' {
				s.position += 4
			} else {
				s.position += 3
			}

			token := NewToken(OpenTag, start, s.position-start, s.ModeStack())
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
		if s.position+2 < s.inputLength && s.input[s.position+2] == '=' &&
			(s.position+3 >= s.inputLength || isWhitespace(s.input[s.position+3])) {
			if s.input[s.position+3] == '\r' && s.position+4 < s.inputLength && s.input[s.position+4] == '\n' {
				s.position += 5
			} else {
				s.position += 4
			}

			token := NewToken(OpenTagEcho, start, s.position-start, s.ModeStack())
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
		if s.position+5 < s.inputLength && strings.ToLower(string(s.input[s.position:s.position+5])) == "<?php" &&
			(s.position+5 >= s.inputLength || isWhitespace(s.input[s.position+5])) {
			if s.input[s.position+5] == '\r' && s.position+6 < s.inputLength && s.input[s.position+6] == '\n' {
				s.position += 7
			} else {
				s.position += 6
			}

			token := NewToken(OpenTag, start, s.position-start, s.ModeStack())
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			return token
		}
	}

	for s.position++; s.position < s.inputLength; s.position++ {
		c = s.input[s.position]
		if c == '<' && (s.position+1 < s.inputLength && s.input[s.position+1] == '?') {
			if s.position+2 >= s.inputLength || isWhitespace(s.input[s.position+2]) {
				break
			}
			if s.position+2 < s.inputLength && s.input[s.position+2] == '=' &&
				(s.position+3 >= s.inputLength || isWhitespace(s.input[s.position+3])) {
				break
			}
			if s.position+5 < s.inputLength && strings.ToLower(string(s.input[s.position:s.position+5])) == "<?php" &&
				(s.position+5 >= s.inputLength || isWhitespace(s.input[s.position+5])) {
				break
			}
		}
	}

	return NewToken(Text, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scripting() *Token {
	start := s.position
	modeStack := s.ModeStack()

	switch s.input[s.position] {
	case ' ', '\t', '\n', '\r':
		for s.position++; s.position < s.inputLength && isWhitespace(s.input[s.position]); s.position++ {
		}

		return NewToken(Whitespace, start, s.position-start, modeStack)
	case '-':
		return s.scriptingMinus()
	case ':':
		s.position++
		if s.position < s.inputLength && s.input[s.position] == ':' {
			s.position++

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
		s.position++
		if s.position < s.inputLength && s.input[s.position] == '=' {
			s.position++

			return NewToken(PercentEquals, start, 2, modeStack)
		}

		return NewToken(Percent, start, 1, modeStack)
	case '&':
		return s.scriptingAmpersand()
	case '|':
		return s.scriptingBar()
	case '^':
		s.position++
		if s.position < s.inputLength && s.input[s.position] == '=' {
			s.position++

			return NewToken(CaretEquals, start, 2, modeStack)
		}

		return NewToken(Caret, start, 1, modeStack)
	case ';':
		s.position++

		return NewToken(Semicolon, start, 1, modeStack)
	case ',':
		s.position++

		return NewToken(Comma, start, 1, modeStack)
	case '[':
		s.position++

		return NewToken(OpenBracket, start, 1, modeStack)
	case ']':
		s.position++

		return NewToken(CloseBracket, start, 1, modeStack)
	case '(':
		return s.scriptingOpenParenthesis()
	case ')':
		s.position++

		return NewToken(CloseParenthesis, start, 1, modeStack)
	case '~':
		s.position++

		return NewToken(Tilde, start, 1, modeStack)
	case '?':
		return s.scriptingQuestion()
	case '@':
		s.position++

		return NewToken(AtSymbol, start, 1, modeStack)
	case '$':
		return s.scriptingDollar()
	case '#':
		s.position++

		return s.scriptingComment(start)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.scriptingNumeric()

	case '{':
		s.position++

		s.modeStack = append(s.modeStack, ModeScripting)

		return NewToken(OpenBrace, start, 1, modeStack)
	case '}':
		s.position++

		if len(s.modeStack) > 1 {
			s.modeStack = s.modeStack[:len(s.modeStack)-1]
		}

		return NewToken(CloseBrace, start, 1, modeStack)
	case '`':
		s.position++

		s.modeStack[len(s.modeStack)-1] = ModeBacktick

		return NewToken(Backtick, start, 1, modeStack)
	case '\\':
		return s.scriptingBackslash()
	case '\'':
		return s.scriptingSingleQuote(start)
	case '"':
		return s.scriptingDoubleQuote(start)
	}

	if isLabelStart(s.input[s.position]) {
		return s.scriptingLabelStart()
	}

	s.position++

	return NewToken(Unknown, start, 1, modeStack)
}

func (s *LexerState) scriptingMinus() *Token {
	start := s.position
	modeStack := s.ModeStack()

	s.position++
	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '>':
			s.position++
			s.modeStack = append(s.modeStack, ModeLookingForProperty)

			return NewToken(Arrow, start, 2, modeStack)
		case '-':
			s.position++

			return NewToken(MinusMinus, start, 2, modeStack)
		case '=':
			s.position++

			return NewToken(MinusEquals, start, 2, modeStack)
		default:
			break
		}

	}

	return NewToken(Minus, start, 1, modeStack)
}

func (s *LexerState) scriptingDot() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		c := s.input[s.position]
		if c == '=' {
			s.position++

			return NewToken(DotEquals, start, 2, s.ModeStack())
		} else if c == '.' && s.position+1 < s.inputLength && s.input[s.position+1] == '.' {
			s.position += 2

			return NewToken(Ellipsis, start, 3, s.ModeStack())
		} else if c >= '0' && c <= '9' {
			// float
			return s.scriptingNumericStartingWithDotOrE(start, true)
		}
	}

	return NewToken(Dot, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingEquals() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '=':
			s.position++
			if s.position < s.inputLength && s.input[s.position] == '=' {
				s.position++

				return NewToken(EqualsEqualsEquals, start, 3, s.ModeStack())
			}

			return NewToken(EqualsEquals, start, 2, s.ModeStack())
		case '>':
			s.position++

			return NewToken(FatArrow, start, 2, s.ModeStack())
		}
	}

	return NewToken(Equals, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingPlus() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '=':
			s.position++

			return NewToken(PlusEquals, start, 2, s.ModeStack())
		case '+':
			s.position++

			return NewToken(PlusPlus, start, 2, s.ModeStack())
		}

	}

	return NewToken(Plus, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingExclamation() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength && s.input[s.position] == '=' {
		s.position++
		if s.position < s.inputLength && s.input[s.position] == '=' {
			s.position++

			return NewToken(ExclamationEqualsEquals, start, 3, s.ModeStack())
		}

		return NewToken(ExclamationEquals, start, 2, s.ModeStack())
	}

	return NewToken(Exclamation, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingLessThan() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {

		switch s.input[s.position] {
		case '>':
			s.position++

			return NewToken(ExclamationEquals, start, 2, s.ModeStack())
		case '<':
			s.position++
			if s.position < s.inputLength {
				if s.input[s.position] == '=' {
					s.position++

					return NewToken(LessThanLessThanEquals, start, 3, s.ModeStack())
				} else if s.input[s.position] == '<' {
					//go back to first <
					s.position -= 2
					heredoc := s.scriptingHeredoc(start)

					if heredoc != nil {
						return heredoc
					}

					s.position += 2
				}
			}

			return NewToken(LessThanLessThan, start, 2, s.ModeStack())
		case '=':
			s.position++

			if s.position < s.inputLength && s.input[s.position] == '>' {
				s.position++

				return NewToken(Spaceship, start, 3, s.ModeStack())
			}

			return NewToken(LessThanEquals, start, 2, s.ModeStack())
		}

	}

	return NewToken(LessThan, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingGreaterThan() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {

		switch s.input[s.position] {
		case '>':
			s.position++

			if s.position < s.inputLength && s.input[s.position] == '=' {
				s.position++

				return NewToken(GreaterThanGreaterThanEquals, start, 3, s.ModeStack())
			}

			return NewToken(GreaterThanGreaterThan, start, 2, s.ModeStack())
		case '=':
			s.position++

			return NewToken(GreaterThanEquals, start, 2, s.ModeStack())
		}
	}

	return NewToken(GreaterThan, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingAsterisk() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '*':
			s.position++

			if s.position < s.inputLength && s.input[s.position] == '=' {
				s.position++

				return NewToken(AsteriskAsteriskEquals, start, 3, s.ModeStack())
			}

			return NewToken(AsteriskAsterisk, start, 2, s.ModeStack())
		case '=':
			s.position++

			return NewToken(AsteriskEquals, start, 2, s.ModeStack())
		}
	}

	return NewToken(Asterisk, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingForwardSlash() *Token {
	start := s.position
	s.position++

	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '=':
			s.position++

			return NewToken(ForwardslashEquals, start, 2, s.ModeStack())
		case '*':
			s.position++

			return s.scriptingInlineCommentOrDocBlock()
		case '/':
			s.position++

			return s.scriptingComment(start)
		}

	}

	return NewToken(ForwardSlash, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingAmpersand() *Token {
	start := s.position
	s.position++

	if s.position < s.inputLength {

		switch s.input[s.position] {
		case '=':
			s.position++

			return NewToken(AmpersandEquals, start, 2, s.ModeStack())
		case '&':
			s.position++

			return NewToken(AmpersandAmpersand, start, 2, s.ModeStack())
		}

	}

	return NewToken(Ampersand, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingBar() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		switch s.input[s.position] {
		case '=':
			s.position++

			return NewToken(BarEquals, start, 2, s.ModeStack())
		case '|':
			s.position++

			return NewToken(BarBar, start, 2, s.ModeStack())
		}
	}

	return NewToken(Bar, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingOpenParenthesis() *Token {
	start := s.position
	k := start

	//check for cast tokens
	k++
	for k < s.inputLength && (s.input[k] == ' ' || s.input[k] == '\t') {
		k++
	}

	keywordStart := k
	for k < s.inputLength &&
		((s.input[k] >= 'A' && s.input[k] <= 'Z') || (s.input[k] >= 'a' && s.input[k] <= 'z')) {
		k++
	}
	keywordEnd := k

	for k < s.inputLength && (s.input[k] == ' ' || s.input[k] == '\t') {
		k++
	}

	//should have a ) here if valid cast token
	if k < s.inputLength && s.input[k] == ')' {
		keyword := strings.ToLower(string(s.input[keywordStart:keywordEnd]))
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
			s.position = k + 1

			return NewToken(tokenType, start, s.position-start, s.ModeStack())
		}
	}

	s.position++

	return NewToken(OpenParenthesis, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingQuestion() *Token {
	start := s.position

	s.position++
	if s.position < s.inputLength {
		if s.input[s.position] == '?' {
			s.position++

			return NewToken(QuestionQuestion, start, 2, s.ModeStack())
		} else if s.input[s.position] == '>' {
			s.position++
			modeStack := s.ModeStack()

			s.modeStack[len(s.modeStack)-1] = ModeInitial

			return NewToken(CloseTag, start, s.position-start, modeStack)
		}
	}

	return NewToken(Question, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingDollar() *Token {
	start := s.position
	k := s.position

	k++

	if k < s.inputLength && isLabelStart(s.input[k]) {
		for k++; k < s.inputLength && isLabelChar(s.input[k]); k++ {
		}
		s.position = k

		return NewToken(VariableName, start, s.position-start, s.ModeStack())
	}

	s.position++

	return NewToken(Dollar, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingComment(start int) *Token {
	//s.position will be on first char after # or //
	//find first newline or closing tag
	var c rune

	for s.position < s.inputLength {
		c = s.input[s.position]
		s.position++

		if c == '\n' ||
			c == '\r' ||
			(c == '?' && s.position < s.inputLength && s.input[s.position] == '>') {
			s.position--

			break
		}
	}

	return NewToken(Comment, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingNumeric() *Token {
	start := s.position
	k := s.position

	if s.input[s.position] == '0' && k < s.inputLength-1 {
		k++
		j := k + 1
		if s.input[k] == 'b' && j < s.inputLength && (s.input[j] == '0' || s.input[j] == '1') {
			for j++; j < s.inputLength && (s.input[j] == '0' || s.input[j] == '1'); j++ {
			}
			s.position = j

			return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
		}
		if s.input[k] == 'x' && j < s.inputLength && isHexDigit(s.input[j]) {
			for j++; j < s.inputLength && isHexDigit(s.input[j]); j++ {
			}
			s.position = j

			return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
		}
	}

	for s.position++; s.position < s.inputLength && s.input[s.position] >= '0' && s.input[s.position] <= '9'; s.position++ {
	}

	if s.position < len(s.input) {
		if s.input[s.position] == '.' {
			s.position++

			return s.scriptingNumericStartingWithDotOrE(start, true)
		} else if s.input[s.position] == 'e' || s.input[s.position] == 'E' {
			return s.scriptingNumericStartingWithDotOrE(start, false)
		}
	}

	return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingBackslash() *Token {
	//single quote, double quote and heredoc open have optional \
	start := s.position
	s.position++
	var t *Token

	if s.position < s.inputLength {
		switch s.input[s.position] {
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
	}

	return NewToken(Backslash, start, 1, s.ModeStack())
}

func (s *LexerState) scriptingSingleQuote(start int) *Token {
	//optional \ already consumed
	//find first unescaped '
	s.position++
	for {
		if s.position < s.inputLength {
			if s.input[s.position] == '\'' {
				s.position++
				break
			} else if s.input[s.position] == '\\' {
				s.position++
				if s.position < s.inputLength {
					s.position++
				}
			} else {
				s.position++
			}

			continue
		}

		return NewToken(EncapsulatedAndWhitespace, start, s.position-start, s.ModeStack())
	}

	return NewToken(StringLiteral, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingDoubleQuote(start int) *Token {
	//optional \ consumed
	//consume until unescaped "
	//if ${LABEL_START}, ${, {$ found or no match return " and consume none
	s.position++
	n := s.position
	var c rune

	for n < s.inputLength {
		c = s.input[n]
		n++
		switch c {
		case '"':
			s.position = n

			return NewToken(StringLiteral, start, s.position-start, s.ModeStack())
		case '$':
			if n < s.inputLength && (isLabelStart(s.input[n]) || s.input[n] == '{') {
				break
			}
			continue
		case '{':
			if n < s.inputLength && s.input[n] == '$' {
				break
			}
			continue
		case '\\':
			if n < s.inputLength {
				n++
			}
			continue
		/* fall through */
		default:
			continue
		}

		n--
		break
	}

	s.doubleQuoteScannedLength = n
	modeStack := s.ModeStack()
	s.modeStack[len(s.modeStack)-1] = ModeDoubleQuotes

	return NewToken(DoubleQuote, start, s.position-start, modeStack)
}

func (s *LexerState) scriptingLabelStart() *Token {
	start := s.position
	for s.position++; s.position < s.inputLength && isLabelChar(s.input[s.position]); s.position++ {

	}

	firstRune := s.input[start]
	text := string(s.input[start:s.position])
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
			return NewToken(tokenType, start, s.position-start, s.ModeStack())
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
		return NewToken(tokenType, start, s.position-start, s.ModeStack())
	}

	return NewToken(Name, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingNumericStartingWithDotOrE(start int, hasDot bool) *Token {
	for ; s.position < s.inputLength && s.input[s.position] >= '0' && s.input[s.position] <= '9'; s.position++ {

	}

	if s.position < s.inputLength && (s.input[s.position] == 'e' || s.input[s.position] == 'E') {
		k := s.position + 1
		if k < s.inputLength && (s.input[k] == '+' || s.input[k] == '-') {
			k++
		}
		if k < s.inputLength && s.input[k] >= '0' && s.input[k] <= '9' {
			for k++; k < s.inputLength && s.input[k] >= '0' && s.input[k] <= '9'; k++ {
			}
			s.position = k

			return NewToken(FloatingLiteral, start, s.position-start, s.ModeStack())
		}
	}

	tokenType := IntegerLiteral
	if hasDot {
		tokenType = FloatingLiteral
	}

	return NewToken(tokenType, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingHeredoc(start int) *Token {
	//pos is on first <
	k := s.position

	var labelStart int
	var labelEnd int

	for posPlus3 := k + 3; k < posPlus3; k++ {
		if k >= s.inputLength || s.input[k] != '<' {
			return nil
		}
	}

	for ; k < s.inputLength && (s.input[k] == ' ' || s.input[k] == '\t'); k++ {
	}

	var quote rune
	if k < s.inputLength && (s.input[k] == '\'' || s.input[k] == '"') {
		quote = s.input[k]
		k++
	}

	labelStart = k

	if k < s.inputLength && isLabelStart(s.input[k]) {
		for k++; k < s.inputLength && isLabelChar(s.input[k]); k++ {
		}
	} else {
		return nil
	}

	labelEnd = k

	if quote != 0 {
		if k < s.inputLength && s.input[k] == quote {
			k++
		} else {
			return nil
		}
	}

	if k < s.inputLength {
		if s.input[k] == '\r' {
			k++
			if s.input[k] == '\n' {
				k++
			}
		} else if s.input[k] == '\n' {
			k++
		} else {
			return nil
		}
	}

	s.position = k
	s.heredocLabel = string(s.input[labelStart:labelEnd])
	t := NewToken(StartHeredoc, start, s.position-start, s.ModeStack())

	if quote == '\'' {
		s.modeStack[len(s.modeStack)-1] = ModeNowDoc
	} else {
		s.modeStack[len(s.modeStack)-1] = ModeHereDoc
	}

	//check for end on next line
	endHereDocLabel := string(
		s.input[s.position+len(s.heredocLabel) : s.position+len(s.heredocLabel)+3])
	isEndOfLine, err := regexp.MatchString("^;?(?:\r\n|\n|\r)", endHereDocLabel)
	if err == nil &&
		string(s.input[s.position:s.position+len(s.heredocLabel)]) == s.heredocLabel &&
		isEndOfLine {
		s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc
	}

	return t
}

func (s *LexerState) scriptingInlineCommentOrDocBlock() *Token {
	// /* already read
	tokenType := Comment
	start := s.position - 2

	if s.position < s.inputLength &&
		s.input[s.position] == '*' &&
		s.position+1 < s.inputLength && s.input[s.position+1] != '/' {
		s.position++
		tokenType = DocumentComment
	}

	//find comment end */
	for s.position < s.inputLength {
		if s.input[s.position] == '*' &&
			s.position+1 < s.inputLength &&
			s.input[s.position+1] == '/' {
			s.position += 2
			break
		}
		s.position++
	}

	//todo WARN unterminated comment
	return NewToken(tokenType, start, s.position-start, s.ModeStack())
}

func (s *LexerState) scriptingYield(start int) *Token {
	//pos will be after yield keyword
	//check for from
	k := s.position

	if k < s.inputLength && isWhitespace(s.input[k]) {
		for k++; k < s.inputLength && isWhitespace(s.input[k]); k++ {
		}

		if strings.ToLower(string(s.input[k:k+4])) == "from" {
			s.position = k + 4

			return NewToken(YieldFrom, start, s.position-start, s.ModeStack())
		}

	}

	return NewToken(Yield, start, s.position-start, s.ModeStack())
}

func (s *LexerState) lookingForProperty() *Token {
	start := s.position
	modeStack := s.ModeStack()
	c := s.input[s.position]

	if isWhitespace(c) {
		for s.position++; s.position < s.inputLength && isWhitespace(s.input[s.position]); s.position++ {
		}

		return NewToken(Whitespace, start, s.position-start, modeStack)
	}

	if isLabelStart(c) {
		for s.position++; s.position < s.inputLength && isLabelChar(s.input[s.position]); s.position++ {
		}
		s.modeStack = s.modeStack[:len(s.modeStack)-1]

		return NewToken(Name, start, s.position-start, modeStack)
	}

	if c == '-' && s.position+1 < s.inputLength && s.input[s.position+1] == '>' {
		s.position += 2

		return NewToken(Arrow, start, 2, modeStack)
	}

	s.modeStack = s.modeStack[:len(s.modeStack)-1]
	return nil
}

func (s *LexerState) doubleQuotes() *Token {
	c := s.input[s.position]
	start := s.position
	modeStack := s.ModeStack()
	var t *Token

	switch c {
	case '$':
		if t = s.encapsulatedDollar(); t != nil {
			return t
		}
	case '{':
		if s.position+1 < s.inputLength && s.input[s.position+1] == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.position++

			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	case '"':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.position++

		return NewToken(DoubleQuote, start, 1, modeStack)
	}

	return s.doubleQuotesAny()
}

func (s *LexerState) encapsulatedDollar() *Token {
	start := s.position
	k := s.position + 1
	modeStack := s.ModeStack()

	if k >= s.inputLength {
		return nil
	}

	if s.input[k] == '{' {
		s.position += 2
		s.modeStack = append(s.modeStack, ModeLookingForVarName)

		return NewToken(DollarCurlyOpen, start, 2, modeStack)
	}

	if !isLabelStart(s.input[k]) {
		return nil
	}

	for k++; k < s.inputLength && isLabelChar(s.input[k]); k++ {
	}

	if k < s.inputLength && s.input[k] == '[' {
		s.modeStack = append(s.modeStack, ModeVarOffset)
		s.position = k

		return NewToken(VariableName, start, s.position-start, modeStack)
	}

	if k < s.inputLength && s.input[k] == '-' {
		if n := k + 1; n < s.inputLength && s.input[n] == '>' {
			if n++; n < s.inputLength && isLabelStart(s.input[n]) {
				s.modeStack = append(s.modeStack, ModeLookingForProperty)
				s.position = k

				return NewToken(VariableName, start, s.position-start, modeStack)
			}
		}
	}

	s.position = k

	return NewToken(VariableName, start, s.position-start, modeStack)
}

func (s *LexerState) doubleQuotesAny() *Token {
	start := s.position

	if s.doubleQuoteScannedLength > 0 {
		//already know position
		s.position = s.doubleQuoteScannedLength
		s.doubleQuoteScannedLength = -1
	} else {
		//find new pos
		n := s.position + 1

		if s.input[s.position] == '\\' && n+1 < s.inputLength {
			n++
		}

		var c rune
		for n < s.inputLength {
			c = s.input[n]
			n++
			switch c {
			case '"':
				break
			case '$':
				if n < s.inputLength && (isLabelStart(s.input[n]) || s.input[n] == '{') {
					break
				}
				continue
			case '{':
				if n < s.inputLength && s.input[n] == '$' {
					break
				}
				continue
			case '\\':
				if n < s.inputLength {
					n++
				}
				continue
			default:
				continue
			}

			n--
			break
		}

		s.position = n
	}

	return NewToken(EncapsulatedAndWhitespace, start, s.position-start, s.ModeStack())
}

func (s *LexerState) nowdoc() *Token {
	//search for label
	start := s.position
	n := start
	var c rune
	modeStack := s.ModeStack()

	for n < s.inputLength {
		c = s.input[n]
		n++
		switch c {
		case '\r', '\n':
			if c == '\r' && n < s.inputLength && s.input[n] == '\n' {
				n++
			}

			/* Check for ending label on the next line */
			heredocLabel := string(s.input[n : n+len(s.heredocLabel)])
			if n < s.inputLength && s.heredocLabel == heredocLabel {
				k := n + len(s.heredocLabel)

				if k < s.inputLength && s.input[k] == ';' {
					k++
				}

				if k < s.inputLength && (s.input[k] == '\n' || s.input[k] == '\r') {

					//set position to whitespace before label
					nl := string(s.input[n-2 : n])
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

	s.position = n

	return NewToken(EncapsulatedAndWhitespace, start, s.position-start, modeStack)
}

func (s *LexerState) heredoc() *Token {
	c := s.input[s.position]
	start := s.position
	modeStack := s.ModeStack()
	var t *Token

	switch c {
	case '$':
		t = s.encapsulatedDollar()
		if t != nil {
			return t
		}
	case '{':
		if s.position+1 < s.inputLength && s.input[s.position+1] == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.position++

			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	}

	return s.heredocAny()
}

func (s *LexerState) heredocAny() *Token {
	start := s.position
	n := start
	var c rune
	modeStack := s.ModeStack()

	for n < s.inputLength {
		c = s.input[n]
		n++

		switch c {
		case '\r', '\n':
			if c == '\r' && n < s.inputLength && s.input[n] == '\n' {
				n++
			}

			/* Check for ending label on the next line */
			heredocLabel := string(s.input[n : n+len(s.heredocLabel)])
			if n < s.inputLength && heredocLabel == s.heredocLabel {
				k := n + len(s.heredocLabel)

				if k < s.inputLength && s.input[k] == ';' {
					k++
				}

				if k < s.inputLength && (s.input[k] == '\n' || s.input[k] == '\r') {
					nl := string(s.input[n-2 : n])
					if nl == "\r\n" {
						n -= 2
					} else {
						n--
					}

					s.position = n
					s.modeStack[len(s.modeStack)-1] = ModeEndHereDoc

					return NewToken(EncapsulatedAndWhitespace, start, s.position-start, modeStack)
				}
			}

			continue
		case '$':
			if n < s.inputLength && (isLabelStart(s.input[n]) || s.input[n] == '{') {
				break
			}
			continue
		case '{':
			if n < s.inputLength && s.input[n] == '$' {
				break
			}
			continue
		case '\\':
			if n < s.inputLength && s.input[n] != '\n' && s.input[n] != '\r' {
				n++
			}
			continue
		default:
			continue
		}

		n--
		break
	}

	s.position = n

	return NewToken(EncapsulatedAndWhitespace, start, s.position-start, modeStack)
}

func (s *LexerState) endHeredoc() *Token {
	start := s.position
	//consume ws
	for s.position++; s.position < s.inputLength && (s.input[s.position] == '\r' || s.input[s.position] == '\n'); s.position++ {
	}

	s.position += len(s.heredocLabel)
	s.heredocLabel = ""
	t := NewToken(EndHeredoc, start, s.position-start, s.ModeStack())
	s.modeStack[len(s.modeStack)-1] = ModeScripting

	return t
}

func (s *LexerState) backtick() *Token {
	c := s.input[s.position]
	start := s.position
	modeStack := s.ModeStack()
	var t *Token

	switch c {
	case '$':
		t = s.encapsulatedDollar()
		if t != nil {
			return t
		}
	case '{':
		if s.position+1 < s.inputLength && s.input[s.position+1] == '$' {
			s.modeStack = append(s.modeStack, ModeScripting)
			s.position++

			return NewToken(CurlyOpen, start, 1, modeStack)
		}
	case '`':
		s.modeStack[len(s.modeStack)-1] = ModeScripting
		s.position++

		return NewToken(Backtick, start, 1, modeStack)
	}

	return s.backtickAny()
}

func (s *LexerState) backtickAny() *Token {
	n := s.position
	var c rune
	start := n

	if s.input[n] == '\\' && n < s.inputLength {
		n++
	}

	for n < s.inputLength {
		c = s.input[n]
		n++

		switch c {
		case '`':
			break
		case '$':
			if n < s.inputLength && (isLabelStart(s.input[n]) || s.input[n] == '{') {
				break
			}
			continue
		case '{':
			if n < s.inputLength && s.input[n] == '$' {
				break
			}
			continue
		case '\\':
			if n < s.inputLength {
				n++
			}
			continue
		default:
			continue
		}

		n--
		break
	}

	s.position = n

	return NewToken(EncapsulatedAndWhitespace, start, s.position-start, s.modeStack)
}

func (s *LexerState) varOffset() *Token {
	start := s.position
	c := s.input[s.position]

	switch s.input[s.position] {
	case '$':
		if s.position+1 < s.inputLength && isLabelStart(s.input[s.position+1]) {
			s.position++
			for s.position++; s.position < s.inputLength && isLabelChar(s.input[s.position]); s.position++ {
			}

			return NewToken(VariableName, start, s.position-start, s.ModeStack())
		}
	case '[':
		s.position++

		return NewToken(OpenBracket, start, 1, s.ModeStack())
	case ']':
		s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
		s.position++

		return NewToken(CloseBracket, start, 1, s.ModeStack())
	case '-':
		s.position++

		return NewToken(Minus, start, 1, s.ModeStack())
	default:
		if c >= '0' && c <= '9' {
			return s.varOffsetNumeric()
		} else if isLabelStart(c) {
			for s.position++; s.position < s.inputLength && isLabelChar(s.input[s.position]); s.position++ {
			}

			return NewToken(Name, start, s.position-start, s.ModeStack())
		}
	}

	//unexpected char
	modeStack := s.ModeStack()
	s.modeStack = s.modeStack[0 : len(s.modeStack)-1]
	s.position++

	return NewToken(Unknown, start, 1, modeStack)
}

func (s *LexerState) varOffsetNumeric() *Token {
	start := s.position
	c := s.input[s.position]

	if c == '0' {
		k := s.position + 1
		if k < s.inputLength && s.input[k] == 'b' {
			if k++; k < s.inputLength && (s.input[k] == '1' || s.input[k] == '0') {
				for k++; k < s.inputLength && (s.input[k] == '1' || s.input[k] == '0'); k++ {
				}
				s.position = k

				return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
			}
		}

		if k < s.inputLength && s.input[k] == 'x' {
			if k++; k < s.inputLength && isHexDigit(s.input[k]) {
				for k++; k < s.inputLength && isHexDigit(s.input[k]); k++ {
				}
				s.position = k

				return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
			}
		}
	}

	for s.position++; s.position < s.inputLength && s.input[s.position] >= '0' && s.input[s.position] <= '9'; s.position++ {
	}

	return NewToken(IntegerLiteral, start, s.position-start, s.ModeStack())
}

func (s *LexerState) lookingForVarName() *Token {
	start := s.position
	modeStack := s.ModeStack()

	if isLabelStart(s.input[s.position]) {
		k := s.position + 1
		for k++; k < s.inputLength && isLabelChar(s.input[k]); k++ {
		}

		if k < s.inputLength && (s.input[k] == '[' || s.input[k] == '}') {
			s.modeStack[len(s.modeStack)-1] = ModeScripting
			s.position = k

			return NewToken(VariableName, start, s.position-start, modeStack)
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
