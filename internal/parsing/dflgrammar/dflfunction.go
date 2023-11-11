// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflexpression.go
// desc: Parsing Grammar to Build Functions of different types
package dflgrammar

import "github.com/alecthomas/participle/v2/lexer"

type FunctionModulePart struct {
	Position  lexer.Position
	Functions []Function `( @@ EOL+ )+`
}

func (modPart FunctionModulePart) ModulePart() {}
func (modPart FunctionModulePart) Pos() lexer.Position {
	return modPart.Position
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
	Position   lexer.Position
	Annotation *string            `"@" (@IDENTIFIER | "@") `
	Type       Type               `( @@ (?= IDENTIFIER ( BEGIN_KEYWORD | EVALS_KEYWORD | "<" ) ) )?`
	Name       FunctionName       `( @IDENTIFIER | @OPERATOR )`
	Inputs     []Input            `( "<" @@ ">")*`
	Definition FunctionDefinition `@@`
}

type ConstexprDefinition struct {
	Position  lexer.Position
	Constexpr []ConstexprExpression `":=" @@`
}

func (expr ConstexprDefinition) FunctionDefinition() {}
func (expr ConstexprDefinition) Pos() lexer.Position {
	return expr.Position
}

type BlockDefinition struct {
	Position     lexer.Position
	Instructions []BlockExpression `BEGIN_KEYWORD EOL* ( @@ (";" | EOL) EOL* )* END_KEYWORD`
}

func (expr BlockDefinition) FunctionDefinition() {}
func (expr BlockDefinition) Pos() lexer.Position {
	return expr.Position
}

type PatternDefinition struct {
	Position lexer.Position

	Patterns []Pattern `EVALS_KEYWORD EOL ( @@ EOL )* END_EVAL`
}

func (expr PatternDefinition) FunctionDefinition() {}
func (expr PatternDefinition) Pos() lexer.Position {
	return expr.Position
}

type Pattern struct {
	Position lexer.Position

	Name       string           `@IDENTIFIER`
	Params     []string         `@IDENTIFIER* "="`
	Definition InlineExpression `@@`
}
