// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflgrammar.go
// desc: Parsing Grammar to Build AST for dfl files
package dflgrammar

import (
	"errors"
	"io"
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/deflemma/internal/types"

	"github.com/tflexsoom/deflemma/internal/parsing/util"
)

// // Lexer
const identifierRegexPattern = `[a-zA-Z][a-zA-Z\d_]*`

var identifierRegex = regexp.MustCompile(identifierRegexPattern)

func getDflLexer() (*lexer.StatefulDefinition, error) {
	return lexer.New(lexer.Rules{
		"Spacing": {
			{Name: "EOL", Pattern: `(\r)?\n`, Action: nil},
			{Name: "WHITESPACE", Pattern: `[ \t]+`, Action: nil},
		},
		"Identity": {
			{Name: "IDENTIFIER", Pattern: identifierRegexPattern, Action: nil},
		},
		"Operator": {
			{Name: "OPERATOR", Pattern: `[^\d\w][^\w]*`, Action: nil},
		},
		"Literal": {
			{Name: "BOOLEAN", Pattern: util.BooleanRegex, Action: nil},
			{Name: util.DecimalTagName, Pattern: util.DecimalRegex, Action: nil},
			{Name: util.IntTagName, Pattern: util.IntRegex, Action: nil},
			{Name: "SINGLE_QUOTED_VAL", Pattern: `'[^']*'`, Action: nil},             // Escape quotes?
			{Name: util.QuotedValTagName, Pattern: util.QuotedValRegex, Action: nil}, // Escape quotes?
		},
		"Expression": {
			{Name: "BACKTICK", Pattern: "`", Action: nil},
			{Name: "EXPR_PUNCTATION", Pattern: `[();]`, Action: nil},
			{Name: "ANNOTATION_SYMBOL", Pattern: `[@]`, Action: nil},
		},
		"Root": {
			lexer.Include("Spacing"),
			lexer.Include("Expression"),
			{Name: "USE_KEYWORD", Pattern: `use`, Action: nil},
			{Name: "FACT_KEYWORD", Pattern: `fact`, Action: nil},
			{Name: "THEORY_KEYWORD", Pattern: `theory`, Action: nil},
			{Name: "PARAM_PUNCTATION", Pattern: `[\[\],<>]`, Action: nil},
			{Name: "BEGIN_KEYWORD", Pattern: `begin`, Action: lexer.Push("Instruction")},
			{Name: "EVALS_KEYWORD", Pattern: `evals`, Action: lexer.Push("Pattern")},
			lexer.Include("Literal"),
			lexer.Include("Identity"),
		},
		"Instruction": {
			lexer.Include("Spacing"),
			lexer.Include("Expression"),
			{Name: "INLINE_IF_KEYWORD", Pattern: `ifthen`, Action: nil},
			{Name: "IF_KEYWORD", Pattern: `if`, Action: nil},
			{Name: "THEN_KEYWORD", Pattern: `then`, Action: lexer.Push("Condition")},
			{Name: "END_KEYWORD", Pattern: `end`, Action: lexer.Pop()},
			lexer.Include("Identity"),
			lexer.Include("Operator"),
		},
		"Condition": {
			lexer.Include("Spacing"),
			{Name: "ELSEIF_KEYWORD", Pattern: `elseif`, Action: lexer.Pop()},
			{Name: "ELSE_KEYWORD", Pattern: `else`, Action: nil},
			{Name: "END_IF_KEYWORD", Pattern: `endif`, Action: lexer.Pop()},
			lexer.Include("Instruction"),
		},
		"Pattern": {
			{Name: "END_EVAL", Pattern: `((\r)?\n)[2]`, Action: lexer.Pop()},
			lexer.Include("Spacing"),
			lexer.Include("Expression"),
			lexer.Include("Identity"),
			lexer.Include("Operator"),
		},
	})
}

type ModuleParser struct {
	Parser *participle.Parser[Module]
}

