package parser

import (
	"go-db/core/query/lexer"
)

type Parser struct {
	tokens lexer.Tokens
	pos    int
}

func parseTokens(p *Parser) {

}

func CreateAST(tokens lexer.Tokens) {
	parser := Parser{
		tokens: tokens,
		pos:    0,
	}
	parseTokens(&parser)
}
