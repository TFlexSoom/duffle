package function

import "github.com/alecthomas/participle/v2/lexer"

type Type struct {
	Name     string `@IDENTIFIER`
	Generics []Type `("[" @@ ("," @@)* "]")?`
}

type Input struct {
	Position lexer.Position

	Type Type   `@@`
	Name string `@IDENTIFIER`
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
