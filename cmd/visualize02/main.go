// Program visualize02 renders the s02 solver's path onto the input grid.
//
// For every committed puzzle 02 input it writes two ASCII files: one for
// Part 1 (walk arrows) and one for Part 2 (walk arrows plus reconstructed
// bomb craters marked with @). The bomb positions are recovered from the
// dive-graph edges: for every dive A -> B the program searches for a 3x3
// crater center that lets a Manhattan-length walk connect A and B, then
// draws that walk.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/liennie/AdventOfCode/pkg/path"
	"github.com/liennie/AdventOfCode/pkg/space"
)

const maxBombs = 3

type state struct {
	r, c, b int
}

func main() {
	outDir := filepath.Join("solutions", "s02", "visualizations")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}
	for i := 1; i <= 5; i++ {
		in := filepath.Join("puzzles", "02", "inputs", fmt.Sprintf("%02d.txt", i))
		input, err := os.ReadFile(in)
		if err != nil {
			panic(err)
		}
		p1 := visualizePart1(input)
		if err := os.WriteFile(filepath.Join(outDir, fmt.Sprintf("%02d.part1.txt", i)), p1, 0o644); err != nil {
			panic(err)
		}
		p2 := visualizePart2(input)
		if err := os.WriteFile(filepath.Join(outDir, fmt.Sprintf("%02d.part2.txt", i)), p2, 0o644); err != nil {
			panic(err)
		}
		fmt.Printf("input %02d: wrote part1 (%d bytes) and part2 (%d bytes)\n", i, len(p1), len(p2))
	}
}

func parseGrid(input []byte) (grid [][]byte, sr, sc, er, ec int) {
	sr, sc, er, ec = -1, -1, -1, -1
	text := string(bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n")))
	for _, line := range splitLines(text) {
		if line == "" {
			continue
		}
		row := []byte(line)
		for c, ch := range row {
			switch ch {
			case 'S':
				sr, sc = len(grid), c
			case 'E':
				er, ec = len(grid), c
			}
		}
		grid = append(grid, row)
	}
	return
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func visualizePart1(input []byte) []byte {
	grid, sr, sc, er, ec := parseGrid(input)
	rows := len(grid)
	cols := len(grid[0])
	inBounds := func(r, c int) bool { return r >= 0 && r < rows && c >= 0 && c < cols }

	nodes, _, err := path.Shortest(path.GraphFunc[space.Point](func(p space.Point) (e []path.Edge[space.Point]) {
		for dir := range space.Orthogonal() {
			np := p.Add(dir)
			if inBounds(np.Y, np.X) && grid[np.Y][np.X] != '#' {
				e = append(e, path.Edge[space.Point]{Len: 1, To: np})
			}
		}
		return
	}), space.Point{X: sc, Y: sr}, path.EndConst(space.Point{X: ec, Y: er}))
	if err != nil {
		return []byte("no path\n")
	}

	out := copyGrid(grid)
	for i := 0; i+1 < len(nodes); i++ {
		a, b := nodes[i], nodes[i+1]
		put(out, grid, a.Y, a.X, arrowChar(b.Y-a.Y, b.X-a.X))
	}
	return render(out)
}

func visualizePart2(input []byte) []byte {
	grid, sr, sc, er, ec := parseGrid(input)
	rows := len(grid)
	cols := len(grid[0])
	inBounds := func(r, c int) bool { return r >= 0 && r < rows && c >= 0 && c < cols }

	graph := path.GraphFunc[state](func(s state) []path.Edge[state] {
		var edges []path.Edge[state]
		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nr, nc := s.r+d[0], s.c+d[1]
			if !inBounds(nr, nc) || grid[nr][nc] == '#' {
				continue
			}
			edges = append(edges, path.Edge[state]{Len: 1, To: state{nr, nc, s.b}})
		}
		if s.b < maxBombs {
			for dr := -3; dr <= 3; dr++ {
				for dc := -3; dc <= 3; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					if (dr == 3 || dr == -3) && (dc == 3 || dc == -3) {
						continue
					}
					nr, nc := s.r+dr, s.c+dc
					if !inBounds(nr, nc) {
						continue
					}
					edges = append(edges, path.Edge[state]{Len: abs(dr) + abs(dc), To: state{nr, nc, s.b + 1}})
				}
			}
		}
		return edges
	})

	start := state{sr, sc, 0}
	end := path.EndFunc[state](func(s state) bool { return s.r == er && s.c == ec })
	nodes, _, err := path.Shortest(graph, start, end)
	if err != nil {
		return []byte("no path\n")
	}

	type arrow struct {
		r, c int
		ch   byte
	}
	var bombs [][2]int
	var arrows []arrow

	for i := 0; i+1 < len(nodes); i++ {
		a, b := nodes[i], nodes[i+1]
		if b.b == a.b {
			arrows = append(arrows, arrow{a.r, a.c, arrowChar(b.r-a.r, b.c-a.c)})
			continue
		}
		br, bc, walk, ok := findBombWalk(grid, a.r, a.c, b.r, b.c, rows, cols)
		if !ok {
			panic(fmt.Sprintf("visualize02: no bomb position for dive (%d,%d)->(%d,%d)", a.r, a.c, b.r, b.c))
		}
		bombs = append(bombs, [2]int{br, bc})
		for j := 0; j+1 < len(walk); j++ {
			p1, p2 := walk[j], walk[j+1]
			arrows = append(arrows, arrow{p1[0], p1[1], arrowChar(p2[0]-p1[0], p2[1]-p1[1])})
		}
	}

	out := copyGrid(grid)
	for _, b := range bombs {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				r, c := b[0]+dr, b[1]+dc
				if r >= 0 && r < rows && c >= 0 && c < cols {
					put(out, grid, r, c, '@')
				}
			}
		}
	}
	for _, a := range arrows {
		put(out, grid, a.r, a.c, a.ch)
	}
	return render(out)
}

