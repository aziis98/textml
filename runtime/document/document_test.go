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
	return regexp.MustCompile(`[ ]*\n[ \n]*`).ReplaceAllString(s, "\n")
}

const doc1 = `
#metadata {
	#uuid { 8bXY0njav9I3bisQfe8K7uzp7WJAjl }
	#title { Example }
	#tags { example, tag-1, tag-2, other }
	#a-dict-value {
		#a { 1 }
		#b { 2 }
	}
}

#title { Short title }

A simple paragraph.

#subtitle { Another example }

Another paragraph #bold{ with } #italic{ some } #underline{ formatting } and
a link to #link{ Wikipedia }{ https://en.wikipedia.org }.
`

var rendered1 string = `
<h1>Short title</h1>

A simple paragraph.

<h2>Another example</h2>

Another paragraph <b>with</b> <i>some</i> <u>formatting</u> and
a link to <a href="https://en.wikipedia.org">Wikipedia</a>.
`

func Test1(t *testing.T) {
	doc, err := textml.ParseDocument(strings.NewReader(doc1))
	assert.Nil(t, err)

	engine := &document.Engine{}

	metadata, nodes, err := engine.Render(doc)
	assert.Nil(t, err)
	assert.Equal(t,
		document.Metadata{
			"uuid":  "8bXY0njav9I3bisQfe8K7uzp7WJAjl",
			"title": "Example",
			"tags":  "example, tag-1, tag-2, other",
			"a-dict-value": map[string]any{
				"a": "1",
				"b": "2",
			},
		},
		metadata,
	)

	htmlString := html.RenderToString(nodes)
	assert.Equal(t,
		simplifyLines(rendered1),
		simplifyLines(htmlString),
	)
}
