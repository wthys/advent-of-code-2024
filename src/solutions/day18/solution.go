package day18

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	PF "github.com/wthys/advent-of-code-2024/pathfinding"
	L "github.com/wthys/advent-of-code-2024/location"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "18"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	locations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	max := 0
	for _, loc := range locations {
		if max < loc.X {
			max = loc.X
		}
		if max < loc.Y {
			max = loc.Y
		}
	}

	start := L.New(0, 0)
	end := L.New(max, max)

	upper := 12
	if len(locations) > 1024 {
		upper = 1024
	}

	excluded := S.New(locations[:upper]...)

	neejberFn := func(loc L.Location) []L.Location {
		neejbers := []L.Location{}
		for _, neejber := range loc.OrthoNeejbers() {
			if excluded.Has(neejber) {
				continue
			}

			if loc.X >= 0 && loc.X <= max && loc.Y >= 0 && loc.Y <= max {
				neejbers = append(neejbers, neejber)
			}
		}
		return neejbers
	}

	pf := PF.ConstructDijkstra(start, neejberFn)

	shortest := pf.ShortestPathLengthTo(end)

	return solver.Solved(shortest)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	locations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	max := 0
	for _, loc := range locations {
		if max < loc.X {
			max = loc.X
		}
		if max < loc.Y {
			max = loc.Y
		}
	}

	start := L.New(0, 0)
	end := L.New(max, max)

	upper := 12
	if len(locations) > 1024 {
		upper = 1024
	}

	lo := upper
	hi := len(locations)-1

	for lo+1 < hi {
		mid := (lo + hi) / 2
		excluded := S.New(locations[:mid]...)
		neejberFn := func(loc L.Location) []L.Location {
			neejbers := []L.Location{}
			for _, neejber := range loc.OrthoNeejbers() {
				if excluded.Has(neejber) {
					continue
				}

				if loc.X >= 0 && loc.X <= max && loc.Y >= 0 && loc.Y <= max {
					neejbers = append(neejbers, neejber)
				}
			}
			return neejbers
		}

		pf := PF.ConstructDijkstra(start, neejberFn)

		shortest := pf.ShortestPathLengthTo(end)
		if shortest == PF.INFINITE {
			hi = mid
		} else {
			lo = mid
		}
	}

	for limit := lo; limit <= hi; limit++ {
		excluded := S.New(locations[:limit]...)
		neejberFn := func(loc L.Location) []L.Location {
			neejbers := []L.Location{}
			for _, neejber := range loc.OrthoNeejbers() {
				if excluded.Has(neejber) {
					continue
				}

				if loc.X >= 0 && loc.X <= max && loc.Y >= 0 && loc.Y <= max {
					neejbers = append(neejbers, neejber)
				}
			}
			return neejbers
		}

		pf := PF.ConstructDijkstra(start, neejberFn)
		shortest := pf.ShortestPathLengthTo(end)
		if shortest == PF.INFINITE {
			return solver.Solved(locations[limit-1])
		}
	}

	return solver.NotImplemented()
}

func parseInput(input []string) (L.Locations, error) {
	locations := L.Locations{}
	for lineno, line := range input {
		vals := util.ExtractInts(line)
		if len(vals) == 0 {
			continue
		}

		if len(vals) < 2 {
			return nil, fmt.Errorf("#%v : not enough values, %q", lineno, line)
		}
		locations = append(locations, L.New(vals[0], vals[1]))
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("no bytes found")
	}

	return locations, nil
}