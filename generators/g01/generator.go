// Package g01 contains an input generator for puzzle 01.
package g01

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Generate produces the input for puzzle 01, seeded deterministically by idx.
//
// The input is a list of 40 fixed-point numbers with 2 decimals, one per line.
// Values are drawn from multiples of 0.25 in [0.25, 6.75] so that:
//   - Part 1 (pairs summing to 7.00) has a healthy count.
//   - Part 2 (signed-subset count summing to 7.00) fits comfortably in int64
//     (3^40 upper bound, sparsified by the 25-cent grid down to ~5e15).
func Generate(idx int) []byte {
	r := rand.New(rand.NewPCG(0x01_5eed_5eed, uint64(idx)+1))

	const (
		n        = 40
		stepCent = 25
		minCent  = 25
		maxCent  = 675 // exclusive of 700 so no zero-sum after negation
	)
	choices := (maxCent-minCent)/stepCent + 1

	var b strings.Builder
	for i := 0; i < n; i++ {
		v := minCent + r.IntN(choices)*stepCent
		fmt.Fprintf(&b, "%d.%02d\n", v/100, v%100)
	}
	return []byte(b.String())
}
