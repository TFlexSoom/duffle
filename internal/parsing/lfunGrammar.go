// author: Tristan Hilbert
// date: 8/29/2023
// filename: lfunGrammar.go
// desc: Parsing Grammar to Build AST for lfun files
package parsing

import (
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// // Lexer
const identifierRegexPattern = `[a-zA-Z][a-zA-Z\d]*`

var identifierRegex = regexp.MustCompile(identifierRegexPattern)

func getLexer() (*lexer.StatefulDefinition, error) {
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
		"Root": {
			lexer.Include("Spacing"),
			{Name: "USE_KEYWORD", Pattern: `use`, Action: nil},
			{Name: "ARROW", Pattern: `[>-][>]`, Action: nil},
			{Name: "TYPE_PUNCTATION", Pattern: `[@\[\](),]`, Action: nil},
			{Name: "BEGIN_KEYWORD", Pattern: `begin`, Action: lexer.Push("Instruction")},
			{Name: "EVALS_KEYWORD", Pattern: `evals`, Action: lexer.Push("Pattern")},
			lexer.Include("Identity"),
		},
		"Expression": {
			lexer.Include("Spacing"),
			{Name: "BACKTICK", Pattern: "`", Action: nil},
			{Name: "EXPR_PUNCTATION", Pattern: `[();]`, Action: nil},
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

func GetParser() (*participle.Parser[Module], error) {
	var lexer, err = getLexer()
	if err != nil {
		return nil, err
	}

	return participle.Build[Module](
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
		participle.Union[Instruction](
			BlockConditionalExpression{},
			InlineConditionalExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
		),
		participle.Union[Parenthetical](
			CaptureExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
		),
		participle.Union[Operand](
			CaptureExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
		),
	)
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
	*u = values[0] == "->"
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
	IsUnique Uniqueness `@( "->" | ">>" )`
	Input    Input      `@@`
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

type FunctionModulePart struct {
	Pos lexer.Position

	Functions []Function `( @@ EOL+ )+`
}

func (modPart FunctionModulePart) modulePart() {}
func (modPart FunctionModulePart) pos() lexer.Position {
	return modPart.Pos
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

	Annotations  []string          `( "@" @IDENTIFIER )*`
	Type         Type              `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "->" ) ) )?`
	Name         FunctionName      `( @IDENTIFIER | @OPERATOR )`
	Inputs       []Input           `( "->" @@)*`
	Instructions []RootInstruction `( BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD )`
	Patterns     []Pattern         `| ( EVALS_KEYWORD EOL ( @@ EOL )* END_EVAL )`
}

type Instruction interface {
	instruction()
	pos() lexer.Position
}

type Parenthetical interface {
	parenthesis()
	pos() lexer.Position
}

type Operand interface {
	operand()
	pos() lexer.Position
}

type RootInstruction struct {
	Pos lexer.Position

	Resolution Instruction `@@`
	Ammendment Ammendment  `@@*`
}

type Ammendment struct {
	Pos lexer.Position

	Op    string  `@OPERATOR`
	Value Operand `@@`
}

type BlockConditionalExpression struct {
	Pos lexer.Position

	Condition Parenthetical     `IF_KEYWORD "(" @@ ")" THEN_KEYWORD EOL+`
	Execution []RootInstruction `(@@ EOL+)* END_IF_KEYWORD`
}

func (expression BlockConditionalExpression) instruction() {}
func (expression BlockConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type InlineConditionalExpression struct {
	Pos lexer.Position

	Condition Parenthetical `INLINE_IF_KEYWORD "(" @@ ")"`
	Execution Operand       `@@`
}

func (expression InlineConditionalExpression) instruction() {}
func (expression InlineConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type ParentheticalExpression struct {
	Pos lexer.Position

	Instance Parenthetical `"(" @@ ")"`
}

func (expression ParentheticalExpression) instruction() {}
func (expression ParentheticalExpression) parenthesis() {}
func (expression ParentheticalExpression) operand()     {}
func (expression ParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type CaptureExpression struct {
	Pos lexer.Position

	Instance Parenthetical `BACKTICK @@ BACKTICK`
}

func (expression CaptureExpression) parenthesis() {}
func (expression CaptureExpression) operand()     {}
func (expression CaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	Name string    `@IDENTIFIER`
	Args []Operand `@@*`
}

func (expression ReferenceExpression) instruction() {}
func (expression ReferenceExpression) parenthesis() {}
func (expression ReferenceExpression) operand()     {}
func (expression ReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type Pattern struct {
	Pos lexer.Position

	Name       string          `@IDENTIFIER`
	Params     []string        `@IDENTIFIER* ASSIGNMENT_OPERATOR`
	Definition RootInstruction `@@`
}
