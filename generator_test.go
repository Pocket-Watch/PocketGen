package main

import (
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

	buffer := bytes.Buffer{}

	js := JavascriptGenerator{defaultOptions()}
	js.generate([]TypeDecl{catDecl}, &buffer)

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

func TestBasicToSnakeCase(t *testing.T) {
	input := "theOldestBook"
	expected := "the_oldest_book"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Expected:", expected, "  Got:", actual)
	}
}

func TestUnicodeToSnakeCase(t *testing.T) {
	input := "onTheЯight"
	expected := "on_the_яight"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Expected:", expected, "  Got:", actual)
	}
}

func TestGoStyleVarToSnakeCase(t *testing.T) {
	input := "GoVariable"
	expected := "go_variable"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Expected:", expected, "  Got:", actual)
	}
}

func TestOverCapitalizationSnakeCase(t *testing.T) {
	input := "v_1_V_3"
	expected := "v_1_v_3"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestNoTrailingSnakeCase(t *testing.T) {
	input := "whatS"
	expected := "what_s"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestTrailingSnakeCase(t *testing.T) {
	input := "AbC_"
	expected := "ab_c_"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestUnicodeSnakeCase(t *testing.T) {
	input := "umlautÜ"
	expected := "umlaut_ü"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestManyUnderscores(t *testing.T) {
	input := "_1_2__3___Z"
	expected := "_1_2__3___z"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestTwoUppercasePascalLetters(t *testing.T) {
	input := "takeABreak"
	expected := "take_a_break"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestConstCase(t *testing.T) {
	input := "simpleXMLFile"
	expected := "simple_xml_file"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestOnlyTwoUppers(t *testing.T) {
	input := "TU"
	expected := "tu"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestEndInTwoUppers(t *testing.T) {
	input := "popTG"
	expected := "pop_tg"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestEndInTwoUppersButUnderscored(t *testing.T) {
	input := "pop_TG"
	expected := "pop_tg"
	actual := toSnakeCase(input)
	if actual != expected {
		t.Error("Input: ", input, "  Expected:", expected, "  Got:", actual)
	}
}

func TestBasicCapitalizeFirstLetter(t *testing.T) {
	input := "cake"
	expected := "Cake"
	actual := capitalizeFirstLetter(input)
	if actual != expected {
		t.Error("Expected:", expected, "  Got:", actual)
	}
}

func TestUnicodeCapitalizeFirstLetter(t *testing.T) {
	input := "тrait"
	expected := "Тrait"
	actual := capitalizeFirstLetter(input)
	if actual != expected {
		t.Error("Expected:", expected, "  Got:", actual)
	}
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
