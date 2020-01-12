package parser

import (
	"errors"
	"reflect"
)

var statementListRecoverSet = []TokenType{Use,
	HaltCompiler,
	Const,
	Function,
	Class,
	Abstract,
	Final,
	Trait,
	Interface,
	OpenBrace,
	If,
	While,
	Do,
	For,
	Switch,
	Break,
	Continue,
	Return,
	Global,
	Static,
	Echo,
	Unset,
	ForEach,
	Declare,
	Try,
	Throw,
	Goto,
	Semicolon,
	CloseTag,
	OpenTagEcho,
	Text,
	OpenTag}

var classMemberDeclarationListRecoverSet = []TokenType{
	Public,
	Protected,
	Private,
	Static,
	Abstract,
	Final,
	Function,
	Var,
	Const,
	Use}

var encapsulatedVariableListRecoverSet = []TokenType{
	EncapsulatedAndWhitespace,
	DollarCurlyOpen,
	CurlyOpen}

type Associativity int

const (
	None Associativity = iota
	Left
	Right
)

func precedenceAssociativityTuple(t *Token) (int, Associativity) {
	switch t.Type {
	case AsteriskAsterisk:
		return 48, Right
	case PlusPlus:
		return 47, Right
	case MinusMinus:
		return 47, Right
	case Tilde:
		return 47, Right
	case IntegerCast:
		return 47, Right
	case FloatCast:
		return 47, Right
	case StringCast:
		return 47, Right
	case ArrayCast:
		return 47, Right
	case ObjectCast:
		return 47, Right
	case BooleanCast:
		return 47, Right
	case UnsetCast:
		return 47, Right
	case AtSymbol:
		return 47, Right
	case InstanceOf:
		return 46, None
	case Exclamation:
		return 45, Right
	case Asterisk:
		return 44, Left
	case ForwardSlash:
		return 44, Left
	case Percent:
		return 44, Left
	case Plus:
		return 43, Left
	case Minus:
		return 43, Left
	case Dot:
		return 43, Left
	case LessThanLessThan:
		return 42, Left
	case GreaterThanGreaterThan:
		return 42, Left
	case LessThan:
		return 41, None
	case GreaterThan:
		return 41, None
	case LessThanEquals:
		return 41, None
	case GreaterThanEquals:
		return 41, None
	case EqualsEquals:
		return 40, None
	case EqualsEqualsEquals:
		return 40, None
	case ExclamationEquals:
		return 40, None
	case ExclamationEqualsEquals:
		return 40, None
	case Spaceship:
		return 40, None
	case Ampersand:
		return 39, Left
	case Caret:
		return 38, Left
	case Bar:
		return 37, Left
	case AmpersandAmpersand:
		return 36, Left
	case BarBar:
		return 35, Left
	case QuestionQuestion:
		return 34, Right
	case Question:
		return 33, Left //?: ternary
	case Equals:
		return 32, Right
	case DotEquals:
		return 32, Right
	case PlusEquals:
		return 32, Right
	case MinusEquals:
		return 32, Right
	case AsteriskEquals:
		return 32, Right
	case ForwardslashEquals:
		return 32, Right
	case PercentEquals:
		return 32, Right
	case AsteriskAsteriskEquals:
		return 32, Right
	case AmpersandEquals:
		return 32, Right
	case BarEquals:
		return 32, Right
	case CaretEquals:
		return 32, Right
	case LessThanLessThanEquals:
		return 32, Right
	case GreaterThanGreaterThanEquals:
		return 32, Right
	case And:
		return 31, Left
	case Xor:
		return 30, Left
	case Or:
		return 29, Left
	}

	return 0, None
}

func binaryOpToPhraseType(t *Token) PhraseType {
	switch t.Type {
	case Question:
		return TernaryExpression
	case Dot, Plus, Minus:
		return AdditiveExpression
	case Bar, Ampersand, Caret:
		return BitwiseExpression
	case Asterisk, ForwardSlash, Percent:
		return MultiplicativeExpression
	case AsteriskAsterisk:
		return ExponentiationExpression
	case LessThanLessThan, GreaterThanGreaterThan:
		return ShiftExpression
	case AmpersandAmpersand, BarBar, And, Or, Xor:
		return LogicalExpression
	case EqualsEqualsEquals,
		ExclamationEqualsEquals,
		EqualsEquals,
		ExclamationEquals:
		return EqualityExpression
	case LessThan,
		LessThanEquals,
		GreaterThan,
		GreaterThanEquals,
		Spaceship:
		return RelationalExpression
	case QuestionQuestion:
		return CoalesceExpression
	case Equals:
		return SimpleAssignmentExpression
	case PlusEquals,
		MinusEquals,
		AsteriskEquals,
		AsteriskAsteriskEquals,
		ForwardslashEquals,
		DotEquals,
		PercentEquals,
		AmpersandEquals,
		BarEquals,
		CaretEquals,
		LessThanLessThanEquals,
		GreaterThanGreaterThanEquals:
		return CompoundAssignmentExpression
	case InstanceOf:
		return InstanceOfExpression
	default:
		return PhraseUnknown
	}
}

type Parser struct {
	tokenOffset     int
	tokenBuffer     []*Token
	phraseStack     []*Phrase
	errorPhrase     *ParseError
	recoverSetStack [][]TokenType
}

func tokenTypeIndexOf(haystack []TokenType, needle TokenType) int {
	for index, tokenType := range haystack {
		if tokenType == needle {
			return index
		}
	}

	return -1
}

func Parse(text string) *Phrase {
	doc := &Parser{0,
		Lex(text),
		make([]*Phrase, 0),
		nil,
		make([][]TokenType, 0)}
	stmtList := doc.statementList([]TokenType{EndOfFile})
	//append trailing hidden tokens
	doc.hidden(stmtList)

	return stmtList
}

func (doc *Parser) popRecover() []TokenType {
	var lastRecoverSet []TokenType

	lastRecoverSet, doc.recoverSetStack = doc.recoverSetStack[len(doc.recoverSetStack)-1],
		doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return lastRecoverSet
}

func (doc *Parser) start(phraseType PhraseType, dontPushHiddenToParent bool) *Phrase {
	//parent node gets hidden tokens between children
	if !dontPushHiddenToParent {
		doc.hidden(nil)
	}

	p := NewPhrase(phraseType, make([]AstNode, 0))

	doc.phraseStack = append(doc.phraseStack, p)

	return p
}

func (doc *Parser) end() *Phrase {
	var result *Phrase

	result, doc.phraseStack =
		doc.phraseStack[len(doc.phraseStack)-1], doc.phraseStack[:len(doc.phraseStack)-1]

	return result
}

func (doc *Parser) hidden(p *Phrase) {
	if p == nil {
		if len(doc.phraseStack) > 0 {
			p = doc.phraseStack[len(doc.phraseStack)-1]
		}
	}

	var t *Token

	for {
		if doc.tokenOffset < len(doc.tokenBuffer) {
			t = doc.tokenBuffer[doc.tokenOffset]
			doc.tokenOffset++
		}

		if t.Type < Comment {
			doc.tokenOffset--
			break
		} else {
			p.Children = append(p.Children, t)
		}
	}
}

func (doc *Parser) optional(tokenType TokenType) *Token {
	if tokenType == doc.peek(0).Type {
		doc.errorPhrase = nil

		return doc.next(false)
	}

	return nil
}

func (doc *Parser) optionalOneOf(tokenTypes []TokenType) *Token {
	if tokenTypeIndexOf(tokenTypes, doc.peek(0).Type) >= 0 {
		doc.errorPhrase = nil

		return doc.next(false)
	}

	return nil
}

func (doc *Parser) next(doNotPush bool) *Token {
	t := doc.tokenBuffer[doc.tokenOffset]
	doc.tokenOffset++

	if t.Type == EndOfFile {
		return t
	}

	lastPhrase := doc.phraseStack[len(doc.phraseStack)-1]
	if t.Type >= Comment {
		//hidden token
		lastPhrase.Children = append(lastPhrase.Children, t)

		return doc.next(doNotPush)
	} else if !doNotPush {
		lastPhrase.Children = append(lastPhrase.Children, t)
	}

	return t
}

func (doc *Parser) expect(tokenType TokenType) *Token {
	t := doc.peek(0)

	if t.Type == tokenType {
		doc.errorPhrase = nil

		return doc.next(false)
	} else if tokenType == Semicolon && t.Type == CloseTag {
		//implicit end statement
		return t
	} else {
		doc.error(tokenType)
		//test skipping a single token to sync
		if doc.peek(1).Type == tokenType {
			doc.skip(func(x *Token) bool {
				return x.Type == tokenType
			})
			doc.errorPhrase = nil
			return doc.next(false) //tokenType
		}
	}

	return nil
}

func (doc *Parser) expectOneOf(tokenTypes []TokenType) *Token {
	t := doc.peek(0)

	if tokenTypeIndexOf(tokenTypes, t.Type) >= 0 {
		doc.errorPhrase = nil

		return doc.next(false)
	} else if tokenTypeIndexOf(tokenTypes, Semicolon) >= 0 && t.Type == CloseTag {
		//implicit end statement
		return t
	}

	doc.error(Undefined)
	//test skipping single token to sync
	if tokenTypeIndexOf(tokenTypes, doc.peek(1).Type) >= 0 {
		doc.skip(func(x *Token) bool {
			return tokenTypeIndexOf(tokenTypes, x.Type) >= 0
		})
		doc.errorPhrase = nil

		return doc.next(false) //tokenType
	}

	return nil
}

func (doc *Parser) peek(n int) *Token {
	k := n + 1
	bufferPos := doc.tokenOffset - 1
	var t *Token

	for {
		bufferPos++
		t = doc.tokenBuffer[bufferPos]

		if t.Type < Comment {
			//not a hidden token
			k--
		}

		if t.Type == EndOfFile || k == 0 {
			break
		}

	}

	return t
}

