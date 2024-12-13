package day11

import (
	"fmt"
	// "math"
	"strconv"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "11"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	stones, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	return solver.Solved(stones.BlinkN(25))
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	stones, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	return solver.Solved(stones.BlinkN(75))
}

type (
	Stones []Stone
	Stone int
)

func (s Stone) Blink() Stones {
	if s == Stone(0) {
		return Stones{Stone(1)}
	}

	digits := fmt.Sprint(s)
	if len(digits) % 2 == 0 {
		halfway := len(digits)/2
		left, _ := strconv.Atoi(digits[:halfway])
		right, _ := strconv.Atoi(digits[halfway:])
		return Stones{Stone(left), Stone(right)}
	}

	return Stones{Stone(s * 2024)}
}

type BlinkNKey struct {
	stone Stone
	n int
}

var (
	blinkNMemo = map[BlinkNKey]int{}
)

func (stones Stones) Blink() Stones {
	newStones := Stones{}
	for _, stone := range stones {
		newStones = append(newStones, stone.Blink()...)
	}
	return newStones
}

func (stones Stones) BlinkN(n int) int {
	catalog := map[Stone]int{}
	for _, s := range stones {
		catalog[s] += 1
	}

	for range n {
		newCatalog := map[Stone]int{}
		for s, occ := range catalog {
			for _, stone := range s.Blink() {
				newCatalog[stone] += occ
			}
		}
		catalog = newCatalog
	}

	count := 0
	for _, occ := range catalog {
		count += occ
	}
	return count
}

func parseInput(input []string) (Stones, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("no stones found")
	}
	values := util.ExtractInts(input[0])
	if len(values) == 0 {
		return nil, fmt.Errorf("no stones found")
	}

	return toStones(values), nil
}

func toStones(values []int) Stones {
	stones := Stones{}
	for _, v := range values {
		stones = append(stones, Stone(v))
	}
	return stones
}

func pow10(n int) int {
	if n <= 0 {
		return 1
	}
	return 10 * pow10(n - 1)
}