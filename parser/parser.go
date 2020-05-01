package parser

import (
	"errors"
	"reflect"

	"github.com/john-nguyen09/go-phpparser/lexer"
	"github.com/john-nguyen09/go-phpparser/phrase"
)

var statementListRecoverSet = []lexer.TokenType{lexer.Use,
	lexer.HaltCompiler,
	lexer.Const,
	lexer.Function,
	lexer.Class,
	lexer.Abstract,
	lexer.Final,
	lexer.Trait,
	lexer.Interface,
	lexer.OpenBrace,
	lexer.If,
	lexer.While,
	lexer.Do,
	lexer.For,
	lexer.Switch,
	lexer.Break,
	lexer.Continue,
	lexer.Return,
	lexer.Global,
	lexer.Static,
	lexer.Echo,
	lexer.Unset,
	lexer.ForEach,
	lexer.Declare,
	lexer.Try,
	lexer.Throw,
	lexer.Goto,
	lexer.Semicolon,
	lexer.CloseTag,
	lexer.OpenTagEcho,
	lexer.Text,
	lexer.OpenTag}

var classMemberDeclarationListRecoverSet = []lexer.TokenType{
	lexer.Public,
	lexer.Protected,
	lexer.Private,
	lexer.Static,
	lexer.Abstract,
	lexer.Final,
	lexer.Function,
	lexer.Var,
	lexer.Const,
	lexer.Use}

var encapsulatedVariableListRecoverSet = []lexer.TokenType{
	lexer.EncapsulatedAndWhitespace,
	lexer.DollarCurlyOpen,
	lexer.CurlyOpen}

type Associativity int

const (
	None Associativity = iota
	Left
	Right
)

func precedenceAssociativityTuple(t *lexer.Token) (int, Associativity) {
	switch t.Type {
	case lexer.AsteriskAsterisk:
		return 48, Right
	case lexer.PlusPlus:
		return 47, Right
	case lexer.MinusMinus:
		return 47, Right
	case lexer.Tilde:
		return 47, Right
	case lexer.IntegerCast:
		return 47, Right
	case lexer.FloatCast:
		return 47, Right
	case lexer.StringCast:
		return 47, Right
	case lexer.ArrayCast:
		return 47, Right
	case lexer.ObjectCast:
		return 47, Right
	case lexer.BooleanCast:
		return 47, Right
	case lexer.UnsetCast:
		return 47, Right
	case lexer.AtSymbol:
		return 47, Right
	case lexer.InstanceOf:
		return 46, None
	case lexer.Exclamation:
		return 45, Right
	case lexer.Asterisk:
		return 44, Left
	case lexer.ForwardSlash:
		return 44, Left
	case lexer.Percent:
		return 44, Left
	case lexer.Plus:
		return 43, Left
	case lexer.Minus:
		return 43, Left
	case lexer.Dot:
		return 43, Left
	case lexer.LessThanLessThan:
		return 42, Left
	case lexer.GreaterThanGreaterThan:
		return 42, Left
	case lexer.LessThan:
		return 41, None
	case lexer.GreaterThan:
		return 41, None
	case lexer.LessThanEquals:
		return 41, None
	case lexer.GreaterThanEquals:
		return 41, None
	case lexer.EqualsEquals:
		return 40, None
	case lexer.EqualsEqualsEquals:
		return 40, None
	case lexer.ExclamationEquals:
		return 40, None
	case lexer.ExclamationEqualsEquals:
		return 40, None
	case lexer.Spaceship:
		return 40, None
	case lexer.Ampersand:
		return 39, Left
	case lexer.Caret:
		return 38, Left
	case lexer.Bar:
		return 37, Left
	case lexer.AmpersandAmpersand:
		return 36, Left
	case lexer.BarBar:
		return 35, Left
	case lexer.QuestionQuestion:
		return 34, Right
	case lexer.Question:
		return 33, Left //?: ternary
	case lexer.Equals:
		return 32, Right
	case lexer.DotEquals:
		return 32, Right
	case lexer.PlusEquals:
		return 32, Right
	case lexer.MinusEquals:
		return 32, Right
	case lexer.AsteriskEquals:
		return 32, Right
	case lexer.ForwardslashEquals:
		return 32, Right
	case lexer.PercentEquals:
		return 32, Right
	case lexer.AsteriskAsteriskEquals:
		return 32, Right
	case lexer.AmpersandEquals:
		return 32, Right
	case lexer.BarEquals:
		return 32, Right
	case lexer.CaretEquals:
		return 32, Right
	case lexer.LessThanLessThanEquals:
		return 32, Right
	case lexer.GreaterThanGreaterThanEquals:
		return 32, Right
	case lexer.And:
		return 31, Left
	case lexer.Xor:
		return 30, Left
	case lexer.Or:
		return 29, Left
	}

	return 0, None
}

func binaryOpToPhraseType(t *lexer.Token) phrase.PhraseType {
	switch t.Type {
	case lexer.Question:
		return phrase.TernaryExpression
	case lexer.Dot, lexer.Plus, lexer.Minus:
		return phrase.AdditiveExpression
	case lexer.Bar, lexer.Ampersand, lexer.Caret:
		return phrase.BitwiseExpression
	case lexer.Asterisk, lexer.ForwardSlash, lexer.Percent:
		return phrase.MultiplicativeExpression
	case lexer.AsteriskAsterisk:
		return phrase.ExponentiationExpression
	case lexer.LessThanLessThan, lexer.GreaterThanGreaterThan:
		return phrase.ShiftExpression
	case lexer.AmpersandAmpersand, lexer.BarBar, lexer.And, lexer.Or, lexer.Xor:
		return phrase.LogicalExpression
	case lexer.EqualsEqualsEquals,
		lexer.ExclamationEqualsEquals,
		lexer.EqualsEquals,
		lexer.ExclamationEquals:
		return phrase.EqualityExpression
	case lexer.LessThan,
		lexer.LessThanEquals,
		lexer.GreaterThan,
		lexer.GreaterThanEquals,
		lexer.Spaceship:
		return phrase.RelationalExpression
	case lexer.QuestionQuestion:
		return phrase.CoalesceExpression
	case lexer.Equals:
		return phrase.SimpleAssignmentExpression
	case lexer.PlusEquals,
		lexer.MinusEquals,
		lexer.AsteriskEquals,
		lexer.AsteriskAsteriskEquals,
		lexer.ForwardslashEquals,
		lexer.DotEquals,
		lexer.PercentEquals,
		lexer.AmpersandEquals,
		lexer.BarEquals,
		lexer.CaretEquals,
		lexer.LessThanLessThanEquals,
		lexer.GreaterThanGreaterThanEquals:
		return phrase.CompoundAssignmentExpression
	case lexer.InstanceOf:
		return phrase.InstanceOfExpression
	default:
		return phrase.Unknown
	}
}

type Parser struct {
	tokenOffset     int
	tokenBuffer     []*lexer.Token
	phraseStack     []*phrase.Phrase
	errorPhrase     *phrase.ParseError
	recoverSetStack [][]lexer.TokenType
}

func tokenTypeIndexOf(haystack []lexer.TokenType, needle lexer.TokenType) int {
	for index, tokenType := range haystack {
		if tokenType == needle {
			return index
		}
	}

	return -1
}

func Parse(source []byte) *phrase.Phrase {
	doc := &Parser{0,
		lexer.Lex(source),
		make([]*phrase.Phrase, 0),
		nil,
		make([][]lexer.TokenType, 0)}
	stmtList := doc.statementList([]lexer.TokenType{lexer.EndOfFile})
	//append trailing hidden tokens
	doc.hidden(stmtList)

	return stmtList
}

func (doc *Parser) popRecover() []lexer.TokenType {
	var lastRecoverSet []lexer.TokenType

	lastRecoverSet, doc.recoverSetStack = doc.recoverSetStack[len(doc.recoverSetStack)-1],
		doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return lastRecoverSet
}

func (doc *Parser) start(phraseType phrase.PhraseType, dontPushHiddenToParent bool) *phrase.Phrase {
	//parent node gets hidden tokens between children
	if !dontPushHiddenToParent {
		doc.hidden(nil)
	}

	p := phrase.NewPhrase(phraseType, make([]phrase.AstNode, 0))

	doc.phraseStack = append(doc.phraseStack, p)

	return p
}

func (doc *Parser) end() *phrase.Phrase {
	var result *phrase.Phrase

	result, doc.phraseStack =
		doc.phraseStack[len(doc.phraseStack)-1], doc.phraseStack[:len(doc.phraseStack)-1]

	return result
}

func (doc *Parser) hidden(p *phrase.Phrase) {
	if p == nil {
		if len(doc.phraseStack) > 0 {
			p = doc.phraseStack[len(doc.phraseStack)-1]
		}
	}

	var t *lexer.Token

	for {
		if doc.tokenOffset < len(doc.tokenBuffer) {
			t = doc.tokenBuffer[doc.tokenOffset]
			doc.tokenOffset++
		} else {
			break
		}

		if t.Type < lexer.Comment {
			doc.tokenOffset--
			break
		} else {
			p.Children = append(p.Children, t)
		}
	}
}

func (doc *Parser) optional(tokenType lexer.TokenType) *lexer.Token {
	if tokenType == doc.peek(0).Type {
		doc.errorPhrase = nil

		return doc.next(false)
	}

	return nil
}

func (doc *Parser) optionalOneOf(tokenTypes []lexer.TokenType) *lexer.Token {
	if tokenTypeIndexOf(tokenTypes, doc.peek(0).Type) >= 0 {
		doc.errorPhrase = nil

		return doc.next(false)
	}

	return nil
}

func (doc *Parser) next(doNotPush bool) *lexer.Token {
	if doc.tokenOffset >= len(doc.tokenBuffer) {
		return doc.tokenBuffer[len(doc.tokenBuffer)-1]
	}
	t := doc.tokenBuffer[doc.tokenOffset]
	doc.tokenOffset++

	if t.Type == lexer.EndOfFile {
		return t
	}

	lastPhrase := doc.phraseStack[len(doc.phraseStack)-1]
	if t.Type >= lexer.Comment {
		//hidden token
		lastPhrase.Children = append(lastPhrase.Children, t)

		return doc.next(doNotPush)
	} else if !doNotPush {
		lastPhrase.Children = append(lastPhrase.Children, t)
	}

	return t
}

func (doc *Parser) expect(tokenType lexer.TokenType) *lexer.Token {
	t := doc.peek(0)

	if t.Type == tokenType {
		doc.errorPhrase = nil

		return doc.next(false)
	} else if tokenType == lexer.Semicolon && t.Type == lexer.CloseTag {
		//implicit end statement
		return t
	} else {
		doc.error(tokenType)
		//test skipping a single token to sync
		if doc.peek(1).Type == tokenType {
			doc.skip(func(x *lexer.Token) bool {
				return x.Type == tokenType
			})
			doc.errorPhrase = nil
			return doc.next(false) //tokenType
		}
	}

	return nil
}

