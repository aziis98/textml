package template_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/runtimes/template"
	"github.com/stretchr/testify/assert"
)

func dedent(s string) string {
	rg := regexp.MustCompile(`(?m)^\s+`)
	return strings.TrimSpace(rg.ReplaceAllString(s, ""))
}

func Test1(t *testing.T) {
	ctx := template.New(&template.Config{})

	doc, err := textml.ParseDocument(strings.NewReader("Lorem ipsum"))
	assert.Nil(t, err)

	r, err := ctx.Evaluate(doc)
	assert.Nil(t, err)
	assert.Equal(t, "Lorem ipsum", r)
}

func Test2(t *testing.T) {
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
