package ast

import (
	"fmt"
	"github.com/binc4t/yinterpreter/identify"
	"strings"
)

type Node interface {
	TokenRaw() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type baseStatement struct {
}

func (b *baseStatement) statementNode() {
}

type baseExpression struct {
}

func (b *baseExpression) expressionNode() {
}

type LetStatement struct {
	baseStatement
	Token *identify.Token
	Left  *IdentExpression
	Right *NormalExpression
}

func (s *LetStatement) TokenRaw() string {
	return s.Token.Raw
}

type ReturnStatement struct {
	baseStatement
	Token *identify.Token
	Exp   Expression
}

func (s *ReturnStatement) TokenRaw() string {
	return s.Token.Raw
}

type IdentExpression struct {
	baseExpression
	Token *identify.Token
}

func (s *IdentExpression) TokenRaw() string {
	return s.Token.Raw
}

type NormalExpression struct {
	baseExpression
	Tokens []*identify.Token
}

func (s *NormalExpression) TokenRaw() string {
	ret := strings.Builder{}
	for _, t := range s.Tokens {
		ret.WriteString(fmt.Sprintf("%v, ", t))
	}
	return ret.String()
}
