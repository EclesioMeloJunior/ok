package parser_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/lexer"
	"github.com/elliotchance/ok/parser"
	"github.com/stretchr/testify/assert"
)

func TestExpr(t *testing.T) {
	for testName, test := range map[string]struct {
		str      string
		expected ast.Node
		errs     []error
	}{
		"literal-true": {
			str: `true`,
			expected: &ast.Literal{
				Kind:  "bool",
				Value: "true",
			},
		},
		"literal-false": {
			str: `false`,
			expected: &ast.Literal{
				Kind:  "bool",
				Value: "false",
			},
		},
		"literal-char": {
			str: `'a'`,
			expected: &ast.Literal{
				Kind:  "char",
				Value: "a",
			},
		},
		"literal-zero-length-char": {
			str: `''`,
			expected: &ast.Literal{
				Kind:  "char",
				Value: "",
			},
			errs: []error{
				errors.New("character literal cannot be empty"),
			},
		},
		"literal-number-zero": {
			str: `0`,
			expected: &ast.Literal{
				Kind:  "number",
				Value: "0",
			},
		},
		"literal-number-negative": {
			str: `-3.20`,
			expected: &ast.Unary{
				Op:   lexer.TokenMinus,
				Expr: ast.NewLiteralNumber("3.20"),
			},
		},
		"numbers-plus": {
			str: `3 + 2`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralNumber("3"),
				Op:    lexer.TokenPlus,
				Right: ast.NewLiteralNumber("2"),
			},
		},
		"numbers-minus": {
			str: `3 - 2`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralNumber("3"),
				Op:    lexer.TokenMinus,
				Right: ast.NewLiteralNumber("2"),
			},
		},
		"numbers-times": {
			str: `3.0*2.1`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralNumber("3.0"),
				Op:    lexer.TokenTimes,
				Right: ast.NewLiteralNumber("2.1"),
			},
		},
		"numbers-divide": {
			str: `3/2.0`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralNumber("3"),
				Op:    lexer.TokenDivide,
				Right: ast.NewLiteralNumber("2.0"),
			},
		},
		"numbers-remainder": {
			str: `3 % 2`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralNumber("3"),
				Op:    lexer.TokenRemainder,
				Right: ast.NewLiteralNumber("2"),
			},
		},
		"expr-3-linear-order": {
			str: `1 + 2 - 3`,
			expected: &ast.Binary{
				Left: ast.NewLiteralNumber("1"),
				Op:   lexer.TokenPlus,
				Right: &ast.Binary{
					Left:  ast.NewLiteralNumber("2"),
					Op:    lexer.TokenMinus,
					Right: ast.NewLiteralNumber("3"),
				},
			},
		},
		"expr-3-non-linear-order": {
			str: `1 * 2 - 3`,
			expected: &ast.Binary{
				Left: &ast.Binary{
					Left:  ast.NewLiteralNumber("1"),
					Op:    lexer.TokenTimes,
					Right: ast.NewLiteralNumber("2"),
				},
				Op:    lexer.TokenMinus,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"expr-3-grouping": {
			str: `1 * (2 - 3)`,
			expected: &ast.Binary{
				Left: ast.NewLiteralNumber("1"),
				Op:   lexer.TokenTimes,
				Right: &ast.Group{
					Expr: &ast.Binary{
						Left:  ast.NewLiteralNumber("2"),
						Op:    lexer.TokenMinus,
						Right: ast.NewLiteralNumber("3"),
					},
				},
			},
		},
		"group-1": {
			str: `(2 - 3)`,
			expected: &ast.Group{
				Expr: &ast.Binary{
					Left:  ast.NewLiteralNumber("2"),
					Op:    lexer.TokenMinus,
					Right: ast.NewLiteralNumber("3"),
				},
			},
		},
		"group-2": {
			str: `(2)`,
			expected: &ast.Group{
				Expr: ast.NewLiteralNumber("2"),
			},
		},
		"bool-and-bool": {
			str: `true and false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenAnd,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-or-bool": {
			str: `true or false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenOr,
				Right: ast.NewLiteralBool(false),
			},
		},
		"not-bool": {
			str: `not true`,
			expected: &ast.Unary{
				Op:   lexer.TokenNot,
				Expr: ast.NewLiteralBool(true),
			},
		},
		"not-not-bool": {
			str: `not not false`,
			expected: &ast.Unary{
				Op: lexer.TokenNot,
				Expr: &ast.Unary{
					Op:   lexer.TokenNot,
					Expr: ast.NewLiteralBool(false),
				},
			},
		},
		"bool-equal-bool": {
			str: `true==false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenEqual,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-not-equal-bool": {
			str: `true != false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenNotEqual,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-greater-than-bool": {
			str: `true>false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenGreaterThan,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-greater-than-equal-bool": {
			str: `true>=false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenGreaterThanEqual,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-less-than-bool": {
			str: `true < false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenLessThan,
				Right: ast.NewLiteralBool(false),
			},
		},
		"bool-less-than-equal-bool": {
			str: `true <= false`,
			expected: &ast.Binary{
				Left:  ast.NewLiteralBool(true),
				Op:    lexer.TokenLessThanEqual,
				Right: ast.NewLiteralBool(false),
			},
		},
		"increment": {
			str: "++a",
			expected: &ast.Unary{
				Op:   "++",
				Expr: &ast.Identifier{Name: "a"},
			},
		},
		"decrement": {
			str: "--a",
			expected: &ast.Unary{
				Op:   "--",
				Expr: &ast.Identifier{Name: "a"},
			},
		},
		"add-assign": {
			str: "a += 3",
			expected: &ast.Binary{
				Left:  &ast.Identifier{Name: "a"},
				Op:    lexer.TokenPlusAssign,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"minus-assign": {
			str: "a -= 3",
			expected: &ast.Binary{
				Left:  &ast.Identifier{Name: "a"},
				Op:    lexer.TokenMinusAssign,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"times-assign": {
			str: "a *= 3",
			expected: &ast.Binary{
				Left:  &ast.Identifier{Name: "a"},
				Op:    lexer.TokenTimesAssign,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"divide-assign": {
			str: "a /= 3",
			expected: &ast.Binary{
				Left:  &ast.Identifier{Name: "a"},
				Op:    lexer.TokenDivideAssign,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"remainder-assign": {
			str: "a %= 3",
			expected: &ast.Binary{
				Left:  &ast.Identifier{Name: "a"},
				Op:    lexer.TokenRemainderAssign,
				Right: ast.NewLiteralNumber("3"),
			},
		},
		"minus-1": {
			str: "1 - 2 + 3",
			expected: &ast.Binary{
				Left: ast.NewLiteralNumber("1"),
				Op:   lexer.TokenMinus,
				Right: &ast.Binary{
					Left:  ast.NewLiteralNumber("2"),
					Op:    lexer.TokenPlus,
					Right: ast.NewLiteralNumber("3"),
				},
			},
		},
		"minus-2": {
			str: "y - y + y",
			expected: &ast.Binary{
				Left: &ast.Identifier{Name: "y"},
				Op:   lexer.TokenMinus,
				Right: &ast.Binary{
					Left:  &ast.Identifier{Name: "y"},
					Op:    lexer.TokenPlus,
					Right: &ast.Identifier{Name: "y"},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			str := fmt.Sprintf("func main() { %s }", test.str)
			p := parser.ParseString(str)

			assertEqualErrors(t, test.errs, p.Errors)
			if test.expected == nil {
				assert.Equal(t, map[string]*ast.Func{}, p.File.Funcs)
			} else {
				assert.Equal(t, map[string]*ast.Func{
					"main": {
						Name:       "main",
						Statements: []ast.Node{test.expected},
					},
				}, p.File.Funcs)
			}
			assert.Nil(t, p.File.Comments)
		})
	}
}
