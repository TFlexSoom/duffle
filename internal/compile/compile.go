package compile

import (
	"log"
	"os"

	"github.com/alecthomas/repr"
	"github.com/tflexsoom/deflemma/internal/discovery"
	"github.com/tflexsoom/deflemma/internal/parsing"
	"github.com/tflexsoom/deflemma/internal/types"
)

type ParserOptions struct {
	ProjectLocations []string
	OutputLocation   string
	FunctionOnly     bool
	DataOnly         bool
	MemoryOnly       bool
	Verbose          bool
}

var parsers = map[types.SourceFileType](func() (types.SourceFileParser, error)){
	types.FunctionFile:  parsing.GetLFunParser,
	types.DataFile:      parsing.GetLDatParser,
	types.StructureFile: parsing.GetLMemParser,
}

func ParseOnly(options ParserOptions) error {
	allFalse := !(options.FunctionOnly || options.DataOnly || options.MemoryOnly)

	fileFilter := map[types.SourceFileType]bool{
		types.FunctionFile:  allFalse || options.FunctionOnly,
		types.DataFile:      allFalse || options.DataOnly,
		types.StructureFile: allFalse || options.MemoryOnly,
	}

	fileMap := make(map[types.SourceFileType][]string)
	for _, location := range options.ProjectLocations {
		discoveredFileMap, err := discovery.DiscoverFiles(location, options.Verbose)
		if err != nil {
			return err
		}

		for i := types.FunctionFile; i < types.SOURCE_FILE_TYPE_LENGTH; i++ {
			fileMap[i] = append(fileMap[i], discoveredFileMap[i][:]...)
		}
	}

	outFile, err := os.OpenFile(
		options.OutputLocation,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}

	for sourceFileType, files := range fileMap {
		if fileFilter[sourceFileType] != true {
			continue
		}

		for _, file := range files {
			reader, err := os.Open(file)
			if err != nil {
				return err
			}

			parser, err := parsers[sourceFileType]()
			if err != nil {
				return err
			}

			ast, err := parser.ParseSourceFile(file, reader)
			if err != nil {
				return err
			}

			num, err := outFile.WriteString(repr.String(ast))
			if err != nil {
				return err
			}

			if options.Verbose {
				log.Printf("%d bytes written!", num)
			}
		}
	}

	return nil
}

type CompilerOptions struct {
	ProjectLocations []string
	OutputLocation   string
	Backend          string
	Verbose          bool
}

func Compile(options CompilerOptions) error {
	return nil
}
