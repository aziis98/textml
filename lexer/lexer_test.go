package lexer_test

import (
	"strings"
	"testing"

	"github.com/aziis98/textml/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLexer1(t *testing.T) {
	s := strings.NewReader("Lorem #node{ipsum} dolor")
	tokens, err := lexer.New(s).AllTokens()

	assert.Nil(t, err)
	assert.Equal(t, []lexer.Token{
		{Type: lexer.TextToken, Value: "Lorem "},
		{Type: lexer.ElementToken, Value: "#node"},
		{Type: lexer.BraceOpenToken, Value: "{"},
		{Type: lexer.TextToken, Value: "ipsum"},
		{Type: lexer.BraceCloseToken, Value: "}"},
		{Type: lexer.TextToken, Value: " dolor"},
		{Type: lexer.EOFToken},
	}, tokens)
}
func TestLexer2(t *testing.T) {
	s := strings.NewReader("Lorem #node{{ipsum}}} dolor")
	tokens, err := lexer.New(s).AllTokens()

	assert.Nil(t, tokens)
	assert.Equal(t, "errors: [too many braces at 21]", err.Error())
}

var example3 = strings.TrimSpace(`
#document {
    #title { A short title }

    This is some text with some #bold{ bold } text
}
`)

func TestLexer3(t *testing.T) {
	s := strings.NewReader(example3)

	tokens, err := lexer.New(s).AllTokens()

	assert.Equal(t, []lexer.Token{
		{lexer.ElementToken, "#document"},
		{lexer.BraceOpenToken, "{"},
		{lexer.TextToken, "\n    "},
		{lexer.ElementToken, "#title"},
		{lexer.BraceOpenToken, "{"},
		{lexer.TextToken, "A short title"},
		{lexer.BraceCloseToken, "}"},
		{lexer.TextToken, "\n\n    This is some text with some "},
		{lexer.ElementToken, "#bold"},
		{lexer.BraceOpenToken, "{"},
		{lexer.TextToken, "bold"},
		{lexer.BraceCloseToken, "}"},
		{lexer.TextToken, " text\n"},
		{lexer.BraceCloseToken, "}"},
		{lexer.EOFToken, ""},
	}, tokens)
	assert.Nil(t, err)
}

var example4 = strings.TrimSpace(`
#code {{
    Some raw #bold{ nodes }
}}
`)

func TestLexer4(t *testing.T) {
	s := strings.NewReader(example4)

	tokens, err := lexer.New(s).AllTokens()

	assert.Equal(t, []lexer.Token{
		{lexer.ElementToken, "#code"},
		{lexer.BraceOpenToken, "{{"},
		{lexer.TextToken, "\n    Some raw #bold{ nodes }\n"},
		{lexer.BraceCloseToken, "}}"},
		{lexer.EOFToken, ""},
	}, tokens)
	assert.Nil(t, err)
}

var example5 = strings.TrimSpace(`
#code {{
    #format {{ js }}
    Some raw #bold{ nodes }
}}
`)

func TestLexer5(t *testing.T) {
	s := strings.NewReader(example5)

	tokens, err := lexer.New(s).AllTokens()

	assert.Equal(t, []lexer.Token{
		{lexer.ElementToken, "#code"},
		{lexer.BraceOpenToken, "{{"},
		{lexer.TextToken, "\n    "},
		{lexer.ElementToken, "#format"},
		{lexer.BraceOpenToken, "{{"},
		{lexer.TextToken, "js"},
		{lexer.BraceCloseToken, "}}"},
		{lexer.TextToken, "\n    Some raw #bold{ nodes }\n"},
		{lexer.BraceCloseToken, "}}"},
		{lexer.EOFToken, ""},
	}, tokens)
	assert.Nil(t, err)
}

func TestLexer6(t *testing.T) {
	s := strings.NewReader(`#sum{ 1 }{ 2 }{ #sum{{ 3 }}{{{ 4 }}} }`)

	tokens, err := lexer.New(s).AllTokens()

	assert.Equal(t, []lexer.Token{
		{lexer.ElementToken, "#sum"},
		{lexer.BraceOpenToken, "{"},
		{lexer.TextToken, "1"},
		{lexer.BraceCloseToken, "}"},
		{lexer.BraceOpenToken, "{"},
		{lexer.TextToken, "2"},
		{lexer.BraceCloseToken, "}"},
		{lexer.BraceOpenToken, "{"},
		{lexer.ElementToken, "#sum"},
		{lexer.BraceOpenToken, "{{"},
		{lexer.TextToken, "3"},
		{lexer.BraceCloseToken, "}}"},
		{lexer.BraceOpenToken, "{{{"},
		{lexer.TextToken, "4"},
		{lexer.BraceCloseToken, "}}}"},
		{lexer.BraceCloseToken, "}"},
		{lexer.EOFToken, ""},
	}, tokens)
	assert.Nil(t, err)
}
