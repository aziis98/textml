package template

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/ast"
)

func FileLoader(ctx *Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	doc, err := textml.ParseDocument(bufio.NewReader(f))
	if err != nil {
		return err
	}

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
	Registry map[string]ast.Block
}

func New(config *Config) *Context {
	return &Context{
		Config:   config,
		Registry: map[string]ast.Block{},
	}
}

func (te *Context) Evaluate(block ast.Block) (string, error) {
	r := &strings.Builder{}
	var extendsDirective *string = nil

	for _, n := range block {
		switch n := n.(type) {
		case *ast.ElementNode:
			switch n.Name {

			case "import":
				if len(n.Arguments) != 1 {
					return "", fmt.Errorf(`#import expected 1 argument, got %d`, len(n.Arguments))
				}

				if te.Config.LoaderFunc == nil {
					return "", fmt.Errorf(`template engine has no module loader`)
				}

				moduleName := n.Arguments[0].TextContent()

				if err := te.Config.LoaderFunc(te, moduleName); err != nil {
					return "", err
				}

			case "extends":
				if len(n.Arguments) != 1 {
					return "", fmt.Errorf(`#import expected 1 argument, got %d`, len(n.Arguments))
				}

				extendsDirective = new(string)
				*extendsDirective = n.Arguments[0].TextContent()
			case "define":
				if len(n.Arguments) != 2 {
					return "", fmt.Errorf(`#define expected 2 arguments, got %d`, len(n.Arguments))
				}

				key := n.Arguments[0].TextContent()
				te.Registry[key] = n.Arguments[1]

			case "":
				if len(n.Arguments) != 1 {
					return "", fmt.Errorf(`anonymous block expected 1 argument, got %d`, len(n.Arguments))
				}

				key := n.Arguments[0].TextContent()
				value := te.Registry[key]

				s, err := te.Evaluate(value)
				if err != nil {
					return "", err
				}

				r.WriteString(s)
			default:
				return "", fmt.Errorf(`unexpected element "#%s"`, n.Name)
			}
		case *ast.TextNode:
			r.WriteString(n.Text)
		}
	}

	if extendsDirective != nil {
		value := te.Registry[*extendsDirective]

		s, err := te.Evaluate(value)
		if err != nil {
			return "", err
		}

		r.WriteString(s)
	}

	return r.String(), nil
}
