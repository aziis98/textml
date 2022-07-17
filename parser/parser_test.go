package parser_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/lexer"
	"github.com/aziis98/textml/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser1(t *testing.T) {
	s := strings.NewReader(`#sum{ 1 }{ 2 }{ #sum{{ 3 }}{{{ 4 }}} }`)

	tokens, err := lexer.New(s).AllTokens()
	assert.Nil(t, err)

	document, err := parser.ParseDocument(tokens)
	assert.Nil(t, err)

	assert.Equal(t, parser.Block{

		&parser.ElementNode{
			Name: "sum",
			Args: []parser.Block{
				{
					&parser.TextNode{
						Text: "1",
					},
				},
				{
					&parser.TextNode{
						Text: "2",
					},
				},
				{
					&parser.ElementNode{
						Name: "sum",
						Args: []parser.Block{
							{
								&parser.TextNode{
									Text: "3",
								},
							},
							{
								&parser.TextNode{
									Text: "4",
								},
							},
						},
					},
				},
			},
		},
	}, document)
}

func TestPaser2(t *testing.T) {
	s := strings.NewReader(`#code{{ #format{{ js }} let x = "#node{ 1 }"; }}`)

	tokens, err := lexer.New(s).AllTokens()
	assert.Nil(t, err)

	document, err := parser.ParseDocument(tokens)
	assert.Nil(t, err)

	assert.Equal(t, parser.Block{
		&parser.ElementNode{
			Name: "code",
			Args: []parser.Block{
				{
					&parser.ElementNode{
						Name: "format",
						Args: []parser.Block{
							{
								&parser.TextNode{
									Text: "js",
								},
							},
						},
					},
					&parser.TextNode{
						Text: ` let x = "#node{ 1 }";`,
					},
				},
			},
		},
	}, document)
}

func TestJsonConversion(t *testing.T) {
	doc, err := textml.ParseDocument(strings.NewReader("#foo{#bar{a}#baz{b}c}"))
	assert.Nil(t, err)

	docJson, err := json.Marshal(doc)
	assert.Nil(t, err)

	assert.Equal(t, `[{"args":[[{"args":[[{"text":"a","type":"text"}]],"name":"bar","type":"element"},{"args":[[{"text":"b","type":"text"}]],"name":"baz","type":"element"},{"text":"c","type":"text"}]],"name":"foo","type":"element"}]`, string(docJson))
}
