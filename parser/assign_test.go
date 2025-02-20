package parser_test

import (
	"fmt"
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/lexer"
	"github.com/elliotchance/ok/parser"

	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	for testName, test := range map[string]struct {
		str      string
		expected *ast.Assign
	}{
		"single-literal": {
			str: "a = 123",
			expected: &ast.Assign{
				Lefts: []ast.Node{
					&ast.Identifier{Name: "a"},
				},
				Rights: []ast.Node{
					asttest.NewLiteralNumber("123"),
				},
			},
		},
		"double-literal": {
			str: `a, b = 123, "foo"`,
			expected: &ast.Assign{
				Lefts: []ast.Node{
					&ast.Identifier{Name: "a"},
					&ast.Identifier{Name: "b"},
				},
				Rights: []ast.Node{
					asttest.NewLiteralNumber("123"),
					asttest.NewLiteralString("foo"),
				},
			},
		},
		"variable": {
			str: `foo = bar`,
			expected: &ast.Assign{
				Lefts: []ast.Node{
					&ast.Identifier{Name: "foo"},
				},
				Rights: []ast.Node{
					&ast.Identifier{Name: "bar"},
				},
			},
		},
		"variable-expr": {
			str: `foo = len(bar) - 1`,
			expected: &ast.Assign{
				Lefts: []ast.Node{
					&ast.Identifier{Name: "foo"},
				},
				Rights: []ast.Node{
					asttest.NewBinary(
						&ast.Call{
							FunctionName: "len",
							Arguments: []ast.Node{
								&ast.Identifier{Name: "bar"},
							},
						},
						lexer.TokenMinus,
						asttest.NewLiteralNumber("1"),
					),
				},
			},
		},
		"assign-anon-func": {
			str: `fn = func() {}`,
			expected: &ast.Assign{
				Lefts: []ast.Node{
					&ast.Identifier{Name: "fn"},
				},
				Rights: []ast.Node{
					&ast.Func{Name: "1"},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			str := fmt.Sprintf("func main() { %s }", test.str)
			p := parser.ParseString(str, "a.ok")

			assert.Nil(t, p.Errors())
			asttest.AssertEqual(t, map[string]*ast.Func{
				"main": newFunc(test.expected),
			}, p.File.Funcs)
			assert.Nil(t, p.File.Comments)
		})
	}
}
