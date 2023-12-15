package command

import (
	"errors"
	"log"
	"os"

	"github.com/alecthomas/repr"
	"github.com/tflexsoom/duffle/internal/discovery"
	"github.com/tflexsoom/duffle/internal/files"
	"github.com/tflexsoom/duffle/internal/parsing/ddatgrammar"
	"github.com/tflexsoom/duffle/internal/parsing/dflgrammar"
	"github.com/tflexsoom/duffle/internal/typing"
)

func getFileMap(projectLocations []string, isVerbose bool) (map[files.SourceFileType][]string, error) {
	fileMap := make(map[files.SourceFileType][]string)
	for _, location := range projectLocations {
		discoveredFileMap, err := discovery.DiscoverFiles(location, isVerbose)
		if err != nil {
			return nil, err
		}

		for i := files.FunctionFile; i < files.SOURCE_FILE_TYPE_LENGTH; i++ {
			fileMap[i] = append(fileMap[i], discoveredFileMap[i][:]...)
		}
	}

	return fileMap, nil
}

func writeOutput(fileName string, data string, isVerbose bool) error {
	outFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}
	defer outFile.Close()

	num, err := outFile.WriteString(data)
	if err != nil {
		return err
	}

	if isVerbose {
		log.Printf("%d bytes written!", num)
	}

	return nil
}

type FileLogicOptions interface {
	GetProjectLocations() []string
	GetFunctionFilesOnly() bool
	GetDataFilesOnly() bool
	GetOutputLocation() string
	IsVerbose() bool
}

func withFileLogicAndOutput(
	fileLogicOptions FileLogicOptions,
	processor func(files.SourceFileType, string, *os.File) (string, error),
) error {
	fileMap, err := getFileMap(fileLogicOptions.GetProjectLocations(), fileLogicOptions.IsVerbose())
	if err != nil {
		return err
	}

	defaultTrue := !(fileLogicOptions.GetFunctionFilesOnly() || fileLogicOptions.GetDataFilesOnly())
	fileFilter := map[files.SourceFileType]bool{
		files.FunctionFile: defaultTrue || fileLogicOptions.GetFunctionFilesOnly(),
		files.DataFile:     defaultTrue || fileLogicOptions.GetDataFilesOnly(),
	}

	tempFileName := fileLogicOptions.GetOutputLocation() + "_temp"
	os.Remove(tempFileName)

	for sourceFileType, files := range fileMap {
		if !fileFilter[sourceFileType] {
			continue
		}

		for _, file := range files {
			reader, err := os.Open(file)
			if err != nil {
				return err
			}

			data, err := processor(sourceFileType, file, reader)
			if err != nil {
				return err
			}

			err = writeOutput(tempFileName, data, fileLogicOptions.IsVerbose())
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(fileLogicOptions.GetOutputLocation())
	if err != nil {
		return err
	}

	err = os.Rename(tempFileName, fileLogicOptions.GetOutputLocation())
	if err != nil {
		return err
	}

	return nil
}

var parsers = map[files.SourceFileType](func() (files.SourceFileParser, error)){
	files.FunctionFile: dflgrammar.GetDflParser,
	files.DataFile:     ddatgrammar.GetDdatParser,
}

func parseProcessor(sourceFileType files.SourceFileType, file string, reader *os.File) (interface{}, error) {
	parser, err := parsers[sourceFileType]()
	if err != nil {
		return nil, err
	}

	ast, err := parser.ParseSourceFile(file, reader)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func parseStringProcessor(sourceFileType files.SourceFileType, file string, reader *os.File) (string, error) {
	ast, err := parseProcessor(sourceFileType, file, reader)
	if err != nil {
		return "", err
	}

	return repr.String(ast), nil
}

type ParserOptions struct {
	ProjectLocations []string
	OutputLocation   string
	FunctionOnly     bool
	DataOnly         bool
	Verbose          bool
}

func (options ParserOptions) GetProjectLocations() []string {
	return options.ProjectLocations
}

func (options ParserOptions) GetFunctionFilesOnly() bool {
	return options.FunctionOnly
}

func (options ParserOptions) GetDataFilesOnly() bool {
	return options.DataOnly
}

func (options ParserOptions) GetOutputLocation() string {
	return options.OutputLocation
}

func (options ParserOptions) IsVerbose() bool {
	return options.Verbose
}

func ParseOnly(options ParserOptions) error {
	return withFileLogicAndOutput(options, parseStringProcessor)
}

func typeCheckProcessor(sourceFileType files.SourceFileType, file string, reader *os.File) (string, error) {
	ast, err := parseProcessor(sourceFileType, file, reader)
	if err != nil {
		return "", err
	}

	casted, isOk := ast.(*dflgrammar.Module)
	if !isOk {
		return "", errors.New("casting module did not work for parsing to ir")
	}

	data, err := typing.TypeCheck(file, *casted)
	if err != nil {
		return "", err
	}

	return data, nil
}

type TypeCheckOptions struct {
	ProjectLocations []string
	OutputLocation   string
	Verbose          bool
}

func (options TypeCheckOptions) GetProjectLocations() []string {
	return options.ProjectLocations
}

func (options TypeCheckOptions) GetFunctionFilesOnly() bool {
	return true
}

func (options TypeCheckOptions) GetDataFilesOnly() bool {
	return false
}

func (options TypeCheckOptions) GetOutputLocation() string {
	return options.OutputLocation
}

func (options TypeCheckOptions) IsVerbose() bool {
	return options.Verbose
}

func TypeCheckOnly(options TypeCheckOptions) error {
	return withFileLogicAndOutput(options, typeCheckProcessor)
}

func compileProcessor(sourceFileType files.SourceFileType, file string, reader *os.File) (string, error) {
	return "", nil
}

type CompilerOptions struct {
	ProjectLocations []string
	OutputLocation   string
	FunctionOnly     bool
	DataOnly         bool
	Backend          string
	Verbose          bool
}

func (options CompilerOptions) GetProjectLocations() []string {
	return options.ProjectLocations
}

func (options CompilerOptions) GetFunctionFilesOnly() bool {
	return options.FunctionOnly
}

func (options CompilerOptions) GetDataFilesOnly() bool {
	return options.DataOnly
}

func (options CompilerOptions) GetOutputLocation() string {
	return options.OutputLocation
}

func (options CompilerOptions) IsVerbose() bool {
	return options.Verbose
}

func Compile(options CompilerOptions) error {
	return withFileLogicAndOutput(options, compileProcessor)
}
