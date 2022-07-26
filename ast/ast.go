package ast

import "encoding/json"

// Node is an abstract syntax tree representation of the TextML without the parsing information and mostly used during transformations.
type Node interface{ sealNode() }

// Block is the top most unit in a TextML document, it's an alias to a slice of [Node]s for easy construction.
type Block []Node

// FirstElement returns the first [*ElementNode] in this block or nil otherwise.
func (b Block) FirstElement() *ElementNode {
	for _, n := range b {
		if elem, ok := n.(*ElementNode); ok {
			return elem
		}
	}

	return nil
}

// TextContent concatenates all [*TextNode] text in this block (for now [*ElementNode]s are skipped)
func (b Block) TextContent() string {
	s := ""

	for _, n := range b {
		if n, ok := n.(*TextNode); ok {
			s += n.Text
		}
	}

	return s
}

// TextNode represents a text node
type TextNode struct {
	Text string
}

func (TextNode) sealNode() {}

func (n *TextNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "text",
		"text": n.Text,
	})
}

// ElementNode represents an element node with its block arguments
type ElementNode struct {
	Name      string
	Arguments []Block
}

func (ElementNode) sealNode() {}

func (n *ElementNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "element",
		"name": n.Name,
		"args": n.Arguments,
	})
}
