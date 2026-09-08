package s01

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
			name:  "tiny",
			input: "1.00\n2.00\n4.00\n",
			// pairs summing to 7.00: none
			// signed subsets summing to 7.00: only {+1,+2,+4} -> 1
			want1: "0",
			want2: "1",
		},
		{
			name:  "example",
			input: "1.00\n2.00\n4.00\n3.50\n3.50",
			want1: "1",
			want2: "4",
		},
		{
			name:  "small mixed",
			input: "3.00\n4.00\n1.00\n2.00\n5.00\n",
			// pairs summing to 7.00: (3,4) and (2,5) -> 2
			// signed subsets summing to 7.00: 9 (verified by hand enumeration)
			want1: "2",
			want2: "9",
		},
		{
			name:  "with comments and blanks",
			input: "# header\n\n3.50\n3.50\n0.50\n6.50\n",
			// pairs summing to 7.00: (3.50,3.50) and (0.50,6.50) -> 2
			// signed subsets to 7.00 (enumerated by hand): 4
			//   {+6.50, +0.50}
			//   {+3.50, +3.50}
			//   {+6.50, +3.50, -3.50, +0.50} and its (a,b) swap of the two 3.50s -> 2
			want1: "2",
			want2: "4",
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

// brutePart1 enumerates every unordered pair by hand.
func brutePart1(cents []int) int64 {
	var n int64
	for i := 0; i < len(cents); i++ {
		for j := i + 1; j < len(cents); j++ {
			if cents[i]+cents[j] == 700 {
				n++
			}
		}
	}
	return n
}

// brutePart2 enumerates all 3^n skip/add/subtract assignments and counts
// those that sum to 700. Exponential, only usable for small n.
func brutePart2(cents []int) int64 {
	var count int64
	var walk func(i, sum int)
	walk = func(i, sum int) {
		if i == len(cents) {
			if sum == 700 {
				count++
			}
			return
		}
		walk(i+1, sum)
		walk(i+1, sum+cents[i])
		walk(i+1, sum-cents[i])
	}
	walk(0, 0)
	return count
}

func formatInputS01(cents []int) string {
	var b strings.Builder
	for _, c := range cents {
		neg := ""
		v := c
		if v < 0 {
			neg = "-"
			v = -v
		}
		fmt.Fprintf(&b, "%s%d.%02d\n", neg, v/100, v%100)
	}
	return b.String()
}

// TestSolutionAgainstBruteforce generates small random inputs whose values
// are all multiples of 25 cents in [-675, 675] (matching the generator's
// range) and compares the DP results against exhaustive enumeration.
func TestSolutionAgainstBruteforce(t *testing.T) {
	r := rand.New(rand.NewPCG(0xB01, 0xB01))
	for trial := 0; trial < 40; trial++ {
		n := 4 + r.IntN(9) // 4..12
		cents := make([]int, n)
		for i := range cents {
			// k in [1..27], sign uniform.
			k := 1 + r.IntN(27)
			v := 25 * k
			if r.IntN(2) == 0 {
				v = -v
			}
			cents[i] = v
		}
		input := formatInputS01(cents)
		got := Solution([]byte(input))
		want1 := brutePart1(cents)
		want2 := brutePart2(cents)
		if got[0] != strconv.FormatInt(want1, 10) {
			t.Errorf("trial %d n=%d cents=%v: part1 got %s want %d",
				trial, n, cents, got[0], want1)
		}
		if got[1] != strconv.FormatInt(want2, 10) {
			t.Errorf("trial %d n=%d cents=%v: part2 got %s want %d",
				trial, n, cents, got[1], want2)
		}
	}
}

func TestSolutionEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want1, want2 string
	}{
		{
			name:  "empty",
			input: "",
			want1: "0",
			want2: "0",
		},
		{
			name: "single 7.00",
			// One 7.00 tile: no pair; only the +7.00 selection hits target.
			input: "7.00\n",
			want1: "0",
			want2: "1",
		},
		{
			name: "three identical 3.50",
			// C(3,2)=3 pairs summing to 7.00; Part 2 also 3 (add any 2 of 3).
			input: "3.50\n3.50\n3.50\n",
			want1: "3",
			want2: "3",
		},
		{
			name: "four identical 3.50",
			// C(4,2)=6 pairs. Part 2: (add,sub)=(2,0):C(4,2)=6 + (3,1):C(4,3)*C(1,1)=4 = 10.
			input: "3.50\n3.50\n3.50\n3.50\n",
			want1: "6",
			want2: "10",
		},
		{
			name: "negative value forces negative t in DP",
			// Values -200, 500, 400. Signed subsets summing to 700:
			//   +(-200)+500+400 and -(-200)+500 (skip 400).
			input: "-2.00\n5.00\n4.00\n",
			want1: "0",
			want2: "2",
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

func TestSolutionRejectsCR(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on CR input, got none")
		}
		if !strings.Contains(fmt.Sprint(r), "CR in input") {
			t.Errorf("panic message: %v", r)
		}
	}()
	Solution([]byte("1.00\r\n2.00\r\n"))
}
