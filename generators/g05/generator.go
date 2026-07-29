// Package g05 contains an input generator for puzzle 05.
package g05

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type topicSpec struct {
	slug string
	cat  string
}

// GT-lore topics. Category controls which suffixes/prefixes may combine with
// each topic (e.g. `unnatural_worlds` is still in development, so no
// `postmortem`; `sauna` and `pond` are places, so no `speedrun`).
var topics = []topicSpec{
	{"bulanci", "game"},
	{"factorio", "game"},
	{"warcraft_3", "game"},
	{"aoe_2", "game"},
	{"gt_nano", "game"},
	{"rts_balance", "game"},
	{"space_tycoon", "game"},
	{"dungeons_and_trolls", "game"},
	{"az_quiz", "game"},
	{"microbots", "game"},
	{"vibe_coding", "game"},
	{"jj_vcs", "game"},
	{"ebpf", "game"},
	{"kcd", "game"},
	{"poe", "game"},

	{"unnatural_worlds", "indev"},
	{"chores_bot", "indev"},

	{"paintball", "outdoor"},
	{"archery", "outdoor"},
	{"roundnet", "outdoor"},
	{"night_game", "outdoor"},
	{"frisbee", "outdoor"},
	{"larp", "outdoor"},

	{"bbq", "food"},
	{"hookah", "food"},
	{"shisha", "food"},
	{"pizza_night", "food"},
	{"pancakes", "food"},

	{"sauna", "place"},
	{"bunker", "place"},
	{"chalupa", "place"},
	{"pond", "place"},
	{"kids_corner", "place"},
	{"whirlpool", "place"},
	{"off_grid_sauna", "place"},

	{"cipher_hunt", "event"},
	{"puzzle_hunt", "event"},
	{"board_games", "event"},
	{"barcamp", "event"},
	{"invite_cipher", "event"},
	{"escape_game", "event"},
	{"trivia", "event"},
	{"lan_party", "event"},
	{"soldering", "event"},
	{"world_travel", "event"},
	{"guess_the_ingredient", "event"},
	{"house_renovation", "event"},
}

type fragmentSpec struct {
	slug string
	// Categories this fragment may combine with. nil means all.
	allowed map[string]bool
}

func catSet(cs ...string) map[string]bool {
	m := make(map[string]bool, len(cs))
	for _, c := range cs {
		m[c] = true
	}
	return m
}

// Trailing fragments joined after a topic.
var suffixes = []fragmentSpec{
	{"deep_dive", nil},
	{"101", nil},
	{"retrospective", catSet("game", "outdoor", "food", "event")},
	{"postmortem", catSet("game", "event")},
	{"lightning_talk", nil},
	{"field_notes", nil},
	{"anti_patterns", catSet("game", "event")},
	{"cheat_sheet", catSet("game", "indev", "outdoor", "food", "event")},
	{"war_stories", catSet("game", "indev", "outdoor", "event")},
	{"crash_course", nil},
	{"speedrun", catSet("game", "indev", "outdoor", "event")},
	{"hot_takes", nil},
	{"from_scratch", catSet("game", "indev", "food", "event")},
	{"considered_harmful", catSet("game", "event")},
	{"done_right", nil},
	{"masterclass", nil},
	{"in_the_wild", nil},
	{"demystified", nil},
	{"bloopers", nil},
	{"patch_notes", catSet("game", "indev", "event")},
	{"office_hours", nil},
	{"workshop", catSet("game", "indev", "food")},
	{"trivia", nil},
}

// Leading fragments joined before a topic.
var prefixes = []fragmentSpec{
	{"advanced", catSet("game", "indev", "outdoor", "food", "event")},
	{"introducing", nil},
	{"practical", catSet("game", "indev", "outdoor", "food", "event")},
	{"rethinking", catSet("game", "indev", "outdoor", "food", "event")},
	{"weekend", nil},
	{"forgotten", nil},
	{"ultimate", nil},
	{"3_slides_of", nil},
	{"zen_of", nil},
	{"state_of", nil},
	{"beyond", nil},
	{"secret_life_of", nil},
}

