package day20

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	L "github.com/wthys/advent-of-code-2024/location"
	S "github.com/wthys/advent-of-code-2024/collections/set"
	PF "github.com/wthys/advent-of-code-2024/pathfinding"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "20"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	start, end, locations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	step0 := Step{start}
	stepN := Step{end}
	neejberFn := func(step Step) []Step {
		return Neejbers(step, locations)
	}

	pf := PF.ConstructDijkstra(step0, neejberFn)

	baselinePath := append([]Step{step0}, pf.ShortestPathTo(stepN)...)
	baseline := pathLength(baselinePath)
	opts.Debugf("__ baseline = %v ps\n", baseline)

	pathLengths := map[int]int{}

	for idx, step := range baselinePath[:baseline-1] {
		for _, cheat := range CheatNeejbers(step, locations) {
			bestLength := idx
			for i, s := range baselinePath {
				if s.Pos == cheat.Pos {
					if i > bestLength {
						bestLength = i
					}
				}
			}
			if bestLength == idx {
				continue
			}

			cpath := append([]Step{}, baselinePath[:idx+1]...)
			cpath = append(cpath, baselinePath[bestLength:]...)
			length := pathLength(cpath)
			pathLengths[length] += 1
		}
	}

	count := 0
	for length, n := range pathLengths {
		if baseline - length >= 100 {
			opts.Debugf("____ %v cheats save %v ps (%v ps)\n", n, baseline - length, length)
			count += n
		}
	}

	return solver.Solved(count)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	start, end, locations, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	step0 := Step{start}
	stepN := Step{end}
	neejberFn := func(step Step) []Step {
		return Neejbers(step, locations)
	}

	pf := PF.ConstructDijkstra(step0, neejberFn)

	baselinePath := append([]Step{step0}, pf.ShortestPathTo(stepN)...)
	baseline := pathLength(baselinePath)
	opts.Debugf("__ baseline = %v ps\n", baseline)

	pathLengths := map[int]int{}

	for idx, step := range baselinePath[:baseline] {
		opts.Debugf("__ checking step %v : %v (%v%%)\n", idx, step, 100 * idx / (baseline + 1))
		for cidx := baseline; cidx > idx; cidx-- {
			cheat := baselinePath[cidx]
			if cheat.Pos.Subtract(step.Pos).Manhattan() <= 20 {
				cpath := append([]Step{}, baselinePath[:idx+1]...)
				cpath = append(cpath, baselinePath[cidx:]...)
				length := pathLength(cpath)
				pathLengths[length] += 1
			}
		}
	}

	count := 0
	for length, n := range pathLengths {
		if baseline - length >= 100 {
			opts.Debugf("____ %v cheats save %v ps (%v ps)\n", n, baseline - length, length)
			count += n
		}
	}

	return solver.Solved(count)
}

type Step struct {
	Pos L.Location
}

func pathLength(steps []Step) int {
	total := 0
	util.PairWiseDo(steps, func(a, b Step) {
		total += b.Pos.Subtract(a.Pos).Manhattan()
	})
	return total
}

func (s Step) String() string {
	return fmt.Sprintf("%v", s.Pos)
}

func (s Step) Move(dir L.Location) Step {
	return Step{s.Pos.Add(dir)}
}

var cheatNeejbers []L.Location = []L.Location{
	L.New( 2,  0),
	L.New( 1,  1),
	L.New( 0,  2),
	L.New(-1,  1),
	L.New(-2,  0),
	L.New(-1, -1),
	L.New( 0, -2),
	L.New( 1, -1),
}

func CheatNeejbers(step Step, validPositions *S.Set[L.Location]) []Step {
	neejbers := []Step{}

	for _, dir := range cheatNeejbers {
		neejber := step.Move(dir)
		if validPositions.Has(neejber.Pos) {
			neejbers = append(neejbers, neejber)
		}
	}

	return neejbers
}

func Neejbers(step Step, validPositions *S.Set[L.Location]) []Step {
	neejbers := []Step{}
	
	for _, dir := range L.New(0,0).OrthoNeejbers() {
		neejber := step.Move(dir)
		if validPositions.Has(neejber.Pos) {
			neejbers = append(neejbers, neejber)
		}
	}

	return neejbers
}

func parseInput(input []string) (L.Location, L.Location, *S.Set[L.Location], error) {
	start := L.New(-1,-1)
	end := L.New(-1,-1)
	path := S.NewFor(L.New(0,0))

	for y, line := range input {
		for x, f := range line {
			pos := L.New(x, y)
			switch string(f) {
			case ".":
				path.Add(pos)
			case "S":
				start = pos
			case "E":
				end = pos
			}
		}
	}

	if start.X < 0 {
		return start, start, nil, fmt.Errorf("no start found")
	}
	
	if end.X < 0 {
		return end, end, nil, fmt.Errorf("no end found")
	}

	if path.Len() == 0 {
		return start, end, nil, fmt.Errorf("no path found")
	}

	return start, end, path.Union(S.New(start, end)), nil
}
