package vm_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/vm"

	"github.com/stretchr/testify/assert"
)

func TestPrint_Execute(t *testing.T) {
	for testName, test := range map[string]struct {
		values         []*ast.Literal
		expectedStdout string
	}{
		"no-args": {nil, "\n"},
		"true": {
			[]*ast.Literal{{Kind: "bool", Value: "true"}},
			"true\n",
		},
		"false": {
			[]*ast.Literal{{Kind: "bool", Value: "false"}},
			"false\n",
		},
		"char": {
			[]*ast.Literal{{Kind: "char", Value: "#"}},
			"#\n",
		},
		"data": {
			[]*ast.Literal{{Kind: "data", Value: "abc"}},
			"abc\n",
		},
		"number": {
			[]*ast.Literal{{Kind: "number", Value: "1.23"}},
			"1.23\n",
		},
		"string": {
			[]*ast.Literal{{Kind: "string", Value: "foo bar"}},
			"foo bar\n",
		},
		"multiple-args": {
			[]*ast.Literal{
				{Kind: "string", Value: "foo"},
				{Kind: "number", Value: "123"},
			},
			"foo 123\n",
		},
		"number-array": {
			[]*ast.Literal{
				{
					Kind: "[]number",
					Array: []*ast.Literal{
						asttest.NewLiteralNumber("123"),
						asttest.NewLiteralNumber("456"),
						asttest.NewLiteralNumber("789"),
					},
				},
			},
			"[123, 456, 789]\n",
		},
		"any-array": {
			[]*ast.Literal{
				{
					Kind: "[]any",
					Array: []*ast.Literal{
						asttest.NewLiteralBool(true),
						asttest.NewLiteralChar('a'),
						asttest.NewLiteralData([]byte("data")),
						asttest.NewLiteralNumber("123"),
						asttest.NewLiteralString("789"),
					},
				},
			},
			"[true, \"a\", \"data\", 123, \"789\"]\n",
		},
		"number-map": {
			[]*ast.Literal{
				{
					Kind: "{}number",
					Map: map[string]*ast.Literal{
						"a": asttest.NewLiteralNumber("123"),
						"b": asttest.NewLiteralNumber("456"),
						"c": asttest.NewLiteralNumber("789"),
					},
				},
			},
			"{\"a\": 123, \"b\": 456, \"c\": 789}\n",
		},
		"any-map": {
			[]*ast.Literal{
				{
					Kind: "{}any",
					Map: map[string]*ast.Literal{
						"a": asttest.NewLiteralBool(true),
						"b": asttest.NewLiteralChar('a'),
						"c": asttest.NewLiteralData([]byte("data")),
						"d": asttest.NewLiteralNumber("123"),
						"e": asttest.NewLiteralString("789"),
					},
				},
			},
			"{\"a\": true, \"b\": \"a\", \"c\": \"data\", \"d\": 123, \"e\": \"789\"}\n",
		},
		"Person": {
			[]*ast.Literal{
				{
					Kind: "Person",
					Map: map[string]*ast.Literal{
						"Foo": asttest.NewLiteralNumber("123"),
						"bar": asttest.NewLiteralNumber("456"),
					},
				},
			},
			"{\"Foo\": 123}\n",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			registers := map[vm.Register]*ast.Literal{}
			var arguments []vm.Register
			for i, value := range test.values {
				register := vm.Register(fmt.Sprintf("%d", i))
				registers[register] = value
				arguments = append(arguments, register)
			}

			buf := bytes.NewBuffer(nil)
			ins := &vm.Print{
				Arguments: arguments,
			}
			vm := &vm.VM{
				Stack:  []map[vm.Register]*ast.Literal{registers},
				Stdout: buf,
			}
			assert.NoError(t, ins.Execute(nil, vm))
			assert.Equal(t, test.expectedStdout, buf.String())
		})
	}
}
