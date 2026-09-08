// Package s03 contains a solution to puzzle 03.
package s03

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/liennie/AdventOfCode/pkg/path"
	"github.com/liennie/AdventOfCode/pkg/space"
)

const maxBombs = 3

type state struct {
	r, c, b int
}

func parseGrid(input []byte) (grid [][]byte, sr, sc, er, ec int) {
	sr, sc, er, ec = -1, -1, -1, -1
	for _, line := range strings.Split(string(input), "\n") {
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

func Solution(input []byte) []string {
	if bytes.IndexByte(input, '\r') >= 0 {
		panic("s03: CR in input, expected LF-only")
	}
	grid, sr, sc, er, ec := parseGrid(input)
	rows := len(grid)
	cols := 0
	if rows > 0 {
		cols = len(grid[0])
	}
	inBounds := func(r, c int) bool { return r >= 0 && r < rows && c >= 0 && c < cols }

	// Part 1: shortest 4-directional walk from S to E, walls block.
	_, part1, err := path.Shortest(path.GraphFunc[space.Point](func(p space.Point) (e []path.Edge[space.Point]) {
		for dir := range space.Orthogonal() {
			np := p.Add(dir)
			if inBounds(np.Y, np.X) && grid[np.Y][np.X] != '#' {
				e = append(e, path.Edge[space.Point]{Len: 1, To: np})
			}
		}
		return
	}), space.Point{X: sc, Y: sr}, path.EndConst(space.Point{X: ec, Y: er}))
	if err != nil {
		part1 = -1
	}

	// Part 2: shortest path with up to 3 "pre-placed" bomb jumps.
	// A jump from (r, c) can reach any (r', c') with |dr| <= 3, |dc| <= 3,
	// (dr, dc) != (0, 0), and (dr, dc) not one of the four extreme corners
	// (+-3, +-3). The destination must be in-bounds. Cost = |dr|+|dc| walk
	// steps plus one bomb charge consumed.
	graph := path.GraphFunc[state](func(s state) []path.Edge[state] {
		var edges []path.Edge[state]
		// Walk moves.
		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nr, nc := s.r+d[0], s.c+d[1]
			if !inBounds(nr, nc) || grid[nr][nc] == '#' {
				continue
			}
			edges = append(edges, path.Edge[state]{Len: 1, To: state{nr, nc, s.b}})
		}
		// Blast/jump moves.
		if s.b < maxBombs {
			for dr := -3; dr <= 3; dr++ {
				for dc := -3; dc <= 3; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					// Exclude the four extreme corners (+-3, +-3).
					if (dr == 3 || dr == -3) && (dc == 3 || dc == -3) {
						continue
					}
					nr, nc := s.r+dr, s.c+dc
					if !inBounds(nr, nc) {
						continue
					}
					manhattan := abs(dr) + abs(dc)
					edges = append(edges, path.Edge[state]{Len: manhattan, To: state{nr, nc, s.b + 1}})
				}
			}
		}
		return edges
	})

	start := state{sr, sc, 0}
	end := path.EndFunc[state](func(s state) bool { return s.r == er && s.c == ec })
	_, part2, err := path.Shortest(graph, start, end)
	if err != nil {
		part2 = -1
	}

	return []string{
		strconv.Itoa(part1),
		strconv.Itoa(part2),
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
