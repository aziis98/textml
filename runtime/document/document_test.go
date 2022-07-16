package document_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/runtime/document"
	"github.com/stretchr/testify/assert"
)

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
`

func Test1(t *testing.T) {
	b := &bytes.Buffer{}
	w := bufio.NewWriter(b)

	doc, err := textml.ParseDocument(strings.NewReader(doc1))
	assert.Nil(t, err)

	engine := &document.Engine{}

	metadata, err := engine.Render(doc, w)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"uuid":  "8bXY0njav9I3bisQfe8K7uzp7WJAjl",
		"title": "Example",
		"tags":  "example, tag-1, tag-2, other",
		"a-dict-value": map[string]any{
			"a": "1",
			"b": "2",
		},
	}, metadata)
}
