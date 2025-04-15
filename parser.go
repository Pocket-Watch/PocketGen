package main

import (
	"fmt"
	"os"
)

type Parser struct {
	filepath string
	lexer    Lexer
	structs  []TypeDecl
	tokenNow Token
}

type TypeDecl struct {
	line    LinePos
	name    string
	fields  []Field
	methods []FuncDecl
}

type FuncDecl struct {
	line       LinePos
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
	line      LinePos
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
	token := lexer.NextToken()

	parser := Parser{
		filepath: path,
		lexer:    lexer,
		structs:  make([]TypeDecl, 0),
		tokenNow: token,
	}

	return parser, true
}

func AdvanceToken(parser *Parser) Token {
	previous := parser.tokenNow
	parser.tokenNow = parser.lexer.NextToken()
	return previous
}

func PeekToken(parser *Parser) Token {
	return parser.tokenNow
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

type ParsingResult struct {
	success bool
	message string
}

func parserOk() ParsingResult {
	result := ParsingResult{
		success: true,
		message: "",
	}

	return result
}

func (parser *Parser) formatExpectedToken(found Token, format string, args ...any) string {
	var foundString string

	switch found.tokenType {
	case TOKEN_KEYWORD:
		foundString = fmt.Sprintf("keyword '%s'", found.tokenValue.string)
	case TOKEN_IDENTIFIER:
		foundString = fmt.Sprintf("identifier '%s'", found.tokenValue.string)
	default:
		foundString = fmt.Sprintf("token '%s'", TokenTypeToString(found.tokenType))
	}

	expectedString := fmt.Sprintf(format, args...)

	line := found.line

	message := fmt.Sprintf("%s:%v:%v Expected %s, but instead %s was found.", parser.filepath, line.number, line.offset, expectedString, foundString)
	return message
}

func (parser *Parser) expectedToken(expectedType TokenType, found Token) ParsingResult {
	expectedString := TokenTypeToString(expectedType)
	message := parser.formatExpectedToken(found, "token %s", expectedString)

	result := ParsingResult{
		success: false,
		message: message,
	}

	return result
}

func (parser *Parser) expectedKeyword(keywordType KeywordType, found Token) ParsingResult {
	message := parser.formatExpectedToken(found, "keyword '%s'", keywordType)

	result := ParsingResult{
		success: false,
		message: message,
	}

	return result
}

func (parser *Parser) parserErrorMessage(found Token, message string) ParsingResult {
	line := found.line
	message = fmt.Sprintf("%s:%v:%v %s", parser.filepath, line.number, line.offset, message)

	result := ParsingResult{
		success: false,
		message: message,
	}

	return result
}

func parseTypeField(parser *Parser, field *Field) ParsingResult {
	token := PeekToken(parser)
	if IsKeyword(token, KEYWORD_CONST) {
		addModifier(field, FIELD_CONST)
		AdvanceToken(parser)
	}

	//
	// Parse the field type
	//
	is_array_type := false
	token = PeekToken(parser)
	if IsType(token, TOKEN_SQUARE_OPEN) {
		is_array_type = true
		addModifier(field, FIELD_ARRAY)
		AdvanceToken(parser)
	}

	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_IDENTIFIER) {
		return parser.expectedToken(TOKEN_IDENTIFIER, token)
	}

	field.line = token.line
	field.typeName = token.tokenValue.string

	token = PeekToken(parser)
	if IsType(token, TOKEN_NULLABLE) {
		addModifier(field, FIELD_NULLABLE)
		AdvanceToken(parser)
	}

	if is_array_type {
		token = AdvanceToken(parser)
		if !IsType(token, TOKEN_SQUARE_CLOSE) {
			return parser.expectedToken(TOKEN_SQUARE_CLOSE, token)
		}
	}

	//
	// Parse the field variable name
	//
	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_IDENTIFIER) {
		return parser.expectedToken(TOKEN_IDENTIFIER, token)
	}

	field.varName = token.tokenValue.string

	return parserOk()
}

func parseFunctionDeclaration(parser *Parser, funcDecl *FuncDecl) ParsingResult {
	token := AdvanceToken(parser)
	funcDecl.line = token.line

	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_IDENTIFIER) {
		return parser.expectedToken(TOKEN_IDENTIFIER, token)
	}

	funcDecl.name = token.tokenValue.string

	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_ROUND_OPEN) {
		return parser.expectedToken(TOKEN_ROUND_OPEN, token)
	}

	token = PeekToken(parser)
	if IsType(token, TOKEN_ROUND_CLOSE) {
		AdvanceToken(parser)
	} else {
		for {
			field := Field{}
			result := parseTypeField(parser, &field)
			funcDecl.fields = append(funcDecl.fields, field)
			if !result.success {
				return result
			}

			token = PeekToken(parser)
			if !IsType(token, TOKEN_COMMA) {
				break
			}

			AdvanceToken(parser)
		}

		token := AdvanceToken(parser)
		if !IsType(token, TOKEN_ROUND_CLOSE) {
			return parser.expectedToken(TOKEN_ROUND_CLOSE, token)
		} 
	}

	token = PeekToken(parser)
	if IsType(token, TOKEN_IDENTIFIER) {
		funcDecl.returnType = token.tokenValue.string
		AdvanceToken(parser)
	}

	return parserOk()
}

func parseTypeDeclaration(parser *Parser, typeDecl *TypeDecl) ParsingResult {
	token := AdvanceToken(parser)
	typeDecl.line = token.line

	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_IDENTIFIER) {
		return parser.expectedToken(TOKEN_IDENTIFIER, token)
	}

	typeDecl.name = token.tokenValue.string

	token = AdvanceToken(parser)
	if !IsType(token, TOKEN_CURLY_OPEN) {
		return parser.expectedToken(TOKEN_CURLY_OPEN, token)
	}

	token = PeekToken(parser)
	if IsType(token, TOKEN_CURLY_CLOSE) {
		AdvanceToken(parser)
		return parserOk()
	}

	for {
		var result ParsingResult

		token := PeekToken(parser)
		if IsKeyword(token, KEYWORD_FUNC) {
			funcDecl := FuncDecl{}
			result = parseFunctionDeclaration(parser, &funcDecl)
			typeDecl.methods = append(typeDecl.methods, funcDecl)
		} else {
			field := Field{}
			result = parseTypeField(parser, &field)
			typeDecl.fields = append(typeDecl.fields, field)
		}

		if !result.success {
			return result
		}

		token = AdvanceToken(parser)
		if !IsType(token, TOKEN_SEMICOLON) {
			return parser.expectedToken(TOKEN_SEMICOLON, token)
		}

		token = PeekToken(parser)
		if IsType(token, TOKEN_CURLY_CLOSE) {
			AdvanceToken(parser)
			break
		}
	}

	return parserOk()
}

func ParseFile(parser *Parser) ParsingResult {
	for {
		token := PeekToken(parser)
		if IsType(token, TOKEN_EOF) {
			break
		}

		var result ParsingResult
		if IsKeyword(token, KEYWORD_TYPE) {
			var typeDecl TypeDecl
			result = parseTypeDeclaration(parser, &typeDecl)
			parser.structs = append(parser.structs, typeDecl)
		} else {
			return parser.expectedKeyword(KEYWORD_TYPE, token)
		}

		if !result.success {
			return result
		}
	}

	return parserOk()
}