func (doc *Parser) expectOneOf(tokenTypes []lexer.TokenType) *lexer.Token {
	t := doc.peek(0)

	if tokenTypeIndexOf(tokenTypes, t.Type) >= 0 {
		doc.errorPhrase = nil

		return doc.next(false)
	} else if tokenTypeIndexOf(tokenTypes, lexer.Semicolon) >= 0 && t.Type == lexer.CloseTag {
		//implicit end statement
		return t
	}

	doc.error(lexer.Undefined)
	//test skipping single token to sync
	if tokenTypeIndexOf(tokenTypes, doc.peek(1).Type) >= 0 {
		doc.skip(func(x *lexer.Token) bool {
			return tokenTypeIndexOf(tokenTypes, x.Type) >= 0
		})
		doc.errorPhrase = nil

		return doc.next(false) //tokenType
	}

	return nil
}

func (doc *Parser) peek(n int) *lexer.Token {
	k := n + 1
	bufferPos := doc.tokenOffset - 1
	var t *lexer.Token

	for {
		bufferPos++
		if bufferPos >= len(doc.tokenBuffer) {
			return doc.tokenBuffer[len(doc.tokenBuffer)-1]
		}
		t = doc.tokenBuffer[bufferPos]

		if t.Type < lexer.Comment {
			//not a hidden token
			k--
		}

		if t.Type == lexer.EndOfFile || k == 0 {
			break
		}

	}

	return t
}

/**
* skipped tokens get pushed to error phrase children
 */
func (doc *Parser) skip(predicate func(*lexer.Token) bool) {

	var t *lexer.Token

	for {
		if doc.tokenOffset < len(doc.tokenBuffer) {
			t = doc.tokenBuffer[doc.tokenOffset]
			doc.tokenOffset++
		} else {
			break
		}

		if predicate(t) || t.Type == lexer.EndOfFile {
			doc.tokenOffset--
			break
		} else {
			doc.errorPhrase.Children = append(doc.errorPhrase.Children, t)
		}
	}

}

func (doc *Parser) error(expected lexer.TokenType) {

	//dont report errors if recovering from another
	if doc.errorPhrase != nil {
		return
	}

	doc.errorPhrase = phrase.NewParseErr(
		phrase.Error, make([]phrase.AstNode, 0), doc.peek(0), expected)

	lastPhrase := doc.phraseStack[len(doc.phraseStack)-1]

	lastPhrase.Children = append(lastPhrase.Children, doc.errorPhrase)
}

func (doc *Parser) list(phraseType phrase.PhraseType,
	elementFunction func() phrase.AstNode,
	elementStartPredicate func(*lexer.Token) bool,
	breakOn []lexer.TokenType,
	recoverSet []lexer.TokenType) *phrase.Phrase {

	p := doc.start(phraseType, false)

	var t *lexer.Token
	recoveryAttempted := false
	var listRecoverSet []lexer.TokenType

	if recoverSet != nil {
		listRecoverSet = recoverSet
	} else {
		listRecoverSet = make([]lexer.TokenType, 0)
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
			doc.error(lexer.Undefined)
			//attempt to sync with token stream
			t = doc.peek(1)
			if elementStartPredicate(t) || tokenTypeIndexOf(breakOn, t.Type) >= 0 {
				doc.skip(func(x *lexer.Token) bool {
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

	mergedRecoverTokenTypeArray := make([]lexer.TokenType, 0)

	for n := len(doc.recoverSetStack) - 1; n >= 0; n-- {
		mergedRecoverTokenTypeArray = append(mergedRecoverTokenTypeArray, doc.recoverSetStack[n]...)
	}

	mergedRecoverTokenTypeSet := make(map[lexer.TokenType]bool, len(mergedRecoverTokenTypeArray))

	for _, recoverTokenType := range mergedRecoverTokenTypeArray {
		if _, ok := mergedRecoverTokenTypeSet[recoverTokenType]; !ok {
			mergedRecoverTokenTypeSet[recoverTokenType] = true
		}
	}

	doc.skip(func(x *lexer.Token) bool {
		_, ok := mergedRecoverTokenTypeSet[x.Type]

		return ok
	})
}

func (doc *Parser) statementList(breakOn []lexer.TokenType) *phrase.Phrase {
	return doc.list(
		phrase.StatementList,
		doc.statement,
		isStatementStart,
		breakOn,
		statementListRecoverSet[:])
}

func (doc *Parser) constDeclaration() *phrase.Phrase {
	p := doc.start(phrase.ConstDeclaration, false)
	doc.next(false) //const

	p.Children = append(p.Children, doc.delimitedList(
		phrase.ConstElementList,
		doc.constElement,
		isConstElementStartToken,
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		false))

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func isClassConstElementStartToken(t *lexer.Token) bool {
	return t.Type == lexer.Name || isSemiReservedToken(t)
}

func isConstElementStartToken(t *lexer.Token) bool {
	return t.Type == lexer.Name
}

func (doc *Parser) constElement() phrase.AstNode {
	p := doc.start(phrase.ConstElement, false)

	doc.expect(lexer.Name)
	doc.expect(lexer.Equals)

	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) expression(minPrecedence int) phrase.AstNode {
	var precedence int
	var associativity Associativity
	var op *lexer.Token
	var p *phrase.Phrase
	var binaryPhraseType phrase.PhraseType

	lhs := doc.expressionAtom(minPrecedence)

	for {
		op = doc.peek(0)
		binaryPhraseType = binaryOpToPhraseType(op)

		if binaryPhraseType == phrase.Unknown {
			break
		}

		precedence, associativity = precedenceAssociativityTuple(op)

		if precedence < minPrecedence {
			break
		}

		if associativity == Left {
			precedence++
		}

		if binaryPhraseType == phrase.TernaryExpression {
			lhs = doc.ternaryExpression(lhs)

			continue
		}

		p = doc.start(binaryPhraseType, true)
		p.Children = append(p.Children, lhs)
		doc.next(false)

		if binaryPhraseType == phrase.InstanceOfExpression {
			p.Children = append(p.Children, doc.typeDesignator(phrase.InstanceofTypeDesignator))
		} else {
			if binaryPhraseType == phrase.SimpleAssignmentExpression &&
				doc.peek(0).Type == lexer.Ampersand {
				doc.next(false) //&
				p.Type = phrase.ByRefAssignmentExpression
			}

			p.Children = append(p.Children, doc.expression(precedence))
		}

		lhs = doc.end()
	}

	return lhs
}

func (doc *Parser) ternaryExpression(testExpr phrase.AstNode) *phrase.Phrase {
	p := doc.start(phrase.TernaryExpression, true)
	p.Children = append(p.Children, testExpr)

	doc.next(false) //?

	if colonToken := doc.optional(lexer.Colon); colonToken != nil {
		p.Children = append(p.Children, doc.expression(0))
	} else {
		p.Children = append(p.Children, doc.expression(0))
		doc.expect(lexer.Colon)
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) variableOrExpression(precedence int) phrase.AstNode {
	part := doc.variableAtom(precedence)
	isVariable := false

	if p, ok := part.(*phrase.Phrase); ok {
		isVariable = p.Type == phrase.SimpleVariable
	}

	if isDereferenceOperator(doc.peek(0)) {
		part = doc.variable(part)
		isVariable = true
	} else {
		switch part.(*phrase.Phrase).Type {
		case phrase.QualifiedName,
			phrase.FullyQualifiedName,
			phrase.RelativeQualifiedName:
			part = doc.constantAccessExpression(part)
		}
	}

	if !isVariable {
		return part
	}

	//check for post increment/decrement
	t := doc.peek(0)
	if t.Type == lexer.PlusPlus {
		return doc.postfixExpression(phrase.PostfixIncrementExpression, part.(*phrase.Phrase))
	} else if t.Type == lexer.MinusMinus {
		return doc.postfixExpression(phrase.PostfixDecrementExpression, part.(*phrase.Phrase))
	} else {
		return part
	}
}

func (doc *Parser) constantAccessExpression(qName phrase.AstNode) *phrase.Phrase {
	p := doc.start(phrase.ConstantAccessExpression, true)
	p.Children = append(p.Children, qName)

	return doc.end()
}

func (doc *Parser) postfixExpression(
	phraseType phrase.PhraseType, variableNode *phrase.Phrase) *phrase.Phrase {
	p := doc.start(phraseType, true)
	p.Children = append(p.Children, variableNode)
	doc.next(false) //operator

	return doc.end()
}

func isDereferenceOperator(t *lexer.Token) bool {
	switch t.Type {
	case lexer.OpenBracket,
		lexer.OpenBrace,
		lexer.Arrow,
		lexer.OpenParenthesis,
		lexer.ColonColon:
		return true
	}

	return false
}

func (doc *Parser) expressionAtom(precedence int) phrase.AstNode {
	t := doc.peek(0)

	switch t.Type {
	case lexer.Static:
		if doc.peek(1).Type == lexer.Function {
			return doc.anonymousFunctionCreationExpression()
		}

		return doc.variableOrExpression(0)
	case lexer.StringLiteral:
		if isDereferenceOperator(doc.peek(1)) {
			return doc.variableOrExpression(0)
		}

		return doc.next(true)
	case lexer.VariableName,
		lexer.Dollar,
		lexer.Array,
		lexer.OpenBracket,
		lexer.Backslash,
		lexer.Name,
		lexer.Namespace,
		lexer.OpenParenthesis:
		return doc.variableOrExpression(precedence)
	case lexer.PlusPlus:
		return doc.unaryExpression(phrase.PrefixIncrementExpression)
	case lexer.MinusMinus:
		return doc.unaryExpression(phrase.PrefixDecrementExpression)
	case lexer.Plus,
		lexer.Minus,
		lexer.Exclamation,
		lexer.Tilde:
		return doc.unaryExpression(phrase.UnaryOpExpression)
	case lexer.AtSymbol:
		return doc.unaryExpression(phrase.ErrorControlExpression)
	case lexer.IntegerCast,
		lexer.FloatCast,
		lexer.StringCast,
		lexer.ArrayCast,
		lexer.ObjectCast,
		lexer.BooleanCast,
		lexer.UnsetCast:
		return doc.unaryExpression(phrase.CastExpression)
	case lexer.List:
		return doc.listIntrinsic()
	case lexer.Clone:
		return doc.cloneExpression()
	case lexer.New:
		return doc.objectCreationExpression()
	case lexer.FloatingLiteral,
		lexer.IntegerLiteral,
		lexer.LineConstant,
		lexer.FileConstant,
		lexer.DirectoryConstant,
		lexer.TraitConstant,
		lexer.MethodConstant,
		lexer.FunctionConstant,
		lexer.NamespaceConstant,
		lexer.ClassConstant:
		return doc.next(true)
	case lexer.StartHeredoc:
		return doc.heredocStringLiteral()
	case lexer.DoubleQuote:
		return doc.doubleQuotedStringLiteral()
	case lexer.Backtick:
		return doc.shellCommandExpression()
	case lexer.Print:
		return doc.printIntrinsic()
	case lexer.Yield:
		return doc.yieldExpression()
	case lexer.YieldFrom:
		return doc.yieldFromExpression()
	case lexer.Function:
		return doc.anonymousFunctionCreationExpression()
	case lexer.Include:
		return doc.scriptInclusion(phrase.IncludeExpression)
	case lexer.IncludeOnce:
		return doc.scriptInclusion(phrase.IncludeOnceExpression)
	case lexer.Require:
		return doc.scriptInclusion(phrase.RequireExpression)
	case lexer.RequireOnce:
		return doc.scriptInclusion(phrase.RequireOnceExpression)
	case lexer.Eval:
		return doc.evalIntrinsic()
	case lexer.Empty:
		return doc.emptyIntrinsic()
	case lexer.Exit:
		return doc.exitIntrinsic()
	case lexer.Isset:
		return doc.issetIntrinsic()
	default:
		//error
		doc.start(phrase.ErrorExpression, false)
		doc.error(lexer.Undefined)

		return doc.end()
	}
}

func (doc *Parser) exitIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.ExitIntrinsic, false)
	doc.next(false) //exit or die
	if t := doc.optional(lexer.OpenParenthesis); t != nil {
		if isExpressionStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.expression(0))
		}

		doc.expect(lexer.CloseParenthesis)
	}

	return doc.end()
}

func (doc *Parser) issetIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.IssetIntrinsic, false)
	doc.next(false) //isset
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.variableList([]lexer.TokenType{lexer.CloseParenthesis}))
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) emptyIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.EmptyIntrinsic, false)
	doc.next(false) //keyword
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) evalIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.EvalIntrinsic, false)
	doc.next(false) //keyword
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) scriptInclusion(phraseType phrase.PhraseType) *phrase.Phrase {
	p := doc.start(phraseType, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) printIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.PrintIntrinsic, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) yieldFromExpression() *phrase.Phrase {
	p := doc.start(phrase.YieldFromExpression, false)
	doc.next(false) //keyword
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) yieldExpression() *phrase.Phrase {
	p := doc.start(phrase.YieldExpression, false)
	doc.next(false) //yield
	if !isExpressionStart(doc.peek(0)) {
		return doc.end()
	}

	keyOrValue := doc.expression(0)
	p.Children = append(p.Children, keyOrValue)

	if doc.optional(lexer.FatArrow) != nil {
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) shellCommandExpression() *phrase.Phrase {
	p := doc.start(phrase.ShellCommandExpression, false)
	doc.next(false) //`
	p.Children = append(p.Children, doc.encapsulatedVariableList(lexer.Backtick))
	doc.expect(lexer.Backtick)

	return doc.end()
}

func (doc *Parser) doubleQuotedStringLiteral() *phrase.Phrase {
	p := doc.start(phrase.DoubleQuotedStringLiteral, false)
	doc.next(false) //"
	p.Children = append(p.Children, doc.encapsulatedVariableList(lexer.DoubleQuote))
	doc.expect(lexer.DoubleQuote)

	return doc.end()
}

func (doc *Parser) encapsulatedVariableList(breakOn lexer.TokenType) *phrase.Phrase {
	return doc.list(
		phrase.EncapsulatedVariableList,
		doc.encapsulatedVariable,
		isEncapsulatedVariableStart,
		[]lexer.TokenType{breakOn},
		encapsulatedVariableListRecoverSet)
}

func isEncapsulatedVariableStart(t *lexer.Token) bool {

	switch t.Type {
	case lexer.EncapsulatedAndWhitespace,
		lexer.VariableName,
		lexer.DollarCurlyOpen,
		lexer.CurlyOpen:
		return true
	}

	return false
}

func (doc *Parser) encapsulatedVariable() phrase.AstNode {

	switch doc.peek(0).Type {
	case lexer.EncapsulatedAndWhitespace:
		return doc.next(true)
	case lexer.VariableName:
		t := doc.peek(1)
		if t.Type == lexer.OpenBracket {
			return doc.encapsulatedDimension()
		} else if t.Type == lexer.Arrow {
			return doc.encapsulatedProperty()
		}

		return doc.simpleVariable()
	case lexer.DollarCurlyOpen:
		return doc.dollarCurlyOpenEncapsulatedVariable()
	case lexer.CurlyOpen:
		return doc.curlyOpenEncapsulatedVariable()
	}

	t := doc.peek(0)

	panic(errors.New("Unexpected token: " + t.Type.String()))
}

func (doc *Parser) curlyOpenEncapsulatedVariable() *phrase.Phrase {
	p := doc.start(phrase.EncapsulatedVariable, false)
	doc.next(false) //{
	p.Children = append(p.Children, doc.variable(doc.variableAtom(0)))
	doc.expect(lexer.CloseBrace)

	return doc.end()
}

func (doc *Parser) dollarCurlyOpenEncapsulatedVariable() *phrase.Phrase {
	p := doc.start(phrase.EncapsulatedVariable, false)
	doc.next(false) //${
	t := doc.peek(0)

	if t.Type == lexer.VariableName {
		if doc.peek(1).Type == lexer.OpenBracket {
			p.Children = append(p.Children, *doc.dollarCurlyEncapsulatedDimension())
		} else {
			doc.start(phrase.SimpleVariable, false)
			doc.next(false)
			p.Children = append(p.Children, doc.end())
		}
	} else if isExpressionStart(t) {
		p.Children = append(p.Children, doc.expression(0))
	} else {
		doc.error(lexer.Undefined)
	}

	doc.expect(lexer.CloseBrace)

	return doc.end()
}

func (doc *Parser) dollarCurlyEncapsulatedDimension() *phrase.Phrase {
	p := doc.start(phrase.SubscriptExpression, false)
	doc.next(false) //VariableName
	doc.next(false) // [
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseBracket)

	return doc.end()
}

