package main

import (
	"fmt"
	"strconv"

	"github.com/john-nguyen09/go-phpparser/lexer"
	"github.com/john-nguyen09/go-phpparser/phrase"
)

func main() {
	astNodes := make([]phrase.AstNode, 0)

	astNodes = append(astNodes, lexer.NewToken(lexer.EndOfFile, 0, 1, nil))
	astNodes = append(astNodes, phrase.NewPhrase(phrase.ElseClause, nil))

	fmt.Println(strconv.Itoa(int(lexer.AsteriskAsterisk)))

	for _, astNode := range astNodes {
		if phrase, ok := astNode.(*phrase.Phrase); ok {
			fmt.Printf("It is a phrase: %v\n", phrase)
		}
	}
}
