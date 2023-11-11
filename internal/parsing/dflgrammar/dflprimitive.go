// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflprimitive.go
// desc: Duffle Primitive Parsing Constructs
package dflgrammar

import (
	"errors"

	"github.com/alecthomas/participle/v2/lexer"
)

type Module struct {
	Pos lexer.Position

	ModuleParts []ModulePart `@@*`
}

type ModulePart interface {
	modulePart()
	pos() lexer.Position
}

type Type struct {
	Name     string `@IDENTIFIER`
	Generics []Type `("[" @@ ("," @@)* "]")?`
}

type Input struct {
	Pos lexer.Position

	Type Type   `@@`
	Name string `@IDENTIFIER`
}

type Char rune

func (charValue *Char) Capture(values []string) error {
	valLen := len(values[0])
	if valLen < 2 {
		return errors.New("char values is less than 1 character")
	}

	if valLen == 4 && values[0][1] == '\\' {
		switch values[0][1] {
		case '\'':
			*charValue = '\''
			return nil
		case '"':
			*charValue = '"'
			return nil
		case '\\':
			*charValue = '\\'
			return nil
		case 'a':
			*charValue = '\a'
			return nil
		case 'b':
			*charValue = '\b'
			return nil
		case 'f':
			*charValue = '\f'
			return nil
		case 'n':
			*charValue = '\n'
			return nil
		case 'r':
			*charValue = '\r'
			return nil
		case 't':
			*charValue = '\t'
			return nil
		case 'v':
			*charValue = '\v'
			return nil
		// TODO Maybe Include Hex Chars?
		// Prob best if those are hexidecimal numerics
		default:
			return errors.New("unrecognized escape character")
		}
	}

	if valLen > 3 {
		return errors.New("char value is more than 1 character")
	}

	*charValue = Char(rune(values[0][1]))

	return nil
}

type CharGrammar struct {
	Position lexer.Position
	Val      Char `@SINGLE_QUOTED_VAL`
}

func (f CharGrammar) Value() {}
func (f CharGrammar) Pos() lexer.Position {
	return f.Position
}
