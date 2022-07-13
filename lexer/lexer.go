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
	TextToken tokenType = iota
	ElementToken
	BraceOpenToken
	BraceCloseToken
	EOFToken
)

const eofRune rune = 0

type Token struct {
	Type  tokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%v, %q}", t.Type, t.Value)
}

type lexer struct {
	io.RuneReader

	buf     []rune
	bufFrom int
	bufTo   int

	pos int

	tokens chan Token
	err    chan error

	bracesStack *utils.Stack[int]
}

func New(rr io.RuneReader) *lexer {
	l := &lexer{
		RuneReader: rr,

		buf:     []rune{},
		bufFrom: 0,
		bufTo:   0,

		pos: 0,

		tokens: make(chan Token, 1),
		err:    make(chan error, 1),

		bracesStack: utils.NewStack(1),
	}

	go l.scan()

	return l
}

func (l *lexer) NextToken() (Token, error) {
	select {
	case t := <-l.tokens:
		return t, nil
	case err := <-l.err:
		return Token{}, err
	}
}

func (l *lexer) AllTokens() ([]Token, error) {
	tokens := []Token{}
	errors := []error{}

	for t := range l.tokens {
		tokens = append(tokens, t)
	}

	for err := range l.err {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("errors: %v", errors)
	}

	return tokens, nil
}

func (l *lexer) scan() {
	lexText(l)
	close(l.tokens)
	close(l.err)
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
			l.err <- err
		}

		return eofRune
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

func (l *lexer) backupAmount(amt int) {
	l.pos -= amt
}

func (l *lexer) peek() rune {
	r := l.next()

	if r != eofRune {
		l.backup()
	}

	return r
}

func (l *lexer) emit(tt tokenType) {
	var value string
	value, l.buf = string(l.bufferSlice(l.bufFrom, l.pos)), l.bufferSlice(l.pos, l.bufTo)
	l.bufFrom = l.pos

	if len(value) > 0 || tt == EOFToken {
		t := Token{tt, value}
		l.tokens <- t
	}
}

func (l *lexer) errorf(format string, args ...any) {
	l.err <- fmt.Errorf(format, args...)
}

func (l *lexer) ignore() {
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

func (l *lexer) acceptSeq(valid string) bool {
	for _, validRune := range valid {
		r := l.next()
		if validRune != r {
			return false
		}
	}

	return true
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
		case eofRune:
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

				if l.acceptAny(" ") { // skip a single whitespace if present after openning brace
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
