// Program computeanswers reads generated puzzle inputs and prints the answers.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"puzzles/solutions/s01"
	"puzzles/solutions/s02"
	"puzzles/solutions/s03"
	"puzzles/solutions/s04"
	"puzzles/solutions/s05"
)

func main() {
	solvers := []struct {
		puzzle string
		fn     func([]byte) []string
	}{
		{"01", s01.Solution},
		{"02", s02.Solution},
		{"03", s03.Solution},
		{"04", s04.Solution},
		{"05", s05.Solution},
	}
	for _, s := range solvers {
		fmt.Printf("=== puzzle %s ===\n", s.puzzle)
		for i := 1; i <= 5; i++ {
			p := filepath.Join("puzzles", s.puzzle, "inputs", fmt.Sprintf("%02d.txt", i))
			b, err := os.ReadFile(p)
			if err != nil {
				fmt.Printf("  %02d: ERR %v\n", i, err)
				continue
			}
			ans := s.fn(b)
			fmt.Printf("  %02d: part1=%q part2=%q\n", i, ans[0], ans[1])
		}
	}
}