/**
* skipped tokens get pushed to error phrase children
 */
func (doc *Parser) skip(predicate func(*Token) bool) {

	var t *Token

	for {
		if doc.tokenOffset < len(doc.tokenBuffer) {
			t = doc.tokenBuffer[doc.tokenOffset]
			doc.tokenOffset++
		}

		if predicate(t) || t.Type == EndOfFile {
			doc.tokenOffset--
			break
		} else {
			doc.errorPhrase.Children = append(doc.errorPhrase.Children, t)
		}
	}

}

func (doc *Parser) error(expected TokenType) {

	//dont report errors if recovering from another
	if doc.errorPhrase != nil {
		return
	}

	doc.errorPhrase = NewParseErr(
		Error, make([]AstNode, 0), doc.peek(0), expected)

	lastPhrase := doc.phraseStack[len(doc.phraseStack)-1]

	lastPhrase.Children = append(lastPhrase.Children, doc.errorPhrase)
}

func (doc *Parser) list(phraseType PhraseType,
	elementFunction func() AstNode,
	elementStartPredicate func(*Token) bool,
	breakOn []TokenType,
	recoverSet []TokenType) *Phrase {

	p := doc.start(phraseType, false)

	var t *Token
	recoveryAttempted := false
	var listRecoverSet []TokenType

	if recoverSet != nil {
		listRecoverSet = recoverSet
	} else {
		listRecoverSet = make([]TokenType, 0)
	}

	if breakOn != nil {
		listRecoverSet = append(listRecoverSet, breakOn...)
	}

	doc.recoverSetStack = append(doc.recoverSetStack, listRecoverSet)

	for {
		t = doc.peek(0)

		if elementStartPredicate(t) {
			recoveryAttempted = false
			p.Children = append(p.Children, elementFunction())
		} else if breakOn == nil ||
			tokenTypeIndexOf(breakOn, t.Type) >= 0 ||
			recoveryAttempted {
			break
		} else {
			doc.error(Undefined)
			//attempt to sync with token stream
			t = doc.peek(1)
			if elementStartPredicate(t) || tokenTypeIndexOf(breakOn, t.Type) >= 0 {
				doc.skip(func(x *Token) bool {
					return reflect.DeepEqual(x, t)
				})
			} else {
				doc.defaultSyncStrategy()
			}
			recoveryAttempted = true
		}

	}

	doc.recoverSetStack = doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return doc.end()
}

func (doc *Parser) defaultSyncStrategy() {

	mergedRecoverTokenTypeArray := make([]TokenType, 0)

	for n := len(doc.recoverSetStack) - 1; n >= 0; n-- {
		mergedRecoverTokenTypeArray = append(mergedRecoverTokenTypeArray, doc.recoverSetStack[n]...)
	}

	mergedRecoverTokenTypeSet := make(map[TokenType]bool, len(mergedRecoverTokenTypeArray))

	for _, recoverTokenType := range mergedRecoverTokenTypeArray {
		if _, ok := mergedRecoverTokenTypeSet[recoverTokenType]; !ok {
			mergedRecoverTokenTypeSet[recoverTokenType] = true
		}
	}

	doc.skip(func(x *Token) bool {
		_, ok := mergedRecoverTokenTypeSet[x.Type]

		return ok
	})
}

func (doc *Parser) statementList(breakOn []TokenType) *Phrase {
	return doc.list(
		StatementList,
		doc.statement,
		isStatementStart,
		breakOn,
		statementListRecoverSet[:])
}

func (doc *Parser) constDeclaration() *Phrase {
	p := doc.start(ConstDeclaration, false)
	doc.next(false) //const

	p.Children = append(p.Children, doc.delimitedList(
		ConstElementList,
		doc.constElement,
		isConstElementStartToken,
		Comma,
		[]TokenType{Semicolon},
		false))

	doc.expect(Semicolon)

	return doc.end()
}

func isClassConstElementStartToken(t *Token) bool {
	return t.Type == Name || isSemiReservedToken(t)
}

func isConstElementStartToken(t *Token) bool {
	return t.Type == Name
}

func (doc *Parser) constElement() AstNode {
	p := doc.start(ConstElement, false)

	doc.expect(Name)
	doc.expect(Equals)

	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) expression(minPrecedence int) AstNode {
	var precedence int
	var associativity Associativity
	var op *Token
	var p *Phrase
	var binaryPhraseType PhraseType

	lhs := doc.expressionAtom(minPrecedence)

	for {
		op = doc.peek(0)
		binaryPhraseType = binaryOpToPhraseType(op)

		if binaryPhraseType == PhraseUnknown {
			break
		}

		precedence, associativity = precedenceAssociativityTuple(op)

		if precedence < minPrecedence {
			break
		}

		if associativity == Left {
			precedence++
		}

		if binaryPhraseType == TernaryExpression {
			lhs = doc.ternaryExpression(lhs)

			continue
		}

		p = doc.start(binaryPhraseType, true)
		p.Children = append(p.Children, lhs)
		doc.next(false)

		if binaryPhraseType == InstanceOfExpression {
			p.Children = append(p.Children, doc.typeDesignator(InstanceofTypeDesignator))
		} else {
			if binaryPhraseType == SimpleAssignmentExpression &&
				doc.peek(0).Type == Ampersand {
				doc.next(false) //&
				p.Type = ByRefAssignmentExpression
			}

			p.Children = append(p.Children, doc.expression(precedence))
		}

		lhs = doc.end()
	}

	return lhs
}

func (doc *Parser) ternaryExpression(testExpr AstNode) *Phrase {
	p := doc.start(TernaryExpression, true)
	p.Children = append(p.Children, testExpr)

	doc.next(false) //?

	if colonToken := doc.optional(Colon); colonToken != nil {
		p.Children = append(p.Children, doc.expression(0))
	} else {
		p.Children = append(p.Children, doc.expression(0))
		doc.expect(Colon)
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) variableOrExpression(precedence int) AstNode {
	part := doc.variableAtom(precedence)
	isVariable := false

	if p, ok := part.(*Phrase); ok {
		isVariable = p.Type == SimpleVariable
	}

	if isDereferenceOperator(doc.peek(0)) {
		part = doc.variable(part)
		isVariable = true
	} else {
		switch part.(*Phrase).Type {
		case QualifiedName,
			FullyQualifiedName,
			RelativeQualifiedName:
			part = doc.constantAccessExpression(part)
		}
	}

	if !isVariable {
		return part
	}

	//check for post increment/decrement
	t := doc.peek(0)
	if t.Type == PlusPlus {
		return doc.postfixExpression(PostfixIncrementExpression, part.(*Phrase))
	} else if t.Type == MinusMinus {
		return doc.postfixExpression(PostfixDecrementExpression, part.(*Phrase))
	} else {
		return part
	}
}

func (doc *Parser) constantAccessExpression(qName AstNode) *Phrase {
	p := doc.start(ConstantAccessExpression, true)
	p.Children = append(p.Children, qName)

	return doc.end()
}

func (doc *Parser) postfixExpression(
	phraseType PhraseType, variableNode *Phrase) *Phrase {
	p := doc.start(phraseType, true)
	p.Children = append(p.Children, variableNode)
	doc.next(false) //operator

	return doc.end()
}

func isDereferenceOperator(t *Token) bool {
	switch t.Type {
	case OpenBracket,
		OpenBrace,
		Arrow,
		OpenParenthesis,
		ColonColon:
		return true
	}

	return false
}

func (doc *Parser) expressionAtom(precedence int) AstNode {
	t := doc.peek(0)

	switch t.Type {
	case Static:
		if doc.peek(1).Type == Function {
			return doc.anonymousFunctionCreationExpression()
		}

		return doc.variableOrExpression(0)
	case StringLiteral:
		if isDereferenceOperator(doc.peek(1)) {
			return doc.variableOrExpression(0)
		}

		return doc.next(true)
	case VariableName,
		Dollar,
		Array,
		OpenBracket,
		Backslash,
		Name,
		Namespace,
		OpenParenthesis:
		return doc.variableOrExpression(precedence)
	case PlusPlus:
		return doc.unaryExpression(PrefixIncrementExpression)
	case MinusMinus:
		return doc.unaryExpression(PrefixDecrementExpression)
	case Plus,
		Minus,
		Exclamation,
		Tilde:
		return doc.unaryExpression(UnaryOpExpression)
	case AtSymbol:
		return doc.unaryExpression(ErrorControlExpression)
	case IntegerCast,
		FloatCast,
		StringCast,
		ArrayCast,
		ObjectCast,
		BooleanCast,
		UnsetCast:
		return doc.unaryExpression(CastExpression)
	case List:
		return doc.listIntrinsic()
	case Clone:
		return doc.cloneExpression()
	case New:
		return doc.objectCreationExpression()
	case FloatingLiteral,
		IntegerLiteral,
		LineConstant,
		FileConstant,
		DirectoryConstant,
		TraitConstant,
		MethodConstant,
		FunctionConstant,
		NamespaceConstant,
		ClassConstant:
		return doc.next(true)
	case StartHeredoc:
		return doc.heredocStringLiteral()
	case DoubleQuote:
		return doc.doubleQuotedStringLiteral()
	case Backtick:
		return doc.shellCommandExpression()
	case Print:
		return doc.printIntrinsic()
	case Yield:
		return doc.yieldExpression()
	case YieldFrom:
		return doc.yieldFromExpression()
	case Function:
		return doc.anonymousFunctionCreationExpression()
	case Include:
		return doc.scriptInclusion(IncludeExpression)
	case IncludeOnce:
		return doc.scriptInclusion(IncludeOnceExpression)
	case Require:
		return doc.scriptInclusion(RequireExpression)
	case RequireOnce:
		return doc.scriptInclusion(RequireOnceExpression)
	case Eval:
		return doc.evalIntrinsic()
	case Empty:
		return doc.emptyIntrinsic()
	case Exit:
		return doc.exitIntrinsic()
	case Isset:
		return doc.issetIntrinsic()
	default:
		//error
		doc.start(ErrorExpression, false)
		doc.error(Undefined)

		return doc.end()
	}
}

func (doc *Parser) exitIntrinsic() *Phrase {
	p := doc.start(ExitIntrinsic, false)
	doc.next(false) //exit or die
	if t := doc.optional(OpenParenthesis); t != nil {
		if isExpressionStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.expression(0))
		}

		doc.expect(CloseParenthesis)
	}

	return doc.end()
}