func (doc *Parser) encapsulatedDimension() *phrase.Phrase {
	p := doc.start(phrase.SubscriptExpression, false)

	p.Children = append(p.Children, doc.simpleVariable()) //T_VARIABLE
	doc.next(false)                                       //[
	switch doc.peek(0).Type {
	case lexer.Name,
		lexer.IntegerLiteral:
		doc.next(false)
	case lexer.VariableName:
		p.Children = append(p.Children, doc.simpleVariable())
	case lexer.Minus:
		doc.start(phrase.UnaryOpExpression, false)
		doc.next(false) //-
		doc.expect(lexer.IntegerLiteral)
		p.Children = append(p.Children, doc.end())
	default:
		//error
		doc.error(lexer.Undefined)
	}

	doc.expect(lexer.CloseBracket)
	return doc.end()
}

func (doc *Parser) encapsulatedProperty() *phrase.Phrase {
	p := doc.start(phrase.PropertyAccessExpression, false)
	p.Children = append(p.Children, doc.simpleVariable())
	doc.next(false) //->
	doc.expect(lexer.Name)

	return doc.end()
}

func (doc *Parser) heredocStringLiteral() *phrase.Phrase {
	p := doc.start(phrase.HeredocStringLiteral, false)
	doc.next(false) //StartHeredoc
	p.Children = append(p.Children, doc.encapsulatedVariableList(lexer.EndHeredoc))
	doc.expect(lexer.EndHeredoc)

	return doc.end()
}

func (doc *Parser) anonymousClassDeclaration() *phrase.Phrase {
	p := doc.start(phrase.AnonymousClassDeclaration, false)
	p.Children = append(p.Children, doc.anonymousClassDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(phrase.ClassDeclarationBody,
		isClassMemberStart, doc.classMemberDeclarationList))

	return doc.end()
}

func (doc *Parser) anonymousClassDeclarationHeader() *phrase.Phrase {
	p := doc.start(phrase.AnonymousClassDeclarationHeader, false)
	doc.next(false) //class
	if doc.optional(lexer.OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}
		doc.expect(lexer.CloseParenthesis)
	}

	if doc.peek(0).Type == lexer.Extends {
		p.Children = append(p.Children, doc.classBaseClause())
	}

	if doc.peek(0).Type == lexer.Implements {
		p.Children = append(p.Children, doc.classInterfaceClause())
	}

	return doc.end()
}

func (doc *Parser) classInterfaceClause() *phrase.Phrase {
	p := doc.start(phrase.ClassInterfaceClause, false)
	doc.next(false) //implements
	p.Children = append(p.Children, doc.qualifiedNameList([]lexer.TokenType{lexer.OpenBrace}))

	return doc.end()
}

func (doc *Parser) classMemberDeclarationList() *phrase.Phrase {
	return doc.list(
		phrase.ClassMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]lexer.TokenType{lexer.CloseBrace},
		classMemberDeclarationListRecoverSet)
}

func isClassMemberStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Public,
		lexer.Protected,
		lexer.Private,
		lexer.Static,
		lexer.Abstract,
		lexer.Final,
		lexer.Function,
		lexer.Var,
		lexer.Const,
		lexer.Use,
		lexer.DocumentCommentStart:
		return true
	}

	return false
}

func (doc *Parser) classMemberDeclaration() phrase.AstNode {
	p := doc.start(phrase.ErrorClassMemberDeclaration, false)
	t := doc.peek(0)

	switch t.Type {
	case lexer.Public,
		lexer.Protected,
		lexer.Private,
		lexer.Static,
		lexer.Abstract,
		lexer.Final:
		modifiers := doc.memberModifierList()
		t = doc.peek(0)
		if t.Type == lexer.VariableName {
			p.Children = append(p.Children, modifiers)

			return doc.propertyDeclaration(p)
		} else if t.Type == lexer.Function {
			return doc.methodDeclaration(p, modifiers)
		} else if t.Type == lexer.Const {
			p.Children = append(p.Children, modifiers)

			return doc.classConstDeclaration(p)
		}

		//error
		p.Children = append(p.Children, modifiers)
		doc.error(lexer.Undefined)

		return doc.end()
	case lexer.Function:
		return doc.methodDeclaration(p, nil)
	case lexer.Var:
		doc.next(false)

		return doc.propertyDeclaration(p)
	case lexer.Const:
		return doc.classConstDeclaration(p)
	case lexer.Use:
		return doc.traitUseClause(p)
	case lexer.DocumentCommentStart:
		doc.end()
		return doc.docComment()
	}

	panic(errors.New("Unexpected token: " + t.Type.String()))
}

func (doc *Parser) traitUseClause(p *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.TraitUseClause
	doc.next(false) //use
	p.Children = append(p.Children,
		doc.qualifiedNameList([]lexer.TokenType{lexer.Semicolon, lexer.OpenBrace}),
		doc.traitUseSpecification())

	return doc.end()

}

func (doc *Parser) traitUseSpecification() *phrase.Phrase {
	p := doc.start(phrase.TraitUseSpecification, false)
	t := doc.expectOneOf([]lexer.TokenType{lexer.Semicolon, lexer.OpenBrace})

	if t != nil && t.Type == lexer.OpenBrace {
		if isTraitAdaptationStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.traitAdaptationList())
		}
		doc.expect(lexer.CloseBrace)
	}

	return doc.end()
}

func (doc *Parser) traitAdaptationList() *phrase.Phrase {
	return doc.list(
		phrase.TraitAdaptationList,
		doc.traitAdaptation,
		isTraitAdaptationStart,
		[]lexer.TokenType{lexer.CloseBrace}, nil)
}

func isTraitAdaptationStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Name,
		lexer.Backslash,
		lexer.Namespace:
		return true
	}

	return isSemiReservedToken(t)
}

