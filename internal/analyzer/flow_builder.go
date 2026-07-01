package analyzer

import (
	"afv/internal/model"
)

func BuildCallGraph(p model.Program) model.Graph {
	g := model.Graph{
		Nodes: []model.Node{},
		Edges: []model.Edge{},
	}

	nodeMap := map[string]bool{}
	funcMap := map[string]bool{}

	Files := p.Files
	for _, file := range Files {
		for _, function := range file.Functions {
			node := model.Node{}
			id := "func:" + file.Path + ":" + function.Name
			node.ID = id

			node.Type = "function"

			if nodeMap[id] == false {
				g.Nodes = append(g.Nodes, node)
				nodeMap[id] = true
			}

			funcMap[file.Path+":"+function.Name] = true
		}

		for _, endpoint := range file.Endpoints {
			node := model.Node{}
			id := "endpoint:" + endpoint.Method + ":" + endpoint.URL
			node.ID = id

			node.Type = "endpoint"

			if nodeMap[id] == false {
				g.Nodes = append(g.Nodes, node)
				nodeMap[id] = true
			}

			ownerName := ""
			bestSpan := -1
			for _, function := range file.Functions {
				if !locContains(function.Location, endpoint.Location) {
					continue
				}

				span := int(function.Location.End.Line) - int(function.Location.Start.Line)
				if bestSpan == -1 || span < bestSpan {
					bestSpan = span
					ownerName = function.Name
				}
			}

			if ownerName != "" {
				g.Edges = append(g.Edges, model.Edge{
					From: "func:" + file.Path + ":" + ownerName,
					To:   id,
					Type: "endpoint",
				})
			}
		}

		for _, call := range file.Calls {
			edge := model.Edge{}

			edge.From = "func:" + file.Path + ":" + call.Caller
			edge.Type = "call"

			if _, ok := funcMap[file.Path+":"+call.Callee]; ok {
				edge.To = "func:" + file.Path + ":" + call.Callee
			} else {
				id := "external:" + call.Callee

				if !nodeMap[id] {
					node := model.Node{
						ID:   id,
						Type: "external",
					}

					g.Nodes = append(g.Nodes, node)
					nodeMap[id] = true
				}

				edge.To = id
			}

			g.Edges = append(g.Edges, edge)
		}

	}

	return g
}

func locContains(outer model.Loc, inner model.Loc) bool {
	if outer.Start.Line > inner.Start.Line {
		return false
	}
	if outer.Start.Line == inner.Start.Line && outer.Start.Column > inner.Start.Column {
		return false
	}
	if outer.End.Line < inner.End.Line {
		return false
	}
	if outer.End.Line == inner.End.Line && outer.End.Column < inner.End.Column {
		return false
	}
	return true
}
