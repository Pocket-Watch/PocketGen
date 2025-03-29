package main

import "fmt"

type TokenType = int

const (
	TOKEN_KEYWORD     = iota
	TOKEN_IDENTIFIER  = iota
	TOKEN_COMMA       = iota
	TOKEN_RIGHT_CURLY = iota
	TOKEN_LEFT_CURLY  = iota
)

type KeywordType = int

const (
	KEYWORD_CONST = "const"
	KEYWORD_FUNC  = "func"
)

var KEYWORDS = []string{KEYWORD_CONST, KEYWORD_FUNC}

const (
	PRIMITIVE_INT    = iota
	PRIMITIVE_FLOAT  = iota
	PRIMITIVE_STRING = iota
	PRIMITIVE_CHAR   = iota
)

type TokenData struct {
	keywordType KeywordType
	stringData  string
	// Placegolder?
	indentifierData string
}

type Token struct {
	tokenType TokenType
	tokenData TokenData
}

type Lexer struct {
	lineNumber int32
	linePos    int32
	totalPos   int32
}

func main() {
	fmt.Println("Hello, World")
}
