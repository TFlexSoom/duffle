package intermediate

type InstructionCode uint64

const (
	NOOP InstructionCode = iota
	PARAM
	VALUE
	//...
)

type FunctionId uint64
type ValueId uint64

const NOT_A_VALUE_ID int64 = int64(-1)

type Value struct {
	LiteralValue []byte
	MemoryValue  ValueId
	CacheValue   uint8
}

type Instruction struct {
	Instruction InstructionCode
	Args        []Value
}

type Function struct {
	Name            string
	Definition      []Instruction
	RequiredSymbols string
	RequiredValues  ValueId
}

type Module struct {
	Functions map[FunctionId]Function
	Values    map[ValueId][]byte
}
