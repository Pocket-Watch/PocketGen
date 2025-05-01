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

type RustGenerator struct {
	options GeneratorOptions
}

func keywordCollisionError(declType string, keyword string, language string, filepath string, pos LinePos) error {
	return fmt.Errorf("ERROR @ %s:%v:%v %v name '%v' is a %v keyword.\n",
		filepath, pos.number, pos.offset, declType, keyword, language)
}

// Checks if keywords collide with type, field or parameter names per given keyword set
func checkKeywords(types []TypeDecl, keywords []string, language string, filepath string) error {
	for _, t := range types {
		if slices.Contains(keywords, t.typeName) {
			return keywordCollisionError("type", t.typeName, language, filepath, t.typeLine)
		}

		// ATP there's no need to check the typename
		for _, field := range t.fields {
			if slices.Contains(keywords, field.varName) {
				return keywordCollisionError("field", field.varName, language, filepath, field.varLine)
			}
		}

		for _, fn := range t.methods {
			if slices.Contains(keywords, fn.name) {
				return keywordCollisionError("method", fn.name, language, filepath, fn.line)
			}

			for _, f := range fn.fields {
				if slices.Contains(keywords, f.varName) {
					return keywordCollisionError("parameter", f.varName, language, filepath, f.varLine)
				}
			}
		}
	}
	return nil
}

// Translate types to a language specific representation using the supplied convert() function.
// Each generator should provide a mapping function.
func translateTypes(types []TypeDecl, convert func(s string) string) {
	for _, t := range types {
		for i := range t.fields {
			field := &t.fields[i]
			field.typeName = convert(field.typeName)
		}
		for i := range t.methods {
			method := &t.methods[i]
			method.returnType = convert(method.returnType)
			for j := range method.fields {
				field := &method.fields[j]
				field.typeName = convert(field.typeName)
			}
		}
	}
}

// Writes Javascript definitions based on type declarations
func (js *JavascriptGenerator) generate(types []TypeDecl, writer *bufio.Writer) {
	indent := js.options.indent
	joiner := newJoiner()
	for _, t := range types {
		if joiner.join() && js.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("class " + t.typeName + " {\n")

		writeIndent(indent, writer)
		js.writeConstructor(t.fields, writer)
		js.writeMethods(t.methods, writer)
	}
	writer.Flush()
}

// Writes Go definitions based on type declarations
func (gen *GoGenerator) generate(types []TypeDecl, writer *bufio.Writer, filepath string) error {
	err := checkKeywords(types, GO_KEYWORDS, "go", filepath)
	if err != nil {
		return err
	}
	translateTypes(types, toGoType)

	writer.WriteString("package " + gen.options.packageName + "\n\n")

	typeJoiner := newJoiner()
	for _, t := range types {
		if typeJoiner.join() && gen.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("type " + t.typeName + " struct {\n")

		gen.writeFields(t.fields, writer)
		writer.WriteString("}\n")
		gen.writeMethods(t, writer)
	}
	writer.Flush()
	return nil
}

// Writes Java definitions based on type declarations
func (java *JavaGenerator) generate(types []TypeDecl, writer *bufio.Writer, filepath string) error {
	err := checkKeywords(types, JAVA_KEYWORDS, "java", filepath)
	if err != nil {
		return err
	}
	translateTypes(types, toJavaType)

	joiner := newJoiner()
	// May require specifying mod name
	for _, t := range types {
		if joiner.join() && java.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("class " + t.typeName + " {\n")

		java.writeFields(t.fields, writer)
		if len(t.fields) > 0 {
			writer.WriteString("\n")
		}
		java.writeConstructor(t, writer)
		java.writeMethods(t, writer)
		writer.WriteString("}\n")
	}
	writer.Flush()
	return nil
}

// Writes Rust definitions based on type declarations
func (rust *RustGenerator) generate(types []TypeDecl, writer *bufio.Writer, filepath string) error {
	err := checkKeywords(types, RUST_KEYWORDS, "rust", filepath)
	if err != nil {
		return err
	}
	translateTypes(types, toRustType)

	joiner := newJoiner()
	// May require specifying package name
	for _, t := range types {
		if joiner.join() && rust.options.separateDefinitions {
			writer.WriteString("\n")
		}
		writer.WriteString("struct " + t.typeName + " {\n")
		rust.writeFields(t.fields, writer)
		writer.WriteString("}\n")

		if len(t.methods) > 0 {
			writer.WriteString("impl " + t.typeName + " {\n")
			rust.writeMethods(t, writer)
			writer.WriteString("}\n")
		}
	}
	writer.Flush()
	return nil
}

func (js *JavascriptGenerator) writeMethods(methods []FuncDecl, writer *bufio.Writer) {
	indent := js.options.indent
	for _, fn := range methods {
		writeIndent(indent, writer)
		writer.WriteString(fn.name + "(")

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
		assignment := "this." + field.varName + " = " + field.varName + ";\n"
		writer.WriteString(assignment)
	}
	writeIndent(indent, writer)
	writer.WriteString("}\n")
}

