package function

import "github.com/alecthomas/participle/v2/lexer"

type Module struct {
	Position lexer.Position

	ModuleParts []ModulePart `@@*`
}

type ModulePart interface {
	ModulePart()
	Pos() lexer.Position
}
