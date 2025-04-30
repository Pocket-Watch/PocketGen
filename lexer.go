package main

import (
	"fmt"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"
)

type TokenType = int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_ERROR
	TOKEN_KEYWORD
	TOKEN_IDENTIFIER
	TOKEN_COMMA
	TOKEN_SEMICOLON
	TOKEN_NULLABLE
	TOKEN_CURLY_OPEN
	TOKEN_CURLY_CLOSE
	TOKEN_ROUND_OPEN
	TOKEN_ROUND_CLOSE
	TOKEN_SQUARE_OPEN
	TOKEN_SQUARE_CLOSE
)

type TokenAsString struct {
	exactType  string
	prettyType string
	value      string
}

var TOKEN_LOOKUP = []TokenAsString{
	{"TOKEN_EOF", "eof", "end of file"},
	{"TOKEN_ERROR", "error", "error"},
	{"TOKEN_KEYWORD", "keyword", "keyword"},
	{"TOKEN_IDENTIFIER", "identifier", "identifier"},
	{"TOKEN_COMMA", "comma", ","},
	{"TOKEN_SEMICOLON", "semicolon", ";"},
	{"TOKEN_NULLABLE", "nullable", "?"},
	{"TOKEN_CURLY_OPEN", "curly open", "{"},
	{"TOKEN_CURLY_CLOSE", "curly close", "}"},
	{"TOKEN_ROUND_OPEN", "round open", "("},
	{"TOKEN_ROUND_CLOSE", "round close", ")"},
	{"TOKEN_SQUARE_OPEN", "square open", "["},
	{"TOKEN_SQUARE_CLOSE", "square close", "]"},
}

func TokenTypeToString(tokenType TokenType) string {
	if tokenType < 0 || tokenType >= len(TOKEN_LOOKUP) {
		return "TOKEN_UNKNOWN"
	}

	return TOKEN_LOOKUP[tokenType].exactType
}

func TokenTypeToStringPretty(tokenType TokenType) string {
	if tokenType < 0 || tokenType >= len(TOKEN_LOOKUP) {
		return "unknown"
	}

	return TOKEN_LOOKUP[tokenType].prettyType
}

func TokenValueToString(token Token) string {
	tokenType := token.tokenType
	if tokenType < 0 || tokenType >= len(TOKEN_LOOKUP) {
		return "unknown"
	}

	switch tokenType {
	case TOKEN_KEYWORD, TOKEN_IDENTIFIER:
		return token.tokenValue.string

	default:
		return TOKEN_LOOKUP[tokenType].value
	}
}

type TokenErrorType = int

const (
	ERROR_INVALID_RUNE_ENCODING TokenErrorType = iota
	ERROR_UNKNOWN_RUNE_SYMBOL
	ERROR_UNCLOSED_STRING
	ERROR_UNCLOSED_BLOCK_COMMENT
)

type KeywordType = string

const (
	KEYWORD_TYPE  KeywordType = "type"
	KEYWORD_CONST KeywordType = "const"
	KEYWORD_FUNC  KeywordType = "func"
)

var KEYWORD_LOOKUP = []string{
	"type",
	"const",
	"func",
}

var KEYWORDS = []KeywordType{KEYWORD_TYPE, KEYWORD_CONST, KEYWORD_FUNC}

var PRIMITIVES = []string{
	"i8", "i16", "i32", "i64",
	"u8", "u16", "u32", "u64",
	"string", "char", "bool",
	"f32", "f64",
}

// TokenValue union, matched with TokenType. Go has no unions so this one is fat.
type TokenValue struct {
	string string
	int    int
}

type Token struct {
	line       LinePos
	tokenType  TokenType
	tokenValue TokenValue
}

// LinePos (debug struct) - number and offset are 1-indexed
type LinePos struct {
	number int
	offset int
}

type Lexer struct {
	data     []byte
	line     LinePos
	pos      int
	runeNow  rune
	runeSize int
}

func CreateLexer(data []byte) Lexer {
	line := LinePos{
		number: 1,
		offset: 1,
	}

	lexer := Lexer{
		data: data,
		line: line,
		pos:  0,
	}

	lexer.startRune()
	return lexer
}

func (lexer *Lexer) startRune() {
	rune, runeSize := utf8.DecodeRune(lexer.data)
	lexer.runeNow = rune
	lexer.runeSize = runeSize
}

func (lexer *Lexer) nextRune() rune {
	if lexer.runeNow == '\n' {
		lexer.line.number += 1
		lexer.line.offset = 1
	} else {
		lexer.line.offset += 1
	}

	lexer.pos += lexer.runeSize
	lexer.runeNow = 0
	lexer.runeSize = 0

	if lexer.pos >= len(lexer.data) {
		return 0
	}

	rune, runeSize := utf8.DecodeRune(lexer.data[lexer.pos:])
	lexer.runeNow = rune
	lexer.runeSize = runeSize
	return rune
}

func (lexer *Lexer) peekRune() rune {
	return lexer.runeNow
}

func makeToken(tokenType TokenType, line LinePos) Token {
	token := Token{
		line:      line,
		tokenType: tokenType,
	}

	return token
}

