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

	document, err := parser.Parse(tokens)
	assert.Nil(t, err)

	assert.Equal(t,
		&parser.Block{
			BeginToken: &lexer.Token{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 0}},
			EndToken:   &lexer.Token{lexer.EOFToken, "", lexer.TokenInfo{0, 38}},
			Children: []parser.Node{
				&parser.ElementNode{
					Token: &lexer.Token{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 0}},
					Name:  "sum",
					Args: []*parser.Block{
						{
							BeginToken: &lexer.Token{lexer.TextToken, "1", lexer.TokenInfo{0, 6}},
							EndToken:   &lexer.Token{lexer.TextToken, "1", lexer.TokenInfo{0, 6}},
							Children: []parser.Node{
								&parser.TextNode{
									Token: &lexer.Token{lexer.TextToken, "1", lexer.TokenInfo{0, 6}},
									Text:  "1",
								},
							},
						},
						{
							BeginToken: &lexer.Token{lexer.TextToken, "2", lexer.TokenInfo{0, 11}},
							EndToken:   &lexer.Token{lexer.TextToken, "2", lexer.TokenInfo{0, 11}},
							Children: []parser.Node{
								&parser.TextNode{
									Token: &lexer.Token{lexer.TextToken, "2", lexer.TokenInfo{0, 11}},
									Text:  "2",
								},
							},
						},
						{
							BeginToken: &lexer.Token{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 16}},
							EndToken:   &lexer.Token{lexer.BraceCloseToken, "}}}", lexer.TokenInfo{0, 33}},
							Children: []parser.Node{
								&parser.ElementNode{
									Token: &lexer.Token{lexer.ElementToken, "#sum", lexer.TokenInfo{0, 16}},
									Name:  "sum",
									Args: []*parser.Block{
										{
											BeginToken: &lexer.Token{lexer.TextToken, "3", lexer.TokenInfo{0, 23}},
											EndToken:   &lexer.Token{lexer.TextToken, "3", lexer.TokenInfo{0, 23}},
											Children: []parser.Node{
												&parser.TextNode{
													Token: &lexer.Token{lexer.TextToken, "3", lexer.TokenInfo{0, 23}},
													Text:  "3",
												}},
										},
										{
											BeginToken: &lexer.Token{lexer.TextToken, "4", lexer.TokenInfo{0, 31}},
											EndToken:   &lexer.Token{lexer.TextToken, "4", lexer.TokenInfo{0, 31}},
											Children: []parser.Node{
												&parser.TextNode{
													Token: &lexer.Token{lexer.TextToken, "4", lexer.TokenInfo{0, 31}},
													Text:  "4",
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

func TestParser2(t *testing.T) {
	s := strings.NewReader(`#code{{ #format{{ js }} let x = "#node{ 1 }"; }}`)

	tokens, err := lexer.New(s).AllTokens()
	assert.Nil(t, err)

	document, err := parser.Parse(tokens)
	assert.Nil(t, err)

	assert.Equal(t,
		&parser.Block{
			BeginToken: &lexer.Token{lexer.ElementToken, "#code", lexer.TokenInfo{0, 0}},
			EndToken:   &lexer.Token{lexer.EOFToken, "", lexer.TokenInfo{0, 48}},
			Children: []parser.Node{
				&parser.ElementNode{
					Token: &lexer.Token{lexer.ElementToken, "#code", lexer.TokenInfo{0, 0}},
					Name:  "code",
					Args: []*parser.Block{
						{
							BeginToken: &lexer.Token{lexer.ElementToken, "#format", lexer.TokenInfo{0, 8}},
							EndToken: &lexer.Token{
								lexer.TextToken, ` let x = "#node{ 1 }";`, lexer.TokenInfo{0, 23},
							},
							Children: []parser.Node{
								&parser.ElementNode{
									Token: &lexer.Token{lexer.ElementToken, "#format", lexer.TokenInfo{0, 8}},
									Name:  "format",
									Args: []*parser.Block{
										{
											BeginToken: &lexer.Token{lexer.TextToken, "js", lexer.TokenInfo{0, 18}},
											EndToken:   &lexer.Token{lexer.TextToken, "js", lexer.TokenInfo{0, 18}},
											Children: []parser.Node{
												&parser.TextNode{
													Token: &lexer.Token{lexer.TextToken, "js", lexer.TokenInfo{0, 18}},
													Text:  "js",
												},
											},
										},
									},
								},
								&parser.TextNode{
									Token: &lexer.Token{
										lexer.TextToken, ` let x = "#node{ 1 }";`, lexer.TokenInfo{0, 23},
									},
									Text: ` let x = "#node{ 1 }";`,
								},
							},
						},
					},
				},
			},
		}, document)
}
