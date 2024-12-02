package day2

import (
	"fmt"
	"regexp"
	"strconv"
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

func (s solution) Part1(input []string) (string, error) {
	reports, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}
	//fmt.Printf("%v\n", reports)

	count := 0
	for _, report := range reports {
		if report.IsSafe() {
			count += 1
		}
	}

	return solver.Solved(count)
}

func (s solution) Part2(input []string) (string, error) {
	reports, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	count := 0
	for _, report := range reports {
		//fmt.Printf("CHECKING %v FOR DAMPENING\n", report)
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
			//	fmt.Printf("  %v 💹\n", dampened)
				safe_versions += 1
				break
			//} else {
			//	fmt.Printf("  %v ⚠\n", dampened)
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

	num := regexp.MustCompile("[0-9]+")

	for _, line := range input {
		nums := num.FindAllString(line, -1)
		if len(nums) == 0 {
			continue
		}
		report := Report{}
		for _, cand := range nums {
			n, _ := strconv.Atoi(cand)
			report = append(report, n)
		}
		reports = append(reports, report)
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
		//fmt.Printf("%v,%v has wrong difference\n", a, b)
		return false, direction
	}

	dir := util.Sign(b-a)
	if direction == 0 {
		return true, dir
	}

	if direction != dir {
		//fmt.Printf("%v,%v has wrong direction\n", a, b)
		return false, dir
	}

	return true, direction
}
