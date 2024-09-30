package identify

import (
	"bytes"
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
	r          io.Reader
	charReader io.Reader
	buf        []byte
	peakBuf    *bytes.Buffer

	size int  // size of the valid buf
	ch   byte // current char
	pos  int  // pos of the next char to read
}

func NewIdentifier(r io.Reader) *Identifier {
	i := &Identifier{
		r:          r,
		charReader: io.LimitReader(r, 1),
		buf:        make([]byte, 1024), // default size
		peakBuf:    new(bytes.Buffer),
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
	i.peakBuf.Reset()
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

func (i *Identifier) PeakChar() (byte, error) {
	if i.pos >= i.size {
		n, err := io.Copy(i.peakBuf, i.charReader)
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, errors.New("PeakChar n is not 1")
		}
	}
	return i.buf[i.pos], nil
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

func (i *Identifier) PeakCharN(n int) (string, error) {
	s := strings.Builder{}
	for j := 0; j < n; j++ {
		c, err := i.PeakChar()
		if err != nil {
			return "", err
		}
		s.WriteByte(c)
	}
	return s.String(), nil
}

func (i *Identifier) NextToken() (*Token, error) {
	if err := i.EatWhiteSpace(); err != nil {
		return nil, err
	}

	var t *Token
	switch {
	case isIdentLetter(i.ch):
		s, err := i.NextItem(isIdentLetter)
		if err != nil {
			return nil, err
		}
		if k := LookupKeywords(s); k != "" {
			return NewTokenString(k, s), nil
		}
		return NewTokenString(IDENT, s), nil
	case isDigit(i.ch):
		s, err := i.NextItem(isDigit)
		if err != nil {
			return nil, err
		}
		return NewTokenString(INT, s), nil
	case i.ch == '!':
		c, err := i.PeakChar()
		if err != nil {
			return NewToken(OPBang, i.ch), nil
		}
		newWord := string(i.ch) + string(c)
		if k := LookupKeywords(newWord); k != "" {
			_ = i.ReadChar()
			t = NewTokenString(k, newWord)
		}
		t = NewToken(OPBang, i.ch)
	case i.ch == '=':
		c, err := i.PeakChar()
		if err != nil {
			return NewToken(OPAssign, i.ch), nil
		}
		newWord := string(i.ch) + string(c)
		if k := LookupKeywords(newWord); k != "" {
			_ = i.ReadChar()
			t = NewTokenString(k, newWord)
		}
		t = NewToken(OPAssign, i.ch)
	default:
		if k := LookupKeywords(string(i.ch)); k != "" {
			t = NewToken(k, i.ch)
		}
		t = NewToken(INVALID, i.ch)
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
