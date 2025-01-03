package day22

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	"github.com/wthys/advent-of-code-2024/collections/list"
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
	initials, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	seen := map[string]int{}
	for _, initial := range initials {
		last := initial
		sequence := list.NewFor(0)
		thisSeen := map[string]int{}
		for n := 0; n < 2000; n++ {
			next := NextSecret(last)
			diff := (next % 10) - (last % 10)
			sequence.Append(diff)
			for sequence.Len() > 4 {
				sequence.PopFront()
			}
			if sequence.Len() == 4 {
				s := sequence.String()
				_, ok := thisSeen[s]
				if !ok {
					opts.Debugf("__ adding %v => %v\n", s, next % 10)
					thisSeen[s] = next % 10
				}
			}
			last = next
		}
		for seq, bananas := range thisSeen {
			seen[seq] += bananas
		}
	}

	most := 0
	for _, bananas := range seen {
		if most < bananas {
			most = bananas
		}
	}
	return solver.Solved(most)
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