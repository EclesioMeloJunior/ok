package parser

import (
	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/lexer"
)

func consumeTest(parser *Parser, offset int) (*ast.Test, int, error) {
	originalOffset := offset
	var err error

	offset, err = consume(parser.File, offset, []string{
		lexer.TokenTest, lexer.TokenStringLiteral})
	if err != nil {
		return nil, originalOffset, err
	}

	t := &ast.Test{
		Name: parser.File.Tokens[offset-1].Value,
		Pos:  parser.File.Pos(originalOffset),
	}

	t.Statements, offset, err = consumeBlock(parser, offset)
	if err != nil {
		return nil, originalOffset, err
	}

	return t, offset, nil
}
