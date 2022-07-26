package lexer

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/aziis98/textml/utils"
)

type tokenType int

const (
	EOFToken tokenType = iota
	TextToken
	ElementToken
	BraceOpenToken
	BraceCloseToken
)

func (t tokenType) GoString() string {
	switch t {
	case EOFToken:
		return "lexer.EOFToken"
	case TextToken:
		return "lexer.TextToken"
	case ElementToken:
		return "lexer.ElementToken"
	case BraceOpenToken:
		return "lexer.BraceOpenToken"
	case BraceCloseToken:
		return "lexer.BraceCloseToken"
	default:
		panic(fmt.Errorf("illegal token type: %d", t))
	}
}

const eof rune = 0

type TokenInfo struct {
	Line, Column int
}

func (ti TokenInfo) GoString() string {
	return fmt.Sprintf("TokenInfo{%d, %d}", ti.Line, ti.Column)
}

type Token struct {
	Type  tokenType
	Value string

	TokenInfo TokenInfo
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%v, %q, %d:%d}",
		t.Type,
		t.Value,
		t.TokenInfo.Line,
		t.TokenInfo.Column,
	)
}

type lexer struct {
	io.RuneReader

	buf     []rune
	bufFrom int
	bufTo   int

	pos int

	over   bool
	done   chan struct{}
	tokens []*Token
	err    error

	bracesStack *utils.Stack[int]

	tokenInfo TokenInfo
}

func New(rr io.RuneReader) *lexer {
	l := &lexer{
		RuneReader: rr,

		buf:     []rune{},
		bufFrom: 0,
		bufTo:   0,

		pos: 0,

		over:   false,
		done:   make(chan struct{}),
		tokens: []*Token{},
		err:    nil,

		bracesStack: utils.NewStack(1),

		tokenInfo: TokenInfo{0, 0},
	}

	go l.scan()

	return l
}

func (l *lexer) AllTokens() ([]*Token, error) {
	if !l.over {
		<-l.done
		l.over = true
		close(l.done)
	}

	if l.err != nil {
		return nil, l.err
	}

	return l.tokens, nil
}

func (l *lexer) scan() {
	lexText(l)
	if l.err == nil {
		l.done <- struct{}{}
	}
}

// : from, to    :                  [-----]
// : source      : [--------------------------------------]
// : buffer      :              [-----------------]
// : cursor      :                        ^
// : bufferReach :                                ^

func (l *lexer) bufferOffset() int {
	return l.bufTo - len(l.buf)
}

func (l *lexer) bufferSlice(from, to int) []rune {
	bufferFrom := from - l.bufferOffset()
	bufferTo := to - l.bufferOffset()

	return l.buf[bufferFrom:bufferTo]
}

func (l *lexer) bufferAt(pos int) rune {
	return l.buf[pos-l.bufferOffset()]
}

func (l *lexer) next() rune {
	if l.pos < l.bufTo {
		r := l.bufferAt(l.pos)
		l.pos++
		return r
	}

	r, _, err := l.ReadRune()
	if err != nil {
		if err != io.EOF {
			l.errorf("%v", err)
		}

		return eof
	}

	l.pos++
	l.bufTo++
	l.buf = append(l.buf, r)

	return r
}

func (l *lexer) cursor() int {
	return l.pos
}

func (l *lexer) backup() {
	l.pos--
}

func (l *lexer) move(pos int) {
	if pos < l.bufFrom {
		panic("cannot backtrack before current working token")
	}

	l.pos = pos
}

func (l *lexer) peek() rune {
	r := l.next()

	if r != eof {
		l.backup()
	}

	return r
}

