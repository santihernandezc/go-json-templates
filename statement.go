package main

type fieldDeclaration struct {
	name  any
	value any
}

func newFieldDeclaration(name any, value any) *fieldDeclaration {
	return &fieldDeclaration{
		name:  name,
		value: value,
	}
}
