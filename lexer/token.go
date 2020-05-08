package lexer

import "strconv"

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

	DocumentCommentStart
	DocumentCommentVersion
	DocumentCommentText
	DocumentCommentUnknown
	DocumentCommentStartline
	DocumentCommentEndline
	DocumentCommentTagName

	DocumentCommentTagNameAnchorStart // Anchor for conveniently detecting tags
	AtAuthor
	AtDeprecated
	AtGlobal
	AtLicense
	AtLink
	AtMethod
	AtParam
	AtProperty
	AtPropertyRead
	AtPropertyWrite
	AtReturn
	AtSince
	AtThrows
	AtVar
	DocumentCommentTagNameAnchorEnd

	DocumentCommentEnd

	//Comments whitespace
	Comment
	Whitespace
)

var /* const */ tokenTypeStrings = []string{
	"Undefined",
	"Unknown",
	"EndOfFile",
	"Abstract",
	"Array",
	"As",
	"Break",
	"Callable",
	"Case",
	"Catch",
	"Class",
	"ClassConstant",
	"Clone",
	"Const",
	"Continue",
	"Declare",
	"Default",
	"Do",
	"Echo",
	"Else",
	"ElseIf",
	"Empty",
	"EndDeclare",
	"EndFor",
	"EndForeach",
	"EndIf",
	"EndSwitch",
	"EndWhile",
	"EndHeredoc",
	"Eval",
	"Exit",
	"Extends",
	"Final",
	"Finally",
	"For",
	"ForEach",
	"Function",
	"Global",
	"Goto",
	"HaltCompiler",
	"If",
	"Implements",
	"Include",
	"IncludeOnce",
	"InstanceOf",
	"InsteadOf",
	"Interface",
	"Isset",
	"List",
	"And",
	"Or",
	"Xor",
	"Namespace",
	"New",
	"Print",
	"Private",
	"Public",
	"Protected",
	"Require",
	"RequireOnce",
	"Return",
	"Static",
	"Switch",
	"Throw",
	"Trait",
	"Try",
	"Unset",
	"Use",
	"Var",
	"While",
	"Yield",
	"YieldFrom",
	"DirectoryConstant",
	"FileConstant",
	"LineConstant",
	"FunctionConstant",
	"MethodConstant",
	"NamespaceConstant",
	"TraitConstant",
	"StringLiteral",
	"FloatingLiteral",
	"EncapsulatedAndWhitespace",
	"Text",
	"IntegerLiteral",
	"Name",
	"VariableName",
	"Equals",
	"Tilde",
	"Colon",
	"Semicolon",
	"Exclamation",
	"Dollar",
	"ForwardSlash",
	"Percent",
	"Comma",
	"AtSymbol",
	"Backtick",
	"Question",
	"DoubleQuote",
	"SingleQuote",
	"LessThan",
	"GreaterThan",
	"Asterisk",
	"AmpersandAmpersand",
	"Ampersand",
	"AmpersandEquals",
	"CaretEquals",
	"LessThanLessThan",
	"LessThanLessThanEquals",
	"GreaterThanGreaterThan",
	"GreaterThanGreaterThanEquals",
	"BarEquals",
	"Plus",
	"PlusEquals",
	"AsteriskAsterisk",
	"AsteriskAsteriskEquals",
	"Arrow",
	"OpenBrace",
	"OpenBracket",
	"OpenParenthesis",
	"CloseBrace",
	"CloseBracket",
	"CloseParenthesis",
	"QuestionQuestion",
	"Bar",
	"BarBar",
	"Caret",
	"Dot",
	"DotEquals",
	"CurlyOpen",
	"MinusMinus",
	"ForwardslashEquals",
	"DollarCurlyOpen",
	"FatArrow",
	"ColonColon",
	"Ellipsis",
	"PlusPlus",
	"EqualsEquals",
	"GreaterThanEquals",
	"EqualsEqualsEquals",
	"ExclamationEquals",
	"ExclamationEqualsEquals",
	"LessThanEquals",
	"Spaceship",
	"Minus",
	"MinusEquals",
	"PercentEquals",
	"AsteriskEquals",
	"Backslash",
	"BooleanCast",
	"UnsetCast",
	"StringCast",
	"ObjectCast",
	"IntegerCast",
	"FloatCast",
	"StartHeredoc",
	"ArrayCast",
	"OpenTag",
	"OpenTagEcho",
	"CloseTag",
	"DocumentCommentStart",
	"DocumentCommentVersion",
	"DocumentCommentText",
	"DocumentCommentUnknown",
	"DocumentCommentStartline",
	"DocumentCommentEndline",
	"DocumentCommentTagName",

	"DocumentCommentTagNameAnchorStart",
	"AtAuthor",
	"AtDeprecated",
	"AtGlobal",
	"AtLicense",
	"AtLink",
	"AtMethod",
	"AtParam",
	"AtProperty",
	"AtPropertyRead",
	"AtPropertyWrite",
	"AtReturn",
	"AtSince",
	"AtThrows",
	"AtVar",
	"DocumentCommentTagNameAnchorEnd",

	"DocumentCommentEnd",
	"Comment",
	"Whitespace",
}

func (tokenType TokenType) String() string {
	if int(tokenType) >= len(tokenTypeStrings) {
		return "Unknown"
	}
	return tokenTypeStrings[int(tokenType)]
}

type Token struct {
	Type   TokenType `json:"TokenType"`
	Offset int       `json:"Offset"`
	Length int       `json:"Length"`
}

func NewToken(tokenType TokenType, offset int, length int) *Token {
	return &Token{tokenType, offset, length}
}

// AstNode is a boilerplate for extending interface
func (token Token) AstNode() {
}

func (token Token) String() string {
	str := token.Type.String() + " " + strconv.Itoa(token.Offset) + " " + strconv.Itoa(token.Length)

	return str
}