func (l *lexer) emit(tt tokenType) {
	var value string
	value, l.buf = string(l.bufferSlice(l.bufFrom, l.pos)), l.bufferSlice(l.pos, l.bufTo)
	l.bufFrom = l.pos

	if len(value) > 0 || tt == EOFToken {
		t := &Token{tt, value, l.tokenInfo}
		l.tokens = append(l.tokens, t)

		for _, r := range value {
			if r == '\n' {
				l.tokenInfo.Line++
				l.tokenInfo.Column = 0
			} else {
				l.tokenInfo.Column++
			}
		}
	}
}

func (l *lexer) errorf(format string, args ...any) {
	l.err = fmt.Errorf(format, args...)
	l.done <- struct{}{}
}

func (l *lexer) ignore() {
	for _, r := range string(l.bufferSlice(l.bufFrom, l.pos)) {
		if r == '\n' {
			l.tokenInfo.Line++
			l.tokenInfo.Column = 0
		} else {
			l.tokenInfo.Column++
		}
	}

	l.buf = l.bufferSlice(l.pos, l.bufTo)
	l.bufFrom = l.pos
}

func (l *lexer) acceptAny(valid string) bool {
	r := l.next()

	if strings.ContainsRune(valid, r) {
		return true
	}

	l.backup()
	return false
}

func (l *lexer) acceptAnyRepeated(valid string) int {
	size := 0

	for {
		r := l.next()
		if !strings.ContainsRune(valid, r) {
			l.backup()
			break
		}

		size++
	}

	return size
}

func (l *lexer) acceptWhile(validFunc func(rune) bool) int {
	size := 0

	for {
		r := l.next()
		if !validFunc(r) {
			l.backup()
			break
		}

		size++
	}

	return size
}

func lexText(l *lexer) {
	for {
		r := l.peek()

		switch r {
		case eof:
			l.emit(TextToken)
			l.emit(EOFToken)
			return
		case '#':
			// Tries to tokenize an element
			elementStart := l.cursor()

			l.next()
			l.acceptWhile(func(r rune) bool {
				return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.'
			})
			elementEnd := l.cursor()

			l.acceptAnyRepeated(" ")
			spacesEnd := l.cursor()

			newDepth := l.acceptAnyRepeated("{")
			bracesEnd := l.cursor()

			depth := l.bracesStack.Top()

			if newDepth >= depth { // if there are enough braces then accept the element token
				l.move(elementStart) // finish previous text token
				l.emit(TextToken)

				l.move(elementEnd) // emit element token
				l.emit(ElementToken)

				l.move(spacesEnd) // skip whitespace
				l.ignore()

				l.move(bracesEnd) // emit new open brace token
				l.emit(BraceOpenToken)

				if l.acceptAny(" ") { // skip a single whitespace if present after opening brace
					l.ignore()
				}

				l.bracesStack.Push(newDepth)
			}
		case '}':
			bracesStart := l.cursor()
			braceCount := l.acceptAnyRepeated("}")

			depth := l.bracesStack.Top()
			if braceCount == depth {
				if l.bufferAt(bracesStart-1) == ' ' {
					l.move(bracesStart - 1)
					l.emit(TextToken)

					l.next() // skip a single whitespace if present before closing brace
					l.ignore()
				} else {
					l.move(bracesStart)
					l.emit(TextToken)
				}

				l.move(bracesStart + braceCount)
				l.emit(BraceCloseToken)

				l.bracesStack.Pop()

				if l.peek() == '{' { // check if there is another argument for this element
					newDepth := l.acceptAnyRepeated("{")
					bracesEnd := l.cursor()

					depth := l.bracesStack.Top()
					if newDepth >= depth { // if there are enough braces then accept the element token

						l.move(bracesEnd)
						l.emit(BraceOpenToken)

						if l.acceptAny(" ") { // skip a single whitespace if present after brace
							l.ignore()
						}

						l.bracesStack.Push(newDepth)
					}
				}
			} else {
				if braceCount > depth {
					l.emit(EOFToken)
					l.errorf("too many braces at %v", l.pos)
					return
				}
			}
		default:
			l.next()
		}
	}
}