func (doc *Parser) traitAdaptation() phrase.AstNode {
	p := doc.start(phrase.ErrorTraitAdaptation, false)
	t := doc.peek(0)
	t2 := doc.peek(1)

	if t.Type == lexer.Namespace ||
		t.Type == lexer.Backslash ||
		(t.Type == lexer.Name &&
			(t2.Type == lexer.ColonColon || t2.Type == lexer.Backslash)) {

		p.Children = append(p.Children, doc.methodReference())

		if doc.peek(0).Type == lexer.InsteadOf {
			doc.next(false)
			return doc.traitPrecedence(p)
		}

	} else if t.Type == lexer.Name || isSemiReservedToken(t) {
		methodRef := doc.start(phrase.MethodReference, false)
		methodRef.Children = append(methodRef.Children, doc.identifier())
		p.Children = append(p.Children, doc.end())
	} else {
		//error
		doc.error(lexer.Undefined)

		return doc.end()
	}

	return doc.traitAlias(p)
}

func (doc *Parser) traitAlias(p *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.TraitAlias
	doc.expect(lexer.As)

	t := doc.peek(0)

	if t.Type == lexer.Name || isReservedToken(t) {
		p.Children = append(p.Children, doc.identifier())
	} else if isMemberModifier(t) {
		doc.next(false)
		t = doc.peek(0)
		if t.Type == lexer.Name || isSemiReservedToken(t) {
			p.Children = append(p.Children, doc.identifier())
		}
	} else {
		doc.error(lexer.Undefined)
	}

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) traitPrecedence(p *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.TraitPrecedence
	p.Children = append(p.Children, doc.qualifiedNameList([]lexer.TokenType{lexer.Semicolon}))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) methodReference() *phrase.Phrase {
	p := doc.start(phrase.MethodReference, false)
	p.Children = append(p.Children, doc.qualifiedName())
	doc.expect(lexer.ColonColon)
	p.Children = append(p.Children, doc.identifier())

	return doc.end()
}

func (doc *Parser) methodDeclarationHeader(memberModifers *phrase.Phrase) *phrase.Phrase {
	p := doc.start(phrase.MethodDeclarationHeader, true)
	if memberModifers != nil {
		p.Children = append(p.Children, memberModifers)
	}
	doc.next(false) //function
	doc.optional(lexer.Ampersand)
	p.Children = append(p.Children, doc.identifier())
	doc.expect(lexer.OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			phrase.ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			lexer.Comma,
			[]lexer.TokenType{lexer.CloseParenthesis}, false))
	}

	doc.expect(lexer.CloseParenthesis)

	if doc.peek(0).Type == lexer.Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func (doc *Parser) methodDeclaration(p *phrase.Phrase, memberModifers *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.MethodDeclaration
	p.Children = append(p.Children, doc.methodDeclarationHeader(memberModifers))
	p.Children = append(p.Children, doc.methodDeclarationBody())

	return doc.end()
}

func (doc *Parser) methodDeclarationBody() *phrase.Phrase {
	p := doc.start(phrase.MethodDeclarationBody, false)

	if doc.peek(0).Type == lexer.Semicolon {
		doc.next(false)
	} else {
		p.Children = append(p.Children, doc.compoundStatement())
	}

	return doc.end()
}

func (doc *Parser) identifier() *phrase.Phrase {
	doc.start(phrase.Identifier, false)
	t := doc.peek(0)
	if t.Type == lexer.Name || isSemiReservedToken(t) {
		doc.next(false)
	} else {
		doc.error(lexer.Undefined)
	}

	return doc.end()
}

func (doc *Parser) interfaceDeclaration() *phrase.Phrase {
	p := doc.start(phrase.InterfaceDeclaration, false)
	p.Children = append(p.Children, doc.interfaceDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		phrase.InterfaceDeclarationBody, isClassMemberStart, doc.interfaceMemberDeclarations))

	return doc.end()
}

func (doc *Parser) typeDeclarationBody(
	phraseType phrase.PhraseType,
	elementStartPredicate func(*lexer.Token) bool,
	listFunction func() *phrase.Phrase) *phrase.Phrase {
	p := doc.start(phraseType, false)
	doc.expect(lexer.OpenBrace)

	if elementStartPredicate(doc.peek(0)) {
		p.Children = append(p.Children, listFunction())
	}

	doc.expect(lexer.CloseBrace)

	return doc.end()
}

func (doc *Parser) interfaceMemberDeclarations() *phrase.Phrase {
	return doc.list(
		phrase.InterfaceMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]lexer.TokenType{lexer.CloseBrace},
		classMemberDeclarationListRecoverSet)
}

func (doc *Parser) interfaceDeclarationHeader() *phrase.Phrase {
	p := doc.start(phrase.InterfaceDeclarationHeader, false)
	doc.next(false) //interface
	doc.expect(lexer.Name)

	if doc.peek(0).Type == lexer.Extends {
		p.Children = append(p.Children, doc.interfaceBaseClause())
	}

	return doc.end()

}

func (doc *Parser) interfaceBaseClause() *phrase.Phrase {
	p := doc.start(phrase.InterfaceBaseClause, false)
	doc.next(false) //extends
	p.Children = append(p.Children, doc.qualifiedNameList([]lexer.TokenType{lexer.OpenBrace}))

	return doc.end()
}

func (doc *Parser) traitDeclaration() *phrase.Phrase {
	p := doc.start(phrase.TraitDeclaration, false)
	p.Children = append(p.Children, doc.traitDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		phrase.TraitDeclarationBody, isClassMemberStart, doc.traitMemberDeclarations))

	return doc.end()
}

func (doc *Parser) traitDeclarationHeader() *phrase.Phrase {
	doc.start(phrase.TraitDeclarationHeader, false)
	doc.next(false) //trait
	doc.expect(lexer.Name)

	return doc.end()
}

func (doc *Parser) traitMemberDeclarations() *phrase.Phrase {
	return doc.list(
		phrase.TraitMemberDeclarationList,
		doc.classMemberDeclaration,
		isClassMemberStart,
		[]lexer.TokenType{lexer.CloseBrace},
		classMemberDeclarationListRecoverSet[:])
}

func (doc *Parser) functionDeclaration() *phrase.Phrase {
	p := doc.start(phrase.FunctionDeclaration, false)
	p.Children = append(p.Children, doc.functionDeclarationHeader())
	p.Children = append(p.Children, doc.functionDeclarationBody())

	return doc.end()
}

func (doc *Parser) functionDeclarationBody() *phrase.Phrase {
	cs := doc.compoundStatement()
	cs.Type = phrase.FunctionDeclarationBody

	return cs
}

func (doc *Parser) functionDeclarationHeader() *phrase.Phrase {
	p := doc.start(phrase.FunctionDeclarationHeader, false)

	doc.next(false) //function
	doc.optional(lexer.Ampersand)
	doc.expect(lexer.Name)
	doc.expect(lexer.OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			phrase.ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			lexer.Comma,
			[]lexer.TokenType{lexer.CloseParenthesis}, false))
	}

	doc.expect(lexer.CloseParenthesis)

	if doc.peek(0).Type == lexer.Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func isParameterStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Ampersand,
		lexer.Ellipsis,
		lexer.VariableName:
		return true
	default:
		return isTypeDeclarationStart(t)
	}
}

func (doc *Parser) classDeclaration() *phrase.Phrase {
	p := doc.start(phrase.ClassDeclaration, false)

	p.Children = append(p.Children, doc.classDeclarationHeader())
	p.Children = append(p.Children, doc.typeDeclarationBody(
		phrase.ClassDeclarationBody, isClassMemberStart, doc.classMemberDeclarationList))

	return doc.end()

}

func (doc *Parser) classDeclarationHeader() *phrase.Phrase {
	p := doc.start(phrase.ClassDeclarationHeader, false)
	doc.optionalOneOf([]lexer.TokenType{lexer.Abstract, lexer.Final})
	doc.expect(lexer.Class)
	doc.expect(lexer.Name)

	if doc.peek(0).Type == lexer.Extends {
		p.Children = append(p.Children, doc.classBaseClause())
	}

	if doc.peek(0).Type == lexer.Implements {
		p.Children = append(p.Children, doc.classInterfaceClause())
	}

	return doc.end()
}

func (doc *Parser) classBaseClause() *phrase.Phrase {
	p := doc.start(phrase.ClassBaseClause, false)
	doc.next(false) //extends
	p.Children = append(p.Children, doc.qualifiedName())

	return doc.end()
}

func (doc *Parser) compoundStatement() *phrase.Phrase {
	p := doc.start(phrase.CompoundStatement, false)
	doc.expect(lexer.OpenBrace)

	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.CloseBrace}))
	}

	doc.expect(lexer.CloseBrace)

	return doc.end()
}

func (doc *Parser) statement() phrase.AstNode {
	t := doc.peek(0)

	switch t.Type {
	case lexer.Namespace:
		return doc.namespaceDefinition()
	case lexer.Use:
		return doc.namespaceUseDeclaration()
	case lexer.HaltCompiler:
		return doc.haltCompilerStatement()
	case lexer.Const:
		return doc.constDeclaration()
	case lexer.Function:
		{
			p1 := doc.peek(1)
			if p1.Type == lexer.OpenParenthesis ||
				(p1.Type == lexer.Ampersand && doc.peek(2).Type == lexer.OpenParenthesis) {
				//anon fn without assignment
				return doc.expressionStatement()
			} else {
				return doc.functionDeclaration()
			}
		}
	case lexer.Class,
		lexer.Abstract,
		lexer.Final:
		return doc.classDeclaration()
	case lexer.Trait:
		return doc.traitDeclaration()
	case lexer.Interface:
		return doc.interfaceDeclaration()
	case lexer.OpenBrace:
		return doc.compoundStatement()
	case lexer.If:
		return doc.ifStatement()
	case lexer.While:
		return doc.whileStatement()
	case lexer.Do:
		return doc.doStatement()
	case lexer.For:
		return doc.forStatement()
	case lexer.Switch:
		return doc.switchStatement()
	case lexer.Break:
		return doc.breakStatement()
	case lexer.Continue:
		return doc.continueStatement()
	case lexer.Return:
		return doc.returnStatement()
	case lexer.Global:
		return doc.globalDeclaration()
	case lexer.Static:
		if doc.peek(1).Type == lexer.VariableName &&
			tokenTypeIndexOf(
				[]lexer.TokenType{lexer.Semicolon, lexer.Comma, lexer.CloseTag, lexer.Equals},
				doc.peek(2).Type) >= 0 {
			return doc.functionStaticDeclaration()
		} else {
			return doc.expressionStatement()
		}
	case lexer.Text,
		lexer.OpenTag,
		lexer.CloseTag:
		return doc.inlineText()
	case lexer.ForEach:
		return doc.foreachStatement()
	case lexer.Declare:
		return doc.declareStatement()
	case lexer.Try:
		return doc.tryStatement()
	case lexer.Throw:
		return doc.throwStatement()
	case lexer.Goto:
		return doc.gotoStatement()
	case lexer.Echo,
		lexer.OpenTagEcho:
		return doc.echoIntrinsic()
	case lexer.Unset:
		return doc.unsetIntrinsic()
	case lexer.Semicolon:
		return doc.nullStatement()
	case lexer.DocumentCommentStart:
		return doc.docComment()
	case lexer.Name:
		if doc.peek(1).Type == lexer.Colon {
			return doc.namedLabelStatement()
		}
		fallthrough
	default:
		return doc.expressionStatement()
	}
}

