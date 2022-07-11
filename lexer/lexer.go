package lexer

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/aziis98/go-text-ml/utils"
)

type tokenType int

const (
	textToken tokenType = iota
	nodeToken
	braceOpenToken
	braceCloseToken
	eofToken
)

const eofRune rune = 0

type Token struct {
	Type  tokenType
	Value string
}

type lexer struct {
	io.RuneReader

	buf     []rune
	bufFrom int
	bufTo   int

	pos int

	tokens chan Token
	err    error

	escapeDepth *utils.Stack[int]
}

func New(rr io.RuneReader) *lexer {
	l := &lexer{
		RuneReader: rr,

		buf:     []rune{},
		bufFrom: 0,
		bufTo:   0,

		pos: 0,

		tokens: make(chan Token),
		err:    nil,

		escapeDepth: utils.NewStack[int](),
	}

	go l.scan()

	return l
}

func (l *lexer) NextToken() (Token, error) {
	select {
	case t := <-l.tokens:
		return t, nil
	default:
		return Token{}, l.err
	}
}

func (l *lexer) AllTokens() ([]Token, error) {
	tokens := []Token{}

	for t := range l.tokens {
		tokens = append(tokens, t)

		if l.err != nil {
			return nil, l.err
		}
	}

	return tokens, nil
}

func (l *lexer) scan() {
	for state := lexTextBlock; state != nil; {
		state = state(l)

		if l.err != nil {
			break
		}
	}

	close(l.tokens)
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
	if l.err != nil {
		return eofRune
	}

	if l.pos < l.bufTo {
		r := l.bufferAt(l.pos)
		l.pos++
		return r
	}

	r, _, err := l.ReadRune()
	if err != nil {
		if err != io.EOF {
			l.err = err
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
		panic("cannot backtrack before current unfinished token")
	}

	l.pos = pos
}

func (l *lexer) backupAmount(amt int) {
	l.pos -= amt
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

func (l *lexer) emit(tt tokenType) {
	var value string
	value, l.buf = string(l.bufferSlice(l.bufFrom, l.pos)), l.bufferSlice(l.pos, l.bufTo)
	l.bufFrom = l.pos

	l.tokens <- Token{tt, value}
}

func (l *lexer) ignore() {
	l.buf = l.bufferSlice(l.pos, l.bufTo)
	l.bufFrom = l.pos
}

func (l *lexer) acceptAny(valid string) bool {
	r := l.next()

	if strings.ContainsRune(valid, r) {
		l.backup()
		return false
	}

	return true
}

func (l *lexer) acceptSeq(valid string) bool {
	for _, validRune := range valid {
		r := l.next()
		if validRune != r {
			l.err = fmt.Errorf(`expected "%s" but got "%v"`, valid, r)
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

type stateFunc func(*lexer) stateFunc

func lexEOF(l *lexer) stateFunc {
	return nil
}

func lexTextBlock(l *lexer) stateFunc {
	// lexBlockNode

	for {
		r := l.peek()

		if r == '#' {
			nodeStart := l.cursor()

			l.acceptSeq("#")
			l.acceptWhile(func(r rune) bool {
				return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
			})
			nodeEnd := l.cursor()

			l.acceptAnyRepeated(" ")
			spacesEnd := l.cursor()

			newDepth := l.acceptAnyRepeated("{")
			bracesEnd := l.cursor()

			depth := l.escapeDepth.Top()
			if newDepth >= depth {
				l.move(nodeStart) // finish previous text token
				l.emit(textToken)

				l.move(nodeEnd) // emit node token
				l.emit(nodeToken)

				l.move(spacesEnd) // skip whitespace
				l.ignore()

				l.move(bracesEnd) // emit new open brace token
				l.emit(braceOpenToken)

				if newDepth > depth {
					l.escapeDepth.Push(newDepth)
				}

				return lexTextInline
			}
		}

		if l.next() == eofRune {
			return lexEOF
		}
	}
}

func lexTextInline(l *lexer) stateFunc {
	return lexEOF
}