func (doc *Parser) issetIntrinsic() *Phrase {
	p := doc.start(IssetIntrinsic, false)
	doc.next(false) //isset
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.variableList([]TokenType{CloseParenthesis}))
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) emptyIntrinsic() *Phrase {
	p := doc.start(EmptyIntrinsic, false)
	doc.next(false) //keyword
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) evalIntrinsic() *Phrase {
	p := doc.start(EvalIntrinsic, false)
	doc.next(false) //keyword
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) scriptInclusion(phraseType PhraseType) *Phrase {
	p := doc.start(phraseType, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) printIntrinsic() *Phrase {
	p := doc.start(PrintIntrinsic, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) yieldFromExpression() *Phrase {
	p := doc.start(YieldFromExpression, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) yieldExpression() *Phrase {
	p := doc.start(YieldExpression, false)
	doc.next(false) //yield
	if !isExpressionStart(doc.peek(0)) {
		return doc.end()
	}

	keyOrValue := doc.expression(0)
	p.Children = append(p.Children, keyOrValue)

	if doc.optional(FatArrow) != nil {
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) shellCommandExpression() *Phrase {
	p := doc.start(ShellCommandExpression, false)
	doc.next(false) //`
	p.Children = append(p.Children, doc.encapsulatedVariableList(Backtick))
	doc.expect(Backtick)

	return doc.end()
}

func (doc *Parser) doubleQuotedStringLiteral() *Phrase {
	p := doc.start(DoubleQuotedStringLiteral, false)
	doc.next(false) //"
	p.Children = append(p.Children, doc.encapsulatedVariableList(DoubleQuote))
	doc.expect(DoubleQuote)

	return doc.end()
}

func (doc *Parser) encapsulatedVariableList(breakOn TokenType) *Phrase {
	return doc.list(
		EncapsulatedVariableList,
		doc.encapsulatedVariable,
		isEncapsulatedVariableStart,
		[]TokenType{breakOn},
		encapsulatedVariableListRecoverSet)
}

func isEncapsulatedVariableStart(t *Token) bool {

	switch t.Type {
	case EncapsulatedAndWhitespace,
		VariableName,
		DollarCurlyOpen,
		CurlyOpen:
		return true
	}

	return false
}

func (doc *Parser) encapsulatedVariable() AstNode {

	switch doc.peek(0).Type {
	case EncapsulatedAndWhitespace:
		return doc.next(true)
	case VariableName:
		t := doc.peek(1)
		if t.Type == OpenBracket {
			return doc.encapsulatedDimension()
		} else if t.Type == Arrow {
			return doc.encapsulatedProperty()
		}

		return doc.simpleVariable()
	case DollarCurlyOpen:
		return doc.dollarCurlyOpenEncapsulatedVariable()
	case CurlyOpen:
		return doc.curlyOpenEncapsulatedVariable()
	}

	t := doc.peek(0)

	panic(errors.New("Unexpected token: " + t.Type.String()))
}

func (doc *Parser) curlyOpenEncapsulatedVariable() *Phrase {
	p := doc.start(EncapsulatedVariable, false)
	doc.next(false) //{
	p.Children = append(p.Children, doc.variable(doc.variableAtom(0)))
	doc.expect(CloseBrace)

	return doc.end()
}

func (doc *Parser) dollarCurlyOpenEncapsulatedVariable() *Phrase {
	p := doc.start(EncapsulatedVariable, false)
	doc.next(false) //${
	t := doc.peek(0)

	if t.Type == VariableName {
		if doc.peek(1).Type == OpenBracket {
			p.Children = append(p.Children, *doc.dollarCurlyEncapsulatedDimension())
		} else {
			doc.start(SimpleVariable, false)
			doc.next(false)
			p.Children = append(p.Children, doc.end())
		}
	} else if isExpressionStart(t) {
		p.Children = append(p.Children, doc.expression(0))
	} else {
		doc.error(Undefined)
	}

	doc.expect(CloseBrace)

	return doc.end()
}

func (doc *Parser) dollarCurlyEncapsulatedDimension() *Phrase {
	p := doc.start(SubscriptExpression, false)
	doc.next(false) //VariableName
	doc.next(false) // [
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseBracket)

	return doc.end()
}

func (doc *Parser) encapsulatedDimension() *Phrase {
	p := doc.start(SubscriptExpression, false)

	p.Children = append(p.Children, doc.simpleVariable()) //T_VARIABLE
	doc.next(false)                                       //[
	switch doc.peek(0).Type {
	case Name,
		IntegerLiteral:
		doc.next(false)
	case VariableName:
		p.Children = append(p.Children, doc.simpleVariable())
	case Minus:
		doc.start(UnaryOpExpression, false)
		doc.next(false) //-
		doc.expect(IntegerLiteral)
		p.Children = append(p.Children, doc.end())
	default:
		//error
		doc.error(Undefined)
	}

	doc.expect(CloseBracket)
	return doc.end()
}

func (doc *Parser) encapsulatedProperty() *Phrase {
	p := doc.start(PropertyAccessExpression, false)
	p.Children = append(p.Children, doc.simpleVariable())
	doc.next(false) //->
	doc.expect(Name)

	return doc.end()
}

func (doc *Parser) heredocStringLiteral() *Phrase {
	p := doc.start(HeredocStringLiteral, false)
	doc.next(false) //StartHeredoc
	p.Children = append(p.Children, doc.encapsulatedVariableList(EndHeredoc))
	doc.expect(EndHeredoc)

	return doc.end()
}

func (doc *Parser) anonymousClassDeclaration() *Phrase {
	p := doc.start(AnonymousClassDeclaration, false)
	p.Children = append(p.Children, doc.anonymousClassDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(ClassDeclarationBody,
		isClassMemberStart, doc.classMemberDeclarationList))

	return doc.end()
}

func (doc *Parser) anonymousClassDeclarationHeader() *Phrase {
	p := doc.start(AnonymousClassDeclarationHeader, false)
	doc.next(false) //class
	if doc.optional(OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}
		doc.expect(CloseParenthesis)
	}

	if doc.peek(0).Type == Extends {
		p.Children = append(p.Children, doc.classBaseClause())
	}

	if doc.peek(0).Type == Implements {
		p.Children = append(p.Children, doc.classInterfaceClause())
	}

	return doc.end()
}

func (doc *Parser) classInterfaceClause() *Phrase {
	p := doc.start(ClassInterfaceClause, false)
	doc.next(false) //implements
	p.Children = append(p.Children, doc.qualifiedNameList([]TokenType{OpenBrace}))

	return doc.end()
}

func (doc *Parser) classMemberDeclarationList() *Phrase {
	return doc.list(
		ClassMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]TokenType{CloseBrace},
		classMemberDeclarationListRecoverSet)
}

func isClassMemberStart(t *Token) bool {
	switch t.Type {
	case Public,
		Protected,
		Private,
		Static,
		Abstract,
		Final,
		Function,
		Var,
		Const,
		Use:
		return true
	}

	return false
}

func (doc *Parser) classMemberDeclaration() AstNode {
	p := doc.start(ErrorClassMemberDeclaration, false)
	t := doc.peek(0)

	switch t.Type {
	case Public,
		Protected,
		Private,
		Static,
		Abstract,
		Final:
		modifiers := doc.memberModifierList()
		t = doc.peek(0)
		if t.Type == VariableName {
			p.Children = append(p.Children, modifiers)

			return doc.propertyDeclaration(p)
		} else if t.Type == Function {
			return doc.methodDeclaration(p, modifiers)
		} else if t.Type == Const {
			p.Children = append(p.Children, modifiers)

			return doc.classConstDeclaration(p)
		}

		//error
		p.Children = append(p.Children, modifiers)
		doc.error(Undefined)

		return doc.end()
	case Function:
		return doc.methodDeclaration(p, nil)
	case Var:
		doc.next(false)

		return doc.propertyDeclaration(p)
	case Const:
		return doc.classConstDeclaration(p)
	case Use:
		return doc.traitUseClause(p)
	}

	panic(errors.New("Unexpected token: " + t.Type.String()))
}

func (doc *Parser) traitUseClause(p *Phrase) *Phrase {
	p.Type = TraitUseClause
	doc.next(false) //use
	p.Children = append(p.Children,
		doc.qualifiedNameList([]TokenType{Semicolon, OpenBrace}),
		doc.traitUseSpecification())

	return doc.end()

}

func (doc *Parser) traitUseSpecification() *Phrase {
	p := doc.start(TraitUseSpecification, false)
	t := doc.expectOneOf([]TokenType{Semicolon, OpenBrace})

	if t != nil && t.Type == OpenBrace {
		if isTraitAdaptationStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.traitAdaptationList())
		}
		doc.expect(CloseBrace)
	}

	return doc.end()
}

func (doc *Parser) traitAdaptationList() *Phrase {
	return doc.list(
		TraitAdaptationList,
		doc.traitAdaptation,
		isTraitAdaptationStart,
		[]TokenType{CloseBrace}, nil)
}

