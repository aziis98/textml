package textml_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aziis98/textml"
	"github.com/stretchr/testify/assert"
)

func TestJsonConversion(t *testing.T) {
	ast, err := textml.ParseDocument(strings.NewReader("#foo{#bar{a}#baz{b}c}"))
	assert.Nil(t, err)

	data, err := json.Marshal(ast)
	assert.Nil(t, err)

	assert.Equal(t, `[{"args":[[{"args":[[{"text":"a","type":"text"}]],"name":"bar","type":"element"},{"args":[[{"text":"b","type":"text"}]],"name":"baz","type":"element"},{"text":"c","type":"text"}]],"name":"foo","type":"element"}]`, string(data))
}
