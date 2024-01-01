package main

type tokenType int

const (
	LEFT_BRACE tokenType = iota
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	COMMA
	COLON
	IDENTIFIER
	STRING
	NUMBER
	TRUE
	FALSE
)

var keywords = map[string]tokenType{
	"true":  TRUE,
	"false": FALSE,
}