func isTraitAdaptationStart(t *Token) bool {
	switch t.Type {
	case Name,
		Backslash,
		Namespace:
		return true
	}

	return isSemiReservedToken(t)
}

func (doc *Parser) traitAdaptation() AstNode {
	p := doc.start(ErrorTraitAdaptation, false)
	t := doc.peek(0)
	t2 := doc.peek(1)

	if t.Type == Namespace ||
		t.Type == Backslash ||
		(t.Type == Name &&
			(t2.Type == ColonColon || t2.Type == Backslash)) {

		p.Children = append(p.Children, doc.methodReference())

		if doc.peek(0).Type == InsteadOf {
			doc.next(false)
			return doc.traitPrecedence(p)
		}

	} else if t.Type == Name || isSemiReservedToken(t) {
		methodRef := doc.start(MethodReference, false)
		methodRef.Children = append(methodRef.Children, doc.identifier())
		p.Children = append(p.Children, doc.end())
	} else {
		//error
		doc.error(Undefined)

		return doc.end()
	}

	return doc.traitAlias(p)
}

func (doc *Parser) traitAlias(p *Phrase) *Phrase {
	p.Type = TraitAlias
	doc.expect(As)

	t := doc.peek(0)

	if t.Type == Name || isReservedToken(t) {
		p.Children = append(p.Children, doc.identifier())
	} else if isMemberModifier(t) {
		doc.next(false)
		t = doc.peek(0)
		if t.Type == Name || isSemiReservedToken(t) {
			p.Children = append(p.Children, doc.identifier())
		}
	} else {
		doc.error(Undefined)
	}

	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) traitPrecedence(p *Phrase) *Phrase {
	p.Type = TraitPrecedence
	p.Children = append(p.Children, doc.qualifiedNameList([]TokenType{Semicolon}))
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) methodReference() *Phrase {
	p := doc.start(MethodReference, false)
	p.Children = append(p.Children, doc.qualifiedName())
	doc.expect(ColonColon)
	p.Children = append(p.Children, doc.identifier())

	return doc.end()
}

func (doc *Parser) methodDeclarationHeader(memberModifers *Phrase) *Phrase {
	p := doc.start(MethodDeclarationHeader, true)
	if memberModifers != nil {
		p.Children = append(p.Children, memberModifers)
	}
	doc.next(false) //function
	doc.optional(Ampersand)
	p.Children = append(p.Children, doc.identifier())
	doc.expect(OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			Comma,
			[]TokenType{CloseParenthesis}, false))
	}

	doc.expect(CloseParenthesis)

	if doc.peek(0).Type == Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func (doc *Parser) methodDeclaration(p *Phrase, memberModifers *Phrase) *Phrase {
	p.Type = MethodDeclaration
	p.Children = append(p.Children, doc.methodDeclarationHeader(memberModifers))
	p.Children = append(p.Children, doc.methodDeclarationBody())

	return doc.end()
}

func (doc *Parser) methodDeclarationBody() *Phrase {
	p := doc.start(MethodDeclarationBody, false)

	if doc.peek(0).Type == Semicolon {
		doc.next(false)
	} else {
		p.Children = append(p.Children, doc.compoundStatement())
	}

	return doc.end()
}

func (doc *Parser) identifier() *Phrase {
	doc.start(Identifier, false)
	t := doc.peek(0)
	if t.Type == Name || isSemiReservedToken(t) {
		doc.next(false)
	} else {
		doc.error(Undefined)
	}

	return doc.end()
}

func (doc *Parser) interfaceDeclaration() *Phrase {
	p := doc.start(InterfaceDeclaration, false)
	p.Children = append(p.Children, doc.interfaceDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		InterfaceDeclarationBody, isClassMemberStart, doc.interfaceMemberDeclarations))

	return doc.end()
}

func (doc *Parser) typeDeclarationBody(
	phraseType PhraseType,
	elementStartPredicate func(*Token) bool,
	listFunction func() *Phrase) *Phrase {
	p := doc.start(phraseType, false)
	doc.expect(OpenBrace)

	if elementStartPredicate(doc.peek(0)) {
		p.Children = append(p.Children, listFunction())
	}

	doc.expect(CloseBrace)

	return doc.end()
}

func (doc *Parser) interfaceMemberDeclarations() *Phrase {
	return doc.list(
		InterfaceMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]TokenType{CloseBrace},
		classMemberDeclarationListRecoverSet)
}

func (doc *Parser) interfaceDeclarationHeader() *Phrase {
	p := doc.start(InterfaceDeclarationHeader, false)
	doc.next(false) //interface
	doc.expect(Name)

	if doc.peek(0).Type == Extends {
		p.Children = append(p.Children, doc.interfaceBaseClause())
	}

	return doc.end()

}

func (doc *Parser) interfaceBaseClause() *Phrase {
	p := doc.start(InterfaceBaseClause, false)
	doc.next(false) //extends
	p.Children = append(p.Children, doc.qualifiedNameList([]TokenType{OpenBrace}))

	return doc.end()
}

func (doc *Parser) traitDeclaration() *Phrase {
	p := doc.start(TraitDeclaration, false)
	p.Children = append(p.Children, doc.traitDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		TraitDeclarationBody, isClassMemberStart, doc.traitMemberDeclarations))

	return doc.end()
}

func (doc *Parser) traitDeclarationHeader() *Phrase {
	doc.start(TraitDeclarationHeader, false)
	doc.next(false) //trait
	doc.expect(Name)

	return doc.end()
}

func (doc *Parser) traitMemberDeclarations() *Phrase {
	return doc.list(
		TraitMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]TokenType{CloseBrace},
		classMemberDeclarationListRecoverSet[:])
}

func (doc *Parser) functionDeclaration() *Phrase {
	p := doc.start(FunctionDeclaration, false)
	p.Children = append(p.Children, doc.functionDeclarationHeader())
	p.Children = append(p.Children, doc.functionDeclarationBody())

	return doc.end()
}

func (doc *Parser) functionDeclarationBody() *Phrase {
	cs := doc.compoundStatement()
	cs.Type = FunctionDeclarationBody

	return cs
}

func (doc *Parser) functionDeclarationHeader() *Phrase {
	p := doc.start(FunctionDeclarationHeader, false)

	doc.next(false) //function
	doc.optional(Ampersand)
	doc.expect(Name)
	doc.expect(OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			Comma,
			[]TokenType{CloseParenthesis}, false))
	}

	doc.expect(CloseParenthesis)

	if doc.peek(0).Type == Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func isParameterStart(t *Token) bool {
	switch t.Type {
	case Ampersand,
		Ellipsis,
		VariableName:
		return true
	default:
		return isTypeDeclarationStart(t)
	}
}

func (doc *Parser) classDeclaration() *Phrase {
	p := doc.start(ClassDeclaration, false)

	p.Children = append(p.Children, doc.classDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		ClassDeclarationBody, isClassMemberStart, doc.classMemberDeclarationList))

	return doc.end()

}

func (doc *Parser) classDeclarationHeader() *Phrase {
	p := doc.start(ClassDeclarationHeader, false)
	doc.optionalOneOf([]TokenType{Abstract, Final})
	doc.expect(Class)
	doc.expect(Name)

	if doc.peek(0).Type == Extends {
		p.Children = append(p.Children, doc.classBaseClause())
	}

	if doc.peek(0).Type == Implements {
		p.Children = append(p.Children, doc.classInterfaceClause())
	}

	return doc.end()
}

func (doc *Parser) classBaseClause() *Phrase {
	p := doc.start(ClassBaseClause, false)
	doc.next(false) //extends
	p.Children = append(p.Children, doc.qualifiedName())

	return doc.end()
}

func (doc *Parser) compoundStatement() *Phrase {
	p := doc.start(CompoundStatement, false)
	doc.expect(OpenBrace)

	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList([]TokenType{CloseBrace}))
	}

	doc.expect(CloseBrace)

	return doc.end()
}

func (doc *Parser) statement() AstNode {
	t := doc.peek(0)

	switch t.Type {
	case Namespace:
		return doc.namespaceDefinition()
	case Use:
		return doc.namespaceUseDeclaration()
	case HaltCompiler:
		return doc.haltCompilerStatement()
	case Const:
		return doc.constDeclaration()
	case Function:
		{
			p1 := doc.peek(1)
			if p1.Type == OpenParenthesis ||
				(p1.Type == Ampersand && doc.peek(2).Type == OpenParenthesis) {
				//anon fn without assignment
				return doc.expressionStatement()
			} else {
				return doc.functionDeclaration()
			}
		}
	case Class,
		Abstract,
		Final:
		return doc.classDeclaration()
	case Trait:
		return doc.traitDeclaration()
	case Interface:
		return doc.interfaceDeclaration()
	case OpenBrace:
		return doc.compoundStatement()
	case If:
		return doc.ifStatement()
	case While:
		return doc.whileStatement()
	case Do:
		return doc.doStatement()
	case For:
		return doc.forStatement()
	case Switch:
		return doc.switchStatement()
	case Break:
		return doc.breakStatement()
	case Continue:
		return doc.continueStatement()
	case Return:
		return doc.returnStatement()
	case Global:
		return doc.globalDeclaration()
	case Static:
		if doc.peek(1).Type == VariableName &&
			tokenTypeIndexOf(
				[]TokenType{Semicolon, Comma, CloseTag, Equals},
				doc.peek(2).Type) >= 0 {
			return doc.functionStaticDeclaration()
		} else {
			return doc.expressionStatement()
		}
	case Text,
		OpenTag,
		CloseTag:
		return doc.inlineText()
	case ForEach:
		return doc.foreachStatement()
	case Declare:
		return doc.declareStatement()
	case Try:
		return doc.tryStatement()
	case Throw:
		return doc.throwStatement()
	case Goto:
		return doc.gotoStatement()
	case Echo,
		OpenTagEcho:
		return doc.echoIntrinsic()
	case Unset:
		return doc.unsetIntrinsic()
	case Semicolon:
		return doc.nullStatement()
	case Name:
		if doc.peek(1).Type == Colon {
			return doc.namedLabelStatement()
		}
		fallthrough
	default:
		return doc.expressionStatement()
	}
}

