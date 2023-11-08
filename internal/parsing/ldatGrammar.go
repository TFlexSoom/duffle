// author: Tristan Hilbert
// date: 8/29/2023
// filename: ldatGrammar.go
// desc: Parsing Grammar to Build AST for ldat files
// notes: Might trade this out for a TOML parser instead!
package parsing

import (
	"io"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/deflemma/internal/types"
)

// // Lexer
func getLDatLexer() (*lexer.StatefulDefinition, error) {
	return lexer.NewSimple([]lexer.SimpleRule{
		{Name: "WHITESPACE", Pattern: `[ \t]+`},
		{Name: "EOL", Pattern: `\r?\n`},
		{Name: "ASSIGNMENT_OP", Pattern: `=`},
		{Name: "ARRAY_START", Pattern: `\[`},
		{Name: "ARRAY_END", Pattern: `\]`},
		{Name: "OBJECT_START", Pattern: `\(`},
		{Name: "OBJECT_END", Pattern: `\)`},
		{Name: "ITEM_SEP", Pattern: `,`},
		{Name: "DECIMAL", Pattern: `[\d]+\.[\d]*`},
		{Name: "INTEGER", Pattern: `[\d]+`},
		{Name: "QUOTED_VAL", Pattern: `"[~"]*"`},
		{Name: "BOOL_VALUE", Pattern: `true|false`},
		{Name: "IDENTIFIER", Pattern: `[a-zA-Z][a-zA-Z\d]*`},
		{Name: "TEXT", Pattern: `[^\w\r\n]+`},
	})
}

type ConfigurationParser struct {
	Parser *participle.Parser[Configuration]
}

func (configParser *ConfigurationParser) ParseSourceFile(
	fileName string,
	reader io.Reader,
) (interface{}, error) {
	return configParser.Parser.Parse(fileName, reader)
}

func GetLDatParser() (types.SourceFileParser, error) {
	var lexer, err = getLDatLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Configuration](
		participle.Lexer(lexer),
		participle.UseLookahead(1),
		participle.Union[AssignmentValue](
			List{},
			Struct{},
			Bool{},
			Float{},
			Int{},
			String{},
		),
	)

	if err != nil {
		return nil, err
	}

	wrapped := ConfigurationParser{
		Parser: parser,
	}

	return &wrapped, nil
}

//// Grammar

type Configuration struct {
	Pos lexer.Position

	Assignments []Assignment `(@@ EOL*)*`
}

type Assignment struct {
	Pos lexer.Position

	Name  string          `@IDENTIFIER`
	Value AssignmentValue `WHITESPACE* "=" WHITESPACE* @@`
}

type AssignmentValue interface {
	value()
	pos() lexer.Position
}

type Boolean bool

func (boolean *Boolean) Capture(values []string) error {
	*boolean = values[0] == "true"
	return nil
}

type Bool struct {
	Pos   lexer.Position
	Value Boolean `@( "true" | "false" ) WHITESPACE* EOL`
}

func (b Bool) value() {}
func (b Bool) pos() lexer.Position {
	return b.Pos
}

type Float struct {
	Pos   lexer.Position
	Value float64 `@DECIMAL WHITESPACE* EOL`
}

func (f Float) value() {}
func (f Float) pos() lexer.Position {
	return f.Pos
}

type Int struct {
	Pos   lexer.Position
	Value int `@INTEGER EOL`
}

func (f Int) value() {}
func (f Int) pos() lexer.Position {
	return f.Pos
}

type StringValue string

func (stringVal *StringValue) Capture(values []string) error {
	if stringVal == nil {
		*stringVal = ""
	}

	*stringVal += StringValue(strings.Join(values, ""))
	return nil
}

type String struct {
	Pos   lexer.Position
	Value StringValue `(@IDENTIFIER | @QUOTED_VAL | @TEXT) (@WHITESPACE* (@IDENTIFIER | @TEXT))* WHITESPACE* EOL`
}

func (f String) value() {}
func (f String) pos() lexer.Position {
	return f.Pos
}

type List struct {
	Pos    lexer.Position
	Values []AssignmentValue `"[" WHITESPACE* EOL? (WHITESPACE* @@ "," EOL?)* WHITESPACE* @@? WHITESPACE* EOL?"]" WHITESPACE* EOL`
}

func (l List) value() {}
func (l List) pos() lexer.Position {
	return l.Pos
}

type Struct struct {
	Pos    lexer.Position
	Fields []AssignmentValue `"(" WHITESPACE* EOL? (WHITESPACE* @@ "," EOL?)* WHITESPACE* @@? WHITESPACE* EOL? ")" WHITESPACE* EOL`
}

func (s Struct) value() {}
func (s Struct) pos() lexer.Position {
	return s.Pos
}
