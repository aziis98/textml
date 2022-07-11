package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aziis98/go-text-ml/lexer"
)

func TestLexer1(t *testing.T) {
	s := strings.NewReader("Lorem #node[ipsum] dolor")
	tokens, err := lexer.New(s).AllTokens()

	fmt.Println(tokens)
	fmt.Println(err)
}
