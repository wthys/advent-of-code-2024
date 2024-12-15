package day14

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	L "github.com/wthys/advent-of-code-2024/location"
	G "github.com/wthys/advent-of-code-2024/grid"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "14"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	robots, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	halfX := MAP_WIDTH/2
	halfY := MAP_HEIGHT/2

	quadrants := []G.Bounds{
		G.Bounds{0,         halfX - 1, 0,         halfY - 1 },
		G.Bounds{0,         halfX - 1, halfY + 1, MAP_HEIGHT},
		G.Bounds{halfX + 1, MAP_WIDTH, 0,         halfY - 1 },
		G.Bounds{halfX + 1, MAP_WIDTH, halfY + 1, MAP_HEIGHT},
	}

	counts := map[int]int{}

	moved := robots.MoveN(100)
	for _, robot := range moved {
		for idx, quadrant := range quadrants {
			if quadrant.Has(robot.Pos) {
				counts[idx] += 1
				break
			}
		}
	}

	securityFactor := 1
	for _, cnt := range counts {
		securityFactor *= cnt
	}

	return solver.Solved(securityFactor)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	robots, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	PATTERN_SIZE := 9
	firstCheck := L.Locations{}
	for n := range 2 {
		for m := range 2 {
			firstCheck = append(firstCheck, L.New(n, m).Scale(PATTERN_SIZE-1))
		}
	}

	secondCheck := L.Locations{}
	for n := range PATTERN_SIZE {
		for m := range PATTERN_SIZE {
			secondCheck = append(secondCheck, L.New(n, m))
		}
	}

	searchPatterns := []*S.Set[L.Location]{}
	for x := range MAP_WIDTH - PATTERN_SIZE {
		for y := range MAP_HEIGHT - PATTERN_SIZE {
			root := L.New(x, y)
			searchPatterns = append(searchPatterns, createSearchPattern(root, firstCheck))
		}
	}

	waitTime := 0
	found := false
	for !found {
		opts.Debugf("__ checking %v __\n", waitTime)
		moved := robots.MoveN(waitTime)

		locs := S.New(moved.Positions()...)

		for _, search := range searchPatterns {
			if search.Subtract(locs).Len() == 0 {
				b := G.BoundsFromSlice(search.Values())

				doubleCheck := createSearchPattern(L.New(b.Xmin, b.Ymin), secondCheck)

				if doubleCheck.Subtract(locs).Len() == 0 {
					found = true
					break
				}
			}
		}

		if found {
			opts.IfDebugDo(func (_ solver.Options) {
				visualizeRobots(moved)
			})
			break
		}
		waitTime++
	}

	return solver.Solved(waitTime)
}


type (
	Robots []Robot
	Robot struct {
		Pos L.Location
		Dir L.Location
	}
)

const (
	MAP_WIDTH = 101
	MAP_HEIGHT = 103
)

func (r Robot) MoveN(n int) Robot {
	newPos := r.Pos.Add(r.Dir.Scale(n))
	for newPos.X < 0 {
		newPos.X += MAP_WIDTH
	}
	newPos.X = newPos.X % MAP_WIDTH

	for newPos.Y < 0 {
		newPos.Y += MAP_HEIGHT
	}
	newPos.Y = newPos.Y % MAP_HEIGHT
	return Robot{newPos, r.Dir}
}

func (r Robot) Move() Robot {
	return r.MoveN(1)
}

func (robots Robots) MoveN(n int) Robots {
	rs := Robots{}
	for _, r := range robots {
		rs = append(rs, r.MoveN(n))
	}
	return rs
}

func (robots Robots) Move() Robots {
	return robots.MoveN(1)
}

func (robots Robots) Positions() L.Locations {
	locs := L.Locations{}
	for _, robot := range robots {
		locs = append(locs, robot.Pos)
	}
	return locs
}

func visualizeRobots(robots Robots) {
	robotLocs := map[L.Location]int{}
	for _, r  := range robots {
		robotLocs[r.Pos] += 1
	}
	g := G.WithDefault(0)
	for loc, n := range robotLocs {
		g.Set(loc, n)
	}

	b := G.Bounds{0, MAP_WIDTH, 0, MAP_HEIGHT}
	g.PrintBoundsFuncWithLoc(b, func(_ L.Location, v int, _ error) string {
		if v == 0 {
			return "⬛"
		}
		return "🟩"
	})
}

func parseInput(input []string) (Robots, error) {
	robots := Robots{}

	for lno, line := range input {
		values := util.ExtractInts(line)
		if len(values) == 0 {
			continue
		}
		if len(values) < 4 {
			return nil, fmt.Errorf("#%v : invalid robot %q", lno, line)
		}

		pos := L.New(values[0], values[1])
		dir := L.New(values[2], values[3])
		robots = append(robots, Robot{pos, dir})
	}
	
	if len(robots) == 0 {
		return nil, fmt.Errorf("no robots found")
	}

	return robots, nil
}

func createSearchPattern(loc L.Location, pattern L.Locations) *S.Set[L.Location] {
	search := S.New[L.Location]()
	for _, p := range pattern {
		search.Add(loc.Add(p))
	}
	return search
}