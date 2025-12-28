package parser

type OperatorKind uint16

const (
	EQUALS OperatorKind = iota
	LESS_THAN
	LESS_THAN_OR_EQUAL
	GREATER_THAN
	GREATER_THAN_OR_EQUAL
)

type column struct {
	name        string
	alias       string
	targetTable string
}

type condition struct {
	column   string
	operator OperatorKind
	value    any
}

type joinTable struct {
	targetTable string
	alias       string
	condition   []condition
}

type SelectAst struct {
	columns     []column
	targetTable string
	alias       string
	joinTables  []joinTable
	WhereClause []condition
	limitOffset int
	limitCount  int
}

func parseSelectStmt(p *Parser) {

}
