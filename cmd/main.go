// def lemma compiler entry file for command line processing.
package main

import (
	"log"
	"os"

	"github.com/tflexsoom/deflemma/internal/parsing"
)


// Entry Function
func main() {
	log.Println("Beginning Parsing!")

	const file = "./example/example0/helloWorld.lfun"
	r, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	_, err = parsing.GetParser().Parse(file, r)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	log.Println("Done!")
}