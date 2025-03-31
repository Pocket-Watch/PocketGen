package main

type TokenType = int

const (
	TOKEN_KEYWORD = iota
	TOKEN_IDENTIFIER
	TOKEN_COMMA
	TOKEN_LEFT_CURLY
	TOKEN_RIGHT_CURLY
	TOKEN_LEFT_SQUARE
	TOKEN_RIGHT_SQUARE
)

type KeywordType = int

const (
	KEYWORD_CONST = "const"
	KEYWORD_FUNC  = "func"
	KEYWORD_TYPE  = "type"
)

var KEYWORDS = []string{KEYWORD_CONST, KEYWORD_FUNC, KEYWORD_TYPE}

const (
	PRIMITIVE_INT = iota
	PRIMITIVE_FLOAT
	PRIMITIVE_STRING
	PRIMITIVE_CHAR
)

type TokenData struct {
	keywordType KeywordType
	stringData  string
	// Placeholder?
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
