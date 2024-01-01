package generator

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/files"
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

func GetDdatParser() (files.SourceFileParser, error) {
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

package ddatgrammar

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/container"
	"github.com/tflexsoom/duffle/internal/intermediate"
	"github.com/tflexsoom/duffle/internal/parsing/util"
)

type DuffleDataValue interface {
	DuffleValue() container.Tree[intermediate.DataValue]
	Pos() lexer.Position
	IsGroup() bool
}

// -- Boolean Grammar
func (b util.BoolGrammar) DuffleValue() container.Tree[intermediate.DataValue] {
	return container.NewGraphTreeCap[intermediate.DataValue](1, 1).AddChild(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_BOOLEAN,
			TextValue: b.Val,
		})

}

func (b util.BoolGrammar) IsGroup() bool {
	return false
}

// -- Float Grammar
func (f util.FloatGrammar) DuffleValue() container.Tree[intermediate.DataValue] {
	return container.NewGraphTreeCap[intermediate.DataValue](1, 1).AddChild(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_DECIMAL,
			TextValue: f.Val,
		})

}

func (f util.FloatGrammar) IsGroup() bool {
	return false
}

// -- Int Grammar

func (f util.IntGrammar) DuffleValue() container.Tree[intermediate.DataValue] {
	return container.NewGraphTreeCap[intermediate.DataValue](1, 1).AddChild(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_INTEGER,
			TextValue: f.Val,
		})
}

func (f util.IntGrammar) IsGroup() bool {
	return false
}

// -- String Grammar

func (f StringGrammar) DuffleValue() container.Tree[intermediate.DataValue] {
	return container.NewGraphTreeCap[intermediate.DataValue](1, 1).AddChild(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_TEXT,
			TextValue: f.Val,
		})
}

func (f StringGrammar) IsGroup() bool {
	return false
}

// -- List Grammar

func (l List) DuffleValue() container.Tree[intermediate.DataValue] {
	result := container.NewGraphTreeCap[intermediate.DataValue](2, uint(len(l.Vals)))
	result.SetValue(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_LIST,
			TextValue: "",
		},
	)

	for _, val := range l.Vals {
		if val.IsGroup() {
			container.AddChildren(result, (val.DuffleValue()))
		} else {
			result.AddChild(val.DuffleValue().GetValue())
		}
	}

	return result
}

func (f List) IsGroup() bool {
	return true
}

// -- Struct Grammar

func (s Struct) DuffleValue() container.Tree[intermediate.DataValue] {
	result := container.NewGraphTreeCap[intermediate.DataValue](2, uint(len(s.Vals)))
	result.SetValue(
		intermediate.DataValue{
			Type:      intermediate.TYPEID_STRUCT,
			TextValue: "",
		},
	)

	for _, val := range s.Vals {
		if val.IsGroup() {
			container.AddChildren(result, (val.DuffleValue()))
		} else {
			result.AddChild(val.DuffleValue().GetValue())
		}
	}

	return result
}

func (s Struct) IsGroup() bool {
	return true
}
