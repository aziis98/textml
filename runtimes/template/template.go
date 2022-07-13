package template

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/parser"
)

func FileLoader(ctx *Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	doc, err := textml.ParseDocument(bufio.NewReader(f))

	if _, err := ctx.Evaluate(doc); err != nil {
		return err
	}

	return nil
}

type Config struct {
	LoaderFunc func(*Context, string) error
}

type Context struct {
	Config   *Config
	Registry map[string]*parser.Block
}

func New(config *Config) *Context {
	return &Context{
		Config:   config,
		Registry: map[string]*parser.Block{},
	}
}

func (te *Context) Evaluate(block *parser.Block) (string, error) {
	r := &strings.Builder{}

	for _, n := range block.Children {
		switch n := n.(type) {
		case *parser.ElementNode:
			switch n.Name {

			case "import":
				if len(n.Args) != 1 {
					return "", fmt.Errorf(`#import expected 1 argument, got %d`, len(n.Args))
				}

				if te.Config.LoaderFunc == nil {
					return "", fmt.Errorf(`template engine has no module loader`)
				}

				moduleName := n.Args[0].TextContent()

				if err := te.Config.LoaderFunc(te, moduleName); err != nil {
					return "", err
				}

			case "define":
				if len(n.Args) != 2 {
					return "", fmt.Errorf(`#define expected 2 arguments, got %d`, len(n.Args))
				}

				key := n.Args[0].TextContent()
				te.Registry[key] = n.Args[1]

			case "":
				if len(n.Args) != 1 {
					return "", fmt.Errorf(`anonymous block expected 1 argument, got %d`, len(n.Args))
				}

				key := n.Args[0].TextContent()
				value := te.Registry[key]

				s, err := te.Evaluate(value)
				if err != nil {
					return "", err
				}

				r.WriteString(s)
			default:
				return "", fmt.Errorf(`unexpected element "#%s"`, n.Name)
			}
		case *parser.TextNode:
			r.WriteString(n.Text)
		}
	}

	return r.String(), nil
}
