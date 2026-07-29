// Package s05 contains a solution to puzzle 05.
package s05

import (
	"strconv"
	"strings"
)

// Input has three blank-line-separated sections:
//   1. WEIGHTS: lines of "<name> <int>"
//   2. DISLIKES: lines of "<name> <name>"
//   3. QUERIES: lines of "<name> <name>"
//
// Part 1: sum of hype[a] + hype[b] over all query pairs (a, b) whose
// talks lie in different connected components of the conflict graph.
// Part 2: maximum total weight of an independent set in the conflict graph
// (a subset of nodes containing no conflict edge).

func Solution(input []byte) []string {
	sections := splitBlank(string(input))
	if len(sections) < 3 {
		return []string{"0", "0"}
	}

	// Parse weights.
	weights := map[string]int64{}
	nameOrder := []string{}
	for _, line := range strings.Split(sections[0], "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		w, err := strconv.ParseInt(f[1], 10, 64)
		if err != nil {
			continue
		}
		if _, seen := weights[f[0]]; !seen {
			nameOrder = append(nameOrder, f[0])
		}
		weights[f[0]] = w
	}

	// Parse dislikes into adjacency and build union-find over names.
	adj := map[string]map[string]bool{}
	uf := map[string]string{}
	var find func(string) string
	find = func(x string) string {
		if uf[x] == "" {
			uf[x] = x
			return x
		}
		if uf[x] == x {
			return x
		}
		r := find(uf[x])
		uf[x] = r
		return r
	}
	union := func(a, b string) {
		ra, rb := find(a), find(b)
		if ra != rb {
			uf[ra] = rb
		}
	}
	addEdge := func(a, b string) {
		if adj[a] == nil {
			adj[a] = map[string]bool{}
		}
		if adj[b] == nil {
			adj[b] = map[string]bool{}
		}
		adj[a][b] = true
		adj[b][a] = true
		union(a, b)
	}
	for _, line := range strings.Split(sections[1], "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		addEdge(f[0], f[1])
	}
	// Make sure every named guest has a UF entry.
	for _, n := range nameOrder {
		find(n)
	}

	// Part 1: sum hype of both endpoints for queries whose talks are in
	// different components.
	part1 := int64(0)
	for _, line := range strings.Split(sections[2], "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		if find(f[0]) != find(f[1]) {
			part1 = addCheckedS05(part1, addCheckedS05(weights[f[0]], weights[f[1]]))
		}
	}

	// Part 2: weighted max independent set on the forest via tree DP.
	// For each tree component, run DP from an arbitrary root:
	//   take[v] = weight[v] + sum over children c of skip[c]
	//   skip[v] = sum over children c of max(take[c], skip[c])
	visited := map[string]bool{}
	var part2 int64
	for _, root := range nameOrder {
		if visited[root] {
			continue
		}
		// Iterative post-order over the tree rooted at `root`. Since the graph
		// is a forest, no cycles; parent map suffices.
		parent := map[string]string{root: ""}
		order := []string{}
		stack := []string{root}
		visited[root] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for u := range adj[v] {
				if visited[u] {
					continue
				}
				visited[u] = true
				parent[u] = v
				stack = append(stack, u)
			}
		}
		// Process in reverse (post-order).
		take := map[string]int64{}
		skip := map[string]int64{}
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			take[v] = weights[v]
			skip[v] = 0
			for u := range adj[v] {
				if u == parent[v] {
					continue
				}
				take[v] = addCheckedS05(take[v], skip[u])
				if take[u] > skip[u] {
					skip[v] = addCheckedS05(skip[v], take[u])
				} else {
					skip[v] = addCheckedS05(skip[v], skip[u])
				}
			}
		}
		best := take[root]
		if skip[root] > best {
			best = skip[root]
		}
		part2 = addCheckedS05(part2, best)
	}

	return []string{
		strconv.FormatInt(part1, 10),
		strconv.FormatInt(part2, 10),
	}
}

// addCheckedS05 adds two non-negative int64 hype sums and panics on
// int64 overflow. Answers here are tiny (80 nodes * 1000 weight upper
// bound = 80000), but the guard makes any future weight-range change
// that would blow past int64 fail loudly instead of silently wrapping.
func addCheckedS05(a, b int64) int64 {
	s := a + b
	if (a > 0 && b > 0 && s < 0) || (a < 0 && b < 0 && s >= 0) {
		panic("s05: int64 overflow accumulating hype sum")
	}
	return s
}

// splitBlank splits s on runs of blank lines, returning the non-empty chunks
// (each chunk preserves internal newlines).
func splitBlank(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	var out []string
	var cur []string
	flush := func() {
		if len(cur) == 0 {
			return
		}
		out = append(out, strings.Join(cur, "\n"))
		cur = nil
	}
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			flush()
		} else {
			cur = append(cur, ln)
		}
	}
	flush()
	return out
}
