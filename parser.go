package textml

import (
	"fmt"
	"io"
	"unicode"

	"github.com/aziis98/go-text-ml/parsenstein"
)

type DocumentNode struct {
	Roots []BlockNode
}

// BlockNode is an enum
type BlockNode struct {
	Element *BlockElementNode
	Text    []*InlineNode
}

type BlockElementNode struct {
	Name string
	Body []BlockNode
}

// InlineNode is an enum
type InlineNode struct {
	Text    string
	Element *InlineElementNode
}

type InlineElementNode struct {
	Name string
	Args [][]*InlineNode
}

//
// Parsing
//

func not(b bool) bool {
	return !b
}

type ParseContext struct {
	column int
	line   int
	indent int
}

type parser struct {
	*parsenstein.Parser[ParseContext]
}

func NewParser(rr io.RuneReader) *parser {
	return &parser{
		parsenstein.New(rr, ParseContext{0, 0, 0}, func(p *ParseContext, r rune) {
			p.column++

			if r == '\n' {
				p.line++
				p.column = 0
			}
		}),
	}
}

func ParseBlockElementNode(p *parser, depth int) (*BlockElementNode, error) {
	p.Context().indent = p.Context().column

	r, err := p.NextRune()
	if err != nil {
		return nil, err
	}
	if r != '#' {
		return nil, fmt.Errorf(`expected "#", but got "%v"`, r)
	}

	p.Begin()
	for {
		r, err := p.PeekRune()
		if err != nil {
			return nil, err
		}

		if not(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') {
			break
		}

		p.Advance()
	}

	name := p.Buffered()
	p.End()

	p.Begin()
	for {
		r, err := p.PeekRune()
		if err != nil {
			return nil, err
		}

		if not(r == ':') {
			break
		}

		p.Advance()
	}
	colonDepth := len(p.Buffered())
	p.End()

	if colonDepth >= depth {
		blocks, err := ParseBlockNode(colonDepth)
		p.End()

		return &BlockElementNode{
			Name: name,
			Body: blocks,
		}
	} else {
		p.Drop()
		return nil, fmt.Errorf("just an escaped block element")
	}
}

func ParseBlockNode(p *parser, depth int) (*BlockNode, error) {
	r, err := p.PeekRune()
	if err != nil {
		return nil, err
	}

	if r == '#' {
		return
	}
}

func ParseDocument(s string) *DocumentNode {
	roots := []BlockNode{}

	return &DocumentNode{roots}
}
