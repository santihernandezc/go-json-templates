package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type scanner struct {
	source []byte
	tokens []*token

	start   int
	current int
}

func newScanner(source []byte) *scanner {
	return &scanner{
		source: source,
	}
}

func (s *scanner) scan() []*token {
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			fmt.Println(err)
		}
	}
	return s.tokens
}

func (s *scanner) scanToken() error {
	r := s.advance()

	switch r {
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case '[':
		s.addToken(LEFT_BRACKET)
	case ']':
		s.addToken(RIGHT_BRACKET)
	case ':':
		s.addToken(COLON)
	case ',':
		s.addToken(COMMA)
	case '"':
		s.string()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number()
	case ' ', '\r', '\t', '\n':
		// Ignore whitespace.
		break
	default:
		if s.isAlpha(r) {
			s.identifier()
			break
		}
		// TODO: cambiar, report error
		return fmt.Errorf("not recognized character %q", string(r))
	}
	return nil
}

func (s *scanner) addToken(ttype tokenType) {
	s.tokens = append(s.tokens, newToken(ttype))
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}

	if s.isAtEnd() {
		// TODO: handle unterminated string.
		fmt.Println("Unterminated string :/")
	}

	// The closing ".
	s.advance()

	tk := newToken(STRING)
	tk.value = string(s.source[s.start+1 : s.current-1])
	s.tokens = append(s.tokens, tk)
}

func (s *scanner) number() error {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	float, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		return err
	}
	tk := newToken(NUMBER)
	tk.value = float
	s.tokens = append(s.tokens, tk)
	return nil
}

func (s *scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	ttype, ok := keywords[string(text)]
	if !ok {
		ttype = IDENTIFIER
	}

	tk := newToken(ttype)

	switch ttype {
	case TRUE:
		tk.value = true
	case FALSE:
		tk.value = false
	default:
		tk.value = string(text)
	}
	s.tokens = append(s.tokens, tk)
}

func (s *scanner) advance() rune {
	r, width := utf8.DecodeRune(s.source[s.current:])
	s.current += width
	return r
}

func (s *scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

func (s *scanner) isAlpha(r rune) bool {
	// TODO: Adding dots, how to handle nested fields?
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '.'
}

func (s *scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		// TODO: ojo null byte
		return 0
	}
	// TODO: consider sing utf8.UTFMax
	r, _ := utf8.DecodeRune(s.source[s.current:])
	return r
}

func (s *scanner) peekNext() rune {
	_, w := utf8.DecodeRune(s.source[s.current:])
	if s.current+w >= len(s.source) {
		return 0
	}
	r, _ := utf8.DecodeRune(s.source[s.current+w:])
	return r
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
