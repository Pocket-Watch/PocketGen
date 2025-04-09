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

		for _, field := range t.fields {
			writeIndent(INDENT, writer)

			if hasModifier(field, FIELD_CONST) {
				writer.WriteString("const ")
			}

			writer.WriteString(field.varName)
			writer.WriteString(";\n")
		}

		for _, fn := range t.methods {
			writeIndent(INDENT, writer)
			writer.WriteString("function ")
			writer.WriteString(fn.name)
			writer.WriteByte('(')

			firstModifier := true
			for _, field := range fn.fields {
				if firstModifier {
					firstModifier = false
				} else {
					writer.WriteString(", ")
				}

				if hasModifier(field, FIELD_CONST) {
					writer.WriteString("const ")
				}

				writer.WriteString(field.varName)
			}
			writer.WriteByte(')')
			if fn.returnType != "" {
				writer.WriteString(" " + fn.returnType)
			}
			writer.WriteString(";\n")
		}
		writer.WriteString("}\n")
		writer.Flush()
	}
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
