package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// How to handle every Write since they can return errors at any point?

type GeneratorOptions struct {
	indent              int
	separateDefinitions bool
	packageName         string
}

func defaultOptions() GeneratorOptions {
	return GeneratorOptions{
		indent:              4,
		separateDefinitions: true,
		packageName:         "main",
	}
}

type JavascriptGenerator struct {
	options GeneratorOptions
}

// Writes Javascript definitions based on type declarations
func (js *JavascriptGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	indent := js.options.indent
	firstType := true
	for _, t := range types {
		if firstType {
			firstType = false
		} else if js.options.separateDefinitions {
			writer.WriteString("\n")
		}
		// TypeDecl has no information if it's an enum or class
		// Enums usually don't have any modifiers? So we probably need EnumDecl but it must be separate from TypeDecl?
		writer.WriteString("class ")
		writer.WriteString(t.name)
		writer.WriteString(" {\n")

		writeIndent(indent, writer)
		js.writeConstructor(t.fields, writer)
		js.writeMethods(t.methods, writer)
	}
	writer.Flush()
}

func (js *JavascriptGenerator) writeMethods(methods []FuncDecl, writer *bufio.Writer) {
	indent := js.options.indent
	for _, fn := range methods {
		writeIndent(indent, writer)
		writer.WriteString(fn.name)
		writer.WriteByte('(')

		first := true
		for _, field := range fn.fields {
			if first {
				first = false
			} else {
				writer.WriteString(", ")
			}
			writer.WriteString(field.varName)
		}
		// TODO: optionally generate TODO("unimplemented")
		writer.WriteString(") {}\n")
	}
	writer.WriteString("}\n")
}

func (js *JavascriptGenerator) writeConstructor(fields []Field, writer *bufio.Writer) {
	indent := js.options.indent
	writer.WriteString("constructor(")
	first := true
	for _, field := range fields {
		if first {
			first = false
		} else {
			writer.WriteString(", ")
		}
		writer.WriteString(field.varName)
	}
	writer.WriteString(") {\n")
	for _, field := range fields {
		writeIndent(2*indent, writer)
		writer.WriteString("this.")
		writer.WriteString(field.varName)
		writer.WriteString(" = ")
		writer.WriteString(field.varName)
		writer.WriteByte(';')
		writer.WriteByte('\n')
	}
	writeIndent(indent, writer)
	writer.WriteString("}\n")
}

type GoGenerator struct {
	options GeneratorOptions
}

// Writes Go definitions based on type declarations
func (gen *GoGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	writer.WriteString("package ")
	writer.WriteString(gen.options.packageName)
	writer.WriteString("\n\n")
	firstType := true
	for _, t := range types {
		if firstType {
			firstType = false
		} else if gen.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("type ")
		writer.WriteString(t.name)
		writer.WriteString(" struct {\n")

		gen.writeFields(t.fields, writer)
		writer.WriteString("}\n")
		gen.writeMethods(t, writer)
	}
	writer.Flush()
}

func (gen *GoGenerator) writeFields(fields []Field, writer *bufio.Writer) {
	indent := gen.options.indent
	for _, field := range fields {
		writeIndent(indent, writer)
		writer.WriteString(field.varName)
		writer.WriteByte(' ')
		goType := gen.toGoType(field.typeName)
		writer.WriteString(goType)
		writer.WriteString("\n")
	}
}

func (gen *GoGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	for _, fn := range typeDecl.methods {
		writer.WriteString("func (")
		receiver := toReceiverName(typeDecl.name)
		writer.WriteString(receiver)
		writer.WriteString(" *")
		writer.WriteString(typeDecl.name)
		writer.WriteString(") ")
		writer.WriteString(fn.name)
		writer.WriteByte('(')
		first := true
		for _, field := range fn.fields {
			if first {
				first = false
			} else {
				writer.WriteString(", ")
			}
			writer.WriteString(field.varName)
			writer.WriteByte(' ')
			goType := gen.toGoType(field.typeName)
			writer.WriteString(goType)
		}
		writer.WriteString(") ")
		if fn.returnType != "" {
			goType := gen.toGoType(fn.returnType)
			writer.WriteString(goType)
			writer.WriteByte(' ')
		}
		writer.WriteString("{\n")
		writeIndent(gen.options.indent, writer)
		writer.WriteString("panic(\"TODO: Unimplemented method\")\n")
		writer.WriteString("}\n")
	}
	writer.WriteString("\n")
}

func (gen *GoGenerator) toGoType(typeName string) string {
	switch typeName {
	case "i8":
		return "int8"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "u8":
		return "uint8"
	case "u16":
		return "uint16"
	case "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "f32":
		return "float32"
	case "f64":
		return "float64"
	case "string":
		return "string"
	case "char":
		return "rune"
	case "bool":
		return "bool"
	default:
		return typeName
	}
}

// This
func toReceiverName(name string) string {
	receiver := strings.ToLower(string(name[0])) + name[1:]
	if slices.Contains(GO_KEYWORDS, receiver) {
		// Generic receiver to prevent collisions with keywords
		return "this"
	}
	return receiver
}

// Constants
var GO_KEYWORDS = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

func writeIndent(indent int, writer *bufio.Writer) {
	for i := 0; i < indent; i++ {
		writer.WriteByte(' ')
	}
}

func openWriter(filename string) *bufio.Writer {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("File already exists:", filename)
		} else {
			fmt.Println("Error opening file:", err)
		}
		return nil
	}

	writer := bufio.NewWriter(file)
	return writer
}
