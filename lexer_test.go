package main

import (
	"testing"
)

func TestValidTokens(t *testing.T) {
	input := `
        # This is a sample .tg file for lexer testing.

        const{ }

        # This is a comment
        , ###
        ranoasdconst  	func
        func123ł}ł ą2 # const

        a#a
        `

	lexer := CreateLexer([]byte(input))

	expectedTokens := asTokens(TOKEN_KEYWORD, TOKEN_CURLY_OPEN, TOKEN_CURLY_CLOSE,
		TOKEN_COMMA, TOKEN_IDENTIFIER, TOKEN_KEYWORD,
		TOKEN_IDENTIFIER, TOKEN_CURLY_CLOSE, TOKEN_IDENTIFIER,
		TOKEN_IDENTIFIER, TOKEN_IDENTIFIER)

	actualTokens := make([]Token, 0)
	// This will probably be a parser function.
	token := NextToken(&lexer)
	for {
		PrintToken(token)

		// Without this, it turns into an infinite loop
		if IsType(token, TOKEN_EOF) {
			break
		}

		if IsType(token, TOKEN_ERROR) {
			t.Error("Token bad, also the parser will handle this.")
			break
		}

		actualTokens = append(actualTokens, token)
		token = NextToken(&lexer)
	}
	compareTokens(expectedTokens, actualTokens, t)
}

func asTokens(types ...TokenType) []Token {
	tokens := make([]Token, 0, len(types))
	for i := 0; i < len(types); i++ {
		pos := LinePos{0, 0}
		token := makeToken(types[i], pos)
		tokens = append(tokens, token)
	}
	return tokens
}

func compareTokens(expected []Token, actual []Token, t *testing.T) {
	if len(expected) != len(actual) {
		t.Errorf("Different number of tokens, expected %v, actual %v", len(expected), len(actual))
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].tokenType != actual[i].tokenType {
			t.Errorf("Different tokens at index %v, expected %v, actual %v", i, expected[i].tokenType, actual[i].tokenType)
			PrintToken(expected[i])
			PrintToken(actual[i])
			return
		}
	}
	return
}
