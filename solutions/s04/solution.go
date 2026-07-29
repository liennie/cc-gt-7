// Package s04 contains a solution to puzzle 04.
package s04

import (
	"math/big"
	"strconv"
	"strings"
)

// The GT-Nano assembly VM: 10 opcodes, six signed-integer registers A..F, a
// 1-D signed-integer memory tape, a signed-integer cursor. Jumps take signed
// non-zero offsets applied to the pc of the jump instruction itself.
//
// The program encodes a linear function of mem[0..n-1] for some small n. Both
// parts treat the program as a black box. Part 1 runs it with mem = 0 and
// reports a positional base-10 weighted sum of the outputs. Part 2 recovers
// the underlying matrix A and constant vector b by probing with mem = 0 and
// each unit vector, solves A*x = -b with big.Rat Gauss elimination, and
// reports the same base-10 weighted sum of x.

type operand struct {
	isReg bool
	reg   byte
	imm   int64
}

type instr struct {
	op   string
	args []operand
}

func parseOperand(s string) operand {
	if len(s) == 1 && s[0] >= 'A' && s[0] <= 'F' {
		return operand{isReg: true, reg: s[0]}
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("s04: bad operand: " + s)
	}
	return operand{imm: n}
}

func parseProgram(src string) []instr {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	var prog []instr
	for _, raw := range strings.Split(src, "\n") {
		if i := strings.IndexByte(raw, ';'); i >= 0 {
			raw = raw[:i]
		}
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		op := strings.ToUpper(fields[0])
		var args []operand
		for _, f := range fields[1:] {
			args = append(args, parseOperand(f))
		}
		prog = append(prog, instr{op: op, args: args})
	}
	return prog
}

// runVM executes prog with the given initial memory patches merged into the
// otherwise-zero tape. Returns the OUT stream.
func runVM(prog []instr, patches map[int64]int64) []int64 {
	var regs [6]int64
	value := func(o operand) int64 {
		if o.isReg {
			return regs[o.reg-'A']
		}
		return o.imm
	}

	mem := make(map[int64]int64, len(patches))
	for k, v := range patches {
		mem[k] = v
	}
	var x int64
	var out []int64
	pc := 0
	for pc < len(prog) {
		here := pc
		ins := prog[pc]
		pc++
		switch ins.op {
		case "LEFT":
			x--
		case "RIGHT":
			x++
		case "LOAD":
			regs[ins.args[0].reg-'A'] = mem[x]
		case "SAVE":
			mem[x] = value(ins.args[0])
		case "SET":
			regs[ins.args[0].reg-'A'] = value(ins.args[1])
		case "ADD":
			ri := ins.args[0].reg - 'A'
			regs[ri] = addChecked(regs[ri], value(ins.args[1]))
		case "SUB":
			ri := ins.args[0].reg - 'A'
			regs[ri] = subChecked(regs[ri], value(ins.args[1]))
		case "JMP":
			off := value(ins.args[0])
			if off == 0 {
				panic("s04: JMP with zero offset")
			}
			pc = here + int(off)
		case "JNZ":
			if regs[ins.args[0].reg-'A'] != 0 {
				off := value(ins.args[1])
				if off == 0 {
					panic("s04: JNZ with zero offset")
				}
				pc = here + int(off)
			}
		case "OUT":
			out = append(out, value(ins.args[0]))
		default:
			panic("s04: unknown opcode: " + ins.op)
		}
	}
	return out
}

// weightedSum returns Σ v[i] * 10^i. Panics on int64 overflow so a
// runaway coefficient explodes at the digest step instead of silently
// wrapping in the reported answer.
func weightedSum(v []int64) int64 {
	var sum, w int64 = 0, 1
	for i, x := range v {
		term := mulChecked(x, w)
		sum = addChecked(sum, term)
		if i < len(v)-1 {
			w = mulChecked(w, 10)
		}
	}
	return sum
}

func addChecked(a, b int64) int64 {
	s := a + b
	if (a > 0 && b > 0 && s < 0) || (a < 0 && b < 0 && s >= 0) {
		panic("s04: int64 overflow in ADD")
	}
	return s
}

func subChecked(a, b int64) int64 {
	s := a - b
	if (b > 0 && s > a) || (b < 0 && s < a) {
		panic("s04: int64 overflow in SUB")
	}
	return s
}

func mulChecked(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	p := a * b
	if p/a != b {
		panic("s04: int64 overflow in MUL")
	}
	return p
}

func Solution(input []byte) []string {
	prog := parseProgram(string(input))
	b := runVM(prog, nil)
	n := len(b)

	// Recover A column-by-column via unit-vector probes: A[:,j] = probe(e_j) - b.
	A := make([][]int64, n)
	for i := range A {
		A[i] = make([]int64, n)
	}
	for j := 0; j < n; j++ {
		col := runVM(prog, map[int64]int64{int64(j): 1})
		if len(col) != n {
			panic("s04: probe changed output length")
		}
		for i := 0; i < n; i++ {
			A[i][j] = col[i] - b[i]
		}
	}

	x := solveLinear(A, b)

	patches := make(map[int64]int64, n)
	for j, v := range x {
		patches[int64(j)] = v
	}
	check := runVM(prog, patches)
	for i, v := range check {
		if v != 0 {
			panic("s04: patched output not zero at " + strconv.Itoa(i) + ": " + strconv.FormatInt(v, 10))
		}
	}

	return []string{
		strconv.FormatInt(weightedSum(b), 10),
		strconv.FormatInt(weightedSum(x), 10),
	}
}

// solveLinear returns x such that A*x = -b using big.Rat Gauss elimination.
// Panics if A is singular or x is not an int64-representable integer vector.
func solveLinear(A [][]int64, b []int64) []int64 {
	n := len(A)
	M := make([][]*big.Rat, n)
	for i := 0; i < n; i++ {
		M[i] = make([]*big.Rat, n+1)
		for j := 0; j < n; j++ {
			M[i][j] = new(big.Rat).SetInt64(A[i][j])
		}
		M[i][n] = new(big.Rat).SetInt64(-b[i])
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
			panic("s04: singular matrix")
		}
		if pivot != k {
			M[k], M[pivot] = M[pivot], M[k]
		}
		for i := k + 1; i < n; i++ {
			if M[i][k].Sign() == 0 {
				continue
			}
			factor := new(big.Rat).Quo(M[i][k], M[k][k])
			for j := k; j <= n; j++ {
				term := new(big.Rat).Mul(factor, M[k][j])
				M[i][j].Sub(M[i][j], term)
			}
		}
	}

	x := make([]*big.Rat, n)
	for i := n - 1; i >= 0; i-- {
		s := new(big.Rat).Set(M[i][n])
		for j := i + 1; j < n; j++ {
			term := new(big.Rat).Mul(M[i][j], x[j])
			s.Sub(s, term)
		}
		x[i] = new(big.Rat).Quo(s, M[i][i])
	}

	out := make([]int64, n)
	for i, r := range x {
		if !r.IsInt() {
			panic("s04: non-integer solution")
		}
		num := r.Num()
		if !num.IsInt64() {
			panic("s04: solution overflows int64")
		}
		out[i] = num.Int64()
	}
	return out
}
