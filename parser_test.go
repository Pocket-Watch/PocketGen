package main

import (
	"fmt"
	"testing"
)

func makeTokenWithValue(tokenType TokenType, val string) Token {
	return Token{
		line:      LinePos{},
		tokenType: tokenType,
		tokenValue: TokenValue{
			string: val,
		},
	}
}

func mockLexer(tokens []Token) Lexer {
	lexer := CreateLexer([]byte(""))
	tokenIndex := 0

	nextToken = func(lexer *Lexer) Token {
		if tokenIndex < len(tokens) {
			tok := tokens[tokenIndex]
			tokenIndex += 1
			return tok
		} else {
			// Panic-less tokenizer
			return makeToken(TOKEN_EOF, LinePos{})
		}
	}
	return lexer
}

func TestCat(t *testing.T) {
	tokens := makeTestTokens(
		makeTokenWithValue(TOKEN_KEYWORD, "type"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "Cat"),
		makeTokenWithValue(TOKEN_CURLY_OPEN, ""),
		// field name
		makeTokenWithValue(TOKEN_KEYWORD, "const"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "string"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "name"),
		makeTokenWithValue(TOKEN_SEMICOLON, ""),
		// field age
		makeTokenWithValue(TOKEN_IDENTIFIER, "u32"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "age"),
		makeTokenWithValue(TOKEN_SEMICOLON, ""),
		// func meow
		makeTokenWithValue(TOKEN_KEYWORD, "func"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "meow"),
		makeTokenWithValue(TOKEN_ROUND_OPEN, ""),
		makeTokenWithValue(TOKEN_IDENTIFIER, "string"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "sound"),
		makeTokenWithValue(TOKEN_COMMA, ""),
		makeTokenWithValue(TOKEN_IDENTIFIER, "u32"),
		makeTokenWithValue(TOKEN_IDENTIFIER, "volume"),
		makeTokenWithValue(TOKEN_ROUND_CLOSE, ""),
		makeTokenWithValue(TOKEN_IDENTIFIER, "string"),
		makeTokenWithValue(TOKEN_SEMICOLON, ""),

		makeTokenWithValue(TOKEN_CURLY_CLOSE, ""),
		makeTokenWithValue(TOKEN_EOF, ""),
	)

	lexer := mockLexer(tokens)
	for i := 0; i < 30; i++ {
		fmt.Println(lexer.NextToken())
	}
}
