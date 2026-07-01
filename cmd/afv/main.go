package main

import (
	"afv/internal/analyzer"
	"afv/internal/input"
	"afv/internal/parser"
	"afv/internal/report"
	"flag"
	"fmt"
)

func main() {

	fmt.Print(`

	   █████████   ███████████ █████   █████
	  ███░░░░░███ ░░███░░░░░░█░░███   ░░███ 
	 ░███    ░███  ░███   █ ░  ░███    ░███ 
	 ░███████████  ░███████    ░███    ░███ 
	 ░███░░░░░███  ░███░░░█    ░░███   ███  
	 ░███    ░███  ░███  ░      ░░░█████░   
	 █████   █████ █████          ░░███     
	░░░░░   ░░░░░ ░░░░░            ░░░      
                                        
        Auth Flow Visualizer: v 1.0.0
          by abysseraphim github
	`)
	fmt.Println()
	targetPath := flag.String("p", "", "Source Code Path")
	targetURL := flag.String("u", "", "Target URL")

	flag.Parse()

	target := ""
	var targetType string

	if *targetPath != "" {
		fmt.Println("[*]Target Path Selected:", *targetPath)
		target = *targetPath
		targetType = "path"
	} else if *targetURL != "" {
		fmt.Println("[*]Target URL Selected:", *targetURL)
		target = *targetURL
		targetType = "url"
	} else {
		fmt.Println("You have to specify At least one Target. Usage:")
		fmt.Println("\"./afv -p /path/to/source/code\"  OR  \"./afv -u https://targetRUL.tld\"")
		return
	}

	fmt.Println("[*]Target loaded:", target, "Type:", targetType)
	jsSources, err := input.InputHandler(target, targetType)
	if err != nil {
		fmt.Println("[e]error occured:", err)
		return
	}
	// fmt.Println("JS file contents:")
	// fmt.Println(jsSources)
	parsedJs, err := parser.Parser(jsSources)
	// fmt.Println(os.Getwd())
	if err != nil {
		fmt.Println("[e]Faled to parse JS:", err)
		return
	}
	// fmt.Println("Parsed JS:")
	// fmt.Println(parsedJs)

	analyzed, err := analyzer.Analyze(parsedJs)
	// fmt.Println("ANALYZER OUTPUT:")
	// fmt.Println(analyzed)

	callFlow := analyzer.BuildCallGraph(analyzed)
	// fmt.Println("Program Call Flow:")
	// fmt.Println(callFlow)

	authFlow := analyzer.BuildAuthFlowGraph(callFlow)
	// fmt.Println("Auth Call Flow:")
	// fmt.Println(authFlow)

	report.GenerateReport(target, callFlow, authFlow)

}
