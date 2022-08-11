package template

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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
	TrimSpaces bool
	LoaderFunc func(*Engine, string) error
}

type Engine struct {
	Config    Config
	Variables map[string]any
	Templates map[string]ast.Block
}

func New(defaultConfig ...Config) *Engine {
	config := Config{
		TrimSpaces: false,
		LoaderFunc: nil,
	}
	if len(defaultConfig) > 0 {
		config = defaultConfig[0]
	}

	return &Engine{
		Config:    config,
		Variables: map[string]any{},
		Templates: map[string]ast.Block{},
	}
}

func (e *Engine) evaluateValue(block ast.Block) (any, error) {
	if len(block) != 1 {
		return nil, fmt.Errorf("invalid value node")
	}

	node, ok := block[0].(*ast.TextNode)
	if !ok {
		return nil, fmt.Errorf("invalid value node")
	}

	text := strings.TrimSpace(node.Text)

	// if first rune is a digit then try parse to float64 or int64
	r, _ := utf8.DecodeRuneInString(text)
	if unicode.IsDigit(r) {
		floatValue, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return float64(floatValue), nil
		}

		intValue, err := strconv.ParseInt(text, 10, 64)
		if err == nil {
			return int(intValue), nil
		}
	}

	// otherwise check if boolean
	if text == "true" {
		return true, nil
	}
	if text == "false" {
		return false, nil
	}

	// otherwise return variable
	value, ok := e.Variables[text]
	if !ok {
		return nil, fmt.Errorf("unknown variable %q", text)
	}

	return value, nil
}

func errInvalidElement(elem *ast.ElementNode) error {
	return fmt.Errorf("invalid template command %q with %d arguments", elem.Name, len(elem.Arguments))
}

var commandCharMap = map[string]string{
	"space":   " ",
	"tab":     "\t",
	"newline": "\n",
}

var regexLineWithIndent = regexp.MustCompile(` *\n\s*`)

func inlineBlock(block ast.Block) ast.Block {
	result := ast.Block{}
	for _, node := range block {
		switch node := node.(type) {
		case *ast.ElementNode:
			arguments := []ast.Block{}

			for _, arg := range node.Arguments {
				arguments = append(arguments, inlineBlock(arg))
			}

			result = append(result, &ast.ElementNode{
				Name:      node.Name,
				Arguments: arguments,
			})
		case *ast.TextNode:
			result = append(result, &ast.TextNode{
				Text: regexLineWithIndent.ReplaceAllString(node.Text, ""),
			})
		}
	}
	return result
}

