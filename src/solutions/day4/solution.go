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
	directions := []L.Location{ {0,1}, {1,1}, {1,0}, {1,-1}, {-1,0}, {-1,-1}, {0,-1}, {-1,1} }
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
	return solver.NotImplemented()
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
	l := '.'
	ok := fmt.Errorf("nothing really")

	l, ok = grid.Get(start)
	if ok != nil || l != 'X' {
		return false
	}

	l, ok = grid.Get(start.Add(direction))
	if ok != nil || l != 'M' {
		return false
	}

	l, ok = grid.Get(start.Add(direction.Scale(2)))
	if ok != nil || l != 'A' {
		return false
	}

	l, ok = grid.Get(start.Add(direction.Scale(3)))
	if ok != nil || l != 'S' {
		return false
	}

	return true
}