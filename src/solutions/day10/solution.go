package day10

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	G "github.com/wthys/advent-of-code-2024/grid"
	L "github.com/wthys/advent-of-code-2024/location"
	S "github.com/wthys/advent-of-code-2024/collections/set"
	PF "github.com/wthys/advent-of-code-2024/pathfinding"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "10"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	grid, trailheads, destinations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	trailheads.ForEach(func (start L.Location) {
		dijkstra := PF.ConstructDijkstra(start, func (loc L.Location) []L.Location {
			neejbers := []L.Location{}
			level, _ := grid.Get(loc)
			for _, neejber := range loc.OrthoNeejbers() {
				lvl, err := grid.Get(neejber)
				if err == nil && lvl == level + 1 {
					neejbers = append(neejbers, neejber)
				}
			}
			return neejbers
		})

		score := 0
		destinations.ForEach(func (end L.Location) {
			path := dijkstra.ShortestPathTo(end)
			if path != nil && len(path) > 0{
				score += 1
			}
		})
		opts.Debugf("%v scores %v\n", start, score)
		total += score
	})

	return solver.Solved(total)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	grid, trailheads, destinations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	trailheads.ForEach(func (start L.Location) {
		bfs := PF.ConstructBreadthFirst(start, func (loc L.Location) []L.Location {
			neejbers := []L.Location{}
			level, _ := grid.Get(loc)
			for _, neejber := range loc.OrthoNeejbers() {
				lvl, err := grid.Get(neejber)
				if err == nil && lvl == level + 1 {
					neejbers = append(neejbers, neejber)
				}
			}
			return neejbers
		})

		rating := 0
		destinations.ForEach(func (end L.Location) {
			bfs.AllPathsFunc(end, func(_ []L.Location) {
				rating += 1
			})
		})
		opts.Debugf("%v rates %v\n", start, rating)
		total += rating
	})

	return solver.Solved(total)
}

func parseInput(input []string) (*G.Grid[int], *S.Set[L.Location], *S.Set[L.Location], error) {
	grid := G.New[int]()
	trailheads := S.New[L.Location]()
	destinations := S.New[L.Location]()

	for y, line := range input {
		levels, err := util.StringsToInts(util.ExtractRegex("[0-9]", line))
		if err != nil {
			return nil, nil, nil, err
		}

		for x, level := range levels {
			loc := L.New(x, y)
			grid.Set(loc, level)
			if level == 0 {
				trailheads.Add(loc)
			} else if level == 9 {
				destinations.Add(loc)
			}
		}
	}

	if grid.Len() == 0 {
		return nil, nil, nil, fmt.Errorf("no map found")
	}

	if trailheads.Len() == 0 {
		return nil, nil, nil, fmt.Errorf("no trailheads found")
	}

	if destinations.Len() == 0 {
		return nil, nil, nil, fmt.Errorf("no destinations found")
	}

	return grid, trailheads, destinations, nil
}