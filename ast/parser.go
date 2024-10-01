package ast

import (
	"github.com/binc4t/yinterpreter/identify"
)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenRaw() string {
	if len(p.Statements) != 0 {
		return p.Statements[0].TokenRaw()
	}
	return ""
}

type Parser struct {
	i *identify.Identifier

	curNode  *identify.Token
	peakNode *identify.Token
}

func NewParser(i *identify.Identifier) *Parser {
	p := &Parser{
		i: i,
	}
	// read
	return p
}

func (p *Parser) nextNode() {
	p.curNode = p.peakNode
	p.peakNode, _ = p.i.NextToken()
}

func (p *Parser) ParseProgram() *Program {

}
