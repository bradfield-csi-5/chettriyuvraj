package scanner

type Type int
type Lexeme int
type Literal int

type Token struct {
	Type 	Type
	Lexeme  Lexeme
	Literal Literal
}

const (
	/* Logical operators */
	AND Type = iota
	OR
	NOT

	EOF
)

const (
	STRING Literal = iota + 1
	INTEGER
	IDENTIFIER
)

/* One Type corresponds to only a single Lexeme */
var TypeAsLexeme [...]string = {
	AND: "AND",
	OR: "OR",
	NOT: "NOT",
}

var LiteralAsStr [...]string = {
	STRING: "STRING"
	INTEGER: "INTEGER"
}

var LexemeAsType map[string]Token = {
	"AND": AND,
	"OR": OR,
	"NOT": NOT
}

