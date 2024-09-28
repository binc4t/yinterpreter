package identify

import (
	"errors"
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
	r   io.Reader
	buf []byte

	size int  // size of the valid buf
	ch   byte // current char
	pos  int  // pos of the next char to read
}

func NewIdentifier(r io.Reader) *Identifier {
	i := &Identifier{
		r:   r,
		buf: make([]byte, 1024), // default size
	}
	return i
}

func (i *Identifier) FillIn() error {
	n, err := io.ReadFull(i.r, i.buf)
	if err != nil && errors.Is(err, io.EOF) {
		return err
	}
	i.ch = i.buf[0]
	i.pos = 1
	i.size = n
	return nil
}

func (i *Identifier) ReadChar() error {
	if i.pos >= i.size {
		err := i.FillIn()
		if err != nil {
			return err
		}
	}
	i.ch = i.buf[i.pos]
	i.pos++
	return nil
}

func (i *Identifier) EatWhiteSpace() error {
	_, err := i.NextItem(isWhiteSpace)
	return err
}

func (i *Identifier) NextItem(fn func(byte) bool) (string, error) {
	s := strings.Builder{}
	for fn(i.ch) {
		s.WriteByte(i.ch)
		if err := i.ReadChar(); err != nil {
			return s.String(), err
		}
	}
	return s.String(), nil
}

func (i *Identifier) NextToken() (*Token, error) {
	if err := i.EatWhiteSpace(); err != nil {
		return nil, err
	}

	b := i.ch
	var t *Token
	switch b {
	case '=':
		t = NewToken(OPAssign, b)
	case '+':
		t = NewToken(OPPlus, b)
	case '-':
		t = NewToken(OPMinus, b)
	case '<':
		t = NewToken(OPLT, b)
	case '>':
		t = NewToken(OPGT, b)
	case '(':
		t = NewToken(LParen, b)
	case ')':
		t = NewToken(RParen, b)
	case '{':
		t = NewToken(LBrace, b)
	case '}':
		t = NewToken(RBrace, b)
	case ',':
		t = NewToken(Comma, b)
	case ';':
		t = NewToken(Semicolon, b)
	default:
		if isIdentLetter(i.ch) {
			s, err := i.NextItem(isIdentLetter)
			if err != nil {
				return nil, err
			}
			return NewTokenString(IDENT, s), nil
		}
		if isDigit(i.ch) {
			s, err := i.NextItem(isDigit)
			if err != nil {
				return nil, err
			}
			return NewTokenString(INT, s), nil
		}
	}
	err := i.ReadChar()
	return t, err
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
