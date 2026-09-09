// Package g02 contains an input generator for puzzle 02.
package g02

import (
	"math/rand/v2"
	"strings"
)

// Generate produces a maze for puzzle 02, seeded deterministically by idx.
//
// The maze is 101x101 (a 50-cell recursive-backtracker maze wrapped in an
// outer wall). Maze construction starts from the centre cell (50, 50) so the
// topology differs significantly from a corner-seeded maze. S is fixed at
// (1, 1), E at (99, 99). A handful of extra walls are punched out to create
// loops so the Part 2 blast model has real shortcuts.
func Generate(idx int) []byte {
	r := rand.New(rand.NewPCG(0x03_5eed_5eed, uint64(idx)+1))

	const n = 50 // cells per side => 101x101 grid
	w, h := 2*n+1, 2*n+1
	grid := make([][]byte, h)
	for i := range grid {
		grid[i] = make([]byte, w)
		for j := range grid[i] {
			grid[i][j] = '#'
		}
	}

	// Iterative recursive-backtracker maze growth starting from the centre.
	type cell struct{ r, c int }
	startR, startC := 2*(n/2)+1, 2*(n/2)+1 // odd coordinate near centre
	stack := []cell{{startR, startC}}
	grid[startR][startC] = '.'
	dirs := []cell{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		r.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
		advanced := false
		for _, d := range dirs {
			nr, nc := top.r+d.r, top.c+d.c
			if nr < 1 || nr >= h-1 || nc < 1 || nc >= w-1 {
				continue
			}
			if grid[nr][nc] == '.' {
				continue
			}
			grid[top.r+d.r/2][top.c+d.c/2] = '.'
			grid[nr][nc] = '.'
			stack = append(stack, cell{nr, nc})
			advanced = true
			break
		}
		if !advanced {
			stack = stack[:len(stack)-1]
		}
	}

	// Knock out ~5% of interior walls to create loops.
	interior := 0
	for i := 2; i < h-1; i++ {
		for j := 2; j < w-1; j++ {
			if grid[i][j] == '#' && ((i%2 == 0) != (j%2 == 0)) {
				interior++
			}
		}
	}
	knock := interior / 20
	for k := 0; k < knock; k++ {
		i := 2 + r.IntN(h-3)
		j := 2 + r.IntN(w-3)
		if grid[i][j] == '#' && ((i%2 == 0) != (j%2 == 0)) {
			grid[i][j] = '.'
		}
	}

	grid[1][0] = 'S'
	grid[h-2][w-2] = 'E'

	var b strings.Builder
	for _, row := range grid {
		b.Write(row)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}
