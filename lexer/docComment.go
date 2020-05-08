package lexer

import (
	"strings"
	"unicode/utf8"
)

func (s *Lexer) scriptingDocBlock() *Token {
	start := s.offset
	switch s.r {
	case ' ', '\t':
		for s.step(); s.r == ' ' || s.r == '\t'; s.step() {
		}
		return NewToken(Whitespace, start, s.offset-start)
	case '\r', '\n':
		s.step()
		return NewToken(DocumentCommentEndline, start, s.offset-start)
	case '@':
		return s.docBlockTagName()
	case '$':
		return s.scriptingDollar()
	case '*':
		if s.peek(1) == '/' {
			s.stepLoop(2)
			s.modeStack = s.modeStack[:len(s.modeStack)-1]
			return NewToken(DocumentCommentEnd, start, s.offset-start)
		}
		for s.step(); isWhitespace(s.r) || s.r == '*'; s.step() {
			if s.r == '*' && s.peek(1) == '/' {
				break
			}
		}
		return NewToken(DocumentCommentStartline, start, s.offset-start)
	default:
		return s.scriptingDocumentBlockLabel()
	}
}

func (s *Lexer) scriptingDocumentBlockLabel() *Token {
	start := s.offset
	c := s.r
	s.step()
	if c == '[' && s.r == ']' {
		s.step()
		return NewToken(Array, start, s.offset-start)
	}
	if c == '|' {
		return NewToken(Bar, start, s.offset-start)
	}
	if c == '/' {
		return NewToken(ForwardSlash, start, s.offset-start)
	}
	if c == '\\' {
		return NewToken(Backslash, start, s.offset-start)
	}
	if c == '<' {
		return NewToken(LessThan, start, s.offset-start)
	}
	if c == '>' {
		return NewToken(GreaterThan, start, s.offset-start)
	}
	if c == '(' {
		return NewToken(OpenParenthesis, start, s.offset-start)
	}
	if c == ')' {
		return NewToken(CloseParenthesis, start, s.offset-start)
	}
	if c == '=' {
		return NewToken(Equals, start, s.offset-start)
	}
	if c == ',' {
		return NewToken(Comma, start, s.offset-start)
	}
	if (c == 's' || c == 'S') &&
		(s.r == 't' || s.r == 'T') &&
		strings.ToLower(s.peekSpanString(0, 4)) == "atic" &&
		s.peek(5) == ' ' {
		s.stepLoop(5)
		return NewToken(Static, start, s.offset-start)
	}
	if isDigit(c) {
		tokenType := IntegerLiteral
		for isDigit(s.r) || s.r == '.' {
			if s.r == '.' {
				tokenType = DocumentCommentVersion
			}
			s.step()
		}
		return NewToken(tokenType, start, s.offset-start)
	}
	if c == '$' && isLabelStart(s.r) {
		return NewToken(Dollar, start, s.offset-start)
	}
	if isLabelStart(c) {
		for ; isLabelChar(s.r); s.step() {
		}
		return NewToken(Name, start, s.offset-start)
	}
	if isDocCommentText(c) {
		for ; isDocCommentText(s.r) && s.r != '[' && s.r != '|' &&
			s.r != '/' && s.r != '\\' &&
			s.r != '<' && s.r != '>' && s.r != '(' && s.r != ')'; s.step() {
		}
		return NewToken(DocumentCommentText, start, s.offset-start)
	}
	for ; !isDocCommentText(s.r); s.step() {
	}
	return NewToken(DocumentCommentUnknown, start, s.offset-start)
}

func (s *Lexer) docBlockTagName() *Token {
	start := s.offset
	startLabel := 1
	endLabel := startLabel
	for ; !isWhitespace(s.peek(endLabel)); endLabel++ {
	}
	tagName := s.peekSpanString(startLabel-1, endLabel-1)
	tokenType := DocumentCommentTagName
	switch tagName {
	case "author":
		tokenType = AtAuthor
	case "deprecated":
		tokenType = AtDeprecated
	case "global":
		tokenType = AtGlobal
	case "link":
		tokenType = AtLink
	case "method":
		tokenType = AtMethod
	case "param":
		tokenType = AtParam
	case "property":
		tokenType = AtProperty
	case "property-read":
		tokenType = AtPropertyRead
	case "property-write":
		tokenType = AtPropertyWrite
	case "return":
		tokenType = AtReturn
	case "since":
		tokenType = AtSince
	case "throws":
		tokenType = AtThrows
	case "var":
		tokenType = AtVar
	}
	s.stepLoop(endLabel)
	return NewToken(tokenType, start, s.offset-start)
}

func isDocCommentText(cp rune) bool {
	return (cp >= '!' && cp <= '~') || cp >= utf8.RuneSelf
}