// Generate produces the input for puzzle 05, seeded deterministically by idx.
//
// The input has three blank-line-separated sections:
//  1. HYPE     : "<slug> <int>" per talk, hype score in [1, 1000].
//  2. CONFLICTS: "<slug> <slug>" per conflict edge. The conflict graph is
//     always a forest of exactly one large tree of 48..54 talks
//     plus 3..5 smaller trees, each 4..10 talks.
//  3. QUERIES  : "<slug> <slug>" per query pair, mixing same-tree and
//     cross-tree pairs.
//
// Talk slugs are unique lowercase identifiers built from GT-lore fragments;
// see the topic/suffix/prefix tables for the compatibility rules that keep
// silly combos (e.g. `bbq_speedrun`) out of the pool.
func Generate(idx int) []byte {
	r := rand.New(rand.NewPCG(0x05_5eed_5eed, uint64(idx)+1))

	const (
		nTalks   = 80
		nQueries = 20
	)

	// Tree sizes: one big tree, then small trees until the budget runs out.
	// The last small tree absorbs any <4 tail so every tree has >=4 nodes.
	bigSize := 48 + r.IntN(7)
	sizes := []int{bigSize}
	remaining := nTalks - bigSize
	for remaining > 0 {
		s := 4 + r.IntN(7)
		if s > remaining || remaining-s < 4 {
			s = remaining
		}
		sizes = append(sizes, s)
		remaining -= s
	}

	titles := composeTitles(r, nTalks)

	// Permute title-to-node-index mapping so slug order doesn't leak the
	// tree grouping.
	perm := r.Perm(nTalks)
	names := make([]string, nTalks)
	for i, p := range perm {
		names[p] = titles[i]
	}

	type edge struct{ a, b int }
	var edges []edge
	parent := make([]int, nTalks)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[ra] = rb
		}
	}

	// Build each tree as a random recursive tree over a contiguous block of
	// node indices; node k attaches to a uniform random earlier node in the
	// same block.
	trees := make([][]int, len(sizes))
	base := 0
	for ti, size := range sizes {
		nodes := make([]int, size)
		for i := 0; i < size; i++ {
			nodes[i] = base + i
			if i > 0 {
				p := base + r.IntN(i)
				edges = append(edges, edge{p, nodes[i]})
				union(p, nodes[i])
			}
		}
		trees[ti] = nodes
		base += size
	}
	r.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })

	weights := make([]int, nTalks)
	for i := range weights {
		weights[i] = 1 + r.IntN(1000)
	}

	type qry struct{ a, b int }
	var queries []qry
	seen := map[[2]int]bool{}
	addQuery := func(a, b int) bool {
		if a == b {
			return false
		}
		key := [2]int{a, b}
		if a > b {
			key = [2]int{b, a}
		}
		if seen[key] {
			return false
		}
		seen[key] = true
		queries = append(queries, qry{a, b})
		return true
	}

	// Randomise the same/cross split so Part 1 answer varies between inputs.
	nSame := 6 + r.IntN(nQueries-11)
	attempts := 0
	for len(queries) < nSame && attempts < 10000 {
		attempts++
		nodes := trees[r.IntN(len(trees))]
		addQuery(nodes[r.IntN(len(nodes))], nodes[r.IntN(len(nodes))])
	}
	attempts = 0
	for len(queries) < nQueries && attempts < 10000 {
		attempts++
		a := r.IntN(nTalks)
		b := r.IntN(nTalks)
		if find(a) == find(b) {
			continue
		}
		addQuery(a, b)
	}
	r.Shuffle(len(queries), func(i, j int) { queries[i], queries[j] = queries[j], queries[i] })

	var b strings.Builder
	for i, n := range names {
		fmt.Fprintf(&b, "%s %d\n", n, weights[i])
	}
	b.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&b, "%s %s\n", names[e.a], names[e.b])
	}
	b.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&b, "%s %s\n", names[q.a], names[q.b])
	}
	return []byte(b.String())
}

// composeTitles returns n unique talk-title slugs. About 75% follow
// "<topic>_<suffix>"; the rest follow "<prefix>_<topic>". Fragment/topic
// combinations are filtered against category rules so nothing like
// `bbq_speedrun` or `unnatural_worlds_postmortem` can appear.
func composeTitles(r *rand.Rand, n int) []string {
	seen := map[string]bool{}
	out := make([]string, 0, n)
	for len(out) < n {
		t := topics[r.IntN(len(topics))]
		var s string
		if r.IntN(4) == 0 {
			s = pickCompatible(r, prefixes, t.cat) + "_" + t.slug
		} else {
			s = t.slug + "_" + pickCompatible(r, suffixes, t.cat)
		}
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

// pickCompatible returns the slug of a uniformly random fragment that allows
// the given category. Every category has at least one compatible fragment in
// both `suffixes` and `prefixes`, so the rejection loop always terminates.
func pickCompatible(r *rand.Rand, frags []fragmentSpec, cat string) string {
	for {
		f := frags[r.IntN(len(frags))]
		if f.allowed == nil || f.allowed[cat] {
			return f.slug
		}
	}
}
