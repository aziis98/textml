package html

import (
	"fmt"
	"io"
	"strings"
)

type WriterTo interface {
	Write(w io.Writer)
}

type Node interface {
	WriterTo
	htmlSeal()
}

type AttributeNode interface {
	WriterTo
	attributeSeal()
}

type AttributeMap map[string]AttributeNode

type EmptyAttribute struct{}

func (EmptyAttribute) attributeSeal() {}

func (EmptyAttribute) Write(w io.Writer) {}

type Attribute struct {
	Value string
}

func (Attribute) attributeSeal() {}

func (a Attribute) Write(w io.Writer) {
	fmt.Fprintf(w, `="%s"`, a.Value)
}

type Element struct {
	TagName    string
	Attributes map[string]AttributeNode
	Children   []Node
}

func NewElementNode(tagName string, attributes map[string]AttributeNode, children []Node) *Element {
	if attributes == nil {
		attributes = map[string]AttributeNode{}
	}
	if children == nil {
		children = []Node{}
	}

	return &Element{tagName, attributes, children}
}

func (Element) htmlSeal() {}

func (n Element) Write(w io.Writer) {
	attrs := &strings.Builder{}
	for k, attr := range n.Attributes {
		fmt.Fprintf(attrs, " %s", k)
		attr.Write(attrs)
	}

	fmt.Fprintf(w, "<%s%s>", n.TagName, attrs.String())

	for _, child := range n.Children {
		child.Write(w)
	}

	fmt.Fprintf(w, "</%s>", n.TagName)
}

type Text struct {
	Value string
}

func NewTextNode(s string) *Text {
	return &Text{s}
}

func (Text) htmlSeal() {}

func (n Text) Write(w io.Writer) {
	fmt.Fprintf(w, `%s`, n.Value)
}

type Comment struct {
	Value string
}

func NewCommentNode(s string) *Comment {
	return &Comment{s}
}

func (Comment) htmlSeal() {}

func (n Comment) Write(w io.Writer) {
	fmt.Fprintf(w, `<!-- %s -->`, n.Value)
}

func RenderToString(nodes []Node) string {
	sb := &strings.Builder{}
	for _, n := range nodes {
		n.Write(sb)
	}
	return sb.String()
}
