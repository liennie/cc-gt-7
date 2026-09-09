// Package g04 contains an input generator for puzzle 04.
package g04

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Generate produces the input for puzzle 04, seeded deterministically by idx.
//
// Structure (never surfaced to the puzzle solver):
//
//	5 independent "groups", each with a prime cycle length drawn from
//	{211..283}. Each group is a random walk over compound states, starting
//	at [Fs_i] and returning to [Fs_i] via a chain of rules that includes an
//	ash-producer step (whose RHS contains Ax) and an ash-consumer step
//	(whose LHS contains Ax and whose RHS is [Fs_i]).
//
// The chain keeps each group's state at [c, F] (one catalyst plus one
// fence) most of the time, connecting these single-catalyst states with
// three "simple" rule shapes — `C > C'`, `CF > C'F`, `CF > C'F'` — and
// occasional two-rule "excursions" that briefly expand the state to two
// or three catalysts before consolidating back to one. Excursions produce
// the growth/decay shapes: expansion (`C > CC`, `C > CCC`, `CF > CCF'`,
// …) paired with consolidation (`CC > D`, `CCF > D F`, `CCCF > D F'`, …).
//
// Every compound (starting token, catalyst, or former "fence") is drawn
// fresh from one global pool, so the 5 groups have strictly disjoint
// token alphabets and Ax is the only compound they share. Each token is
// consumed by exactly one rule, so no rule LHS ever appears as a
// substring of an earlier state — exactly one rule per group fires per
// tick. Group cycle length therefore equals the total number of rules in
// the group. That count is then padded up to L_i by distributing extra
// smolder ticks across a large majority of the rules with mildly
// randomised shares so that individual rule timings stay modest and
// spread out.
//
// The ash-consumer rule keeps smolder 1 so ash lives for exactly one
// tick per group per cycle.
func Generate(idx int) []byte {
	r := rand.New(rand.NewPCG(0x02_5eed_5eed, uint64(idx)+1))

	// 5 prime cycle lengths from a shared pool. LCM lands in ~10¹¹–10¹².
	primePool := []int{
		211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283,
	}
	r.Shuffle(len(primePool), func(i, j int) {
		primePool[i], primePool[j] = primePool[j], primePool[i]
	})
	Ls := primePool[:5]

	// One global pool of fresh 2-letter tokens serves every non-Ax draw:
	// each group's starting compound and every "catalyst" / "fence" token
	// it uses. Draws are strictly sequential from a shuffled pool, so the
	// 5 groups never share any token.
	catalystPool := makeCatalystPool()
	r.Shuffle(len(catalystPool), func(i, j int) {
		catalystPool[i], catalystPool[j] = catalystPool[j], catalystPool[i]
	})
	catIdx := 0
	drawCat := func() string {
		s := catalystPool[catIdx]
		catIdx++
		return s
	}
	// Kept as an alias so the buildGroup helpers can preserve their
	// original two-closure signatures without change.
	drawFence := drawCat

	startingFences := make([]string, 5)
	for i := range startingFences {
		startingFences[i] = drawCat()
	}

	var allRules []genRule
	for i := 0; i < 5; i++ {
		Fs := startingFences[i]
		L := Ls[i]
		group := buildGroup(r, Fs, L, drawCat, drawFence)
		allRules = append(allRules, group...)
	}

	// Shuffle every rule across the whole input.
	r.Shuffle(len(allRules), func(i, j int) {
		allRules[i], allRules[j] = allRules[j], allRules[i]
	})

	// Shuffle blend order.
	blendOrder := []int{0, 1, 2, 3, 4}
	r.Shuffle(len(blendOrder), func(i, j int) {
		blendOrder[i], blendOrder[j] = blendOrder[j], blendOrder[i]
	})

	var b strings.Builder
	for _, i := range blendOrder {
		b.WriteString(startingFences[i])
	}
	b.WriteByte('\n')
	b.WriteByte('\n')
	for _, ru := range allRules {
		b.WriteString(strings.Join(ru.lhs, ""))
		fmt.Fprintf(&b, " %d>", ru.smolder)
		if len(ru.rhs) > 0 {
			b.WriteByte(' ')
			b.WriteString(strings.Join(ru.rhs, ""))
		}
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

type genRule struct {
	lhs, rhs []string
	smolder  int
}

// buildGroup constructs one group's rule set: a walk from [Fs] back to
// [Fs] whose length (= number of rules) is boosted to the target prime L.
func buildGroup(
	r *rand.Rand, Fs string, L int,
	drawCat, drawFence func() string,
) []genRule {
	var rules []genRule

	// Step 0: Fs → [c_1, ..., c_n, F_1] with n = 1..3. Larger n biases the
	// chain toward an expansion-flavored opening.
	n0 := 1 + r.IntN(3)
	firstF := drawFence()
	init := make([]string, 0, n0+1)
	for i := 0; i < n0; i++ {
		init = append(init, drawCat())
	}
	init = append(init, firstF)
	rules = append(rules, genRule{
		lhs: []string{Fs}, rhs: append([]string(nil), init...), smolder: 1,
	})
	state := append([]string(nil), init...)

	// If the starting rule produced more than one catalyst, sprinkle a few
	// strict-subgroup bubbles (rewriting one catalyst at a time in place,
	// LHS narrower than the full catalyst slice) before consolidating.
	if len(state) > 2 {
		var bubs []genRule
		state, bubs = bubbleChain(r, state, drawCat, 1+r.IntN(3))
		rules = append(rules, bubs...)
		var ru genRule
		state, ru = consolidateStep(r, state, drawCat, drawFence)
		rules = append(rules, ru)
	}

	// Middle phase: mix of simple steps (1 rule), short excursions (2
	// rules with a temporary 2- or 3-catalyst state), and long excursions
	// (3 rules that pass through a mid walk-two, walk-three, or middle
	// decay stage).
	steps := 10 + r.IntN(6) // 10..15 middle steps
	for k := 0; k < steps; k++ {
		phase := float64(k) / float64(steps)
		var newRules []genRule
		state, newRules = pickAndAdvance(r, phase, state, drawCat, drawFence)
		rules = append(rules, newRules...)
	}

	// Ash producer: state == [C, F]. Emit `C F > finalC Ax`.
	finalC := drawCat()
	rules = append(rules, genRule{
		lhs:     []string{state[0], state[1]},
		rhs:     []string{finalC, "Ax"},
		smolder: 1,
	})

	// Ash consumer: MUST stay at smolder 1 (ash lifetime = 1 tick).
	rules = append(rules, genRule{
		lhs:     []string{finalC, "Ax"},
		rhs:     []string{Fs},
		smolder: 1,
	})

	distributeSmolder(r, rules, L)
	return rules
}

// simpleStep advances a single-catalyst state [C, F] by one rule of shape
// `C > C'`, `CF > C'F`, or `CF > C'F'`. Returns the new state and the rule.
func simpleStep(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, genRule) {
	C, F := state[0], state[1]
	D := drawCat()
	switch pickSimpleKind(r, phase) {
	case simpleWalk:
		return []string{D, F}, genRule{
			lhs: []string{C}, rhs: []string{D}, smolder: 1,
		}
	case simpleWalkNearFence:
		return []string{D, F}, genRule{
			lhs: []string{C, F}, rhs: []string{D, F}, smolder: 1,
		}
	default: // simpleFenceTrans
		F2 := drawFence()
		return []string{D, F2}, genRule{
			lhs: []string{C, F}, rhs: []string{D, F2}, smolder: 1,
		}
	}
}

// excursionStep advances a single-catalyst state [C, F] via a pair of
// rules: an expansion (state grows to 2 or 3 catalysts) followed by a
// consolidation (state returns to a single catalyst).
func excursionStep(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, genRule, genRule) {
	C, F := state[0], state[1]

	// Choose intermediate catalyst count: 2 by default, 3 in early phase.
	excSize := 2
	if phase < 0.5 && r.IntN(3) == 0 {
		excSize = 3
	}
	inc := make([]string, excSize)
	for i := range inc {
		inc[i] = drawCat()
	}

	// Expansion: three shapes. Kind 0 keeps LHS at just [C]; kinds 1 and 2
	// include the fence in LHS (with kind 2 also transitioning it).
	var expRule genRule
	var midFence string
	switch r.IntN(3) {
	case 0:
		expRule = genRule{
			lhs: []string{C}, rhs: append([]string(nil), inc...), smolder: 1,
		}
		midFence = F
	case 1:
		rhs := append(append([]string(nil), inc...), F)
		expRule = genRule{lhs: []string{C, F}, rhs: rhs, smolder: 1}
		midFence = F
	default:
		F2 := drawFence()
		rhs := append(append([]string(nil), inc...), F2)
		expRule = genRule{lhs: []string{C, F}, rhs: rhs, smolder: 1}
		midFence = F2
	}

	// Consolidation: three shapes mirroring the expansion.
	D := drawCat()
	var consRule genRule
	var newFence string
	switch r.IntN(3) {
	case 0:
		consRule = genRule{
			lhs: append([]string(nil), inc...), rhs: []string{D}, smolder: 1,
		}
		newFence = midFence
	case 1:
		lhs := append(append([]string(nil), inc...), midFence)
		consRule = genRule{
			lhs: lhs, rhs: []string{D, midFence}, smolder: 1,
		}
		newFence = midFence
	default:
		F2 := drawFence()
		lhs := append(append([]string(nil), inc...), midFence)
		consRule = genRule{
			lhs: lhs, rhs: []string{D, F2}, smolder: 1,
		}
		newFence = F2
	}

	return []string{D, newFence}, expRule, consRule
}

// consolidateStep collapses a state with more than one catalyst back to
// [D, F] via one rule (either dropping catalysts with fence preserved, or
// dropping catalysts with a fence transition). Used to normalise the
// state right after a multi-catalyst starting rule.
func consolidateStep(
	r *rand.Rand, state []string,
	drawCat, drawFence func() string,
) ([]string, genRule) {
	F := state[len(state)-1]
	cats := state[:len(state)-1]
	D := drawCat()
	if r.IntN(2) == 0 {
		return []string{D, F}, genRule{
			lhs: append([]string(nil), cats...), rhs: []string{D}, smolder: 1,
		}
	}
	F2 := drawFence()
	lhs := append([]string(nil), state...)
	return []string{D, F2}, genRule{
		lhs: lhs, rhs: []string{D, F2}, smolder: 1,
	}
}

// pickAndAdvance dispatches to a step kind for a single middle-phase
// iteration, returning the new state and the rules emitted.
func pickAndAdvance(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	x := r.IntN(100)
	switch {
	case x < 30:
		s, ru := simpleStep(r, phase, state, drawCat, drawFence)
		return s, []genRule{ru}
	case x < 55:
		s, exp, cons := excursionStep(r, phase, state, drawCat, drawFence)
		return s, []genRule{exp, cons}
	default:
		return longExcursionStep(r, phase, state, drawCat, drawFence)
	}
}

// longExcursionStep advances a single-catalyst state via a 3-to-8-rule
// chain that visits a temporary 2- or 3-catalyst state. Sub-kinds:
//   - walk2:      C > C1C2       [bubbles]  C1C2 > D E    [bubbles]  D E > X
//   - walk3:      C > C1C2C3     [bubbles]  C1C2C3 > DEF  [bubbles]  DEF > X
//   - decay3:     C > C1C2C3     C2 >       [bubbles]     C1C3 > D
//   - bubble:     mid-shift, right-fence, wait-decay, or pure same-shape chain
//
// Each sub-step independently picks a fence variant (LHS excludes fence,
// LHS includes fence and keeps it, or LHS includes fence and transitions).
func longExcursionStep(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	switch r.IntN(4) {
	case 0:
		return longWalkExcursion(r, phase, state, 2, drawCat, drawFence)
	case 1:
		return longWalkExcursion(r, phase, state, 3, drawCat, drawFence)
	case 2:
		return decayExcursion(r, phase, state, drawCat, drawFence)
	default:
		return bubbleExcursion(r, phase, state, drawCat, drawFence)
	}
}

// longWalkExcursion emits: expand C to `size` fresh catalysts, optionally
// bubble a few of them in place, walk to `size` fresh catalysts (possibly
// with a fence transition), optionally bubble again, then consolidate.
func longWalkExcursion(
	r *rand.Rand, phase float64, state []string, size int,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]

	mid1 := make([]string, size)
	for i := range mid1 {
		mid1[i] = drawCat()
	}
	expRule, midFence1 := expandFromOne(r, C, F, mid1, drawFence)
	rules := []genRule{expRule}

	work := append(append([]string(nil), mid1...), midFence1)
	work, bubs := maybeBubbleChain(r, work, drawCat, 0, 3)
	rules = append(rules, bubs...)
	mid1Now := append([]string(nil), work[:size]...)

	mid2 := make([]string, size)
	for i := range mid2 {
		mid2[i] = drawCat()
	}
	walkRule, midFence2 := walkSameSize(r, mid1Now, midFence1, mid2, drawFence)
	rules = append(rules, walkRule)

	work = append(append([]string(nil), mid2...), midFence2)
	work, bubs = maybeBubbleChain(r, work, drawCat, 0, 3)
	rules = append(rules, bubs...)
	mid2Now := append([]string(nil), work[:size]...)

	D := drawCat()
	consRule, finalFence := consolidateToOne(r, mid2Now, midFence2, D, drawFence)
	rules = append(rules, consRule)

	return []string{D, finalFence}, rules
}

// decayExcursion emits three rules: expand C to 3 fresh catalysts, decay
// the middle catalyst (empty RHS), then consolidate the two survivors to
// a single fresh catalyst. Decaying the MIDDLE catalyst is required so
// that the consolidation rule's LHS [C1, C3] is not a contiguous
// substring of the pre-decay state [C1, C2, C3, F].
func decayExcursion(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]

	mid := []string{drawCat(), drawCat(), drawCat()}
	expRule, midFence := expandFromOne(r, C, F, mid, drawFence)
	rules := []genRule{expRule}

	decayRule := genRule{lhs: []string{mid[1]}, rhs: nil, smolder: 1}
	rules = append(rules, decayRule)

	D := drawCat()
	surv := []string{mid[0], mid[2]}
	consRule, finalFence := consolidateToOne(r, surv, midFence, D, drawFence)
	rules = append(rules, consRule)

	return []string{D, finalFence}, rules
}

