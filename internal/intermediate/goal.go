package intermediate

import "github.com/tflexsoom/duffle/internal/container"

type TypeId uint32
type OpCode uint32

const (
	OPCODE_NOOP OpCode = iota
	OPCODE_CONST
	OPCODE_CALL
)

type SenimentExpression struct {
	TypeId TypeId
	Op     OpCode
	Value  []byte
}

type SentimentInput struct {
	Name   string
	TypeId TypeId
}

type Sentiment struct {
	Annotations []string
	Name        string
	Inputs      SentimentInput
	Definition  container.Tree[SenimentExpression]
}

type Goal struct {
	Sentments map[string]Sentiment
	Types     map[TypeId]string
}
