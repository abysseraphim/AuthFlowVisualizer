package input

import (
	"afv/internal/collector"
	"errors"
	"fmt"
)

func InputHandler(target string, inputType string) ([]collector.SourceFile, error) {
	switch inputType {
	case "path":
		JsFiles, err := LocalHandler(target)
		if err != nil {
			return nil, err
		}
		fmt.Println("[*]JavaScript Source Files:")
		for _, file := range JsFiles {
			println(file)
		}

		jsSources, err := collector.FileCollector(JsFiles)
		if err != nil {
			fmt.Println("[e]error:", err)
			return nil, err
		}
		return jsSources, nil

	case "url":
		TempFiles, err := URLHandler(target)
		if err != nil {
			return nil, err
		}
		fmt.Println("[*]TEMP FILES:")
		for _, tmpfile := range TempFiles {
			fmt.Println(tmpfile)
		}

		jsSources, err := collector.FileCollector(TempFiles)
		if err != nil {
			fmt.Println("[e]error:", err)
			return nil, err
		}
		return jsSources, nil

	default:
		fmt.Println("[e]Unknown Type, exitting function")
		return nil, errors.New("[e]unknown input type")
	}
}
