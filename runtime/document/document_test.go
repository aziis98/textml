package document_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/html"
	"github.com/aziis98/textml/runtime/document"
	"github.com/stretchr/testify/assert"
)

func simplifyLines(s string) string {
	return strings.TrimSpace(
		regexp.MustCompile(`[ ]*\n[ \n]*`).ReplaceAllString(s, "\n"),
	)
}

const doc1 = `
#metadata {
	#uuid { 8bXY0njav9I3bisQfe8K7uzp7WJAjl }
	#title { Example }
	#tags { example, tag-1, tag-2, other }
	#a-dict-value {
		#dict{
			#a { 1 }
			#b { 2 }
		}
	}
}

#title { Short title }

A simple paragraph.

#subtitle { Another example }

Another Paragraph.
`

func Test1(t *testing.T) {
	doc, err := textml.ParseDocument(strings.NewReader(doc1))
	assert.Nil(t, err)

	engine := &document.Engine{}

	metadata, nodes, err := engine.Render(doc)
	assert.Equal(t, document.Metadata{
		"uuid":  "8bXY0njav9I3bisQfe8K7uzp7WJAjl",
		"title": "Example",
		"tags":  "example, tag-1, tag-2, other",
		"a-dict-value": map[string]any{
			"a": "1",
			"b": "2",
		},
	}, metadata)
	assert.Nil(t, err)

	htmlString := html.RenderToString(nodes)
	assert.Equal(t, "<h1>Short title</h1>\nA simple paragraph.\n<h2>Another example</h2>\nAnother Paragraph.", simplifyLines(htmlString))
}
