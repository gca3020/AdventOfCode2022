package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	rope    []Point
	visited map[Point]int
}

type Point struct {
	x, y int
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

	// Part1
	runPuzzle(lines, 2)

	// Part 2
	runPuzzle(lines, 10)
}

func runPuzzle(steps []string, ropeLength int) {
	p := &Puzzle{
		rope:    make([]Point, ropeLength),
		visited: make(map[Point]int),
	}

	for _, step := range steps {
		fields := strings.Fields(step)
		count, _ := strconv.Atoi(fields[1])
		p.runStep(fields[0], count)
	}

	fmt.Println("Visited", len(p.visited), "with rope of length", ropeLength)
}

func (p *Puzzle) runStep(direction string, count int) {
	newX, newY := p.rope[0].x, p.rope[0].y
	stepX, stepY := 0, 0
	switch direction {
	case "U":
		newY += count
		stepY = 1
	case "D":
		newY -= count
		stepY = -1
	case "R":
		newX += count
		stepX = 1
	case "L":
		newX -= count
		stepX = -1
	default:
		fmt.Println("Invalid Direction:", direction)
	}

	// While the head is not in its new destination
	for p.rope[0].x != newX || p.rope[0].y != newY {
		// Move it one step
		p.rope[0].x += stepX
		p.rope[0].y += stepY

		// Then go down the tail, moving each part of the tail in turn
		for i := 1; i < len(p.rope); i++ {
			// If we ever have one part of the tail not move, then don't bother moving the rest
			if !p.rope[i].follow(p.rope[i-1]) {
				break
			}
		}

		// Record the position of the last part of the tail
		p.visited[p.rope[len(p.rope)-1]] = 1
	}
}

func (p *Point) follow(other Point) bool {
	if p.isAdjacent(other) {
		return false
	}

	if other.x-p.x > 0 {
		p.x++
	} else if other.x-p.x < 0 {
		p.x--
	}

	if other.y-p.y > 0 {
		p.y++
	} else if other.y-p.y < 0 {
		p.y--
	}
	return true
}

func (p *Point) isAdjacent(other Point) bool {
	return absDiff(p.x, other.x) <= 1 && absDiff(p.y, other.y) <= 1
}

func absDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
