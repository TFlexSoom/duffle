// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflprimitive.go
// desc: Duffle Primitive Parsing Constructs
package dflgrammar

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/container"
	"github.com/tflexsoom/duffle/internal/intermediate"
)

type Module struct {
	Position lexer.Position

	ModuleParts []ModulePart `@@*`
}

type ModulePart interface {
	ModulePart()
	Pos() lexer.Position
}

type Type struct {
	Name     string `@IDENTIFIER`
	Generics []Type `("[" @@ ("," @@)* "]")?`
}

type Input struct {
	Position lexer.Position

	Type Type   `@@`
	Name string `@IDENTIFIER`
}

type Char struct {
	Position lexer.Position
	Val      string `@SINGLE_QUOTED_VAL`
}

func (f Char) Value() container.Tree[intermediate.DataValue] {
	return container.NewGraphTreeCap[intermediate.DataValue](1, 1).AddChild(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_CHAR,
			TextValue: f.Val,
		})
}
func (f Char) Pos() lexer.Position {
	return f.Position
}
func (f Char) IsGroup() bool {
	return false
}
