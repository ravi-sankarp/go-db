package lexer

func Tokenize(query string) Lexer {
	lexer := Lexer{
		Tokens: make([]Token, 0, 20),
		len:    len(query),
		query:  query,
		pos:    0,
	}

	for lexer.pos < uint8(lexer.len) {
		lexer.nextToken()
	}
	return lexer
}
