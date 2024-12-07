package day7

import (
	"fmt"
	"strconv"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "7"
}

func (s solution) Part1(input []string) (string, error) {
	calibrations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, calibration := range calibrations {
		valid := 0
		util.CombinationDo([]string{"+","*"}, len(calibration.nums)-1, func(cand []string) {
			if valid > 0 {
				return
			}
			
			if calibration.value == evaluate(cand, calibration.nums) {
				valid += 1
			}
		})

		if valid > 0 {
			// fmt.Printf("%v has %v valid interpretations!\n", calibration, valid)
			total += calibration.value
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	calibrations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, calibration := range calibrations {
		valid := 0
		util.CombinationDo([]string{"+","*","|"}, len(calibration.nums)-1, func(cand []string) {
			if valid > 0 {
				return
			}

			if calibration.value == evaluate(cand, calibration.nums) {
				valid += 1
			}
		})

		if valid > 0 {
			// fmt.Printf("%v has %v valid interpretations!\n", calibration, valid)
			total += calibration.value
		}
	}

	return solver.Solved(total)
}

type Calibrations []Calibration

type Calibration struct {
	value int
	nums []int
}

func operate(oper string, left int, right int) int {
	// fmt.Printf("______ %v %v %v\n", left, oper, right)
	if oper == "+" {
		return left + right
	}

	if oper == "|" {
		n, _ := strconv.Atoi(fmt.Sprintf("%v%v", left, right))
		return n
	}

	return left * right
}

func evaluate(operators []string, nums []int) int {
	total := nums[0]
	for idx, oper := range operators {
		total = operate(oper, total, nums[idx+1])
	}
	return total
}

func parseInput(input []string) (Calibrations, error) {
	calibrations := Calibrations{}

	for no, line := range input {
		nums := util.ExtractInts(line)
		if len(nums) == 0 {
			continue
		}
		if len(nums) < 3 {
			return nil, fmt.Errorf("#%v : not enough values in %q", no, line)
		}

		calibrations = append(calibrations, Calibration{nums[0], nums[1:]})
	}

	if len(calibrations) == 0 {
		return nil, fmt.Errorf("no calibrations found")
	}

	return calibrations, nil
}

