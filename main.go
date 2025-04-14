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

func LexerDebuggingThing() {
	data, err := os.ReadFile("test/cat.tg")

	if err != nil {
		fmt.Printf("ERROR: Failed to open sample tg file: %v\n", err)
		return
	}

	lexer := CreateLexer(data)

	// This will probably be a parser function.
	token := NextToken(&lexer)
	for {
		PrintToken(token)

		if IsType(token, TOKEN_EOF) {
			break
		}

		if IsType(token, TOKEN_ERROR) {
			println("Token bad, also the parser will handle this.")
			break
		}

		token = NextToken(&lexer)
	}
}

func main() {
	LexerDebuggingThing()

	parser, success := CreateParser("test/cat.tg")
	if !success {
		os.Exit(1)
	}

	result := ParseFile(&parser)
	if !result.success {
		fmt.Println()
		fmt.Println(result.message)
		fmt.Println()
	}

	exec, err := os.Executable()
	if err != nil {
		exec = ""
	}
	fmt.Println("Usage:")
	fmt.Printf("  %v\n", exec)
}