// findBombWalk searches for a 3x3 crater center that lets a 4-directional
// walk of Manhattan length reach (er,ec) from (sr,sc). It returns the
// center and the reconstructed walk (start .. end inclusive).
func findBombWalk(grid [][]byte, sr, sc, er, ec, rows, cols int) (int, int, [][2]int, bool) {
	target := abs(er-sr) + abs(ec-sc)
	rLo, rHi := min(sr, er)-2, max(sr, er)+2
	cLo, cHi := min(sc, ec)-2, max(sc, ec)+2
	for br := rLo; br <= rHi; br++ {
		for bc := cLo; bc <= cHi; bc++ {
			if walk, ok := bfsWalk(grid, sr, sc, er, ec, br, bc, target, rows, cols); ok {
				return br, bc, walk, true
			}
		}
	}
	return 0, 0, nil, false
}

func bfsWalk(grid [][]byte, sr, sc, er, ec, br, bc, target, rows, cols int) ([][2]int, bool) {
	walkable := func(r, c int) bool {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false
		}
		if abs(r-br) <= 1 && abs(c-bc) <= 1 {
			return true
		}
		return grid[r][c] != '#'
	}
	if !walkable(sr, sc) || !walkable(er, ec) {
		return nil, false
	}
	type key = [2]int
	dist := map[key]int{{sr, sc}: 0}
	prev := map[key]key{}
	queue := []key{{sr, sc}}
	for i := 0; i < len(queue); i++ {
		cur := queue[i]
		if cur == (key{er, ec}) {
			break
		}
		if dist[cur] >= target {
			continue
		}
		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nr, nc := cur[0]+d[0], cur[1]+d[1]
			if !walkable(nr, nc) {
				continue
			}
			np := key{nr, nc}
			if _, seen := dist[np]; seen {
				continue
			}
			dist[np] = dist[cur] + 1
			prev[np] = cur
			queue = append(queue, np)
		}
	}
	d, ok := dist[key{er, ec}]
	if !ok || d != target {
		return nil, false
	}
	var walk [][2]int
	cur := key{er, ec}
	for cur != (key{sr, sc}) {
		walk = append(walk, cur)
		cur = prev[cur]
	}
	walk = append(walk, key{sr, sc})
	for i, j := 0, len(walk)-1; i < j; i, j = i+1, j-1 {
		walk[i], walk[j] = walk[j], walk[i]
	}
	return walk, true
}

func copyGrid(grid [][]byte) [][]byte {
	out := make([][]byte, len(grid))
	for r := range grid {
		out[r] = make([]byte, len(grid[r]))
		copy(out[r], grid[r])
	}
	return out
}

// put writes ch at (r,c) in out unless the original cell is S or E.
func put(out, orig [][]byte, r, c int, ch byte) {
	if o := orig[r][c]; o == 'S' || o == 'E' {
		return
	}
	out[r][c] = ch
}

func render(grid [][]byte) []byte {
	var buf bytes.Buffer
	for r := range grid {
		buf.Write(grid[r])
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func arrowChar(dr, dc int) byte {
	switch {
	case dr < 0 && dc == 0:
		return '^'
	case dr > 0 && dc == 0:
		return 'v'
	case dr == 0 && dc < 0:
		return '<'
	case dr == 0 && dc > 0:
		return '>'
	}
	return '?'
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
