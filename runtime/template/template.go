package template

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/ast"
)

func FileLoader(ctx *Engine, filename string) error {
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
	LoaderFunc func(*Engine, string) error
}

type Engine struct {
	Config   Config
	Registry map[string]any
}

func New(defaultConfig ...Config) *Engine {
	config := Config{
		LoaderFunc: nil,
	}
	if len(defaultConfig) > 0 {
		config = defaultConfig[0]
	}

	return &Engine{
		Config:   config,
		Registry: map[string]any{},
	}
}

func (te *Engine) ActuallyEvaluateExpression(expr ast.Block) (any, error) {
	sb := &strings.Builder{}
	var result any

	for _, node := range expr {
		switch node := node.(type) {
		case *ast.ElementNode:
			switch node.Name {
			case "if":
				if len(node.Arguments) != 3 {
					return nil, fmt.Errorf("if expected 3 args but got %d", len(node.Arguments))
				}

			default:
				return nil, fmt.Errorf("unexpected expression %q", node.Name)
			}
		case *ast.TextNode:
			varParts := regexp.MustCompile(`\s+`).Split(node.Text, -1)
			for _, varPart := range varParts {
				varName := strings.TrimSpace(varPart)
				varValue, ok := te.Registry[varName]
				if !ok {
					return nil, fmt.Errorf("no variable named %q", varName)
				}

				switch value := varValue.(type) {
				case ast.Block:
					s, err := te.Evaluate(value)
					if err != nil {
						return nil, err
					}

					sb.WriteString(s)
				default:
					if _, err := fmt.Fprintf(sb, "%v", value); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if result == nil {
		return sb.String(), nil
	}

	return result, nil
}

func (te *Engine) EvaluateExpression(expr ast.Block) (string, error) {
	v, err := te.ActuallyEvaluateExpression(expr)
	if err != nil {
		return "", err
	}

	if s, ok := v.(string); ok {
		return s, nil
	}

	return fmt.Sprintf("%v", v), nil
}

func (te *Engine) Evaluate(block ast.Block) (string, error) {
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

				s, err := te.EvaluateExpression(n.Arguments[0])
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

		s, err := te.Evaluate(value.(ast.Block))
		if err != nil {
			return "", err
		}

		r.WriteString(s)
	}

	return r.String(), nil
}
