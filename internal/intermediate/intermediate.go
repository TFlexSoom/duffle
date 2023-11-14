package intermediate

import (
	"github.com/tflexsoom/duffle/internal/parsing/dflgrammar"
)

/***
  TODO
  With the current implementation we are stripping all of the lexer
  Positions. Here are a few options

  - Add it to every structure below (doesn't seem necessary since most things won't error)
  - Just use functions for line/col references (easy but bad info)
  - Add it to some structures below
  - Create an auxiliary Data Structure

  I think the main issue/unknown is how will we recover from an error?
  - Will we flag the bad variable?
  - Will we flag the bad call?
  - Will we flag everything for the user?

  This is why this is being kept as a todo for now.
***/

type SymbolPosition struct {
	FileName   string
	LineNumber int
	ColNumber  int
}

type Type struct {
	Name string

	// TODO Make a 1-D Structure Tree
	GenericParamsTree []Type
}

type ExpressionTreeNode struct {
	Type string
	Name string

	// TODO Make a 1-D Structure Tree
	Children []ExpressionTreeNode
}

type Input struct {
	Name string
	Type Type
}

type Reference struct {
	UniqueId   uint64
	Name       string
	Inputs     []Input
	ReturnType Type
	Definition DefinitionMonad
}

type DefinitionMonad struct {
	IsStruct    bool
	IsFunction  bool
	IsOperator  bool
	Annotation  string
	Expressions []ExpressionTreeNode
}

type ImportName string

var atomicCounter uint64 = 0

func GetUniqueId() uint64 {
	atomicCounter += 1
	return atomicCounter
}

func GetIR(fileName string, ast dflgrammar.Module, symbolPositions *map[uint64]SymbolPosition) ([]ImportName, []Reference, error) {
	imports := make([]string, 0, max(1024, len(ast.ModuleParts)))
	references := make([]Reference, 0, max(1024, len(ast.ModuleParts)))

	addSymbol := func(uid uint64, lineNumber int, colNumber int) {
		(*symbolPositions)[uid] = SymbolPosition{
			FileName:   fileName,
			LineNumber: lineNumber,
			ColNumber:  colNumber,
		}
	}

	for _, part := range ast.ModuleParts {
		importable, isOk := part.(dflgrammar.ImportModulePart)
		if isOk {
			appendImports(&imports, importable)
		}

		structable, isOk := part.(dflgrammar.StructModulePart)
		if isOk {
			appendStructRefs(&references, structable, addSymbol)
		}

		functionable, isOk := part.(dflgrammar.FunctionModulePart)
		if isOk {
			appendFunctionRefs(&references, functionable, addSymbol)
		}
	}

	importsRetyped := make([]ImportName, 0, len(imports))
	for _, importName := range imports {
		importsRetyped = append(importsRetyped, ImportName(importName))
	}

	return importsRetyped, references, nil
}

func appendImports(imports *[]string, node dflgrammar.ImportModulePart) {
	for _, useCall := range node.Imports {
		*imports = append(*imports, useCall.ImportVal()...)
	}
}

func grammarToIRTypes(grammarTypes []dflgrammar.Type) []Type {
	if len(grammarTypes) == 0 {
		return nil
	}

	result := make([]Type, 0, len(grammarTypes))
	for _, typeRef := range grammarTypes {
		result = append(result, Type{
			Name:              typeRef.Name,
			GenericParamsTree: grammarToIRTypes(typeRef.Generics),
		})
	}

	return result
}

func grammarToIRInputs(grammarInputs []dflgrammar.Input) []Input {
	result := make([]Input, 0, len(grammarInputs))
	for _, input := range grammarInputs {
		result = append(result, Input{
			Name: input.Name,
			Type: Type{
				Name:              input.Type.Name,
				GenericParamsTree: grammarToIRTypes(input.Type.Generics),
			},
		})
	}

	return result
}

func appendStructRefs(references *[]Reference, node dflgrammar.StructModulePart, addSymbol func(uint64, int, int)) {
	for _, structRef := range node.Structs {
		uid := GetUniqueId()
		*references = append(*references, Reference{
			UniqueId: uid,
			Name:     structRef.Name,
			Inputs:   grammarToIRInputs(structRef.Fields),
			ReturnType: Type{
				Name:              structRef.Name,
				GenericParamsTree: nil,
			},
			Definition: DefinitionMonad{
				IsStruct: true,
			},
		})

		addSymbol(uid, structRef.Position.Line, structRef.Position.Column)
	}
}

func generateExpressionTree(definition dflgrammar.FunctionDefinition) []ExpressionTreeNode {
	// TODO
	return []ExpressionTreeNode{}
}

func appendFunctionRefs(references *[]Reference, node dflgrammar.FunctionModulePart, addSymbol func(uint64, int, int)) {
	for _, funcRef := range node.Functions {
		annotation := ""
		if funcRef.Annotation != nil {
			annotation = *funcRef.Annotation
		}

		uid := GetUniqueId()

		*references = append(*references, Reference{
			UniqueId: uid,
			Name:     funcRef.Name.Name,
			Inputs:   grammarToIRInputs(funcRef.Inputs),
			ReturnType: Type{
				Name:              funcRef.Type.Name,
				GenericParamsTree: grammarToIRTypes(funcRef.Type.Generics),
			},
			Definition: DefinitionMonad{
				IsFunction:  true,
				IsOperator:  funcRef.Name.IsOperator,
				Annotation:  annotation,
				Expressions: generateExpressionTree(funcRef.Definition),
			},
		})

		addSymbol(uid, funcRef.Position.Line, funcRef.Position.Column)
	}
}
