package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type puzzle struct {
	unvisited map[location]node
	end       location
}

type location struct {
	row, col int
}

type node struct {
	elevation int
	distance  int
}

func main() {
	f, err := os.Open("input.txt")
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

	part1 := &puzzle{nil, location{}}
	part1.parse(lines, false)
	fmt.Printf("The shortest Part 1 path contains %d steps\n", part1.findShortestPath())

	part2 := &puzzle{nil, location{}}
	part2.parse(lines, true)
	fmt.Printf("The shortest Part 2 path contains %d steps\n", part2.findShortestPath())
}

func (p *puzzle) parse(lines []string, part2 bool) {
	p.unvisited = make(map[location]node)
	for i, line := range lines {
		for j, r := range line {
			dist := math.MaxInt
			loc := location{i, j}
			if r == 'S' {
				dist = 0
				r = 'a'
			} else if r == 'E' {
				p.end = loc
				r = 'z'
			}

			// For part 2, we want to find the shortest path from any 'a' elevation
			if part2 && r == 'a' {
				dist = 0
			}
			p.unvisited[location{i, j}] = node{int(r), dist}
		}
	}
}

func (p *puzzle) findShortestPath() int {
	for len(p.unvisited) > 0 {
		// Find the shortest unvisited node
		curLoc := p.closestUnvisited()
		curNode := p.unvisited[curLoc]

		// Early return if this is the node we care about
		if curLoc == p.end {
			return curNode.distance
		}

		// Update the distances of all of its neighbors
		for _, l := range curLoc.neighbors() {
			nextDistance := curNode.distance + 1
			if neighbor, ok := p.unvisited[l]; ok {
				// If this is a valid neighbor
				if neighbor.elevation <= curNode.elevation+1 && nextDistance < neighbor.distance {
					// Update its distance in the map
					neighbor.distance = nextDistance
					p.unvisited[l] = neighbor
				}
			}
		}

		// Remove the current node from the set
		delete(p.unvisited, curLoc)
	}
	return 0
}

func (p *puzzle) closestUnvisited() location {
	closestDist := math.MaxInt
	closestLoc := location{}
	for loc, node := range p.unvisited {
		if node.distance < closestDist {
			closestLoc = loc
			closestDist = node.distance
		}
	}
	return closestLoc
}

func (l *location) neighbors() []location {
	n := make([]location, 0, 4)
	n = append(n, location{l.row + 1, l.col}, location{l.row - 1, l.col}, location{l.row, l.col - 1}, location{l.row, l.col + 1})
	return n
}
