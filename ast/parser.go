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

func (p *Parser) peekTokenIs(tokenType string) bool {
	if p.peekToken != nil && p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	}
	return false
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
	if !p.peekTokenIs(identify.IDENT) {
		return nil
	}

	s := &LetStatement{
		Token: p.curToken,
	}

	left := p.parseIdentExpression()

	if !p.peekTokenIs(identify.OPAssign) {
		return nil
	}

	right := p.parseNormalExpression()

	s.Left = left
	s.Right = right
	return s
}

func (p *Parser) parseIdentExpression() *IdentExpression {
	return &IdentExpression{Token: p.curToken}
}

func (p *Parser) parseNormalExpression() *NormalExpression {
	for ; p.curToken.Type != identify.Semicolon; p.nextToken() {
	}
	return &NormalExpression{Token: p.curToken}
}