func (doc *Parser) inlineText() *Phrase {
	doc.start(InlineText, false)

	doc.optional(CloseTag)
	doc.optional(Text)
	doc.optional(OpenTag)

	return doc.end()
}

func (doc *Parser) nullStatement() *Phrase {
	doc.start(NullStatement, false)
	doc.next(false) //;

	return doc.end()
}

func (doc *Parser) tryStatement() *Phrase {
	p := doc.start(TryStatement, false)
	doc.next(false) //try
	p.Children = append(p.Children, doc.compoundStatement())

	t := doc.peek(0)

	if t.Type == Catch {
		p.Children = append(p.Children, doc.list(
			CatchClauseList,
			doc.catchClause,
			func(t *Token) bool { return t.Type == Catch }, nil, nil))
	} else if t.Type != Finally {
		doc.error(Undefined)
	}

	if doc.peek(0).Type == Finally {
		p.Children = append(p.Children, doc.finallyClause())
	}

	return doc.end()

}

func (doc *Parser) finallyClause() *Phrase {
	p := doc.start(FinallyClause, false)
	doc.next(false) //finally
	p.Children = append(p.Children, doc.compoundStatement())

	return doc.end()
}

func (doc *Parser) catchClause() AstNode {
	p := doc.start(CatchClause, false)
	doc.next(false) //catch
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.delimitedList(
		CatchNameList,
		doc.qualifiedName,
		isQualifiedNameStart,
		Bar,
		[]TokenType{VariableName}, false))
	doc.expect(VariableName)
	doc.expect(CloseParenthesis)
	p.Children = append(p.Children, doc.compoundStatement())

	return doc.end()
}

func (doc *Parser) declareDirective() *Phrase {
	doc.start(DeclareDirective, false)
	doc.expect(Name)
	doc.expect(Equals)
	doc.expectOneOf([]TokenType{
		IntegerLiteral, FloatingLiteral, StringLiteral})

	return doc.end()
}

func (doc *Parser) declareStatement() *Phrase {
	p := doc.start(DeclareStatement, false)
	doc.next(false) //declare
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.declareDirective())
	doc.expect(CloseParenthesis)

	t := doc.peek(0)

	if t.Type == Colon {
		doc.next(false) //:
		p.Children = append(p.Children, doc.statementList([]TokenType{EndDeclare}))
		doc.expect(EndDeclare)
		doc.expect(Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else if t.Type == Semicolon {
		doc.next(false)
	} else {
		doc.error(Undefined)
	}

	return doc.end()
}

func (doc *Parser) switchStatement() *Phrase {
	p := doc.start(SwitchStatement, false)
	doc.next(false) //switch
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)

	t := doc.expectOneOf([]TokenType{Colon, OpenBrace})
	tCase := doc.peek(0)

	if tCase.Type == Case || tCase.Type == Default {
		caseTokenType := CloseBrace
		if t != nil && t.Type == Colon {
			caseTokenType = EndSwitch
		}

		p.Children = append(p.Children, doc.caseStatements(caseTokenType))
	}

	if t != nil && t.Type == Colon {
		doc.expect(EndSwitch)
		doc.expect(Semicolon)
	} else {
		doc.expect(CloseBrace)
	}

	return doc.end()

}

func (doc *Parser) caseStatements(breakOn TokenType) *Phrase {
	p := doc.start(CaseStatementList, false)
	var t *Token
	caseBreakOn := []TokenType{Case, Default, breakOn}

	for {
		t = doc.peek(0)

		if t.Type == Case {
			p.Children = append(p.Children, doc.caseStatement(caseBreakOn))
		} else if t.Type == Default {
			p.Children = append(p.Children, doc.defaultStatement(caseBreakOn))
		} else if breakOn == t.Type {
			break
		} else {
			doc.error(Undefined)
			break
		}

	}

	return doc.end()
}

func (doc *Parser) caseStatement(breakOn []TokenType) *Phrase {
	p := doc.start(CaseStatement, false)
	doc.next(false) //case
	p.Children = append(p.Children, doc.expression(0))
	doc.expectOneOf([]TokenType{Colon, Semicolon})
	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList(breakOn))
	}

	return doc.end()
}

func (doc *Parser) defaultStatement(breakOn []TokenType) *Phrase {
	p := doc.start(DefaultStatement, false)
	doc.next(false) //default
	doc.expectOneOf([]TokenType{Colon, Semicolon})
	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList(breakOn))
	}

	return doc.end()
}

func (doc *Parser) namedLabelStatement() *Phrase {
	doc.start(NamedLabelStatement, false)
	doc.next(false) //name
	doc.next(false) //:

	return doc.end()
}

func (doc *Parser) gotoStatement() *Phrase {
	doc.start(GotoStatement, false)
	doc.next(false) //goto
	doc.expect(Name)
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) throwStatement() *Phrase {
	p := doc.start(ThrowStatement, false)
	doc.next(false) //throw
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) foreachCollection() *Phrase {
	p := doc.start(ForeachCollection, false)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) foreachKeyOrValue() *Phrase {
	p := doc.start(ForeachValue, false)
	p.Children = append(p.Children, doc.expression(0))
	if doc.peek(0).Type == FatArrow {
		doc.next(false)
		p.Type = ForeachKey
	}

	return doc.end()
}

func (doc *Parser) foreachValue() *Phrase {
	p := doc.start(ForeachValue, false)
	doc.optional(Ampersand)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) foreachStatement() *Phrase {
	p := doc.start(ForeachStatement, false)
	doc.next(false) //foreach
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.foreachCollection())
	doc.expect(As)
	var keyOrValue *Phrase
	if doc.peek(0).Type == Ampersand {
		keyOrValue = doc.foreachValue()
	} else {
		keyOrValue = doc.foreachKeyOrValue()
	}
	p.Children = append(p.Children, keyOrValue)

	if keyOrValue.Type == ForeachKey {
		p.Children = append(p.Children, doc.foreachValue())
	}

	doc.expect(CloseParenthesis)

	t := doc.peek(0)

	if t.Type == Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]TokenType{EndForeach}))
		doc.expect(EndForeach)
		doc.expect(Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(Undefined)
	}

	return doc.end()

}

func isVariableStart(t *Token) bool {
	switch t.Type {
	case VariableName,
		Dollar,
		OpenParenthesis,
		Array,
		OpenBracket,
		StringLiteral,
		Static,
		Name,
		Namespace,
		Backslash:
		return true
	}

	return false
}

func (doc *Parser) variableInitial() AstNode {
	return doc.variable(doc.variableAtom(0))
}

func (doc *Parser) variableList(breakOn []TokenType) *Phrase {
	return doc.delimitedList(
		VariableList,
		doc.variableInitial,
		isVariableStart,
		Comma,
		breakOn,
		false)
}

func (doc *Parser) unsetIntrinsic() *Phrase {
	p := doc.start(UnsetIntrinsic, false)
	doc.next(false) //unset
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.variableList([]TokenType{CloseParenthesis}))
	doc.expect(CloseParenthesis)
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) expressionInitial() AstNode {
	return doc.expression(0)
}

func (doc *Parser) echoIntrinsic() *Phrase {
	p := doc.start(EchoIntrinsic, false)
	doc.next(false) //echo or <?=
	p.Children = append(p.Children, doc.delimitedList(
		ExpressionList,
		doc.expressionInitial,
		isExpressionStart,
		Comma,
		nil,
		false))
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) functionStaticDeclaration() *Phrase {
	p := doc.start(FunctionStaticDeclaration, false)
	doc.next(false) //static
	p.Children = append(p.Children, doc.delimitedList(
		StaticVariableDeclarationList,
		doc.staticVariableDeclaration,
		func(t *Token) bool { return t.Type == VariableName },
		Comma,
		[]TokenType{Semicolon},
		false))

	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) globalDeclaration() *Phrase {
	p := doc.start(GlobalDeclaration, false)
	doc.next(false) //global
	p.Children = append(p.Children, doc.delimitedList(
		VariableNameList,
		doc.simpleVariable,
		isSimpleVariableStart,
		Comma,
		[]TokenType{Semicolon},
		false))
	doc.expect(Semicolon)

	return doc.end()
}

func isSimpleVariableStart(t *Token) bool {
	switch t.Type {
	case VariableName, Dollar:
		return true
	}

	return false
}

func (doc *Parser) staticVariableDeclaration() AstNode {
	p := doc.start(StaticVariableDeclaration, false)
	doc.expect(VariableName)

	if doc.peek(0).Type == Equals {
		p.Children = append(p.Children, doc.functionStaticInitialiser())
	}

	return doc.end()
}

func (doc *Parser) functionStaticInitialiser() *Phrase {
	p := doc.start(FunctionStaticInitialiser, false)
	doc.next(false) //=
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) continueStatement() *Phrase {
	p := doc.start(ContinueStatement, false)
	doc.next(false) //break/continue
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) breakStatement() *Phrase {
	p := doc.start(BreakStatement, false)
	doc.next(false) //break/continue
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) returnStatement() *Phrase {
	p := doc.start(ReturnStatement, false)
	doc.next(false) //return
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}

	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) forExpressionGroup(
	phraseType PhraseType, breakOn []TokenType) *Phrase {

	return doc.delimitedList(
		phraseType,
		doc.expressionInitial,
		isExpressionStart,
		Comma,
		breakOn,
		false)
}