// bubbleExcursion emits a 4- to 7-rule cascade that keeps the state at 3
// catalysts across multiple ticks while mutating only part of it at each
// step. Every intermediate rule's LHS includes a token freshly produced
// by the previous rule, which guarantees strict tick-by-tick sequencing
// (no earlier state matches the rule, no later state does either).
//
// Sub-kinds:
//   - midShift:    C > c1c2c3   c2 > c2'         c2'c3 > c2'c3'    (consolidate)
//   - rightFence:  C > c1c2c3   c3F > c3'F'      c2c3' > c2'c3”   (consolidate)
//   - waitDecay:   C > c1c2c3   c2 > c2'         c2' > (empty)     c1c3 > D
//   - sameShape:   C > c1c2c3   [2-4 same-shape bubbles]           (consolidate)
func bubbleExcursion(
	r *rand.Rand, phase float64, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	switch r.IntN(4) {
	case 0:
		return bubbleMidShift(r, state, drawCat, drawFence)
	case 1:
		return bubbleRightFence(r, state, drawCat, drawFence)
	case 2:
		return bubbleWaitDecay(r, state, drawCat, drawFence)
	default:
		return bubbleSameShape(r, state, drawCat, drawFence)
	}
}

// bubbleMidShift:
//
//	R0: C >          c1 c2 c3      state: [c1 c2 c3 F]
//	R1: c2 >         c2'            state: [c1 c2' c3 F]
//	R2: c2' c3 >     c2' c3'        state: [c1 c2' c3' F]
//	R3: c1 c2' c3' > D              state: [D F']
//
// R2 anchors on c2' (only exists post-R1). R3 anchors on both c2' and c3'
// (c3' only exists post-R2). So each rule uniquely fires in its own tick.
func bubbleMidShift(
	r *rand.Rand, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]
	c1, c2, c3 := drawCat(), drawCat(), drawCat()
	c2p, c3p := drawCat(), drawCat()
	D := drawCat()

	expRule, midFence := expandFromOne(r, C, F, []string{c1, c2, c3}, drawFence)
	r1 := genRule{lhs: []string{c2}, rhs: []string{c2p}, smolder: 1}
	r2 := genRule{lhs: []string{c2p, c3}, rhs: []string{c2p, c3p}, smolder: 1}
	consRule, finalFence := consolidateToOne(
		r, []string{c1, c2p, c3p}, midFence, D, drawFence,
	)
	return []string{D, finalFence}, []genRule{expRule, r1, r2, consRule}
}

