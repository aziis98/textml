package transpile

import (
	"io"

	"github.com/aziis98/textml/ast"
)

type StringTranspiler interface {
	Transpile(block ast.Block) (string, error)
}

type WriteTranspiler interface {
	Transpile(w io.Writer, block ast.Block) error
}

var Registry = map[string]any{
	// Go Repr
	"repr": &Repr{},
	// HTML
	"html":        &Html{Inline: false},
	"html.inline": &Html{Inline: true},
	// Json
	"json":        &Json{Inline: false},
	"json.inline": &Json{Inline: true},
}
