package main

import (
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

func createParserWithLexer(lexer Lexer) Parser {
	firstToken := lexer.NextToken()
	parser := Parser{
		filepath: "test",
		lexer:    lexer,
		structs:  make([]TypeDecl, 0),
		tokenNow: firstToken,
	}
	return parser
}

func TestSingleCatType(t *testing.T) {
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
	parser := createParserWithLexer(lexer)

	result := ParseFile(&parser)
	if !result.success {
		t.Error(result.message)
		return
	}

	if len(parser.structs) != 1 {
		t.Errorf("Expected 1 type, found %v", len(parser.structs))
		return
	}

	cat := &parser.structs[0]

	if cat.name != "Cat" {
		t.Errorf("Expected different type name")
		return
	}

	if len(cat.fields) != 2 {
		t.Errorf("Expected 2 fields, found %v", len(cat.fields))
		return
	}
	nameField := &cat.fields[0]
	if !nameField.hasModifier(FIELD_CONST) || nameField.typeName != "string" || nameField.varName != "name" {
		t.Errorf("Expected field[0] to be 'const string name'")
		return
	}

	ageField := &cat.fields[1]
	if ageField.modifiers != 0 {
		t.Errorf("Expected field[1] to have no modifiers, actual value %v", ageField.modifiers)
		return
	}
	if ageField.typeName != "u32" || ageField.varName != "age" {
		t.Errorf("Expected field[1] to be 'u32 age', found typeName: %v, varName: %v",
			ageField.typeName, ageField.varName)
		return
	}

	if len(cat.methods) != 1 {
		t.Errorf("Expected 1 method, found %v", len(cat.methods))
		return
	}

	meow := &cat.methods[0]
	if meow.name != "meow" {
		t.Errorf("Expected 'meow' as method name")
		return
	}
	if meow.returnType != "string" {
		t.Errorf("Expected 'string' as return type of method")
		return
	}

	if len(meow.fields) != 2 {
		t.Errorf("Expected 2 arguments to method meow")
		return
	}

	sound := &meow.fields[0]
	if sound.typeName != "string" || sound.varName != "sound" {
		t.Errorf("Expected 'string sound' as 1st argument to method meow")
		return
	}

	volume := &meow.fields[1]
	if volume.typeName != "u32" || volume.varName != "volume" {
		t.Errorf("Expected 'u32 volume' as 2nd argument to method meow")
		return
	}
}
