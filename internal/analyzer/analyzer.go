package analyzer

import (
	"afv/internal/model"
	"afv/internal/parser"
)

func Analyze(files []parser.ParsedFile) (model.Program, error) {
	program := model.Program{}
	for _, file := range files {
		programFile := model.ProgramFile{}
		programFile.Path = file.Path
		programFile.Functions = ExtractFunctions(file.AST)
		programFile.Endpoints = ExtractEndpoints(file.AST)
		programFile.Calls = ExtractCallgraphs(file.AST)

		program.Files = append(program.Files, programFile)
	}
	return program, nil
}
