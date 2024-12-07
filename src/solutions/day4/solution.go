package day4

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	L "github.com/wthys/advent-of-code-2024/location"
	G "github.com/wthys/advent-of-code-2024/grid"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "4"
}

func (s solution) Part1(input []string) (string, error) {
	grid, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0
	directions := L.Locations{ {0,1}, {1,1}, {1,0}, {1,-1}, {-1,0}, {-1,-1}, {0,-1}, {-1,1} }
	grid.ForEach(func (loc L.Location, letter rune) {
		if letter != 'X' {
			return
		}

		for _, dir := range directions {
			if matchXmas(grid, loc, dir) {
				count += 1
			}
		}
	})

	return solver.Solved(count)
}

func (s solution) Part2(input []string) (string, error) {
	grid, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0

	forwardDirs := L.Locations{{1,1},{-1,-1}}
	backDirs := L.Locations{{-1,1},{1,-1}}
	grid.ForEach(func (start L.Location, letter rune) {
		if letter != 'A' {
			return
		}

		for _, fdir := range forwardDirs {
			if !matchMas(grid, start.Add(fdir.Scale(-1)), fdir) {
				continue
			}
			for _, bdir := range backDirs {
				if matchMas(grid, start.Add(bdir.Scale(-1)), bdir) {
					count += 1
				}
			}
		}
	})

	return solver.Solved(count)
}

func parseInput(input []string) (*G.Grid[rune], error) {
	g := G.WithDefault(rune('.'))

	for y, line := range input {
		for x, letter := range line {
			if !(letter == 'X' || letter == 'M' || letter == 'A' || letter == 'S') {
				return nil, fmt.Errorf("letter '%v' on line %v is not valid", letter, y+1)
			}
			g.Set(L.New(x, y), letter)
		}
	}

	return g, nil
}

func matchXmas(grid *G.Grid[rune], start, direction L.Location) bool {
	return matchWord(grid, "XMAS", start, direction)
}

func matchMas(grid *G.Grid[rune], start, direction L.Location) bool {
	return matchWord(grid, "MAS", start, direction)
}

func matchWord(grid *G.Grid[rune], word string, start, direction L.Location) bool {
	for n, letter := range word {
		l, ok := grid.Get(start.Add(direction.Scale(n)))
		if ok != nil || l != letter {
			return false
		}
	}

	return true
}