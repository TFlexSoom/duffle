// author: Tristan Hilbert
// date: 8/29/2023
// filename: lfunGrammar.go
// desc: Parsing Grammar to Build AST for lfun files
package parsing

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// // Lexer
func getLDatLexer() (*lexer.StatefulDefinition, error) {
	return lexer.NewSimple([]lexer.SimpleRule{
		{Name: "WHITESPACE", Pattern: `[ \t]+`},
		{Name: "EOL", Pattern: `(\r)?\n`},
		{Name: "ARRAY_START", Pattern: `\[`},
		{Name: "ARRAY_END", Pattern: `\]`},
		{Name: "OBJECT_START", Pattern: `\(`},
		{Name: "OBJECT_END", Pattern: `\)`},
		{Name: "ITEM_SEP", Pattern: `,`},
		{Name: "DECIMAL", Pattern: `[\d]+\.[\d]*`},
		{Name: "INTEGER", Pattern: `[\d]+`},
		{Name: "QUOTED_VAL", Pattern: `"[~"]*"`},
		{Name: "IDENTIFIER", Pattern: `[a-zA-Z][a-zA-Z\d]*`},
	})
}

func GetLDatParser() (*participle.Parser[Configuration], error) {
	var lexer, err = getLDatLexer()
	if err != nil {
		return nil, err
	}

	return participle.Build[Configuration](
		participle.Lexer(lexer),
		participle.Elide("WHITESPACE"),
		participle.UseLookahead(1),
	)
}

//// Grammar

type Configuration struct {
	Pos lexer.Position

	Assignments []Assignment `@@*`
}

type Assignment struct {
	Pos lexer.Position

	Name  string          `@IDENTITY`
	Value AssignmentValue `= @@`
}

type AssignmentValue interface {
	value()
	pos() lexer.Position
}

type Float struct {
	Pos   lexer.Position
	Value float64 `@DECIMAL`
}

func (f Float) value() {}
func (f Float) pos() lexer.Position {
	return f.Pos
}

type Int struct {
	Pos   lexer.Position
	Value int `@INTEGER`
}

func (f Int) value() {}
func (f Int) pos() lexer.Position {
	return f.Pos
}

type String struct {
	Pos   lexer.Position
	Value string `@IDENTITY | @QUOTED_VAL`
}

func (f String) value() {}
func (f String) pos() lexer.Position {
	return f.Pos
}

type List struct {
	Pos    lexer.Position
	Values []AssignmentValue `"[" (@@ ",")* @@? "]"`
}

func (l List) value() {}
func (l List) pos() lexer.Position {
	return l.Pos
}

type Struct struct {
	Pos    lexer.Position
	Fields []AssignmentValue `"(" (@@ ",")* @@? ")"`
}

func (s Struct) value() {}
func (s Struct) pos() lexer.Position {
	return s.Pos
}
