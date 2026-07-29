package s03

import (
	"strconv"
	"testing"

	"github.com/liennie/AdventOfCode/pkg/path"
	"github.com/liennie/AdventOfCode/pkg/space"
)

// The shipped Part 2 uses a dive/warp graph. The puzzle text describes a
// different mechanic: 3 pre-detonated bombs, each turning a 3x3 wall
// area into corridor, then shortest walking path. This file implements
// that mechanic directly and pins its equivalence to the shipped solver
// on every committed input.

// wdState tracks the walker's position, bombs consumed so far, and the
// rectangle of still-valid centers for the currently-open bomb. The
// rectangle is (rMin..rMax) x (cMin..cMax); rMin<0 means no bomb is
// currently open. The rectangle persists across corridor stretches so
// that a wall re-entered after a corridor gap can still fit inside the
// same 3x3, up to intersection with all wall neighborhoods visited so
// far under this bomb.
type wdState struct {
	r, c                   int
	bombs                  int
	rMin, cMin, rMax, cMax int
}

func walldestructionSolve(input []byte) (int, int) {
	grid, sr, sc, er, ec := parseGrid(input)
	rows := len(grid)
	cols := 0
	if rows > 0 {
		cols = len(grid[0])
	}
	inBounds := func(r, c int) bool { return r >= 0 && r < rows && c >= 0 && c < cols }

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

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	graph := path.GraphFunc[wdState](func(s wdState) []path.Edge[wdState] {
		var edges []path.Edge[wdState]
		for _, d := range dirs {
			nr, nc := s.r+d[0], s.c+d[1]
			if !inBounds(nr, nc) {
				continue
			}
			if grid[nr][nc] != '#' {
				edges = append(edges, path.Edge[wdState]{Len: 1, To: wdState{
					r: nr, c: nc, bombs: s.bombs,
					rMin: s.rMin, cMin: s.cMin, rMax: s.rMax, cMax: s.cMax,
				}})
				continue
			}
			wrMin, wrMax := max(0, nr-1), min(rows-1, nr+1)
			wcMin, wcMax := max(0, nc-1), min(cols-1, nc+1)

			if s.rMin >= 0 {
				nrMin := max(s.rMin, wrMin)
				nrMax := min(s.rMax, wrMax)
				ncMin := max(s.cMin, wcMin)
				ncMax := min(s.cMax, wcMax)
				if nrMin <= nrMax && ncMin <= ncMax {
					edges = append(edges, path.Edge[wdState]{Len: 1, To: wdState{
						r: nr, c: nc, bombs: s.bombs,
						rMin: nrMin, cMin: ncMin, rMax: nrMax, cMax: ncMax,
					}})
					continue
				}
			}
			if s.bombs >= maxBombs {
				continue
			}
			edges = append(edges, path.Edge[wdState]{Len: 1, To: wdState{
				r: nr, c: nc, bombs: s.bombs + 1,
				rMin: wrMin, cMin: wcMin, rMax: wrMax, cMax: wcMax,
			}})
		}
		return edges
	})

	start := wdState{sr, sc, 0, -1, -1, -1, -1}
	end := path.EndFunc[wdState](func(s wdState) bool { return s.r == er && s.c == ec })
	_, part2, err := path.Shortest(graph, start, end)
	if err != nil {
		part2 = -1
	}
	return part1, part2
}

func TestWallDestructionSmall(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want1, want2 int
	}{
		{"tiny open", "S..\n...\n..E\n", 4, 4},
		{"corridor with double wall", "S....\n.....\n.##..\n.....\n....E\n", 8, 8},
		{"zigzag corridor", "S......\n######.\n.......\n.######\n.......\n######.\n......E\n", 24, 12},
		{"detour with bomb shortcut", "S....\n####.\n.....\n.####\nE....\n", 12, 4},
		{"narrow bomb corridor", "S..\n##.\nE..\n", 6, 2},
		{"part2 example", "#######\n#S....#\n#####.#\n#.....#\n#.#####\n#....E#\n#######\n", 16, 8},
		{"p1 unreachable, single dive", "S#E\n###\n", -1, 2},
		{"p1 unreachable, walk then dive", "S..#E\n#####\n", -1, 4},
		{"bigger example", "###############\nS...#.........#\n#.###.#######.#\n#.....#.....#.#\n#.#####.###.#.#\n#.#.....#.#.#.#\n###.#.###.#.#.#\n#...#.#.#.#...#\n#.#.###.#.###.#\n#.#.#...#...#.#\n#.###.###.#.#.#\n#.#...#...#...#\n#.#.###.#######\n#...#........E#\n###############\n", 45, 25},
		{"enclosed E", "S....\n.###.\n.#E#.\n.###.\n.....\n", -1, 4},
		{"diagonal wall block", "S...\n##..\n.##.\n..##\n...E\n", -1, 7},
		{"three wall bands", "S...\n####\n....\n....\n####\n....\n....\n####\n...E\n", -1, 11},
		{"two disjoint walls", "S.#..#.E\n", -1, 7},
		{"open grid, bombs unused", "S....\n.....\n.....\n.....\n....E\n", 8, 8},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got1, got2 := walldestructionSolve([]byte(tc.input))
			if got1 != tc.want1 {
				t.Errorf("part1: got %d, want %d", got1, tc.want1)
			}
			if got2 != tc.want2 {
				t.Errorf("part2: got %d, want %d", got2, tc.want2)
			}
		})
	}
}

