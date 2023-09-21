// author: Tristan Hilbert
// date: 8/29/2023
// filename: lmemGrammar.go
// desc: Parsing Grammar to Build AST for lmem files
package parsing

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/tflexsoom/deflemma/internal/types"
)

// // Lexer
func getLMemLexer() (*lexer.StatefulDefinition, error) {
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

type StructureParser struct {
	Parser *participle.Parser[Structure]
}

func (structureParser *StructureParser) ParseSourceFile(
	fileName string,
	reader io.Reader,
) (interface{}, error) {
	return structureParser.Parser.Parse(fileName, reader)
}

func GetLMemParser() (types.SourceFileParser, error) {
	var lexer, err = getLMemLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Structure](
		participle.Lexer(lexer),
		participle.Elide("WHITESPACE"),
		participle.UseLookahead(1),
	)
	if err != nil {
		return nil, err
	}

	wrapped := StructureParser{
		Parser: parser,
	}

	return &wrapped, nil
}

//// Grammar

type Structure struct {
	Pos lexer.Position
}