func (gen *GoGenerator) writeFields(fields []Field, writer *bufio.Writer) {
	indent := gen.options.indent
	for _, field := range fields {
		writeIndent(indent, writer)
		writer.WriteString(field.varName + " ")
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteString(field.typeName + "\n")
	}
}

func (rust *RustGenerator) writeFields(fields []Field, writer *bufio.Writer) {
	indent := rust.options.indent
	for _, field := range fields {
		writeIndent(indent, writer)
		writer.WriteString(field.varName + ": ")
		rust.writeFieldType(field, writer)
		// Trailing comma is probably fine
		writer.WriteString(",\n")
	}
}

func (rust *RustGenerator) writeFieldType(field Field, writer *bufio.Writer) {
	isList := field.hasModifier(FIELD_ARRAY)
	if isList {
		writer.WriteString("Vec<")
	}
	writer.WriteString(field.typeName)
	if isList {
		writer.WriteString(">")
	}
}

func (gen *GoGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	for _, fn := range typeDecl.methods {
		receiver := gen.toReceiverName(typeDecl.typeName)

		funcHeader := "func (" + receiver + " *" + typeDecl.typeName + ") " + fn.name + "("
		writer.WriteString(funcHeader)

		joiner := newJoiner()
		for _, field := range fn.fields {
			if joiner.join() {
				writer.WriteString(", ")
			}
			writer.WriteString(field.varName + " " + field.typeName)
		}
		writer.WriteString(") ")
		if fn.returnType != "" {
			writer.WriteString(fn.returnType + " ")
		}
		writer.WriteString("{\n")
		writeIndent(gen.options.indent, writer)
		writer.WriteString("panic(\"TODO: Unimplemented method\")\n}\n")
	}
}

func (rust *RustGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	indent := rust.options.indent
	for _, fn := range typeDecl.methods {
		writeIndent(indent, writer)
		writer.WriteString("fn " + fn.name + "(&self")

		for _, field := range fn.fields {
			writer.WriteString(", ")
			writer.WriteString(field.varName + ": " + field.typeName)
		}
		writer.WriteString(") ")
		if fn.returnType != "" {
			writer.WriteString("-> " + fn.returnType + " ")
		}
		writer.WriteString("{\n")
		writeIndent(2*indent, writer)
		writer.WriteString("panic!(\"TODO: Unimplemented method\")\n")
		writeIndent(indent, writer)
		writer.WriteString("}\n")
	}
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

func capitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func (java *JavaGenerator) writeFields(fields []Field, writer *bufio.Writer) {
	indent := java.options.indent
	for _, field := range fields {
		writeIndent(indent, writer)
		if field.hasModifier(FIELD_CONST) {
			writer.WriteString("final ")
		}
		writer.WriteString(field.typeName)
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteString(" " + field.varName + ";\n")
	}
}

func (java *JavaGenerator) writeConstructor(t TypeDecl, writer *bufio.Writer) {
	indent := java.options.indent
	writeIndent(indent, writer)
	writer.WriteString(t.typeName + "(")
	join := newJoiner()
	for _, field := range t.fields {
		if join.join() {
			writer.WriteString(", ")
		}
		writer.WriteString(field.typeName)
		if field.hasModifier(FIELD_ARRAY) {
			writer.WriteString("[]")
		}
		writer.WriteString(" " + field.varName)
	}
	writer.WriteString(") {\n")
	for _, field := range t.fields {
		writeIndent(2*indent, writer)
		varName := field.varName
		assignment := "this." + varName + " = " + varName + ";\n"
		writer.WriteString(assignment)
	}
	writeIndent(indent, writer)
	writer.WriteString("}\n")
}

func (java *JavaGenerator) writeMethods(typeDecl TypeDecl, writer *bufio.Writer) {
	indent := java.options.indent
	for _, fn := range typeDecl.methods {
		writeIndent(indent, writer)
		writer.WriteString(fn.returnType + " " + fn.name)

		writer.WriteByte('(')
		fieldJoiner := newJoiner()
		for _, field := range fn.fields {
			if fieldJoiner.join() {
				writer.WriteString(", ")
			}
			writer.WriteString(field.typeName)
			if field.hasModifier(FIELD_ARRAY) {
				writer.WriteString("[]")
			}
			writer.WriteString(" " + field.varName)
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

// Type mappers

func toGoType(typeName string) string {
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

func toJavaType(typeName string) string {
	switch typeName {
	case "":
		return "void"
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

func toRustType(typeName string) string {
	switch typeName {
	case "string":
		return "String"
	default:
		return typeName
	}
}

// Keywords constants

var GO_KEYWORDS = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
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

var RUST_KEYWORDS = []string{
	"as", "break", "const", "continue", "crate", "else", "enum", "extern", "false",
	"fn", "for", "if", "impl", "in", "let", "loop", "match", "mod", "move", "mut",
	"pub", "ref", "return", "self", "Self", "static", "struct", "super", "trait",
	"true", "type", "unsafe", "use", "where", "while", "async", "await", "dyn",
}
