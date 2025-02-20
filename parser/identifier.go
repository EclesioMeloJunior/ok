package parser

import (
	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/lexer"
)

func consumeIdentifier(parser *Parser, offset int) (*ast.Identifier, int, error) {
	originalOffset := offset

	var err error
	offset, err = consume(parser.File, offset, []string{lexer.TokenIdentifier})
	if err != nil {
		return nil, originalOffset, err
	}

	return &ast.Identifier{
		Name: parser.File.Tokens[originalOffset].Value,
		Pos:  parser.File.Pos(originalOffset),
	}, offset, nil
}
