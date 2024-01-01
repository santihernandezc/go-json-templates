package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type interpreter struct {
	statements []any
	values     map[string]any
}

func newInterpreter(statements []any, values map[string]any) *interpreter {
	return &interpreter{
		statements: statements,
		values:     values,
	}
}

func (i *interpreter) interpret() ([]byte, error) {
	data := make(map[string]any)
	for _, stmt := range i.statements {
		switch v := stmt.(type) {
		case *fieldDeclaration:
			name, value, err := i.fieldDeclaration(v)
			if err != nil {
				return nil, err
			}
			data[name] = value
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// fieldDeclaration returns the value for a field and an error.
// TODO: change name.
func (i *interpreter) fieldDeclaration(declaration *fieldDeclaration) (string, any, error) {
	var name string
	switch n := declaration.name.(type) {
	case *literal:
		str, ok := n.value.(string)
		if !ok {
			return "", nil, fmt.Errorf("expected string as field name")
		}
		name = str

	case *identifier:
		keys := strings.Split(n.value, ".")
		resolved, err := resolve(i.values, keys...)
		if err != nil {
			return "", nil, err
		}

		// Field name must be a string.
		str, ok := resolved.(string)
		if !ok {
			return "", nil, fmt.Errorf("expected string as field name")
		}
		name = str
	default:
		return "", nil, fmt.Errorf("invalid name for field")
	}

	switch v := declaration.value.(type) {
	case *literal:
		return name, v.value, nil

	case *identifier:
		keys := strings.Split(v.value, ".")
		resolved, err := resolve(i.values, keys...)
		if err != nil {
			return "", nil, err
		}
		return name, resolved, nil

	case *array:
		values, err := i.array(v.values)
		if err != nil {
			return "", nil, err
		}
		return name, values, nil

		// TODO: block
	case *block:
		object, err := i.object(v.statements)
		if err != nil {
			return "", nil, err
		}
		return name, object, nil

	default:
		return "", nil, fmt.Errorf("unexpected value for field")
	}
}

// resolve returns the value for a key from the values map.
func resolve(values map[string]any, keys ...string) (any, error) {
	var resolved any
	var ok bool
	for _, k := range keys {
		resolved, ok = values[k]
		if !ok {
			return nil, fmt.Errorf("undefined value for %s", k)
		}

		// If it's a map, we can take values from there.
		if v, ok := resolved.(map[string]any); ok {
			values = v
		}
	}
	return resolved, nil
}

func (i *interpreter) array(elements []any) ([]any, error) {
	var finalValues []any
	for _, e := range elements {
		switch v := e.(type) {
		case *literal:
			finalValues = append(finalValues, v.value)

		case *identifier:
			resolved, err := resolve(i.values, v.value)
			if err != nil {
				return nil, err
			}
			finalValues = append(finalValues, resolved)

		case *array:
			arr, err := i.array(v.values)
			if err != nil {
				return nil, err
			}
			finalValues = append(finalValues, arr)

		case *block:
			obj, err := i.object(v.statements)
			if err != nil {
				return nil, err
			}
			finalValues = append(finalValues, obj)
			// TODO: add objects.
		}
	}
	return finalValues, nil
}

func (i *interpreter) object(statements []any) (map[string]any, error) {
	obj := make(map[string]any)
	for _, stmt := range statements {
		switch v := stmt.(type) {
		case *fieldDeclaration:
			name, value, err := i.fieldDeclaration(v)
			if err != nil {
				return nil, err
			}
			obj[name] = value
		}
	}
	return obj, nil
}
