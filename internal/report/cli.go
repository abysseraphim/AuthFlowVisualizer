package report

import (
	"afv/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateCLIReport(target string, graph model.Graph, auth model.Graph) error {

	err := os.MkdirAll(filepath.Join("output", target), 0755)
	if err != nil {
		return err
	}
	reportFile := filepath.Join("output", target, "report.txt")

	var builder strings.Builder

	builder.WriteString("==============================\n")
	builder.WriteString("     AUTH FLOW VISUALIZER\n")
	builder.WriteString("==============================\n\n")

	builder.WriteString("Target:\n")
	builder.WriteString(target)
	builder.WriteString("\n\n")

	builder.WriteString("********************************\n")
	builder.WriteString("PROGRAM FLOW\n")
	builder.WriteString("********************************\n\n")

	for _, node := range graph.Nodes {
		if node.Type == "function" {
			builder.WriteString("Function:\n")
			sof := strings.Split(node.ID, ":")
			builder.WriteString(sof[2] + "()\n\n")

			builder.WriteString("Calls:\n")

			for _, edge := range graph.Edges {
				if edge.From == node.ID {
					sot := strings.Split(edge.To, ":")
					if sot[0] == "external" {
						builder.WriteString("  |--> " + sot[1] + "()\n")
					} else if sot[0] == "func" {
						builder.WriteString("  |--> " + sot[len(sot)-1] + "()\n")
					}

				}
			}
			builder.WriteString("\n")

			for _, edge := range graph.Edges {
				if edge.From != node.ID {
					continue
				}

				if !strings.HasPrefix(edge.To, "endpoint:") {
					continue
				}

				ep := strings.Split(edge.To, ":")

				if len(ep) < 3 {
					continue
				}

				builder.WriteString("Endpoint:\n")
				builder.WriteString("  " + ep[1] + " " + ep[2] + "\n\n")
			}

			builder.WriteString("--------------------------------\n\n")
		}
	}

	builder.WriteString("********************************\n")
	builder.WriteString("AUTH FLOW\n")
	builder.WriteString("********************************\n\n")
	for _, node := range auth.Nodes {
		if node.Type == "function" {
			builder.WriteString("Entry:\n")
			sof := strings.Split(node.ID, ":")
			builder.WriteString(sof[2] + "()\n\n")

			builder.WriteString("Calls:\n")

			for _, edge := range auth.Edges {
				if edge.From == node.ID {
					sot := strings.Split(edge.To, ":")
					if len(sot) >= 2 && sot[0] == "external" {
						builder.WriteString("  |--> " + sot[1] + "()\n")
					}
				}
			}

			builder.WriteString("\n")

			for _, edge := range auth.Edges {
				if edge.From != node.ID {
					continue
				}

				if !strings.HasPrefix(edge.To, "endpoint:") {
					continue
				}

				ep := strings.Split(edge.To, ":")

				if len(ep) < 3 {
					continue
				}

				builder.WriteString("Endpoint:\n")
				builder.WriteString("  " + ep[1] + " " + ep[2] + "\n\n")
			}

			builder.WriteString("--------------------------------\n\n")
		}
	}
	builder.WriteString("********************************\n")
	builder.WriteString("SUMMARY\n")
	builder.WriteString("********************************\n\n")

	functions := 0
	endpoints := 0

	for _, node := range graph.Nodes {
		switch node.Type {
		case "function":
			functions++
		case "endpoint":
			endpoints++
		}
	}

	builder.WriteString(fmt.Sprintf("Functions : %d\n", functions))
	builder.WriteString(fmt.Sprintf("Endpoints : %d\n", endpoints))
	builder.WriteString(fmt.Sprintf("Calls     : %d\n\n", len(graph.Edges)))

	builder.WriteString("[✔] Analysis completed successfully.\n")

	content := builder.String()

	err = os.WriteFile(reportFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	fmt.Println("[*]Report generated:", reportFile)
	return nil
}