func (doc *Parser) forStatement() *Phrase {
	p := doc.start(ForStatement, false)
	doc.next(false) //for
	doc.expect(OpenParenthesis)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			ForInitialiser, []TokenType{Semicolon}))
	}

	doc.expect(Semicolon)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			ForControl, []TokenType{Semicolon}))
	}

	doc.expect(Semicolon)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			ForEndOfLoop, []TokenType{CloseParenthesis}))
	}

	doc.expect(CloseParenthesis)

	t := doc.peek(0)

	if t.Type == Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]TokenType{EndFor}))
		doc.expect(EndFor)
		doc.expect(Semicolon)
	} else if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(Undefined)
	}

	return doc.end()

}

func (doc *Parser) doStatement() *Phrase {
	p := doc.start(DoStatement, false)
	doc.next(false) // do
	p.Children = append(p.Children, doc.statement())
	doc.expect(While)
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) whileStatement() *Phrase {
	p := doc.start(WhileStatement, false)
	doc.next(false) //while
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)

	t := doc.peek(0)

	if t.Type == Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]TokenType{EndWhile}))
		doc.expect(EndWhile)
		doc.expect(Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		//error
		doc.error(Undefined)
	}

	return doc.end()
}

func (doc *Parser) elseIfClause1() AstNode {
	p := doc.start(ElseIfClause, false)
	doc.next(false) //elseif
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)
	p.Children = append(p.Children, doc.statement())

	return doc.end()
}

func (doc *Parser) elseIfClause2() AstNode {
	p := doc.start(ElseIfClause, false)
	doc.next(false) //elseif
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)
	doc.expect(Colon)
	p.Children = append(p.Children, doc.statementList(
		[]TokenType{EndIf, Else, ElseIf}))

	return doc.end()
}

func (doc *Parser) elseClause1() *Phrase {
	p := doc.start(ElseClause, false)
	doc.next(false) //else
	p.Children = append(p.Children, doc.statement())

	return doc.end()
}

func (doc *Parser) elseClause2() *Phrase {
	p := doc.start(ElseClause, false)
	doc.next(false) //else
	doc.expect(Colon)
	p.Children = append(p.Children, doc.statementList([]TokenType{EndIf}))

	return doc.end()
}

func isElseIfClauseStart(t *Token) bool {
	return t.Type == ElseIf
}

func (doc *Parser) ifStatement() *Phrase {
	p := doc.start(IfStatement, false)
	doc.next(false) //if
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(CloseParenthesis)

	t := doc.peek(0)
	elseIfClausefunc := doc.elseIfClause1
	elseClausefunc := doc.elseClause1
	expectEndIf := false

	if t.Type == Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList(
			[]TokenType{ElseIf, Else, EndIf}))
		elseIfClausefunc = doc.elseIfClause2
		elseClausefunc = doc.elseClause2
		expectEndIf = true
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(Undefined)
	}

	if doc.peek(0).Type == ElseIf {
		p.Children = append(p.Children, doc.list(
			ElseIfClauseList,
			elseIfClausefunc,
			isElseIfClauseStart,
			nil,
			nil))
	}

	if doc.peek(0).Type == Else {
		p.Children = append(p.Children, elseClausefunc())
	}

	if expectEndIf {
		doc.expect(EndIf)
		doc.expect(Semicolon)
	}

	return doc.end()

}

func (doc *Parser) expressionStatement() *Phrase {
	p := doc.start(ExpressionStatement, false)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) returnType() *Phrase {
	p := doc.start(ReturnType, false)
	doc.next(false) //:
	p.Children = append(p.Children, doc.typeDeclaration())

	return doc.end()
}

func (doc *Parser) typeDeclaration() *Phrase {
	p := doc.start(TypeDeclaration, false)
	doc.optional(Question)

	switch doc.peek(0).Type {
	case Callable, Array:
		doc.next(false)
	case Name, Namespace, Backslash:
		p.Children = append(p.Children, doc.qualifiedName())
	default:
		doc.error(Undefined)
	}

	return doc.end()

}

func (doc *Parser) classConstDeclaration(p *Phrase) *Phrase {
	p.Type = ClassConstDeclaration
	doc.next(false) //const
	p.Children = append(p.Children, doc.delimitedList(
		ClassConstElementList,
		doc.classConstElement,
		isClassConstElementStartToken,
		Comma,
		[]TokenType{Semicolon},
		false))

	doc.expect(Semicolon)

	return doc.end()
}

func isExpressionStart(t *Token) bool {

	switch t.Type {
	case VariableName,
		Dollar,
		Array,
		OpenBracket,
		StringLiteral,
		Backslash,
		Name,
		Namespace,
		OpenParenthesis,
		Static,
		PlusPlus,
		MinusMinus,
		Plus,
		Minus,
		Exclamation,
		Tilde,
		AtSymbol,
		IntegerCast,
		FloatCast,
		StringCast,
		ArrayCast,
		ObjectCast,
		BooleanCast,
		UnsetCast,
		List,
		Clone,
		New,
		FloatingLiteral,
		IntegerLiteral,
		LineConstant,
		FileConstant,
		DirectoryConstant,
		TraitConstant,
		MethodConstant,
		FunctionConstant,
		NamespaceConstant,
		ClassConstant,
		StartHeredoc,
		DoubleQuote,
		Backtick,
		Print,
		Yield,
		YieldFrom,
		Function,
		Include,
		IncludeOnce,
		Require,
		RequireOnce,
		Eval,
		Empty,
		Isset,
		Exit:
		return true
	}

	return false
}

func (doc *Parser) classConstElement() AstNode {
	p := doc.start(ClassConstElement, false)
	p.Children = append(p.Children, doc.identifier())
	doc.expect(Equals)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func isPropertyElementStart(t *Token) bool {
	return t.Type == VariableName
}

func (doc *Parser) propertyDeclaration(p *Phrase) *Phrase {
	p.Type = PropertyDeclaration
	p.Children = append(p.Children, doc.delimitedList(
		PropertyElementList,
		doc.propertyElement,
		isPropertyElementStart,
		Comma,
		[]TokenType{Semicolon},
		false))
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) propertyElement() AstNode {
	p := doc.start(PropertyElement, false)
	doc.expect(VariableName)

	if doc.peek(0).Type == Equals {
		p.Children = append(p.Children, doc.propertyInitialiser())
	}

	return doc.end()
}

func (doc *Parser) propertyInitialiser() *Phrase {
	p := doc.start(PropertyInitialiser, false)
	doc.next(false) //equals
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) memberModifierList() *Phrase {
	doc.start(MemberModifierList, false)

	for isMemberModifier(doc.peek(0)) {
		doc.next(false)
	}

	return doc.end()
}

func isMemberModifier(t *Token) bool {
	switch t.Type {
	case Public,
		Protected,
		Private,
		Static,
		Abstract,
		Final:
		return true
	}

	return false
}

func (doc *Parser) qualifiedNameList(breakOn []TokenType) *Phrase {

	return doc.delimitedList(
		QualifiedNameList,
		doc.qualifiedName,
		isQualifiedNameStart,
		Comma,
		breakOn,
		false)
}

func (doc *Parser) objectCreationExpression() *Phrase {
	p := doc.start(ObjectCreationExpression, false)
	doc.next(false) //new
	if doc.peek(0).Type == Class {
		p.Children = append(p.Children, doc.anonymousClassDeclaration())

		return doc.end()
	}

	p.Children = append(p.Children, doc.typeDesignator(ClassTypeDesignator))

	if doc.optional(OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}

		doc.expect(CloseParenthesis)
	}

	return doc.end()

}

func (doc *Parser) typeDesignator(phraseType PhraseType) *Phrase {
	p := doc.start(phraseType, false)
	part := doc.classTypeDesignatorAtom()

	for {
		switch doc.peek(0).Type {
		case OpenBracket:
			part = doc.subscriptExpression(part, CloseBracket)
			continue
		case OpenBrace:
			part = doc.subscriptExpression(part, CloseBrace)
			continue
		case Arrow:
			part = doc.propertyAccessExpression(part)
			continue
		case ColonColon:
			staticPropNode := doc.start(ScopedPropertyAccessExpression, false)
			staticPropNode.Children = append(staticPropNode.Children, part)
			doc.next(false) //::
			staticPropNode.Children = append(staticPropNode.Children,
				doc.restrictedScopedMemberName())
			part = doc.end()
			continue
		}

		break
	}

	p.Children = append(p.Children, part)

	return doc.end()
}

func (doc *Parser) restrictedScopedMemberName() *Phrase {
	p := doc.start(ScopedMemberName, false)
	t := doc.peek(0)

	switch t.Type {
	case VariableName:
		//Spec says this should be SimpleVariable
		//leaving as a token as this avoids confusion between
		//static property names and simple variables
		doc.next(false)
	case Dollar:
		p.Children = append(p.Children, doc.simpleVariable())
	default:
		doc.error(Undefined)
	}

	return doc.end()
}

func (doc *Parser) classTypeDesignatorAtom() AstNode {
	t := doc.peek(0)

	switch t.Type {
	case Static:
		return doc.relativeScope()
	case VariableName, Dollar:
		return doc.simpleVariable()
	case Name, Namespace, Backslash:
		return doc.qualifiedName()
	default:
		doc.start(ErrorClassTypeDesignatorAtom, false)
		doc.error(Undefined)

		return doc.end()
	}
}

