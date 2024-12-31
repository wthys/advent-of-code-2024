package pathfinding

import (
	"fmt"

	"github.com/wthys/advent-of-code-2024/collections/set"
)

type (
	DistMap[T comparable] map[T]int
	PrevMap[T comparable] map[T]*T
	FromMap[T comparable] map[T]*set.Set[T]

	Dijkstra[T comparable] interface {
		ShortestPathTo(end T) []T
		ShortestPathLengthTo(end T) int
		ShortestPathToFunc(end T, complete PathConsumer[T])
		ForEachNode(doer func(node T) bool)
	}

	SimpleDijkstra[T comparable] struct {
		Start T
		dist  DistMap[T]
		prev  PrevMap[T]
		from  FromMap[T]
	}

	BreadthFirst[T comparable] interface {
		AllPathsTo(end T) [][]T
		AllPathsFunc(end T, complete PathConsumer[T])
		AllPathsFuncVisited(end T, complete PathConsumer[T], visited VisitedFunc[T])
	}

	SimpleBFS[T comparable] struct {
		Start T
		neejberFunc NeejberFunc[T]
	}

	PathConsumer[T comparable]   func(path []T)
	NeejberFunc[T comparable]    func(node T) []T
	ExitFunc[T comparable]       func(node T) bool
	EdgeWeightFunc[T comparable] func(in T, out T) int
	VisitedFunc[T comparable]    func(path []T, node T) bool
)

const (
	INFINITE = int((^uint(0)) >> 1)
)

func NeverVisited[T comparable](_ []T, _ T) bool {
	return false
}

func ControlledDijkstra[T comparable](start T, neejbers NeejberFunc[T], exitters ...ExitFunc[T]) Dijkstra[T] {
	dist := DistMap[T]{}
	prev := PrevMap[T]{}
	from := FromMap[T]{}
	visited := set.New[T]()

	prev[start] = nil
	dist[start] = 0
	queue := set.New(start)

	for queue.Len() > 0 {
		node, err := closest(queue, dist)
		if err != nil {
			panic(err)
		}
		queue = queue.Remove(node)
		visited.Add(node)

		stop := false
		for _, exit := range exitters {
			if exit(node) {
				stop = true
			}
		}
		if stop {
			break
		}

		for _, neejber := range neejbers(node) {
			if visited.Has(neejber) {
				continue
			}
			queue.Add(neejber)
			alt := dist[node] + 1
			ndist, ok := dist[neejber]
			if !ok || alt < ndist {
				dist[neejber] = alt
				from[node] = set.New(node)
				prev[neejber] = &node
			} else if alt <= ndist {
				from[neejber].Add(node)
			}
		}
	}

	return SimpleDijkstra[T]{start, dist, prev, from}

}

func ConstructWeightedDijkstra[T comparable](start T, neejbers NeejberFunc[T], weigh EdgeWeightFunc[T]) Dijkstra[T] {
	dist := DistMap[T]{}
	prev := PrevMap[T]{}
	from := FromMap[T]{}
	visited := set.New[T]()

	prev[start] = nil
	dist[start] = 0
	queue := set.New(start)

	for queue.Len() > 0 {
		node, err := closest(queue, dist)
		if err != nil {
			panic(err)
		}
		queue = queue.Remove(node)
		visited.Add(node)

		for _, neejber := range neejbers(node) {
			if visited.Has(neejber) {
				continue
			}
			queue.Add(neejber)
			alt := dist[node] + weigh(node, neejber)
			ndist, ok := dist[neejber]
			if !ok || alt < ndist {
				dist[neejber] = alt
				from[neejber] = set.New(node)
				prev[neejber] = &node
			} else if alt <= ndist {
				from[neejber].Add(node)
			}
		}
	}

	return SimpleDijkstra[T]{start, dist, prev, from}
}

func WeightConstant[T comparable](value int) EdgeWeightFunc[T] {
	return func (_, _ T) int {
		return value
	}
}

func ConstructDijkstra[T comparable](start T, neejbers NeejberFunc[T]) Dijkstra[T] {
	return ConstructWeightedDijkstra(start, neejbers, WeightConstant[T](1))
}

func (d SimpleDijkstra[T]) ShortestPathTo(end T) []T {
	path := []T{}
	node := &end
	for node != nil && *node != d.Start {
		path = append([]T{*node}, path...)
		ok := true
		node, ok = d.prev[*node]
		if !ok {
			node = nil
		}
	}

	if node == nil {
		return nil
	}

	return path
}

func (d SimpleDijkstra[T]) ForEachNode(forEach func(node T) bool) {
	for node, _ := range d.prev {
		if !forEach(node) {
			return
		}
	}
}

func (d SimpleDijkstra[T]) ShortestPathLengthTo(end T) int {
	dist, ok := d.dist[end]
	if !ok {
		return INFINITE
	}
	return dist
}

func (d SimpleDijkstra[T]) ShortestPathToFunc(end T, complete PathConsumer[T]) {
	from, ok := d.from[end]
	if !ok {
		return
	}

	from.ForEach(func(prev T) {
		if prev == d.Start {
			complete([]T{prev, end})
		} else {
			d.ShortestPathToFunc(prev, func(path []T) {
				complete(append(path, end))
			})
		}
	})
}

func ShortestPath[T comparable](start, end T, neejbers NeejberFunc[T]) ([]T, error) {

	d := ControlledDijkstra(start, neejbers, func(node T) bool { return node == end })

	path := d.ShortestPathTo(end)
	if path == nil {
		return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
	}

	return path, nil
}

func closest[T comparable](Q *set.Set[T], dist DistMap[T]) (T, error) {
	shortest := INFINITE
	snode := *new(T)
	found := false
	Q.ForEach(func(node T) {
		d := dist[node]
		if d < shortest {
			shortest = d
			snode = node
			found = true
		}
	})

	if !found {
		return snode, fmt.Errorf("could not find closest node")
	}

	return snode, nil
}

func ConstructBreadthFirst[T comparable](start T, neejbers NeejberFunc[T]) BreadthFirst[T] {
	return SimpleBFS[T]{start, neejbers}
}


func (bfs SimpleBFS[T]) AllPathsTo(end T) [][]T {
	paths := [][]T{}
	bfs.AllPathsFunc(end, func(path []T) {
		paths = append(paths, path)
	})
	return paths
}

func (bfs SimpleBFS[T]) AllPathsFunc(end T, complete PathConsumer[T]) {
	bfs.seek([]T{bfs.Start}, end, complete, NeverVisited)
}

func (bfs SimpleBFS[T]) AllPathsFuncVisited(end T, complete PathConsumer[T], visited VisitedFunc[T]) {
	bfs.seek([]T{bfs.Start}, end, complete, visited)
}

func (bfs SimpleBFS[T]) seek(path []T, end T, complete PathConsumer[T], visited VisitedFunc[T]) {
	last := path[len(path)-1]
	for _, neejber := range bfs.neejberFunc(last) {
		if neejber == end {
			complete(append(path, neejber))
			continue
		}

		if visited(path, neejber) {
			continue
		}

		bfs.seek(append(path, neejber), end, complete, visited)
	}
}