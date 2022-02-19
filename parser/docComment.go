package parser

import (
	"github.com/john-nguyen09/go-phpparser/lexer"
	"github.com/john-nguyen09/go-phpparser/phrase"
)

func (doc *Parser) docComment() *phrase.Phrase {
	p := doc.list(phrase.DocumentComment,
		doc.docCommentStatement,
		isDocumentCommentStatementStart,
		[]lexer.TokenType{
			lexer.DocumentCommentEnd,
		},
		[]lexer.TokenType{
			lexer.DocumentCommentStartline,
			lexer.DocumentCommentEndline,
		})
	doc.hidden(p)
	p.Children = append(p.Children, doc.next(true)) // DocumentCommentEnd
	return p
}

func (doc *Parser) docCommentStatement() phrase.AstNode {
	t := doc.peek(0)

	switch t.Type {
	case lexer.DocumentCommentStart:
		return doc.next(true)
	case lexer.DocumentCommentStartline, lexer.Name, lexer.DocumentCommentText:
		return doc.documentCommentStatementStart()
	case lexer.DocumentCommentEndline, lexer.Whitespace:
		return doc.next(true)
	}
	if isTagName(t) {
		return doc.docCommentTag()
	}
	return nil
}

func (doc *Parser) docCommentDescription() *phrase.Phrase {
	doc.start(phrase.DocumentCommentDescription, false)
	if !isDescriptionStart(doc.peek(0)) {
		return doc.end()
	}
	doc.next(false)
	t := doc.peek(0)
	for {
		if t.Type == lexer.DocumentCommentEndline &&
			(doc.peek(1).Type != lexer.DocumentCommentStartline || !isDescriptionStart(doc.peek(2))) {
			doc.next(false)
			break
		}
		if t.Type == lexer.DocumentCommentEnd {
			break
		}
		doc.next(false)
		t = doc.peek(0)
	}
	return doc.end()
}

func (doc *Parser) docCommentTypeName() phrase.AstNode {
	p := doc.start(phrase.TypeDeclaration, false)
	switch doc.peek(0).Type {
	case lexer.Callable, lexer.Array:
		doc.next(false)
	case lexer.Name, lexer.Backslash:
		p.Children = append(p.Children, doc.qualifiedName())
		for doc.peek(0).Type == lexer.Array {
			doc.next(false)
		}
	case lexer.VariableName:
		doc.next(false)
	default:
		doc.end()
		return nil
	}
	return doc.end()
}

func (doc *Parser) docCommentTypeUnionOrTypeDeclaration() phrase.AstNode {
	p := doc.start(phrase.TypeUnion, false)
	typeName := doc.docCommentTypeName()
	if typeName == nil {
		doc.end()
		return nil
	}
	p.Children = append(p.Children, typeName)
	if doc.peek(0).Type != lexer.Bar {
		doc.end()
		return typeName
	}
	for doc.peek(0).Type == lexer.Bar {
		doc.next(false)
		typeName = doc.docCommentTypeName()
		if typeName != nil {
			p.Children = append(p.Children, typeName)
		} else {
			break
		}
	}
	return doc.end()
}

func (doc *Parser) documentCommentStatementStart() phrase.AstNode {
	if isTagName(doc.peek(1)) {
		return doc.docCommentTag()
	}
	return doc.docCommentDescription()
}

func (doc *Parser) docCommentTag() phrase.AstNode {
	p := doc.start(phrase.DocumentCommentTag, false)
	if doc.peek(0).Type == lexer.DocumentCommentStartline {
		doc.next(false)
	}
	t := doc.next(false)
	switch t.Type {
	case lexer.AtAuthor:
		doc.authorTag(p)
	case lexer.AtDeprecated:
		doc.deprecatedTag(p)
	case lexer.AtGlobal:
		doc.globalTag(p)
	case lexer.AtMethod:
		doc.methodTag(p)
	case lexer.AtParam:
		doc.paramTag(p)
	case lexer.AtProperty, lexer.AtPropertyRead, lexer.AtPropertyWrite:
		doc.propertyTag(p)
	case lexer.AtReturn:
		doc.returnTag(p)
	case lexer.AtThrows:
		doc.throwsTag(p)
	case lexer.AtVar:
		doc.varTag(p)
	default:
		p.Children = append(p.Children, doc.docCommentDescription())
	}
	return doc.end()
}

