package ast

// FirstElement returns the first [ast.ElementNode] in this block or nil otherwise.
func (b Block) FirstElement() *ElementNode {
	for _, n := range b {
		if elem, ok := n.(*ElementNode); ok {
			return elem
		}
	}
	return nil
}

// TextContent concatenates all [ast.TextNode] text in this block (for now [ast.ElementNode]s are skipped)
func (b Block) TextContent() string {
	s := ""
	for _, n := range b {
		if n, ok := n.(*TextNode); ok {
			s += n.Text
		}
	}
	return s
}

// WalkNodes walks the AST using depth-first pre-order traversal (first the visit the node itself and then all its children). The traversal finishes if the visit function returns a non nil error.
func (b Block) Walk(visitFunc func(Node) error) error {
	for _, node := range b {
		err := visitFunc(node)
		if err != nil {
			return err
		}

		if elem, ok := node.(*ElementNode); ok {
			for _, arg := range elem.Arguments {
				arg.Walk(visitFunc)
			}
		}
	}

	return nil
}

// WalkTypes calls the right visit function based on node type. For traversal information see [ast.Walk].
func (b Block) WalkTypes(
	visitElemFunc func(*ElementNode) error,
	visitTextFunc func(*TextNode) error,
) error {
	for _, node := range b {
		switch node := node.(type) {
		case *ElementNode:
			err := visitElemFunc(node)
			if err != nil {
				return err
			}

			for _, arg := range node.Arguments {
				arg.WalkTypes(visitElemFunc, visitTextFunc)
			}
		case *TextNode:
			err := visitTextFunc(node)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
