// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflexpression.go
// desc: Parsing Grammar to Build Expressions of Different Types
package dflgrammar

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/tflexsoom/deflemma/internal/parsing/util"
)

type BlockExpression interface {
	Block()
	Pos() lexer.Position
}

type InlineExpression interface {
	Inline()
	Pos() lexer.Position
}

type ConstexprExpression interface {
	Constexpr()
	Pos() lexer.Position
}

type BlockConditionalExpression struct {
	Position lexer.Position

	Condition      BlockExpression       `IF_KEYWORD "(" @@ ")"`
	Execution      []InlineExpression    `THEN_KEYWORD EOL+ (@@ EOL+)*`
	SubConditional []SubBlockConditional `@@*`
	Alternative    []InlineExpression    `(ELSE_KEYWORD EOL+ (@@ EOL+)*)? END_IF_KEYWORD`
}

func (expression BlockConditionalExpression) Block() {}
func (expression BlockConditionalExpression) Pos() lexer.Position {
	return expression.Position
}

type SubBlockConditional struct {
	Position lexer.Position

	Condition InlineExpression   `ELSEIF_KEYWORD "(" @@ ")"`
	Execution []InlineExpression `THEN_KEYWORD EOL+ (@@ EOL+)*`
}

type LabelExpression struct {
	Position lexer.Position

	Label      string           `@IDENTIFIER`
	Resolution InlineExpression `":=" @@`
}

func (expression LabelExpression) Block() {}
func (expression LabelExpression) Pos() lexer.Position {
	return expression.Position
}

type InlineConditionalExpression struct {
	Position lexer.Position

	Condition          InlineExpression `INLINE_IF_KEYWORD "(" @@ ")"`
	ConditionExecution InlineExpression `@@`
}

func (expression InlineConditionalExpression) Block() {}
func (expression InlineConditionalExpression) Pos() lexer.Position {
	return expression.Position
}

type ParentheticalExpression struct {
	Position lexer.Position

	Execution     InlineExpression `"(" EOL* @@ EOL* ")"`
	NextExecution InlineExpression `@@?`
}

func (expression ParentheticalExpression) Block()  {}
func (expression ParentheticalExpression) Inline() {}
func (expression ParentheticalExpression) Pos() lexer.Position {
	return expression.Position
}

type ConstexprParentheticalExpression struct {
	Position lexer.Position

	Execution     ConstexprExpression `"(" @@ ")"`
	NextExecution ConstexprExpression `@@?`
}

func (expression ConstexprParentheticalExpression) Constexpr() {}
func (expression ConstexprParentheticalExpression) Pos() lexer.Position {
	return expression.Position
}

type CaptureExpression struct {
	Position lexer.Position

	Execution     InlineExpression `BACKTICK @@ BACKTICK`
	NextExecution InlineExpression `@@?`
}

func (expression CaptureExpression) Block()  {}
func (expression CaptureExpression) Inline() {}
func (expression CaptureExpression) Pos() lexer.Position {
	return expression.Position
}

type ConstexprCaptureExpression struct {
	Position lexer.Position

	Execution     ConstexprExpression `BACKTICK @@ BACKTICK`
	NextExecution ConstexprExpression `@@?`
}

func (expression ConstexprCaptureExpression) Constexpr() {}
func (expression ConstexprCaptureExpression) Pos() lexer.Position {
	return expression.Position
}

type ReferenceExpression struct {
	Position lexer.Position

	ReferenceGroup []string         `@IDENTIFIER+`
	NextExecution  InlineExpression `@@?`
}

func (expression ReferenceExpression) Block()  {}
func (expression ReferenceExpression) Inline() {}
func (expression ReferenceExpression) Pos() lexer.Position {
	return expression.Position
}

type ConstexprReferenceExpression struct {
	Position lexer.Position

	ReferenceGroup []string            `@IDENTIFIER+`
	NextExecution  ConstexprExpression `@@?`
}

func (expression ConstexprReferenceExpression) Constexpr() {}
func (expression ConstexprReferenceExpression) Pos() lexer.Position {
	return expression.Position
}

type OperatorExpression struct {
	Position lexer.Position

	ReferenceGroup []string         `@OPERATOR`
	NextExecution  InlineExpression `@@?`
}

func (expression OperatorExpression) Inline() {}
func (expression OperatorExpression) Pos() lexer.Position {
	return expression.Position
}

type ConstexprOperatorExpression struct {
	Position lexer.Position

	ReferenceGroup []string            `@OPERATOR`
	NextExecution  ConstexprExpression `@@?`
}

func (expression ConstexprOperatorExpression) Constexpr() {}
func (expression ConstexprOperatorExpression) Pos() lexer.Position {
	return expression.Position
}

type LiteralExpression struct {
	Position lexer.Position

	Value util.Value `@@`
}

func (expression LiteralExpression) Constexpr() {}
func (expression LiteralExpression) Pos() lexer.Position {
	return expression.Position
}
