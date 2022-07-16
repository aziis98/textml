package document

import (
	"fmt"
	"io"

	"github.com/aziis98/textml/parser"
)

func parseDictEntries(ast *parser.Block) (map[string]any, error) {
	m := map[string]any{}

	for _, n := range ast.Children {
		if n, ok := n.(*parser.ElementNode); ok {
			val, err := parseDictValue(n.Args[0])
			if err != nil {
				return nil, err
			}

			m[n.Name] = val
		}
	}

	return m, nil
}

func parseDictValue(ast *parser.Block) (any, error) {
	if elem := ast.FirstElement(); elem != nil {
		switch elem.Name {
		case "dict":
			return parseDictEntries(elem.Args[0])
		default:
			return nil, fmt.Errorf(`invalid dict value with identifier %q`, elem.Name)
		}
	} else {
		return ast.TextContent(), nil
	}
}

// Engine is a Markdown like format that transpiles TextML to HTML
type Engine struct{}

func (t *Engine) Render(ast *parser.Block, w io.Writer) (map[string]any, error) {
	documentMetadata := map[string]any{}

	for _, n := range ast.Children {
		switch n := n.(type) {
		case *parser.ElementNode:
			switch n.Name {
			case "metadata":
				metadata, err := parseDictEntries(n.Args[0])
				if err != nil {
					return nil, err
				}

				for k, v := range metadata {
					documentMetadata[k] = v
				}
			default:
				return nil, fmt.Errorf(`unexpected identifier %q`, n.Name)
			}
		case *parser.TextNode:
			if _, err := fmt.Fprintf(w, "%s", n.Text); err != nil {
				return nil, err
			}
		default:
			panic("illegal state")
		}
	}

	return documentMetadata, nil
}
