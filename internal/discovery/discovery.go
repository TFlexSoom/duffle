package discovery

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tflexsoom/deflemma/internal/types"
)

func DiscoverFiles(path string, verbose bool) (map[types.SourceFileType][]string, error) {
	if verbose {
		log.Default().Printf("READING DIR %v", path)
	}

	fileMap := make(map[types.SourceFileType][]string)

	filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if verbose {
			log.Default().Printf("SCANNING %v WITH ERR %v", path, err)
		}

		if err != nil {
			return err
		}

		sourceFileType := types.SourceFileEnding[GetFileEnding(path)]
		if sourceFileType != types.UNKNOWN_SOURCE_FILE {
			fileMap[sourceFileType] = append(fileMap[sourceFileType], path)
		}

		return nil
	})

	if verbose {
		log.Default().Printf("Result: %v", fileMap)
	}

	return fileMap, nil
}

func GetFileEnding(filename string) string {
	if filename[0:2] == "./" {
		filename = filename[2:]
	}

	fileParts := strings.Split(filename, ".")
	return fileParts[len(fileParts)-1]
}
