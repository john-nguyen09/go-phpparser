package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/john-nguyen09/go-phpparser/lexer"
	"github.com/john-nguyen09/go-phpparser/parser"
)

func TestParserAndLexer(t *testing.T) {
	dir := "cases"
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Println("Folder not found: " + dir)
		t.FailNow()
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".php") {
			continue
		}

		filePath := dir + "/" + file.Name()
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		t.Run(strings.TrimSuffix(file.Name(), path.Ext(file.Name())), func(t *testing.T) {
			lexerState := lexer.NewLexerState(string(data), nil, 0)
			tokens := []*lexer.Token{}
			token := lexerState.Lex()
			for {
				tokens = append(tokens, token)
				if token.Type == lexer.EndOfFile {
					break
				}
				token = lexerState.Lex()
			}
			rootNode := parser.Parse(string(data))
			cupaloy.SnapshotT(t, tokens, rootNode)
		})
	}
}
