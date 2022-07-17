package document

import (
	"fmt"

	"github.com/aziis98/textml/html"
	"github.com/aziis98/textml/parser"
)

func parseDictEntries(ast parser.Block) (map[string]any, error) {
	m := map[string]any{}

	for _, n := range ast {
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

func parseDictValue(ast parser.Block) (any, error) {
	if elem := ast.FirstElement(); elem != nil {
		return parseDictEntries(ast)
	} else {
		return ast.TextContent(), nil
	}
}

// Engine is a Markdown like format that transpiles TextML to HTML
type Engine struct{}

type Metadata map[string]any

func checkArgCount(el *parser.ElementNode, count int) error {
	if len(el.Args) != count {
		return fmt.Errorf(`invalid argument count, expected %d but got %d`, count, len(el.Args))
	}

	return nil
}

var directTranslationMap = map[string]string{
	"title":          "h1",
	"subtitle":       "h2",
	"subsubtitle":    "h3",
	"subsubsubtitle": "h4",

	"bold":          "b",
	"italic":        "i",
	"underline":     "u",
	"strikethrough": "s",

	"code": "code",
}

func (t *Engine) RenderElement(el *parser.ElementNode) ([]html.Node, error) {
	// Direct translations
	if tagName, found := directTranslationMap[el.Name]; found {
		if err := checkArgCount(el, 1); err != nil {
			return nil, err
		}

		children, err := t.RenderBlock(el.Args[0])
		if err != nil {
			return nil, err
		}

		return []html.Node{
			html.NewElementNode(tagName, nil, children),
		}, nil
	}

	nodes := []html.Node{}

	switch el.Name {
	case "link":
		if err := checkArgCount(el, 2); err != nil {
			return nil, err
		}

		children, err := t.RenderBlock(el.Args[0])
		if err != nil {
			return nil, err
		}

		linkTarget := el.Args[1].TextContent()

		return []html.Node{
			html.NewElementNode(
				"a",
				html.AttributeMap{
					"href": &html.Attribute{Value: linkTarget},
				},
				children,
			),
		}, nil
	}

	return nodes, nil
}

func (t *Engine) RenderBlock(ast parser.Block) ([]html.Node, error) {
	// TODO: Add automatic paragraph splitting after "\n\n", this requires distinguishing between inline and block elements...
	nodes := []html.Node{}

	for _, n := range ast {
		switch n := n.(type) {
		case *parser.ElementNode:
			children, err := t.RenderElement(n)
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, children...)
		case *parser.TextNode:
			nodes = append(nodes, &html.Text{Value: n.Text})
		default:
			panic("illegal state")
		}
	}

	return nodes, nil
}

func (t *Engine) Render(ast parser.Block) (Metadata, []html.Node, error) {
	documentMetadata := Metadata{}
	nodes := []html.Node{}

	for _, n := range ast {
		switch n := n.(type) {
		case *parser.ElementNode:
			switch n.Name {
			case "metadata":
				metadata, err := parseDictEntries(n.Args[0])
				if err != nil {
					return nil, nil, err
				}

				for k, v := range metadata {
					documentMetadata[k] = v
				}
			default:
				children, err := t.RenderElement(n)
				if err != nil {
					return nil, nil, err
				}

				nodes = append(nodes, children...)
			}
		case *parser.TextNode:
			nodes = append(nodes, &html.Text{Value: n.Text})
		default:
			panic("illegal state")
		}
	}

	return documentMetadata, nodes, nil
}
