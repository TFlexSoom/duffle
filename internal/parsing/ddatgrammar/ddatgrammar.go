// author: Tristan Hilbert
// date: 8/29/2023
// filename: ldatGrammar.go
// desc: Parsing Grammar to Build AST for ldat files
// notes: Might trade this out for a TOML parser instead!
package ddatgrammar

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/deflemma/internal/parsing/util"
	"github.com/tflexsoom/deflemma/internal/types"
)

// // Lexer
func getDdatLexer() (*lexer.StatefulDefinition, error) {
	return lexer.NewSimple([]lexer.SimpleRule{
		{Name: "WHITESPACE", Pattern: `[ \t]+`},
		{Name: "EOL", Pattern: `\r?\n`},
		{Name: "ASSIGNMENT_OP", Pattern: `=`},
		{Name: "ARRAY_START", Pattern: `\[`},
		{Name: "ARRAY_END", Pattern: `\]`},
		{Name: "OBJECT_START", Pattern: `\(`},
		{Name: "OBJECT_END", Pattern: `\)`},
		{Name: "ITEM_SEP", Pattern: `,`},
		{Name: util.DecimalTagName, Pattern: util.DecimalRegex},
		{Name: util.IntTagName, Pattern: util.IntRegex},
		{Name: util.QuotedValTagName, Pattern: util.QuotedValRegex},
		{Name: "BOOLEAN", Pattern: util.BooleanRegex},
		{Name: "DOT_OPERATOR", Pattern: `\.`},
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

func GetDdatParser() (types.SourceFileParser, error) {
	var lexer, err = getDdatLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Configuration](
		participle.Lexer(lexer),
		participle.UseLookahead(1),
		participle.Union[util.Value](
			List{},
			Struct{},
			util.BoolGrammar{},
			util.FloatGrammar{},
			util.IntGrammar{},
			StringGrammar{},
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

	FirstName  *string    `@IDENTIFIER`
	SecondName string     `("." @IDENTIFIER )?`
	Value      util.Value `WHITESPACE* "=" WHITESPACE* @@ WHITESPACE* EOL`
}

type StringGrammar struct {
	Position lexer.Position
	Val      util.String `(@IDENTIFIER | @QUOTED_VAL | @TEXT) (@WHITESPACE* (@IDENTIFIER | @TEXT))*`
}

func (f StringGrammar) Value() {}
func (f StringGrammar) Pos() lexer.Position {
	return f.Position
}

type List struct {
	Position lexer.Position
	Values   []util.Value `"[" WHITESPACE* EOL? WHITESPACE* @@? ("," EOL? WHITESPACE* @@)* WHITESPACE* EOL? WHITESPACE*"]"`
}

func (l List) Value() {}
func (l List) Pos() lexer.Position {
	return l.Position
}

type Struct struct {
	Position lexer.Position
	Fields   []util.Value `"(" WHITESPACE* EOL? WHITESPACE* @@? ("," EOL? WHITESPACE* @@)* WHITESPACE* EOL? WHITESPACE* ")"`
}

func (s Struct) Value() {}
func (s Struct) Pos() lexer.Position {
	return s.Position
}
