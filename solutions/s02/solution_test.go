package s02

import (
	"fmt"
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
			name: "tiny open",
			// 3x3 open grid, S at (0,0), E at (2,2).
			// Part 1: shortest 4-way walk = 4.
			// Part 2: one blast jump with dr=2, dc=2 (Manhattan 4, cost 4,
			// consumes 1 bomb) — same 4 steps. Alternatively 2 walk + 1
			// diagonal blast (dr=1,dc=1 Manhattan 2 cost 2, plus 2 walk = 4);
			// or two blasts each dr=1,dc=1 (2 bombs, 2+2=4). All give 4.
			// The blast (dr=2,dc=2) from (0,0) → (2,2) has Manhattan 4 = cost 4
			// but is a single move and reaches E in 4 total.
			// So Part 2 = 4.
			input: "S..\n...\n..E\n",
			want1: "4",
			want2: "4",
		},
		{
			name: "corridor with double wall",
			// 5x5. S at (0,0). E at (4,4). Two walls at row 2 block the
			// corridor; blasting one bomb at (2,2) clears both walls at once
			// (i.e. jump from (1,2) to (3,2) with Manhattan 2, cost 2,
			// consumes 1 bomb).
			// Layout:
			//   S...
			//   ....
			//   .##.  (two walls at row 2 cols 1..2)
			//   ....
			//   ...E
			// Part 1: must go around the walls. Shortest walk = 8 (5 rows -1
			// then 5 cols -1 = 8 in an open 5x5). Since walls only affect
			// cols 1-2 at row 2, an open detour still exists.
			// Verified walk length = 8.
			// Part 2: jump from (1,1) with dr=2,dc=1 to (3,2) Manhattan 3,
			// or S (0,0) -> (3,3) diagonal Manhattan 6. Cheaper: walk to
			// (1,1), blast to (3,3) Manhattan 4 = 1 + 4 = 5, then walk 2 to
			// (4,4)? No E is (4,4).
			// Optimal blast from (0,0) -> (3,3) Manhattan 6 cost 6, then walk
			// 2 to E: total 8. No gain. Try blast (0,0) -> (2,2) Manhattan 4
			// cost 4, then walk (2,2)->(4,4)=4 total 8. Two blasts: (0,0)
			// -> (2,2) Manhattan 4 + (2,2) -> (4,4) Manhattan 4 = 8.
			// So Part 2 = 8 (same as Part 1) for this trivial example.
			input: "S....\n.....\n.##..\n.....\n....E\n",
			want1: "8",
			want2: "8",
		},
		{
			name: "zigzag corridor",
			// 7x7 zig-zag matching the Part 1 example. Row 0 is open, then
			// three wall bands at rows 1, 3, 5 each cover 6 of the 7 columns
			// but leave a single-cell notch at alternating ends so a walker
			// can snake between them:
			//   S......
			//   ######.  (notch at col 6)
			//   .......
			//   .######  (notch at col 0)
			//   .......
			//   ######.  (notch at col 6)
			//   ......E
			// Part 1: the shortest walk must traverse all three open rows in
			// full plus the six vertical joins = 6+1+1+6+1+1+6+1+1 = 24.
			// Part 2: Manhattan(S,E) = 12, and a monotone right/down bomb
			// chain achieves it, e.g. bomb (0,0)->(3,0) cost 3, walk
			// (3,0)->(4,0) 1, walk (4,0)->(4,6) 6, bomb (4,6)->(6,6) cost 2.
			// Uses 2 bombs, total 12.
			input: "S......\n######.\n.......\n.######\n.......\n######.\n......E\n",
			want1: "24",
			want2: "12",
		},
		{
			name: "detour with bomb shortcut",
			// 5x5, S at (0,0), E at (4,0). Two wall bands at rows 1 and 3
			// force the walker into a long C-shaped detour, but a single
			// well-placed bomb cuts straight down the left edge:
			//   S....
			//   ####.
			//   .....
			//   .####
			//   E....
			// Part 1: right 4, down 2, left 4, down 2 = 12 (Manhattan is
			// only 4, so the maze forces a big detour).
			// Part 2: one bomb from (0,0) to (3,0) costs 3 (dr=3, dc=0,
			// target is `.`), then walk (3,0)->(4,0) = 1. Total 4.
			input: "S....\n####.\n.....\n.####\nE....\n",
			want1: "12",
			want2: "4",
		},
		{
			name: "narrow bomb corridor",
			// 3x3 with a wall band at row 1 covering the middle two cells.
			// Walking is blocked; bombs are the only way through.
			//   S..
			//   ##.
			//   E..
			// Part 1: (0,0)->(0,2)=2, (0,2)->(1,2) dot, (1,2)->(2,2)=2,
			// (2,2)->(2,0)=2 -> total 6.
			// Part 2: single bomb (0,0)->(2,0), dr=2, dc=0, cost 2.
			input: "S..\n##.\nE..\n",
			want1: "6",
			want2: "2",
		},
		{
			name: "part2 example",
			// 7x7 with outer walls, matches the puzzle example in
			// puzzles/02/part2.md. Part 1: forced spiral of 16 steps.
			// Part 2: bomb (1,1)->(4,1) (dr=3, dc=0, cost 3) plus 5 walk
			// steps = 8, which equals Manhattan((1,1),(5,5)).
			input: "#######\n#S....#\n#####.#\n#.....#\n#.#####\n#....E#\n#######\n",
			want1: "16",
			want2: "8",
		},
		{
			name: "p1 unreachable, single dive",
			// 2x3 with S trapped by walls; walking cannot reach E.
			// Part 1: no walk exists -> -1.
			// Part 2: single dive (0,0)->(0,2) at Manhattan 2 lands on E.
			input: "S#E\n###\n",
			want1: "-1",
			want2: "2",
		},
		{
			name: "p1 unreachable, walk then dive",
			// 2x5; row 1 solid wall, wall at (0,3) blocks the walker.
			// Part 1: walker stops at (0,2) -> -1.
			// Part 2: walk 2 to (0,2), dive (0,2)->(0,4) at Manhattan 2 = 4.
			input: "S..#E\n#####\n",
			want1: "-1",
			want2: "4",
		},
		{
			name: "bigger example",
			// 15x15 recursive-backtracker maze matching the puzzle example
			// in puzzles/02/{part1,part2}.md. S at (1,0), E at (13,13),
			// outer border. Part 1: shortest walk 45. Part 2: two bombs at
			// (1,4) and (12,13) open row 1 and the approach to E for a
			// 25-step straight walk.
			input: `
###############
S...#.........#
#.###.#######.#
#.....#.....#.#
#.#####.###.#.#
#.#.....#...#.#
###.#.###.#.#.#
#...#...#.#...#
#.#.###.#.###.#
#.#.....#...#.#
#.###.###.#.#.#
#.#...#...#...#
#.#.###.#######
#............E#
###############
`,
			want1: "45",
			want2: "25",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Solution([]byte(strings.TrimSpace(tc.input)))
			if len(got) != 2 {
				t.Fatalf("want 2 answers, got %d", len(got))
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

func TestSolutionAdjacentSE(t *testing.T) {
	// 1x2 grid, S and E on the same row. Walk takes one step; Part 2's
	// dive graph doesn't beat that.
	got := Solution([]byte("SE\n"))
	if got[0] != "1" || got[1] != "1" {
		t.Errorf("SE adjacent: got %v, want [1 1]", got)
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
	Solution([]byte("S..\r\n...\r\n..E\r\n"))
}
