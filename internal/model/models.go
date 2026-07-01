package model

type ASTNode struct {
	Type     string
	Role     string
	Name     string
	Value    any
	Children []ASTNode
	Location Loc
}

type Position struct {
	Line   uint32
	Column uint32
}

type Loc struct {
	Start Position
	End   Position
}
