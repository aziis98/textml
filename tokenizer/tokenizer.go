package tokenizer

import (
	"io"
	"unicode"

	"github.com/aziis98/go-text-ml/scanner"
)

type TokenType string

var (
	TextType        TokenType = "text"
	NodeType        TokenType = "node"
	OpenBlockType   TokenType = "open-block"
	CloseBlockType  TokenType = "close-block"
	OpenInlineType  TokenType = "open-inline"
	CloseInlineType TokenType = "close-inline"
)

type Token struct {
	Source string
	Type   TokenType

	From, To int
}

// TokenReader
//
//
type TokenReader interface {
	ReadToken() (Token, error)
}

func ReadAllTokens(tr TokenReader) ([]Token, error) {
	tokens := []Token{}

	for {
		t, err := tr.ReadToken()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			break
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

//
type textmlTokenizer struct {
	StackScanner *scanner.StackScanner

	tokens chan Token
	err    chan error
}

func NewTextmlTokenizer(rr io.RuneReader) *textmlTokenizer {
	t := &textmlTokenizer{
		StackScanner: scanner.New(rr),

		tokens: make(chan Token),
		err:    make(chan error),
	}

	go func() {
		t.err <- t.scan()
	}()

	return t
}

func (t *textmlTokenizer) scan() error {
	for {
		r, err := t.StackScanner.Peek()
		if err != nil {
			return err
		}

		if r == '#' {
			t.scanBlock()
			continue
		}
	}
}

func (t *textmlTokenizer) scanBlock() error {
	if _, err := t.StackScanner.Next(); err != nil {
		return err
	}

	for {
		r, err := t.StackScanner.Next()
		if err != nil {
			return err
		}

		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			break
		}
	}
}

func (t *textmlTokenizer) ReadToken() (Token, error) {
	select {
	case t := <-t.tokens:
		return t, nil
	case err := <-t.err:
		return Token{}, err
	}
}