func (doc *Parser) cloneExpression() *Phrase {
	p := doc.start(CloneExpression, false)
	doc.next(false) //clone
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) listIntrinsic() *Phrase {
	p := doc.start(ListIntrinsic, false)
	doc.next(false) //list
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.arrayInitialiserList(CloseParenthesis))
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) unaryExpression(phraseType PhraseType) *Phrase {
	p := doc.start(phraseType, false)
	op := doc.next(false) //op

	switch phraseType {
	case PrefixDecrementExpression, PrefixIncrementExpression:
		p.Children = append(p.Children, doc.variable(doc.variableAtom(0)))
	default:
		precendence, _ := precedenceAssociativityTuple(op)
		p.Children = append(p.Children, doc.expression(precendence))
	}

	return doc.end()
}

func (doc *Parser) anonymousFunctionHeader() *Phrase {
	p := doc.start(AnonymousFunctionHeader, false)
	doc.optional(Static)
	doc.next(false) //function
	doc.optional(Ampersand)
	doc.expect(OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			Comma,
			[]TokenType{CloseParenthesis},
			false))
	}

	doc.expect(CloseParenthesis)

	if doc.peek(0).Type == Use {
		p.Children = append(p.Children, doc.anonymousFunctionUseClause())
	}

	if doc.peek(0).Type == Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func (doc *Parser) anonymousFunctionCreationExpression() *Phrase {
	p := doc.start(AnonymousFunctionCreationExpression, false)

	p.Children = append(p.Children, doc.anonymousFunctionHeader())
	p.Children = append(p.Children, doc.functionDeclarationBody())

	return doc.end()
}

func isAnonymousFunctionUseVariableStart(t *Token) bool {
	return t.Type == VariableName || t.Type == Ampersand
}

func (doc *Parser) anonymousFunctionUseClause() *Phrase {
	p := doc.start(AnonymousFunctionUseClause, false)
	doc.next(false) //use
	doc.expect(OpenParenthesis)
	p.Children = append(p.Children, doc.delimitedList(
		ClosureUseList,
		doc.anonymousFunctionUseVariable,
		isAnonymousFunctionUseVariableStart,
		Comma,
		[]TokenType{CloseParenthesis},
		false))
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) anonymousFunctionUseVariable() AstNode {
	doc.start(AnonymousFunctionUseVariable, false)
	doc.optional(Ampersand)
	doc.expect(VariableName)

	return doc.end()
}

func isTypeDeclarationStart(t *Token) bool {
	switch t.Type {
	case Backslash,
		Name,
		Namespace,
		Question,
		Array,
		Callable:
		return true
	}

	return false
}

func (doc *Parser) parameterDeclaration() AstNode {
	p := doc.start(ParameterDeclaration, false)

	if isTypeDeclarationStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.typeDeclaration())
	}

	doc.optional(Ampersand)
	doc.optional(Ellipsis)
	doc.expect(VariableName)

	if doc.peek(0).Type == Equals {
		doc.next(false)
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) variable(variableAtomNode AstNode) AstNode {
	count := 0

	for {
		count++
		switch doc.peek(0).Type {
		case ColonColon:
			variableAtomNode = doc.scopedAccessExpression(variableAtomNode)
			continue
		case Arrow:
			variableAtomNode = doc.propertyOrMethodAccessExpression(variableAtomNode)
			continue
		case OpenBracket:
			variableAtomNode = doc.subscriptExpression(variableAtomNode, CloseBracket)
			continue
		case OpenBrace:
			variableAtomNode = doc.subscriptExpression(variableAtomNode, CloseBrace)
			continue
		case OpenParenthesis:
			variableAtomNode = doc.functionCallExpression(variableAtomNode)
			continue
		default:
			//only simple variable atoms qualify as variables
			p := variableAtomNode.(*Phrase)

			if count == 1 && p.Type != SimpleVariable {
				errNode := doc.start(ErrorVariable, true)
				errNode.Children = append(errNode.Children, variableAtomNode)
				doc.error(Undefined)

				return doc.end()
			}
		}

		break
	}

	return variableAtomNode
}

func (doc *Parser) functionCallExpression(lhs AstNode) *Phrase {
	p := doc.start(FunctionCallExpression, true)
	p.Children = append(p.Children, lhs)
	doc.expect(OpenParenthesis)
	if isArgumentStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.argumentList())
	}
	doc.expect(CloseParenthesis)

	return doc.end()
}

func (doc *Parser) scopedAccessExpression(lhs AstNode) *Phrase {
	p := doc.start(ErrorScopedAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //::
	p.Children = append(p.Children, doc.scopedMemberName(p))

	if doc.optional(OpenParenthesis) != nil {
		p.Type = ScopedCallExpression
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}

		doc.expect(CloseParenthesis)

		return doc.end()
	} else if p.Type == ScopedCallExpression {
		//error
		doc.error(Undefined)
	}

	return doc.end()

}

func (doc *Parser) scopedMemberName(parent *Phrase) *Phrase {
	p := doc.start(ScopedMemberName, false)
	t := doc.peek(0)

	switch t.Type {
	case OpenBrace:
		parent.Type = ScopedCallExpression
		p.Children = append(p.Children, doc.encapsulatedExpression(
			OpenBrace, CloseBrace))
	case VariableName:
		//Spec says this should be SimpleVariable
		//leaving as a token as this avoids confusion between
		//static property names and simple variables
		parent.Type = ScopedPropertyAccessExpression
		doc.next(false)
	case Dollar:
		p.Children = append(p.Children, doc.simpleVariable())
		parent.Type = ScopedPropertyAccessExpression
	default:
		if t.Type == Name || isSemiReservedToken(t) {
			p.Children = append(p.Children, doc.identifier())
			parent.Type = ClassConstantAccessExpression
		} else {
			//error
			doc.error(Undefined)
		}
	}

	return doc.end()

}

func (doc *Parser) propertyAccessExpression(lhs AstNode) *Phrase {
	p := doc.start(PropertyAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //->
	p.Children = append(p.Children, doc.memberName())

	return doc.end()
}

func (doc *Parser) propertyOrMethodAccessExpression(lhs AstNode) *Phrase {

	p := doc.start(PropertyAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //->
	p.Children = append(p.Children, doc.memberName())

	if doc.optional(OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}
		p.Type = MethodCallExpression
		doc.expect(CloseParenthesis)
	}

	return doc.end()

}

func (doc *Parser) memberName() *Phrase {
	p := doc.start(MemberName, false)

	switch doc.peek(0).Type {
	case Name:
		doc.next(false)
	case OpenBrace:
		p.Children = append(p.Children, doc.encapsulatedExpression(
			OpenBrace, CloseBrace))
	case Dollar, VariableName:
		p.Children = append(p.Children, doc.simpleVariable())
	default:
		doc.error(Undefined)
	}

	return doc.end()
}

func (doc *Parser) subscriptExpression(
	lhs AstNode, closeTokenType TokenType) *Phrase {

	p := doc.start(SubscriptExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) // [ or {
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}

	doc.expect(closeTokenType)

	return doc.end()
}

func (doc *Parser) argumentList() *Phrase {
	return doc.delimitedList(
		ArgumentExpressionList,
		doc.argumentExpression,
		isArgumentStart,
		Comma,
		[]TokenType{CloseParenthesis},
		false)
}

func isArgumentStart(t *Token) bool {
	return t.Type == Ellipsis || isExpressionStart(t)
}

func (doc *Parser) variadicUnpacking() *Phrase {
	p := doc.start(VariadicUnpacking, false)
	doc.next(false) //...
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) argumentExpression() AstNode {
	if doc.peek(0).Type == Ellipsis {
		return doc.variadicUnpacking()
	} else {
		return doc.expression(0)
	}
}

func (doc *Parser) qualifiedName() AstNode {
	p := doc.start(QualifiedName, false)
	t := doc.peek(0)

	if t.Type == Backslash {
		doc.next(false)
		p.Type = FullyQualifiedName
	} else if t.Type == Namespace {
		p.Type = RelativeQualifiedName
		doc.next(false)
		doc.expect(Backslash)
	}

	p.Children = append(p.Children, doc.namespaceName())

	return doc.end()
}

func isQualifiedNameStart(t *Token) bool {
	switch t.Type {
	case Backslash, Name, Namespace:
		return true
	}

	return false
}

func (doc *Parser) shortArrayCreationExpression(precedence int) *Phrase {
	p := doc.start(ArrayCreationExpression, false)
	doc.next(false) //[
	if isArrayElementStart(doc.peek(0)) || (precedence == 0 && doc.peek(0).Type == Comma) {
		p.Children = append(p.Children, doc.arrayInitialiserList(CloseBracket))
	}
	doc.expect(CloseBracket)

	return doc.end()
}

func (doc *Parser) longArrayCreationExpression() *Phrase {
	p := doc.start(ArrayCreationExpression, false)
	doc.next(false) //array
	doc.expect(OpenParenthesis)

	if isArrayElementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.arrayInitialiserList(CloseParenthesis))
	}

	doc.expect(CloseParenthesis)

	return doc.end()
}

func isArrayElementStart(t *Token) bool {
	return t.Type == Ampersand || isExpressionStart(t)
}

func (doc *Parser) arrayInitialiserList(breakOn TokenType) *Phrase {

	p := doc.start(ArrayInitialiserList, false)
	var t *Token

	arrayInitialiserListRecoverSet := []TokenType{breakOn, Comma}
	doc.recoverSetStack = append(doc.recoverSetStack, arrayInitialiserListRecoverSet)

	for {
		//an array can have empty elements
		if isArrayElementStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.arrayElement())
		}

		t = doc.peek(0)

		if t.Type == Comma {
			doc.next(false)
		} else if t.Type == breakOn {
			break
		} else {
			doc.error(Undefined)
			//check for missing delimeter
			if isArrayElementStart(t) {
				continue
			} else {
				//skip until recover token
				doc.defaultSyncStrategy()
				t = doc.peek(0)
				if t.Type == Comma || t.Type == breakOn {
					continue
				}
			}

			break
		}

	}

	doc.recoverSetStack = doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return doc.end()
}

