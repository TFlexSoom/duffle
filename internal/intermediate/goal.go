package intermediate

import "github.com/tflexsoom/duffle/internal/container"

type OpCode uint32

type SenimentExpression struct {
	TypeId TypeId
	Op     OpCode
	Value  []string
}

type SentimentInput struct {
	Name   string
	TypeId TypeId
}

type Sentiment struct {
	Annotations []string
	Name        string
	Inputs      []SentimentInput
	Definition  container.Tree[SenimentExpression]
}

type Goal struct {
	Sentments map[string]Sentiment
	Types     map[TypeId]string
}

const (
	OPCODE_NOOP OpCode = iota
	OPCODE_CONST
	OPCODE_CALL
)
