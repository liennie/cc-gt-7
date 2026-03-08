package main

import (
	"fmt"
	"os"
	"path/filepath"

	"puzzles/generators/g01"
	"puzzles/generators/g02"
	"puzzles/generators/g03"
	"puzzles/generators/g04"
	"puzzles/generators/g05"
)

const inputCount = 5

var generators = map[string]func() []byte{
	"puzzles/01/inputs": g01.Generate,
	"puzzles/02/inputs": g02.Generate,
	"puzzles/03/inputs": g03.Generate,
	"puzzles/04/inputs": g04.Generate,
	"puzzles/05/inputs": g05.Generate,
}

func main() {
	for outputDir, generator := range generators {
		for idx := range inputCount {
			input := generator()
			err := os.WriteFile(filepath.Join(filepath.FromSlash(outputDir), fmt.Sprintf("%02d.txt", idx+1)), input, 0666)
			if err != nil {
				panic(err)
			}
		}
	}
}
