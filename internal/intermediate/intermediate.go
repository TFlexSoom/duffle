package intermediate

import "github.com/tflexsoom/duffle/internal/container"

type Instruction uint64

const (
	NOOP Instruction = iota
	VALUE
	//...
)

type FunctionId uint64
type ValueId uint64

const NOT_A_VALUE_ID int64 = int64(-1)

type IValue struct {
	IsLiteral    bool
	LiteralValue string
	MemoryValue  ValueId
}

type IExpression struct {
	Instruction Instruction
	Value       IValue
}

type IModule struct {
	Functions       map[FunctionId]string
	Definitions     map[FunctionId]container.GraphTree[IExpression]
	RequiredSymbols map[FunctionId]string
	RequiredValues  map[FunctionId]ValueId
	Values          map[ValueId][]byte
}
