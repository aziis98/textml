package transpile

import (
	"encoding/json"
	"io"
	"log"

	"github.com/alecthomas/repr"
	"github.com/aziis98/textml/ast"
)

type Repr struct{}

func (Repr) Transpile(w io.Writer, block ast.Block) error {
	repr.New(w).Println(block)
	return nil
}

type Json struct{ Inline bool }

func (t *Json) Transpile(w io.Writer, block ast.Block) error {
	enc := json.NewEncoder(w)

	if !t.Inline {
		enc.SetIndent("", "    ")
	}

	if err := enc.Encode(block); err != nil {
		log.Fatal(err)
	}

	return nil
}
