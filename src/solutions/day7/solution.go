package day7

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
		max := 1 << (len(calibration.nums)-1)
		operators := 0
		for operators < max {
			if calibration.value == evaluate(operators, calibration.nums) {
				valid += 1
			}
			operators += 1
		}
		if valid > 0 {
			fmt.Printf("%v has %v valid interpretations!\n", calibration, valid)
			total += calibration.value
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	return solver.NotImplemented()
}

type Calibrations []Calibration

type Calibration struct {
	value int
	nums []int
}

func getOperators(operators int, n int) []string {
	opers := []string{}
	n -= 1
	for n >= 0 {
		opers = append(opers, getOperator(operators, n))
		n -= 1
	}
	return opers
}

func getOperator(operators int, idx int) string {
	opermap := map[bool]string{
		false: "+",
		true: "*",
	}
	return opermap[(operators & (1<<idx)) > 0]
}

func operate(oper string, left int, right int) int {
	// fmt.Printf("______ %v %v %v\n", left, oper, right)
	if oper == "+" {
		return left + right
	}

	return left * right
}

func evaluate(operators int, nums []int) int {
	// fmt.Printf("evaluate(%v, %v) = ", getOperators(operators, len(nums)-1), nums)
	opers := len(nums)-2
	idx := 1
	total := nums[0]
	for true {
		total = operate(getOperator(operators, opers), total, nums[idx])
		opers -= 1
		idx += 1
		if opers < 0 {
			break
		}
	}
	// fmt.Println(total)
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

