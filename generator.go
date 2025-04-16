package main

import (
	"bufio"
	"fmt"
	"os"
)

const INDENT = 4

// How to handle every Write since they can all return errors?
func generateJS(types []TypeDecl, writer *bufio.Writer) {
	for _, t := range types {
		// TypeDecl has no information if it's an enum or class
		// Enums usually don't have any modifiers? So we probably need EnumDecl but it must be separate from TypeDecl?
		writer.WriteString("class ")
		writer.WriteString(t.name)
		writer.WriteString(" {\n")

		writeIndent(INDENT, writer)
		writeJsConstructor(t.fields, INDENT, writer)
		writeJsMethods(t.methods, INDENT, writer)

		writer.Flush()
	}
}

func writeJsMethods(methods []FuncDecl, currentIndent int, writer *bufio.Writer) {
	for _, fn := range methods {
		writeIndent(currentIndent, writer)
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

func writeJsConstructor(fields []Field, currentIndent int, writer *bufio.Writer) {
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
		writeIndent(currentIndent+INDENT, writer)
		writer.WriteString("this.")
		writer.WriteString(field.varName)
		writer.WriteString(" = ")
		writer.WriteString(field.varName)
		writer.WriteByte(';')
		writer.WriteByte('\n')
	}
	writeIndent(currentIndent, writer)
	writer.WriteString("}\n")
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
	defer file.Close()

	writer := bufio.NewWriter(file)
	return writer
}
