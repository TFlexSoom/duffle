// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflgrammar.go
// desc: Parsing Grammar to Build AST for dfl files
package dflgrammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/tflexsoom/deflemma/internal/types"

	"github.com/tflexsoom/deflemma/internal/parsing/util"
)

func GetDflParser() (types.SourceFileParser, error) {
	var lexer, err = getDflLexer()
	if err != nil {
		return nil, err
	}

	parser, err := participle.Build[Module](
		participle.Lexer(lexer),
		participle.Elide("WHITESPACE"),
		participle.UseLookahead(1),
		participle.Union[ModulePart](
			ImportModulePart{},
			StructModulePart{},
			FunctionModulePart{},
		),
		participle.Union[Import](
			ListImport{},
			SingleImport{},
		),
		participle.Union[FunctionDefinition](
			ConstexprDefinition{},
			BlockDefinition{},
			PatternDefinition{},
		),
		participle.Union[ConstexprExpression](
			ConstexprCaptureExpression{},
			ConstexprParentheticalExpression{},
			ConstexprReferenceExpression{},
			ConstexprOperatorExpression{},
			LiteralExpression{},
		),
		participle.Union[BlockExpression](
			BlockConditionalExpression{},
			BlockCaptureExpression{},
			InlineCaptureExpression{},
			InlineConditionalExpression{},
			ParentheticalExpression{},
			LabelExpression{},
			ReferenceExpression{},
		),
		participle.Union[InlineExpression](
			InlineCaptureExpression{},
			ParentheticalExpression{},
			ReferenceExpression{},
			OperatorExpression{},
		),
		participle.Union[util.Value](
			util.BoolGrammar{},
			util.FloatGrammar{},
			util.IntGrammar{},
			util.StringGrammar{},
			CharGrammar{},
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
