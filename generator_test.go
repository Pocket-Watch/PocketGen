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
			{"name", "string", FIELD_CONST},
			{"age", "u32", FIELD_NONE},
		},
		methods: []FuncDecl{
			{
				"meow",
				[]Field{
					{"sound", "string", FIELD_NONE},
					{"volume", "u32", FIELD_NONE},
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
