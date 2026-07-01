package analyzer

import (
	"afv/internal/model"
	"strings"
)

// whats the main purpose of this file? to extract authentication flow graph from whole program graph
// input is program graph including:
// Nodes: functions, endpoints, external calls
// Edges: ralation between nodes
// output has also the same type
func BuildAuthFlowGraph(g model.Graph) model.Graph {
	ag := model.Graph{ // create an empty graph from model
		Nodes: []model.Node{},
		Edges: []model.Edge{},
	}

	adj := map[string][]string{} // to store all routes from a node
	// a -> b
	// a -> c
	// b -> d
	// will be:
	// a: [b, c]
	// b: [c]
	nodeMap := map[string]model.Node{} // to reach from ID to comlete node.
	visited := map[string]bool{}       // to avoid infinite loops

	for _, edge := range g.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To) // for every edge, store the routes to adj
		// for example for 'login -> validate' and 'login -> fetch' we will have:
		// adj["login"] = ["validate", "fetch"]
	}

	for _, node := range g.Nodes {
		nodeMap[node.ID] = node // filling nodeMaps with main graph data
	}

	for _, node := range g.Nodes {

		if strings.Contains(node.ID, "login") { // starting point, check all IDs to findout which one contains login keyword

			queue := []string{node.ID} // initialize with node.ID

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]

				if visited[current] {
					continue
				}
				visited[current] = true // mark visited

				n, ok := nodeMap[current]
				if ok {
					ag.Nodes = append(ag.Nodes, n) // adding node to AuthFlowGraph.Nodes
				}

				for _, neighbor := range adj[current] { // where can i go from this node? (all possible routes, one at a time). validate for example as first value of list
					queue = append(queue, neighbor) // add it queue

					ag.Edges = append(ag.Edges, model.Edge{
						From: current,  // login e.g.
						To:   neighbor, // validate
						Type: "auth_flow",
					})
				}
			}
		}
	}

	return ag // return the graph
}
