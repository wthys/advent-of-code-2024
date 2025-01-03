package day22

import (
	"testing"
)

type caseNextSecret struct {
	input int
	expected int
}

func TestNextSecret(t *testing.T) {
	cases := []caseNextSecret{
		{123, 15887950},
		{15887950, 16495136},
		{16495136, 527345},
		{527345, 704524},
		{704524,1553684},
		{1553684,12683156},
		{12683156,11100544},
		{11100544,12249484},
		{12249484,7753432},
		{7753432,5908254},
	}

	for _, cs := range cases {
		actual := NextSecret(cs.input)
		if cs.expected != actual {
			t.Errorf("NextSecret(%v) should be %v, got %v", cs.input, cs.expected, actual)
		}
	}
}