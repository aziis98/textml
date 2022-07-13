package parser

import (
	"encoding/json"
	"fmt"

	"github.com/aziis98/go-text-ml/lexer"
)

type Node interface {
	nodeType()
}

// Block is a list of nodes
type Block struct {
	Children []Node `json:"children"`
}

type TextNode struct {
	Text string
}

func (_ *TextNode) nodeType() {}

func (n *TextNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "text",
		"text": n.Text,
	})
}

type ElementNode struct {
	Name string
	Args []*Block
}

func (_ *ElementNode) nodeType() {}

func (n *ElementNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "element",
		"name": n.Name,
		"args": n.Args,
	})
}

func ParseDocument(ts []lexer.Token) (*Block, error) {
	children := []Node{}

	t := ts[0]
	for len(ts) > 0 && t.Type != lexer.EOFToken {
		switch t.Type {
		case lexer.TextToken:
			ts = ts[1:]
			children = append(children, &TextNode{Text: t.Value})
		case lexer.ElementToken:
			var elt Node
			var err error

			elt, ts, err = ParseElement(ts)
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

	return &Block{children}, nil
}

func ParseElement(ts []lexer.Token) (Node, []lexer.Token, error) {
	if len(ts) == 0 {
		return nil, ts, fmt.Errorf("[element] not enough tokens")
	}

	t := ts[0]
	if t.Type != lexer.ElementToken {
		return nil, ts, fmt.Errorf("[element] expected element, got: %v", t)
	}

	name := t.Value[1:]

	ts = ts[1:]
	t = ts[0]

	blocks := []*Block{}

	for len(ts) > 0 && t.Type == lexer.BraceOpenToken {
		var blk *Block
		var err error

		blk, ts, err = ParseArgument(ts)
		if err != nil {
			return nil, ts, err
		}

		blocks = append(blocks, blk)

		t = ts[0]
	}

	return &ElementNode{
		Name: name,
		Args: blocks,
	}, ts, nil
}

func ParseArgument(ts []lexer.Token) (*Block, []lexer.Token, error) {
	if ts[0].Type != lexer.BraceOpenToken {
		return nil, ts, fmt.Errorf("[argument] expected openning brace, got: %v", ts[0])
	}

	ts = ts[1:]

	children := []Node{}

	for {
		if len(ts) == 0 {
			return nil, ts, fmt.Errorf("[argument] unbalanced block")
		}

		t := ts[0]

		if t.Type == lexer.BraceCloseToken {
			ts = ts[1:]
			return &Block{children}, ts, nil
		}

		switch t.Type {
		case lexer.TextToken:
			ts = ts[1:]
			children = append(children, &TextNode{t.Value})
		case lexer.ElementToken:
			var elt Node
			var err error

			elt, ts, err = ParseElement(ts)
			if err != nil {
				return nil, ts, err
			}

			children = append(children, elt)
		default:
			return nil, ts, fmt.Errorf("[argument] expected text or element, got: %v", t)
		}
	}
}