func (doc *Parser) authorTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentAuthorTag

	doc.start(phrase.DocumentCommentAuthor, false)
	doc.next(false) // Tag name
	t := doc.peek(0)
	for t.Type != lexer.DocumentCommentEndline && t.Type != lexer.LessThan {
		doc.next(false)
		t = doc.peek(0)
	}
	p.Children = append(p.Children, doc.end())

	if t.Type == lexer.LessThan {
		doc.start(phrase.DocumentCommentEmail, false)
		doc.next(false)
		t := doc.peek(0)
		for t.Type != lexer.DocumentCommentEndline && t.Type != lexer.GreaterThan {
			doc.next(false)
			t = doc.peek(0)
		}
		doc.next(false)
		p.Children = append(p.Children, doc.end())
	}
}

func (doc *Parser) deprecatedTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentDeprecatedTag

	doc.optional(lexer.DocumentCommentVersion)
	if doc.peek(0).Type != lexer.DocumentCommentEndline {
		desc := doc.docCommentDescription()
		if len(desc.Children) > 0 {
			p.Children = append(p.Children, desc)
		}
	}
}

func (doc *Parser) globalTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentGlobalTag

	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	} else {
		doc.error(lexer.Name)
	}

	doc.optional(lexer.VariableName)
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) docCommentParameterDeclaration() phrase.AstNode {
	p := doc.start(phrase.ParameterDeclaration, false)
	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	}
	doc.optional(lexer.Ampersand)
	doc.optional(lexer.Ellipsis)
	doc.expect(lexer.VariableName)
	if doc.peek(0).Type == lexer.Equals {
		doc.next(false)
		p.Children = append(p.Children, doc.docCommentParameterValue())
	}
	return doc.end()
}

func (doc *Parser) docCommentParameterValue() *phrase.Phrase {
	doc.start(phrase.ParameterValue, false)
	t := doc.peek(0)
	for t.Type != lexer.CloseParenthesis && t.Type != lexer.Comma {
		doc.next(false)
		t = doc.peek(0)
	}
	return doc.end()
}

func (doc *Parser) methodTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentMethodTag

	doc.optional(lexer.Static)
	if doc.peek(1).Type != lexer.OpenParenthesis {
		typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
		if typeListOrName != nil {
			p.Children = append(p.Children, typeListOrName)
		}
	}
	p.Children = append(p.Children, doc.identifier())
	doc.expect(lexer.OpenParenthesis)
	if isParameterStart(doc.peek(0)) {
		p.Children = append(p.Children, doc.delimitedList(
			phrase.ParameterDeclarationList,
			doc.docCommentParameterDeclaration,
			isParameterStart,
			lexer.Comma,
			[]lexer.TokenType{lexer.CloseParenthesis}, false, true))
	}
	doc.expect(lexer.CloseParenthesis)
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) paramTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentParamTag

	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	} else {
		doc.error(lexer.Name)
	}
	doc.expect(lexer.VariableName)
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) propertyTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentPropertyTag

	if doc.peek(0).Type != lexer.VariableName {
		typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
		if typeListOrName != nil {
			p.Children = append(p.Children, typeListOrName)
		} else {
			doc.error(lexer.Name)
		}
	}
	doc.expect(lexer.VariableName)
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) returnTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentReturnTag

	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	} else {
		doc.error(lexer.Name)
	}
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) throwsTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentThrowsTag

	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	} else {
		doc.error(lexer.Name)
	}
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func (doc *Parser) varTag(p *phrase.Phrase) {
	p.Type = phrase.DocumentCommentVarTag

	typeListOrName := doc.docCommentTypeUnionOrTypeDeclaration()
	if typeListOrName != nil {
		p.Children = append(p.Children, typeListOrName)
	} else {
		doc.error(lexer.Name)
	}
	doc.optional(lexer.VariableName)
	desc := doc.docCommentDescription()
	if len(desc.Children) > 0 {
		p.Children = append(p.Children, desc)
	}
}

func isDocumentCommentStatementStart(t *lexer.Token) bool {
	switch t.Type {
	case lexer.DocumentCommentStart,
		lexer.DocumentCommentStartline,
		lexer.DocumentCommentEndline,
		lexer.Name,
		lexer.DocumentCommentText:
		return true
	}
	return isTagName(t)
}

func isTagName(t *lexer.Token) bool {
	return (t.Type > lexer.DocumentCommentTagNameAnchorStart && t.Type < lexer.DocumentCommentTagNameAnchorEnd) ||
		t.Type == lexer.DocumentCommentTagName
}

func isDescriptionStart(t *lexer.Token) bool {
	return !isTagName(t) && t.Type != lexer.DocumentCommentEndline && t.Type != lexer.DocumentCommentEnd
}
