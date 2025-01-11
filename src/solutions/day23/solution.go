package day23

import (
	"fmt"
	"strings"
	"slices"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "23"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	links, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	trios := S.NewFor("")
	for _, host := range links.Hosts() {
		if !strings.HasPrefix(host, "t") {
			continue
		}

		neejbers := links.Neejbers(host)
		util.CombinationNoRepeatDo(neejbers, 2, func(couple []string) {
			a, b := couple[0], couple[1]
			if links.Has(a, b) {
				trios.Add(uniform(host, a, b))
			}
		})
	}

	opts.IfDebugDo(func(_ solver.Options) {
		trios.ForEach(func (s string) {
			opts.Debugf("__ %v\n", s)
		})
	})

	return solver.Solved(trios.Len())
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	links, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	largest := ""
	checked := S.NewFor("")
	for _, host := range links.Hosts() {
		connected := append(links.Neejbers(host), host)
		if len(connected) < len(largest)/3+1 {
			continue
		}
		
		key := uniform(connected...)
		if checked.Has(key) { continue }
		slices.Sort(connected)

		// for n := util.Max(len(largest)/3+2, 2); n <= len(connected); n++ {
		for n := len(connected); n >= util.Max(len(largest)/3+2, 2); n-- {
			// opts.Debugf("__ checking %v / %v\n", connected, n)
			found := false
			
			util.CombinationNoRepeatDo(connected, n, func(neejbers []string) {
				if found { return }

				lan := uniform(neejbers...)
				if checked.Has(lan) { return }
				checked.Add(lan)

				if n != S.New(neejbers...).Len() { return }

				fault := false

				util.CombinationNoRepeatDo(neejbers, 2, func(couple []string) {
					if couple[0] == couple[1] { return }
					if fault { return }

					fault = !links.Has(couple[0], couple[1])
				})

				if !fault {
					opts.Debugf(" > checking %v\n", lan)
					if len(lan) > len(largest) {
						largest = lan
						opts.Debugf("___ largest : %v\n", largest)
						found = true
					}
				}
			})
		}
		checked.Add(key)
	}

	return solver.Solved(largest)
}

type (
	Links struct {
		links map[string]*S.Set[string]
	}
)

func uniform(hosts ...string) string {
	vals := append([]string{}, hosts...)
	slices.Sort(vals)
	return strings.Join(vals, ",")
}

func (l Links) AddLink(from, to string) {
	dests, ok := l.links[from]
	if !ok {
		l.links[from] = S.New(to)
	} else {
		dests.Add(to)
	}
}

func (l Links) Len() int {
	return len(l.links)
}

func (l Links) Has(a, b string) bool {
	dest, ok := l.links[a]
	if !ok {
		return false
	}
	if dest.Has(b) {
		return true
	}
	dest, ok = l.links[b]
	if !ok {
		return false
	}
	return dest.Has(a)
}

func (l Links) Hosts() []string {
	hosts := []string{}
	for host, _ := range l.links {
		hosts = append(hosts, host)
	}
	return hosts
}

func (l Links) Neejbers(a string) []string {
	dests, ok := l.links[a]
	if !ok {
		return []string{}
	}
	return dests.Values()
}

func parseInput(input []string) (*Links, error) {
	links := &Links{map[string]*S.Set[string]{}}

	for lineno, line := range input {
		ends := util.ExtractRegex("[^-]+", line)
		if len(ends) == 0 {
			continue
		}

		if len(ends) < 2 {
			return nil, fmt.Errorf("#%v : not enough ends : %q", lineno, line)
		}

		links.AddLink(ends[0], ends[1])
		links.AddLink(ends[1], ends[0])
	}

	if links.Len() == 0 {
		return nil, fmt.Errorf("no network links found")
	}

	return links, nil
}