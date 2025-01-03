package day19

import (
	"testing"
)

type caseCount struct {
	target string
	patterns []string
	expected int
}

func TestCounterCount(t *testing.T) {
	patterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}
	cases := []caseCount{
		{"brwrr", patterns, 2},
		{"bggr", patterns, 1},
		{"gbbr", patterns, 4},
		{"rrbgbr", patterns, 6},
		{"bwurrg", patterns, 1},
		{"brgr", patterns, 2},
		{"ubwu", patterns, 0},
		{"bbrgwb", patterns, 0},
	}

	for _, cs := range cases {
		counter := NewCounter(cs.patterns)
		actual := counter.Count(cs.target)
		if actual != cs.expected {
			t.Errorf("%v should have %v ways, got %v", cs.target, cs.expected, actual)
		}
	}
}