func (e *Engine) evaluateElement(elem *ast.ElementNode) (any, error) {
	switch elem.Name {
	case "import":
		if len(elem.Arguments) != 1 {
			return nil, errInvalidElement(elem)
		}
		if e.Config.LoaderFunc == nil {
			return "", fmt.Errorf(`template engine has no module loader`)
		}
		moduleName := elem.Arguments[0].TextContent()
		if err := e.Config.LoaderFunc(e, moduleName); err != nil {
			return "", err
		}

		return nil, nil

	case "template":
		if len(elem.Arguments) != 2 {
			return "", errInvalidElement(elem)
		}
		tmplName := elem.Arguments[0].TextContent()
		e.Templates[tmplName] = elem.Arguments[1]

		return nil, nil

	case "define":
		if len(elem.Arguments) != 2 {
			return "", errInvalidElement(elem)
		}
		varName := elem.Arguments[0].TextContent()

		varValue, err := e.evaluateBlock(elem.Arguments[1])
		if err != nil {
			return nil, err
		}

		e.Variables[varName] = varValue

		return nil, nil

	case "extends":
		if len(elem.Arguments) != 2 {
			return "", fmt.Errorf(`#define expected 2 arguments, got %d`, len(elem.Arguments))
		}

		key := elem.Arguments[0].TextContent()
		extendingTemplate, ok := e.Templates[key]
		if !ok {
			return nil, fmt.Errorf("no binding for %q", key)
		}

		_, err := e.evaluateBlock(elem.Arguments[1])
		if err != nil {
			return nil, err
		}

		result, err := e.evaluateBlock(extendingTemplate)
		if err != nil {
			return nil, err
		}

		return result, nil

	case "if", "unless":
		if len(elem.Arguments) < 2 && len(elem.Arguments) > 3 {
			return nil, errInvalidElement(elem)
		}

		invertCond := elem.Name == "unless"

		cond := elem.Arguments[0]
		ifTrue := elem.Arguments[1]

		condValue, err := e.evaluateBlock(cond)
		if err != nil {
			return nil, err
		}

		c, ok := condValue.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid #if condition type %T", condValue)
		}

		if invertCond {
			c = !c
		}

		if c {
			return e.evaluateBlock(ifTrue)
		} else if len(elem.Arguments) == 3 {
			return e.evaluateBlock(elem.Arguments[2])
		} else {
			return nil, nil
		}

	case "foreach":
		if len(elem.Arguments) != 3 {
			return nil, errInvalidElement(elem)
		}

		loopVarName := strings.TrimSpace(elem.Arguments[0].TextContent())
		itemsVarName := strings.TrimSpace(elem.Arguments[1].TextContent())

		itemsValue, ok := e.Variables[itemsVarName]
		if !ok {
			return nil, fmt.Errorf("no binding for %q", itemsVarName)
		}

		if reflect.TypeOf(itemsValue).Kind() != reflect.Slice {
			return nil, fmt.Errorf("the type %T given to #foreach cannot be iterated", itemsValue)
		}

		sb := &strings.Builder{}

		s := reflect.ValueOf(itemsValue)
		for i := 0; i < s.Len(); i++ {
			itemValue := s.Index(i).Interface()
			e.Variables[loopVarName] = itemValue

			result, err := e.evaluateBlock(elem.Arguments[2])
			if err != nil {
				return nil, err
			}

			fmt.Fprintf(sb, "%v", result)
		}

		return sb.String(), nil

	case "intersperse":
		if len(elem.Arguments) != 2 {
			return nil, errInvalidElement(elem)
		}

		itemsVarName := strings.TrimSpace(elem.Arguments[0].TextContent())
		intersperseStr := elem.Arguments[1].TextContent()

		itemsValue, ok := e.Variables[itemsVarName]
		if !ok {
			return nil, fmt.Errorf("no binding for %q", itemsVarName)
		}

		if reflect.TypeOf(itemsValue).Kind() != reflect.Slice {
			return nil, fmt.Errorf("the type %T given to #foreach cannot be iterated", itemsValue)
		}

		sb := &strings.Builder{}

		s := reflect.ValueOf(itemsValue)
		for i := 0; i < s.Len(); i++ {
			if i != 0 {
				sb.WriteString(intersperseStr)
			}

			itemValue := s.Index(i).Interface()
			fmt.Fprintf(sb, "%v", itemValue)
		}

		return sb.String(), nil

	case "char":
		if len(elem.Arguments) != 1 {
			return nil, errInvalidElement(elem)
		}

		charName := elem.Arguments[0].TextContent()
		str, ok := commandCharMap[charName]
		if !ok {
			return nil, fmt.Errorf("invalid char %q", charName)
		}

		return str, nil

	case "inline":
		if len(elem.Arguments) != 1 {
			return nil, errInvalidElement(elem)
		}

		inlinedAst := inlineBlock(elem.Arguments[0])

		return e.evaluateBlock(inlinedAst)

	case "":
		if len(elem.Arguments) != 1 {
			return nil, errInvalidElement(elem)
		}

		return e.evaluateBlock(elem.Arguments[0])

	default:
		return nil, fmt.Errorf("invalid template command %q", elem.Name)
	}
}

func (e *Engine) evaluateBlock(block ast.Block) (any, error) {

	// first try evaluate value literals and variables
	value, err := e.evaluateValue(block)
	if err == nil {
		return value, nil
	}

	sb := &strings.Builder{}
	var resultValue any = nil

	elemNodeCount := 0
	nonBlankTextNodeCount := 0

	for _, node := range block {
		switch node := node.(type) {
		case *ast.ElementNode:
			elemNodeCount++

			result, err := e.evaluateElement(node)
			if err != nil {
				return nil, err
			}

			resultValue = result

			if result != nil {
				fmt.Fprintf(sb, "%v", result)
			}
		case *ast.TextNode:
			if len(strings.TrimSpace(node.Text)) > 0 {
				nonBlankTextNodeCount++
			}

			text := node.Text
			if e.Config.TrimSpaces {
				text = strings.Trim(text, " ")
			}

			fmt.Fprintf(sb, "%s", text)
		}
	}

	if elemNodeCount == 1 && nonBlankTextNodeCount == 0 {
		return resultValue, nil
	}

	return sb.String(), nil
}

type EngineEvaluateOption func(e *Engine)

func WithTemplate(name string, tmpl ast.Block) EngineEvaluateOption {
	return func(e *Engine) {
		e.Templates[name] = tmpl
	}
}

func WithVariable(key string, value any) EngineEvaluateOption {
	return func(e *Engine) {
		e.Variables[key] = value
	}
}

func WithContext(vars map[string]any) EngineEvaluateOption {
	return func(e *Engine) {
		for k, v := range vars {
			e.Variables[k] = v
		}
	}
}

func (e *Engine) Evaluate(block ast.Block, options ...EngineEvaluateOption) (string, error) {
	for _, opt := range options {
		opt(e)
	}

	result, err := e.evaluateBlock(block)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}

	return fmt.Sprintf("%v", result), nil
}
