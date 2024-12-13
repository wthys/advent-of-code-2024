package day12

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
	return "12"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	garden, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	regs := regions(garden)

	total := 0
	for _, region := range regs {
		count := map[L.Location]int{}
		edge := edge(region)
		edge.ForEach(func (loc L.Location) {
			for _, neejber := range loc.OrthoNeejbers() {
				count[neejber] += 1
			}
		})

		A := region.Len()
		S := 0
		region.ForEach(func (loc L.Location) {
			n, _ := count[loc]
			S += n
		})
		cost := A * S
		total += cost
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	garden, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	regs := regions(garden)

	total := 0
	for _, region := range regs {
		A := region.Len()
		S := 0
		for _, dir := range L.New(0,0).OrthoNeejbers() {
			S += countDirectedSides(region, dir)
		}

		cost := A * S

		opts.IfDebugDo(func(_ solver.Options) {
			b, _ := garden.Bounds()
			E := edge(region)
			E.ForEach(func (loc L.Location) {
				b = b.Accomodate(loc)
			})
			garden.PrintBoundsFuncWithLoc(b, func (loc L.Location, v rune, _ error) string {
				if region.Has(loc) {
					return string(v)
				}

				if E.Has(loc) {
					return "+"
				}

				return "."
			})
			fmt.Printf("== A=%v, S=%v, COST=%v ==\n", A, S, cost)
		})
		total += cost
	}

	return solver.Solved(total)
}

type Regions []*S.Set[L.Location]

func parseInput(input []string) (*G.Grid[rune], error) {
	garden := G.WithDefault('.')

	for y, line := range input {
		for x, plot := range line {
			loc := L.New(x, y)
			garden.Set(loc, plot)
		}
	}

	if garden.Len() == 0 {
		return nil, fmt.Errorf("no garden found")
	}

	return garden, nil
}

func regions(g *G.Grid[rune]) Regions {
	regions := Regions{}
	findRegion := func (loc L.Location) *S.Set[L.Location] {
		for _, reg := range regions {
			if reg.Has(loc) {
				return reg
			}
		}
		return nil
	}

	g.ForEach(func (loc L.Location, v rune) {
		if v == '.' {
			return
		}

		if findRegion(loc) != nil {
			return
		}

		region := S.New(loc)
		for _, neejber := range loc.OrthoNeejbers() {
			p, _ := g.Get(neejber)
			if p != v {
				continue
			}
			reg := findRegion(neejber)
			if reg != nil {
				region.AddAll(reg.Values())
			}
		}

		newRegions := Regions{}
		for _, reg := range regions {
			if reg.Subtract(region).Len() > 0 {
				newRegions = append(newRegions, reg)
			}
		}

		regions = append(newRegions, region)
	})
	return regions
}

func edge(region *S.Set[L.Location]) *S.Set[L.Location] {
	regionEdge := S.New[L.Location]()
	region.ForEach(func (loc L.Location) {
		regionEdge.AddAll(loc.OrthoNeejbers())
	})

	return regionEdge.Subtract(region)
}

func countDirectedSides(region *S.Set[L.Location], dir L.Location) int {
	grid := G.WithDefault('.')
	E := edge(region)
	
	bounds, _ := grid.Bounds()
	E.ForEach(func (loc L.Location) {
		bounds = bounds.Accomodate(loc)
	})

	detect := func (prev int, loc L.Location) (int, bool) {
		if region.Has(loc) {
			return 2, prev == 1
		}

		if E.Has(loc) {
			return 1, false
		}

		return 0, false
	}

	switch dir.Unit() {
	case L.New(1,0): //LR
		for y := bounds.Ymin; y <= bounds.Ymax; y++ {
			prev := 0
			side := false
			for x := bounds.Xmin; x <= bounds.Xmax; x++ {
				loc := L.New(x, y)
				prev, side = detect(prev, loc)
				if side {
					grid.Set(loc, '1')
				}
			}
		}
	case L.New(-1,0): //RL
		for y := bounds.Ymin; y <= bounds.Ymax; y++ {
			prev := 0
			side := false
			for x := bounds.Xmax; x >= bounds.Xmin; x-- {
				loc := L.New(x, y)
				prev, side = detect(prev, loc)
				if side {
					grid.Set(loc, '1')
				}
			}
		}
	case L.New(0,1): //TB
		for x := bounds.Xmin; x <= bounds.Xmax; x++ {
			prev := 0
			side := false
			for y := bounds.Ymin; y <= bounds.Ymax; y++ {
				loc := L.New(x, y)
				prev, side = detect(prev, loc)
				if side {
					grid.Set(loc, '1')
				}
			}
		}
	case L.New(0,-1): //BT
		for x := bounds.Xmin; x <= bounds.Xmax; x++ {
			prev := 0
			side := false
			for y := bounds.Ymax; y >= bounds.Ymin; y-- {
				loc := L.New(x, y)
				prev, side = detect(prev, loc)
				if side {
					grid.Set(loc, '1')
				}
			}
		}
	}

	return len(regions(grid))
}