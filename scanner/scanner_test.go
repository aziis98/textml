package scanner_test

import (
	"io"
	"strings"
	"testing"

	"github.com/aziis98/go-text-ml/scanner"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var r rune

	s := scanner.New(strings.NewReader("abcd"))
	{

		r, _ = s.Next()
		assert.Equal(t, 'a', r)

		r, _ = s.Next()
		assert.Equal(t, 'b', r)

		s.RaiseCursor()

		r, _ = s.Next()
		assert.Equal(t, 'c', r)

		r, _ = s.Next()
		assert.Equal(t, 'd', r)

		s.DropCursor()

		r, _ = s.Next()
		assert.Equal(t, 'c', r)

		r, _ = s.Next()
		assert.Equal(t, 'd', r)

		_, err := s.Next()
		assert.Equal(t, io.EOF, err)
	}
}
