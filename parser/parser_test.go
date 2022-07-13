package parser_test

import (
	"strings"
	"testing"

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

	assert.Equal(t, &parser.Block{
		Children: []parser.Node{
			&parser.ElementNode{
				Name: "sum",
				Args: []*parser.Block{
					{
						Children: []parser.Node{
							&parser.TextNode{
								Text: "1",
							},
						},
					},
					{
						Children: []parser.Node{
							&parser.TextNode{
								Text: "2",
							},
						},
					},
					{
						Children: []parser.Node{
							&parser.ElementNode{
								Name: "sum",
								Args: []*parser.Block{
									{
										Children: []parser.Node{
											&parser.TextNode{
												Text: "3",
											},
										},
									},
									{
										Children: []parser.Node{
											&parser.TextNode{
												Text: "4",
											},
										},
									},
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

	assert.Equal(t, &parser.Block{
		Children: []parser.Node{
			&parser.ElementNode{
				Name: "code",
				Args: []*parser.Block{
					{
						Children: []parser.Node{
							&parser.ElementNode{
								Name: "format",
								Args: []*parser.Block{
									{
										Children: []parser.Node{
											&parser.TextNode{
												Text: "js",
											},
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
			},
		},
	}, document)
}
