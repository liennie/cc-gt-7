package s05

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

func TestSolution(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want1, want2 string
	}{
		{
			name: "path of 3, one query",
			// weights: a=5, b=100, c=5. Edges: a-b, b-c (path a-b-c).
			// Max IS: {a, c} = 10 (skip b) OR just {b} = 100. Best = 100.
			// Query "a c": same thread -> contributes 0. Part 1 = 0.
			input: `a 5
b 100
c 5

a b
b c

a c
`,
			want1: "0",
			want2: "100",
		},
		{
			name: "two trees, mixed queries",
			// Tree 1: a-b, a-c. Weights a=10, b=8, c=7. Take a (10) skips b,c.
			//   take_a = 10 + skip_b + skip_c = 10 + 0 + 0 = 10
			//   skip_a = max(take_b, skip_b) + max(take_c, skip_c) = 8 + 7 = 15
			// Best tree 1 = 15.
			// Tree 2: d-e. Weights d=3, e=20. Best = 20.
			// Part 2 = 15 + 20 = 35.
			// Queries: "a b" same thread -> 0. "a d" cross-thread -> 10+3 = 13.
			// "b c" same thread -> 0. Part 1 = 13.
			input: `a 10
b 8
c 7
d 3
e 20

a b
a c
d e

a b
a d
b c
`,
			want1: "13",
			want2: "35",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Solution([]byte(tc.input))
			if len(got) != 2 {
				t.Fatalf("expected 2 answers, got %d", len(got))
			}
			if got[0] != tc.want1 {
				t.Errorf("part1: got %q, want %q", got[0], tc.want1)
			}
			if got[1] != tc.want2 {
				t.Errorf("part2: got %q, want %q", got[1], tc.want2)
			}
		})
	}
}

// randomForestInputS05 builds a random forest of up to `n` nodes as
// several trees, each grown by attaching new nodes to a random existing
// node (random recursive tree). Returns the puzzle-format input string
// plus the ground-truth Part 1 answer (sum over cross-tree queries) and
// the exact Part 2 answer (2^n brute-force max independent set).
func randomForestInputS05(r *rand.Rand, n int) (string, int64, int64) {
	names := make([]string, n)
	weights := make([]int64, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("t%02d", i)
		weights[i] = int64(1 + r.IntN(50))
	}
	adj := make([][]int, n)
	// Split into trees of size 1..min(remaining, 5). Track tree id per node.
	treeID := make([]int, n)
	pos := 0
	nt := 0
	for pos < n {
		size := 1 + r.IntN(5)
		if pos+size > n {
			size = n - pos
		}
		for k := 0; k < size; k++ {
			treeID[pos+k] = nt
			if k > 0 {
				parent := pos + r.IntN(k)
				adj[parent] = append(adj[parent], pos+k)
				adj[pos+k] = append(adj[pos+k], parent)
			}
		}
		pos += size
		nt++
	}

	// Queries: mix of cross-tree and same-tree pairs.
	nQueries := 4 + r.IntN(6)
	var qb strings.Builder
	var wantP1 int64
	for q := 0; q < nQueries; q++ {
		a := r.IntN(n)
		b := r.IntN(n)
		if a == b {
			continue
		}
		fmt.Fprintf(&qb, "%s %s\n", names[a], names[b])
		if treeID[a] != treeID[b] {
			wantP1 += weights[a] + weights[b]
		}
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%s %d\n", names[i], weights[i])
	}
	sb.WriteString("\n")
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			if u < v {
				fmt.Fprintf(&sb, "%s %s\n", names[u], names[v])
			}
		}
	}
	sb.WriteString("\n")
	sb.WriteString(qb.String())

	wantP2 := bruteMISS05(adj, weights)
	return sb.String(), wantP1, wantP2
}

// bruteMISS05 enumerates all 2^n vertex subsets and returns the max total
// weight of an independent set (no two adjacent vertices selected).
// Exponential; use only for small n (<= 20).
func bruteMISS05(adj [][]int, weights []int64) int64 {
	n := len(adj)
	if n == 0 {
		return 0
	}
	if n > 22 {
		panic("bruteMISS05: n too large")
	}
	var best int64
	for mask := 0; mask < (1 << n); mask++ {
		ok := true
		for u := 0; u < n && ok; u++ {
			if mask&(1<<u) == 0 {
				continue
			}
			for _, v := range adj[u] {
				if v > u && mask&(1<<v) != 0 {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		var sum int64
		for u := 0; u < n; u++ {
			if mask&(1<<u) != 0 {
				sum += weights[u]
			}
		}
		if sum > best {
			best = sum
		}
	}
	return best
}

// TestSolutionAgainstBruteforce generates random small forests and
// verifies both parts against a 2^n brute-force MIS and an exhaustive
// cross-tree hype sum.
func TestSolutionAgainstBruteforce(t *testing.T) {
	r := rand.New(rand.NewPCG(0x505, 0x505))
	for trial := 0; trial < 30; trial++ {
		n := 6 + r.IntN(9) // 6..14
		input, wantP1, wantP2 := randomForestInputS05(r, n)
		got := Solution([]byte(input))
		if got[0] != strconv.FormatInt(wantP1, 10) {
			t.Errorf("trial %d n=%d: part1 got %s want %d\ninput:\n%s",
				trial, n, got[0], wantP1, input)
		}
		if got[1] != strconv.FormatInt(wantP2, 10) {
			t.Errorf("trial %d n=%d: part2 got %s want %d\ninput:\n%s",
				trial, n, got[1], wantP2, input)
		}
	}
}