func (modParser *ModuleParser) ParseSourceFile(
	fileName string,
	reader io.Reader,
) (interface{}, error) {
	return modParser.Parser.Parse(fileName, reader)
}

func GetDflParser() (types.SourceFileParser, error) {
	var lexer, err = getDflLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Module](
		participle.Lexer(lexer),
		participle.Elide("WHITESPACE"),
		participle.UseLookahead(2),
		participle.Union[ModulePart](
			ImportModulePart{},
			ConfigModulePart{},
			FunctionModulePart{},
		),
		participle.Union[Import](
			ListImport{},
			SingleImport{},
		),
		participle.Union[ConfigInstruction](
			ConfCaptureExpression{},
			ConfParentheticalExpression{},
			ConfReferenceExpression{},
			ConfOperatorExpression{},
			LiteralExpression{},
		),
		participle.Union[BlockInstruction](
			BlockConditionalExpression{},
			CaptureExpression{},
			InlineConditionalExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
		),
		participle.Union[InlineInstruction](
			CaptureExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
			OperatorExpression{},
		),
		participle.Union[util.Value](
			util.BoolGrammar{},
			util.FloatGrammar{},
			util.IntGrammar{},
			util.StringGrammar{},
			CharGrammar{},
		),
	)

	if err != nil {
		return nil, err
	}

	wrapped := ModuleParser{
		Parser: parser,
	}

	return &wrapped, nil
}

//// Grammar

type Module struct {
	Pos lexer.Position

	ModuleParts []ModulePart `@@*`
}

type ModulePart interface {
	modulePart()
	pos() lexer.Position
}

type ImportModulePart struct {
	Pos lexer.Position

	Imports []Import `( USE_KEYWORD @@ EOL+ )+`
}

func (modPart ImportModulePart) modulePart() {}
func (modPart ImportModulePart) pos() lexer.Position {
	return modPart.Pos
}

type Import interface {
	importVal()
	pos() lexer.Position
}

type ListImport struct {
	Pos lexer.Position

	Value []string `"(" EOL+ (@IDENTIFIER EOL+)+ ")"`
}

func (listImport ListImport) importVal() {}
func (listImport ListImport) pos() lexer.Position {
	return listImport.Pos
}

type SingleImport struct {
	Pos lexer.Position

	Value []string `@IDENTIFIER`
}

func (singleImport SingleImport) importVal() {}
func (singleImport SingleImport) pos() lexer.Position {
	return singleImport.Pos
}

type Uniqueness bool

func (u *Uniqueness) Capture(values []string) error {
	*u = values[0] == "theory"
	return nil
}

type ConfigModulePart struct {
	Pos lexer.Position

	Configs []*Config `( @@ EOL+ )+`
}

func (modPart ConfigModulePart) modulePart() {}
func (modPart ConfigModulePart) pos() lexer.Position {
	return modPart.Pos
}

type Config struct {
	Pos      lexer.Position
	IsUnique Uniqueness        `@( THEORY_KEYWORD | FACT_KEYWORD )`
	Input    string            `@IDENTIFIER`
	Value    ConfigInstruction `@@`
}

type FunctionModulePart struct {
	Pos lexer.Position

	Functions []Function `( @@ EOL+ )+`
}

func (modPart FunctionModulePart) modulePart() {}
func (modPart FunctionModulePart) pos() lexer.Position {
	return modPart.Pos
}

type Input struct {
	Pos lexer.Position

	Type Type   `@@`
	Name string `@IDENTIFIER`
}

type Type struct {
	Name     string `@IDENTIFIER`
	Generics []Type `("[" @@ ("," @@)* "]")?`
}

type FunctionName struct {
	Name       string
	IsOperator bool
}

func (fname *FunctionName) Capture(values []string) error {
	fname.Name = values[0]
	fname.IsOperator = !identifierRegex.MatchString(values[0])
	return nil
}

