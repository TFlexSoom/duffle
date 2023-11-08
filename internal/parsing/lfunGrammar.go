// author: Tristan Hilbert
// date: 8/29/2023
// filename: lfunGrammar.go
// desc: Parsing Grammar to Build AST for lfun files
package parsing

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/deflemma/internal/types"
)

// // Lexer
const identifierRegexPattern = `[a-zA-Z][a-zA-Z\d]*`

var identifierRegex = regexp.MustCompile(identifierRegexPattern)

func getLFunLexer() (*lexer.StatefulDefinition, error) {
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
			{Name: "DECIMAL", Pattern: `[\d]+\.[\d]*`, Action: nil},
			{Name: "INTEGER", Pattern: `[\d]+`, Action: nil},
			{Name: "SINGLE_QUOTED_VAL", Pattern: `'[~']*'`, Action: nil}, // Escape quotes?
			{Name: "QUOTED_VAL", Pattern: `"[~"]*"`, Action: nil},        // Escape quotes?
			{Name: "BOOL_VALUE", Pattern: `true|false`, Action: nil},
		},
		"Expression": {
			lexer.Include("Spacing"),
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

func GetLFunParser() (types.SourceFileParser, error) {
	var lexer, err = getLFunLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Module](
		participle.Lexer(lexer),
		participle.Elide("WHITESPACE"),
		participle.UseLookahead(1),
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
			CaptureExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
			OperatorExpression{},
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
		participle.Union[LiteralValue](
			ExprBool{},
			ExprFloat{},
			ExprInt{},
			ExprString{},
			ExprChar{},
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

type ConfigInstruction interface {
	config()
	pos() lexer.Position
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
	Type         Type               `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "->" ) ) )?`
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
func (expression ParentheticalExpression) config() {}
func (expression ParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type CaptureExpression struct {
	Pos lexer.Position

	Execution     InlineInstruction `BACKTICK @@ BACKTICK`
	NextExecution InlineInstruction `@@?`
}

func (expression CaptureExpression) block()  {}
func (expression CaptureExpression) inline() {}
func (expression CaptureExpression) config() {}
func (expression CaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@IDENTIFIER+`
	NextExecution  InlineInstruction `@@?`
}

func (expression ReferenceExpression) block()  {}
func (expression ReferenceExpression) inline() {}
func (expression ReferenceExpression) config() {}
func (expression ReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type OperatorExpression struct {
	Pos lexer.Position

	ReferenceGroup []string          `@OPERATOR`
	NextExecution  InlineInstruction `@@?`
}

func (expression OperatorExpression) inline() {}
func (expression OperatorExpression) config() {}
func (expression OperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type LiteralExpression struct {
	Pos lexer.Position

	Value LiteralValue `@@`
}

func (expression LiteralExpression) config() {}
func (expression LiteralExpression) pos() lexer.Position {
	return expression.Pos
}

type LiteralValue interface {
	literal()
	pos() lexer.Position
}

type ExprBoolean bool

func (boolean *ExprBoolean) Capture(values []string) error {
	*boolean = values[0] == "true"
	return nil
}

type ExprBool struct {
	Pos   lexer.Position
	Value Boolean `@( "true" | "false" )`
}

func (b ExprBool) literal() {}
func (b ExprBool) pos() lexer.Position {
	return b.Pos
}

type ExprFloat struct {
	Pos   lexer.Position
	Value float64 `@DECIMAL`
}

func (f ExprFloat) literal() {}
func (f ExprFloat) pos() lexer.Position {
	return f.Pos
}

type ExprInt struct {
	Pos   lexer.Position
	Value int `@INTEGER`
}

func (f ExprInt) literal() {}
func (f ExprInt) pos() lexer.Position {
	return f.Pos
}

type ExprStringValue string

func (stringVal *ExprStringValue) Capture(values []string) error {
	if stringVal == nil {
		*stringVal = ""
	}

	*stringVal += ExprStringValue(strings.Join(values, ""))
	return nil
}

type ExprString struct {
	Pos   lexer.Position
	Value ExprStringValue `@QUOTED_VAL`
}

func (f ExprString) literal() {}
func (f ExprString) pos() lexer.Position {
	return f.Pos
}

type ExprCharValue rune

func (stringVal *ExprCharValue) Capture(values []string) error {
	if len(values[0]) > 1 {
		return errors.New("char value is more than 1 character")
	} else if len(values[0]) < 1 {
		return errors.New("char values is less than 1 character")
	}

	*stringVal = ExprCharValue(rune(values[0][0]))

	return nil
}

type ExprChar struct {
	Pos   lexer.Position
	Value ExprCharValue `@SINGLE_QUOTED_VAL`
}

func (f ExprChar) literal() {}
func (f ExprChar) pos() lexer.Position {
	return f.Pos
}

type Pattern struct {
	Pos lexer.Position

	Name       string            `@IDENTIFIER`
	Params     []string          `@IDENTIFIER* "="`
	Definition InlineInstruction `@@`
}
