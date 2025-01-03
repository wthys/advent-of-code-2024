package day19

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
	return "19"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	available, desired, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0
	counter := NewCounterOpts(available, &opts)
	for _, display := range desired {
		if counter.Count(display) > 0 {
			opts.Debugf("__ %q is possible\n", display)
			count++
		}
	}

	return solver.Solved(count)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	available, desired, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	counter := NewCounterOpts(available, &opts)

	total := 0
	for _, display := range desired {
		inter := counter.Count(display)
		opts.Debugf("__ %q = %v\n", display, inter)
		total += inter
	}

	return solver.Solved(total)
}

func parseInput(input []string) ([]string, []string, error) {
	available := util.ExtractRegex("[wubrg]+", input[0])

	if len(available) == 0 {
		return nil, nil, fmt.Errorf("No available patterns")
	}

	desired := []string{}
	for _, line := range input[2:] {
		if len(line) == 0 {
			continue
		}

		desired = append(desired, line)
	}

	if len(desired) == 0 {
		return nil, nil, fmt.Errorf("No desired patterns")
	}

	return available, desired, nil
}

type Counter struct {
	patterns []string
	cache map[string]int
	opts solver.Options
}

func NewCounter(patterns []string) Counter {
	return NewCounterOpts(patterns, nil)
}

func NewCounterOpts(patterns []string, opts *solver.Options) Counter {
	cache := map[string]int{"": 1}
	if opts == nil {
		return Counter{patterns, cache, solver.DefaultOptions()}
	}
	return Counter{patterns, cache, *opts}
}

func (c Counter) Count(target string) int {
	c.opts.Debugf("Counter(%q)\n", target)
	cached, ok := c.cache[target]

	if ok {
		return cached
	}

	total := 0
	for _, p := range c.patterns {
		if strings.HasPrefix(target, p) {
			c.opts.Debugf("removing %q => ", p)
			total += c.Count(target[len(p):])
		}
	}

	c.cache[target] = total
	return total
}