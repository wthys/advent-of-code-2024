package day6

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	G "github.com/wthys/advent-of-code-2024/grid"
	L "github.com/wthys/advent-of-code-2024/location"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "6"
}

func (s solution) Part1(input []string) (string, error) {
	guard, grid, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	guards, _ := simulateGuard(guard, grid, L.Locations{})

	return solver.Solved(guards.PositionSet().Len())
}

func (s solution) Part2(input []string) (string, error) {
	guard, grid, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	guards, _ := simulateGuard(guard, grid, L.Locations{})

	candidates := S.New[L.Location]()
	for _, g := range guards {
		if (g.pos == guard.pos) {
			continue
		}

		_, looped := simulateGuard(guard, grid, L.Locations{g.pos})
		if looped {
			candidates.Add(g.pos)
		}
	}

	return solver.Solved(candidates.Len())
}

func simulateGuard(guard Guard, grid *G.Grid[string], obstacles L.Locations) (Guards, bool) {
	gg := G.New[string]()
	grid.ForEach(func (loc L.Location, v string) {
		gg.Set(loc, v)
	})

	for _, obs := range obstacles {
		gg.Set(obs, "O")
	}

	visited := S.New(guard)
	newG, ok := guard.Move(gg)
	for ok {
		if visited.Has(newG) {
			return visited.Values(), true
		}
		visited.Add(newG)
		newG, ok = newG.Move(gg)
	}

	return visited.Values(), false
}

type Guard struct {
	pos L.Location
	dir L.Location
}

type Guards []Guard

func rotate(g Guard) Guard {
	rotated := map[L.Location]L.Location {
		L.New(1,0): L.New(0,1),
		L.New(0,1): L.New(-1,0),
		L.New(-1,0): L.New(0,-1),
		L.New(0,-1): L.New(1,0),
	}

	return Guard{g.pos, rotated[g.dir]}
}

func (g Guard) Move(grid *G.Grid[string]) (Guard, bool) {
	newPos := g.pos.Add(g.dir)
	v, err := grid.Get(newPos)
	if err != nil {
		return Guard{}, false
	}
	if (v != ".") {
		return rotate(g), true
	}
	return Guard{newPos, g.dir}, true
}

func (guards Guards) PositionSet() *S.Set[L.Location] {
	return S.New(guards.Positions()...)
}

func (guards Guards) Positions() L.Locations {
	locs := L.Locations{}
	for _, g := range guards {
		locs = append(locs, g.pos)
	}
	return locs
}

func parseInput(input []string) (Guard, *G.Grid[string], error) {
	guard := Guard{}
	guardSet := false
	grid := G.New[string]()

	dirs := map[rune]L.Location{
		'^': L.New(0,-1),
		'>': L.New(1,0),
		'v': L.New(0,1),
		'<': L.New(-1,0),
	}

	for y, line := range input {
		for x, v := range line {
			pos := L.New(x,y)
			grid.Set(pos, string(v))
			
			dir, ok := dirs[v]
			if ok {
				guard = Guard{pos, dir}
				grid.Set(pos, ".")
				guardSet = true
			}
		}
	}

	if grid.Len() == 0 {
		return Guard{}, nil, fmt.Errorf("no map provided")
	}

	if !guardSet {
		return Guard{}, nil, fmt.Errorf("no guard present")
	}

	return guard, grid, nil
}

type Region struct {
	display string
	locations *S.Set[L.Location]
}
type RegionMap []Region

func visualiseRegions(grid *G.Grid[string], regionMap RegionMap) {
	grid.PrintFuncWithLoc(func (loc L.Location, v string, _ error) string {
		for _, region := range regionMap {
			if region.locations.Has(loc) {
				return region.display
			}
		}
		return v
	})
}