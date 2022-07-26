package ast

import (
	"fmt"

	"github.com/aziis98/textml/parser"
)

// Compile a [*parser.Block] into a [Block] instance (for now just strips token information).
func Compile(block *parser.Block) Block {
	nodes := Block{}
	for _, child := range block.Children {
		switch child := child.(type) {
		case *parser.TextNode:
			nodes = append(nodes, &TextNode{
				Text: child.Text,
			})

		case *parser.ElementNode:
			args := []Block{}
			for _, block := range child.Args {
				args = append(args, Compile(block))
			}

			nodes = append(nodes, &ElementNode{
				Name:      child.Name,
				Arguments: args,
			})

		default:
			panic(fmt.Errorf("unexpected node of type: %T", child))
		}
	}

	return nodes
}