func (doc *Parser) inlineText() *phrase.Phrase {
	doc.start(phrase.InlineText, false)

	doc.optional(lexer.CloseTag)
	doc.optional(lexer.Text)
	doc.optional(lexer.OpenTag)

	return doc.end()
}

func (doc *Parser) nullStatement() *phrase.Phrase {
	doc.start(phrase.NullStatement, false)
	doc.next(false) //;

	return doc.end()
}

func (doc *Parser) tryStatement() *phrase.Phrase {
	p := doc.start(phrase.TryStatement, false)
	doc.next(false) //try
	p.Children = append(p.Children, doc.compoundStatement())

	t := doc.peek(0)

	if t.Type == lexer.Catch {
		p.Children = append(p.Children, doc.list(
			phrase.CatchClauseList,
			doc.catchClause,
			func(t *lexer.Token) bool { return t.Type == lexer.Catch }, nil, nil))
	} else if t.Type != lexer.Finally {
		doc.error(lexer.Undefined)
	}

	if doc.peek(0).Type == lexer.Finally {
		p.Children = append(p.Children, doc.finallyClause())
	}

	return doc.end()

}

func (doc *Parser) finallyClause() *phrase.Phrase {
	p := doc.start(phrase.FinallyClause, false)
	doc.next(false) //finally
	p.Children = append(p.Children, doc.compoundStatement())

	return doc.end()
}

func (doc *Parser) catchClause() phrase.AstNode {
	p := doc.start(phrase.CatchClause, false)
	doc.next(false) //catch
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.delimitedList(
		phrase.CatchNameList,
		doc.qualifiedName,
		isQualifiedNameStart,
		lexer.Bar,
		[]lexer.TokenType{lexer.VariableName}, false))
	doc.expect(lexer.VariableName)
	doc.expect(lexer.CloseParenthesis)
	p.Children = append(p.Children, doc.compoundStatement())

	return doc.end()
}

func (doc *Parser) declareDirective() *phrase.Phrase {
	doc.start(phrase.DeclareDirective, false)
	doc.expect(lexer.Name)
	doc.expect(lexer.Equals)
	doc.expectOneOf([]lexer.TokenType{
		lexer.IntegerLiteral, lexer.FloatingLiteral, lexer.StringLiteral})

	return doc.end()
}

func (doc *Parser) declareStatement() *phrase.Phrase {
	p := doc.start(phrase.DeclareStatement, false)
	doc.next(false) //declare
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.declareDirective())
	doc.expect(lexer.CloseParenthesis)

	t := doc.peek(0)

	if t.Type == lexer.Colon {
		doc.next(false) //:
		p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.EndDeclare}))
		doc.expect(lexer.EndDeclare)
		doc.expect(lexer.Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else if t.Type == lexer.Semicolon {
		doc.next(false)
	} else {
		doc.error(lexer.Undefined)
	}

	return doc.end()
}

func (doc *Parser) switchStatement() *phrase.Phrase {
	p := doc.start(phrase.SwitchStatement, false)
	doc.next(false) //switch
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)

	t := doc.expectOneOf([]lexer.TokenType{lexer.Colon, lexer.OpenBrace})
	tCase := doc.peek(0)

	if tCase.Type == lexer.Case || tCase.Type == lexer.Default {
		caseTokenType := lexer.CloseBrace
		if t != nil && t.Type == lexer.Colon {
			caseTokenType = lexer.EndSwitch
		}

		p.Children = append(p.Children, doc.caseStatements(caseTokenType))
	}

	if t != nil && t.Type == lexer.Colon {
		doc.expect(lexer.EndSwitch)
		doc.expect(lexer.Semicolon)
	} else {
		doc.expect(lexer.CloseBrace)
	}

	return doc.end()

}

func (doc *Parser) caseStatements(breakOn lexer.TokenType) *phrase.Phrase {
	p := doc.start(phrase.CaseStatementList, false)
	var t *lexer.Token
	caseBreakOn := []lexer.TokenType{lexer.Case, lexer.Default, breakOn}

	for {
		t = doc.peek(0)

		if t.Type == lexer.Case {
			p.Children = append(p.Children, doc.caseStatement(caseBreakOn))
		} else if t.Type == lexer.Default {
			p.Children = append(p.Children, doc.defaultStatement(caseBreakOn))
		} else if breakOn == t.Type {
			break
		} else {
			doc.error(lexer.Undefined)
			break
		}

	}

	return doc.end()
}

func (doc *Parser) caseStatement(breakOn []lexer.TokenType) *phrase.Phrase {
	p := doc.start(phrase.CaseStatement, false)
	doc.next(false) //case
	p.Children = append(p.Children, doc.expression(0))
	doc.expectOneOf([]lexer.TokenType{lexer.Colon, lexer.Semicolon})
	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList(breakOn))
	}

	return doc.end()
}

func (doc *Parser) defaultStatement(breakOn []lexer.TokenType) *phrase.Phrase {
	p := doc.start(phrase.DefaultStatement, false)
	doc.next(false) //default
	doc.expectOneOf([]lexer.TokenType{lexer.Colon, lexer.Semicolon})
	if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statementList(breakOn))
	}

	return doc.end()
}

func (doc *Parser) namedLabelStatement() *phrase.Phrase {
	doc.start(phrase.NamedLabelStatement, false)
	doc.next(false) //name
	doc.next(false) //:

	return doc.end()
}

func (doc *Parser) gotoStatement() *phrase.Phrase {
	doc.start(phrase.GotoStatement, false)
	doc.next(false) //goto
	doc.expect(lexer.Name)
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) throwStatement() *phrase.Phrase {
	p := doc.start(phrase.ThrowStatement, false)
	doc.next(false) //throw
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) foreachCollection() *phrase.Phrase {
	p := doc.start(phrase.ForeachCollection, false)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) foreachKeyOrValue() *phrase.Phrase {
	p := doc.start(phrase.ForeachValue, false)
	p.Children = append(p.Children, doc.expression(0))
	if doc.peek(0).Type == lexer.FatArrow {
		doc.next(false)
		p.Type = phrase.ForeachKey
	}

	return doc.end()
}

func (doc *Parser) foreachValue() *phrase.Phrase {
	p := doc.start(phrase.ForeachValue, false)
	doc.optional(lexer.Ampersand)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) foreachStatement() *phrase.Phrase {
	p := doc.start(phrase.ForeachStatement, false)
	doc.next(false) //foreach
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.foreachCollection())
	doc.expect(lexer.As)
	var keyOrValue *phrase.Phrase
	if doc.peek(0).Type == lexer.Ampersand {
		keyOrValue = doc.foreachValue()
	} else {
		keyOrValue = doc.foreachKeyOrValue()
	}
	p.Children = append(p.Children, keyOrValue)

	if keyOrValue.Type == phrase.ForeachKey {
		p.Children = append(p.Children, doc.foreachValue())
	}

	doc.expect(lexer.CloseParenthesis)

	t := doc.peek(0)

	if t.Type == lexer.Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.EndForeach}))
		doc.expect(lexer.EndForeach)
		doc.expect(lexer.Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(lexer.Undefined)
	}

	return doc.end()

}

func isVariableStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.VariableName,
		lexer.Dollar,
		lexer.OpenParenthesis,
		lexer.Array,
		lexer.OpenBracket,
		lexer.StringLiteral,
		lexer.Static,
		lexer.Name,
		lexer.Namespace,
		lexer.Backslash:
		return true
	}

	return false
}

func (doc *Parser) variableInitial() phrase.AstNode {
	return doc.variable(doc.variableAtom(0))
}

func (doc *Parser) variableList(breakOn []lexer.TokenType) *phrase.Phrase {
	return doc.delimitedList(
		phrase.VariableList,
		doc.variableInitial,
		isVariableStart,
		lexer.Comma,
		breakOn,
		false)
}

func (doc *Parser) unsetIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.UnsetIntrinsic, false)
	doc.next(false) //unset
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.variableList([]lexer.TokenType{lexer.CloseParenthesis}))
	doc.expect(lexer.CloseParenthesis)
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) expressionInitial() phrase.AstNode {
	return doc.expression(0)
}

func (doc *Parser) echoIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.EchoIntrinsic, false)
	doc.next(false) //echo or <?=
	p.Children = append(p.Children, doc.delimitedList(
		phrase.ExpressionList,
		doc.expressionInitial,
		isExpressionStart,
		lexer.Comma,
		nil,
		false))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) functionStaticDeclaration() *phrase.Phrase {
	p := doc.start(phrase.FunctionStaticDeclaration, false)
	doc.next(false) //static
	p.Children = append(p.Children, doc.delimitedList(
		phrase.StaticVariableDeclarationList,
		doc.staticVariableDeclaration,
		func(t *lexer.Token) bool { return t.Type == lexer.VariableName },
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		false))

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) globalDeclaration() *phrase.Phrase {
	p := doc.start(phrase.GlobalDeclaration, false)
	doc.next(false) //global
	p.Children = append(p.Children, doc.delimitedList(
		phrase.VariableNameList,
		doc.simpleVariable,
		isSimpleVariableStart,
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		false))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func isSimpleVariableStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.VariableName, lexer.Dollar:
		return true
	}

	return false
}

func (doc *Parser) staticVariableDeclaration() phrase.AstNode {
	p := doc.start(phrase.StaticVariableDeclaration, false)
	doc.expect(lexer.VariableName)

	if doc.peek(0).Type == lexer.Equals {
		p.Children = append(p.Children, doc.functionStaticInitialiser())
	}

	return doc.end()
}

func (doc *Parser) functionStaticInitialiser() *phrase.Phrase {
	p := doc.start(phrase.FunctionStaticInitialiser, false)
	doc.next(false) //=
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) continueStatement() *phrase.Phrase {
	p := doc.start(phrase.ContinueStatement, false)
	doc.next(false) //break/continue
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) breakStatement() *phrase.Phrase {
	p := doc.start(phrase.BreakStatement, false)
	doc.next(false) //break/continue
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) returnStatement() *phrase.Phrase {
	p := doc.start(phrase.ReturnStatement, false)
	doc.next(false) //return
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) forExpressionGroup(
	phraseType phrase.PhraseType, breakOn []lexer.TokenType) *phrase.Phrase {

	return doc.delimitedList(
		phraseType,
		doc.expressionInitial,
		isExpressionStart,
		lexer.Comma,
		breakOn,
		false)
}

