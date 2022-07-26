package textml

import (
	"io"

	"github.com/aziis98/textml/ast"
	"github.com/aziis98/textml/lexer"
	"github.com/aziis98/textml/parser"
)

func ParseDocument(r io.RuneReader) (ast.Block, error) {
	tokens, err := lexer.New(r).AllTokens()
	if err != nil {
		return nil, err
	}

	doc, err := parser.ParseDocument(tokens)
	if err != nil {
		return nil, err
	}

	return ast.Compile(doc), nil
}
