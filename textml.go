package textml

import (
	"io"

	"github.com/aziis98/textml/ast"
	"github.com/aziis98/textml/lexer"
	"github.com/aziis98/textml/parser"
)

// ParseDocument tokenizes the input using [lexer] and then parses it with [parser.ParseDocument]
func ParseDocument(r io.RuneReader) (ast.Block, error) {
	tokens, err := lexer.New(r).AllTokens()
	if err != nil {
		return nil, err
	}

	doc, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return ast.Compile(doc), nil
}