func (doc *Parser) forStatement() *phrase.Phrase {
	p := doc.start(phrase.ForStatement, false)
	doc.next(false) //for
	doc.expect(lexer.OpenParenthesis)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			phrase.ForInitialiser, []lexer.TokenType{lexer.Semicolon}))
	}

	doc.expect(lexer.Semicolon)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			phrase.ForControl, []lexer.TokenType{lexer.Semicolon}))
	}

	doc.expect(lexer.Semicolon)

	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.forExpressionGroup(
			phrase.ForEndOfLoop, []lexer.TokenType{lexer.CloseParenthesis}))
	}

	doc.expect(lexer.CloseParenthesis)

	t := doc.peek(0)

	if t.Type == lexer.Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.EndFor}))
		doc.expect(lexer.EndFor)
		doc.expect(lexer.Semicolon)
	} else if isStatementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(lexer.Undefined)
	}

	return doc.end()

}

func (doc *Parser) doStatement() *phrase.Phrase {
	p := doc.start(phrase.DoStatement, false)
	doc.next(false) // do
	p.Children = append(p.Children, doc.statement())
	doc.expect(lexer.While)
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) whileStatement() *phrase.Phrase {
	p := doc.start(phrase.WhileStatement, false)
	doc.next(false) //while
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)

	t := doc.peek(0)

	if t.Type == lexer.Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.EndWhile}))
		doc.expect(lexer.EndWhile)
		doc.expect(lexer.Semicolon)
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		//error
		doc.error(lexer.Undefined)
	}

	return doc.end()
}

func (doc *Parser) elseIfClause1() phrase.AstNode {
	p := doc.start(phrase.ElseIfClause, false)
	doc.next(false) //elseif
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)
	p.Children = append(p.Children, doc.statement())

	return doc.end()
}

func (doc *Parser) elseIfClause2() phrase.AstNode {
	p := doc.start(phrase.ElseIfClause, false)
	doc.next(false) //elseif
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)
	doc.expect(lexer.Colon)
	p.Children = append(p.Children, doc.statementList(
		[]lexer.TokenType{lexer.EndIf, lexer.Else, lexer.ElseIf}))

	return doc.end()
}

func (doc *Parser) elseClause1() *phrase.Phrase {
	p := doc.start(phrase.ElseClause, false)
	doc.next(false) //else
	p.Children = append(p.Children, doc.statement())

	return doc.end()
}

func (doc *Parser) elseClause2() *phrase.Phrase {
	p := doc.start(phrase.ElseClause, false)
	doc.next(false) //else
	doc.expect(lexer.Colon)
	p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.EndIf}))

	return doc.end()
}

func isElseIfClauseStart(t *lexer.Token) bool {
	return t.Type == lexer.ElseIf
}

func (doc *Parser) ifStatement() *phrase.Phrase {
	p := doc.start(phrase.IfStatement, false)
	doc.next(false) //if
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.CloseParenthesis)

	t := doc.peek(0)
	elseIfClausefunc := doc.elseIfClause1
	elseClausefunc := doc.elseClause1
	expectEndIf := false

	if t.Type == lexer.Colon {
		doc.next(false)
		p.Children = append(p.Children, doc.statementList(
			[]lexer.TokenType{lexer.ElseIf, lexer.Else, lexer.EndIf}))
		elseIfClausefunc = doc.elseIfClause2
		elseClausefunc = doc.elseClause2
		expectEndIf = true
	} else if isStatementStart(t) {
		p.Children = append(p.Children, doc.statement())
	} else {
		doc.error(lexer.Undefined)
	}

	if doc.peek(0).Type == lexer.ElseIf {
		p.Children = append(p.Children, doc.list(
			phrase.ElseIfClauseList,
			elseIfClausefunc,
			isElseIfClauseStart,
			nil,
			nil))
	}

	if doc.peek(0).Type == lexer.Else {
		p.Children = append(p.Children, elseClausefunc())
	}

	if expectEndIf {
		doc.expect(lexer.EndIf)
		doc.expect(lexer.Semicolon)
	}

	return doc.end()

}

func (doc *Parser) expressionStatement() *phrase.Phrase {
	p := doc.start(phrase.ExpressionStatement, false)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) returnType() *phrase.Phrase {
	p := doc.start(phrase.ReturnType, false)
	doc.next(false) //:
	p.Children = append(p.Children, doc.typeDeclaration())

	return doc.end()
}

func (doc *Parser) typeDeclaration() *phrase.Phrase {
	p := doc.start(phrase.TypeDeclaration, false)
	doc.optional(lexer.Question)

	switch doc.peek(0).Type {
	case lexer.Callable, lexer.Array:
		doc.next(false)
	case lexer.Name, lexer.Namespace, lexer.Backslash:
		p.Children = append(p.Children, doc.qualifiedName())
	default:
		doc.error(lexer.Undefined)
	}

	return doc.end()

}

func (doc *Parser) classConstDeclaration(p *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.ClassConstDeclaration
	doc.next(false) //const
	p.Children = append(p.Children, doc.delimitedList(
		phrase.ClassConstElementList,
		doc.classConstElement,
		isClassConstElementStartToken,
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		false))

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func isExpressionStart(t *lexer.Token) bool {

	switch t.Type {
	case lexer.VariableName,
		lexer.Dollar,
		lexer.Array,
		lexer.OpenBracket,
		lexer.StringLiteral,
		lexer.Backslash,
		lexer.Name,
		lexer.Namespace,
		lexer.OpenParenthesis,
		lexer.Static,
		lexer.PlusPlus,
		lexer.MinusMinus,
		lexer.Plus,
		lexer.Minus,
		lexer.Exclamation,
		lexer.Tilde,
		lexer.AtSymbol,
		lexer.IntegerCast,
		lexer.FloatCast,
		lexer.StringCast,
		lexer.ArrayCast,
		lexer.ObjectCast,
		lexer.BooleanCast,
		lexer.UnsetCast,
		lexer.List,
		lexer.Clone,
		lexer.New,
		lexer.FloatingLiteral,
		lexer.IntegerLiteral,
		lexer.LineConstant,
		lexer.FileConstant,
		lexer.DirectoryConstant,
		lexer.TraitConstant,
		lexer.MethodConstant,
		lexer.FunctionConstant,
		lexer.NamespaceConstant,
		lexer.ClassConstant,
		lexer.StartHeredoc,
		lexer.DoubleQuote,
		lexer.Backtick,
		lexer.Print,
		lexer.Yield,
		lexer.YieldFrom,
		lexer.Function,
		lexer.Include,
		lexer.IncludeOnce,
		lexer.Require,
		lexer.RequireOnce,
		lexer.Eval,
		lexer.Empty,
		lexer.Isset,
		lexer.Exit:
		return true
	}

	return false
}

func (doc *Parser) classConstElement() phrase.AstNode {
	p := doc.start(phrase.ClassConstElement, false)
	p.Children = append(p.Children, doc.identifier())
	doc.expect(lexer.Equals)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func isPropertyElementStart(t *lexer.Token) bool {
	return t.Type == lexer.VariableName
}

func (doc *Parser) propertyDeclaration(p *phrase.Phrase) *phrase.Phrase {
	p.Type = phrase.PropertyDeclaration
	p.Children = append(p.Children, doc.delimitedList(
		phrase.PropertyElementList,
		doc.propertyElement,
		isPropertyElementStart,
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		false))
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) propertyElement() phrase.AstNode {
	p := doc.start(phrase.PropertyElement, false)
	doc.expect(lexer.VariableName)

	if doc.peek(0).Type == lexer.Equals {
		p.Children = append(p.Children, doc.propertyInitialiser())
	}

	return doc.end()
}

func (doc *Parser) propertyInitialiser() *phrase.Phrase {
	p := doc.start(phrase.PropertyInitialiser, false)
	doc.next(false) //equals
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) memberModifierList() *phrase.Phrase {
	doc.start(phrase.MemberModifierList, false)

	for isMemberModifier(doc.peek(0)) {
		doc.next(false)
	}

	return doc.end()
}

func isMemberModifier(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Public,
		lexer.Protected,
		lexer.Private,
		lexer.Static,
		lexer.Abstract,
		lexer.Final:
		return true
	}

	return false
}

func (doc *Parser) qualifiedNameList(breakOn []lexer.TokenType) *phrase.Phrase {

	return doc.delimitedList(
		phrase.QualifiedNameList,
		doc.qualifiedName,
		isQualifiedNameStart,
		lexer.Comma,
		breakOn,
		false)
}

func (doc *Parser) objectCreationExpression() *phrase.Phrase {
	p := doc.start(phrase.ObjectCreationExpression, false)
	doc.next(false) //new
	if doc.peek(0).Type == lexer.Class {
		p.Children = append(p.Children, doc.anonymousClassDeclaration())

		return doc.end()
	}

	p.Children = append(p.Children, doc.typeDesignator(phrase.ClassTypeDesignator))

	if doc.optional(lexer.OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}

		doc.expect(lexer.CloseParenthesis)
	}

	return doc.end()

}

func (doc *Parser) typeDesignator(phraseType phrase.PhraseType) *phrase.Phrase {
	p := doc.start(phraseType, false)
	part := doc.classTypeDesignatorAtom()

	for {
		switch doc.peek(0).Type {
		case lexer.OpenBracket:
			part = doc.subscriptExpression(part, lexer.CloseBracket)
			continue
		case lexer.OpenBrace:
			part = doc.subscriptExpression(part, lexer.CloseBrace)
			continue
		case lexer.Arrow:
			part = doc.propertyAccessExpression(part)
			continue
		case lexer.ColonColon:
			staticPropNode := doc.start(phrase.ScopedPropertyAccessExpression, false)
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

func (doc *Parser) restrictedScopedMemberName() *phrase.Phrase {
	p := doc.start(phrase.ScopedMemberName, false)
	t := doc.peek(0)

	switch t.Type {
	case lexer.VariableName:
		//Spec says this should be SimpleVariable
		//leaving as a token as this avoids confusion between
		//static property names and simple variables
		doc.next(false)
	case lexer.Dollar:
		p.Children = append(p.Children, doc.simpleVariable())
	default:
		doc.error(lexer.Undefined)
	}

	return doc.end()
}

func (doc *Parser) classTypeDesignatorAtom() phrase.AstNode {
	t := doc.peek(0)

	switch t.Type {
	case lexer.Static:
		return doc.relativeScope()
	case lexer.VariableName, lexer.Dollar:
		return doc.simpleVariable()
	case lexer.Name, lexer.Namespace, lexer.Backslash:
		return doc.qualifiedName()
	default:
		doc.start(phrase.ErrorClassTypeDesignatorAtom, false)
		doc.error(lexer.Undefined)

		return doc.end()
	}
}

func (doc *Parser) cloneExpression() *phrase.Phrase {
	p := doc.start(phrase.CloneExpression, false)
	doc.next(false) //clone
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) listIntrinsic() *phrase.Phrase {
	p := doc.start(phrase.ListIntrinsic, false)
	doc.next(false) //list
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.arrayInitialiserList(lexer.CloseParenthesis))
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) unaryExpression(phraseType phrase.PhraseType) *phrase.Phrase {
	p := doc.start(phraseType, false)
	op := doc.next(false) //op

	switch phraseType {
	case phrase.PrefixDecrementExpression, phrase.PrefixIncrementExpression:
		p.Children = append(p.Children, doc.variable(doc.variableAtom(0)))
	default:
		precendence, _ := precedenceAssociativityTuple(op)
		p.Children = append(p.Children, doc.expression(precendence))
	}

	return doc.end()
}

func (doc *Parser) anonymousFunctionHeader() *phrase.Phrase {
	p := doc.start(phrase.AnonymousFunctionHeader, false)
	doc.optional(lexer.Static)
	doc.next(false) //function
	doc.optional(lexer.Ampersand)
	doc.expect(lexer.OpenParenthesis)

	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			phrase.ParameterDeclarationList,
			doc.parameterDeclaration,
			isParameterStart,
			lexer.Comma,
			[]lexer.TokenType{lexer.CloseParenthesis},
			false))
	}

	doc.expect(lexer.CloseParenthesis)

	if doc.peek(0).Type == lexer.Use {
		p.Children = append(p.Children, doc.anonymousFunctionUseClause())
	}

	if doc.peek(0).Type == lexer.Colon {
		p.Children = append(p.Children, doc.returnType())
	}

	return doc.end()
}

