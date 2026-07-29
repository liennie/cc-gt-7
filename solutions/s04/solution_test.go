package s04

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

// TestMultiplyLoop covers the SET-body-SUB-JNZ pattern the generator emits
// for every column: A = 0 + B * C = 0 + 7 * 5 = 35.
func TestMultiplyLoop(t *testing.T) {
	prog := parseProgram(`SET A 0
SET B 7
SET C 5
ADD A B
SUB C 1
JNZ C -2
OUT A`)
	got := runVM(prog, nil)
	if len(got) != 1 || got[0] != 35 {
		t.Errorf("got %v, want [35]", got)
	}
}

// SAVE and LOAD at negative cursor positions round-trip.
func TestSaveLoadNegativeCursor(t *testing.T) {
	prog := parseProgram(`SAVE 7
LEFT
SAVE -3
RIGHT
LOAD A
LEFT
LOAD B
OUT A
OUT B`)
	got := runVM(prog, nil)
	if len(got) != 2 || got[0] != 7 || got[1] != -3 {
		t.Errorf("got %v, want [7 -3]", got)
	}
}

// JMP with positive offset skips the SET, jumps from pc=1 to pc=3.
func TestJMPForward(t *testing.T) {
	prog := parseProgram(`SET A 1
JMP 2
SET A 999
OUT A`)
	got := runVM(prog, nil)
	if len(got) != 1 || got[0] != 1 {
		t.Errorf("got %v, want [1]", got)
	}
}

// JNZ falls through when the register is zero.
func TestJNZNotTaken(t *testing.T) {
	prog := parseProgram(`SET A 0
JNZ A 3
SET A 42
OUT A`)
	got := runVM(prog, nil)
	if len(got) != 1 || got[0] != 42 {
		t.Errorf("got %v, want [42]", got)
	}
}

// Initial memory patches are visible to LOAD before the program runs.
func TestPatchesMerged(t *testing.T) {
	prog := parseProgram(`LOAD A
RIGHT
LOAD B
OUT A
OUT B`)
	got := runVM(prog, map[int64]int64{0: 11, 1: 22})
	if len(got) != 2 || got[0] != 11 || got[1] != 22 {
		t.Errorf("got %v, want [11 22]", got)
	}
}

// End-to-end 2-var / 2-eqn linear system, matching the worked example in
// puzzles/04/part1.md + part2.md.
//
// V0 = 1 + 2*m_0 + m_1;  V1 = -1 + m_0 + m_1
// Part 1 (mem = 0):  outputs [1, -1],  digest = 1 + (-1)*10 = -9
// Part 2 (mem = [-2, 3]):  outputs [0, 0],  digest = -2 + 3*10 = 28
func TestSolutionTwoVar(t *testing.T) {
	input := `SET A 1
LOAD B
SET C 2
ADD A B
SUB C 1
JNZ C -2
RIGHT
LOAD B
ADD A B
OUT A
SET A -1
LOAD B
ADD A B
LEFT
LOAD B
ADD A B
OUT A
`
	got := Solution([]byte(input))
	if len(got) != 2 || got[0] != "-9" || got[1] != "28" {
		t.Errorf("got %v, want [-9 28]", got)
	}
}

// emitProgramS04 emits an assembly program shaped exactly like the ones
// g04.Generate produces: n equation blocks each of the form
//
//	SET A b_i
//	<for j in sweep order>:
//	  LOAD B; ADD/SUB B offset; SET C |A[i][j]|; ADD/SUB A B; SUB C 1; JNZ C -2
//	  RIGHT|LEFT (skipped after last column)
//	OUT A
//
// out_reg = A, var_reg = B, ctr_reg = C. Sweep direction alternates so
// no explicit cursor reset is needed between blocks (even i sweeps LTR,
// odd i sweeps RTL).
func emitProgramS04(A, offset [][]int64, b []int64) string {
	n := len(A)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "SET A %d\n", b[i])
		ltr := i%2 == 0
		for step := 0; step < n; step++ {
			var j int
			if ltr {
				j = step
			} else {
				j = n - 1 - step
			}
			sb.WriteString("LOAD B\n")
			if offset[i][j] > 0 {
				fmt.Fprintf(&sb, "ADD B %d\n", offset[i][j])
			} else if offset[i][j] < 0 {
				fmt.Fprintf(&sb, "SUB B %d\n", -offset[i][j])
			}
			a := A[i][j]
			absA := a
			if absA < 0 {
				absA = -absA
			}
			fmt.Fprintf(&sb, "SET C %d\n", absA)
			if a >= 0 {
				sb.WriteString("ADD A B\n")
			} else {
				sb.WriteString("SUB A B\n")
			}
			sb.WriteString("SUB C 1\n")
			sb.WriteString("JNZ C -2\n")
			if step < n-1 {
				if ltr {
					sb.WriteString("RIGHT\n")
				} else {
					sb.WriteString("LEFT\n")
				}
			}
		}
		sb.WriteString("OUT A\n")
	}
	return sb.String()
}

