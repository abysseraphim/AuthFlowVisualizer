package collector

import (
	"fmt"
	"os"
)

type SourceFile struct {
	Path    string
	Content []byte
}

func FileCollector(paths []string) ([]SourceFile, error) {
	sources := []SourceFile{}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("warning: failed to read file: %q: %v\n", path, err)
			continue
		}

		sourceFile := SourceFile{
			Path:    path,
			Content: data,
		}
		sources = append(sources, sourceFile)

	}
	return sources, nil
}
