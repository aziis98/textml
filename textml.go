package textml

import (
	"io"

	"github.com/aziis98/textml/lexer"
	"github.com/aziis98/textml/parser"
)

func ParseDocument(r io.RuneReader) (parser.Block, error) {
	tokens, err := lexer.New(r).AllTokens()
	if err != nil {
		return nil, err
	}

	return parser.ParseDocument(tokens)
}
