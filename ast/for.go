package ast

// For represents a for loop.
type For struct {
	// All of Init, Condition and Next may be nil.
	Init, Condition, Next Node

	// Statements may be nil.
	Statements []Node
}

// In represents an "in" expression in for loops.
type In struct {
	Key, Value string
	Expr       Node
}
