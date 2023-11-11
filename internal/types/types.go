package types

import (
	"io"
)

type SourceFileType int

const (
	UNKNOWN_SOURCE_FILE SourceFileType = iota
	FunctionFile
	DataFile
	StructureFile
	SOURCE_FILE_TYPE_LENGTH
)

var SourceFileEnding = map[string]SourceFileType{
	"dfl":  FunctionFile,
	"ddat": DataFile,
}

type SourceFileParser interface {
	ParseSourceFile(string, io.Reader) (interface{}, error)
}
