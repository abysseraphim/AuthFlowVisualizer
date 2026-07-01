package report

import (
	"afv/internal/model"
)

func GenerateReport(target string, programGraph model.Graph, authGraph model.Graph) error {

	GenerateCLIReport(target, programGraph, authGraph)
	GenerateHTMLReport(target, programGraph, authGraph)

	return nil
}
