package parser

// In Go, an AST (Abstract Syntax Tree) is a hierarchical, tree-like data structure that represents the semantic and syntactic structure of source code.
// We are going to use a parser to get prepared AST structure.

import (
	"afv/internal/collector"
	"afv/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ParsedFile struct {
	Path string
	AST  model.ASTNode
}

func Parser(JsSource []collector.SourceFile) ([]ParsedFile, error) {
	JSParsedFiles := []ParsedFile{}

	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	parserScript := filepath.Join(filepath.Dir(execPath), "internal", "parser-side", "parser.js")

	if _, err := os.Stat(parserScript); err != nil {
		parserScript = filepath.Join("internal", "parser-side", "parser.js")
	}

	for _, fileData := range JsSource {
		ParserOutput, err := exec.Command("node", parserScript, fileData.Path).Output()
		if err != nil {
			fmt.Printf("[e]failed to parse a file: %s error: %v", fileData.Path, err)
			continue
		}

		var ast any
		err = json.Unmarshal(ParserOutput, &ast)
		if err != nil {
			fmt.Println("[e]Failed to decode json data")
			return nil, err
		}

		// Normalize here
		NormalizedAST, err := Normalizer(ast)
		if err != nil {
			return nil, err
		}

		parsedFile := ParsedFile{
			Path: fileData.Path,
			AST:  NormalizedAST,
		}

		JSParsedFiles = append(JSParsedFiles, parsedFile)
	}
	return JSParsedFiles, nil
}
