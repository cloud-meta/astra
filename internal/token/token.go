package token

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	// 控制结构
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	IDENT   TokenType = "IDENT" // resource, variable, etc
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	BOOL    TokenType = "BOOL"

	// 分隔符
	ASSIGN TokenType = "="
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	LBRACK TokenType = "["
	RBRACK TokenType = "]"
	COMMA  TokenType = ","
	DOT    TokenType = "."
	COLON  TokenType = ":"

	// 运算符
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	QUESTION TokenType = "?"
	EQ       TokenType = "=="
	NEQ      TokenType = "!="
	LT       TokenType = "<"
	GT       TokenType = ">"
	LE       TokenType = "<="
	GE       TokenType = ">="

	// 关键字
	RESOURCE TokenType = "RESOURCE"
	WHEN     TokenType = "WHEN"
	DERIVE   TokenType = "DERIVE"
	MAP      TokenType = "MAP"
	POLICY   TokenType = "POLICY"
	ENFORCE  TokenType = "ENFORCE"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
)

var (
	silence  = struct{}{}
	keywords = map[string]struct{}{
		"resource": silence,
		"service":  silence,
		"model":    silence,
		"provider": silence,
		"extends":  silence,
		"func":     silence,
		"abstract": silence,
	}
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	length int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, line: 1, col: 0, length: len(input)}
}

func (l *Lexer) NextToken() Token {
	for l.pos < l.length {
		ch := l.input[l.pos]

		if ch == ' ' || ch == '\t' {
			l.pos++
			l.col++
			continue
		}

		if ch == '\n' {
			l.pos++
			l.line++
			l.col = 0
			continue
		}

		if unicode.IsLetter(rune(ch)) {
			start := l.pos
			startCol := l.col
			for l.pos < l.length && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos]))) {
				l.pos++
				l.col++
			}

			val := l.input[start:l.pos]
			if _, ok := keywords[val]; ok {
				return Token{Type: TokenKeyword, Literal: val, Line: l.line, Column: startCol}
			}

			return Token{Type: TokenIdentifier, Literal: val, Line: l.line, Column: startCol}
		}

		// 数字
		if unicode.IsDigit(rune(ch)) {
			start := l.pos
			startCol := l.col
			for l.pos < l.length && unicode.IsDigit(rune(l.input[l.pos])) {
				l.pos++
				l.col++
			}
			return Token{Type: TokenNumber, Literal: l.input[start:l.pos], Line: l.line, Column: startCol}
		}

		// 符号
		switch ch {
		case '{', '}', ':', '=', '[', ']', ',':
			l.pos++
			l.col++
			return Token{Type: TokenSymbol, Literal: string(ch), Line: l.line, Column: l.col - 1}
		case '"':
			startCol := l.col
			l.pos++
			l.col++
			start := l.pos
			for l.pos < l.length && l.input[l.pos] != '"' {
				l.pos++
				l.col++
			}
			val := l.input[start:l.pos]
			l.pos++
			l.col++
			return Token{Type: TokenString, Literal: val, Line: l.line, Column: startCol}
		default:
			panic(fmt.Sprintf("Unexpected char: %c", ch))
		}
	}

	return Token{Type: TokenEOF, Literal: "", Line: l.line, Column: l.col}
}
