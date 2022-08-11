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
			dedent(s), "\n",
		),
	)
}

func removeNewlines(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
}

func renderTemplate(e *template.Engine, s string, opts ...template.EngineEvaluateOption) (string, error) {
	doc, err := textml.ParseDocument(strings.NewReader(dedent(s)))
	if err != nil {
		return "", err
	}

	rendered, err := e.Evaluate(doc, opts...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(rendered), nil
}

func TestBasicEvaluation(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		Lorem ipsum
	`)
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
		#template{ example }{
			Hello, #{ name }!
		}
		
		#extends{ example }{
			#define{ name }{ John }
		}
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, John!", a)
}

func TestNestedExtends(t *testing.T) {
	a, err := renderTemplate(template.New(), `
		#template{ base-layout }{
			Document(#{ body })
		}

		#template{ article-layout }{
			#extends{ base-layout }{
				#define{ body }{
					Article(#{ article.body })
				}
			}
		}
		
		#extends{ article-layout }{
			#define{ article.body }{ Example }
		}
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Document(Article(Example))", removeNewlines(a))
}

func TestExpressionIf2(t *testing.T) {
	eng := template.New()
	eng.Variables["c"] = false
	eng.Variables["x"] = "Nope"

	a, err := renderTemplate(eng, `
		Hidden: #if{ c }{ #{ x } }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "Hidden:", a)
}

func TestExpressionIf3(t *testing.T) {
	eng := template.New()
	eng.Variables["c"] = false
	eng.Variables["x"] = 123
	eng.Variables["y"] = 234

	a, err := renderTemplate(eng, `
		#if{ c }{ #{ x } }{ #{ y } }
	`)
	assert.Nil(t, err)
	assert.Equal(t, "234", a)
}

func TestExpressionForEach(t *testing.T) {
	eng := template.New()
	eng.Variables["names"] = []string{"Adam", "Billy", "John", "Rose"}

	a, err := renderTemplate(eng, `
		#foreach{ name }{ names }{
		Hi, #{ name }! }
	`)
	assert.Nil(t, err)
	assert.Equal(t, simplify(`
		Hi, Adam!
		Hi, Billy!
		Hi, John!
		Hi, Rose!
	`), a)
}

func TestExpressionForEachChar(t *testing.T) {
	eng := template.New()
	eng.Variables["names"] = []string{"Adam", "Billy", "John", "Rose"}

	a, err := renderTemplate(eng, `
		#foreach{ name }{ names }{ Hi, #{ name }!#char{ newline } }
	`)
	assert.Nil(t, err)
	assert.Equal(t, simplify(`
		Hi, Adam!
		Hi, Billy!
		Hi, John!
		Hi, Rose!
	`), a)
}

func TestExpressionIntersperse(t *testing.T) {
	eng := template.New()
	eng.Variables["things"] = []string{"aaa", "bbb", "ccc", "ddd"}

	a, err := renderTemplate(eng, `
		#intersperse{ things }{ ,  }
	`)
	assert.Nil(t, err)
	assert.Equal(t, `aaa, bbb, ccc, ddd`, a)
}
