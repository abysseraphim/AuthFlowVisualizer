package analyzer

import (
	"afv/internal/model"
)

func ExtractFunctions(root model.ASTNode) []model.Function {
	functions := []model.Function{}

	Walk(root, nil, func(node model.ASTNode, parent *model.ASTNode) {
		if !isFunction(node) {
			return
		}
		funcData := extractFunction(node, parent)

		functions = append(functions, funcData)
	})

	return functions
}

func isFunction(node model.ASTNode) bool {
	switch node.Type {
	case "FunctionDeclaration":
		return true
	case "FunctionExpression":
		return true
	case "ArrowFunctionExpression":
		return true
	case "ClassMethod":
		return true
	case "ObjectMethod":
		return true
	default:
		return false
	}
}

func extractFunction(root model.ASTNode, parent *model.ASTNode) model.Function {
	funcData := model.Function{}

	location := root.Location
	funcData.Location = location

	kind := root.Type
	funcData.Kind = kind

	name := getFunctionName(root, parent)
	funcData.Name = name

	return funcData
}

func getFunctionName(node model.ASTNode, parent *model.ASTNode) string {
	if node.Type == "FunctionDeclaration" {
		for _, child := range node.Children {
			name := child.Name
			if child.Type == "Identifier" && child.Role == "id" && name != "" {
				return name
			}
		}
	} else if node.Type == "ClassMethod" || node.Type == "ObjectMethod" {
		for _, child := range node.Children {
			name := child.Name
			if child.Type == "Identifier" && child.Role == "key" && name != "" {
				return name
			}
		}
	} else if node.Type == "FunctionExpression" || node.Type == "ArrowFunctionExpression" {
		if parent != nil && parent.Type == "VariableDeclarator" { // function is declared inside a variable
			for _, child := range parent.Children {
				if child.Type == "Identifier" && child.Role == "id" && child.Name != "" {
					return child.Name
				}
			}
		}
	}

	return ""
}
