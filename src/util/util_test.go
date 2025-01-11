package util

import (
	"fmt"
	"testing"
)

type (
	testInt struct {
		Input int
		Want  int
	}

	testFloat struct {
		Input float64
		Want  float64
	}
)

func TestSignInt(t *testing.T) {
	cases := []testInt{
		{int(5), int(1)},
		{int(0), int(0)},
		{int(-245), int(-1)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Sign(cs.Input)
			if s != cs.Want {
				t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestSignFloat(t *testing.T) {
	cases := []testFloat{
		{float64(3.14), float64(1.0)},
		{float64(0.0), float64(0.0)},
		{float64(-5.256), float64(-1.0)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Sign(cs.Input)
			if s != cs.Want {
				t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestAbsInt(t *testing.T) {
	cases := []testInt{
		{int(5), int(5)},
		{int(0), int(0)},
		{int(-245), int(245)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Abs(cs.Input)
			if s != cs.Want {
				t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestAbsFloat(t *testing.T) {
	cases := []testFloat{
		{float64(3.14), float64(3.14)},
		{float64(0.0), float64(0.0)},
		{float64(-5.256), float64(5.256)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Abs(cs.Input)
			if s != cs.Want {
				t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func hash(array []int) int {
	h := 0
	for _, value := range array {
		h = 13*h + value
	}
	return h
}

func TestPermutationDo(t *testing.T) {
	array := []int{1, 2, 3}

	check := map[int]bool{
		hash([]int{1, 2, 3}): false,
		hash([]int{1, 3, 2}): false,
		hash([]int{2, 1, 3}): false,
		hash([]int{2, 3, 1}): false,
		hash([]int{3, 2, 1}): false,
		hash([]int{3, 1, 2}): false,
	}

	PermutationDo(3, array, func(perm []int) {
		check[hash(perm)] = true
	})

	for value, seen := range check {
		if !seen {
			t.Fatalf("PermutationDo(3, %v, ...) should produce %v, but was not seen", array, value)
		}
	}

}

func TestForEach(t *testing.T) {
	array := []int{1, 2, 3}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
	}

	ForEach(array, func(value int) {
		check[value] = true
	})

	for value, seen := range check {
		if !seen {
			t.Fatalf("Do(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}

func TestForEachError(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
	}

	err := ForEachError(array, func(value int) error {
		if value > 3 {
			return fmt.Errorf("value %v is too large", value)
		}
		check[value] = true
		return nil
	})

	if err == nil {
		t.Fatalf("ForEachError(%v, ...) should produce an error but none was seen", array)
	}

	for value, seen := range check {
		if !seen && value <= 3 {
			t.Fatalf("ForEachError(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}

func TestForEachStopping(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
	}

	stopped := ForEachStopping(array, func(value int) bool {
		if value > 3 {
			return false
		}
		check[value] = true
		return true
	})

	if !stopped {
		t.Fatalf("ForEachStopping(%v, ...) should produce %v but was %v", array, true, stopped)
	}

	for value, seen := range check {
		if !seen && value <= 3 {
			t.Fatalf("ForEachStopping(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}

func TestExtractInts(t *testing.T) {
	line := "1,2,+3-4++6fdsfdsfgdsg7,88,-999,19234"
	expected := []int{1,2,3,-4,6,7,88,-999,19234}

	actual := ExtractInts(line)
	if len(actual) != len(expected) {
		t.Fatalf("ExtractInts(%q) should have %v element, got %v", line, len(expected), len(actual))
	}

	for idx, a := range actual {
		e := expected[idx]

		if a != e {
			t.Fatalf("ExtractInts(%q) should have %v at index %v, got %v", line, e, idx, a)
		}
	}
}

type testComb struct {
	vals []int
	k int
	expected [][]int
}

func combCase(vals []int, k int, expected [][]int) testComb {
	return testComb{
		vals,
		k,
		expected,
	}
}

func TestCombinationDo(t * testing.T) {
	cases := []testComb{
		combCase([]int{1,2}, 1, [][]int{{1},{2}}),
		combCase([]int{3,4}, 2, [][]int{{3,3},{3,4},{4,3},{4,4}}),
		combCase([]int{5,6}, 3, [][]int{{5,5,5}, {5,5,6}, {5,6,5}, {5,6,6}, {6,5,5}, {6,5,6}, {6,6,5}, {6,6,6}}),
		combCase([]int{7,8,9}, 2, [][]int{{7,7}, {7,8}, {7,9}, {8,7}, {8,8}, {8,9}, {9,7}, {9,8}, {9,9}}),
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("CombinationDo(%v, %v)", cs.vals, cs.k), func(tt *testing.T) {
			seen := map[int]int{}
			for _, val := range cs.expected {
				seen[hash(val)] = 0
			}

			CombinationDo(cs.vals, cs.k, func (cand []int) {
				seen[hash(cand)] += 1
			})

			for _, vals := range cs.expected {
				n, _ := seen[hash(vals)]
				if n != 1 {
					tt.Errorf("%v was seen %v times, expected 1", vals, n)
				}
			}
		})
	}
}

func TestCombinationNoRepeatDo(t *testing.T) {
	cases := []testComb{
		combCase([]int{1,2}, 1, [][]int{{1}, {2}}),
		combCase([]int{1,2}, 2, [][]int{{1,2}}),
		combCase([]int{1,2,3}, 2, [][]int{{1,2}, {1,3}, {2,3}}),
		combCase([]int{1,2,3,4}, 2, [][]int{{1,2}, {1,3}, {1,4}, {2,3}, {2,4}, {3,4}}),
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("CombinationNoRepeatDo(%v, %v)", cs.vals, cs.k), func(tt *testing.T) {
			seen := map[int]int{}
			for _, val := range cs.expected {
				seen[hash(val)] = 0
			}

			CombinationNoRepeatDo(cs.vals, cs.k, func(cand []int) {
				seen[hash(cand)] += 1
			})

			for _, vals := range cs.expected {
				n, _ := seen[hash(vals)]
				if n != 1 {
					tt.Errorf("%v was seen %v times, expected 1", vals, n)
				}
			}
		})
	}
}

type caseStoI struct {
	input []string
	expected []int
	err bool
}

func TestStringsToInts(t *testing.T) {
	cases := []caseStoI{
		{[]string{"1","2","3"}, []int{1,2,3}, false},
		{[]string{"1","2","a"}, nil, true},
		{[]string{"-11","+22","33"}, []int{-11,22,33}, false},
		{[]string{}, []int{}, false},
	}

	for _, cs := range cases {
		actual, err := StringsToInts(cs.input)
		if cs.err && err == nil {
			t.Fatalf("StringsToInts(%v) expected an error, got none", cs.input)
		}

		if len(cs.expected) != len(actual) {
			t.Fatalf("StringsToInts(%v) expected to have %v elements, got %v", cs.input, len(cs.expected), len(actual))
		}

		for idx, e := range cs.expected {
			if e != actual[idx] {
				t.Fatalf("StringsToInts(%v) expected %v @ %v, got %v", cs.input, e, idx, actual[idx])
			}
		}
	}
}

type caseExRgx struct {
	pattern string
	line string
	expected []string
}

func TestExtractRegex(t *testing.T) {
	cases := []caseExRgx{
		{"ab", "abaaabbabbaba", []string{"ab", "ab", "ab", "ab"}},
		{"[0-9]", "1jh2ki54jhh9", []string{"1", "2", "5", "4", "9"}},
		{"hello", "world", []string{}},
		{"[-+]?[0-9]+", "123,+45adf67,-8,9", []string{"123", "+45", "67", "-8", "9"}},
	}

	for _, cs := range cases {
		actual := ExtractRegex(cs.pattern, cs.line)
		if len(actual) != len(cs.expected) {
			t.Fatalf("ExtractRegex(%q, %q) expected %v elements, got %v", cs.pattern, cs.line, len(cs.expected), len(actual))
		}

		for idx, e := range cs.expected {
			a := actual[idx]
			if e != a {
				t.Fatalf("ExtractRegex(%q, %q) expected %v @ %v, got %v", cs.pattern, cs.line, e, idx, a)
			}
		}
	}
}