package day2

import (
	"fmt"
	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
)

type solution struct{}

type Report []int

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "2"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	reports, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0
	for _, report := range reports {
		if report.IsSafe() {
			count += 1
		}
	}

	return solver.Solved(count)
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	reports, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0
	for _, report := range reports {
		if report.IsSafe() {
			count += 1
			continue
		}

		safe_versions := 0
		for idx, _ := range report {
			dampened := Report{}
			if len(report[:idx]) > 0 {
				dampened = append(dampened, report[:idx]...)
			}

			if len(report[idx+1:]) > 0 {
				dampened = append(dampened, report[idx+1:]...)
			}

			if dampened.IsSafe() {
				safe_versions += 1
				break
			}
		}

		if safe_versions > 0 {
			count += 1
		}
	}
	return solver.Solved(count)
}

func parseInput(input []string) ([]Report, error) {
	reports := []Report{}

	for _, line := range input {
		nums := util.ExtractInts(line)
		if len(nums) == 0 {
			continue
		}

		reports = append(reports, Report(nums))
	}

	if len(reports) == 0 {
		return nil, fmt.Errorf("No reports found")
	}

	return reports, nil
}

func (r Report) IsSafe() bool {
	direction := 0
	safe := true
	util.PairWiseDo(r, func (a, b int) {
		if safe {
			safe, direction = isSafePair(a, b, direction)
		}
	})

	return safe
}

func isSafePair(a, b int, direction int) (bool, int) {
	diff := util.Abs(b-a)
	if diff < 1 || diff > 3 {
		return false, direction
	}

	dir := util.Sign(b-a)
	if direction == 0 {
		return true, dir
	}

	if direction != dir {
		return false, dir
	}

	return true, direction
}
