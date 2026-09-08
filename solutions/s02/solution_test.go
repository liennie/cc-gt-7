package s02

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSolutionDocExample(t *testing.T) {
	// Puzzle-text example: 2 groups with smolder 2 and 3.
	// Group 1: Ap → SmMn (2 ticks) → SmAx (1) → Ap (1). Cycle 4.
	// Group 2: Vn → HzLc (3 ticks) → HzAx (1) → Vn (1). Cycle 5.
	// First Ax at tick 3, position 2 → Part 1 = 6.
	// LCM(4,5) = 20 → Part 2 = 19.
	input := `ApVn

Ap 2> SmMn
SmMn 1> SmAx
SmAx 1> Ap
Vn 3> HzLc
HzLc 1> HzAx
HzAx 1> Vn
`
	got := Solution([]byte(input))
	if len(got) != 2 {
		t.Fatalf("want 2 answers, got %d", len(got))
	}
	if got[0] != "6" {
		t.Errorf("part1: got %q, want %q", got[0], "6")
	}
	if got[1] != "19" {
		t.Errorf("part2: got %q, want %q", got[1], "19")
	}
}

func TestSolutionSmolderOne(t *testing.T) {
	// Single-group case, all smolders = 1. Chain:
	// Fs=Ap → Sm Mn → Sm Ax → Ap. Cycle 3. First Ax at tick 2, position 2.
	// Part 1 = 4. Only 1 group, so Part 2 = LCM(3) - 1 = 2.
	input := `Ap

Ap 1> SmMn
SmMn 1> SmAx
SmAx 1> Ap
`
	got := Solution([]byte(input))
	if got[0] != "4" {
		t.Errorf("part1: got %q, want %q", got[0], "4")
	}
	if got[1] != "2" {
		t.Errorf("part2: got %q, want %q", got[1], "2")
	}
}

// TestBlendMatchesConcatenatedGroups asserts the generator's independent-
// group invariant: simulating the full initial blend produces, at every
// tick, exactly the concatenation of simulating each initial-blend token
// on its own. If any rule ever crossed group boundaries this would fail.
func TestBlendMatchesConcatenatedGroups(t *testing.T) {
	const ticks = 3000
	for i := 1; i <= 5; i++ {
		i := i
		t.Run(fmt.Sprintf("input%02d", i), func(t *testing.T) {
			path := filepath.Join("..", "..", "puzzles", "02", "inputs", fmt.Sprintf("%02d.txt", i))
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			initial, rules := parseInputForTest(data)
			if len(initial) == 0 {
				t.Fatalf("empty initial blend")
			}

			whole := newSim(append([]string(nil), initial...), rules)
			parts := make([]*sim, len(initial))
			for k, tok := range initial {
				parts[k] = newSim([]string{tok}, rules)
			}

			concat := func() []string {
				var out []string
				for _, p := range parts {
					out = append(out, p.blend...)
				}
				return out
			}
			check := func(tick int64) {
				if !slices.Equal(whole.blend, concat()) {
					t.Fatalf("tick %d: whole blend %v != concat of groups %v",
						tick, whole.blend, concat())
				}
			}

			check(0)
			whole.matchRules(0)
			for _, p := range parts {
				p.matchRules(0)
			}
			for tick := int64(1); tick <= ticks; tick++ {
				whole.applyCompletions(tick)
				for _, p := range parts {
					p.applyCompletions(tick)
				}
				check(tick)
				whole.matchRules(tick)
				for _, p := range parts {
					p.matchRules(tick)
				}
			}
		})
	}
}

// naivePart2 simulates the blend tick by tick until at least half of the
// compounds are Ax. Bounded by maxTicks — returns -1 if never reached.
// Purely for testing against the closed-form LCM-based Solution on tiny
// setups whose cycle LCM is small enough to enumerate.
func naivePart2(initial []string, rules []Rule, maxTicks int64) int64 {
	half := func(blend []string) bool {
		n := 0
		for _, t := range blend {
			if t == "Ax" {
				n++
			}
		}
		return n*2 >= len(blend)
	}
	s := newSim(append([]string(nil), initial...), rules)
	if half(s.blend) {
		return 0
	}
	s.matchRules(0)
	for tick := int64(1); tick <= maxTicks; tick++ {
		s.applyCompletions(tick)
		if half(s.blend) {
			return tick
		}
		s.matchRules(tick)
	}
	return -1
}