type Function struct {
	Pos lexer.Position

	Annotations  []string           `( "@" @IDENTIFIER )*`
	Type         Type               `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "<" ) ) )?`
	Name         FunctionName       `( @IDENTIFIER | @OPERATOR )`
	Inputs       []Input            `( "<" @@ ">")*`
	Instructions []BlockInstruction `( BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD )`
	Patterns     []Pattern          `| ( EVALS_KEYWORD EOL ( @@ EOL )* END_EVAL )`
}

type BlockInstruction interface {
	block()
	pos() lexer.Position
}

type InlineInstruction interface {
	inline()
	pos() lexer.Position
}

type ConfigInstruction interface {
	config()
	pos() lexer.Position
}

type BlockConditionalExpression struct {
	Pos lexer.Position

	Condition      InlineInstruction     `IF_KEYWORD "(" @@ ")"`
	Execution      []InlineInstruction   `THEN_KEYWORD EOL+ (@@ EOL*)*`
	SubConditional []SubBlockConditional `@@*`
	Alternative    []InlineInstruction   `(ELSE_KEYWORD EOL+ (@@ EOL+)*)? END_IF_KEYWORD`
}

func (expression BlockConditionalExpression) block() {}
func (expression BlockConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type SubBlockConditional struct {
	Pos lexer.Position

	Condition InlineInstruction   `ELSEIF_KEYWORD "(" @@ ")"`
	Execution []InlineInstruction `THEN_KEYWORD EOL+ (@@ EOL+)*`
}

type InlineConditionalExpression struct {
	Pos lexer.Position

	Condition     InlineInstruction `INLINE_IF_KEYWORD "(" @@ ")"`
	NextExecution InlineInstruction `@@`
}

func (expression InlineConditionalExpression) block() {}
func (expression InlineConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type ParentheticalExpression struct {
	Pos lexer.Position

	Execution     InlineInstruction `"(" @@ ")"`
	NextExecution InlineInstruction `@@?`
}

func (expression ParentheticalExpression) block()  {}
func (expression ParentheticalExpression) inline() {}
func (expression ParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type ConfParentheticalExpression struct {
	Pos lexer.Position

	Execution     ConfigInstruction `"(" @@ ")"`
	NextExecution ConfigInstruction `@@?`
}

func (expression ConfParentheticalExpression) config() {}
func (expression ConfParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type CaptureExpression struct {
	Pos lexer.Position

	Execution     InlineInstruction `BACKTICK @@ BACKTICK`
	NextExecution InlineInstruction `@@?`
}

func (expression CaptureExpression) block()  {}
func (expression CaptureExpression) inline() {}
func (expression CaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ConfCaptureExpression struct {
	Pos lexer.Position

	Execution     ConfigInstruction `BACKTICK @@ BACKTICK`
	NextExecution ConfigInstruction `@@?`
}

func (expression ConfCaptureExpression) config() {}
func (expression ConfCaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@IDENTIFIER+`
	NextExecution  InlineInstruction `@@?`
}

func (expression ReferenceExpression) block()  {}
func (expression ReferenceExpression) inline() {}
func (expression ReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type ConfReferenceExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@IDENTIFIER+`
	NextExecution  ConfigInstruction `@@?`
}

func (expression ConfReferenceExpression) config() {}
func (expression ConfReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type OperatorExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@OPERATOR`
	NextExecution  InlineInstruction `@@?`
}

func (expression OperatorExpression) inline() {}
func (expression OperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type ConfOperatorExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@OPERATOR`
	NextExecution  ConfigInstruction `@@?`
}

func (expression ConfOperatorExpression) config() {}
func (expression ConfOperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type LiteralExpression struct {
	Pos lexer.Position

	Value util.Value `@@`
}

func (expression LiteralExpression) config() {}
func (expression LiteralExpression) pos() lexer.Position {
	return expression.Pos
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

type Pattern struct {
	Pos lexer.Position

	Name       string            `@IDENTIFIER`
	Params     []string          `@IDENTIFIER* "="`
	Definition InlineInstruction `@@`
}
