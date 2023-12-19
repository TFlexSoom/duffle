// author: Tristan Hilbert
// date: 8/29/2023
// filename: grammar.go
// desc: Grammar Utilities for All Parsers in Duffle
// notes: Might trade this out for a TOML parser instead!
package util

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/container"
)

//// Grammar

type Value interface {
	Value() container.Tree[string]
	Pos() lexer.Position
	IsGroup() bool
}

const BooleanRegex = `true|false`

type BoolGrammar struct {
	Position lexer.Position
	Val      string `@( "true" | "false" )`
}

func (b BoolGrammar) Value() container.Tree[string] {
	return container.NewGraphTreeCap[string](1, 1).AddChild(b.Val)

}
func (b BoolGrammar) Pos() lexer.Position {
	return b.Position
}
func (b BoolGrammar) IsGroup() bool {
	return false
}

const DecimalTagName = "DECIMAL"
const DecimalRegex = `[\d]+\.[\d]*`

type FloatGrammar struct {
	Position lexer.Position
	Val      string `@DECIMAL`
}

func (f FloatGrammar) Value() container.Tree[string] {
	return container.NewGraphTreeCap[string](1, 1).AddChild(f.Val)
}
func (f FloatGrammar) Pos() lexer.Position {
	return f.Position
}
func (f FloatGrammar) IsGroup() bool {
	return false
}

const IntTagName = "INTEGER"
const IntRegex = `[\d]+`

type IntGrammar struct {
	Position lexer.Position
	Val      string `@INTEGER`
}

func (f IntGrammar) Value() container.Tree[string] {
	return container.NewGraphTreeCap[string](1, 1).AddChild(f.Val)
}
func (f IntGrammar) Pos() lexer.Position {
	return f.Position
}
func (f IntGrammar) IsGroup() bool {
	return false
}

const QuotedValTagName = "QUOTED_VAL"
const QuotedValRegex = `"[^"]*"`

type StringGrammar struct {
	Position lexer.Position
	Val      string `@QUOTED_VAL`
}

func (f StringGrammar) Value() container.Tree[string] {
	return container.NewGraphTreeCap[string](1, 1).AddChild(f.Val)
}

func (f StringGrammar) Pos() lexer.Position {
	return f.Position
}
func (f StringGrammar) IsGroup() bool {
	return false
}
