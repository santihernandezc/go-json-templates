package main

type identifier struct {
	value string
}

func newIdentifier(value string) *identifier {
	return &identifier{value: value}
}

type literal struct {
	value any
}

func newLiteral(value any) *literal {
	return &literal{value: value}
}

type array struct {
	values []any
}

func newArray(values []any) *array {
	return &array{values: values}
}

type block struct {
	statements []any
}

func newBlock(statements []any) *block {
	return &block{statements: statements}
}
