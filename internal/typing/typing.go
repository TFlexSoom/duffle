package typing

import (
	"github.com/tflexsoom/duffle/internal/parsing/dflgrammar"
)

const PASS_STRING = "PASS"

func TypeCheck(fileName string, ast dflgrammar.Module) (string, error) {
	return PASS_STRING, nil
}
