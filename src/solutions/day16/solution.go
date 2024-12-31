package day16

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	PF "github.com/wthys/advent-of-code-2024/pathfinding"
	L "github.com/wthys/advent-of-code-2024/location"
	G "github.com/wthys/advent-of-code-2024/grid"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "16"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	start, end, walkable, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	step0 := Step{start, L.New(1,0)}

	pf := initDijkstra(step0, walkable)

	minLength := PF.INFINITE
	for _, dir := range L.New(0,0).OrthoNeejbers() {
		stepN := Step{end, dir}
		path := pf.ShortestPathTo(stepN)
		if len(path) == 0 {
			continue
		}
		path = append(Steps{step0}, path...)
		length := pf.ShortestPathLengthTo(stepN)
		if length < minLength {
			minLength = length
		}
		// opts.IfDebugDo(func (_ solver.Options) {
		// 	opts.Debugf("=== %v points ===\n", length)
		// 	visualisePath(step0, stepN, Steps(path), walkable)
		// })
	}

	return solver.Solved(minLength)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	start, end, walkable, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	step0 := Step{start, L.New(1,0)}

	pf := initDijkstra(step0, walkable)

	minLength := PF.INFINITE
	for _, dir := range L.New(0,0).OrthoNeejbers() {
		stepN := Step{end, dir}
		path := pf.ShortestPathTo(stepN)
		if len(path) == 0 {
			continue
		}
		path = append(Steps{step0}, path...)
		length := pf.ShortestPathLengthTo(stepN)
		if length < minLength {
			minLength = length
		}
	}

	spots := S.NewFor(step0.Pos)

	for _, dir := range L.New(0,0).OrthoNeejbers() {
		stepN := Step{end, dir}
		pf.ShortestPathToFunc(stepN, func(path []Step) {
			if Steps(path).Cost() != minLength {
				return
			}

			for _, step := range path {
				spots.Add(step.Pos)
			}
		})
	}

	opts.IfDebugDo(func(_ solver.Options) {
		visualiseSpots(spots, walkable)
	})

	return solver.Solved(spots.Len())
}

type (
	Step struct {
		Pos L.Location
		Dir L.Location
	}
	Steps []Step
)

var (
	DIRMAP = map[L.Location]string{
		L.New(1,0): ">",
		L.New(-1,0): "<",
		L.New(0,1): "v",
		L.New(0,-1): "^",
	}
)

func initDijkstra(step0 Step, walkable *S.Set[L.Location]) PF.Dijkstra[Step] {
	neejberFn := func (from Step) []Step {
		nbrs := []Step{}
		for _, dir := range L.New(0,0).OrthoNeejbers() {
			to := from.Move(dir)
			if walkable.Has(to.Pos) {
				nbrs = append(nbrs, to)
			}
		}
		return nbrs
	}
	weightFn := func (a, b Step) int {
		cost, err := a.MoveCost(b)
		if err != nil {
			return PF.INFINITE
		}
		return cost
	}

	return PF.ConstructWeightedDijkstra(step0, neejberFn, weightFn)
}

func (from Step) MoveCost(to Step) (int, error) {
	dist := to.Pos.Subtract(from.Pos)

	if dist.X != 0 && dist.Y != 0 {
		return 0, fmt.Errorf("not a straight line (%v -> %v)", from, to)
	}

	if from.Dir == to.Dir {
		return dist.Manhattan(), nil
	}

	if from.Dir.Add(to.Dir) == L.New(0,0) {
		return 2000 + dist.Manhattan(), nil
	}

	return 1000 + dist.Manhattan(), nil
}

func (step Step) String() string {
	return fmt.Sprintf("%v=%v", step.Pos, DIRMAP[step.Dir])
}

func moveCost(from Step, to Step) int {
	cost, err := from.MoveCost(to)
	if err != nil {
		panic(err)
	}
	return cost
}

func (from Step) Move(dir L.Location) Step {
	udir := dir.Unit()
	return Step{from.Pos.Add(udir), udir}
}

func (steps Steps) Cost() int {
	total := 0
	prev := steps[0]
	for _, step := range steps[1:] {
		total += moveCost(prev, step)
		prev = step
	}
	return total
}

func parseInput(input []string) (L.Location, L.Location, *S.Set[L.Location], error) {
	start := L.New(0,0)
	end := start
	walkable := S.NewFor(start)

	lnil := L.New(-1,-1)

	for y, line := range input {
		for x, space := range line {
			pos := L.New(x, y)
			switch space {
			case '.':
				walkable.Add(pos)
			case 'S':
				start = pos
			case 'E':
				end = pos
			}
		}
	}

	if start.Manhattan() == 0 {
		return lnil, lnil, nil, fmt.Errorf("no start found")
	}

	if end.Manhattan() == 0 {
		return lnil, lnil, nil, fmt.Errorf("no end found")
	}

	if walkable.Len() == 0 {
		return lnil, lnil, nil, fmt.Errorf("no walkable spaces found")
	}

	walkable.Add(start).Add(end)

	return start, end, walkable, nil
}

func visualisePath(start Step, end Step, steps Steps, walkable *S.Set[L.Location]) {
	g := G.WithDefault("#")
	walkable.ForEach(func (loc L.Location) {
		g.Set(loc, ".")
	})

	for _, step := range steps {
		g.Set(step.Pos, DIRMAP[step.Dir])
	}

	g.Set(start.Pos, "S")
	g.Set(end.Pos, "E")

	mapper := map[string]string {
		"#": "■",
		".": " ",
		">": "→",
		"<": "←",
		"^": "↑",
		"v": "↓",
		"S": "S",
		"E": "E",
	}

	g.PrintFunc(func(v string, _ error) string {
		return mapper[v]
	})
}

func visualiseSpots(spots *S.Set[L.Location], walkable *S.Set[L.Location]) {
	g := G.WithDefault("#")
	walkable.ForEach(func (loc L.Location) {
		g.Set(loc, ".")
	})

	spots.ForEach(func (loc L.Location) {
		g.Set(loc, "O")
	})

	g.Print()
}