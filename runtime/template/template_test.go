package template_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/runtime/template"
	"github.com/stretchr/testify/assert"
)

var regexLineIndent = regexp.MustCompile(`(?m)^\s+`)
var regexConsecutiveBlankLines = regexp.MustCompile(`[ ]*\n[ \n]*`)

func dedent(s string) string {
	return strings.TrimSpace(regexLineIndent.ReplaceAllString(s, ""))
}

func simplify(s string) string {
	return strings.TrimSpace(
		regexConsecutiveBlankLines.ReplaceAllString(
			strings.TrimSpace(s), "\n",
		),
	)
}

func removeNewlines(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
}

func renderTemplate(e *template.Engine, s string) (string, error) {
	doc, err := textml.ParseDocument(strings.NewReader(dedent(s)))
	if err != nil {
		return "", err
	}

	rendered, err := e.Evaluate(doc)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(rendered), nil
}

func TestBasicEvaluation(t *testing.T) {
	a, err := renderTemplate(template.New(), "Lorem ipsum")
	assert.Nil(t, err)
	assert.Equal(t, "Lorem ipsum", a)
}

func TestBasicDefine(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		#define{ x }{ Lorem ipsum }
		#{ x }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Lorem ipsum", a)
}

func TestLastWins(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		#define{ x }{ 1 }
		#define{ x }{ 2 }
		#{ x }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "2", a)
}

func TestExtends(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		#define{ layout-1 }{
			Hello, #{ name }!
		}
		
		#extends{ layout-1 }
		#define{ name }{ John }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, John!", a)
}

func TestNestedExtends(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		#define{ base-layout }{
			Document(#{ body })
		}

		#define{ article-layout }{
			#extends{ base-layout }

			#define{ body }{
				Article(#{ article.body })
			}
		}
		
		#extends{ article-layout }
		#define{ article.body }{ Example }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Document(Article(Example))", removeNewlines(a))
}

func TestExpressionMultipleVariables(t *testing.T) {
	eng := template.New()
	eng.Registry["x"] = 123
	eng.Registry["y"] = 234

	a, err := renderTemplate(eng, `#{ x y }`)
	assert.Nil(t, err)
	assert.Equal(t, "123234", a)
}

func TestExpressionIfs(t *testing.T) {
	eng := template.New()
	eng.Registry["c"] = false
	eng.Registry["x"] = 123
	eng.Registry["y"] = 234

	a, err := renderTemplate(eng, `#{ #if{ c }{ x }{ y } }`)
	assert.Nil(t, err)
	assert.Equal(t, "234", a)
}
