package day9

import (
	"testing"
)

type caseCompact struct {
	input Sectors
	expected Sectors
}

func TestSectorsCompact(t *testing.T) {
	cases := []caseCompact{
		{
			Sectors{{1,1}, {1,1}, {1,1}},
			Sectors{{3,1}},
		},
		{
			Sectors{{4,-1}, {3,-1}, {2,1}, {1,1}, {2,-1}, {23,-1}},
			Sectors{{7,-1}, {3,1}, {25,-1}},
		},
	}

	for _, cs := range cases {
		actual := cs.input.Compact()

		if len(actual) != len(cs.expected) {
			t.Fatalf("Compact should return %v elements, got %v", len(cs.expected), len(actual))
		}

		for idx, e := range cs.expected {
			a := actual[idx]
			if a != e {
				t.Fatalf("Compact should return %v @ %v, got %v", e, idx, a)
			}
		}
	}
}