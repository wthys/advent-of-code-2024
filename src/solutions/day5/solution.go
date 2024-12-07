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

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	rules, updates, err := parseInput(input)
	if err !=  nil {
		return solver.Error(err)
	}

	total := 0
	for _, update := range updates {
		if checkAll(update, rules) {
			total += update.Middle()
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	rules, updates, err := parseInput(input)
	if err !=  nil {
		return solver.Error(err)
	}

	total := 0
	for _, update := range updates {
		if !checkAll(update, rules) {
			filtered := filterRules(update, rules)
			corrected := rearrange(update, filtered)
			total += corrected.Middle()
		}
	}

	return solver.Solved(total)
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

	return rightIdx > leftIdx
}

func (r Rule) String() string {
	return fmt.Sprintf("%v|%v", r.left, r.right)
}

func (r Rule) AppliesTo(u Update) bool {
	leftFound := false
	rightFound := false

	for _, nr := range u {
		leftFound = leftFound || r.left == nr
		rightFound = rightFound || r.right == nr
		if leftFound && rightFound {
			return true
		}
	}

	return false
}

func (u Update) Middle() int {
	return u[len(u) / 2]
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
			if len(values) % 2 != 1 {
				return nil, nil, fmt.Errorf("#%v : not an odd number of updates : %v", no+1, line)
			}
			updates = append(updates, Update(values))
		}
	}

	return rules, updates, nil
}

func rearrange(u Update, rules Rules) Update {
	extended := extendUpdate(Update{}, u, rules)
	return extended[0]
}

func extendUpdate(u Update, rem Update, rules Rules) Updates {
	if len(rem) == 0 {
		if checkAll(u, rules) {
			return Updates{u}
		}
		return Updates{}
	}

	if !checkAll(u, rules) {
		return Updates{}
	}

	x := rem[0]
	if len(u) == 0 {
		return extendUpdate(Update{x}, rem[1:], rules)
	}

	tests := Updates{}
	for idx, _ := range u {
		cand := append(append(append(Update{}, u[:idx]...), x), u[idx:]...)
		tests = append(tests, extendUpdate(cand, rem[1:], rules)...)
	}
	tests = append(tests, extendUpdate(append(u, x), rem[1:], rules)...)

	return tests
}

func checkAll(u Update, rules Rules) bool {
	for _, rule := range rules {
		if !rule.Check(u) {
			return false
		}
	}
	return true
}

func filterRules(u Update, rules Rules) Rules {
	rs := Rules{}
	for _, rule := range rules {
		if rule.AppliesTo(u) {
			rs = append(rs, rule)
		}
	}
	return rs
}