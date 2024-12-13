package day11

import (
	"testing"
)

type caseSB struct {
	input Stone
	expected Stones
}

func TestStoneBlink(t *testing.T) {
	cases := []caseSB{
		{Stone(0), Stones{Stone(1)}},
		{Stone(1), Stones{Stone(2024)}},
		{Stone(11), Stones{Stone(1), Stone(1)}},
		{Stone(101), Stones{Stone(101*2024)}},
	}

	for _, cs := range cases {
		actual := cs.input.Blink()
		if len(actual) != len(cs.expected) {
			t.Fatalf("%v.Blink() expected %v elements, got %v (%v)", cs.input, len(cs.expected), len(actual), actual)
		}

		for idx, e := range cs.expected {
			a := actual[idx]
			if e != a {
				t.Fatalf("%v.Blink() expected %v @ %v, got %v", cs.input, e, idx, a)
			}
		}
	}
}

func TestStonesBlinkBlinkN(t *testing.T) {
	stones := Stones{Stone(0), Stone(1)}
	for n := range 25 {
		expected := blinkN(stones, n)
		actual := stones.BlinkN(n)

		if len(expected) != actual {
			t.Fatalf("BlinkN does not match Blink after %v blinks, expected %v, got %v", n, len(expected), actual)
		}
	}
}

func blinkN(stones Stones, n int) Stones {
	if n == 0 {
		return stones
	}

	return blinkN(stones.Blink(), n - 1)
}