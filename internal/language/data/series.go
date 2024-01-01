package data

import "github.com/alecthomas/participle/v2/lexer"

type List struct {
	Position lexer.Position
	Vals     []DuffleDataValue `"[" WHITESPACE* EOL? WHITESPACE* @@? ("," EOL? WHITESPACE* @@)* WHITESPACE* EOL? WHITESPACE*"]"`
}

func (l List) Pos() lexer.Position {
	return l.Position
}

type Struct struct {
	Position lexer.Position
	Vals     []DuffleDataValue `"(" WHITESPACE* EOL? WHITESPACE* @@? ("," EOL? WHITESPACE* @@)* WHITESPACE* EOL? WHITESPACE* ")"`
}

func (s Struct) Pos() lexer.Position {
	return s.Position
}
