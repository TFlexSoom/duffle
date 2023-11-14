// author: Tristan Hilbert
// date: 8/29/2023
// filename: dflimport.go
// desc: Parsing Grammar for Imports
package dflgrammar

import "github.com/alecthomas/participle/v2/lexer"

type ImportModulePart struct {
	Position lexer.Position

	Imports []Import `( USE_KEYWORD @@ EOL+ )+`
}

func (modPart ImportModulePart) ModulePart() {}
func (modPart ImportModulePart) Pos() lexer.Position {
	return modPart.Position
}

type Import interface {
	ImportVal() []string
	Pos() lexer.Position
}

type ListImport struct {
	Position lexer.Position

	Value []string `"(" EOL+ (@IDENTIFIER EOL+)+ ")"`
}

func (listImport ListImport) ImportVal() []string {
	return listImport.Value
}
func (listImport ListImport) Pos() lexer.Position {
	return listImport.Position
}

type SingleImport struct {
	Position lexer.Position

	Value string `@IDENTIFIER`
}

func (singleImport SingleImport) ImportVal() []string {
	return []string{
		singleImport.Value,
	}
}
func (singleImport SingleImport) Pos() lexer.Position {
	return singleImport.Position
}
