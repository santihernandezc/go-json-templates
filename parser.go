package main

import (
	"fmt"
)

type parser struct {
	tokens  []*token
	current int
}

func newParser(tokens []*token) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) parse() ([]any, error) {
	var statements []any
	for !p.isAtEnd() {
		if p.match(LEFT_BRACE) {
			stmts, err := p.block()
			if err != nil {
				return nil, fmt.Errorf("error parsing object: %w", err)
			}
			statements = append(statements, stmts...)
			continue
		}
		return nil, fmt.Errorf("nothing")
	}
	return statements, nil
}

// Assuming I'm in an object declaration.
// TODO: change for interface?
func (p *parser) block() ([]any, error) {
	var statements []any
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		// TODO: add statements
		stmt, err := p.fieldDeclaration()
		if err != nil {
			fmt.Println("Hubo err!", err)
			return nil, err
		}

		// If that was the last statement, there should be no comma.
		if !p.check(RIGHT_BRACE) {
			if _, err := p.consume(COMMA); err != nil {
				return nil, fmt.Errorf("expected comma after field declaration")
			}
		}
		statements = append(statements, stmt)
	}

	if _, err := p.consume(RIGHT_BRACE); err != nil {
		return nil, err
	}

	return statements, nil
}

// check returns whether the current tokes is of certain type.
func (p *parser) check(tTypes ...tokenType) bool {
	if p.isAtEnd() {
		return false
	}

	for _, t := range tTypes {
		if p.peek()._type == t {
			return true
		}
	}
	return false
}

func (p *parser) peek() *token {
	return p.tokens[p.current]
}

func (p *parser) advance() *token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *parser) isAtEnd() bool {
	// TODO: add EOF token?
	return p.current >= len(p.tokens)
}

func (p *parser) consume(tType tokenType) (*token, error) {
	if p.check(tType) {
		return p.advance(), nil
	}

	return nil, fmt.Errorf("token type not found")
}

// match checks for a type and returns 'true' if it matches, discarding the token.
func (p *parser) match(ttypes ...tokenType) bool {
	if p.check(ttypes...) {
		p.advance()
		return true
	}
	return false
}

func (p *parser) fieldDeclaration() (*fieldDeclaration, error) {
	var left any
	switch {
	case p.check(STRING):
		left = newLiteral(p.advance().value)

	case p.check(IDENTIFIER):
		identifier, err := p.identifier()
		if err != nil {
			return nil, err
		}
		left = identifier

	default:
		// TODO: better errors.
		return nil, fmt.Errorf("expected string or identifier as field name")
	}

	if _, err := p.consume(COLON); err != nil {
		return nil, fmt.Errorf("expected ':' after field name")
	}

	var rightValue any
	switch {
	case p.match(LEFT_BRACE):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		rightValue = newBlock(stmts)

	case p.match(LEFT_BRACKET):
		array, err := p.array()
		if err != nil {
			return nil, err
		}
		// TODO: could be an array of arrays/objects.
		rightValue = newArray(array)

	case p.check(STRING, NUMBER, TRUE, FALSE):
		rightValue = newLiteral(p.advance().value)

	case p.check(IDENTIFIER):
		identifier, err := p.identifier()
		if err != nil {
			return nil, err
		}
		rightValue = identifier

	default:
		return nil, fmt.Errorf("expected string, number, boolean, array or object as field value")
	}

	return newFieldDeclaration(left, rightValue), nil
}

func (p *parser) array() ([]any, error) {
	var array []any
	for !p.isAtEnd() && !p.check(RIGHT_BRACKET) {
		switch {
		case p.check(STRING, NUMBER, TRUE, FALSE):
			array = append(array, newLiteral(p.advance().value))

		case p.check(IDENTIFIER):
			identifier, err := p.identifier()
			if err != nil {
				return nil, err
			}
			array = append(array, identifier)

		case p.match(LEFT_BRACE):
			stmts, err := p.block()
			if err != nil {
				return nil, err
			}
			array = append(array, newBlock(stmts))

		case p.match(LEFT_BRACKET):
			a, err := p.array()
			if err != nil {
				return nil, err
			}
			array = append(array, newArray(a))
		}

		// If that was the last value, there should be no comma.
		if !p.check(RIGHT_BRACKET) {
			if _, err := p.consume(COMMA); err != nil {
				return nil, fmt.Errorf("expected comma after array element")
			}
		}
	}

	if _, err := p.consume(RIGHT_BRACKET); err != nil {
		return nil, fmt.Errorf("expected ']' at the end of array")
	}
	return array, nil
}

func (p *parser) identifier() (*identifier, error) {
	i, err := p.consume(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	name, ok := i.value.(string)
	if !ok {
		return nil, fmt.Errorf("identifier must be a string")
	}

	return newIdentifier(name), nil
}
