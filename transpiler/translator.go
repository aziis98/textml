package transpiler

import (
	"io"

	"github.com/aziis98/go-text-ml/parser"
)

type Transpiler interface {
	Transpile(w io.Writer, block *parser.Block) error
}

func BlockTextContent(b *parser.Block) string {
	s := ""

	for _, n := range b.Children {
		if n.Type == parser.NodeTypes.Text {
			s += n.Text
		}
	}

	return s
}
