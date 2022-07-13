package transpiler

import (
	"io"

	"github.com/aziis98/textml/parser"
)

type Transpiler interface {
	Transpile(w io.Writer, block *parser.Block) error
}
