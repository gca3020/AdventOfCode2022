package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type puzzle struct {
	grid                   map[coordinate]rune
	xMin, xMax, yMin, yMax int
}

type coordinate struct {
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

	// Parse the input
	p := puzzle{make(map[coordinate]rune), 0, 0, 0, 0}
	p.parse(lines)

	// Compute the bounds of the playing field and then draw it
	p.computeBounds()
	p.draw()

	sandAdded := 0

	// Part 1 - Add sand until it starts falling off the bottom
	for p.addSand(false) {
		sandAdded++
	}
	p.draw()
	fmt.Printf("Added %d grains of sand in part 1\n", sandAdded)

	// Part 2 - Add sand but assuming there's a floor at yMax+2
	for p.addSand(true) {
		sandAdded++
	}
	p.computeBounds()
	p.draw()
	fmt.Printf("Added %d grains of sand total\n", sandAdded+1)
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		from := coordinate{0, 0}
		pairs := strings.Split(line, " -> ")
		for _, pair := range pairs {
			coordStr := strings.Split(pair, ",")
			x, _ := strconv.Atoi(coordStr[0])
			y, _ := strconv.Atoi(coordStr[1])
			to := coordinate{x, y}

			// If we have a previous coordinate, draw the line between them
			if from.x != 0 || from.y != 0 {
				p.drawLine(from, to)
			}
			from = to
		}
	}
}

func (p *puzzle) drawLine(from, to coordinate) {
	//fmt.Printf("(%d,%d) -> (%d,%d): ", from.x, from.y, to.x, to.y)
	xStep, yStep := 1, 1
	if to.x < from.x {
		xStep = -1
	}
	if to.y < from.y {
		yStep = -1
	}

	for x := from.x; x != to.x; x += xStep {
		//fmt.Printf("(%d,%d) ", x, from.y)
		p.grid[coordinate{x, from.y}] = '#'
	}

	for y := from.y; y != to.y; y += yStep {
		//fmt.Printf("(%d,%d) ", from.x, y)
		p.grid[coordinate{from.x, y}] = '#'
	}

	//fmt.Printf("(%d,%d)\n", to.x, to.y)
	p.grid[to] = '#'
}

func (p *puzzle) computeBounds() {
	xMin, xMax := math.MaxInt, math.MinInt
	yMin, yMax := math.MaxInt, math.MinInt
	for coordinate := range p.grid {
		if coordinate.x < xMin {
			xMin = coordinate.x
		}
		if coordinate.x > xMax {
			xMax = coordinate.x
		}
		if coordinate.y < yMin {
			yMin = coordinate.y
		}
		if coordinate.y > yMax {
			yMax = coordinate.y
		}
	}
	p.xMin = xMin
	p.xMax = xMax
	p.yMin = yMin
	p.yMax = yMax
}

func (p *puzzle) draw() {
	for y := 0; y <= p.yMax; y++ {
		for x := p.xMin; x <= p.xMax; x++ {
			r, ok := p.grid[coordinate{x, y}]
			if ok {
				fmt.Printf("%c", r)
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func (p *puzzle) addSand(part2 bool) bool {
	sand := coordinate{500, 0}
	for {
		next := sand.next(p, part2)
		if next == sand {
			//fmt.Printf("Sand has come to a rest at (%d,%d)\n", sand.x, sand.y)
			p.grid[sand] = 'o'
			if sand.y == 0 && sand.x == 500 {
				// If the sand ever covers the origin, return false indicating that we're done
				return false
			}
			return true
		}
		if sand.y > p.yMax+2 {
			//fmt.Println("Sand has fallen into oblivion")
			break
		}
		sand = next
	}
	return false
}

func (c coordinate) next(p *puzzle, part2 bool) coordinate {
	// If we're already at the bottom, just return this same position
	if part2 && c.y == p.yMax+1 {
		return c
	}

	down := coordinate{c.x, c.y + 1}
	if _, ok := p.grid[down]; !ok {
		return down
	}

	downLeft := coordinate{c.x - 1, c.y + 1}
	if _, ok := p.grid[downLeft]; !ok {
		return downLeft
	}

	downRight := coordinate{c.x + 1, c.y + 1}
	if _, ok := p.grid[downRight]; !ok {
		return downRight
	}

	return c
}
