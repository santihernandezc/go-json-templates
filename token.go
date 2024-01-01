package main

type token struct {
	_type tokenType
	value any
}

func newToken(ttype tokenType) *token {
	return &token{_type: ttype}
}
