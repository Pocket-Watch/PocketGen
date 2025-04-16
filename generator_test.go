package main

import (
	"bufio"
	"bytes"
	"strings"
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

	output := buffer.String()
	t.Log("\n" + output)
	lines := strings.Split(output, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	expectedLines := []string{
		"class Cat {",
		"constructor(name, age) {",
		"this.name = name;",
		"this.age = age;",
		"}",
		"meow(sound, volume) {}",
		"}",
		"",
	}

	compareLines(expectedLines, lines, t)
}

func compareLines(expectedLines []string, lines []string, t *testing.T) {
	if len(lines) != len(expectedLines) {
		t.Errorf("Different number of lines: actual = %d, expected = %d", len(lines), len(expectedLines))
		return
	}
	for i := 0; i < len(lines); i++ {
		actual := lines[i]
		expected := expectedLines[i]
		if actual != expected {
			t.Errorf("Lines differ at index %v: actual = %s, expected = %s", i, actual, expected)
			return
		}
	}
}
