package main

import (
	"fmt"
	"os"
)

type Parser struct {
	filepath string
	lexer    Lexer
	structs  []TypeDecl
}

type TypeDecl struct {
	name    string
	fields  []Field
	methods []FuncDecl
}

type FuncDecl struct {
	name       string
	fields     []Field
	returnType string
}

type FieldModifier = uint32

const (
	FIELD_NONE  FieldModifier = 0
	FIELD_CONST FieldModifier = (1 << iota)
	FIELD_ARRAY
	FIELD_NULLABLE
	FIELD_PRIMITIVE
)

type Field struct {
	varName   string
	typeName  string
	modifiers FieldModifier
}

func CreateParser(path string) (Parser, bool) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("ERROR: Failed to open file %v: %v.\n", path, err)
		return Parser{}, false
	}

	lexer := CreateLexer(data)

	parser := Parser{
		filepath: path,
		lexer:    lexer,
		structs:  make([]TypeDecl, 0),
	}

	return parser, true
}

// Returns true if modifier was already applied. This check might be useful when more modifier keywords are added
// and we don't want a modifier to repeat multiple times. For example:
//
//	public const public [string] name;
//
// Note that, this problem is irrelevant if the parser enforces particular modifier order.
func addModifier(field *Field, modifier FieldModifier) bool {
	alreadyApplied := field.modifiers&modifier != 0
	field.modifiers |= modifier
	return alreadyApplied
}

func hasModifier(field Field, modifier FieldModifier) bool {
	return field.modifiers&modifier != 0
}

type ParsingError struct {
	expected      Token
	found         Token
	customMessage string
}

func parserOk() (bool, ParsingError) {
	return true, ParsingError{}
}

func expectedToken(expectedType TokenType, found Token) (bool, ParsingError) {
	expected := Token{
		tokenType: expectedType,
	}

	error := ParsingError{
		expected:      expected,
		found:         found,
		customMessage: "",
	}

	return false, error
}

func expectedKeyword(keywordType KeywordType, found Token) (bool, ParsingError) {
	expected := Token{
		tokenType:  TOKEN_KEYWORD,
		tokenValue: TokenValue{string: keywordType},
	}

	error := ParsingError{
		expected:      expected,
		found:         found,
		customMessage: "",
	}

	return false, error
}

func parserErrorMessage(found Token, message string) (bool, ParsingError) {
	error := ParsingError{
		expected:      Token{},
		found:         found,
		customMessage: message,
	}

	return false, error

}

func parseTypeDeclaration(parser *Parser) (bool, ParsingError) {
	typeDecl := TypeDecl{}

	token := NextToken(&parser.lexer)
	if IsKeyword(token, KEYWORD_TYPE) {
		return expectedKeyword(KEYWORD_TYPE, token)
	}

	token = NextToken(&parser.lexer)
	if IsType(token, TOKEN_IDENTIFIER) {
		return expectedToken(TOKEN_IDENTIFIER, token)
	}

	typeDecl.name = token.tokenValue.string

	token = NextToken(&parser.lexer)
	if IsType(token, TOKEN_CURLY_OPEN) {
		return expectedToken(TOKEN_CURLY_OPEN, token)
	}

	// Parse fields here...

	token = NextToken(&parser.lexer)
	if IsType(token, TOKEN_CURLY_CLOSE) {
		return expectedToken(TOKEN_CURLY_CLOSE, token)
	}

	parser.structs = append(parser.structs, typeDecl)
	return parserOk()
}

func ParseFile(parser *Parser) (bool, ParsingError) {
	// 1. Peek token.
	// 2. Check end of file.
	//      yes -> parsing ended.
	//      no -> parser type declaration.
	return parserOk()
}
