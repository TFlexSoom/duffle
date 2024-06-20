// author: Tristan Hilbert
// date: 10/27/2023
// filename: config.go
// desc: Configuration Grammar Components Where Assignments lead to Global Constant Symbols
package config

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/intermediate"
)

type Configuration struct {
	Pos lexer.Position

	Assignments []Assignment `(@@ EOL*)*`
}

type Assignment struct {
	Pos lexer.Position

	FirstName  string          `@IDENTIFIER`
	SecondName *string         `("." @IDENTIFIER )?`
	Value      DuffleDataValue `WHITESPACE* "=" WHITESPACE* @@ WHITESPACE* EOL`
}

func (a Assignment) GetDataConfig() intermediate.DataConfig {
	firstName := ""
	secondName := a.FirstName
	if a.SecondName != nil {
		firstName = secondName
		secondName = *a.SecondName
	}

	return intermediate.DataConfig{
		FirstName:  firstName,
		SecondName: secondName,
		Values:     a.Value.Value(),
	}
}
