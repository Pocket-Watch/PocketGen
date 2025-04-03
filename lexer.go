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
	TOKEN_CURLY_OPEN
	TOKEN_CURLY_CLOSE
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
	KEYWORD_CONST KeywordType = "const"
	KEYWORD_FUNC  KeywordType = "func"
)

var KEYWORDS = []KeywordType{KEYWORD_CONST, KEYWORD_FUNC}

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

type LinePos struct {
	number int
	offset int
}

type Lexer struct {
	data []byte
	line LinePos
	pos  int

	// TODO(kihau):
	//     Store previous and current token to allow implementation of Peek() and PeekNext() functions.
	//     This could also be a token ring buffer when larger token sequence is needed for parsing.
	// prevToken Token
	// currToken Token
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

	return lexer
}

func nextRune(lexer *Lexer) rune {
	if lexer.pos >= len(lexer.data) {
		return 0
	}

	rune, runeSize := utf8.DecodeRune(lexer.data[lexer.pos:])
	if rune == utf8.RuneError {
		return rune
	}

	lexer.pos += runeSize

	if rune == '\n' {
		lexer.line.number += 1
		lexer.line.offset = 1
	} else {
		lexer.line.offset += 1
	}

	return rune
}

func peekRune(lexer *Lexer) rune {
	if lexer.pos >= len(lexer.data) {
		return 0
	}

	rune, _ := utf8.DecodeRune(lexer.data[lexer.pos:])
	return rune
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
	for rune != '\n' {
		nextRune(lexer)
		rune = peekRune(lexer)
	}
}

func IsType(token Token, tokenType TokenType) bool {
	return token.tokenType == tokenType
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

		case ',':
			nextRune(lexer)
			return makeToken(TOKEN_COMMA, line)

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

func PrintToken(token Token) {
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

	case TOKEN_CURLY_OPEN:
		name = "CURLY OPEN"
		value = "{"

	case TOKEN_CURLY_CLOSE:
		name = "CURLY CLOSE"
		value = "}"

	default:
		name = "<UNKNOWN TOKEN>"
		value = fmt.Sprintf("%v", token.tokenType)
	}

	line := fmt.Sprintf("%v:%v ", token.line.number, token.line.offset)
	fmt.Printf("%-14s %-6s - %s\n", name, line, value)
}
