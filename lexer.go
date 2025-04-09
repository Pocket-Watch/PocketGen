package main

import (
	"fmt"
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
)

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

var KEYWORDS = []KeywordType{KEYWORD_TYPE, KEYWORD_CONST, KEYWORD_FUNC}

const (
	PRIMITIVE_INT = iota
	PRIMITIVE_FLOAT
	PRIMITIVE_STRING
	PRIMITIVE_CHAR
)

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

	// TODO(kihau):
	//     Store previous and current token to allow implementation of Peek() and PeekNext() functions.
	//     This could also be a token ring buffer when larger token sequence is needed for parsing.
	// tokenNow  Token
	// tokenNext Token
}

/*
func startToken(lexer *Lexer) {
	lexer.tokenNow  = NextToken(lexer)
	lexer.tokenNext = NextToken(lexer)
}

func PeekToken(lexer *Lexer) Token {
	return lexer.tokenNow
}

func PeekNextToken(lexer *Lexer) Token {
	return lexer.tokenNext
}

func AdvanceToken(lexer *Lexer) {
	lexer.tokenNow = lexer.tokenNext
	lexer.tokenNext = NextToken(lexer)
}
*/

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

	startRune(&lexer)
	return lexer
}

func startRune(lexer *Lexer) {
	rune, runeSize := utf8.DecodeRune(lexer.data)
	lexer.runeNow = rune
	lexer.runeSize = runeSize
}

func nextRune(lexer *Lexer) rune {
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

func peekRune(lexer *Lexer) rune {
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

	rune := peekRune(lexer)
	if !unicode.IsLetter(rune) {
		return "", false
	}

	startPos := lexer.pos
	nextRune(lexer)

	for {
		rune := peekRune(lexer)
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '_' {
			break
		}

		nextRune(lexer)
	}

	wordSlice := lexer.data[startPos:lexer.pos]

	// NOTE(kihau): Maybe the byte slice should be returned instead of the string? Does this impact runtime performance?
	return string(wordSlice), true
}

func skipComment(lexer *Lexer) {
	rune := peekRune(lexer)
	for rune != '\n' && rune != 0 {
		rune = nextRune(lexer)
	}
}

func IsType(token Token, tokenType TokenType) bool {
	return token.tokenType == tokenType
}

func IsKeyword(token Token, keywordType KeywordType) bool {
	return token.tokenType == TOKEN_KEYWORD && token.tokenValue.string == keywordType
}

func NextToken(lexer *Lexer) Token {
	for {
		rune := peekRune(lexer)
		line := lexer.line

		switch rune {

		case 0:
			return makeToken(TOKEN_EOF, line)

		case utf8.RuneError:
			return makeError(ERROR_INVALID_RUNE_ENCODING, line)

		// TODO(kihau): Maybe add other white-space characters?
		case ' ', '\t', '\v', '\r', '\n', '\f':
			nextRune(lexer)
			continue

		case '{':
			nextRune(lexer)
			return makeToken(TOKEN_CURLY_OPEN, line)

		case '}':
			nextRune(lexer)
			return makeToken(TOKEN_CURLY_CLOSE, line)

		case '(':
			nextRune(lexer)
			return makeToken(TOKEN_ROUND_OPEN, line)

		case ')':
			nextRune(lexer)
			return makeToken(TOKEN_ROUND_CLOSE, line)

		case ',':
			nextRune(lexer)
			return makeToken(TOKEN_COMMA, line)

		case ';':
			nextRune(lexer)
			return makeToken(TOKEN_SEMICOLON, line)

		case '?':
			nextRune(lexer)
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
	var name string
	var value string

	switch token.tokenType {
	case TOKEN_EOF:
		name = "EOF"
		value = "end of file"

	case TOKEN_ERROR:
		name = "ERROR"
		value = fmt.Sprintf("%v", token.tokenValue.int)

	case TOKEN_KEYWORD:
		name = "KEYWORD"
		value = fmt.Sprintf("%v", token.tokenValue.string)

	case TOKEN_IDENTIFIER:
		name = "IDENTIFIER"
		value = fmt.Sprintf("'%v'", token.tokenValue.string)

	case TOKEN_COMMA:
		name = "COMMA"
		value = ","

	case TOKEN_SEMICOLON:
		name = "SEMICOLON"
		value = ";"

	case TOKEN_NULLABLE:
		name = "NULLABLE"
		value = "?"

	case TOKEN_CURLY_OPEN:
		name = "CURLY OPEN"
		value = "{"

	case TOKEN_CURLY_CLOSE:
		name = "CURLY CLOSE"
		value = "}"

	case TOKEN_ROUND_OPEN:
		name = "ROUND OPEN"
		value = "("

	case TOKEN_ROUND_CLOSE:
		name = "ROUND CLOSE"
		value = ")"

	default:
		name = "<UNKNOWN TOKEN>"
		value = fmt.Sprintf("%v", token.tokenType)
	}

	line := fmt.Sprintf("%v:%v ", token.line.number, token.line.offset)
	string := fmt.Sprintf("%-14s %-6s - %s", name, line, value)
	return string
}

func PrintToken(token Token) {
	tokenString := TokenToString(token)
	fmt.Printf("%v\n", tokenString)
}

func TokenTypeToString(tokenType TokenType) string {
	switch tokenType {
	case TOKEN_EOF:
		return "TOKEN_EOF"

	case TOKEN_ERROR:
		return "TOKEN_ERROR"

	case TOKEN_KEYWORD:
		return "TOKEN_KEYWORD"

	case TOKEN_IDENTIFIER:
		return "TOKEN_IDENTIFIER"

	case TOKEN_COMMA:
		return "TOKEN_COMMA"

	case TOKEN_SEMICOLON:
		return "TOKEN_SEMICOLON"

	case TOKEN_NULLABLE:
		return "TOKEN_NULLABLE"

	case TOKEN_CURLY_OPEN:
		return "TOKEN_CURLY_OPEN"

	case TOKEN_CURLY_CLOSE:
		return "TOKEN_CURLY_CLOSE"

	case TOKEN_ROUND_OPEN:
		return "TOKEN_ROUND_OPEN"

	case TOKEN_ROUND_CLOSE:
		return "TOKEN_ROUND_CLOSE"

	default:
		return "TOKEN_UNKNOWN"
	}
}

func GenerateTestTokens(data []byte) {
	lexer := CreateLexer(data)

	fmt.Println("expectedTokens := makeTestTokens(")
	for {
		token := NextToken(&lexer)

		typeString := TokenTypeToString(token.tokenType)
		fmt.Printf("\ttestToken(%v, \"%v\", %v, %v, %v),\n", typeString, token.tokenValue.string, token.tokenValue.int, token.line.number, token.line.offset)

		if IsType(token, TOKEN_EOF) || IsType(token, TOKEN_ERROR) {
			break
		}
	}
	fmt.Println(")")
}
