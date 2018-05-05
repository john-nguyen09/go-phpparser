package lexer_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/john-nguyen09/phpparser/lexer"
)

func TestLexer(t *testing.T) {
	dir := "./cases"
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
			fmt.Printf("File not found: %s\n", filePath)
			continue
		}

		lexerState := lexer.NewLexerState(string(data), nil, 0)
		outFile, err := os.Create(filePath + ".lexed")

		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(outFile)

		for token := lexerState.Lex(); token.Type != lexer.EndOfFile; token = lexerState.Lex() {
			fmt.Fprintln(writer, token.String())
		}

		writer.Flush()
		outFile.Close()
	}

}
