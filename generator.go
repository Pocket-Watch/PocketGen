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
	indent               int
	separateDefinitions  bool
	packageName          string
	receiverNameFallback string
}

// These require validation against specific languages
func defaultOptions() GeneratorOptions {
	return GeneratorOptions{
		indent:               4,
		separateDefinitions:  true,
		packageName:          "main",
		receiverNameFallback: "this",
	}
}

type JavascriptGenerator struct {
	options GeneratorOptions
}

type GoGenerator struct {
	options GeneratorOptions
}

type JavaGenerator struct {
	options GeneratorOptions
}

// Writes Javascript definitions based on type declarations
func (js *JavascriptGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	indent := js.options.indent
	joiner := newJoiner()
	for _, t := range types {
		if joiner.join() && js.options.separateDefinitions {
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

		joiner := newJoiner()
		for _, field := range fn.fields {
			if joiner.join() {
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
	joiner := newJoiner()
	for _, field := range fields {
		if joiner.join() {
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

// Writes Go definitions based on type declarations
func (gen *GoGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	writer.WriteString("package ")
	writer.WriteString(gen.options.packageName)
	writer.WriteString("\n\n")
	typeJoiner := newJoiner()
	for _, t := range types {
		if typeJoiner.join() && gen.options.separateDefinitions {
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
		writer.WriteString(gen.toFieldName(field.varName))
		writer.WriteByte(' ')
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteString(gen.toGoType(field.typeName))
		writer.WriteString("\n")
	}
}

func (gen *GoGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	for _, fn := range typeDecl.methods {
		writer.WriteString("func (")
		receiver := gen.toReceiverName(typeDecl.name)
		writer.WriteString(receiver)
		writer.WriteString(" *")
		writer.WriteString(typeDecl.name)
		writer.WriteString(") ")
		writer.WriteString(fn.name)
		if slices.Contains(GO_KEYWORDS, fn.name) {
			// Prevent keyword collision error, append type name to method name
			writer.WriteString(typeDecl.name)
		}

		writer.WriteByte('(')
		joiner := newJoiner()
		for _, field := range fn.fields {
			if joiner.join() {
				writer.WriteString(", ")
			}
			writer.WriteString(gen.toFieldName(field.varName))
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

func (*GoGenerator) toGoType(typeName string) string {
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

func (*JavaGenerator) toJavaType(typeName string) string {
	switch typeName {
	case "i8", "u8":
		return "byte"
	case "i16", "u16":
		return "short"
	case "i32", "u32":
		return "int"
	case "i64", "u64":
		return "long"
	case "f32":
		return "float"
	case "f64":
		return "double"
	case "string":
		return "String"
	case "char":
		return "char"
	case "bool":
		return "boolean"
	default:
		return typeName
	}
}

// Need a strategy to resolve collisions with keywords
func (*JavaGenerator) toFieldName(name string) string {
	if slices.Contains(JAVA_KEYWORDS, name) {
		return name + "0"
	}
	return name
}

// Need a strategy to resolve collisions with keywords
func (*GoGenerator) toFieldName(name string) string {
	if slices.Contains(GO_KEYWORDS, name) {
		return name + "0"
	}
	return name
}

// This
func (gen *GoGenerator) toReceiverName(name string) string {
	receiver := strings.ToLower(string(name[0])) + name[1:]
	if slices.Contains(GO_KEYWORDS, receiver) {
		// Generic receiver to prevent collisions with keywords
		return gen.options.receiverNameFallback
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

// Writes Javascript definitions based on type declarations
func (java *JavaGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	joiner := newJoiner()
	// May require specifying package name
	for _, t := range types {
		if joiner.join() && java.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("class ")
		writer.WriteString(t.name)
		writer.WriteString(" {\n")

		java.writeFields(t.fields, writer)
		writer.WriteString("\n")
		java.writeConstructor(t, writer)
		java.writeMethods(t, writer)
		writer.WriteString("}\n")
	}
	writer.Flush()
}

func (java *JavaGenerator) writeFields(fields []Field, writer *bufio.Writer) {
	indent := java.options.indent
	for _, field := range fields {
		writeIndent(indent, writer)
		if field.hasModifier(FIELD_CONST) {
			writer.WriteString("final ")
		}
		javaType := java.toJavaType(field.typeName)
		writer.WriteString(javaType)
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteByte(' ')
		writer.WriteString(java.toFieldName(field.varName))
		writer.WriteString(";\n")
	}
}

func (java *JavaGenerator) writeConstructor(t TypeDecl, writer *bufio.Writer) {
	indent := java.options.indent
	writeIndent(indent, writer)
	writer.WriteString(t.name)
	writer.WriteString("(")
	join := newJoiner()
	for _, field := range t.fields {
		if join.join() {
			writer.WriteString(", ")
		}
		javaType := java.toJavaType(field.typeName)
		writer.WriteString(javaType)
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteByte(' ')
		writer.WriteString(java.toFieldName(field.varName))
	}
	writer.WriteString(") {\n")
	for _, field := range t.fields {
		writeIndent(2*indent, writer)
		varName := java.toFieldName(field.varName)
		writer.WriteString("this.")
		writer.WriteString(varName)
		writer.WriteString(" = ")
		writer.WriteString(varName)
		writer.WriteByte(';')
		writer.WriteByte('\n')
	}
	writeIndent(indent, writer)
	writer.WriteString("}\n")
}

func (java *JavaGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	indent := java.options.indent
	for _, fn := range typeDecl.methods {
		writeIndent(indent, writer)
		if fn.returnType == "" {
			writer.WriteString("void")
		} else {
			writer.WriteString(java.toJavaType(fn.returnType))
		}
		writer.WriteByte(' ')
		writer.WriteString(fn.name)
		if slices.Contains(JAVA_KEYWORDS, fn.name) {
			// Prevent keyword collision error, append type name to method name
			writer.WriteString(typeDecl.name)
		}

		writer.WriteByte('(')
		fieldJoiner := newJoiner()
		for _, field := range fn.fields {
			if fieldJoiner.join() {
				writer.WriteString(", ")
			}
			javaType := java.toJavaType(field.typeName)
			writer.WriteString(javaType)
			if field.hasModifier(FIELD_ARRAY) {
				writer.WriteString("[]")
			}
			writer.WriteByte(' ')
			writer.WriteString(java.toFieldName(field.varName))
		}
		writer.WriteString(") {\n")
		writeIndent(2*indent, writer)
		writer.WriteString("throw new RuntimeException(\"TODO: Unimplemented method\");\n")
		writeIndent(indent, writer)
		writer.WriteString("}\n")
	}
}

// Joiner abstracts the logic of applying separators
type Joiner struct {
	firstCall bool
}

func newJoiner() Joiner {
	return Joiner{firstCall: true}
}

// Rejects the first call and allows every subsequent one
func (j *Joiner) join() bool {
	if j.firstCall {
		j.firstCall = false
		return false
	}
	return true
}

func (j *Joiner) reset() {
	j.firstCall = true
}

var JAVA_KEYWORDS = []string{
	"abstract", "continue", "for", "new", "switch",
	"assert", "default", "goto", "package", "synchronized",
	"boolean", "do", "if", "private", "this",
	"break", "double", "implements", "protected", "throw",
	"byte", "else", "import", "public", "throws",
	"case", "enum", "instanceof", "return", "transient",
	"catch", "extends", "int", "short", "try",
	"char", "final", "interface", "static", "void",
	"class", "finally", "long", "strictfp", "volatile",
	"const", "float", "native", "super", "while",
}

func writeIndent(indent int, writer *bufio.Writer) {
	for i := 0; i < indent; i++ {
		writer.WriteByte(' ')
	}
}

func openWriter(filename string) *bufio.Writer {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	return bufio.NewWriter(file)
}