// bubbleRightFence:
//
//	R0: C >          c1 c2 c3      state: [c1 c2 c3 F]
//	R1: c3 F >       c3' F'         state: [c1 c2 c3' F']
//	R2: c2 c3' >     c2' c3''       state: [c1 c2' c3'' F']
//	R3: c1 c2' c3'' > D             state: [D F']
//
// R1 is a right-end rule that also swaps the fence. R2 anchors on c3'
// (post-R1), R3 anchors on c2' and c3” (post-R2).
func bubbleRightFence(
	r *rand.Rand, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]
	c1, c2, c3 := drawCat(), drawCat(), drawCat()
	c3p := drawCat()
	Fp := drawFence()
	c2p, c3pp := drawCat(), drawCat()
	D := drawCat()

	expRule, midFence := expandFromOne(r, C, F, []string{c1, c2, c3}, drawFence)
	r1 := genRule{
		lhs: []string{c3, midFence}, rhs: []string{c3p, Fp}, smolder: 1,
	}
	r2 := genRule{
		lhs: []string{c2, c3p}, rhs: []string{c2p, c3pp}, smolder: 1,
	}
	consRule, finalFence := consolidateToOne(
		r, []string{c1, c2p, c3pp}, Fp, D, drawFence,
	)
	return []string{D, finalFence}, []genRule{expRule, r1, r2, consRule}
}

