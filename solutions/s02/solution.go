// Package s02 contains a solution to puzzle 02.
package s02

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Rule is a single rewrite rule: LHS transforms into RHS after smoldering
// for Smolder ticks.
type Rule struct {
	LHS     []string
	RHS     []string
	Smolder int
}

// Solution computes the answers to Puzzle 02.
//
// Input:
//  1. Initial blend, a single line of concatenated 2-letter compounds.
//  2. A blank line.
//  3. Rules, one per line, in the form "LHS N> RHS" where N is a positive
//     integer smolder time. LHS and RHS are concatenated 2-letter compounds.
//     RHS may be empty.
//
// Part 1 finds the first tick at which any Ax appears in the blend and
// returns tick × its 1-indexed position.
//
// Part 2 finds the first tick at which at least half of the compounds in
// the blend are Ax. In this puzzle each initial-blend token seeds an
// independent cycle; the answer is LCM(cycle lengths) − 1.
func Solution(input []byte) []string {
	if bytes.IndexByte(input, '\r') >= 0 {
		panic("s02: CR in input, expected LF-only")
	}
	lines := strings.Split(string(input), "\n")

	var initial []string
	var rules []Rule
	seenBlend := false
	for _, raw := range lines {
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

	part1 := simulatePart1(initial, rules)
	part2 := simulatePart2(initial, rules)

	return []string{
		strconv.FormatInt(part1, 10),
		strconv.FormatInt(part2, 10),
	}
}

// simulatePart1 runs the blend forward tick by tick until Ax first appears.
func simulatePart1(initial []string, rules []Rule) int64 {
	s := newSim(append([]string(nil), initial...), rules)
	if p := findAx(s.blend); p >= 0 {
		assertSingleAx(s.blend, 0)
		return 0
	}
	s.matchRules(0)
	for tick := int64(1); ; tick++ {
		s.applyCompletions(tick)
		if p := findAx(s.blend); p >= 0 {
			assertSingleAx(s.blend, tick)
			pos1 := int64(p + 1)
			ans := tick * pos1
			if pos1 != 0 && ans/pos1 != tick {
				panic("s02: int64 overflow computing tick*position")
			}
			return ans
		}
		s.matchRules(tick)
	}
}

// simulatePart2 exploits the puzzle's independent-group structure: each
// initial-blend token seeds a closed cycle. Returns LCM of cycle lengths − 1.
func simulatePart2(initial []string, rules []Rule) int64 {
	l := int64(1)
	for _, tok := range initial {
		l = lcm(l, cycleLength(tok, rules))
	}
	return l - 1
}

// cycleLength returns the number of ticks from initial state [tok] until
// the blend next returns to the single-token state [tok] with no rule
// currently smoldering.
func cycleLength(tok string, rules []Rule) int64 {
	s := newSim([]string{tok}, rules)
	s.matchRules(0)
	for tick := int64(1); ; tick++ {
		s.applyCompletions(tick)
		if len(s.blend) == 1 && s.blend[0] == tok && len(s.commits) == 0 {
			return tick
		}
		s.matchRules(tick)
	}
}

func findAx(blend []string) int {
	for i, t := range blend {
		if t == "Ax" {
			return i
		}
	}
	return -1
}

// assertSingleAx verifies the generator's invariant that Ax first appears
// alone. Solvers may take the leftmost Ax without checking, but this
// solution panics so any generator change that breaks the invariant is
// caught by TestEvent.
func assertSingleAx(blend []string, tick int64) {
	n := 0
	for _, t := range blend {
		if t == "Ax" {
			n++
		}
	}
	if n > 1 {
		panic(fmt.Sprintf("s02: %d Ax at first appearance tick %d, expected exactly 1", n, tick))
	}
}

// --- simulator ---

type commit struct {
	pos        int
	length     int
	rhs        []string
	completeAt int64
}

type sim struct {
	blend   []string
	rules   []Rule
	commits []commit
}

func newSim(blend []string, rules []Rule) *sim {
	return &sim{blend: blend, rules: rules}
}

// applyCompletions runs any commitments whose completeAt equals tick.
// It processes them right-to-left so that splicing doesn't shift the
// positions of pending completions to their left; positions of *not*
// completing commits to the right of a splice are adjusted by delta.
func (s *sim) applyCompletions(tick int64) {
	sort.Slice(s.commits, func(i, j int) bool { return s.commits[i].pos < s.commits[j].pos })
	for i := len(s.commits) - 1; i >= 0; i-- {
		c := s.commits[i]
		if c.completeAt != tick {
			continue
		}
		newBlend := make([]string, 0, len(s.blend)-c.length+len(c.rhs))
		newBlend = append(newBlend, s.blend[:c.pos]...)
		newBlend = append(newBlend, c.rhs...)
		newBlend = append(newBlend, s.blend[c.pos+c.length:]...)
		s.blend = newBlend
		delta := len(c.rhs) - c.length
		for j := range s.commits {
			if j != i && s.commits[j].pos > c.pos {
				s.commits[j].pos += delta
			}
		}
		s.commits = append(s.commits[:i], s.commits[i+1:]...)
	}
}

// matchRules scans uncommitted regions and commits every rule whose LHS
// matches. Because inputs are generated so that rules never overlap in a
// single tick, this panics if two matches would share any position.
func (s *sim) matchRules(tick int64) {
	mask := make([]bool, len(s.blend))
	for _, c := range s.commits {
		for k := c.pos; k < c.pos+c.length; k++ {
			if k >= 0 && k < len(mask) {
				mask[k] = true
			}
		}
	}
	type match struct {
		pos, length, smolder int
		rhs                  []string
	}
	var matches []match
	for pos := 0; pos < len(s.blend); pos++ {
		if mask[pos] {
			continue
		}
		for i := range s.rules {
			ru := &s.rules[i]
			l := len(ru.LHS)
			if pos+l > len(s.blend) {
				continue
			}
			collided := false
			for k := pos; k < pos+l; k++ {
				if mask[k] {
					collided = true
					break
				}
			}
			if collided {
				continue
			}
			ok := true
			for k := 0; k < l; k++ {
				if s.blend[pos+k] != ru.LHS[k] {
					ok = false
					break
				}
			}
			if ok {
				matches = append(matches, match{pos, l, ru.Smolder, ru.RHS})
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].pos != matches[j].pos {
			return matches[i].pos < matches[j].pos
		}
		return matches[i].length > matches[j].length
	})
	for i := 1; i < len(matches); i++ {
		if matches[i].pos < matches[i-1].pos+matches[i-1].length {
			panic(fmt.Sprintf(
				"s02: overlapping rule matches at tick %d: [pos %d len %d] and [pos %d len %d]",
				tick, matches[i-1].pos, matches[i-1].length,
				matches[i].pos, matches[i].length,
			))
		}
	}
	for _, m := range matches {
		s.commits = append(s.commits, commit{
			pos:        m.pos,
			length:     m.length,
			rhs:        m.rhs,
			completeAt: tick + int64(m.smolder),
		})
	}
}

// --- parsing ---

func parseRule(line string) (Rule, bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return Rule{}, false
	}
	if !strings.HasSuffix(fields[1], ">") {
		return Rule{}, false
	}
	n, err := strconv.Atoi(strings.TrimSuffix(fields[1], ">"))
	if err != nil || n < 1 {
		return Rule{}, false
	}
	lhs := parseCompounds(fields[0])
	if len(lhs) == 0 {
		return Rule{}, false
	}
	var rhs []string
	if len(fields) >= 3 {
		rhs = parseCompounds(fields[2])
		if rhs == nil {
			return Rule{}, false
		}
	}
	return Rule{LHS: lhs, RHS: rhs, Smolder: n}, true
}

func parseCompounds(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	if len(s)%2 != 0 {
		return nil
	}
	out := make([]string, len(s)/2)
	for i := range out {
		out[i] = s[i*2 : i*2+2]
	}
	return out
}

// --- math ---

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	q := a / gcd(a, b)
	// q * b overflow guard: with 5 primes near 200-280, the true product
	// fits comfortably (~1e12), but a broken generator that grew that
	// count could silently wrap. Fail loudly instead.
	p := q * b
	if b != 0 && p/b != q {
		panic("s02: int64 overflow computing LCM")
	}
	return p
}
