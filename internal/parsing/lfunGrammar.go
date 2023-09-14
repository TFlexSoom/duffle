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
		"Root": {
			lexer.Include("Spacing"),
			{Name: "USE_KEYWORD", Pattern: `use`, Action: nil},
			{Name: "ARROW", Pattern: `[>-][>]`, Action: nil},
			{Name: "TYPE_PUNCTATION", Pattern: `[@\[\](),]`, Action: nil},
			{Name: "BEGIN_KEYWORD", Pattern: `begin`, Action: lexer.Push("Expression")},
			lexer.Include("Identity"),
		},
		"Expression": {
			lexer.Include("Spacing"),
			{Name: "BACKTICK", Pattern: "`", Action: nil},
			{Name: "EXPR_PUNCTATION", Pattern: `[();]`, Action: nil},
			{Name: "IF_KEYWORD", Pattern: `if`, Action: nil},
			{Name: "THEN_KEYWORD", Pattern: `then`, Action: lexer.Push("Condition")},
			{Name: "EVALS_KEYWORD", Pattern: `evals`, Action: nil},
			{Name: "END_KEYWORD", Pattern: `end`, Action: lexer.Pop()},
			lexer.Include("Identity"),
			{Name: "OPERATOR", Pattern: `[^@\d\w][^\w]*`, Action: nil},
		},
		"Condition": {
			lexer.Include("Spacing"),
			{Name: "END_IF_KEYWORD", Pattern: `endif`, Action: lexer.Pop()},
			lexer.Include("Expression"),
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
		participle.Union[ModulePart](
			ImportModulePart{},
			ConfigModulePart{},
			FunctionModulePart{},
		),
		participle.Union[Import](
			ListImport{},
			SingleImport{},
		),
		participle.Union[RootExpression](
			ConditionalExpression{},
			ParentheticalExpression{},
			OperatorExpression{},
			ReferenceExpression{},
		),
		participle.Union[ContainedExpression](
			ParentheticalExpression{},
			OperatorExpression{},
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

	Annotations []string         `( "@" @IDENTIFIER )*`
	Type        Type             `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "->" ) ) )?`
	Name        FunctionName     `( @IDENTIFIER | @OPERATOR )`
	Inputs      []Input          `( "->" @@)*`
	Expressions []RootExpression `( BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD )`
	Patterns    []Pattern        `| ( EVALS_KEYWORD EOL ( @@ EOL )* EOL )`
}

type RootExpression interface {
	root()
	pos() lexer.Position
}

type ContainedExpression interface {
	contained()
	pos() lexer.Position
}

type ConditionalExpression struct {
	Pos lexer.Position

	Condition ContainedExpression `IF_KEYWORD "(" @@ ")"`
	Execution RootExpression      `( (THEN_KEYWORD EOL @@ EOL+ END_IF_KEYWORD) | @@  )`
}

func (expression ConditionalExpression) root() {}
func (expression ConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type ParentheticalExpression struct {
	Pos lexer.Position

	Instance ContainedExpression `"(" @@ ")"`
}

func (expression ParentheticalExpression) root()      {}
func (expression ParentheticalExpression) contained() {}
func (expression ParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type CallableExpression struct {
	Pos lexer.Position

	Instance ContainedExpression `BACKTICK @@ BACKTICK`
}

func (expression CallableExpression) contained() {}
func (expression CallableExpression) pos() lexer.Position {
	return expression.Pos
}

type OperatorExpression struct {
	Pos lexer.Position

	LInstance ContainedExpression `@@`
	Name      string              `@OPERATOR`
	RInstance ContainedExpression `@@`
}

func (expression OperatorExpression) root()      {}
func (expression OperatorExpression) contained() {}
func (expression OperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	Name string              `@IDENTIFIER`
	Args ContainedExpression `@@*`
}

func (expression ReferenceExpression) root()      {}
func (expression ReferenceExpression) contained() {}
func (expression ReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type Pattern struct {
	Pos lexer.Position

	Name       string         `@IDENTIFIER`
	Params     []string       `@IDENTIFIER* "="`
	Definition RootExpression `@@`
}
