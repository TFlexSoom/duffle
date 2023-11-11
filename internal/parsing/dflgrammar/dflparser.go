// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflparser.go
// desc: Parsing Interface Implementation for easy programming usage
package dflgrammar

import (
	"io"

	"github.com/alecthomas/participle/v2"
)

type ModuleParser struct {
	Parser *participle.Parser[Module]
}

func (modParser *ModuleParser) ParseSourceFile(
	fileName string,
	reader io.Reader,
) (interface{}, error) {
	return modParser.Parser.Parse(fileName, reader)
}
