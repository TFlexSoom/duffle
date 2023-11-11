// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflimport.go
// desc: Parsing Grammar for Imports
package dflgrammar

import "github.com/alecthomas/participle/v2/lexer"

type ImportModulePart struct {
	Pos lexer.Position

	Imports []Import `( USE_KEYWORD @@ EOL+ )+`
}

func (modPart ImportModulePart) modulePart() {}
func (modPart ImportModulePart) pos() lexer.Position {
	return modPart.Pos
}

type Import interface {
	importVal()
	pos() lexer.Position
}

type ListImport struct {
	Pos lexer.Position

	Value []string `"(" EOL+ (@IDENTIFIER EOL+)+ ")"`
}

func (listImport ListImport) importVal() {}
func (listImport ListImport) pos() lexer.Position {
	return listImport.Pos
}

type SingleImport struct {
	Pos lexer.Position

	Value []string `@IDENTIFIER`
}

func (singleImport SingleImport) importVal() {}
func (singleImport SingleImport) pos() lexer.Position {
	return singleImport.Pos
}
