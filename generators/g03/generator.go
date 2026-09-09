// Package g03 contains an input generator for puzzle 03.
package g03

import (
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"
)

// Generate produces the input for puzzle 03, seeded deterministically by idx.
//
// The input is a single assembly program of 6 equation-blocks. Each block
// implements one linear equation over mem[0..5] using repeated add/sub loops
// (multiplication by a positive counter). Blocks alternate sweep direction so
// no cursor-reset instructions are needed.
//
// The reference solver treats the program as a black box: one baseline run
// plus 6 unit-vector probes reveal A and b for the system A*x = -b, which is
// then solved with big.Rat Gauss elimination.
func Generate(idx int) []byte {
	const n = 6
	r := rand.New(rand.NewPCG(0x04_5eed_5eed, uint64(idx)+1))

	// Reject-sampling loop: draw x*, A, offsets until A is invertible over
	// the rationals and every -Σ A[i][j]*x*[j] is nonzero.
	var xstar []int64
	var A [][]int64
	var offset [][]int64
	for {
		xstar = make([]int64, n)
		for j := 0; j < n; j++ {
			xstar[j] = int64(drawNonzero(r, 100))
		}
		A = make([][]int64, n)
		for i := 0; i < n; i++ {
			A[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				A[i][j] = int64(drawNonzero(r, 100))
			}
		}
		offset = make([][]int64, n)
		for i := 0; i < n; i++ {
			offset[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				offset[i][j] = int64(r.IntN(201) - 100)
			}
		}
		if !isInvertible(A) {
			continue
		}
		bad := false
		for i := 0; i < n && !bad; i++ {
			var s int64
			for j := 0; j < n; j++ {
				s += A[i][j] * xstar[j]
			}
			if s == 0 {
				bad = true
			}
		}
		if bad {
			continue
		}
		break
	}

	// Emitted immediates: b_i = -Σ A[i][j] * (x*_j + offset[i][j]).
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		var s int64
		for j := 0; j < n; j++ {
			s += A[i][j] * (xstar[j] + offset[i][j])
		}
		b[i] = -s
	}

	// Register pool. out_reg is fixed for the whole program; var_reg/ctr_reg
	// are chosen per emitted (i, j) block from the 5 non-out registers.
	outReg := byte('A' + r.IntN(6))
	nonOut := make([]byte, 0, 5)
	for c := byte('A'); c <= byte('F'); c++ {
		if c != outReg {
			nonOut = append(nonOut, c)
		}
	}

	var buf strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&buf, "SET %c %d\n", outReg, b[i])
		sweep := sweepOrder(i, n)
		for k, j := range sweep {
			vi := r.IntN(5)
			ci := r.IntN(4)
			if ci >= vi {
				ci++
			}
			varReg := nonOut[vi]
			ctrReg := nonOut[ci]

			a := A[i][j]
			mag := a
			if mag < 0 {
				mag = -mag
			}
			fmt.Fprintf(&buf, "LOAD %c\n", varReg)
			if offset[i][j] < 0 {
				fmt.Fprintf(&buf, "SUB %c %d\n", varReg, -offset[i][j])
			} else {
				fmt.Fprintf(&buf, "ADD %c %d\n", varReg, offset[i][j])
			}
			fmt.Fprintf(&buf, "SET %c %d\n", ctrReg, mag)
			if a > 0 {
				fmt.Fprintf(&buf, "ADD %c %c\n", outReg, varReg)
			} else {
				fmt.Fprintf(&buf, "SUB %c %c\n", outReg, varReg)
			}
			fmt.Fprintf(&buf, "SUB %c 1\n", ctrReg)
			fmt.Fprintf(&buf, "JNZ %c -2\n", ctrReg)
			if k < n-1 {
				if i%2 == 0 {
					buf.WriteString("RIGHT\n")
				} else {
					buf.WriteString("LEFT\n")
				}
			}
		}
		fmt.Fprintf(&buf, "OUT %c\n", outReg)
	}
	return []byte(buf.String())
}

// sweepOrder returns the column visit order for equation i. Even rows sweep
// left-to-right (0..n-1), odd rows sweep right-to-left (n-1..0), so each
// equation's final cursor position lines up with the next equation's start.
func sweepOrder(i, n int) []int {
	out := make([]int, n)
	for k := 0; k < n; k++ {
		if i%2 == 0 {
			out[k] = k
		} else {
			out[k] = n - 1 - k
		}
	}
	return out
}

// drawNonzero draws an integer from [-bound, bound] \ {0} biased toward the
// extremes: the magnitude is the max of two uniform draws on [1, bound], and
// the sign is fair.
func drawNonzero(r *rand.Rand, bound int) int {
	a := r.IntN(bound) + 1
	b := r.IntN(bound) + 1
	v := a
	if b > v {
		v = b
	}
	if r.IntN(2) == 1 {
		v = -v
	}
	return v
}

// isInvertible reports whether A is invertible over the rationals. Uses
// big.Rat Gauss elimination; returns false on the first zero-pivot column.
func isInvertible(A [][]int64) bool {
	n := len(A)
	M := make([][]*big.Rat, n)
	for i := 0; i < n; i++ {
		M[i] = make([]*big.Rat, n)
		for j := 0; j < n; j++ {
			M[i][j] = new(big.Rat).SetInt64(A[i][j])
		}
	}
	for k := 0; k < n; k++ {
		pivot := -1
		for i := k; i < n; i++ {
			if M[i][k].Sign() != 0 {
				pivot = i
				break
			}
		}
		if pivot < 0 {
			return false
		}
		if pivot != k {
			M[k], M[pivot] = M[pivot], M[k]
		}
		for i := k + 1; i < n; i++ {
			if M[i][k].Sign() == 0 {
				continue
			}
			factor := new(big.Rat).Quo(M[i][k], M[k][k])
			for j := k; j < n; j++ {
				term := new(big.Rat).Mul(factor, M[k][j])
				M[i][j].Sub(M[i][j], term)
			}
		}
	}
	return true
}