func makeError(errorType TokenErrorType, line LinePos) Token {
	value := TokenValue{int: errorType}
	token := Token{
		line:       line,
		tokenType:  TOKEN_ERROR,
		tokenValue: value,
	}

	return token
}

func makeKeyword(keyword KeywordType, line LinePos) Token {
	value := TokenValue{string: keyword}
	token := Token{
		line:       line,
		tokenType:  TOKEN_KEYWORD,
		tokenValue: value,
	}

	return token
}

func makeIdentifier(identifier string, line LinePos) Token {
	value := TokenValue{string: identifier}
	token := Token{
		line:       line,
		tokenType:  TOKEN_IDENTIFIER,
		tokenValue: value,
	}

	return token
}

func parseWord(lexer *Lexer) (string, bool) {
	// NOTE(kihau):
	//     A lexer word must always begin with a unicode letter. This applies to both identifiers and keywords.
	//     A word might contain digits and underscores within its string. Here are some examples of correct lexer words:
	//         - hello
	//         - hi123_
	//         - a_b_c_d
	//         - a123_2
	//         - z______________

	rune := lexer.peekRune()
	if !unicode.IsLetter(rune) {
		return "", false
	}

	startPos := lexer.pos
	lexer.nextRune()

	for {
		rune := lexer.peekRune()
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '_' {
			break
		}

		lexer.nextRune()
	}

	wordSlice := lexer.data[startPos:lexer.pos]

	// NOTE(kihau): Maybe the byte slice should be returned instead of the string? Does this impact runtime performance?
	return string(wordSlice), true
}

func skipComment(lexer *Lexer) {
	rune := lexer.peekRune()
	for rune != '\n' && rune != 0 {
		rune = lexer.nextRune()
	}
}

func IsType(token Token, tokenType TokenType) bool {
	return token.tokenType == tokenType
}

func IsKeyword(token Token, keywordType KeywordType) bool {
	return token.tokenType == TOKEN_KEYWORD && token.tokenValue.string == keywordType
}

func (lexer *Lexer) NextToken() Token {
	return nextToken(lexer)
}

var nextToken = func(lexer *Lexer) Token {
	for {
		rune := lexer.peekRune()
		line := lexer.line

		switch rune {

		case 0:
			return makeToken(TOKEN_EOF, line)

		case utf8.RuneError:
			return makeError(ERROR_INVALID_RUNE_ENCODING, line)

		// TODO(kihau): Maybe add other white-space characters?
		case ' ', '\t', '\v', '\r', '\n', '\f':
			lexer.nextRune()
			continue

		case '{':
			lexer.nextRune()
			return makeToken(TOKEN_CURLY_OPEN, line)

		case '}':
			lexer.nextRune()
			return makeToken(TOKEN_CURLY_CLOSE, line)

		case '(':
			lexer.nextRune()
			return makeToken(TOKEN_ROUND_OPEN, line)

		case ')':
			lexer.nextRune()
			return makeToken(TOKEN_ROUND_CLOSE, line)

		case '[':
			lexer.nextRune()
			return makeToken(TOKEN_SQUARE_OPEN, line)

		case ']':
			lexer.nextRune()
			return makeToken(TOKEN_SQUARE_CLOSE, line)

		case ',':
			lexer.nextRune()
			return makeToken(TOKEN_COMMA, line)

		case ';':
			lexer.nextRune()
			return makeToken(TOKEN_SEMICOLON, line)

		case '?':
			lexer.nextRune()
			return makeToken(TOKEN_NULLABLE, line)

		case '#':
			skipComment(lexer)

		default:
			word, ok := parseWord(lexer)
			if !ok {
				return makeError(ERROR_UNKNOWN_RUNE_SYMBOL, line)
			}

			if slices.Contains(KEYWORDS, word) {
				return makeKeyword(word, line)
			} else {
				return makeIdentifier(word, line)
			}
		}
	}
}

func TokenToString(token Token) string {
	typeString := TokenTypeToString(token.tokenType)
	valueString := TokenValueToString(token)

	line := fmt.Sprintf("%v:%v ", token.line.number, token.line.offset)
	string := fmt.Sprintf("%-14s %-6s - %s", typeString, line, valueString)
	return string
}

func PrintToken(token Token) {
	tokenString := TokenToString(token)
	fmt.Printf("%v\n", tokenString)
}

func GenerateTestTokens(data []byte) {
	lexer := CreateLexer(data)

	fmt.Println("expectedTokens := makeTestTokens(")
	for {
		token := lexer.NextToken()

		typeString := TokenTypeToString(token.tokenType)
		fmt.Printf("\ttestToken(%v, \"%v\", %v, %v, %v),\n", typeString, token.tokenValue.string, token.tokenValue.int, token.line.number, token.line.offset)

		if IsType(token, TOKEN_EOF) || IsType(token, TOKEN_ERROR) {
			break
		}
	}
	fmt.Println(")")
}

func RunLexerScratch(path string) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("ERROR: Failed to open sample tg file: %v\n", err)
		return
	}

	lexer := CreateLexer(data)

	token := lexer.NextToken()
	for {
		PrintToken(token)

		if IsType(token, TOKEN_EOF) {
			break
		}

		if IsType(token, TOKEN_ERROR) {
			println("Token bad, also the parser will handle this.")
			break
		}

		token = lexer.NextToken()
	}
}
