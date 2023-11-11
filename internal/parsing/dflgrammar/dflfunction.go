// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflexpression.go
// desc: Parsing Grammar to Build Functions of different types
package dflgrammar

import "github.com/alecthomas/participle/v2/lexer"

type FunctionModulePart struct {
	Pos       lexer.Position
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

type FunctionDefinition interface {
	FunctionDefinition()
	Pos() lexer.Position
}

type Function struct {
	Pos        lexer.Position
	Annotation *string            `"@" @IDENTIFIER? `
	Type       Type               `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "<" ) ) )?`
	Name       FunctionName       `( @IDENTIFIER | @OPERATOR )`
	Inputs     []Input            `( "<" @@ ">")*`
	Definition FunctionDefinition `@@`
}

type ConstexprDefinition struct {
	Constexpr []ConstexprExpression `":=" @@`
}

func (constexprDef ConstexprDefinition) FunctionDefinition() {}
func (constexprDef ConstexprDefinition) Pos() lexer.Position

type BlockDefinition struct {
	Instructions []BlockExpression `BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD`
}

func (constexprDef BlockDefinition) FunctionDefinition() {}
func (constexprDef BlockDefinition) Pos() lexer.Position

type PatternDefinition struct {
	Patterns []Pattern `EVALS_KEYWORD EOL ( @@ EOL )* END_EVAL`
}

func (constexprDef PatternDefinition) FunctionDefinition() {}
func (constexprDef PatternDefinition) Pos() lexer.Position

type Pattern struct {
	Pos lexer.Position

	Name       string           `@IDENTIFIER`
	Params     []string         `@IDENTIFIER* "="`
	Definition InlineExpression `@@`
}
