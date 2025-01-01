package day17

import (
	"fmt"
	"strings"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "17"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	computer, program, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	computer.Run(program)

	return solver.Solved(formatOutput(computer))
}

const (
	INFINITE = int((^uint(0)) >> 1)
)

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	computer, program, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	index := len(program)-1
	indexCandidates := map[int][]int{index+1: []int{0}}
	watchdog := func(c Computer) bool {
		return len(program) >= len(c.Output)
	}

	for index >= 0 {
		for n := range 8 {
			tail := Program(program[index:])
			for _, v := range indexCandidates[index+1] {
				inst := computer
				inst.RegA = n + 8 * v
				ok := inst.RunWatchdog(program, watchdog)

				if ok && tail.Equals(Program(inst.Output)) {
					indexCandidates[index] = append(indexCandidates[index], n + 8 * v)
				}
			}
		}
		index--
	}

	lowest := INFINITE
	for _, v := range indexCandidates[0] {
		if v < lowest {
			lowest = v
		}
	}

	return solver.Solved(lowest)
}

type (
	Instr struct {
		opcode int
		operand int
	}
	Program []int
	Computer struct {
		IP int
		RegA int
		RegB int
		RegC int
		Output []int
	}

	ComputerBuilder struct {
		ip int
		regA int
		regB int
		regC int
		output []int
	}
)

func (cb ComputerBuilder) RegA(value int) ComputerBuilder {
	new := cb
	new.regA = value
	return new
}

func (cb ComputerBuilder) RegB(value int) ComputerBuilder {
	new := cb
	new.regB = value
	return new
}

func (cb ComputerBuilder) RegC(value int) ComputerBuilder {
	new := cb
	new.regC = value
	return new
}

func (cb ComputerBuilder) IP(value int) ComputerBuilder {
	new := cb
	new.ip = value
	return new
}

func (cb ComputerBuilder) Output(value []int) ComputerBuilder {
	new := cb
	new.output = value
	return new
}

func (cb ComputerBuilder) Build() Computer {
	return Computer{cb.ip, cb.regA, cb.regB, cb.regC, cb.output}
}

func (c Computer) resolveOperand(operand int) int {
	switch operand {
	case 4:
		return c.RegA
	case 5:
		return c.RegB
	case 6:
		return c.RegC
	case 7:
		panic("invalid operand")
	default:
		return operand
	}
}

func (c *Computer) AddOutput(value int) {
	c.Output = append(c.Output, value)
}

func (c *Computer) Run(program Program) {
	c.RunWatchdog(program, func(_ Computer) bool {
		return true
	})
}

func (c *Computer) RunWatchdog(program Program, watchdog func(Computer) bool) bool {
	ok := watchdog(*c)
	for c.IP < len(program) && ok {
		instr := program.GetInstruction(c.IP)
		instr.Execute(c)
		ok = watchdog(*c)
	}
	return ok
}

func formatOutput(c Computer) string {
	strs := []string{}
	for _, n := range c.Output {
		strs = append(strs, fmt.Sprint(n))
	}
	return strings.Join(strs, ",")
}

func (p Program) GetInstruction(index int) Instr {
	return Instr{p[index], p[index+1]}
}

func (p Program) Equals(o Program) bool {
	if len(p) != len(o) {
		return false
	}

	for idx, _ := range p {
		if p[idx] != o[idx] {
			return false
		}
	}

	return true
}

var (
	OPCODEMAPPER = []string{"adv", "bxl", "bst", "jnz", "bxc", "out", "bdv", "cdv"}
)

func (i Instr) String() string {
	return fmt.Sprintf("%v,%v", OPCODEMAPPER[i.opcode], i.operand)
}

func (i Instr) Execute(comp *Computer) {
	switch i.opcode {
	case 0: //adv
		comp.RegA = div(comp.RegA, comp.resolveOperand(i.operand))
		comp.IP += 2
	case 1: //bxl
		val := comp.RegB
		mask := i.operand
		comp.RegB = val ^ mask
		comp.IP += 2
	case 2: //bst
		val := comp.resolveOperand(i.operand)
		comp.RegB = val % 8
		comp.IP += 2
	case 3: //jnz
		if comp.RegA == 0 {
			comp.IP += 2
		} else {
			comp.IP = i.operand
		}
	case 4: //bxc
		comp.RegB = comp.RegB ^ comp.RegC
		comp.IP += 2
	case 5: //out
		val := comp.resolveOperand(i.operand) % 8
		comp.AddOutput(val)
		comp.IP += 2
	case 6: //bdv
		comp.RegB = div(comp.RegA, comp.resolveOperand(i.operand))
		comp.IP += 2
	case 7: //cdv
		comp.RegC = div(comp.RegA, comp.resolveOperand(i.operand))
		comp.IP += 2
	}
}

func div(numerator int, denumPower int) int {
	return numerator / pow(2, denumPower)
}

func pow(base, exp int) int {
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}

func parseInput(input []string) (Computer, Program, error) {
	builder := ComputerBuilder{}
	
	computer := builder.RegA(util.ExtractInts(input[0])[0]).RegB(util.ExtractInts(input[1])[0]).RegC(util.ExtractInts(input[2])[0]).Build()
	
	program := Program(util.ExtractInts(input[4]))

	if len(program) == 0 {
		return Computer{}, nil, fmt.Errorf("no program found")
	}

	return computer, program, nil
}