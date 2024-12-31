package day16

import (
	"testing"
	L "github.com/wthys/advent-of-code-2024/location"
)

type caseMoveCost struct {
	from Step
	to Step
	expected int
	err bool
}

func TestMoveCost(t *testing.T) {
	cases := []caseMoveCost{
		{
			Step{L.New(0,0), L.New(1,0)},
			Step{L.New(1,0), L.New(1,0)},
			1,
			false,
		},
		{
			Step{L.New(0,0), L.New(0,1)},
			Step{L.New(1,0), L.New(1,0)},
			1001,
			false,
		},
		{
			Step{L.New(0,0), L.New(-1,0)},
			Step{L.New(1,0), L.New(1,0)},
			2001,
			false,
		},
		{
			Step{L.New(0,0), L.New(-1,0)},
			Step{L.New(1,1), L.New(1,0)},
			0,
			true,
		},
	}

	for _, cs := range cases {
		actual, err := cs.from.MoveCost(cs.to)

		if err != nil && !cs.err {
			t.Fatalf("%v.MoveCost(%v) was not expected to fail, got %v", cs.from, cs.to, err)
		}
		if err == nil && cs.err {
			t.Fatalf("%v.MoveCost(%v) was expected to fail and dit not", cs.from, cs.to)
		}
		if cs.expected != actual {
			t.Fatalf("%v.MoveCost(%v) should return %v, got %v", cs.from, cs.to, cs.expected, actual)
		}
	}
}

type caseMove struct {
	from Step
	dir L.Location
	expected Step
}

func TestMove(t *testing.T) {
	cases := []caseMove{
		{
			Step{L.New(1,1), L.New(0,1)},
			L.New(1,0),
			Step{L.New(2,1), L.New(1,0)},
		},
		{
			Step{L.New(1,1), L.New(0,1)},
			L.New(2,0),
			Step{L.New(2,1), L.New(1,0)},
		},
		{
			Step{L.New(1,1), L.New(1,0)},
			L.New(1,0),
			Step{L.New(2,1), L.New(1,0)},
		},
	}

	for _, cs := range cases {
		actual := cs.from.Move(cs.dir)
		if actual.Dir != cs.expected.Dir {
			t.Fatalf("%v.Move(%v) result direction should be %v, got %v", cs.from, cs.dir, cs.expected.Dir, actual.Dir)
		}
		if actual.Pos != cs.expected.Pos {
			t.Fatalf("%v.Move(%v) result position should be %v, got %v", cs.from, cs.dir, cs.expected.Pos, actual.Pos)
		}
	}
}