// naivePart1 simulates the blend tick by tick until any Ax appears.
func naivePart1(initial []string, rules []Rule, maxTicks int64) int64 {
	s := newSim(append([]string(nil), initial...), rules)
	if p := findAx(s.blend); p >= 0 {
		return 0
	}
	s.matchRules(0)
	for tick := int64(1); tick <= maxTicks; tick++ {
		s.applyCompletions(tick)
		if p := findAx(s.blend); p >= 0 {
			return tick * int64(p+1)
		}
		s.matchRules(tick)
	}
	return -1
}

// TestSolutionAgainstNaiveSimulation builds small multi-group rule sets
// with tiny cycle lengths and confirms Solution's LCM-based Part 2 (and
// its Part 1) matches straight tick-by-tick simulation.
func TestSolutionAgainstNaiveSimulation(t *testing.T) {
	// Each entry: (initial blend, per-group rule cycles). Every group
	// takes exactly its cycle length in ticks to return to Fs. Fs is a
	// distinct starting token; the intermediate tokens are unique per
	// group so the groups are independent.
	tests := []struct {
		name    string
		initial string
		rules   string
	}{
		{
			name:    "cycle 3 and 5",
			initial: "ApVn",
			rules: `Ap 1> SmMn
SmMn 1> SmAx
SmAx 1> Ap
Vn 3> HzLc
HzLc 1> HzAx
HzAx 1> Vn`,
		},
		{
			name:    "cycles 4 and 5 (doc example)",
			initial: "ApVn",
			rules: `Ap 2> SmMn
SmMn 1> SmAx
SmAx 1> Ap
Vn 3> HzLc
HzLc 1> HzAx
HzAx 1> Vn`,
		},
		{
			name:    "cycles 3 5 and 7",
			initial: "ApVnBt",
			rules: `Ap 1> SmMn
SmMn 1> SmAx
SmAx 1> Ap
Vn 3> HzLc
HzLc 1> HzAx
HzAx 1> Vn
Bt 5> CdEf
CdEf 1> CdAx
CdAx 1> Bt`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.initial + "\n\n" + tc.rules + "\n"
			got := Solution([]byte(input))
			initial, rules := parseInputForTest([]byte(input))
			want1 := naivePart1(initial, rules, 10_000)
			want2 := naivePart2(initial, rules, 10_000)
			if want1 < 0 || want2 < 0 {
				t.Fatalf("naive simulation didn't converge: p1=%d p2=%d", want1, want2)
			}
			if got[0] != fmt.Sprint(want1) {
				t.Errorf("part1: got %s, naive %d", got[0], want1)
			}
			if got[1] != fmt.Sprint(want2) {
				t.Errorf("part2: got %s, naive %d", got[1], want2)
			}
		})
	}
}

func parseInputForTest(input []byte) ([]string, []Rule) {
	var initial []string
	var rules []Rule
	seenBlend := false
	for _, raw := range strings.Split(string(input), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if !seenBlend {
			initial = parseCompounds(line)
			seenBlend = true
			continue
		}
		if r, ok := parseRule(line); ok {
			rules = append(rules, r)
		}
	}
	return initial, rules
}

func TestSolutionEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want1, want2 string
	}{
		{
			name: "cycle 2 single group",
			// [Fs] -> [Ax] -> [Fs]. Cycle 2. Ax at pos1=1 on tick 1.
			input: "Fs\n\nFs 1> Ax\nAx 1> Fs\n",
			want1: "1",
			want2: "1",
		},
		{
			name: "size-change with empty-RHS decay",
			// [Fs] -> [Ax, Tm] -> [Fs]. Decay Tm 1> and Ax 1> Fs fire in
			// parallel at tick 1; applyCompletions R->L handles both.
			input: "Fs\n\nFs 1> AxTm\nAx 1> Fs\nTm 1>\n",
			want1: "1",
			want2: "1",
		},
		{
			name: "size expansion 1->3",
			// [Fs] -> [Ax, Bp, Cp] -> [Fs]. Multi-token LHS on the return leg.
			input: "Fs\n\nFs 1> AxBpCp\nAxBpCp 1> Fs\n",
			want1: "1",
			want2: "1",
		},
		{
			name: "multi-group tick-0 fires, Ax at pos 2",
			// Two groups, cycles 3 and 2. At tick 0 both starters fire in
			// parallel; group 2's smolder-1 delivers its Ax at tick 1 to
			// position 2 (0-indexed 1) of the whole blend [Fs, Ax, Is].
			// Part 1 = tick * pos1 = 1 * 2 = 2. Part 2 = LCM(3,2) - 1 = 5.
			input: "FsGs\n\nFs 2> AxHs\nAxHs 1> Fs\nGs 1> AxIs\nAxIs 1> Gs\n",
			want1: "2",
			want2: "5",
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
	Solution([]byte("Fs\r\n\r\nFs 1> Ax\r\n"))
}
