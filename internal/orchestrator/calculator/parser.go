package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	tokens []string
	pos    int
}

func NewParser(expression string) *Parser {
	tokens := tokenize(expression)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func tokenize(expression string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range expression {
		if unicode.IsSpace(char) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else if isOperator(string(char)) || char == '(' || char == ')' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func (p *Parser) Parse() (ASTNode, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	return p.parseExpression()
}

func (p *Parser) parseExpression() (ASTNode, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "+" || p.tokens[p.pos] == "-") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseTerm() (ASTNode, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "*" || p.tokens[p.pos] == "/") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseFactor() (ASTNode, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	token := p.tokens[p.pos]
	p.pos++

	if token == "(" {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		
		p.pos++ // Consume ")"
		return expr, nil
	}

	num, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", token)
	}

	return &NumberNode{Value: num}, nil
}