func (doc *Parser) anonymousFunctionCreationExpression() *phrase.Phrase {
	p := doc.start(phrase.AnonymousFunctionCreationExpression, false)

	p.Children = append(p.Children, doc.anonymousFunctionHeader())
	p.Children = append(p.Children, doc.functionDeclarationBody())

	return doc.end()
}

func isAnonymousFunctionUseVariableStart(t *lexer.Token) bool {
	return t.Type == lexer.VariableName || t.Type == lexer.Ampersand
}

func (doc *Parser) anonymousFunctionUseClause() *phrase.Phrase {
	p := doc.start(phrase.AnonymousFunctionUseClause, false)
	doc.next(false) //use
	doc.expect(lexer.OpenParenthesis)
	p.Children = append(p.Children, doc.delimitedList(
		phrase.ClosureUseList,
		doc.anonymousFunctionUseVariable,
		isAnonymousFunctionUseVariableStart,
		lexer.Comma,
		[]lexer.TokenType{lexer.CloseParenthesis},
		false))
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) anonymousFunctionUseVariable() phrase.AstNode {
	doc.start(phrase.AnonymousFunctionUseVariable, false)
	doc.optional(lexer.Ampersand)
	doc.expect(lexer.VariableName)

	return doc.end()
}

func isTypeDeclarationStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Backslash,
		lexer.Name,
		lexer.Namespace,
		lexer.Question,
		lexer.Array,
		lexer.Callable:
		return true
	}

	return false
}

func (doc *Parser) parameterDeclaration() phrase.AstNode {
	p := doc.start(phrase.ParameterDeclaration, false)

	if isTypeDeclarationStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.typeDeclaration())
	}

	doc.optional(lexer.Ampersand)
	doc.optional(lexer.Ellipsis)
	doc.expect(lexer.VariableName)

	if doc.peek(0).Type == lexer.Equals {
		doc.next(false)
		p.Children = append(p.Children, doc.expression(0))
	}

	return doc.end()
}

func (doc *Parser) variable(variableAtomNode phrase.AstNode) phrase.AstNode {
	count := 0

	for {
		count++
		switch doc.peek(0).Type {
		case lexer.ColonColon:
			variableAtomNode = doc.scopedAccessExpression(variableAtomNode)
			continue
		case lexer.Arrow:
			variableAtomNode = doc.propertyOrMethodAccessExpression(variableAtomNode)
			continue
		case lexer.OpenBracket:
			variableAtomNode = doc.subscriptExpression(variableAtomNode, lexer.CloseBracket)
			continue
		case lexer.OpenBrace:
			variableAtomNode = doc.subscriptExpression(variableAtomNode, lexer.CloseBrace)
			continue
		case lexer.OpenParenthesis:
			variableAtomNode = doc.functionCallExpression(variableAtomNode)
			continue
		default:
			//only simple variable atoms qualify as variables
			p := variableAtomNode.(*phrase.Phrase)

			if count == 1 && p.Type != phrase.SimpleVariable {
				errNode := doc.start(phrase.ErrorVariable, true)
				errNode.Children = append(errNode.Children, variableAtomNode)
				doc.error(lexer.Undefined)

				return doc.end()
			}
		}

		break
	}

	return variableAtomNode
}

func (doc *Parser) functionCallExpression(lhs phrase.AstNode) *phrase.Phrase {
	p := doc.start(phrase.FunctionCallExpression, true)
	p.Children = append(p.Children, lhs)
	doc.expect(lexer.OpenParenthesis)
	if isArgumentStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.argumentList())
	}
	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func (doc *Parser) scopedAccessExpression(lhs phrase.AstNode) *phrase.Phrase {
	p := doc.start(phrase.ErrorScopedAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //::
	p.Children = append(p.Children, doc.scopedMemberName(p))

	if doc.optional(lexer.OpenParenthesis) != nil {
		p.Type = phrase.ScopedCallExpression
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}

		doc.expect(lexer.CloseParenthesis)

		return doc.end()
	} else if p.Type == phrase.ScopedCallExpression {
		//error
		doc.error(lexer.Undefined)
	}

	return doc.end()

}

func (doc *Parser) scopedMemberName(parent *phrase.Phrase) *phrase.Phrase {
	p := doc.start(phrase.ScopedMemberName, false)
	t := doc.peek(0)

	switch t.Type {
	case lexer.OpenBrace:
		parent.Type = phrase.ScopedCallExpression
		p.Children = append(p.Children, doc.encapsulatedExpression(
			lexer.OpenBrace, lexer.CloseBrace))
	case lexer.VariableName:
		//Spec says this should be SimpleVariable
		//leaving as a token as this avoids confusion between
		//static property names and simple variables
		parent.Type = phrase.ScopedPropertyAccessExpression
		doc.next(false)
	case lexer.Dollar:
		p.Children = append(p.Children, doc.simpleVariable())
		parent.Type = phrase.ScopedPropertyAccessExpression
	default:
		if t.Type == lexer.Name || isSemiReservedToken(t) {
			p.Children = append(p.Children, doc.identifier())
			parent.Type = phrase.ClassConstantAccessExpression
		} else {
			//error
			doc.error(lexer.Undefined)
		}
	}

	return doc.end()

}

func (doc *Parser) propertyAccessExpression(lhs phrase.AstNode) *phrase.Phrase {
	p := doc.start(phrase.PropertyAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //->
	p.Children = append(p.Children, doc.memberName())

	return doc.end()
}

func (doc *Parser) propertyOrMethodAccessExpression(lhs phrase.AstNode) *phrase.Phrase {

	p := doc.start(phrase.PropertyAccessExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) //->
	p.Children = append(p.Children, doc.memberName())

	if doc.optional(lexer.OpenParenthesis) != nil {
		if isArgumentStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.argumentList())
		}
		p.Type = phrase.MethodCallExpression
		doc.expect(lexer.CloseParenthesis)
	}

	return doc.end()

}

func (doc *Parser) memberName() *phrase.Phrase {
	p := doc.start(phrase.MemberName, false)

	switch doc.peek(0).Type {
	case lexer.Name:
		doc.next(false)
	case lexer.OpenBrace:
		p.Children = append(p.Children, doc.encapsulatedExpression(
			lexer.OpenBrace, lexer.CloseBrace))
	case lexer.Dollar, lexer.VariableName:
		p.Children = append(p.Children, doc.simpleVariable())
	default:
		doc.error(lexer.Undefined)
	}

	return doc.end()
}

func (doc *Parser) subscriptExpression(
	lhs phrase.AstNode, closeTokenType lexer.TokenType) *phrase.Phrase {

	p := doc.start(phrase.SubscriptExpression, true)
	p.Children = append(p.Children, lhs)
	doc.next(false) // [ or {
	if isExpressionStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.expression(0))
	}

	doc.expect(closeTokenType)

	return doc.end()
}

func (doc *Parser) argumentList() *phrase.Phrase {
	return doc.delimitedList(
		phrase.ArgumentExpressionList,
		doc.argumentExpression,
		isArgumentStart,
		lexer.Comma,
		[]lexer.TokenType{lexer.CloseParenthesis},
		false)
}

func isArgumentStart(t *lexer.Token) bool {
	return t.Type == lexer.Ellipsis || isExpressionStart(t)
}

func (doc *Parser) variadicUnpacking() *phrase.Phrase {
	p := doc.start(phrase.VariadicUnpacking, false)
	doc.next(false) //...
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) argumentExpression() phrase.AstNode {
	if doc.peek(0).Type == lexer.Ellipsis {
		return doc.variadicUnpacking()
	} else {
		return doc.expression(0)
	}
}

func (doc *Parser) qualifiedName() phrase.AstNode {
	p := doc.start(phrase.QualifiedName, false)
	t := doc.peek(0)

	if t.Type == lexer.Backslash {
		doc.next(false)
		p.Type = phrase.FullyQualifiedName
	} else if t.Type == lexer.Namespace {
		p.Type = phrase.RelativeQualifiedName
		doc.next(false)
		doc.expect(lexer.Backslash)
	}

	p.Children = append(p.Children, doc.namespaceName())

	return doc.end()
}

func isQualifiedNameStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Backslash, lexer.Name, lexer.Namespace:
		return true
	}

	return false
}

func (doc *Parser) shortArrayCreationExpression(precedence int) *phrase.Phrase {
	p := doc.start(phrase.ArrayCreationExpression, false)
	doc.next(false) //[
	if isArrayElementStart(doc.peek(0)) || (precedence == 0 && doc.peek(0).Type == lexer.Comma) {
		p.Children = append(p.Children, doc.arrayInitialiserList(lexer.CloseBracket))
	}
	doc.expect(lexer.CloseBracket)

	return doc.end()
}

func (doc *Parser) longArrayCreationExpression() *phrase.Phrase {
	p := doc.start(phrase.ArrayCreationExpression, false)
	doc.next(false) //array
	doc.expect(lexer.OpenParenthesis)

	if isArrayElementStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.arrayInitialiserList(lexer.CloseParenthesis))
	}

	doc.expect(lexer.CloseParenthesis)

	return doc.end()
}

func isArrayElementStart(t *lexer.Token) bool {
	return t.Type == lexer.Ampersand || isExpressionStart(t)
}

func (doc *Parser) arrayInitialiserList(breakOn lexer.TokenType) *phrase.Phrase {

	p := doc.start(phrase.ArrayInitialiserList, false)
	var t *lexer.Token

	arrayInitialiserListRecoverSet := []lexer.TokenType{breakOn, lexer.Comma}
	doc.recoverSetStack = append(doc.recoverSetStack, arrayInitialiserListRecoverSet)

	for {
		//an array can have empty elements
		if isArrayElementStart(doc.peek(0)) {
			p.Children = append(p.Children, doc.arrayElement())
		}

		t = doc.peek(0)

		if t.Type == lexer.Comma {
			doc.next(false)
		} else if t.Type == breakOn {
			break
		} else {
			doc.error(lexer.Undefined)
			//check for missing delimeter
			if isArrayElementStart(t) {
				continue
			} else {
				//skip until recover token
				doc.defaultSyncStrategy()
				t = doc.peek(0)
				if t.Type == lexer.Comma || t.Type == breakOn {
					continue
				}
			}

			break
		}

	}

	doc.recoverSetStack = doc.recoverSetStack[:len(doc.recoverSetStack)-1]

	return doc.end()
}

