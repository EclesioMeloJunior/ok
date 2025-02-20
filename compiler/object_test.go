package compiler_test

import (
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/compiler"
	"github.com/elliotchance/ok/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	for testName, test := range map[string]struct {
		node     *ast.Func
		expected []vm.Instruction
		err      error
	}{
		"empty-object": {
			node: &ast.Func{
				Name:    "Foo",
				Returns: []string{"Foo"},
			},
			expected: []vm.Instruction{
				// return instance
				&vm.Return{
					Results: []vm.Register{vm.StateRegister},
				},
			},
		},
		"one-variable": {
			node: &ast.Func{
				Name:    "Foo",
				Returns: []string{"Foo"},
				Statements: []ast.Node{
					&ast.Assign{
						Lefts: []ast.Node{
							&ast.Identifier{Name: "bar"},
						},
						Rights: []ast.Node{
							asttest.NewLiteralNumber("123"),
						},
					},
				},
			},
			expected: []vm.Instruction{
				// bar = 123
				&vm.Assign{
					VariableName: "1",
					Value:        asttest.NewLiteralNumber("123"),
				},
				&vm.Assign{
					VariableName: "bar",
					Register:     "1",
				},

				// return instance
				&vm.Return{
					Results: []vm.Register{vm.StateRegister},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			compiledFunc, err := compiler.CompileFunc(test.node,
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
