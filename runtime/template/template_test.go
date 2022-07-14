package template_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/runtime/template"
	"github.com/stretchr/testify/assert"
)

func dedent(s string) string {
	rg := regexp.MustCompile(`(?m)^\s+`)
	return strings.TrimSpace(rg.ReplaceAllString(s, ""))
}

func simplify(s string) string {
	rg := regexp.MustCompile(`[ ]*\n[ \n]*`)

	return strings.TrimSpace(
		rg.ReplaceAllString(
			strings.TrimSpace(s), "\n",
		),
	)
}

func removeNewlines(s string) string {
	rg := regexp.MustCompile(`\n`)

	return strings.TrimSpace(
		rg.ReplaceAllString(
			strings.TrimSpace(s), "",
		),
	)
}

func TestBasicEvaluation(t *testing.T) {
	ctx := template.New(&template.Config{})

	doc, err := textml.ParseDocument(strings.NewReader("Lorem ipsum"))
	assert.Nil(t, err)

	r, err := ctx.Evaluate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "Lorem ipsum", r)
}

func TestBasicDefine(t *testing.T) {
	ctx := template.New(&template.Config{})

	doc, err := textml.ParseDocument(strings.NewReader(dedent(`
		#define{ x }{ Lorem ipsum }
		#{ x }
	`)))
	assert.Nil(t, err)

	r, err := ctx.Evaluate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "\nLorem ipsum", r)
}

func TestLastWins(t *testing.T) {
	ctx := template.New(&template.Config{})

	doc, err := textml.ParseDocument(strings.NewReader(dedent(`
		#define{ x }{ 1 }
		#define{ x }{ 2 }
		#{ x }
	`)))
	assert.Nil(t, err)

	r, err := ctx.Evaluate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "2", strings.TrimSpace(r))
}

func TestExtends(t *testing.T) {
	ctx := template.New(&template.Config{})

	doc, err := textml.ParseDocument(strings.NewReader(dedent(`
		#define{ layout-1 }{
			Hello, #{ name }!
		}
		
		#extends{ layout-1 }
		#define{ name }{ John }
	`)))
	assert.Nil(t, err)

	r, err := ctx.Evaluate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, John!", strings.TrimSpace(r))
}

func TestNestedExtends(t *testing.T) {

	doc, err := textml.ParseDocument(strings.NewReader(dedent(`
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
	`)))
	assert.Nil(t, err)

	ctx := template.New(&template.Config{})
	r, err := ctx.Evaluate(doc)

	assert.Nil(t, err)
	assert.Equal(t, "Document(Article(Example))", removeNewlines(r))
}
