package puzzles

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"puzzles/solutions/s01"
	"puzzles/solutions/s02"
	"puzzles/solutions/s03"
	"puzzles/solutions/s04"
	"puzzles/solutions/s05"
)

// TestRealInputTiming runs every solver against each of its five committed
// inputs and logs the elapsed time. Reported via t.Log so it shows up only
// under `go test -v`. This is a diagnostic, not a hard threshold — a
// generator or solver change that regresses throughput will show up here.
func TestRealInputTiming(t *testing.T) {
	solvers := []struct {
		name string
		dir  string
		fn   func([]byte) []string
	}{
		{"s01", "01", s01.Solution},
		{"s02", "02", s02.Solution},
		{"s03", "03", s03.Solution},
		{"s04", "04", s04.Solution},
		{"s05", "05", s05.Solution},
	}
	for _, s := range solvers {
		var total time.Duration
		var worst time.Duration
		for j := 1; j <= 5; j++ {
			path := filepath.Join("puzzles", s.dir, "inputs", fmt.Sprintf("%02d.txt", j))
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%s read %s: %v", s.name, path, err)
			}
			start := time.Now()
			out := s.fn(data)
			elapsed := time.Since(start)
			total += elapsed
			if elapsed > worst {
				worst = elapsed
			}
			if len(out) != 2 {
				t.Errorf("%s input %02d: expected 2 answers, got %d", s.name, j, len(out))
			}
			t.Logf("%s input %02d: %s", s.name, j, elapsed)
		}
		t.Logf("%s total: %s (worst single input: %s)", s.name, total, worst)
	}
}
