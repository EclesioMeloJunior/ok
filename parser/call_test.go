package parser_test

import (
	"fmt"
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/parser"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	for testName, test := range map[string]struct {
		str      string
		expected *ast.Call
	}{
		"no-args": {
			str: "foo()",
			expected: &ast.Call{
				FunctionName: "foo",
			},
		},
		"one-arg": {
			str: `bar("baz")`,
			expected: &ast.Call{
				FunctionName: "bar",
				Arguments: []ast.Node{
					asttest.NewLiteralString("baz"),
				},
			},
		},
		"math-abs": {
			str: `math.abs(123)`,
			expected: &ast.Call{
				FunctionName: "math.abs",
				Arguments: []ast.Node{
					asttest.NewLiteralNumber("123"),
				},
			},
		},
		"cast-string": {
			str: `string 'a'`,
			expected: &ast.Call{
				FunctionName: "string",
				Arguments: []ast.Node{
					asttest.NewLiteralChar('a'),
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
