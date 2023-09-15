// def lemma compiler entry file for command line processing.
package main

import (
	"log"
	"os"

	"github.com/tflexsoom/deflemma/internal/parsing"

	"github.com/alecthomas/repr"
)

// Entry Function
func main() {
	log.Println("Beginning Parsing!")

	const file = "./example/example3/students.lfun"
	r, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	parser, err := parsing.GetParser()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	ast, err := parser.Parse(file, r)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	repr.Println(ast)

	log.Println("Done!")
}
