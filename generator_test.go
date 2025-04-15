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
			{LinePos{0, 0}, "name", "string", FIELD_CONST},
			{LinePos{0, 0}, "age", "u32", FIELD_NONE},
		},
		methods: []FuncDecl{
			{
				LinePos{0, 0},
				"meow",
				[]Field{
					{LinePos{0, 0}, "sound", "string", FIELD_NONE},
					{LinePos{0, 0}, "volume", "u32", FIELD_NONE},
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
