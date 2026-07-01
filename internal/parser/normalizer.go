package parser

import (
	"afv/internal/model"
	"errors"
	"sort"
)

func Normalizer(rawAST any) (model.ASTNode, error) {
	return normalizeWithRole(rawAST, "")
}

func normalizeWithRole(rawAST any, role string) (model.ASTNode, error) {

	ast, ok := rawAST.(map[string]any)
	if !ok {
		return model.ASTNode{}, errors.New("invalid AST node")
	}

	var node model.ASTNode

	t, ok := ast["type"].(string)
	if !ok {
		return node, errors.New("node has no type")
	}

	node.Type = t
	node.Role = role

	if name, ok := ast["name"].(string); ok {
		node.Name = name
	}

	if value, exists := ast["value"]; exists {
		node.Value = value
	}

	if loc, ok := ast["loc"].(map[string]any); ok {
		if start, ok := loc["start"].(map[string]any); ok {
			if line, ok := start["line"].(float64); ok {
				node.Location.Start.Line = uint32(line)
			}
			if column, ok := start["column"].(float64); ok {
				node.Location.Start.Column = uint32(column)
			}
		}
		if end, ok := loc["end"].(map[string]any); ok {
			if line, ok := end["line"].(float64); ok {
				node.Location.End.Line = uint32(line)
			}
			if column, ok := end["column"].(float64); ok {
				node.Location.End.Column = uint32(column)
			}
		}
	}

	keys := make([]string, 0, len(ast))
	for k := range ast {
		keys = append(keys, k) // store AST field names in a slice of strings (like json keys)
	}
	sort.Strings(keys) // sort that slice of strings

	for _, k := range keys { // for every field...
		value := ast[k] // get the value of this field

		if childObj, ok := value.(map[string]any); ok { // is our field value a node itself? like callee
			if _, exists := childObj["type"]; exists { // if has a type...
				child, err := normalizeWithRole(childObj, k) // run function on it again (recursion)
				if err == nil {
					node.Children = append(node.Children, child)
				}
			}
			continue // move on to next key
		}

		if arr, ok := value.([]any); ok { // if value is an array...
			for _, item := range arr { // for each item in array...
				childObj, ok := item.(map[string]any) // check if our item is a node
				if !ok {
					continue // if not, move on to next field
				}
				if _, exists := childObj["type"]; !exists {
					continue // if has no type, move on to next field
				}
				child, err := normalizeWithRole(childObj, k) // if is a node and has a type, run funciton on it again (recursion)
				if err == nil {
					node.Children = append(node.Children, child) // if there was no error, append child to node children slice
				}
			}
		}
	}

	return node, nil
}
