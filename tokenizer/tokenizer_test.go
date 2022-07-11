package tokenizer_test

import (
	"strings"
	"testing"

	"github.com/aziis98/go-text-ml/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestTokenizer1(t *testing.T) {
	source := `Example text`
	tr := tokenizer.NewRuneTokenizer(strings.NewReader(source))

	tokens, err := tokenizer.ReadAllTokens(tr)
	assert.Nil(t, err)

	assert.Equal(t, []tokenizer.Token{
		{
			Source: "Example text",
			Type:   tokenizer.Text,
			From:   0,
			To:     12,
		},
	}, tokens)
}
