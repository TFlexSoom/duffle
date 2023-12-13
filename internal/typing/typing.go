package typing

import (
	"github.com/tflexsoom/duffle/internal/intermediate"
	"github.com/tflexsoom/duffle/internal/parsing/dflgrammar"
)

const PASS_STRING = "PASS"

func TypeCheck(fileName string, ast dflgrammar.Module) (string, error) {
	symbols := make(map[uint64]intermediate.SymbolPosition)
	_, _, err := intermediate.GetIR(fileName, ast, &symbols)
	if err != nil {
		return "", err
	}

	return PASS_STRING, nil
}