// bruteforceSolveS04 enumerates every mem vector in [-M..M]^n and returns
// the (unique) one that makes every OUT zero. Panics if none or multiple.
func bruteforceSolveS04(prog []instr, n, M int) []int64 {
	var found []int64
	cand := make([]int64, n)
	var walk func(i int)
	walk = func(i int) {
		if i == n {
			patches := make(map[int64]int64, n)
			for j, v := range cand {
				patches[int64(j)] = v
			}
			out := runVM(prog, patches)
			if len(out) != n {
				return
			}
			for _, v := range out {
				if v != 0 {
					return
				}
			}
			if found != nil {
				panic("bruteforceSolveS04: multiple solutions")
			}
			found = append([]int64(nil), cand...)
			return
		}
		for v := -M; v <= M; v++ {
			cand[i] = int64(v)
			walk(i + 1)
		}
	}
	walk(0)
	if found == nil {
		panic("bruteforceSolveS04: no solution")
	}
	return found
}

// TestSolutionAgainstBruteforce generates small random 2- and 3-variable
// programs that share the exact emission shape of g04.Generate and
// verifies Solution's recovered x* by brute-force enumeration over a
// bounded mem-value cube.
func TestSolutionAgainstBruteforce(t *testing.T) {
	r := rand.New(rand.NewPCG(0x504, 0x504))
	for trial := 0; trial < 20; trial++ {
		n := 2 + trial%2 // 2 or 3
		const bound = 3
		var A, offset [][]int64
		var xstar []int64
		for {
			xstar = make([]int64, n)
			for j := 0; j < n; j++ {
				v := 1 + r.IntN(bound)
				if r.IntN(2) == 0 {
					v = -v
				}
				xstar[j] = int64(v)
			}
			A = make([][]int64, n)
			for i := 0; i < n; i++ {
				A[i] = make([]int64, n)
				for j := 0; j < n; j++ {
					v := 1 + r.IntN(bound)
					if r.IntN(2) == 0 {
						v = -v
					}
					A[i][j] = int64(v)
				}
			}
			offset = make([][]int64, n)
			for i := 0; i < n; i++ {
				offset[i] = make([]int64, n)
				for j := 0; j < n; j++ {
					offset[i][j] = int64(r.IntN(2*bound+1) - bound)
				}
			}
			// require invertibility and nonzero baseline outputs.
			if n == 2 && det2(A) == 0 {
				continue
			}
			if n == 3 && det3(A) == 0 {
				continue
			}
			bad := false
			for i := 0; i < n; i++ {
				var s int64
				for j := 0; j < n; j++ {
					s += A[i][j] * xstar[j]
				}
				if s == 0 {
					bad = true
					break
				}
			}
			if !bad {
				break
			}
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			var s int64
			for j := 0; j < n; j++ {
				s += A[i][j] * (xstar[j] + offset[i][j])
			}
			b[i] = -s
		}
		src := emitProgramS04(A, offset, b)
		got := Solution([]byte(src))

		// Expected digest of Part 2 is Σ xstar_j * 10^j.
		var wantP2 int64
		w := int64(1)
		for _, v := range xstar {
			wantP2 += v * w
			w *= 10
		}
		if got[1] != strconv.FormatInt(wantP2, 10) {
			t.Errorf("trial %d n=%d: part2 got %s want %d (x*=%v)",
				trial, n, got[1], wantP2, xstar)
			continue
		}

		// Independently verify via brute enumeration.
		prog := parseProgram(src)
		brute := bruteforceSolveS04(prog, n, bound)
		for j, v := range brute {
			if v != xstar[j] {
				t.Errorf("trial %d n=%d: brute x*[%d]=%d, wanted %d",
					trial, n, j, v, xstar[j])
			}
		}
	}
}

func det2(A [][]int64) int64 {
	if len(A) != 2 {
		return 0
	}
	return A[0][0]*A[1][1] - A[0][1]*A[1][0]
}

func det3(A [][]int64) int64 {
	if len(A) != 3 {
		return 0
	}
	return A[0][0]*(A[1][1]*A[2][2]-A[1][2]*A[2][1]) -
		A[0][1]*(A[1][0]*A[2][2]-A[1][2]*A[2][0]) +
		A[0][2]*(A[1][0]*A[2][1]-A[1][1]*A[2][0])
}
