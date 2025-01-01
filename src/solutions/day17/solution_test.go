package day17

import (
	"testing"
)

type caseExec struct {
	start Computer
	input Program
	validator func(c Computer)
}

func TestExecute(t *testing.T) {
	cases := []caseExec{
		{ComputerBuilder{}.RegC(9).Build(), Program{2,6}, func(c Computer) {
			assertRegB(c, 1, t)
		}},
		{ComputerBuilder{}.RegA(10).Build(), Program{5,0,5,1,5,4}, func(c Computer) {
			assertOutput(c, []int{0,1,2}, t)
		}},
		{ComputerBuilder{}.RegA(2024).Build(), Program{0,1,5,4,3,0}, func(c Computer) {
			assertOutput(c, []int{4,2,5,6,7,7,7,7,3,1,0}, t)
			assertRegA(c, 0, t)
		}},
		{ComputerBuilder{}.RegB(29).Build(), Program{1,7}, func(c Computer) {
			assertRegB(c, 26, t)
		}},
		{ComputerBuilder{}.RegB(2024).RegC(43690).Build(), Program{4,0}, func(c Computer) {
			assertRegB(c, 44354, t)
		}},
	}

	for _, cs := range cases {
		c := cs.start
		c.Run(cs.input)
		cs.validator(c)
	}
}

func assertRegA(c Computer, expected int, t *testing.T) {
	if c.RegA != expected {
		t.Fatalf("expected %v in A register, got %v", expected, c.RegA)
	}
}

func assertRegB(c Computer, expected int, t *testing.T) {
	if c.RegB != expected {
		t.Fatalf("expected %v in B register, got %v", expected, c.RegB)
	}
}

func assertRegC(c Computer, expected int, t *testing.T) {
	if c.RegC != expected {
		t.Fatalf("expected %v in C register, got %v", expected, c.RegC)
	}
}

func assertOutput(c Computer, expected []int, t *testing.T) {
	if len(c.Output) != len(expected) {
		t.Fatalf("expected %v outputs, got %v", len(expected), len(c.Output))
	}
	for idx, e := range expected {
		actual := c.Output[idx]
		if e != actual {
			t.Fatalf("expected %v @ %v, got %v", e, idx, actual)
		}
	}
}

type caseProgramEquals struct {
	left Program
	right Program
	expected bool
}

func TestProgramEquals(t *testing.T) {
	cases := []caseProgramEquals{
		{Program{1,2,3}, Program{1}, false},
		{Program{}, Program{}, true},
		{Program{1,2}, Program{1,2}, true},
		{Program{1,2}, Program{2,1}, false},
		{Program{1}, Program{1,2}, false},
	}

	for _, cs := range cases {
		actual := cs.left.Equals(cs.right)
		if cs.expected != actual {
			t.Errorf("%v.Equals(%v) should be %v, got %v", cs.left, cs.right, cs.expected, actual)
		}
	}
}