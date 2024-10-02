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

	curToken  *identify.Token
	peekToken *identify.Token
}

func NewParser(i *identify.Identifier) *Parser {
	p := &Parser{
		i: i,
	}

	// init curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.i.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{
		Statements: make([]Statement, 0),
	}

	for {
		if p.curToken == nil || p.curToken.Type == identify.EOF {
			break
		}

		s := p.parseStatement()
		prog.Statements = append(prog.Statements, s)

	}

	return prog
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case identify.LET:
		return p.parseLetStatement()
	default:
	}

	return nil
}

func (p *Parser) parseLetStatement() *LetStatement {
	s := &LetStatement{
		Token: p.curToken,
	}

	if !p.nextIfPeekTokenIs(identify.IDENT) {
		return nil
	}

	left := p.parseIdentExpression()

	if !p.nextIfPeekTokenIs(identify.OPAssign) {
		return nil
	}

	right := p.parseNormalExpression()
	p.nextToken()

	s.Left = left
	s.Right = right
	return s
}

func (p *Parser) parseIdentExpression() *IdentExpression {
	return &IdentExpression{Token: p.curToken}
}

func (p *Parser) parseNormalExpression() *NormalExpression {
	ret := &NormalExpression{}
	for ; p.curToken.Type != identify.Semicolon; p.nextToken() {
		ret.Tokens = append(ret.Tokens, p.curToken)
	}
	return ret
}

func (p *Parser) peekTokenIs(tokenType string) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) curTokenIs(tokenType string) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) nextIfPeekTokenIs(tokenType string) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	return false
}
