package parser

import (
	"fmt"

	"github.com/aziis98/go-text-ml/lexer"
)

var NodeTypes = struct {
	Text    string
	Element string
}{
	Text:    "text",
	Element: "element",
}

// Block is a list of nodes
type Block struct {
	Children []*Node `json:"children"`
}

type TextNode struct {
	Text string `json:"text"`
}

type ElementNode struct {
	Name string   `json:"name"`
	Args []*Block `json:"arguments"`
}

// Node is an enum of text or element nodes
type Node struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

func ParseDocument(ts []lexer.Token) (*Block, error) {
	children := []*Node{}

	t := ts[0]
	for len(ts) > 0 && t.Type != lexer.EOFToken {
		switch t.Type {
		case lexer.TextToken:
			ts = ts[1:]
			children = append(children, &Node{
				Type:  NodeTypes.Text,
				Value: &TextNode{t.Value},
			})
		case lexer.ElementToken:
			var elt *ElementNode
			var err error

			elt, ts, err = ParseElement(ts)
			if err != nil {
				fmt.Printf("rest: %v\n", ts)
				return nil, err
			}

			children = append(children, &Node{
				Type:  NodeTypes.Element,
				Value: elt,
			})
		default:
			fmt.Printf("rest: %v\n", ts)
			return nil, fmt.Errorf("[document] expected text or element, got: %v", t)
		}

		t = ts[0]
	}

	return &Block{children}, nil
}

func ParseElement(ts []lexer.Token) (*ElementNode, []lexer.Token, error) {
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

	for len(ts) > 0 && t.Type != lexer.TextToken && t.Type != lexer.EOFToken {
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

	children := []*Node{}

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
			children = append(children, &Node{
				Type:  NodeTypes.Text,
				Value: &TextNode{t.Value},
			})
		case lexer.ElementToken:
			var elt *ElementNode
			var err error

			elt, ts, err = ParseElement(ts)
			if err != nil {
				return nil, ts, err
			}

			children = append(children, &Node{
				Type:  NodeTypes.Element,
				Value: elt,
			})
		default:
			return nil, ts, fmt.Errorf("[argument] expected text or element, got: %v", t)
		}
	}
}
