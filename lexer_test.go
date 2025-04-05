package main

import (
	"os"
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

func validateLexerTokens(t *testing.T, testPath string, expectedTokens []Token) {
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Errorf("Failed to open lexer1 test file: %v.", err)
		return
	}

	lexer := CreateLexer(data)
	counter := 0

	token := NextToken(&lexer)
	for {
		// NOTE(kihau):
		//    This is a sanity check. Because all token streams end with EOF, the loop should always
		//    exit before expectedTokens run out and because of that branch is expected to always be 'false'.
		if counter >= len(expectedTokens) {
			t.Errorf("Number of actual tokens is greater than those expected in the test.\n")
			t.Errorf("Expected: End of file.\n")
			t.Errorf("Found:    %v.\n", TokenToString(token))
			break
		}

		expected := expectedTokens[counter]
		if !compareTokens(t, expected, token, counter) {
			break
		}

		if IsType(token, TOKEN_EOF) {
			break
		}

		counter += 1
		token = NextToken(&lexer)
	}
}

func TestLexer1(t *testing.T) {
	// NOTE(kihau): Generated using the GenerateTestTokens() function.
	expectedTokens := makeTestTokens(
		testToken(TOKEN_KEYWORD, "const", 0, 3, 1),
		testToken(TOKEN_CURLY_OPEN, "", 0, 3, 6),
		testToken(TOKEN_CURLY_CLOSE, "", 0, 3, 8),
		testToken(TOKEN_COMMA, "", 0, 6, 1),
		testToken(TOKEN_IDENTIFIER, "ranoasdconst", 0, 7, 3),
		testToken(TOKEN_KEYWORD, "func", 0, 7, 18),
		testToken(TOKEN_IDENTIFIER, "func123ł", 0, 8, 1),
		testToken(TOKEN_CURLY_CLOSE, "", 0, 8, 9),
		testToken(TOKEN_IDENTIFIER, "ł", 0, 8, 10),
		testToken(TOKEN_IDENTIFIER, "ą2", 0, 8, 12),
		testToken(TOKEN_IDENTIFIER, "a", 0, 10, 1),
		testToken(TOKEN_EOF, "", 0, 11, 1),
	)

	validateLexerTokens(t, "test/lexer1.tg", expectedTokens)
}

func compareTokens(t *testing.T, expected Token, actual Token, number int) bool {
	if expected.tokenType != actual.tokenType {
		t.Errorf("Expected token at index %v has different token type than the actual token.", number)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	if expected.tokenValue.int != actual.tokenValue.int || expected.tokenValue.string != actual.tokenValue.string {
		t.Errorf("Expected token at index %v has different value content than the actual token.", number)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	if expected.line.number != actual.line.number || expected.line.offset != actual.line.offset {
		t.Errorf("Expected token at index %v is at line different position than the actual token.", number)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	return true
}
