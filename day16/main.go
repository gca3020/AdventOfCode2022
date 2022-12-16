package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type puzzle struct {
	valves map[string]valve
}

type valve struct {
	name      string
	rate      int
	neighbors []string
	distance  map[string]int
}

func main() {
	file := "input.txt"
	//file := "input_test"

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	lines := make([]string, 0)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	// Parse the input
	p := puzzle{make(map[string]valve)}
	p.parse(lines)

	// Do some optimizations to the input in order to speed up processing
	p.optimizeInput()

	// Run the solution
	p.part1()
	p.part2()
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		before, after, _ := strings.Cut(line, "; ")
		name, rate := "", 0
		fmt.Sscanf(before, "Valve %s has flow rate=%d", &name, &rate)
		neighbors := strings.Split(strings.SplitN(after, " ", 5)[4], ", ")
		p.valves[name] = NewValve(name, rate, neighbors)
	}
}

func (p *puzzle) optimizeInput() {
	for _, valve := range p.valves {
		// Compute the shortest path to each "important" node
		for _, dest := range p.valves {
			if dest.rate != 0 {
				valve.distance[dest.name] = len(p.getPath(valve.name, dest.name))
				p.valves[valve.name] = valve
			}
		}
	}
}

func (p *puzzle) part1() {
	important := make([]string, 0)
	for _, valve := range p.valves {
		if valve.rate != 0 {
			important = append(important, valve.name)
		}
	}
	fmt.Printf("The best path is %#v\n", p.bestScore(30, "AA", important, nil))
}

func (p *puzzle) part2() {
	important := make([]string, 0)
	for _, valve := range p.valves {
		if valve.rate != 0 {
			important = append(important, valve.name)
		}
	}

	// This is really dumb, but we just shuffle the list and split it down the middle, and run two copies of the
	// part 1 algorithm. Every time we find a better score, we print it. If this iteration count is sufficiently large,
	// eventually we will find the right answer.
	best := 0
	for i := 0; i < 10000; i++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(important), func(i, j int) { important[i], important[j] = important[j], important[i] })
		middle := len(important) / 2
		p1 := p.bestScore(26, "AA", important[:middle], nil)
		p2 := p.bestScore(26, "AA", important[middle:], nil)
		if p1+p2 > best {
			best = p1 + p2
			fmt.Printf("New Best Score %d found using search spaces %v, %v!\n", best, important[:middle], important[middle:])
		}
	}
}

func (p *puzzle) bestScore(remainingTime int, currentValve string, important, open []string) int {
	if remainingTime < 1 {
		return 0
	} else {
		best := 0
		for _, v := range important {
			if !contains(open, v) {
				newTime := remainingTime - p.valves[currentValve].distance[v] - 1
				newScore := (newTime * p.valves[v].rate) + p.bestScore(newTime, v, important, append(open, v))
				if newScore > best {
					best = newScore
				}
			}
		}
		return best
	}
}

// getPath gets the shortest path from one node to another, shamelessly cribbing from
// the Wikipedia source code for Dijkstra's algorithm
func (p puzzle) getPath(from, to string) []string {
	dist := make(map[string]int)
	prev := make(map[string]string)
	q := make(map[string]bool)

	for v := range p.valves {
		dist[v] = math.MaxInt
		prev[v] = ""
		q[v] = true
	}
	dist[from] = 0

	for len(q) != 0 {
		// Find the next closest node that's still in Q
		closestDist := math.MaxInt
		u := ""
		for name := range q {
			if dist[name] < closestDist {
				closestDist = dist[name]
				u = name
			}
		}

		// If it's our destination, we can break early
		if u == to {
			break
		}

		// Remove the node from the Q
		delete(q, u)

		// Iterate over any of its neighbors still in Q, updating their distances
		for _, v := range p.valves[u].neighbors {
			if _, ok := q[v]; ok {
				alt := dist[u] + 1
				if alt < dist[v] {
					dist[v] = alt
					prev[v] = u
				}
			}
		}
	}

	// Reverse iterate to get the path
	reversePath, u := make([]string, 0), to
	for u != "" {
		if u == from {
			// Swap the order before returning
			for i, j := 0, len(reversePath)-1; i < j; i, j = i+1, j-1 {
				reversePath[i], reversePath[j] = reversePath[j], reversePath[i]
			}
			return reversePath
		}
		reversePath = append(reversePath, u)
		u = prev[u]
	}
	return nil
}

func NewValve(name string, rate int, neighbors []string) valve {
	return valve{name, rate, neighbors, make(map[string]int)}
}

func contains(array []string, key string) bool {
	for _, v := range array {
		if v == key {
			return true
		}
	}
	return false
}
