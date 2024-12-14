package day13

import (
	"fmt"
	"math"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	L "github.com/wthys/advent-of-code-2024/location"
	"gonum.org/v1/gonum/mat"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "13"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	machines, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, machine := range machines {
		a, b := machine.Presses()
		close := a >= 0 && b >= 0 && machine.A.Scale(a).Add(machine.B.Scale(b)) == machine.Prize
		if close {
			total += 3*a + b
		}
		opts.Debugf("Machine: %v, presses=Ax%v, Bx%v, SUCCESS=%v\n", machine, a, b, close)
	}


	return solver.Solved(total)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	machines, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}
	
	correction := L.New(10000000000000,10000000000000)

	total := 0
	for _, machine := range machines {
		machine.Prize = machine.Prize.Add(correction)
		a, b := machine.Presses()
		close := a >= 0 && b >= 0 && machine.A.Scale(a).Add(machine.B.Scale(b)) == machine.Prize
		if close {
			total += 3*a + b
		}
		opts.Debugf("Machine: %v, presses=Ax%v, Bx%v, SUCCESS=%v\n", machine, a, b, close)
	}


	return solver.Solved(total)
}

type (
	Machines []Machine
	Machine struct {
		A L.Location
		B L.Location
		Prize L.Location
	}
)

func (m Machine) Presses() (int, int) {
	p := mat.NewDense(2, 2, []float64{
		float64(m.A.X), float64(m.B.X),
		float64(m.A.Y), float64(m.B.Y),
	})

	var pI mat.Dense
	pI.Inverse(p)

	x := mat.NewDense(2, 1, []float64{
		float64(m.Prize.X),
		float64(m.Prize.Y),
	})

	var xB mat.Dense
	xB.Mul(&pI, x)

	a := xB.At(0, 0)
	b := xB.At(1, 0)
	// fmt.Printf("DEBUG: (%v, %v) -> (%v, %v) / (%v, %v)\n", a, b, int(a), int(b), int(math.Round(a)), int(math.Round(b)))
	return int(math.Round(a)), int(math.Round(b))
}

func parseInput(input []string) (Machines, error) {
	machines := Machines{}

	machine := Machine{}
	lastStart := 0

	for idx, line := range input {
		values := util.ExtractInts(line)
		if len(values) == 0 {
			lastStart = idx + 1
			continue
		}

		loc := L.New(values[0], values[1])
		switch idx - lastStart {
			case 0:
				machine.A = loc
			case 1:
				machine.B = loc
			case 2:
				machine.Prize = loc
				machines = append(machines, machine)
			default:
				lastStart = idx + 1
		} 
	}

	if len(machines) == 0 {
		return nil, fmt.Errorf("no machines found")
	}

	return machines, nil
}