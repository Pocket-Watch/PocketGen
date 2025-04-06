package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestJSGen(t *testing.T) {
	// Mirror of cat.tg
	catDecl := TypeDecl{
		name: "Cat",
		fields: []Field{
			{"name", "string", []FieldModifier{FIELD_CONST}},
			{"age", "u32", []FieldModifier{}},
		},
		methods: []FuncDecl{
			{
				"meow",
				[]Field{
					{"sound", "string", []FieldModifier{}},
					{"volume", "u32", []FieldModifier{}},
				},
				"string",
			},
		},
	}

	// Mock in-memory buffer
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	generateJS([]TypeDecl{catDecl}, writer)

	t.Log("\n" + buffer.String())
}
