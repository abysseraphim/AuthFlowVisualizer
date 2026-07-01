package analyzer

import "afv/internal/model"

// Walk on a Normalized AST.
func Walk(node model.ASTNode, parent *model.ASTNode, visitor func(model.ASTNode, *model.ASTNode)) {
	visitor(node, parent)
	for _, child := range node.Children {
		Walk(child, &node, visitor)
	}
}
