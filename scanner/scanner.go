package scanner

import (
	"fmt"
	"io"
)

type StackScanner struct {
	io.RuneReader
	buffer      []rune
	cursorStack []int
}

func New(r io.RuneReader) *StackScanner {
	return &StackScanner{
		RuneReader:  r,
		buffer:      []rune{},
		cursorStack: []int{0},
	}
}

func (rs *StackScanner) cursor() *int {
	return &rs.cursorStack[len(rs.cursorStack)-1]
}

func (rs *StackScanner) Peek() (rune, error) {
	if *rs.cursor() >= len(rs.buffer) {
		r, _, err := rs.ReadRune()
		if err != nil {
			return 0, err
		}

		rs.buffer = append(rs.buffer, r)
	}

	return rs.buffer[*rs.cursor()], nil
}

func (rs *StackScanner) Next() (rune, error) {
	fmt.Printf("%+v\n", rs)

	r, err := rs.Peek()
	if err != nil {
		return 0, err
	}

	*rs.cursor() = *rs.cursor() + 1
	return r, nil
}

// RaiseCursor duplicates the cursor position on top of the stack
func (rs *StackScanner) RaiseCursor() {
	rs.cursorStack = append(rs.cursorStack, *rs.cursor())
}

// DropCursor pops the top cursor from the stack
func (rs *StackScanner) DropCursor() {
	rs.cursorStack = rs.cursorStack[:len(rs.cursorStack)-1]
}

// LowerCursor takes the top cursor and puts it at the next lower stack level
func (rs *StackScanner) LowerCursor() (string, int, int) {
	last := *rs.cursor()
	rs.cursorStack = rs.cursorStack[:len(rs.cursorStack)-1]
	prev := *rs.cursor()

	*rs.cursor() = last

	return string(rs.buffer[prev:last]), prev, last
}
