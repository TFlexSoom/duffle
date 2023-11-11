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
	block()
	pos() lexer.Position
}

type InlineExpression interface {
	inline()
	pos() lexer.Position
}

type ConstexprExpression interface {
	constexpr()
	pos() lexer.Position
}

type BlockConditionalExpression struct {
	Pos lexer.Position

	Condition      BlockExpression       `IF_KEYWORD "(" @@ ")"`
	Execution      []InlineExpression    `THEN_KEYWORD EOL+ (@@ EOL*)*`
	SubConditional []SubBlockConditional `@@*`
	Alternative    []InlineExpression    `(ELSE_KEYWORD EOL+ (@@ EOL+)*)? END_IF_KEYWORD`
}

func (expression BlockConditionalExpression) block() {}
func (expression BlockConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type SubBlockConditional struct {
	Pos lexer.Position

	Condition InlineExpression   `ELSEIF_KEYWORD "(" @@ ")"`
	Execution []InlineExpression `THEN_KEYWORD EOL+ (@@ EOL+)*`
}

type InlineConditionalExpression struct {
	Pos lexer.Position

	Condition     InlineExpression `INLINE_IF_KEYWORD "(" @@ ")"`
	NextExecution InlineExpression `@@`
}

func (expression InlineConditionalExpression) block() {}
func (expression InlineConditionalExpression) pos() lexer.Position {
	return expression.Pos
}

type ParentheticalExpression struct {
	Pos lexer.Position

	Execution     InlineExpression `"(" @@ ")"`
	NextExecution InlineExpression `@@?`
}

func (expression ParentheticalExpression) block()  {}
func (expression ParentheticalExpression) inline() {}
func (expression ParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type ConstexprParentheticalExpression struct {
	Pos lexer.Position

	Execution     ConstexprExpression `"(" @@ ")"`
	NextExecution ConstexprExpression `@@?`
}

func (expression ConstexprParentheticalExpression) constexpr()
func (expression ConstexprParentheticalExpression) pos() lexer.Position {
	return expression.Pos
}

type CaptureExpression struct {
	Pos lexer.Position

	Execution     InlineExpression `BACKTICK @@ BACKTICK`
	NextExecution InlineExpression `@@?`
}

func (expression CaptureExpression) block()  {}
func (expression CaptureExpression) inline() {}
func (expression CaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ConstexprCaptureExpression struct {
	Pos lexer.Position

	Execution     ConstexprExpression `BACKTICK @@ BACKTICK`
	NextExecution ConstexprExpression `@@?`
}

func (expression ConstexprCaptureExpression) constexpr() {}
func (expression ConstexprCaptureExpression) pos() lexer.Position {
	return expression.Pos
}

type ReferenceExpression struct {
	Pos lexer.Position

	ReferenceGroup []string         `@IDENTIFIER+`
	NextExecution  InlineExpression `@@?`
}

func (expression ReferenceExpression) block()  {}
func (expression ReferenceExpression) inline() {}
func (expression ReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type ConstexprReferenceExpression struct {
	Pos lexer.Position

	ReferenceGroup []string            `@IDENTIFIER+`
	NextExecution  ConstexprExpression `@@?`
}

func (expression ConstexprReferenceExpression) constexpr() {}
func (expression ConstexprReferenceExpression) pos() lexer.Position {
	return expression.Pos
}

type OperatorExpression struct {
	Pos lexer.Position

	ReferenceGroup []string         `@OPERATOR`
	NextExecution  InlineExpression `@@?`
}

func (expression OperatorExpression) inline() {}
func (expression OperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type ConstexprOperatorExpression struct {
	Pos lexer.Position

	ReferenceGroup []string            `@OPERATOR`
	NextExecution  ConstexprExpression `@@?`
}

func (expression ConstexprOperatorExpression) constexpr() {}
func (expression ConstexprOperatorExpression) pos() lexer.Position {
	return expression.Pos
}

type LiteralExpression struct {
	Pos lexer.Position

	Value util.Value `@@`
}

func (expression LiteralExpression) constexpr() {}
func (expression LiteralExpression) pos() lexer.Position {
	return expression.Pos
}
