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
	assert.Equal(t, []*lexer.Token{
		{lexer.TextToken, "Lorem ", lexer.TokenInfo{0, 0}},
		{lexer.ElementToken, "#node", lexer.TokenInfo{0, 6}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{0, 11}},
		{lexer.TextToken, "ipsum", lexer.TokenInfo{0, 12}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{0, 17}},
		{lexer.TextToken, " dolor", lexer.TokenInfo{0, 18}},
		{lexer.EOFToken, "", lexer.TokenInfo{0, 24}},
	}, tokens)
}
func TestLexer2(t *testing.T) {
	s := strings.NewReader("Lorem #node{{ipsum}}} dolor")
	tokens, err := lexer.New(s).AllTokens()

	assert.Nil(t, tokens)
	assert.Equal(t, "too many braces at 21", err.Error())
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

	assert.Equal(t, []*lexer.Token{
		{lexer.ElementToken, "#document", lexer.TokenInfo{Line: 0, Column: 0}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{Line: 0, Column: 10}},
		{lexer.TextToken, "\n    ", lexer.TokenInfo{Line: 0, Column: 11}},
		{lexer.ElementToken, "#title", lexer.TokenInfo{Line: 1, Column: 4}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{Line: 1, Column: 11}},
		{lexer.TextToken, "A short title", lexer.TokenInfo{Line: 1, Column: 13}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{Line: 1, Column: 27}},
		{lexer.TextToken, "\n\n    This is some text with some ", lexer.TokenInfo{Line: 1, Column: 28}},
		{lexer.ElementToken, "#bold", lexer.TokenInfo{Line: 3, Column: 32}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{Line: 3, Column: 37}},
		{lexer.TextToken, "bold", lexer.TokenInfo{Line: 3, Column: 39}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{Line: 3, Column: 44}},
		{lexer.TextToken, " text\n", lexer.TokenInfo{Line: 3, Column: 45}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{Line: 4, Column: 0}},
		{lexer.EOFToken, "", lexer.TokenInfo{Line: 4, Column: 1}},
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

	assert.Equal(t, []*lexer.Token{
		{lexer.ElementToken, "#code", lexer.TokenInfo{0, 0}},
		{lexer.BraceOpenToken, "{{", lexer.TokenInfo{0, 6}},
		{lexer.TextToken, "\n    Some raw #bold{ nodes }\n", lexer.TokenInfo{0, 8}},
		{lexer.BraceCloseToken, "}}", lexer.TokenInfo{2, 0}},
		{lexer.EOFToken, "", lexer.TokenInfo{2, 2}},
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

	assert.Equal(t, []*lexer.Token{
		{lexer.ElementToken, "#code", lexer.TokenInfo{0, 0}},
		{lexer.BraceOpenToken, "{{", lexer.TokenInfo{0, 6}},
		{lexer.TextToken, "\n    ", lexer.TokenInfo{0, 8}},
		{lexer.ElementToken, "#format", lexer.TokenInfo{1, 4}},
		{lexer.BraceOpenToken, "{{", lexer.TokenInfo{1, 12}},
		{lexer.TextToken, "js", lexer.TokenInfo{1, 15}},
		{lexer.BraceCloseToken, "}}", lexer.TokenInfo{1, 18}},
		{lexer.TextToken, "\n    Some raw #bold{ nodes }\n", lexer.TokenInfo{1, 20}},
		{lexer.BraceCloseToken, "}}", lexer.TokenInfo{3, 0}},
		{lexer.EOFToken, "", lexer.TokenInfo{3, 2}},
	}, tokens)
	assert.Nil(t, err)
}

func TestLexer6(t *testing.T) {
	s := strings.NewReader(`#sum{ 1 }{ 2 }{ #sum{{ 3 }}{{{ 4 }}} }`)

	tokens, err := lexer.New(s).AllTokens()

	assert.Equal(t, []*lexer.Token{
		{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 0}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{0, 4}},
		{lexer.TextToken, "1", lexer.TokenInfo{0, 6}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{0, 8}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{0, 9}},
		{lexer.TextToken, "2", lexer.TokenInfo{0, 11}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{0, 13}},
		{lexer.BraceOpenToken, "{", lexer.TokenInfo{0, 14}},
		{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 16}},
		{lexer.BraceOpenToken, "{{", lexer.TokenInfo{0, 20}},
		{lexer.TextToken, "3", lexer.TokenInfo{0, 23}},
		{lexer.BraceCloseToken, "}}", lexer.TokenInfo{0, 25}},
		{lexer.BraceOpenToken, "{{{", lexer.TokenInfo{0, 27}},
		{lexer.TextToken, "4", lexer.TokenInfo{0, 31}},
		{lexer.BraceCloseToken, "}}}", lexer.TokenInfo{0, 33}},
		{lexer.BraceCloseToken, "}", lexer.TokenInfo{0, 37}},
		{lexer.EOFToken, "", lexer.TokenInfo{0, 38}},
	}, tokens)
	assert.Nil(t, err)
}
