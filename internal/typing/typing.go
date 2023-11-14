package typing

import (
	"errors"

	"github.com/tflexsoom/duffle/internal/intermediate"
	"github.com/tflexsoom/duffle/internal/parsing/dflgrammar"
)

const PASS_STRING = "PASS"

func TypeCheck(fileName string, ast interface{}) (string, error) {
	vfunc, isOk := ast.(dflgrammar.Module)
	if !isOk {
		return "", errors.New("unknown ast type")
	}

	symbols := make(map[uint64]intermediate.SymbolPosition)
	_, _, err := intermediate.GetIR(fileName, vfunc, &symbols)
	if err != nil {
		return "", err
	}

	return PASS_STRING, nil
}
