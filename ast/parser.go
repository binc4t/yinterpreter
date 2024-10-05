package ast

import (
	"fmt"
	"github.com/binc4t/yinterpreter/identify"
	"github.com/binc4t/yinterpreter/libs"
	"strings"
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

func (p *Program) String() string {
	b := strings.Builder{}
	for _, s := range p.Statements {
		b.WriteString(s.String())
	}
	return b.String()
}

type Parser struct {
	i *identify.Identifier

	curToken  *identify.Token
	peekToken *identify.Token
	errors    []error
}

func NewParser(i *identify.Identifier) *Parser {
	p := &Parser{
		i:      i,
		errors: make([]error, 0),
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
		if p.curToken.Type == identify.EOF {
			break
		}

		s := p.parseStatement()
		if !libs.IsNil(s) {
			prog.Statements = append(prog.Statements, s)
		}
		p.nextToken()
	}

	return prog
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case identify.LET:
		return p.parseLetStatement()
	case identify.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
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
	p.nextToken()

	right := p.parseNormalExpression()
	p.nextToken()

	s.Left = left
	s.Right = right
	return s
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	s := &ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()

	if exp := p.parseNormalExpression(); exp != nil {
		s.Exp = exp
	}

	return s
}

func (p *Parser) parseIdentExpression() *IdentExpression {
	return &IdentExpression{Token: p.curToken}
}

func (p *Parser) parseNormalExpression() *NormalExpression {
	ret := &NormalExpression{}
	for {
		if p.curToken.Type == identify.EOF || p.curToken.Type == identify.Semicolon {
			break
		}
		ret.Tokens = append(ret.Tokens, p.curToken)
		p.nextToken()
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
	p.peekError(tokenType)
	return false
}

func (p *Parser) peekError(tokenType string) {
	p.errors = append(p.errors, fmt.Errorf("expect peek to be %s, but got %s", tokenType, p.peekToken.Type))
}
