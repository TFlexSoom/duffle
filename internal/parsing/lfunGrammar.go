// author: Tristan Hilbert
// date: 8/29/2023
// filename: lfunGrammar.go
// desc: Parsing Grammar to Build AST for lfun files
package parsing

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// // Lexer
func getLexer() (*lexer.StatefulDefinition, error) {
	return lexer.NewSimple([]lexer.SimpleRule{
		{"IDENTIFIER", `[a-zA-Z]+[a-zA-Z\d]+\w*`},
		{"WHITESPACE", `[ \t]+`},
		{"EOL", `(\r)?\n`},
		// {"OPERATOR", `[^\d\w]+[^\w]*`},
		{"BACKTICK", "`"},
		{"PUNCTUATION", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
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
		// participle.UseLookahead(5),
		participle.Union[Expression](
			// DelayedExpression{},
			ResolvingExpression{},
			ReferenceExpression{},
			// AssignmentExpression{},
			// OperatorExpression{},
		))
}

//// Grammar

type Module struct {
	Pos lexer.Position

	Imports   []*Import   `("use" @@ EOL+)*`
	Configs   []*Config   `(@@ EOL+)*`
	Functions []*Function `(@@ EOL+)*`
}

type Import struct {
	Pos lexer.Position

	ListImport   []string `"(" ( EOL+ @IDENTIFIER EOL+ ( "," EOL+  @IDENTIFIER )* )? ")"`
	SingleImport string   `| @IDENTIFIER`
}

type Uniqueness bool

func (u *Uniqueness) Capture(values []string) error {
	*u = values[0] == "-"
	return nil
}

type Config struct {
	Pos      lexer.Position
	IsUnique *Uniqueness `@("-" | ">" ) ">"`
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

type Function struct {
	Pos lexer.Position

	Annotations []string      `("@" @IDENTIFIER)*`
	Type        *Type         `(@@ (?= IDENTIFIER "{" | ":" | "~" ))?`
	Name        string        `@IDENTIFIER`
	Inputs      []*Input      `( "~" @@)*`
	Expressions []*Expression `( "{" EOL* ( @@ (";" | EOL) EOL* )* "}" )?`
	Patterns    []*Pattern    `( ":" EOL ( @@ EOL )* EOL )?`
}

type Expression interface {
	expression()
	lexPosition() lexer.Position
}

func (config Config) expression() {}
func (config Config) lexPosition() lexer.Position {
	return config.Pos
}

// type DelayedExpression struct {
// 	Pos lexer.Position

// 	Expression *Expression `BACKTICK @@ BACKTICK`
// }

// func (delayed DelayedExpression) expression() {}
// func (delayed DelayedExpression) lexPosition() lexer.Position {
// 	return delayed.Pos
// }

type ResolvingExpression struct {
	Pos lexer.Position

	Name string        `@IDENTIFIER`
	Args []*Expression `@@+`
}

func (resolving ResolvingExpression) expression() {}
func (resolving ResolvingExpression) lexPosition() lexer.Position {
	return resolving.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	Name string `@IDENTIFIER`
}

func (reference ReferenceExpression) expression() {}
func (reference ReferenceExpression) lexPosition() lexer.Position {
	return reference.Pos
}

// type Shallowness bool

// func (isShallow *Shallowness) Capture(values []string) error {
// 	*isShallow = values[0] == "<-"
// 	return nil
// }

// type AssignmentExpression struct {
// 	Pos lexer.Position

// 	Name      string       `@IDENTIFIER`
// 	IsShallow *Shallowness `@("<-" | "<<")`
// 	Value     *Expression  `@@`
// }

// func (assigning AssignmentExpression) expression() {}
// func (assigning AssignmentExpression) lexPosition() lexer.Position {
// 	return assigning.Pos
// }

// type OperatorExpression struct {
// 	Pos lexer.Position

// 	Left  *Expression `@@`
// 	Op    string      `@IDENTIFIER`
// 	Right *Expression `@@`
// }

// func (operator OperatorExpression) expression() {}
// func (operator OperatorExpression) lexPosition() lexer.Position {
// 	return operator.Pos
// }

type Pattern struct {
	Pos lexer.Position

	Name       string      `@IDENTIFIER`
	Inputs     []*Input    `@@* "="`
	Definition *Expression `@@`
}
