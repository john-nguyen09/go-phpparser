package parser_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/john-nguyen09/go-phpparser/lexer"
	"github.com/john-nguyen09/go-phpparser/parser"
	"github.com/john-nguyen09/go-phpparser/phrase"
)

func TestParser(t *testing.T) {
	dir := "../cases"
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

}

func TestPerformance(t *testing.T) {
	dir := "../cases/moodle"

	start := time.Now()
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(path, ".php") {
			if path != "..\\cases\\moodle\\admin\\tests\\behat\\behat_admin.php" {
				return nil
			}

			data, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			parser.Parse(string(data))
		}

		return nil
	})
	elapsed := time.Since(start)

	fmt.Printf("Parser took %s to finish\n", elapsed)
}

func traverse(writer *bufio.Writer, node phrase.AstNode, depth int) {
	var p *phrase.Phrase
	var err *phrase.ParseError
	var isPhrase, isParseError bool
	indent := ""

	for i := 0; i < depth; i++ {
		indent += "-"
	}

	if p, isPhrase = node.(*phrase.Phrase); isPhrase {
		fmt.Fprintln(writer, indent+p.Type.String())
	} else if _, isToken := node.(*lexer.Token); isToken {
		fmt.Fprintln(writer, indent+node.(*lexer.Token).String())
	} else if err, isParseError = node.(*phrase.ParseError); isParseError {
		fmt.Fprintln(writer, indent+err.Type.String())
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
