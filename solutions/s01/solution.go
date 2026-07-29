// Package s01 contains a solution to puzzle 01.
package s01

import (
	"strconv"
	"strings"
)

// parseCents converts a fixed-point number with exactly two decimals (e.g. "-3.14")
// into an integer count of cents (-314).
func parseCents(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	neg := false
	switch s[0] {
	case '-':
		neg = true
		s = s[1:]
	case '+':
		s = s[1:]
	}
	dot := strings.IndexByte(s, '.')
	if dot < 0 {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, false
		}
		if neg {
			n = -n
		}
		return n * 100, true
	}
	intPart := s[:dot]
	fracPart := s[dot+1:]
	if len(fracPart) != 2 {
		return 0, false
	}
	if intPart == "" {
		intPart = "0"
	}
	i, err := strconv.Atoi(intPart)
	if err != nil {
		return 0, false
	}
	f, err := strconv.Atoi(fracPart)
	if err != nil {
		return 0, false
	}
	v := i*100 + f
	if neg {
		v = -v
	}
	return v, true
}

func Solution(input []byte) []string {
	const target = 700 // 7.00 in cents

	var nums []int
	for _, line := range strings.Split(string(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		v, ok := parseCents(line)
		if !ok {
			continue
		}
		nums = append(nums, v)
	}

	// Part 1: count unordered pairs (i<j) with nums[i]+nums[j] == 700.
	pairs := int64(0)
	freq := map[int]int64{}
	for _, v := range nums {
		pairs = addCheckedS01(pairs, freq[target-v])
		freq[v]++
	}

	// Part 2: count signed subsets (each element omitted, added, or subtracted)
	// summing to 700. f(i, t) = number of ways over nums[i:]; empty selection
	// contributes only to t=0, so no adjustment needed for target=700.
	n := len(nums)
	suffAbs := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		a := nums[i]
		if a < 0 {
			a = -a
		}
		suffAbs[i] = suffAbs[i+1] + a
	}

	type key struct {
		i, t int
	}
	memo := map[key]int64{}
	var f func(i, t int) int64
	f = func(i, t int) int64 {
		if t > suffAbs[i] || -t > suffAbs[i] {
			return 0
		}
		if i == n {
			if t == 0 {
				return 1
			}
			return 0
		}
		k := key{i, t}
		if v, ok := memo[k]; ok {
			return v
		}
		a := f(i+1, t)
		b := f(i+1, t-nums[i])
		c := f(i+1, t+nums[i])
		v := addCheckedS01(addCheckedS01(a, b), c)
		memo[k] = v
		return v
	}
	total := f(0, target)

	return []string{
		strconv.FormatInt(pairs, 10),
		strconv.FormatInt(total, 10),
	}
}

// addCheckedS01 adds two non-negative int64 count values and panics on
// int64 overflow. The Part 2 count is bounded by 3^n over all signed
// subsets; committed inputs sit well below 2^63-1, but the guard makes
// any future generator change that would blow past int64 explode loudly
// instead of silently wrapping.
func addCheckedS01(a, b int64) int64 {
	s := a + b
	if (a > 0 && b > 0 && s < 0) || (a < 0 && b < 0 && s >= 0) {
		panic("s01: int64 overflow in count accumulation")
	}
	return s
}
