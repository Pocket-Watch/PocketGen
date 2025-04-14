package main

import (
	"os"
	"testing"
)

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

func TestValidDefinitions(t *testing.T) {
	expectedTokens := makeTestTokens(
		testToken(TOKEN_KEYWORD, "type", 0, 1, 1),
		testToken(TOKEN_IDENTIFIER, "Cat", 0, 1, 6),
		testToken(TOKEN_CURLY_OPEN, "", 0, 1, 10),
		testToken(TOKEN_KEYWORD, "const", 0, 2, 5),
		testToken(TOKEN_IDENTIFIER, "string", 0, 2, 11),
		testToken(TOKEN_IDENTIFIER, "name", 0, 2, 18),
		testToken(TOKEN_SEMICOLON, "", 0, 2, 22),
		testToken(TOKEN_IDENTIFIER, "u32", 0, 3, 5),
		testToken(TOKEN_IDENTIFIER, "age", 0, 3, 9),
		testToken(TOKEN_SEMICOLON, "", 0, 3, 12),
		testToken(TOKEN_KEYWORD, "func", 0, 5, 5),
		testToken(TOKEN_IDENTIFIER, "meow", 0, 5, 10),
		testToken(TOKEN_ROUND_OPEN, "", 0, 5, 14),
		testToken(TOKEN_IDENTIFIER, "string", 0, 5, 15),
		testToken(TOKEN_IDENTIFIER, "sound", 0, 5, 22),
		testToken(TOKEN_COMMA, "", 0, 5, 27),
		testToken(TOKEN_IDENTIFIER, "u32", 0, 5, 29),
		testToken(TOKEN_IDENTIFIER, "volume", 0, 5, 33),
		testToken(TOKEN_ROUND_CLOSE, "", 0, 5, 39),
		testToken(TOKEN_IDENTIFIER, "string", 0, 5, 41),
		testToken(TOKEN_SEMICOLON, "", 0, 5, 47),
		testToken(TOKEN_CURLY_CLOSE, "", 0, 6, 1),
		testToken(TOKEN_EOF, "", 0, 6, 17),
	)
	validateLexerTokens(t, "test/cat.tg", expectedTokens)
}

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

	token := lexer.NextToken()
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

		if IsType(token, TOKEN_EOF) {
			break
		}

		expected := expectedTokens[counter]
		if !compareTokens(t, expected, token, counter) {
			break
		}

		counter += 1
		token = lexer.NextToken()
	}
}

func compareTokens(t *testing.T, expected Token, actual Token, tokenIndex int) bool {
	if expected.tokenType != actual.tokenType {
		t.Errorf("Tokens at index %v are of different types.", tokenIndex)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	if expected.tokenValue.int != actual.tokenValue.int || expected.tokenValue.string != actual.tokenValue.string {
		t.Errorf("Tokens at index %v have different token values (int or string).", tokenIndex)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	if expected.line.number != actual.line.number || expected.line.offset != actual.line.offset {
		t.Errorf("Tokens at index %v are at different positions (line or offset).", tokenIndex)
		t.Errorf("Expected: %v", TokenToString(expected))
		t.Errorf("Actual:   %v", TokenToString(actual))

		return false
	}

	return true
}
