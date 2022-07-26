package transpile

import (
	"fmt"
	"html"
	"strings"

	"github.com/aziis98/textml/ast"
)

type Html struct {
	Inline bool
}

var htmlElements = map[string]string{
	"html.head":    "head",
	"html.title":   "title",
	"html.body":    "body",
	"html.main":    "main",
	"html.section": "section",
	"html.p":       "p",
	"html.h1":      "h1",
	"html.h2":      "h2",
	"html.h3":      "h3",
	"html.h4":      "h4",
	"html.h5":      "h5",
	"html.h6":      "h6",
	"html.ul":      "ul",
	"html.ol":      "ol",
	"html.li":      "li",
	"html.strong":  "strong",
	"html.em":      "em",
	"html.b":       "b",
	"html.i":       "i",
	"html.u":       "u",
	"html.code":    "code",
	"html.pre":     "pre",
	"html.div":     "div",
	"html.span":    "span",
	"html.img":     "img",
	"html.figure":  "figure",
}

func (h *Html) writeElement(elem, s string) string {
	args := []any{elem, html.EscapeString(strings.TrimSpace(s)), elem}

	if h.Inline {
		return fmt.Sprintf("<%s>%s</%s>", args...)
	}
	return fmt.Sprintf("<%s>\n%s\n</%s>\n", args...)
}

func (h *Html) writeElementOpening(elem string, attrs map[string]string) string {
	f := "<%s%s>\n"
	if h.Inline {
		f = "<%s%s>"
	}

	attrString := ""
	for k, v := range attrs {
		attrString += fmt.Sprintf(` %s="%s"`, k, v)
	}

	return fmt.Sprintf(f, elem, attrString)
}

func (h *Html) writeElementClosing(elem string) string {
	f := "</%s>\n"
	if h.Inline {
		f = "</%s>"
	}

	return fmt.Sprintf(f, elem)
}

func (h *Html) TranspileElement(node *ast.ElementNode) (string, error) {
	element, ok := htmlElements[node.Name]
	if !ok {
		return "", fmt.Errorf("invalid html element with name %q", node.Name)
	}

	args := node.Arguments
	if len(args) > 2 {
		return "", fmt.Errorf(`invalid number of arguments for element "%q"`, element)
	}

	htmlAttributes := map[string]string{}

	if len(args) == 2 {
		var attrs ast.Block
		attrs, args = args[0], args[1:]

		for _, n := range attrs {
			if elm, ok := n.(*ast.ElementNode); ok {
				key := elm.Name

				value := ""
				if len(elm.Arguments) > 0 {
					value = elm.Arguments[0].TextContent()
				}

				htmlAttributes[key] = value
			}
		}
	}

	htmlBlock, err := h.TranspileBlock(args[0])
	if err != nil {
		return "", err
	}

	sb := &strings.Builder{}

	fmt.Fprintf(sb, h.writeElementOpening(element, htmlAttributes))
	fmt.Fprintf(sb, htmlBlock)
	if !h.Inline {
		fmt.Fprintln(sb)
	}
	fmt.Fprintf(sb, h.writeElementClosing(element))

	return sb.String(), nil
}

func (h *Html) TranspileBlock(b ast.Block) (string, error) {
	sb := &strings.Builder{}

	for _, node := range b {
		switch node := node.(type) {
		case *ast.ElementNode:
			htmlElem, err := h.TranspileElement(node)
			if err != nil {
				return "", err
			}

			fmt.Fprintln(sb, htmlElem)

		case *ast.TextNode:
			if len(strings.TrimSpace(node.Text)) > 0 {
				fmt.Fprintln(sb, html.EscapeString(strings.TrimSpace(node.Text)))
			}

		default:
			panic("invalid node type")
		}
	}

	if !h.Inline {
		fmt.Fprintln(sb)
	}

	return sb.String(), nil
}

func (h *Html) Transpile(b ast.Block) (string, error) {

	htmlBlock, err := h.TranspileBlock(b)
	if err != nil {
		return "", err
	}

	sb := &strings.Builder{}

	fmt.Fprintln(sb, "<!DOCTYPE html>")
	fmt.Fprintln(sb, "<html>")
	fmt.Fprintln(sb, htmlBlock)
	fmt.Fprintln(sb, "</html>")

	return sb.String(), nil
}
