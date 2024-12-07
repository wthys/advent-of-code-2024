package day1

import (
	"fmt"
	"slices"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

type Notes struct {
	left []int
	right []int
}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "1"
}

func (s solution) Part1(input []string) (string, error) {
	notes, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	left := slices.Clone(notes.left)
	right := slices.Clone(notes.right)
	slices.Sort(left)
	slices.Sort(right)

	dist := 0
	for i, l := range left {
		dist += util.Abs(l - right[i])
	}


	return solver.Solved(dist)
}

func (s solution) Part2(input []string) (string, error) {
	notes, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	rCounter := map[int]int{}
	for _, b := range notes.right {
		_, ok := rCounter[b]
		if ok {
			rCounter[b] += 1
		} else {
			rCounter[b] = 1
		}
	}

	total := 0
	for _, a := range notes.left {
		n, ok := rCounter[a]
		if ok {
			total += a * n
		}
	}

	return solver.Solved(total)
}

func parseInput(input []string) (*Notes, error) {
	notes := &Notes{[]int{},[]int{}}
	for _, line := range input {
		nums := util.ExtractInts(line)
		if len(nums) == 0 {
			continue
		}
		if len(nums) < 2 {
			return nil, fmt.Errorf("not enough numbers")
		}
		notes.left = append(notes.left, nums[0])
		notes.right = append(notes.right, nums[1])
	}

	return notes, nil
}
