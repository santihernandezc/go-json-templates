package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	testValues := map[string]any{
		"test_string": "test",
		"test_int":    3,
		"test_float":  3.14,
		"test_obj": map[string]any{
			"string": "test",
			"test":   true,
		},
	}

	tests := []struct {
		name string
		data string
		exp  string
	}{
		{
			"empty",
			"",
			"{}",
		},
		{
			"braces",
			"{}",
			"{}",
		},
		{
			"string, float, and integer values",
			`{
				"test_string": "test",
				"test_int": 12345,
				"test_float": 123.45
			}`,
			`{
				"test_string": "test",
				"test_int": 12345,
				"test_float": 123.45
			}`,
		},
		{
			"array",
			`{
				"test": [1, 2, 3, 4, 5]
			}`,
			`{
				"test": [1, 2, 3, 4, 5]
			}`,
		},
		{
			"array of integers, floats, and strings",
			`{
				"test": [1, 2.5, "test", "4", "5.7"]
			}`,
			`{
				"test": [1, 2.5, "test", "4", "5.7"]
			}`,
		},
		{
			"variables in array",
			`{
				"test": [test_string, test_int, test_float, test_obj]
			}`,
			`{
				"test": ["test", 3, 3.14, {"string": "test", "test": true}]
			}`,
		},
		{
			"objects in array",
			`{
				"test": [{"test": true}, {"test": "test"}]
			}`,
			`{
				"test": [{"test": true}, {"test": "test"}]
			}`,
		},
		{
			"arrays in array",
			`{
				"test": ["test", true, 1234, 1.234, false, [test_string]]
			}`,
			`{
				"test": ["test", true, 1234, 1.234, false, ["test"]]
			}`,
		},
		{
			"variables in field values",
			`{
				"test_string": test_string,
				"test_int": test_int,
				"test_float": test_float
			}`,
			`{
				"test_string": "test",
				"test_int": 3,
				"test_float": 3.14
			}`,
		},
		{
			"variables in field names",
			`{
				test_string: "test"
			}`,
			`{
				"test": "test"
			}`,
		},
		{
			"variables in field names and values",
			`{
				test_string: test_string
			}`,
			`{
				"test": "test"
			}`,
		},
		{
			"booleans in field values",
			`{
				"true": true,
				"false": false
			}`,
			`{
				"true": true,
				"false": false
			}`,
		},
		{
			"external object in field values",
			`{
				"object": test_obj
			}`,
			`{
				"object": {
					"test": true,
					"string": "test"
				}
			}`,
		},
		{
			"nested fields in field value and field name",
			`{
				test_obj.string: test_obj.test
			}`,
			`{
				"test": true
			}`,
		},
		{
			"object in field value",
			`{
				"test_object": {
					"test": true
				}
			}`,
			`{
				"test_object": {
					"test": true
				}
			}`,
		},
		{
			"nested objects in field value",
			`{
				"test_object": {
					"test": {
						"test": {
							"test": true
						}
					}
				}
			}`,
			`{
				"test_object": {
					"test": {
						"test": {
							"test": true
						}
					}
				}
			}`,
		},
		{
			"nested objects and arrays in field value",
			`{
				"test_object": {
					"test": [
						{
							"test": {
								"test": true
							}
						},
						{
							"test": [1, 2, 3, 4]
						}
					]
				}
			}`,
			`{
				"test_object": {
					"test": [
						{
							"test": {
								"test": true
							}
						},
						{
							"test": [1, 2, 3, 4]
						}
					]
				}
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			scanner := newScanner([]byte(test.data))
			tokens := scanner.scan()

			parser := newParser(tokens)
			statements, err := parser.parse()
			if err != nil {
				tt.Fatalf("parser.parse() err = %v, want nil", err)
			}

			i := newInterpreter(statements, testValues)
			res, err := i.interpret()
			if err != nil {
				tt.Fatalf("interpreter.interpret() err = %v, want nil", err)
			}

			require.JSONEq(tt, test.exp, string(res))
		})
	}
}

func BenchmarkInterpreter(b *testing.B) {
	b.ReportAllocs()
	testValues := map[string]any{
		"test_value": map[string]any{
			"test": map[string]any{
				"testcito": map[string]any{
					"test": true,
				},
			},
		},
	}

	testTemplate := `
{
	"string": "test",
	"int": 12345,
	"float": 123.45,
	"true": true,
	"false": false,
	"array": [1, 2, 3, 4, 5],
	"nested_obj": {
		"test": {
			"test": {
				"test": true
			}
		}
	},
	"obj": test_value,
	"nested_field": test_value.test,
	"nested_field_2": test_value.test.testcito
}
`

	for i := 0; i < b.N; i++ {
		statements, _ := newParser(newScanner([]byte(testTemplate)).scan()).parse()
		if _, err := newInterpreter(statements, testValues).interpret(); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonnet(b *testing.B) {
	b.ReportAllocs()
	testValue := map[string]any{
		"test": map[string]any{
			"testcito": map[string]any{
				"test": true,
			},
		},
	}

	testTemplate := `
		function(test_value) {
		string: "test",
		int: 12345,
		float: 123.45,
		"true": true,
		"false": false,
		array: [1, 2, 3, 4, 5],
		nested_obj: {
			"test": {
				"test": {
					"test": true
				}
			}
		},
		obj: test_value,
		nested_field: test_value["test"],
		nested_field_2: test_value["test"]["testcito"]
	}`
	bytes, _ := json.Marshal(testValue)

	vm := jsonnet.MakeVM()
	vm.TLACode("test_value", string(bytes))

	for i := 0; i < b.N; i++ {
		if _, err := vm.EvaluateAnonymousSnippet("test.jsonnet", testTemplate); err != nil {
			panic(err)
		}
	}
}
