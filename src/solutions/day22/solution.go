package day22

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "22"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	initials, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, initial := range initials {
		nth := NthSecret(initial, 2000)
		opts.Debugf("__ %v --(x2000)--> %v\n", initial, nth)
		total += nth
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	return solver.NotImplemented()
}

func parseInput(input []string) ([]int, error) {
	seeds := []int{}
	for _, line := range input {
		seeds = append(seeds, util.ExtractInts(line)...)
	}

	if len(seeds) == 0 {
		return nil, fmt.Errorf("no initial numbers found")
	}

	return seeds, nil
}

const (
	MODULO = 16777216
)

func NextSecret(current int) int {
	stage1 := ((current * 64) ^ current) % MODULO
	stage2 := ((stage1 / 32) ^ stage1) % MODULO
	stage3 := ((stage2 * 2048) ^ stage2) % MODULO
	return stage3
}

func NthSecret(secret int, n int) int {
	if n <= 0 {
		return secret
	}

	return NthSecret(NextSecret(secret), n - 1)
}