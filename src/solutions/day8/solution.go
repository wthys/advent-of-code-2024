package day8

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	G "github.com/wthys/advent-of-code-2024/grid"
	L "github.com/wthys/advent-of-code-2024/location"
	S "github.com/wthys/advent-of-code-2024/collections/set"

)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "8"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	grid, catalog, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	opts.Debugf("catalog : %v\n", catalog)

	antinodes := S.New[L.Location]()
	for antenna, list := range catalog {
		opts.Debugf("__ %v -> %v\n", antenna, list)
		util.CombinationDo(list, 2, func (locs []L.Location) {
			a := locs[0]
			b := locs[1]

			if a == b {
				return
			}

			opts.Debugf("____ checking %v - %v\n", a, b)

			diff := b.Subtract(a)
			cands := L.Locations{
				b.Add(diff),
				a.Subtract(diff),
			}
			for _, cand := range cands {
				_, ok := grid.Get(cand)
				if ok == nil {
					antinodes.Add(cand)
				}
			}
		})
	}

	opts.IfDebugDo(func (_ solver.Options) {
		grid.PrintFuncWithLoc(func (loc L.Location, v string, _ error) string {
			if antinodes.Has(loc) {
				return "#"
			}

			return v
		})
	})

	return solver.Solved(antinodes.Len())
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	grid, catalog, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	opts.Debugf("catalog : %v\n", catalog)

	antinodes := S.New[L.Location]()
	for antenna, list := range catalog {
		opts.Debugf("__ %v -> %v\n", antenna, list)
		util.CombinationDo(list, 2, func (locs []L.Location) {
			a := locs[0]
			b := locs[1]

			if a == b {
				return
			}

			opts.Debugf("____ checking %v - %v\n", a, b)

			diff := b.Subtract(a)

			done := false
			forward := b
			backward := a
			antinodes.Add(a)
			antinodes.Add(b)
			for !done {
				forward = forward.Add(diff)
				backward = backward.Subtract(diff)
				cands := L.Locations{forward, backward}
				added := false
				for _, cand := range cands{
					_, ok := grid.Get(cand)
					if ok == nil {
						antinodes.Add(cand)
						added = true
					}
				}
				if !added {
					done = true
				}
			}
		})
	}

	opts.IfDebugDo(func (_ solver.Options) {
		grid.PrintFuncWithLoc(func (loc L.Location, v string, _ error) string {
			if antinodes.Has(loc) {
				return "#"
			}

			return v
		})
	})

	return solver.Solved(antinodes.Len())
}

type AntennaCatalog map[string]L.Locations

func parseInput(input []string) (*G.Grid[string], AntennaCatalog, error) {
	grid := G.New[string]()
	tempCatalog := map[rune]*S.Set[L.Location]{}

	for y, line := range input {
		for x, v := range line {
			loc := L.New(x, y)

			grid.Set(loc, string(v))
			if v != '.' {
				lst, ok := tempCatalog[v]
				if !ok {
					lst = S.New(loc)
					tempCatalog[v] = lst
				}
				lst.Add(loc)
			}
		}
	}

	if grid.Len() == 0 {
		return nil, nil, fmt.Errorf("no map to be found")
	}

	catalog := AntennaCatalog{}
	for v, locs := range tempCatalog {
		catalog[string(v)] = L.Locations(locs.Values())
	}

	return grid, catalog, nil
}