func (doc *Parser) arrayValue() *phrase.Phrase {
	p := doc.start(phrase.ArrayValue, false)
	doc.optional(lexer.Ampersand)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) arrayKey() *phrase.Phrase {
	p := doc.start(phrase.ArrayKey, false)
	p.Children = append(p.Children, doc.expression(0))

	return doc.end()
}

func (doc *Parser) arrayElement() *phrase.Phrase {
	p := doc.start(phrase.ArrayElement, false)

	if doc.peek(0).Type == lexer.Ampersand {
		p.Children = append(p.Children, doc.arrayValue())

		return doc.end()
	}

	keyOrValue := doc.arrayKey()
	p.Children = append(p.Children, keyOrValue)

	if doc.optional(lexer.FatArrow) == nil {
		keyOrValue.Type = phrase.ArrayValue

		return doc.end()
	}

	p.Children = append(p.Children, doc.arrayValue())

	return doc.end()
}

func (doc *Parser) encapsulatedExpression(
	openTokenType lexer.TokenType, closeTokenType lexer.TokenType) *phrase.Phrase {
	p := doc.start(phrase.EncapsulatedExpression, false)
	doc.expect(openTokenType)
	p.Children = append(p.Children, doc.expression(0))
	doc.expect(closeTokenType)

	return doc.end()
}

func (doc *Parser) relativeScope() *phrase.Phrase {
	doc.start(phrase.RelativeScope, false)
	doc.next(false)

	return doc.end()
}

func (doc *Parser) variableAtom(precedence int) phrase.AstNode {
	t := doc.peek(0)
	switch t.Type {
	case lexer.VariableName, lexer.Dollar:
		return doc.simpleVariable()
	case lexer.OpenParenthesis:
		return doc.encapsulatedExpression(lexer.OpenParenthesis, lexer.CloseParenthesis)
	case lexer.Array:
		return doc.longArrayCreationExpression()
	case lexer.OpenBracket:
		return doc.shortArrayCreationExpression(precedence)
	case lexer.StringLiteral:
		return doc.next(true)
	case lexer.Static:
		return doc.relativeScope()
	case lexer.Name, lexer.Namespace, lexer.Backslash:
		return doc.qualifiedName()
	default:
		//error
		doc.start(phrase.ErrorVariableAtom, false)
		doc.error(lexer.Undefined)

		return doc.end()
	}
}

func (doc *Parser) simpleVariable() phrase.AstNode {
	p := doc.start(phrase.SimpleVariable, false)
	t := doc.expectOneOf([]lexer.TokenType{lexer.VariableName, lexer.Dollar})

	if t != nil && t.Type == lexer.Dollar {
		t = doc.peek(0)
		if t.Type == lexer.OpenBrace {
			p.Children = append(p.Children, doc.encapsulatedExpression(
				lexer.OpenBrace, lexer.CloseBrace))
		} else if t.Type == lexer.Dollar || t.Type == lexer.VariableName {
			p.Children = append(p.Children, doc.simpleVariable())
		} else {
			doc.error(lexer.Undefined)
		}
	}

	return doc.end()
}

func (doc *Parser) haltCompilerStatement() *phrase.Phrase {
	doc.start(phrase.HaltCompilerStatement, false)
	doc.next(false) // __halt_compiler
	doc.expect(lexer.OpenParenthesis)
	doc.expect(lexer.CloseParenthesis)
	doc.expect(lexer.Semicolon)

	return doc.end()
}

func (doc *Parser) namespaceUseDeclaration() *phrase.Phrase {
	p := doc.start(phrase.NamespaceUseDeclaration, false)
	doc.next(false) //use
	doc.optionalOneOf([]lexer.TokenType{lexer.Function, lexer.Const})
	doc.optional(lexer.Backslash)
	nsNameNode := doc.namespaceName()
	t := doc.peek(0)

	if t.Type == lexer.Backslash || t.Type == lexer.OpenBrace {
		p.Children = append(p.Children, nsNameNode)
		doc.expect(lexer.Backslash)
		doc.expect(lexer.OpenBrace)
		p.Children = append(p.Children, doc.delimitedList(
			phrase.NamespaceUseGroupClauseList,
			doc.namespaceUseGroupClause,
			isNamespaceUseGroupClauseStartToken,
			lexer.Comma,
			[]lexer.TokenType{lexer.CloseBrace},
			false))
		doc.expect(lexer.CloseBrace)
		doc.expect(lexer.Semicolon)

		return doc.end()
	}

	p.Children = append(p.Children, doc.delimitedList(
		phrase.NamespaceUseClauseList,
		doc.namespaceUseClauseFunction(nsNameNode),
		isNamespaceUseClauseStartToken,
		lexer.Comma,
		[]lexer.TokenType{lexer.Semicolon},
		true))

	doc.expect(lexer.Semicolon)

	return doc.end()
}

func isNamespaceUseClauseStartToken(t *lexer.Token) bool {
	return t.Type == lexer.Name || t.Type == lexer.Backslash
}

func (doc *Parser) namespaceUseClauseFunction(nsName *phrase.Phrase) func() phrase.AstNode {

	return func() phrase.AstNode {
		p := doc.start(phrase.NamespaceUseClause, nsName != nil)

		if nsName != nil {
			p.Children = append(p.Children, nsName)
			nsName = nil
		} else {
			p.Children = append(p.Children, doc.namespaceName())
		}

		if doc.peek(0).Type == lexer.As {
			p.Children = append(p.Children, doc.namespaceAliasingClause())
		}

		return doc.end()
	}
}

func (doc *Parser) delimitedList(
	phraseType phrase.PhraseType,
	elementFunction func() phrase.AstNode,
	elementStartPredicate func(*lexer.Token) bool,
	delimiter lexer.TokenType,
	breakOn []lexer.TokenType,
	doNotPushHiddenToParent bool) *phrase.Phrase {

	p := doc.start(phraseType, doNotPushHiddenToParent)
	var t *lexer.Token
	var delimitedListRecoverSet []lexer.TokenType
	if breakOn != nil {
		delimitedListRecoverSet = make([]lexer.TokenType, len(breakOn))
		copy(delimitedListRecoverSet, breakOn)
	} else {
		delimitedListRecoverSet = make([]lexer.TokenType, 0)
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
			doc.error(lexer.Undefined)
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

func isNamespaceUseGroupClauseStartToken(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Const, lexer.Function, lexer.Name:
		return true
	}

	return false
}

func (doc *Parser) namespaceUseGroupClause() phrase.AstNode {
	p := doc.start(phrase.NamespaceUseGroupClause, false)
	doc.optionalOneOf([]lexer.TokenType{lexer.Function, lexer.Const})
	p.Children = append(p.Children, doc.namespaceName())

	if doc.peek(0).Type == lexer.As {
		p.Children = append(p.Children, doc.namespaceAliasingClause())
	}

	return doc.end()
}

func (doc *Parser) namespaceAliasingClause() *phrase.Phrase {
	doc.start(phrase.NamespaceAliasingClause, false)
	doc.next(false) //as
	doc.expect(lexer.Name)

	return doc.end()
}

func (doc *Parser) namespaceDefinition() *phrase.Phrase {
	p := doc.start(phrase.NamespaceDefinition, false)
	doc.next(false) //namespace
	if doc.peek(0).Type == lexer.Name {

		p.Children = append(p.Children, doc.namespaceName())
		t := doc.expectOneOf([]lexer.TokenType{lexer.Semicolon, lexer.OpenBrace})
		if t == nil || t.Type != lexer.OpenBrace {
			return doc.end()
		}

	} else {
		doc.expect(lexer.OpenBrace)
	}

	p.Children = append(p.Children, doc.statementList([]lexer.TokenType{lexer.CloseBrace}))
	doc.expect(lexer.CloseBrace)

	return doc.end()
}

func (doc *Parser) namespaceName() *phrase.Phrase {
	doc.start(phrase.NamespaceName, false)
	doc.expect(lexer.Name)

	for {
		if doc.peek(0).Type == lexer.Backslash &&
			doc.peek(1).Type == lexer.Name {
			doc.next(false)
			doc.next(false)
		} else {
			break
		}
	}

	return doc.end()
}

func isReservedToken(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Include,
		lexer.IncludeOnce,
		lexer.Eval,
		lexer.Require,
		lexer.RequireOnce,
		lexer.Or,
		lexer.Xor,
		lexer.And,
		lexer.InstanceOf,
		lexer.New,
		lexer.Clone,
		lexer.Exit,
		lexer.If,
		lexer.ElseIf,
		lexer.Else,
		lexer.EndIf,
		lexer.Echo,
		lexer.Do,
		lexer.While,
		lexer.EndWhile,
		lexer.For,
		lexer.EndFor,
		lexer.ForEach,
		lexer.EndForeach,
		lexer.Declare,
		lexer.EndDeclare,
		lexer.As,
		lexer.Try,
		lexer.Catch,
		lexer.Finally,
		lexer.Throw,
		lexer.Use,
		lexer.InsteadOf,
		lexer.Global,
		lexer.Var,
		lexer.Unset,
		lexer.Isset,
		lexer.Empty,
		lexer.Continue,
		lexer.Goto,
		lexer.Function,
		lexer.Const,
		lexer.Return,
		lexer.Print,
		lexer.Yield,
		lexer.List,
		lexer.Switch,
		lexer.EndSwitch,
		lexer.Case,
		lexer.Default,
		lexer.Break,
		lexer.Array,
		lexer.Callable,
		lexer.Extends,
		lexer.Implements,
		lexer.Namespace,
		lexer.Trait,
		lexer.Interface,
		lexer.Class,
		lexer.ClassConstant,
		lexer.TraitConstant,
		lexer.FunctionConstant,
		lexer.MethodConstant,
		lexer.LineConstant,
		lexer.FileConstant,
		lexer.DirectoryConstant,
		lexer.NamespaceConstant:
		return true
	}

	return false
}

func isSemiReservedToken(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Static,
		lexer.Abstract,
		lexer.Final,
		lexer.Private,
		lexer.Protected,
		lexer.Public:
		return true
	}

	return isReservedToken(t)
}

func isStatementStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.Namespace,
		lexer.Use,
		lexer.HaltCompiler,
		lexer.Const,
		lexer.Function,
		lexer.Class,
		lexer.Abstract,
		lexer.Final,
		lexer.Trait,
		lexer.Interface,
		lexer.OpenBrace,
		lexer.If,
		lexer.While,
		lexer.Do,
		lexer.For,
		lexer.Switch,
		lexer.Break,
		lexer.Continue,
		lexer.Return,
		lexer.Global,
		lexer.Static,
		lexer.Echo,
		lexer.Unset,
		lexer.ForEach,
		lexer.Declare,
		lexer.Try,
		lexer.Throw,
		lexer.Goto,
		lexer.Name,
		lexer.Semicolon,
		lexer.CloseTag,
		lexer.Text,
		lexer.OpenTag,
		lexer.OpenTagEcho,
		lexer.DocumentCommentStart:
		return true
	}

	return isExpressionStart(t)
}