// bubbleWaitDecay:
//
//	R0: C >          c1 c2 c3      state: [c1 c2 c3 F]
//	R1: c2 >         c2'            state: [c1 c2' c3 F]
//	R2: c2' >        (empty)        state: [c1 c3 F]
//	R3: c1 c3 >      D              state: [D F']
//
// R3's LHS [c1 c3] is not contiguous in any earlier state (c2 or c2' sits
// between them). R2 uses c2' as anchor so it can't fire before R1. This
// is the "rule that waits for a middle catalyst to decay" pattern.
func bubbleWaitDecay(
	r *rand.Rand, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]
	c1, c2, c3 := drawCat(), drawCat(), drawCat()
	c2p := drawCat()
	D := drawCat()

	expRule, midFence := expandFromOne(r, C, F, []string{c1, c2, c3}, drawFence)
	r1 := genRule{lhs: []string{c2}, rhs: []string{c2p}, smolder: 1}
	r2 := genRule{lhs: []string{c2p}, rhs: nil, smolder: 1}
	consRule, finalFence := consolidateToOne(
		r, []string{c1, c3}, midFence, D, drawFence,
	)
	return []string{D, finalFence}, []genRule{expRule, r1, r2, consRule}
}

// bubbleSameShape:
//
//	R0: C >          c1 c2 c3       state: [c1 c2 c3 F]
//	R1..Rk: strict-subgroup bubbles, each rewrites ONE catalyst in place
//	Rn: <consolidate current cats> > D
//
// The pure "bubbling" excursion: 3-5 in-place single-catalyst rewrites on
// a 3-catalyst state, each with a strict-subgroup LHS (size 1 or 2 of 3),
// before consolidating.
func bubbleSameShape(
	r *rand.Rand, state []string,
	drawCat, drawFence func() string,
) ([]string, []genRule) {
	C, F := state[0], state[1]
	c1, c2, c3 := drawCat(), drawCat(), drawCat()
	expRule, midFence := expandFromOne(r, C, F, []string{c1, c2, c3}, drawFence)
	rules := []genRule{expRule}

	work := []string{c1, c2, c3, midFence}
	var bubs []genRule
	work, bubs = bubbleChain(r, work, drawCat, 3+r.IntN(3))
	rules = append(rules, bubs...)

	D := drawCat()
	cats := []string{work[0], work[1], work[2]}
	consRule, finalFence := consolidateToOne(r, cats, midFence, D, drawFence)
	rules = append(rules, consRule)
	return []string{D, finalFence}, rules
}

