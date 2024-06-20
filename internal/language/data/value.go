// author: Tristan Hilbert
// date: 8/29/2023
// filename: grammar.go
// desc: Grammar Utilities for All Parsers in Duffle
// notes: Might trade this out for a TOML parser instead!
package data

import "github.com/alecthomas/participle/v2/lexer"

const (
	VAL_BOOLEAN = SingleGrammar{
		name:    "BOOLEAN",
		regex:   "true|false",
		grammar: `@( "true" | "false")`,
	}

	VAL_FLOAT = SingleGrammar{
		name:    "DECIMAL",
		regex:   `[\d]+\.[\d]*`,
		grammar: `@DECIMAL`,
	}

	VAL_INTEGER = SingleGrammar{
		name:    "INTEGER",
		regex:   `[\d]+`,
		grammar: `@INTEGER`,
	}

	VAL_FUNC_STRING = SingleGrammar{
		name:    "QUOTED_VAL",
		regex:   `"[^"]*"`,
		grammar: `@QUOTED_VAL`,
	}

	VAL_FUNC_CHAR = SingleGrammar{
		name:    "SINGLE_QUOTED_VAL",
		regex:   `'[^']*'`,
		grammar: `@SIGNLE_QUOTED_VAL`,
	}
)

type ConfigValue struct {
	Position lexer.Position
	Val      string `(@IDENTIFIER | @QUOTED_VAL | @TEXT) (@WHITESPACE* (@IDENTIFIER | @TEXT))*`
}

func (f ConfigValue) Pos() lexer.Position {
	return f.Position
}
