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
const identifierRegexPattern = `[a-zA-Z]+[a-zA-Z\d]+\w*`

var identifierRegex = regexp.MustCompile(identifierRegexPattern)

func getLexer() (*lexer.StatefulDefinition, error) {
	return lexer.NewSimple([]lexer.SimpleRule{
		{"EOL", `(\r)?\n`},
		{"WHITESPACE", `[ \t]+`},
		{"BACKTICK", "`"},
		{"ARROW", `[>-][>]`},
		{"PUNCTUATION", `[@()\[\];,]`},
		{"BEGIN_KEYWORD", `begin`},
		{"END_KEYWORD", `end`},
		{"EVALS_KEYWORD", `evals`},
		{"IDENTIFIER", identifierRegexPattern},
		{"OPERATOR", `[^@\d\w]+[^\w]*`},
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

	Imports []*Import `( "use" @@ EOL+ )+`
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

	Value []string `"(" ( EOL+ @IDENTIFIER EOL+ ( "," EOL+  @IDENTIFIER )* )? ")"`
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
	IsUnique *Uniqueness `@( "->" | ">>" )`
	Input    *Input      `@@`
}

type Input struct {
	Pos lexer.Position

	Type *Type  `@@`
	Name string `@IDENTIFIER`
}

type Type struct {
	Name string `@IDENTIFIER`
}

type FunctionModulePart struct {
	Pos lexer.Position

	Functions []*Function `( @@ EOL+ )+`
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

	Annotations []string      `( "@" @IDENTIFIER )*`
	Type        *Type         `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "->" ) ) )?`
	Name        FunctionName  `( @IDENTIFIER | @OPERATOR )`
	Inputs      []*Input      `( "->" @@)*`
	Expressions []*Expression `( BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD )?`
	Patterns    []*Pattern    `( EVALS_KEYWORD EOL ( @@ EOL )* EOL )?`
}

type Expression struct {
	Pos lexer.Position

	Name string      `( @IDENTIFIER | @OPERATOR )`
	Arg  *Expression `(( "(" @@ ")" ) | @@ )*`
}

type Pattern struct {
	Pos lexer.Position

	Name       string      `@IDENTIFIER`
	Params     []string    `@IDENTIFIER* "="`
	Definition *Expression `@@`
}