// subgroupBubble emits one bubble rule whose LHS is a strict subgroup of
// the catalyst slice (size 1..n-1). The window must contain `anchor`
// (the position freshly rewritten by the immediately preceding rule), so
// that the LHS uniquely matches only the current tick. If anchor < 0 the
// window may sit anywhere in the slice (used for the first bubble in a
// chain, right after every catalyst was freshly introduced).
//
// state must be [c0..c_{n-1}, F] with n >= 2. Returns the new state, the
// position of the freshly-introduced catalyst (the new anchor), and the
// rule.
func subgroupBubble(
	r *rand.Rand, state []string, anchor int, drawCat func() string,
) ([]string, int, genRule) {
	n := len(state) - 1
	sizeMax := n - 1
	if sizeMax < 1 {
		sizeMax = 1
	}
	size := 1 + r.IntN(sizeMax)
	var lo int
	if anchor < 0 {
		lo = r.IntN(n - size + 1)
	} else {
		lo0 := anchor - size + 1
		if lo0 < 0 {
			lo0 = 0
		}
		lo1 := anchor
		if lo1 > n-size {
			lo1 = n - size
		}
		lo = lo0 + r.IntN(lo1-lo0+1)
	}
	hi := lo + size - 1
	p := lo + r.IntN(size)
	fresh := drawCat()
	lhs := append([]string(nil), state[lo:hi+1]...)
	rhs := append([]string(nil), state[lo:hi+1]...)
	rhs[p-lo] = fresh
	newState := append([]string(nil), state...)
	newState[p] = fresh
	return newState, p, genRule{lhs: lhs, rhs: rhs, smolder: 1}
}

