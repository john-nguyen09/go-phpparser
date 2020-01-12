package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/john-nguyen09/go-phpparser/parser"
)

func main() {
	dir := "cases"
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Println("Folder not found: " + dir)
		return
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

		writeLex(filePath, data)
		writeParseTree(filePath, data)
	}
}

func writeLex(filePath string, data []byte) {
	lexerState := parser.NewLexerState(string(data), nil, 0)
	outFile, err := os.Create(filePath + ".lexed")

	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(outFile)

	for token := lexerState.Lex(); token.Type != parser.EndOfFile; token = lexerState.Lex() {
		fmt.Fprintln(writer, token.String())
	}

	writer.Flush()
	outFile.Close()
}

func writeParseTree(filePath string, data []byte) {
	outFile, err := os.Create(filePath + ".parsed")

	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(outFile)
	rootNode := parser.Parse(string(data))

	traverse(writer, rootNode, 0)

	writer.Flush()
	outFile.Close()
}

func traverse(writer *bufio.Writer, node parser.AstNode, depth int) {
	var p *parser.Phrase
	var err *parser.ParseError
	var isPhrase, isParseError bool
	indent := ""

	for i := 0; i < depth; i++ {
		indent += "-"
	}

	if p, isPhrase = node.(*parser.Phrase); isPhrase {
		fmt.Fprintln(writer, indent+p.Type.String()+"[Phrase]")
	} else if _, isToken := node.(*parser.Token); isToken {
		fmt.Fprintln(writer, indent+node.(*parser.Token).String()+"[Token]")
	} else if err, isParseError = node.(*parser.ParseError); isParseError {
		fmt.Fprintln(writer, indent+err.Type.String()+"[ParseError]")
		thisIndent := indent + "-"
		if len(err.Children) == 0 {
			fmt.Fprintln(writer, thisIndent+"Unexpected: "+err.Unexpected.String())
		} else {
			for _, child := range err.Children {
				if t, ok := child.(*parser.Token); ok {
					fmt.Fprintln(writer, thisIndent+t.Type.String())
				}
			}
		}
	}

	if p != nil && p.Children != nil {
		for _, child := range p.Children {
			traverse(writer, child, depth+1)
		}
	}

	if err != nil && err.Children != nil {
		for _, child := range err.Children {
			traverse(writer, child, depth+1)
		}
	}
}
