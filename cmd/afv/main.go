package main

import (
	"afv/internal/analyzer"
	"afv/internal/input"
	"afv/internal/parser"
	"afv/internal/progress"
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
                                        
        Auth Flow Visualizer: v 1.1.0
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
		fmt.Println("\"./afv -p /path/to/source/code\"  OR  \"./afv -u https://targetURL.tld\"")
		return
	}

	fmt.Println("[*]Target loaded:", target, "Type:", targetType)

	// Stage 1: Collect JS sources
	spin := progress.New("Collecting JS sources")
	jsSources, err := input.InputHandler(target, targetType)
	if err != nil {
		spin.Stop("failed.")
		fmt.Println("[e]error occured:", err)
		return
	}
	spin.Stop(fmt.Sprintf("done. (%d files)", len(jsSources)))

	// Stage 2: Parse JS
	spin = progress.New("Parsing JavaScript")
	parsedJs, err := parser.Parser(jsSources)
	if err != nil {
		spin.Stop("failed.")
		fmt.Println("[e]Failed to parse JS:", err)
		return
	}
	spin.Stop(fmt.Sprintf("done. (%d files parsed)", len(parsedJs)))

	// Stage 3: Analyze AST
	spin = progress.New("Analyzing AST")
	analyzed, err := analyzer.Analyze(parsedJs)
	if err != nil {
		spin.Stop("failed.")
		fmt.Println("[e]Analysis failed:", err)
		return
	}
	spin.Stop("done.")

	// Stage 4: Build call graph
	spin = progress.New("Building call graph")
	callFlow := analyzer.BuildCallGraph(analyzed)
	spin.Stop(fmt.Sprintf("done. (%d nodes, %d edges)", len(callFlow.Nodes), len(callFlow.Edges)))

	// Stage 5: Build auth flow
	spin = progress.New("Extracting auth flow")
	authFlow := analyzer.BuildAuthFlowGraph(callFlow)
	spin.Stop(fmt.Sprintf("done. (%d auth nodes)", len(authFlow.Nodes)))

	// Stage 6: Generate report
	spin = progress.New("Generating report")
	report.GenerateReport(target, callFlow, authFlow)
	spin.Stop("done.")
}