// bubbleChain emits `count` strict-subgroup bubbles in a row. The first
// bubble may sit anywhere in the catalyst slice; each subsequent bubble
// must overlap the previously-rewritten position. For n == 2 the window
// is always size 1 and the anchor pins successive bubbles to that one
// position (a repeated single-catalyst rotation).
func bubbleChain(
	r *rand.Rand, state []string, drawCat func() string, count int,
) ([]string, []genRule) {
	if len(state) < 3 || count < 1 {
		return state, nil
	}
	rules := make([]genRule, 0, count)
	anchor := -1
	for i := 0; i < count; i++ {
		var ru genRule
		state, anchor, ru = subgroupBubble(r, state, anchor, drawCat)
		rules = append(rules, ru)
	}
	return state, rules
}

// maybeBubbleChain picks a bubble count in [minB, maxB] and emits that
// many strict-subgroup bubbles via bubbleChain.
func maybeBubbleChain(
	r *rand.Rand, state []string, drawCat func() string, minB, maxB int,
) ([]string, []genRule) {
	n := minB
	if maxB > minB {
		n += r.IntN(maxB - minB + 1)
	}
	return bubbleChain(r, state, drawCat, n)
}

// expandFromOne emits a rule that turns [C, F] into [inc..., midFence].
// LHS is either [C] (fence untouched) or [C, F] (fence kept or transitioned).
func expandFromOne(
	r *rand.Rand, C, F string, inc []string,
	drawFence func() string,
) (genRule, string) {
	switch r.IntN(3) {
	case 0:
		return genRule{
			lhs: []string{C}, rhs: append([]string(nil), inc...), smolder: 1,
		}, F
	case 1:
		rhs := append(append([]string(nil), inc...), F)
		return genRule{lhs: []string{C, F}, rhs: rhs, smolder: 1}, F
	default:
		F2 := drawFence()
		rhs := append(append([]string(nil), inc...), F2)
		return genRule{lhs: []string{C, F}, rhs: rhs, smolder: 1}, F2
	}
}

// walkSameSize emits a rule that turns [mid1..., inFence] into
// [mid2..., outFence], where len(mid1) == len(mid2). LHS is either
// [mid1...] (fence untouched) or [mid1..., inFence] (fence kept or
// transitioned).
func walkSameSize(
	r *rand.Rand, mid1 []string, inFence string, mid2 []string,
	drawFence func() string,
) (genRule, string) {
	switch r.IntN(3) {
	case 0:
		return genRule{
			lhs:     append([]string(nil), mid1...),
			rhs:     append([]string(nil), mid2...),
			smolder: 1,
		}, inFence
	case 1:
		lhs := append(append([]string(nil), mid1...), inFence)
		rhs := append(append([]string(nil), mid2...), inFence)
		return genRule{lhs: lhs, rhs: rhs, smolder: 1}, inFence
	default:
		F2 := drawFence()
		lhs := append(append([]string(nil), mid1...), inFence)
		rhs := append(append([]string(nil), mid2...), F2)
		return genRule{lhs: lhs, rhs: rhs, smolder: 1}, F2
	}
}

