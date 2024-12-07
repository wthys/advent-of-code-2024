package day3

import (
	"fmt"
	"regexp"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "3"
}

func (s solution) Part1(input []string) (string, error) {
	instructions, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, instr := range instructions {
		r := instr.Result()
		if instr.IsMul() {
			total += r
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	instructions, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	enabled := true
	for _, instr := range instructions {
		// fmt.Printf("enabled=%v, %v\n", enabled, instr)
		r := instr.Result()
		if instr.IsMul() && enabled {
				total += r
		}

		if enabled {
			enabled = !instr.IsDont()
		} else {
			enabled = instr.IsDo()
		}
	}

	return solver.Solved(total)
}

type Instructions []Instruction
type Instruction struct {
	oper string
	left int
	right int
}

func (instr Instruction) Result() int {
	return instr.left * instr.right
}

func (instr Instruction) IsDo() bool {
	return instr.oper == "do"
}

func (instr Instruction) IsDont() bool {
	return instr.oper == "don't"
}

func (instr Instruction) IsMul() bool {
	return instr.oper == "mul"
}

func (instr Instruction) String() string {
	if instr.IsMul() {
		return fmt.Sprintf("mul(%v,%v)", instr.left, instr.right)
	} else {
		return fmt.Sprintf("%v()", instr.oper)
	}
}

func parseInput(input []string) (Instructions, error) {
	instructions := Instructions{}

	reInstr := regexp.MustCompile("(mul|do|don't)[(]([-+]?[0-9]+,[-+]?[0-9]+)?[)]")

	for _, line := range input {
		matches := reInstr.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			oper := match[1]
			if oper == "mul" {
				args := util.ExtractInts(match[2])
				instructions = append(instructions, Instruction{oper, args[0], args[1]})
			} else {
				instructions = append(instructions, Instruction{oper, 0, 0})
			}
		}
	}

	if len(instructions) == 0 {
		return nil, fmt.Errorf("no instructions found")
	}

	return instructions, nil
}