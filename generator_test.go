package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestJSGen(t *testing.T) {
	// Mirror of cat.tg (except the line number)
	catDecl := TypeDecl{
		line:     LinePos{0, 0},
		typeName: "Cat",
		typeLine: LinePos{0, 0},
		fields: []Field{
			{
				varName:   "name",
				varLine:   LinePos{0, 0},
				typeName:  "string",
				typeLine:  LinePos{0, 0},
				modifiers: FIELD_CONST,
			},
			{
				varName:   "age",
				varLine:   LinePos{0, 0},
				typeName:  "u32",
				typeLine:  LinePos{0, 0},
				modifiers: FIELD_NONE,
			},
		},
		methods: []FuncDecl{
			{
				line: LinePos{0, 0},
				name: "meow",
				fields: []Field{
					{
						varName:   "sound",
						varLine:   LinePos{0, 0},
						typeName:  "string",
						typeLine:  LinePos{0, 0},
						modifiers: FIELD_NONE,
					},
					{
						varName:   "volume",
						varLine:   LinePos{0, 0},
						typeName:  "u32",
						typeLine:  LinePos{0, 0},
						modifiers: FIELD_NONE,
					},
				},
				returnType: "string",
				returnLine: LinePos{0, 0},
			},
		},
	}

	// Mock in-memory buffer
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	js := JavascriptGenerator{defaultOptions()}
	js.generate([]TypeDecl{catDecl}, writer)

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
