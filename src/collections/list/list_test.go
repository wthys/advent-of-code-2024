package list

import (
	"fmt"
	"testing"
)

type caseLen struct {
	input *List[int]
	expected int
}

func TestLen(t *testing.T) {
	cases := []caseLen{
		{NewFor(0), 0},
		{Empty[int](), 0},
		{New(1), 1},
		{New(1,2,3), 3},
	}
	for _, cs := range cases {
		actual := cs.input.Len()
		if actual != cs.expected {
			t.Errorf("%v expected to have length %v, got %v", cs.input, cs.expected, actual)
		}
	}
}

type caseEquals struct {
	left *List[int]
	right *List[int]
	expected bool
}

func TestEquals(t *testing.T) {
	same := New(1,2,3,4)
	cases := []caseEquals{
		{NewFor(0), Empty[int](), true},
		{New(1), New(2), false},
		{same, same, true},
		{New(1), New(1), true},
		{New(1), New(1,2,3), false},
		{New(2), nil, false},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v.Equals(%v)", cs.left, cs.right), func(tt *testing.T){
			actual := cs.left.Equals(cs.right)
			if actual != cs.expected {
				tt.Errorf("%v.Equals(%v) should be %v, got %v", cs.left, cs.right, cs.expected, actual)
			}
		})
		if cs.right != nil {
			t.Run(fmt.Sprintf("%v.Equals(%v)", cs.right, cs.left), func(tt *testing.T){
				revAct := cs.right.Equals(cs.left)
				if revAct != cs.expected {
					tt.Errorf("%v.Equals(%v) should be %v, got %v", cs.right, cs.left, cs.expected, revAct)
				}
			})
		}
	}
}

type caseGet struct {
	input *List[int]
	index int
	expected int
	found bool
}

func TestGet(t *testing.T) {
	cases := []caseGet{
		{New(1,2,3), 1, 2, true},
		{NewFor(0), 1, 0, false},
		{New(1,2,3), -1, 0, false},
		{New(1,2,3), 5, 0, false},
		{New(1,2,3), 0, 1, true},
		{New(1,2,3), 1, 2, true},
		{New(1,2,3), 2, 3, true},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v.Get(%v)", cs.input, cs.index), func(tt *testing.T) {
			actual, ok := cs.input.Get(cs.index)
			if ok != cs.found {
				tt.Fatalf("expected second return to be %v, got %v", cs.found, ok)
			}
			if actual != cs.expected {
				tt.Errorf("expected value to be %v, got %v", cs.expected, actual)
			}
		})
	}
}

type caseString struct {
	input *List[int]
	expected string
}

func TestString(t *testing.T) {
	cases := []caseString{
		{NewFor(0), "[ ]"},
		{New(1), "[ 1 ]"},
		{New(1,2,3), "[ 1 2 3 ]"},
	}

	for _, cs := range cases {
		actual := cs.input.String()
		if actual != cs.expected {
			t.Errorf("%v.String() should return %q, got %q", cs.input, cs.expected, actual)
		}
	}
}