func (doc *Parser) arrayValue() *Phrase {
	p := doc.start(ArrayValue, false)
	doc.optional(Ampersand)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) arrayKey() *Phrase {
	p := doc.start(ArrayKey, false)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) arrayElement() *Phrase {
	p := doc.start(ArrayElement, false)

	if doc.peek(0).Type == Ampersand {
		p.Children = append(p.Children, doc.arrayValue())

		return doc.end()
	}

	keyOrValue := doc.arrayKey()
	p.Children = append(p.Children, keyOrValue)

	if doc.optional(FatArrow) == nil {
		keyOrValue.Type = ArrayValue

		return doc.end()
	}

	p.Children = append(p.Children, doc.arrayValue())

	return doc.end()
}

func (doc *Parser) encapsulatedExpression(
	openTokenType TokenType, closeTokenType TokenType) *Phrase {
	p := doc.start(EncapsulatedExpression, false)
	doc.expect(openTokenType)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(closeTokenType)

	return doc.end()
}

func (doc *Parser) relativeScope() *Phrase {
	doc.start(RelativeScope, false)
	doc.next(false)

	return doc.end()
}

func (doc *Parser) variableAtom(precedence int) AstNode {
	t := doc.peek(0)
	switch t.Type {
	case VariableName, Dollar:
		return doc.simpleVariable()
	case OpenParenthesis:
		return doc.encapsulatedExpression(OpenParenthesis, CloseParenthesis)
	case Array:
		return doc.longArrayCreationExpression()
	case OpenBracket:
		return doc.shortArrayCreationExpression(precedence)
	case StringLiteral:
		return doc.next(true)
	case Static:
		return doc.relativeScope()
	case Name, Namespace, Backslash:
		return doc.qualifiedName()
	default:
		//error
		doc.start(ErrorVariableAtom, false)
		doc.error(Undefined)

		return doc.end()
	}
}

func (doc *Parser) simpleVariable() AstNode {
	p := doc.start(SimpleVariable, false)
	t := doc.expectOneOf([]TokenType{VariableName, Dollar})

	if t != nil && t.Type == Dollar {
		t = doc.peek(0)
		if t.Type == OpenBrace {
			p.Children = append(p.Children, doc.encapsulatedExpression(
				OpenBrace, CloseBrace))
		} else if t.Type == Dollar || t.Type == VariableName {
			p.Children = append(p.Children, doc.simpleVariable())
		} else {
			doc.error(Undefined)
		}
	}

	return doc.end()
}

func (doc *Parser) haltCompilerStatement() *Phrase {
	doc.start(HaltCompilerStatement, false)
	doc.next(false) // __halt_compiler
	doc.expect(OpenParenthesis)
	doc.expect(CloseParenthesis)
	doc.expect(Semicolon)

	return doc.end()
}

func (doc *Parser) namespaceUseDeclaration() *Phrase {
	p := doc.start(NamespaceUseDeclaration, false)
	doc.next(false) //use
	doc.optionalOneOf([]TokenType{Function, Const})
	doc.optional(Backslash)
	nsNameNode := doc.namespaceName()
	t := doc.peek(0)

	if t.Type == Backslash || t.Type == OpenBrace {
		p.Children = append(p.Children, nsNameNode)
		doc.expect(Backslash)
		doc.expect(OpenBrace)
		p.Children = append(p.Children, doc.delimitedList(
			NamespaceUseGroupClauseList,
			doc.namespaceUseGroupClause,
			isNamespaceUseGroupClauseStartToken,
			Comma,
			[]TokenType{CloseBrace},
			false))
		doc.expect(CloseBrace)
		doc.expect(Semicolon)

		return doc.end()
	}

	p.Children = append(p.Children, doc.delimitedList(
		NamespaceUseClauseList,
		doc.namespaceUseClauseFunction(nsNameNode),
		isNamespaceUseClauseStartToken,
		Comma,
		[]TokenType{Semicolon},
		true))

	doc.expect(Semicolon)

	return doc.end()
}

func isNamespaceUseClauseStartToken(t *Token) bool {
	return t.Type == Name || t.Type == Backslash
}

func (doc *Parser) namespaceUseClauseFunction(nsName *Phrase) func() AstNode {

	return func() AstNode {
		p := doc.start(NamespaceUseClause, nsName != nil)

		if nsName != nil {
			p.Children = append(p.Children, nsName)
			nsName = nil
		} else {
			p.Children = append(p.Children, doc.namespaceName())
		}

		if doc.peek(0).Type == As {
			p.Children = append(p.Children, doc.namespaceAliasingClause())
		}

		return doc.end()
	}
}

func (doc *Parser) delimitedList(
	phraseType PhraseType,
	elementFunction func() AstNode,
	elementStartPredicate func(*Token) bool,
	delimiter TokenType,
	breakOn []TokenType,
	doNotPushHiddenToParent bool) *Phrase {

	p := doc.start(phraseType, doNotPushHiddenToParent)
	var t *Token
	var delimitedListRecoverSet []TokenType
	if breakOn != nil {
		delimitedListRecoverSet = make([]TokenType, len(breakOn))
		copy(delimitedListRecoverSet, breakOn)
	} else {
		delimitedListRecoverSet = make([]TokenType, 0)
	}

	delimitedListRecoverSet = append(delimitedListRecoverSet, delimiter)
	doc.recoverSetStack = append(doc.recoverSetStack, delimitedListRecoverSet)

	for {
		p.Children = append(p.Children, elementFunction())
		t = doc.peek(0)

		if t.Type == delimiter {
			doc.next(false)
		} else if breakOn == nil || tokenTypeIndexOf(breakOn, t.Type) >= 0 {
			break
		} else {
			doc.error(Undefined)
			//check for missing delimeter
			if elementStartPredicate(t) {
				continue
			} else if breakOn != nil {
				//skip until breakOn or delimiter token or whatever else is in recover set
				doc.defaultSyncStrategy()
				if doc.peek(0).Type == delimiter {
					continue
				}
			}

			break
		}

	}

	doc.recoverSetStack = doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return doc.end()
}

func isNamespaceUseGroupClauseStartToken(t *Token) bool {
	switch t.Type {
	case Const, Function, Name:
		return true
	}

	return false
}

func (doc *Parser) namespaceUseGroupClause() AstNode {
	p := doc.start(NamespaceUseGroupClause, false)
	doc.optionalOneOf([]TokenType{Function, Const})
	p.Children = append(p.Children, doc.namespaceName())

	if doc.peek(0).Type == As {
		p.Children = append(p.Children, doc.namespaceAliasingClause())
	}

	return doc.end()
}

func (doc *Parser) namespaceAliasingClause() *Phrase {
	doc.start(NamespaceAliasingClause, false)
	doc.next(false) //as
	doc.expect(Name)

	return doc.end()
}

func (doc *Parser) namespaceDefinition() *Phrase {
	p := doc.start(NamespaceDefinition, false)
	doc.next(false) //namespace
	if doc.peek(0).Type == Name {

		p.Children = append(p.Children, doc.namespaceName())
		t := doc.expectOneOf([]TokenType{Semicolon, OpenBrace})
		if t == nil || t.Type != OpenBrace {
			return doc.end()
		}

	} else {
		doc.expect(OpenBrace)
	}

	p.Children = append(p.Children, doc.statementList([]TokenType{CloseBrace}))
	doc.expect(CloseBrace)

	return doc.end()
}

func (doc *Parser) namespaceName() *Phrase {
	doc.start(NamespaceName, false)
	doc.expect(Name)

	for {
		if doc.peek(0).Type == Backslash &&
			doc.peek(1).Type == Name {
			doc.next(false)
			doc.next(false)
		} else {
			break
		}
	}

	return doc.end()
}

func isReservedToken(t *Token) bool {
	switch t.Type {
	case Include,
		IncludeOnce,
		Eval,
		Require,
		RequireOnce,
		Or,
		Xor,
		And,
		InstanceOf,
		New,
		Clone,
		Exit,
		If,
		ElseIf,
		Else,
		EndIf,
		Echo,
		Do,
		While,
		EndWhile,
		For,
		EndFor,
		ForEach,
		EndForeach,
		Declare,
		EndDeclare,
		As,
		Try,
		Catch,
		Finally,
		Throw,
		Use,
		InsteadOf,
		Global,
		Var,
		Unset,
		Isset,
		Empty,
		Continue,
		Goto,
		Function,
		Const,
		Return,
		Print,
		Yield,
		List,
		Switch,
		EndSwitch,
		Case,
		Default,
		Break,
		Array,
		Callable,
		Extends,
		Implements,
		Namespace,
		Trait,
		Interface,
		Class,
		ClassConstant,
		TraitConstant,
		FunctionConstant,
		MethodConstant,
		LineConstant,
		FileConstant,
		DirectoryConstant,
		NamespaceConstant:
		return true
	}

	return false
}

func isSemiReservedToken(t *Token) bool {
	switch t.Type {
	case Static,
		Abstract,
		Final,
		Private,
		Protected,
		Public:
		return true
	}

	return isReservedToken(t)
}

func isStatementStart(t *Token) bool {
	switch t.Type {
	case Namespace,
		Use,
		HaltCompiler,
		Const,
		Function,
		Class,
		Abstract,
		Final,
		Trait,
		Interface,
		OpenBrace,
		If,
		While,
		Do,
		For,
		Switch,
		Break,
		Continue,
		Return,
		Global,
		Static,
		Echo,
		Unset,
		ForEach,
		Declare,
		Try,
		Throw,
		Goto,
		Name,
		Semicolon,
		CloseTag,
		Text,
		OpenTag,
		OpenTagEcho:
		return true
	}

	return isExpressionStart(t)
}
