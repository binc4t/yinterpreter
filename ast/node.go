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
	Left  Statement
	Right Statement
}
