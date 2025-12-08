package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	lexer := Tokenize("SELECT name, age from users;")
	if len(lexer.Tokens) != 7 {
		t.Log(lexer)
		t.Fatal("The number of output token should be 7")
	}
	t.Logf("%v", lexer)
}