// consolidateToOne emits a rule that turns [cats..., inFence] into
// [D, outFence]. LHS is either [cats...] (fence untouched) or
// [cats..., inFence] (fence kept or transitioned).
func consolidateToOne(
	r *rand.Rand, cats []string, inFence, D string,
	drawFence func() string,
) (genRule, string) {
	switch r.IntN(3) {
	case 0:
		return genRule{
			lhs: append([]string(nil), cats...), rhs: []string{D}, smolder: 1,
		}, inFence
	case 1:
		lhs := append(append([]string(nil), cats...), inFence)
		return genRule{
			lhs: lhs, rhs: []string{D, inFence}, smolder: 1,
		}, inFence
	default:
		F2 := drawFence()
		lhs := append(append([]string(nil), cats...), inFence)
		return genRule{lhs: lhs, rhs: []string{D, F2}, smolder: 1}, F2
	}
}

type simpleKind int

const (
	simpleWalk simpleKind = iota
	simpleWalkNearFence
	simpleFenceTrans
)

// pickSimpleKind returns one of the three simple-step kinds, weighted by
// phase (fence transitions are somewhat more common in mid-phase).
func pickSimpleKind(r *rand.Rand, phase float64) simpleKind {
	var w [3]int
	switch {
	case phase < 0.33:
		w = [3]int{4, 3, 3}
	case phase < 0.7:
		w = [3]int{3, 3, 4}
	default:
		w = [3]int{4, 4, 2}
	}
	total := w[0] + w[1] + w[2]
	x := r.IntN(total)
	for i, wt := range w {
		if x < wt {
			return simpleKind(i)
		}
		x -= wt
	}
	return simpleWalk
}

// distributeSmolder spreads L − numRules extra ticks over a large majority
// of the non-consumer rules using randomised shares, so that timings stay
// modest and are visible across most of the chain instead of piling up on
// three or four outlier rules.
func distributeSmolder(r *rand.Rand, rules []genRule, L int) {
	numRules := len(rules)
	extra := L - numRules
	if extra <= 0 {
		return
	}
	eligibleCount := numRules - 1 // exclude the ash consumer at the end
	// Boost 65-90% of eligible rules.
	numBoost := (eligibleCount*(65+r.IntN(26)) + 99) / 100
	if numBoost < 3 {
		numBoost = 3
	}
	if numBoost > eligibleCount {
		numBoost = eligibleCount
	}
	eligible := make([]int, eligibleCount)
	for i := range eligible {
		eligible[i] = i
	}
	r.Shuffle(len(eligible), func(a, b int) {
		eligible[a], eligible[b] = eligible[b], eligible[a]
	})
	boost := eligible[:numBoost]

	weights := make([]int, numBoost)
	totalW := 0
	for i := range weights {
		weights[i] = 3 + r.IntN(8) // random weight 3..10
		totalW += weights[i]
	}
	used := 0
	for i, idx := range boost {
		var share int
		if i == numBoost-1 {
			share = extra - used
		} else {
			share = (extra * weights[i]) / totalW
		}
		used += share
		rules[idx].smolder += share
	}
}

// makeCatalystPool returns a list of pronounceable-ish 2-letter tokens
// (uppercase letter followed by lowercase letter), excluding Ax. The
// pool must be large enough for every fresh token drawn across all
// groups (starting compounds, catalysts, and former fences).
func makeCatalystPool() []string {
	reserved := map[string]bool{"Ax": true}
	const (
		upper = "BCDFGHJKLMNPQRSTVWXYZ"
		lower = "abcdefghijklmnopqrstuvwxyz"
	)
	var pool []string
	add := func(s string) {
		if !reserved[s] {
			reserved[s] = true
			pool = append(pool, s)
		}
	}
	for _, u := range upper {
		for _, l := range lower {
			add(string(u) + string(l))
		}
	}
	const (
		uVowels     = "AEIOU"
		lConsonants = "bdfglmnprstvz"
	)
	for _, v := range uVowels {
		for _, c := range lConsonants {
			add(string(v) + string(c))
		}
	}
	return pool
}
