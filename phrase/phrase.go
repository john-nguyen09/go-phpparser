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
)

func (phraseType PhraseType) String() string {
	switch phraseType {
	case Unknown:
		return "Unknown"
	case AdditiveExpression:
		return "AdditiveExpression"
	case AnonymousClassDeclaration:
		return "AnonymousClassDeclaration"
	case AnonymousClassDeclarationHeader:
		return "AnonymousClassDeclarationHeader"
	case AnonymousFunctionCreationExpression:
		return "AnonymousFunctionCreationExpression"
	case AnonymousFunctionHeader:
		return "AnonymousFunctionHeader"
	case AnonymousFunctionUseClause:
		return "AnonymousFunctionUseClause"
	case AnonymousFunctionUseVariable:
		return "AnonymousFunctionUseVariable"
	case ArgumentExpressionList:
		return "ArgumentExpressionList"
	case ArrayCreationExpression:
		return "ArrayCreationExpression"
	case ArrayElement:
		return "ArrayElement"
	case ArrayInitialiserList:
		return "ArrayInitialiserList"
	case ArrayKey:
		return "ArrayKey"
	case ArrayValue:
		return "ArrayValue"
	case BitwiseExpression:
		return "BitwiseExpression"
	case BreakStatement:
		return "BreakStatement"
	case ByRefAssignmentExpression:
		return "ByRefAssignmentExpression"
	case CaseStatement:
		return "CaseStatement"
	case CaseStatementList:
		return "CaseStatementList"
	case CastExpression:
		return "CastExpression"
	case CatchClause:
		return "CatchClause"
	case CatchClauseList:
		return "CatchClauseList"
	case CatchNameList:
		return "CatchNameList"
	case ClassBaseClause:
		return "ClassBaseClause"
	case ClassConstantAccessExpression:
		return "ClassConstantAccessExpression"
	case ClassConstDeclaration:
		return "ClassConstDeclaration"
	case ClassConstElement:
		return "ClassConstElement"
	case ClassConstElementList:
		return "ClassConstElementList"
	case ClassDeclaration:
		return "ClassDeclaration"
	case ClassDeclarationBody:
		return "ClassDeclarationBody"
	case ClassDeclarationHeader:
		return "ClassDeclarationHeader"
	case ClassInterfaceClause:
		return "ClassInterfaceClause"
	case ClassMemberDeclarationList:
		return "ClassMemberDeclarationList"
	case ClassModifiers:
		return "ClassModifiers"
	case ClassTypeDesignator:
		return "ClassTypeDesignator"
	case CloneExpression:
		return "CloneExpression"
	case ClosureUseList:
		return "ClosureUseList"
	case CoalesceExpression:
		return "CoalesceExpression"
	case CompoundAssignmentExpression:
		return "CompoundAssignmentExpression"
	case CompoundStatement:
		return "CompoundStatement"
	case TernaryExpression:
		return "TernaryExpression"
	case ConstantAccessExpression:
		return "ConstantAccessExpression"
	case ConstDeclaration:
		return "ConstDeclaration"
	case ConstElement:
		return "ConstElement"
	case ConstElementList:
		return "ConstElementList"
	case ContinueStatement:
		return "ContinueStatement"
	case DeclareDirective:
		return "DeclareDirective"
	case DeclareStatement:
		return "DeclareStatement"
	case DefaultStatement:
		return "DefaultStatement"
	case DoStatement:
		return "DoStatement"
	case DoubleQuotedStringLiteral:
		return "DoubleQuotedStringLiteral"
	case EchoIntrinsic:
		return "EchoIntrinsic"
	case ElseClause:
		return "ElseClause"
	case ElseIfClause:
		return "ElseIfClause"
	case ElseIfClauseList:
		return "ElseIfClauseList"
	case EmptyIntrinsic:
		return "EmptyIntrinsic"
	case EncapsulatedExpression:
		return "EncapsulatedExpression"
	case EncapsulatedVariable:
		return "EncapsulatedVariable"
	case EncapsulatedVariableList:
		return "EncapsulatedVariableList"
	case EqualityExpression:
		return "EqualityExpression"
	case Error:
		return "Error"
	case ErrorClassMemberDeclaration:
		return "ErrorClassMemberDeclaration"
	case ErrorClassTypeDesignatorAtom:
		return "ErrorClassTypeDesignatorAtom"
	case ErrorControlExpression:
		return "ErrorControlExpression"
	case ErrorExpression:
		return "ErrorExpression"
	case ErrorScopedAccessExpression:
		return "ErrorScopedAccessExpression"
	case ErrorTraitAdaptation:
		return "ErrorTraitAdaptation"
	case ErrorVariable:
		return "ErrorVariable"
	case ErrorVariableAtom:
		return "ErrorVariableAtom"
	case EvalIntrinsic:
		return "EvalIntrinsic"
	case ExitIntrinsic:
		return "ExitIntrinsic"
	case ExponentiationExpression:
		return "ExponentiationExpression"
	case ExpressionList:
		return "ExpressionList"
	case ExpressionStatement:
		return "ExpressionStatement"
	case FinallyClause:
		return "FinallyClause"
	case ForControl:
		return "ForControl"
	case ForeachCollection:
		return "ForeachCollection"
	case ForeachKey:
		return "ForeachKey"
	case ForeachStatement:
		return "ForeachStatement"
	case ForeachValue:
		return "ForeachValue"
	case ForEndOfLoop:
		return "ForEndOfLoop"
	case ForExpressionGroup:
		return "ForExpressionGroup"
	case ForInitialiser:
		return "ForInitialiser"
	case ForStatement:
		return "ForStatement"
	case FullyQualifiedName:
		return "FullyQualifiedName"
	case FunctionCallExpression:
		return "FunctionCallExpression"
	case FunctionDeclaration:
		return "FunctionDeclaration"
	case FunctionDeclarationBody:
		return "FunctionDeclarationBody"
	case FunctionDeclarationHeader:
		return "FunctionDeclarationHeader"
	case FunctionStaticDeclaration:
		return "FunctionStaticDeclaration"
	case FunctionStaticInitialiser:
		return "FunctionStaticInitialiser"
	case GlobalDeclaration:
		return "GlobalDeclaration"
	case GotoStatement:
		return "GotoStatement"
	case HaltCompilerStatement:
		return "HaltCompilerStatement"
	case HeredocStringLiteral:
		return "HeredocStringLiteral"
	case Identifier:
		return "Identifier"
	case IfStatement:
		return "IfStatement"
	case IncludeExpression:
		return "IncludeExpression"
	case IncludeOnceExpression:
		return "IncludeOnceExpression"
	case InlineText:
		return "InlineText"
	case InstanceOfExpression:
		return "InstanceOfExpression"
	case InstanceofTypeDesignator:
		return "InstanceofTypeDesignator"
	case InterfaceBaseClause:
		return "InterfaceBaseClause"
	case InterfaceDeclaration:
		return "InterfaceDeclaration"
	case InterfaceDeclarationBody:
		return "InterfaceDeclarationBody"
	case InterfaceDeclarationHeader:
		return "InterfaceDeclarationHeader"
	case InterfaceMemberDeclarationList:
		return "InterfaceMemberDeclarationList"
	case IssetIntrinsic:
		return "IssetIntrinsic"
	case ListIntrinsic:
		return "ListIntrinsic"
	case LogicalExpression:
		return "LogicalExpression"
	case MemberModifierList:
		return "MemberModifierList"
	case MemberName:
		return "MemberName"
	case MethodCallExpression:
		return "MethodCallExpression"
	case MethodDeclaration:
		return "MethodDeclaration"
	case MethodDeclarationBody:
		return "MethodDeclarationBody"
	case MethodDeclarationHeader:
		return "MethodDeclarationHeader"
	case MethodReference:
		return "MethodReference"
	case MultiplicativeExpression:
		return "MultiplicativeExpression"
	case NamedLabelStatement:
		return "NamedLabelStatement"
	case NamespaceAliasingClause:
		return "NamespaceAliasingClause"
	case NamespaceDefinition:
		return "NamespaceDefinition"
	case NamespaceName:
		return "NamespaceName"
	case NamespaceUseClause:
		return "NamespaceUseClause"
	case NamespaceUseClauseList:
		return "NamespaceUseClauseList"
	case NamespaceUseDeclaration:
		return "NamespaceUseDeclaration"
	case NamespaceUseGroupClause:
		return "NamespaceUseGroupClause"
	case NamespaceUseGroupClauseList:
		return "NamespaceUseGroupClauseList"
	case NullStatement:
		return "NullStatement"
	case ObjectCreationExpression:
		return "ObjectCreationExpression"
	case ParameterDeclaration:
		return "ParameterDeclaration"
	case ParameterDeclarationList:
		return "ParameterDeclarationList"
	case PostfixDecrementExpression:
		return "PostfixDecrementExpression"
	case PostfixIncrementExpression:
		return "PostfixIncrementExpression"
	case PrefixDecrementExpression:
		return "PrefixDecrementExpression"
	case PrefixIncrementExpression:
		return "PrefixIncrementExpression"
	case PrintIntrinsic:
		return "PrintIntrinsic"
	case PropertyAccessExpression:
		return "PropertyAccessExpression"
	case PropertyDeclaration:
		return "PropertyDeclaration"
	case PropertyElement:
		return "PropertyElement"
	case PropertyElementList:
		return "PropertyElementList"
	case PropertyInitialiser:
		return "PropertyInitialiser"
	case QualifiedName:
		return "QualifiedName"
	case QualifiedNameList:
		return "QualifiedNameList"
	case RelationalExpression:
		return "RelationalExpression"
	case RelativeQualifiedName:
		return "RelativeQualifiedName"
	case RelativeScope:
		return "RelativeScope"
	case RequireExpression:
		return "RequireExpression"
	case RequireOnceExpression:
		return "RequireOnceExpression"
	case ReturnStatement:
		return "ReturnStatement"
	case ReturnType:
		return "ReturnType"
	case ScopedCallExpression:
		return "ScopedCallExpression"
	case ScopedMemberName:
		return "ScopedMemberName"
	case ScopedPropertyAccessExpression:
		return "ScopedPropertyAccessExpression"
	case ShellCommandExpression:
		return "ShellCommandExpression"
	case ShiftExpression:
		return "ShiftExpression"
	case SimpleAssignmentExpression:
		return "SimpleAssignmentExpression"
	case SimpleVariable:
		return "SimpleVariable"
	case StatementList:
		return "StatementList"
	case StaticVariableDeclaration:
		return "StaticVariableDeclaration"
	case StaticVariableDeclarationList:
		return "StaticVariableDeclarationList"
	case SubscriptExpression:
		return "SubscriptExpression"
	case SwitchStatement:
		return "SwitchStatement"
	case ThrowStatement:
		return "ThrowStatement"
	case TraitAdaptationList:
		return "TraitAdaptationList"
	case TraitAlias:
		return "TraitAlias"
	case TraitDeclaration:
		return "TraitDeclaration"
	case TraitDeclarationBody:
		return "TraitDeclarationBody"
	case TraitDeclarationHeader:
		return "TraitDeclarationHeader"
	case TraitMemberDeclarationList:
		return "TraitMemberDeclarationList"
	case TraitPrecedence:
		return "TraitPrecedence"
	case TraitUseClause:
		return "TraitUseClause"
	case TraitUseSpecification:
		return "TraitUseSpecification"
	case TryStatement:
		return "TryStatement"
	case TypeDeclaration:
		return "TypeDeclaration"
	case UnaryOpExpression:
		return "UnaryOpExpression"
	case UnsetIntrinsic:
		return "UnsetIntrinsic"
	case VariableList:
		return "VariableList"
	case VariableNameList:
		return "VariableNameList"
	case VariadicUnpacking:
		return "VariadicUnpacking"
	case WhileStatement:
		return "WhileStatement"
	case YieldExpression:
		return "YieldExpression"
	case YieldFromExpression:
		return "YieldFromExpression"
	}

	return ""
}

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

func NewPhrase(phraseType PhraseType, children []AstNode) *Phrase {
	return &Phrase{phraseType, children}
}

func NewParseErr(
	phraseType PhraseType,
	children []AstNode,
	unexpected *lexer.Token,
	expected lexer.TokenType) *ParseError {
	phrase := NewPhrase(phraseType, children)

	return &ParseError{*phrase, unexpected, expected}
}
