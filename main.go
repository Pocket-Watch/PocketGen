package main

import (
	"fmt"
	"os"
)

// Metagen format grammar (TODO):
//   newline        := // The U+000A character.
//   unicode_char   := // Any unicode character except newline.
//   unicode_letter := // Any unicode character categorized as "Letter".
//
//   letter     := unicode_letter | "_"
//   identifier := letter { letter | digit }
//   string     := '"' { unicode_char } '"'
//
//   metagen := typeDecl
//   		 | funcDecl
//   		 | varDecl
//
//   typeDecl := "type" identifier "{" { varDecl | funcDecl } "}"
//
//   funcDecl := "func" identifier "(" { varDecl comma } ")"
//
//   varDecl := identifier ( "[" "]" ) identifier

func LexerDebuggingThing(path string) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("ERROR: Failed to open sample tg file: %v\n", err)
		return
	}

	lexer := CreateLexer(data)

	// This will probably be a parser function.
	token := lexer.NextToken()
	for {
		PrintToken(token)

		if IsType(token, TOKEN_EOF) {
			break
		}

		if IsType(token, TOKEN_ERROR) {
			println("Token bad, also the parser will handle this.")
			break
		}

		token = lexer.NextToken()
	}
}

func main() {
	exec, err := os.Executable()
	if err != nil {
		exec = ""
	}
	fmt.Println("Usage:")
	fmt.Printf("  %v\n", exec)
	fmt.Println()

	LexerDebuggingThing("test/cat.tg")

	fmt.Println()

	parser, success := CreateParser("test/cat.tg")
	if !success {
		os.Exit(1)
	}

	result := ParseFile(&parser)
	if !result.success {
		fmt.Println(result.message)
		os.Exit(1)
	}

	result = TypecheckFile(&parser)
	if !result.success {
		fmt.Println(result.message)
		os.Exit(1)
	}
}
