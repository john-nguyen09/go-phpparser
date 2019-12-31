package main

import (
	"log"
	"testing"

	"github.com/john-nguyen09/go-phpparser/lexer"
)

func TestLexerSync(t *testing.T) {
	data := "<?php echo 'Hello world';"
	text := []rune(data)
	tokens := lexer.Lex(data)
	max := len(data)
	for i := 0; i < max; i++ {
		hasErr := false
		change := lexer.ChangeEvent{
			Start: i,
			End:   i,
			Text:  []rune("r"),
		}
		newText := append(text[:0:0], text[0:change.Start]...)
		newText = append(newText, change.Text...)
		newText = append(newText, text[change.End:]...)
		log.Printf("%d %s", i, string(newText))
		changedTokens := lexer.Sync(string(newText), change, tokens)
		// for _, token := range changedTokens {
		// 	log.Println(token)
		// }
		newTokens := lexer.Lex(string(newText))
		if len(changedTokens) != len(newTokens) {
			t.Errorf("Length is not the same, %d != %d", len(changedTokens), len(newTokens))
			hasErr = true
		}
		for i := range changedTokens {
			if !tokensEqual(*changedTokens[i], *newTokens[i]) {
				t.Errorf("%v != %v", changedTokens[i], newTokens[i])
				hasErr = true
			}
		}
		if hasErr {
			t.FailNow()
		}
	}
}

func tokensEqual(t1 lexer.Token, t2 lexer.Token) bool {
	return t1.Type == t2.Type && t1.Offset == t2.Offset && t1.Length == t2.Length
}

func TestLexerSyncDelete(t *testing.T) {
	data := "<?php echo 'Hello world';"
	text := []rune(data)
	tokens := lexer.Lex(data)
	hasErr := false
	change := lexer.ChangeEvent{
		Start: 10,
		End:   11,
		Text:  []rune("print("),
	}
	newText := append(text[:0:0], text[0:change.Start]...)
	newText = append(newText, change.Text...)
	newText = append(newText, text[change.End:]...)
	log.Printf("%s", string(newText))
	changedTokens := lexer.Sync(string(newText), change, tokens)
	// for _, token := range changedTokens {
	// 	log.Println(token)
	// }
	newTokens := lexer.Lex(string(newText))
	if len(changedTokens) != len(newTokens) {
		t.Errorf("Length is not the same, %d != %d", len(changedTokens), len(newTokens))
		hasErr = true
	}
	for i := range changedTokens {
		if !tokensEqual(*changedTokens[i], *newTokens[i]) {
			t.Errorf("%v != %v", changedTokens[i], newTokens[i])
			hasErr = true
		}
	}
	if hasErr {
		t.FailNow()
	}
}

func TestLexerSyncNew(t *testing.T) {
	data := "<?php echo 'Hello world';"
	change := lexer.ChangeEvent{
		Start: 0,
		End:   0,
		Text:  []rune(data),
	}
	newTokens := lexer.Lex(data)
	changedTokens := lexer.Sync(data, change, nil)
	hasErr := false
	log.Printf("%s", data)
	// for _, token := range changedTokens {
	// 	log.Println(token)
	// }
	if len(changedTokens) != len(newTokens) {
		t.Errorf("Length is not the same, %d != %d", len(changedTokens), len(newTokens))
		hasErr = true
	}
	for i := range changedTokens {
		if !tokensEqual(*changedTokens[i], *newTokens[i]) {
			t.Errorf("%v != %v", changedTokens[i], newTokens[i])
			hasErr = true
		}
	}
	if hasErr {
		t.FailNow()
	}
}
