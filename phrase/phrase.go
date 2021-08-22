package phrase

import (
	"bytes"

	"github.com/john-nguyen09/go-phpparser/lexer"
)

type PhraseType uint8

const (
	Unknown PhraseType = iota
	AdditiveExpression
	AnonymousClassDeclaration
	AnonymousClassDeclarationHeader
	AnonymousFunctionCreationExpression
	AnonymousFunctionHeader
	AnonymousFunctionUseClause
	AnonymousFunctionUseVariable
	ArrowFunctionCreationExpression
	ArrowFunctionHeader
	ArrowFunctionUseClause
	ArrowFunctionUseVariable
	ArgumentExpressionList
	ArrayCreationExpression
	ArrayElement
	ArrayInitialiserList
	ArrayKey
	ArrayValue
	BitwiseExpression
	BreakStatement
	ByRefAssignmentExpression
	CaseStatement
	CaseStatementList
	CastExpression
	CatchClause
	CatchClauseList
	CatchNameList
	ClassBaseClause
	ClassConstantAccessExpression
	ClassConstDeclaration
	ClassConstElement
	ClassConstElementList
	ClassDeclaration
	ClassDeclarationBody
	ClassDeclarationHeader
	ClassInterfaceClause
	ClassMemberDeclarationList
	ClassModifiers
	ClassTypeDesignator
	CloneExpression
	ClosureUseList
	CoalesceExpression
	CompoundAssignmentExpression
	CompoundStatement
	TernaryExpression
	ConstantAccessExpression
	ConstDeclaration
	ConstElement
	ConstElementList
	ContinueStatement
	DeclareDirective
	DeclareStatement
	DefaultStatement
	DoStatement
	DoubleQuotedStringLiteral
	EchoIntrinsic
	ElseClause
	ElseIfClause
	ElseIfClauseList
	EmptyIntrinsic
	EncapsulatedExpression
	EncapsulatedVariable
	EncapsulatedVariableList
	EqualityExpression
	Error
	ErrorClassMemberDeclaration
	ErrorClassTypeDesignatorAtom
	ErrorControlExpression
	ErrorExpression
	ErrorScopedAccessExpression
	ErrorTraitAdaptation
	ErrorVariable
	ErrorVariableAtom
	EvalIntrinsic
	ExitIntrinsic
	ExponentiationExpression
	ExpressionList
	ExpressionStatement
	FinallyClause
	ForControl
	ForeachCollection
	ForeachKey
	ForeachStatement
	ForeachValue
	ForEndOfLoop
	ForExpressionGroup
	ForInitialiser
	ForStatement
	FullyQualifiedName
	FunctionCallExpression
	FunctionDeclaration
	FunctionDeclarationBody
	FunctionDeclarationHeader
	FunctionStaticDeclaration
	FunctionStaticInitialiser
	GlobalDeclaration
	GotoStatement
	HaltCompilerStatement
	HeredocStringLiteral
	Identifier
	IfStatement
	IncludeExpression
	IncludeOnceExpression
	InlineText
	InstanceOfExpression
	InstanceofTypeDesignator
	InterfaceBaseClause
	InterfaceDeclaration
	InterfaceDeclarationBody
	InterfaceDeclarationHeader
	InterfaceMemberDeclarationList
	IssetIntrinsic
	ListIntrinsic
	LogicalExpression
	MemberModifierList
	MemberName
	MethodCallExpression
	MethodDeclaration
	MethodDeclarationBody
	MethodDeclarationHeader
	MethodReference
	MultiplicativeExpression
	NamedLabelStatement
	NamespaceAliasingClause
	NamespaceDefinition
	NamespaceName
	NamespaceUseClause
	NamespaceUseClauseList
	NamespaceUseDeclaration
	NamespaceUseGroupClause
	NamespaceUseGroupClauseList
	NullStatement
	ObjectCreationExpression
	ParameterDeclaration
	ParameterDeclarationList
	PostfixDecrementExpression
	PostfixIncrementExpression
	PrefixDecrementExpression
	PrefixIncrementExpression
	PrintIntrinsic
	PropertyAccessExpression
	PropertyDeclaration
	PropertyElement
	PropertyElementList
	PropertyInitialiser
	QualifiedName
	QualifiedNameList
	RelationalExpression
	RelativeQualifiedName
	RelativeScope
	RequireExpression
	RequireOnceExpression
	ReturnStatement
	ReturnType
	ScopedCallExpression
	ScopedMemberName
	ScopedPropertyAccessExpression
	ShellCommandExpression
	ShiftExpression
	SimpleAssignmentExpression
	SimpleVariable
	StatementList
	StaticVariableDeclaration
	StaticVariableDeclarationList
	SubscriptExpression
	SwitchStatement
	ThrowStatement
	TraitAdaptationList
	TraitAlias
	TraitDeclaration
	TraitDeclarationBody
	TraitDeclarationHeader
	TraitMemberDeclarationList
	TraitPrecedence
	TraitUseClause
	TraitUseSpecification
	TryStatement
	TypeDeclaration
	UnaryOpExpression
	UnsetIntrinsic
	VariableList
	VariableNameList
	VariadicUnpacking
	WhileStatement
	YieldExpression
	YieldFromExpression
	DocumentComment
	DocumentCommentDescription
	DocumentCommentAuthor
	DocumentCommentEmail

	DocumentCommentTagAnchorStart
	DocumentCommentTag
	DocumentCommentAuthorTag
	DocumentCommentDeprecatedTag
	DocumentCommentGlobalTag
	DocumentCommentMethodTag
	DocumentCommentParamTag
	DocumentCommentPropertyTag
	DocumentCommentReturnTag
	DocumentCommentThrowsTag
	DocumentCommentVarTag
	DocumentCommentTagAnchorEnd

	TypeUnion
	ParameterValue
)

//go:generate stringer -type=PhraseType

type AstNode interface {
	AstNode()
}

type Phrase struct {
	Type     PhraseType `json:"PhraseType"`
	Children []AstNode
}

type ParseError struct {
	Phrase
	Unexpected *lexer.Token
	Expected   lexer.TokenType
}

// AstNode is to extend interface
func (p Phrase) AstNode() {
}

func (phraseType *PhraseType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(phraseType.String())
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func NewPhrase(pool *Pool, phraseType PhraseType, children []AstNode) *Phrase {
	p := pool.Get()
	p.Type = phraseType
	p.Children = children
	return p
}

func NewParseErr(
	pool *Pool,
	phraseType PhraseType,
	children []AstNode,
	unexpected *lexer.Token,
	expected lexer.TokenType) *ParseError {
	phrase := NewPhrase(pool, phraseType, children)

	return &ParseError{*phrase, unexpected, expected}
}
