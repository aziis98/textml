package transpiler

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/aziis98/go-text-ml/parser"
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
}

func (h *Html) writeElement(w io.Writer, elem, s string) error {
	f := "<%s>\n%s\n</%s>\n"
	if h.Inline {
		f = "<%s>%s</%s>"
	}

	_, err := fmt.Fprintf(w, f, elem, html.EscapeString(strings.TrimSpace(s)), elem)

	if err != nil {
		return err
	}

	return nil
}

func (h *Html) writeElementOpenning(w io.Writer, elem string, attrs map[string]string) error {
	f := "<%s%s>\n"
	if h.Inline {
		f = "<%s%s>\n"
	}

	attrString := ""
	for k, v := range attrs {
		attrString += fmt.Sprintf(` %s="%s"`, k, v)
	}

	if _, err := fmt.Fprintf(w, f, elem, attrString); err != nil {
		return err
	}

	return nil
}

func (h *Html) writeElementClosing(w io.Writer, elem string) error {
	f := "</%s>\n"
	if h.Inline {
		f = "</%s>"
	}

	if _, err := fmt.Fprintf(w, f, elem); err != nil {
		return err
	}

	return nil
}

func (h *Html) TranspileElement(w io.Writer, element string, node *parser.Node) error {
	args := node.Args
	if len(args) > 2 {
		return fmt.Errorf(`invalid number of arguments for element "%q"`, element)
	}

	htmlAttributes := map[string]string{}

	if len(args) == 2 {
		var attrs *parser.Block
		attrs, args = args[0], args[1:]

		for _, n := range attrs.Children {
			if n.Type == parser.NodeTypes.Element {
				key := n.Name

				value := ""
				if len(n.Args) > 0 {
					value = BlockTextContent(n.Args[0])
				}

				htmlAttributes[key] = value
			}
		}
	}

	if err := h.writeElementOpenning(w, element, htmlAttributes); err != nil {
		return err
	}

	insideBlock := args[0]
	if err := h.TranspileBlock(w, insideBlock); err != nil {
		return err
	}

	if err := h.writeElementClosing(w, element); err != nil {
		return err
	}

	return nil
}

func (h *Html) TranspileBlock(w io.Writer, block *parser.Block) error {
	for _, node := range block.Children {
		if node.Type == parser.NodeTypes.Element {
			if htmlName, ok := htmlElements[node.Name]; ok {
				if err := h.TranspileElement(w, htmlName, node); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid html element with name %q", node.Name)
			}
		} else {
			if len(strings.TrimSpace(node.Text)) > 0 {
				if _, err := fmt.Fprintf(w, "%s\n",
					html.EscapeString(strings.TrimSpace(node.Text))); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (h *Html) Transpile(w io.Writer, block *parser.Block) error {
	if _, err := fmt.Fprintf(w, "<!DOCTYPE html>\n<html>\n"); err != nil {
		return err
	}

	if err := h.TranspileBlock(w, block); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "</html>\n"); err != nil {
		return err
	}

	return nil
}
