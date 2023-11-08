// author: Tristan Hilbert
// date: 8/29/2023
// filename: grammar.go
// desc: Grammar Utilities for All Parsers in Duffle
// notes: Might trade this out for a TOML parser instead!
package util

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

//// Grammar

type Value interface {
	Value()
	Pos() lexer.Position
}

const BooleanRegex = `true|false`

type Bool bool

func (boolean *Bool) Capture(values []string) error {
	*boolean = values[0] == "true"
	return nil
}

type BoolGrammar struct {
	Position lexer.Position
	Val      Bool `@( "true" | "false" )`
}

func (b BoolGrammar) Value() {}
func (b BoolGrammar) Pos() lexer.Position {
	return b.Position
}

const DecimalTagName = "DECIMAL"
const DecimalRegex = `[\d]+\.[\d]*`

type FloatGrammar struct {
	Position lexer.Position
	Val      float64 `@DECIMAL`
}

func (f FloatGrammar) Value() {}
func (f FloatGrammar) Pos() lexer.Position {
	return f.Position
}

const IntTagName = "INTEGER"
const IntRegex = `[\d]+`

type IntGrammar struct {
	Position lexer.Position
	Val      int `@INTEGER`
}

func (f IntGrammar) Value() {}
func (f IntGrammar) Pos() lexer.Position {
	return f.Position
}

const QuotedValTagName = "QUOTED_VAL"
const QuotedValRegex = `"[^"]*"`

type String string

func (stringVal *String) Capture(values []string) error {
	if stringVal == nil {
		*stringVal = ""
	}

	*stringVal += String(strings.Join(values, ""))
	return nil
}

type StringGrammar struct {
	Position lexer.Position
	Val      String `@QUOTED_VAL`
}

func (f StringGrammar) Value() {}
func (f StringGrammar) Pos() lexer.Position {
	return f.Position
}
