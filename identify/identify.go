package identify

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token struct {
	Type string
	Raw  string
}

func NewToken(tokenType string, raw byte) *Token {
	return &Token{Type: tokenType, Raw: string(raw)}
}

func NewTokenString(tokenType string, rawString string) *Token {
	return &Token{Type: tokenType, Raw: rawString}
}

type Identifier struct {
	r       io.Reader
	scanner *bufio.Scanner
	buf     *bytes.Buffer

	size int  // size of the valid buf
	ch   byte // current char
	pos  int  // pos of the next char to read
}

func NewIdentifier(r io.Reader) *Identifier {
	raw := make([]byte, 0, 1024)
	i := &Identifier{
		r:       r,
		scanner: bufio.NewScanner(r),
		buf:     bytes.NewBuffer(raw),
	}
	return i
}

func (i *Identifier) FillIn() bool {
	if !i.scanner.Scan() {
		return false
	}
	i.buf.Reset()
	i.buf.WriteString(i.scanner.Text())
	if i.buf.Len() != 0 {
		i.ch = i.buf.Bytes()[0]
	} else {
		i.ch = 0
	}
	i.pos = 1
	i.size = i.buf.Len()
	return true
}

func (i *Identifier) ReadChar() {
	if i.pos >= i.size {
		i.ch = 0
	} else {
		i.ch = i.buf.Bytes()[i.pos]
	}
	i.pos++
}

func (i *Identifier) PeekChar() byte {
	if i.pos >= i.size {
		return 0
	}
	return i.buf.Bytes()[i.pos]
}

func (i *Identifier) EatWhiteSpace() {
	_ = i.NextItem(isWhiteSpace)
}

func (i *Identifier) NextItem(fn func(byte) bool) string {
	s := strings.Builder{}
	for fn(i.ch) {
		s.WriteByte(i.ch)
		i.ReadChar()
	}
	return s.String()
}

func (i *Identifier) PeekCharN(n int) string {
	s := strings.Builder{}
	for j := 0; j < n; j++ {
		c := i.PeekChar()
		s.WriteByte(c)
	}
	return s.String()
}

func (i *Identifier) NextToken() *Token {
	i.EatWhiteSpace()

	var t *Token
	switch {
	case isIdentLetter(i.ch):
		s := i.NextItem(isIdentLetter)
		if k := LookupKeywords(s); k != "" {
			return NewTokenString(k, s)
		}
		return NewTokenString(IDENT, s)
	case isDigit(i.ch):
		s := i.NextItem(isDigit)
		return NewTokenString(INT, s)
	case i.ch == '!':
		c := i.PeekChar()
		newWord := string(i.ch) + string(c)
		if k := LookupKeywords(newWord); k != "" {
			i.ReadChar()
			t = NewTokenString(k, newWord)
		} else {
			t = NewToken(OPBang, i.ch)
		}
	case i.ch == '=':
		c := i.PeekChar()
		newWord := string(i.ch) + string(c)
		if k := LookupKeywords(newWord); k != "" {
			i.ReadChar()
			t = NewTokenString(k, newWord)
		} else {
			t = NewToken(OPAssign, i.ch)
		}
	case i.ch == 0:
		return NewToken(EOF, 0)
	default:
		if k := LookupKeywords(string(i.ch)); k != "" {
			t = NewToken(k, i.ch)
		} else {
			t = NewToken(INVALID, i.ch)
		}
	}

	i.ReadChar()
	return t
}

func isIdentLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b == '_')
}

func isDigit(b byte) bool {
	return b > '0' && b < '9'
}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n'
}
