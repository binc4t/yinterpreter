package ast

import "github.com/binc4t/yinterpreter/identify"

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

type LetStatement struct {
	Token *identify.Token
	Left  *IdentExpression
	Right *NormalExpression
}

func (s *LetStatement) TokenRaw() string {
	return s.Token.Raw
}

func (s *LetStatement) statementNode() {

}

type IdentExpression struct {
	Token *identify.Token
}

func (s *IdentExpression) TokenRaw() string {
	return s.Token.Raw
}

func (s *IdentExpression) expressionNode() {
}

type NormalExpression struct {
	Token *identify.Token
}

func (s *NormalExpression) TokenRaw() string {
	return s.Token.Raw
}

func (s *NormalExpression) expressionNode() {
}
