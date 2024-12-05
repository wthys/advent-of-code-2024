package day5

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
	return "5"
}

func (s solution) Part1(input []string) (string, error) {
	rules, updates, err := parseInput(input)
	if err !=  nil {
		return solver.Error(err)
	}

	total := 0
	for _, update := range updates {
		correct := true
		for _, rule := range rules {
			if !rule.Check(update) {
				correct = false
				break
			}
		}
		if correct {
			total += update.Middle()
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	return solver.NotImplemented()
}

type Rules []Rule
type Rule struct {
	left int
	right int
}

type Updates []Update
type Update []int

func (r Rule) Check(u Update) bool {
	leftIdx := -1
	rightIdx := len(u)
	for idx, nr := range u {
		if nr == r.left {
			leftIdx = idx
		}
		if nr ==  r.right {
			rightIdx = idx
		}
	}

	// fmt.Printf("DEBUG: %v : %v - %v : %v\n", r, leftIdx, rightIdx, rightIdx > leftIdx)
	return rightIdx > leftIdx
}

func (r Rule) String() string {
	return fmt.Sprintf("%v|%v", r.left, r.right)
}

func (u Update) Middle() int {
	return u[len(u)/2]
}

func parseInput(input []string) (Rules, Updates, error) {
	rules := Rules{}
	updates := Updates{}

	rulesDone := false

	for no, line := range input {
		values := util.ExtractInts(line)
		if len(values) == 0 {
			rulesDone = true
			continue
		}		

		if !rulesDone {
			if len(values) != 2 {
				return nil, nil, fmt.Errorf("#%v : invalid sorting rule : %v", no+1, line)
			}
			rules = append(rules, Rule{values[0], values[1]})
		} else {
			if len(values)%2 != 1 {
				return nil, nil, fmt.Errorf("#%v : not an odd number of updates : %v", no+1, line)
			}
			updates = append(updates, Update(values))
		}
	}

	return rules, updates, nil
}