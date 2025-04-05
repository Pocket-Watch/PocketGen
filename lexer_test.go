package main

import (
	"testing"
)

func testToken(tokenType TokenType, tokenString string, tokenInt int, tokenLine int, tokenOffset int) Token {
	line := LinePos{
		number: tokenLine,
		offset: tokenOffset,
	}

	value := TokenValue{
		string: tokenString,
		int:    tokenInt,
	}

	token := Token{
		tokenType:  tokenType,
		tokenValue: value,
		line:       line,
	}

	return token
}

func makeTestTokens(tokens ...Token) []Token {
	expected := make([]Token, 0)
	for _, token := range tokens {
		expected = append(expected, token)
	}

	return expected
}

func TestValidTokens(t *testing.T) {
	input := `# This is a sample .tg file for lexer testing.

const{ }

# This is a comment
, ###
ranoasdconst  	func
func123ł}ł ą2 # const

a#a
`

	lexer := CreateLexer([]byte(input))

	// TODO(kihau): Lexer should have a function to generate those automatically.
	expectedTokens := makeTestTokens(
		testToken(TOKEN_KEYWORD,     "const",        0, 3, 1),
		testToken(TOKEN_CURLY_OPEN,  "",             0, 3, 6),
		testToken(TOKEN_CURLY_CLOSE, "",             0, 3, 8),
		testToken(TOKEN_COMMA,       "",             0, 6, 1),
		testToken(TOKEN_IDENTIFIER,  "ranoasdconst", 0, 7, 3),
		testToken(TOKEN_KEYWORD,     "func",         0, 7, 18),
		testToken(TOKEN_IDENTIFIER,  "func123ł",     0, 8, 1),
		testToken(TOKEN_CURLY_CLOSE, "",             0, 8, 9),
		testToken(TOKEN_IDENTIFIER,  "ł",            0, 8, 10),
		testToken(TOKEN_IDENTIFIER,  "ą2",           0, 8, 12),
		testToken(TOKEN_IDENTIFIER,  "a",            0, 10, 1),
	)

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
