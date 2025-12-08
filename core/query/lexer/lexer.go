package lexer

import (
	"strings"
	"unicode"
)

type TokenKind uint8

const (
	KEYWORD TokenKind = iota
	IDENTIFIER
	NUMBER
	STRING
	QUOTES
	COMMA
	SEMICOLON
	ASTERIK
	OPERATOR
)

type Keyword struct {
	Create string
	Where  string
	Select string
	Insert string
	Update string
}

var Keywords Keyword = Keyword{
	Create: "CREATE",
	Where:  "WHERE",
	Select: "SELECT",
	Insert: "INSERT",
	Update: "UPDATE",
}

type Token struct {
	Kind  TokenKind
	Value string
}

type Lexer struct {
	Tokens []Token
	pos    uint8
	len    int
	query  string
}

func (lexer *Lexer) nextToken() {
	str := ""
	for i, charRune := range lexer.query[lexer.pos:] {
		char := string(charRune)
		if unicode.IsSpace(charRune) && len(str) != 0 {
			lexer.addToken(str)
			lexer.advanceN(i + 1)
			break
		}
		if char == "," || char == `"` || char == `*` || char == `;` {
			lexer.addToken(str)
			lexer.advanceN(i + 1)
			lexer.addToken(char)
			break
		}
		str += char
	}
}

func (lexer *Lexer) addToken(str string) {
	formattedStr := strings.ToUpper(str)
	switch formattedStr {
	case Keywords.Create, Keywords.Where, Keywords.Insert, Keywords.Select, Keywords.Update:
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  KEYWORD,
				Value: formattedStr,
			})
		}
	case ",":
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  COMMA,
				Value: formattedStr,
			})
		}
	case `"`:
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  QUOTES,
				Value: formattedStr,
			})
		}
	case `*`:
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  ASTERIK,
				Value: formattedStr,
			})
		}
	case `;`:
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  SEMICOLON,
				Value: formattedStr,
			})
		}
	default:
		{
			lexer.Tokens = append(lexer.Tokens, Token{
				Kind:  IDENTIFIER,
				Value: str,
			})
		}
	}
}

func (lexer *Lexer) advanceN(n int) {
	lexer.pos += uint8(n)
}