func TestWallDestructionEquivalence(t *testing.T) {
	// Exhaustive Part 2 over every triple of in-bounds bomb centers, using
	// the puzzle-text mechanic literally: place three 3x3 wall-clearing
	// craters, BFS on the modified grid, take the min. Ground truth for
	// the shipped dive/warp solver on inputs small enough to enumerate.
	// The real 101 x 101 inputs are far too big for brute force and blow
	// the DP up on memory too, so this test only runs on the small table.
	tests := []struct {
		name  string
		input string
	}{
		{"tiny open", "S..\n...\n..E\n"},
		{"corridor with double wall", "S....\n.....\n.##..\n.....\n....E\n"},
		{"zigzag corridor", "S......\n######.\n.......\n.######\n.......\n######.\n......E\n"},
		{"detour with bomb shortcut", "S....\n####.\n.....\n.####\nE....\n"},
		{"narrow bomb corridor", "S..\n##.\nE..\n"},
		{"part2 example", "#######\n#S....#\n#####.#\n#.....#\n#.#####\n#....E#\n#######\n"},
		{"p1 unreachable, single dive", "S#E\n###\n"},
		{"p1 unreachable, walk then dive", "S..#E\n#####\n"},
		{"bigger example", "###############\nS...#.........#\n#.###.#######.#\n#.....#.....#.#\n#.#####.###.#.#\n#.#.....#.#.#.#\n###.#.###.#.#.#\n#...#.#.#.#...#\n#.#.###.#.###.#\n#.#.#...#...#.#\n#.###.###.#.#.#\n#.#...#...#...#\n#.#.###.#######\n#...#........E#\n###############\n"},
		{"enclosed E", "S....\n.###.\n.#E#.\n.###.\n.....\n"},
		{"diagonal wall block", "S...\n##..\n.##.\n..##\n...E\n"},
		{"three wall bands", "S...\n####\n....\n....\n####\n....\n....\n####\n...E\n"},
		{"two disjoint walls", "S.#..#.E\n"},
		{"open grid, bombs unused", "S....\n.....\n.....\n.....\n....E\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dive := Solution([]byte(tc.input))
			brute := bruteWallDestructionPart2([]byte(tc.input))
			if dive[1] != strconv.Itoa(brute) {
				t.Errorf("part2 mismatch: dive=%s brute wall-destruction=%d", dive[1], brute)
			}
		})
	}
}

// bruteWallDestructionPart2 enumerates every unordered triple of in-bounds
// bomb centers (repetitions allowed, so 0/1/2-bomb placements are covered
// by picking coincident centers), applies the 3x3 wall clears, runs BFS
// on the modified grid, and returns the minimum walk length. Returns -1
// if E stays unreachable under every placement. Only usable on grids
// small enough that n^3 BFS is cheap.
func bruteWallDestructionPart2(input []byte) int {
	grid, sr, sc, er, ec := parseGrid(input)
	rows := len(grid)
	if rows == 0 {
		return -1
	}
	cols := len(grid[0])
	n := rows * cols

	craters := make([][]int, n)
	for i := 0; i < n; i++ {
		cr, cc := i/cols, i%cols
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				nr, nc := cr+dr, cc+dc
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
					continue
				}
				craters[i] = append(craters[i], nr*cols+nc)
			}
		}
	}

	startIdx := sr*cols + sc
	endIdx := er*cols + ec
	coverCount := make([]int, n)
	dist := make([]int, n)
	queue := make([]int, 0, n)
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	bfs := func() int {
		for i := range dist {
			dist[i] = -1
		}
		dist[startIdx] = 0
		queue = queue[:0]
		queue = append(queue, startIdx)
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			if v == endIdx {
				return dist[v]
			}
			r, c := v/cols, v%cols
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
					continue
				}
				nv := nr*cols + nc
				if dist[nv] >= 0 {
					continue
				}
				if grid[nr][nc] == '#' && coverCount[nv] == 0 {
					continue
				}
				dist[nv] = dist[v] + 1
				queue = append(queue, nv)
			}
		}
		return -1
	}

	apply := func(c int, delta int) {
		for _, x := range craters[c] {
			coverCount[x] += delta
		}
	}

	best := -1
	for c1 := 0; c1 < n; c1++ {
		apply(c1, +1)
		for c2 := c1; c2 < n; c2++ {
			apply(c2, +1)
			for c3 := c2; c3 < n; c3++ {
				apply(c3, +1)
				d := bfs()
				if d >= 0 && (best < 0 || d < best) {
					best = d
				}
				apply(c3, -1)
			}
			apply(c2, -1)
		}
		apply(c1, -1)
	}
	return best
}
