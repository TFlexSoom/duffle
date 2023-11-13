// author: Tristan Hilbert
// date: 8/29/2023
// filename: dfllexer.go
// desc: Lexer for dflgrammar
package dflgrammar

import (
	"regexp"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/duffle/internal/parsing/util"
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
			{Name: "FUNCTION_SYMBOL", Pattern: `@`, Action: nil},
			{Name: "CONSTEXPR_OPERATOR", Pattern: `:=`, Action: nil},
		},
		"Root": {
			lexer.Include("Spacing"),
			lexer.Include("Expression"),
			{Name: "USE_KEYWORD", Pattern: `use`, Action: nil},
			{Name: "STRUCT_KEYWORD", Pattern: `struct`, Action: nil},
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
