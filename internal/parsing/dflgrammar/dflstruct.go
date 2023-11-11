// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflstruct.go
// desc: Duffle struct keyword parsing grammar
package dflgrammar

import "github.com/alecthomas/participle/v2/lexer"

type StructModulePart struct {
	Position lexer.Position

	Name   string  `STRUCT_KEYWORD @IDENTIFIER`
	Fields []Input `"(" EOL ( "<" @@ ">" EOL )+ ")" EOL+`
}

func (modPart StructModulePart) ModulePart() {}
func (modPart StructModulePart) Pos() lexer.Position {
	return modPart.Position
}
