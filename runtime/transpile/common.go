package transpile

import (
	"encoding/json"
	"io"
	"log"

	"github.com/alecthomas/repr"
	"github.com/aziis98/textml/parser"
)

type Repr struct{}

func (_ Repr) Transpile(w io.Writer, ast *parser.Block) error {
	repr.New(w).Println(ast)

	return nil
}

type Json struct{ Inline bool }

func (t *Json) Transpile(w io.Writer, ast *parser.Block) error {
	enc := json.NewEncoder(w)

	if !t.Inline {
		enc.SetIndent("", "    ")
	}

	if err := enc.Encode(ast); err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}

	return nil
}
