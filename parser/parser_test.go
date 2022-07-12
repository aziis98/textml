package parser_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aziis98/go-text-ml/lexer"
	"github.com/aziis98/go-text-ml/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser1(t *testing.T) {
	s := strings.NewReader(`#sum{ 1 }{ 2 }{ #sum{{ 3 }}{{{ 4 }}} }`)

	tokens, err := lexer.New(s).AllTokens()
	assert.Nil(t, err)

	document, err := parser.ParseDocument(tokens)
	assert.Nil(t, err)

	doc, _ := json.MarshalIndent(document, "", "  ")
	fmt.Println(string(doc))
}

func TestPaser2(t *testing.T) {
	s := strings.NewReader(`#code{{ #format{{ js }} let x = "#node{ 1 }"; }}`)

	tokens, err := lexer.New(s).AllTokens()
	assert.Nil(t, err)

	document, err := parser.ParseDocument(tokens)
	assert.Nil(t, err)

	doc, _ := json.MarshalIndent(document, "", "  ")
	fmt.Println(string(doc))
}
