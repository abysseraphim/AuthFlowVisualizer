package analyzer

import (
	"afv/internal/model"
)

func ExtractCallgraphs(root model.ASTNode) []model.Call {
	calls := []model.Call{}

	Walk(root, nil, func(node model.ASTNode, parent *model.ASTNode) {
		if !isFunction(node) { // only operate on functions
			return
		}

		caller := getFunctionName(node, parent)
		if caller == "" {
			return
		}

		var body *model.ASTNode
		for i := range node.Children {
			child := &node.Children[i]

			if child.Type == "BlockStatement" && child.Role == "body" { // blockstatement of a function, has role: body
				body = child
				break
			}
		}

		if body == nil {
			return
		}

		walkCalls(*body, caller, &calls)

	})

	return calls
}

func walkCalls(node model.ASTNode, caller string, calls *[]model.Call) {
	if isFunction(node) {
		return
	}

	if node.Type == "CallExpression" {
		callee := ""
		for _, child := range node.Children {
			if child.Role != "callee" {
				continue
			}

			if child.Type == "Identifier" {
				callee = child.Name
			}
			if child.Type == "MemberExpression" {
				var object string
				var prop string
				for _, member := range child.Children {
					if member.Role == "object" {
						object = member.Name
					}
					if member.Role == "property" {
						prop = member.Name
					}
				}
				callee = object + "." + prop
			}

		}

		if callee != "" {
			*calls = append(*calls, model.Call{
				Caller:   caller,
				Callee:   callee,
				Location: node.Location,
			})

			for _, child := range node.Children {
				walkCalls(child, caller, calls)
			}
			return
		}
	}

	for _, child := range node.Children {
		walkCalls(child, caller, calls)
	}
}
