package parser

import (
	"go-db/core/query/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	tokens := lexer.Tokenize("SELECT name, age from users;")
	CreateAST(tokens)
}
