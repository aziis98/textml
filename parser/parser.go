package parser

import (
	"fmt"

	"github.com/aziis98/textml/lexer"
)

type Node interface {
	sealNode()
}

func (TextNode) sealNode()    {}
func (ElementNode) sealNode() {}

// Block represents a sequence of parsed node in a document, holds references to the first and last parsed tokens.
type Block struct {
	BeginToken, EndToken *lexer.Token

	Children []Node
}

// TextNode represents a text node with token information.
type TextNode struct {
	*lexer.Token
	Text string
}

// ElementNode represents an element node with name, arguments and token information.
type ElementNode struct {
	*lexer.Token
	Name string
	Args []*Block
}

// Parse creates a parse AST, this keeps token information if one wants to do low level processing after the parse.
func Parse(ts []*lexer.Token) (*Block, error) {
	children := []Node{}
	begin := ts[0]
	t := ts[0]
	for len(ts) > 0 && t.Type != lexer.EOFToken {
		switch t.Type {
		case lexer.TextToken:
			ts = ts[1:]
			children = append(children, &TextNode{Text: t.Value})
		case lexer.ElementToken:
			var elt Node
			var err error

			elt, ts, err = parseElement(ts)
			if err != nil {
				fmt.Printf("rest: %v\n", ts)
				return nil, err
			}

			children = append(children, elt)
		default:
			fmt.Printf("rest: %v\n", ts)
			return nil, fmt.Errorf("[document] expected text or element, got: %v", t)
		}

		t = ts[0]
	}

	return &Block{begin, t, children}, nil
}

// parseElement expects an element and parses it into a [*parser.ElementNode].
func parseElement(ts []*lexer.Token) (*ElementNode, []*lexer.Token, error) {
	if len(ts) == 0 {
		return nil, ts, fmt.Errorf("[element] not enough tokens")
	}

	t := ts[0]
	if t.Type != lexer.ElementToken {
		return nil, ts, fmt.Errorf("[element] expected element, got: %v", t)
	}

	elemToken := t
	name := t.Value[1:]

	ts = ts[1:]
	t = ts[0]

	blocks := []*Block{}

	for len(ts) > 0 && t.Type == lexer.BraceOpenToken {
		var blk *Block
		var err error

		blk, ts, err = parseElementArgument(ts)
		if err != nil {
			return nil, ts, err
		}

		blocks = append(blocks, blk)

		t = ts[0]
	}

	return &ElementNode{
		Token: elemToken,

		Name: name,
		Args: blocks,
	}, ts, nil
}

func parseElementArgument(ts []*lexer.Token) (*Block, []*lexer.Token, error) {
	if ts[0].Type != lexer.BraceOpenToken {
		return nil, ts, fmt.Errorf("[argument] expected opening brace, got: %v", ts[0])
	}

	ts = ts[1:]
	begin := ts[0] // first token after brace
	end := ts[0]   // last token before brace

	children := []Node{}

	for {
		if len(ts) == 0 {
			return nil, ts, fmt.Errorf("[argument] unbalanced block")
		}

		t := ts[0]

		if t.Type == lexer.BraceCloseToken {
			ts = ts[1:]
			return &Block{begin, end, children}, ts, nil
		}

		switch t.Type {
		case lexer.TextToken:
			ts = ts[1:]
			children = append(children, &TextNode{t, t.Value})
			end = t

		case lexer.ElementToken:
			elt, tss, err := parseElement(ts)
			if err != nil {
				return nil, ts, err
			}

			children = append(children, elt)

			end = ts[len(ts)-len(tss)-1]
			ts = tss

		default:
			return nil, ts, fmt.Errorf("[argument] expected text or element, got: %v", t)

		}
	}
}
