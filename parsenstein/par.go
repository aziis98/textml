package parsenstein

import "io"

// parseFrame is a stack frame of the parser, the T better not be a pointer value
type parseFrame[T any] struct {
	from    int
	cursor  int
	context T
}

// Parser is a parser with a custom context T, this better not be a pointer value (just for mutable data that follows the cursor and is pushed/popped/copied with it)
type Parser[T any] struct {
	io.RuneReader
	buffer []rune
	stack  []parseFrame[T]

	advanceContextFunc func(*T, rune)
}

func New[T any](rr io.RuneReader, initialContext T, advanceFunc func(*T, rune)) *Parser[T] {
	return &Parser[T]{
		RuneReader: rr,
		buffer:     []rune{},
		stack: []parseFrame[T]{
			{0, 0, initialContext},
		},
		advanceContextFunc: advanceFunc,
	}
}

func (p *Parser[T]) PeekRune() (rune, error) {
	c := p.Cursor()
	if c >= len(p.buffer) {
		r, _, err := p.RuneReader.ReadRune()
		if err != nil {
			return 0, err
		}

		p.buffer = append(p.buffer, r)
	}

	return p.buffer[c], nil
}

func (p *Parser[T]) NextRune() (rune, error) {
	r, err := p.PeekRune()
	if err != nil {
		return 0, err
	}

	p.Advance()

	return r, nil
}

func (p *Parser[T]) Cursor() int {
	return p.stack[len(p.stack)-1].cursor
}

func (p *Parser[T]) Context() *T {
	return &p.stack[len(p.stack)-1].context
}

// Advance error can be ignored if the call is preceded by [Peek]
func (p *Parser[T]) Advance() error {
	r, err := p.PeekRune()
	if err != nil {
		return err
	}

	p.advanceContextFunc(p.Context(), r)
	p.stack[len(p.stack)-1].cursor++

	return nil
}

func (p *Parser[T]) Begin() {
	frame := p.stack[len(p.stack)-1]
	frame.from = frame.cursor
	p.stack = append(p.stack, frame)
}

func (p *Parser[T]) End() {
	topFrame := p.stack[len(p.stack)-1]
	frame := &p.stack[len(p.stack)-2]

	frame.context = topFrame.context
	frame.cursor = topFrame.cursor

	p.stack = p.stack[:len(p.stack)-1]
}

func (p *Parser[T]) Drop() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *Parser[T]) Buffered() string {
	frame := p.stack[len(p.stack)-1]
	return string(p.buffer[frame.from:frame.cursor])
}
