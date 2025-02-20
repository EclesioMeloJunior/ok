package compiler_test

import (
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/compiler"
	"github.com/elliotchance/ok/lexer"
	"github.com/elliotchance/ok/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnary(t *testing.T) {
	for testName, test := range map[string]struct {
		nodes    []ast.Node
		expected []vm.Instruction
		err      error
	}{
		"increment-variable": {
			nodes: []ast.Node{
				&ast.Assign{
					Lefts: []ast.Node{
						&ast.Identifier{Name: "i"},
					},
					Rights: []ast.Node{
						asttest.NewLiteralNumber("0"),
					},
				},
				&ast.Unary{
					Op:   "++",
					Expr: &ast.Identifier{Name: "i"},
				},
			},
			expected: []vm.Instruction{
				&vm.Assign{
					VariableName: "1",
					Value:        asttest.NewLiteralNumber("0"),
				},
				&vm.Assign{
					VariableName: "i",
					Register:     "1",
				},
				&vm.Assign{
					VariableName: "2",
					Value:        asttest.NewLiteralNumber("1"),
				},
				&vm.Add{
					Left:   "i",
					Right:  "2",
					Result: "i",
				},
			},
		},
		"decrement-variable": {
			nodes: []ast.Node{
				&ast.Assign{
					Lefts: []ast.Node{
						&ast.Identifier{Name: "i"},
					},
					Rights: []ast.Node{
						asttest.NewLiteralNumber("0"),
					},
				},
				&ast.Unary{
					Op:   "--",
					Expr: &ast.Identifier{Name: "i"},
				},
			},
			expected: []vm.Instruction{
				&vm.Assign{
					VariableName: "1",
					Value:        asttest.NewLiteralNumber("0"),
				},
				&vm.Assign{
					VariableName: "i",
					Register:     "1",
				},
				&vm.Assign{
					VariableName: "2",
					Value:        asttest.NewLiteralNumber("1"),
				},
				&vm.Subtract{
					Left:   "i",
					Right:  "2",
					Result: "i",
				},
			},
		},
		"not-true": {
			nodes: []ast.Node{
				&ast.Unary{
					Op:   lexer.TokenNot,
					Expr: asttest.NewLiteralBool(true),
				},
			},
			expected: []vm.Instruction{
				&vm.Assign{
					VariableName: "1",
					Value:        asttest.NewLiteralBool(true),
				},
				&vm.Not{
					Left:   "1",
					Result: "2",
				},
			},
		},
		"not-false": {
			nodes: []ast.Node{
				&ast.Unary{
					Op:   lexer.TokenNot,
					Expr: asttest.NewLiteralBool(false),
				},
			},
			expected: []vm.Instruction{
				&vm.Assign{
					VariableName: "1",
					Value:        asttest.NewLiteralBool(false),
				},
				&vm.Not{
					Left:   "1",
					Result: "2",
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			compiledFunc, err := compiler.CompileFunc(newFunc(test.nodes...),
				&compiler.Compiled{})
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, compiledFunc.Instructions)
			}
		})
	}
}
