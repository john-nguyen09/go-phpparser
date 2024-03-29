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
		return NewToken(s.pool, Whitespace, start, s.offset-start)
	case '\r', '\n':
		s.step()
		return NewToken(s.pool, DocumentCommentEndline, start, s.offset-start)
	case '@':
		return s.docBlockTagName()
	case '$':
		return s.scriptingDollar()
	case '*':
		if s.peek(1) == '/' {
			s.stepLoop(2)
			s.modeStack = s.modeStack[:len(s.modeStack)-1]
			return NewToken(s.pool, DocumentCommentEnd, start, s.offset-start)
		}
		for s.step(); isWhitespace(s.r) || s.r == '*'; s.step() {
			if s.r == '*' && s.peek(1) == '/' {
				break
			}
		}
		return NewToken(s.pool, DocumentCommentStartline, start, s.offset-start)
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
		return NewToken(s.pool, Array, start, s.offset-start)
	}
	if c == '|' {
		return NewToken(s.pool, Bar, start, s.offset-start)
	}
	if c == '/' {
		return NewToken(s.pool, ForwardSlash, start, s.offset-start)
	}
	if c == '\\' {
		return NewToken(s.pool, Backslash, start, s.offset-start)
	}
	if c == '<' {
		return NewToken(s.pool, LessThan, start, s.offset-start)
	}
	if c == '>' {
		return NewToken(s.pool, GreaterThan, start, s.offset-start)
	}
	if c == '(' {
		return NewToken(s.pool, OpenParenthesis, start, s.offset-start)
	}
	if c == ')' {
		return NewToken(s.pool, CloseParenthesis, start, s.offset-start)
	}
	if c == '=' {
		return NewToken(s.pool, Equals, start, s.offset-start)
	}
	if c == ',' {
		return NewToken(s.pool, Comma, start, s.offset-start)
	}
	if (c == 's' || c == 'S') &&
		(s.r == 't' || s.r == 'T') &&
		strings.ToLower(s.peekSpanString(0, 4)) == "atic" &&
		s.peek(5) == ' ' {
		s.stepLoop(5)
		return NewToken(s.pool, Static, start, s.offset-start)
	}
	if isDigit(c) {
		tokenType := IntegerLiteral
		for isDigit(s.r) || s.r == '.' {
			if s.r == '.' {
				tokenType = DocumentCommentVersion
			}
			s.step()
		}
		return NewToken(s.pool, tokenType, start, s.offset-start)
	}
	if c == '$' && isLabelStart(s.r) {
		return NewToken(s.pool, Dollar, start, s.offset-start)
	}
	if isLabelStart(c) {
		for ; isLabelChar(s.r); s.step() {
		}
		return NewToken(s.pool, Name, start, s.offset-start)
	}
	if isDocCommentText(c, s.r) {
		for ; isDocCommentText(s.r, s.peek(1)) && s.r != '[' && s.r != '|' &&
			s.r != '/' && s.r != '\\' &&
			s.r != '<' && s.r != '>' && s.r != '(' && s.r != ')'; s.step() {
		}
		return NewToken(s.pool, DocumentCommentText, start, s.offset-start)
	}
	for ; !isDocCommentText(s.r, s.peek(1)); s.step() {
	}
	return NewToken(s.pool, DocumentCommentUnknown, start, s.offset-start)
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
	return NewToken(s.pool, tokenType, start, s.offset-start)
}

func isDocCommentText(cp, next rune) bool {
	if cp == '*' && next == '/' {
		return false
	}
	return (cp >= '!' && cp <= '~') || cp >= utf8.RuneSelf
}
