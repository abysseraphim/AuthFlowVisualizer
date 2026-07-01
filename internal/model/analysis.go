package model

type Program struct {
	Files []ProgramFile
}

type ProgramFile struct {
	Path      string
	Functions []Function
	Endpoints []Endpoint
	Calls     []Call
}

type Function struct {
	Name     string
	Kind     string
	Location Loc
}

type Endpoint struct {
	URL      string
	Method   string
	Location Loc
}

type Call struct {
	Caller   string
	Callee   string
	Location Loc
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID   string
	Type string
}

type Edge struct { // by Edge i mean relation between nodes
	From string
	To   string
	Type string
}
