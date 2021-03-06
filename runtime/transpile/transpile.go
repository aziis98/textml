package transpile

import (
	"io"

	"github.com/aziis98/textml/parser"
)

type Transpiler interface {
	Transpile(block parser.Block, w io.Writer) error
}

var Registry = map[string]Transpiler{
	// Go Repr
	"repr": &Repr{},
	// HTML
	"html":        &Html{Inline: false},
	"html.inline": &Html{Inline: true},
	// Json
	"json":        &Json{Inline: false},
	"json.inline": &Json{Inline: true},
}
