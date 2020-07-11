package ast

// On a error handler for an ErrorScope.
type On struct {
	// Type is the type name, like "MyError".
	Type string

	// Statement may be nil.
	Statements []Node

	Pos string
}

// Position returns the position.
func (node *On) Position() string {
	return node.Pos
}

// ErrorScope represents the try/on error scope.
type ErrorScope struct {
	// Statements is what will be run in this scope. It is allowed to be nil.
	Statements []Node

	On  []*On
	Pos string
}

// Position returns the position.
func (node *ErrorScope) Position() string {
	return node.Pos
}
