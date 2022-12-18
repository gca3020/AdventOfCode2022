package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type boundingBox struct {
	min, max coordinate
}

type coordinate struct {
	x, y, z int
}

type puzzle struct {
	grid map[coordinate]bool
	gas  map[coordinate]bool
}

func main() {
	//file := "input_test"
	file := "input.txt"

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

	p := puzzle{make(map[coordinate]bool), make(map[coordinate]bool)}
	p.parse(lines)

	p.part1()
	p.part2()
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		c := coordinate{0, 0, 0}
		fmt.Sscanf(line, "%d,%d,%d", &c.x, &c.y, &c.z)
		p.grid[c] = true
	}
}

func (p *puzzle) part1() {
	exposed := 0
	for c := range p.grid {
		exposed += c.surfaceArea(p.grid)
	}
	fmt.Println("Part 1 - The object has a total surface area of", exposed)
}

func (p *puzzle) part2() {
	// Compute the bounding box of the shape, then expand it by 2 in all 6 dimensions
	bounds := p.getBounds()
	gasBounds := bounds.resize(2)

	// Recursively fill the large box with gas
	p.fillGas(gasBounds.min, gasBounds)

	// Now compute the surface area of the gas, but only those in the smaller grid.
	// The surface area of the object that gets exposed to the gas is the same as
	// the "surface area" for the gas inside the smaller bounding box
	exposed := 0
	gasBounds = bounds.resize(1)
	for c := range p.gas {
		if c.inside(gasBounds) {
			exposed += c.surfaceArea(p.gas)
		}
	}

	fmt.Println("Part 2 - The object has an exterior surface area of", exposed)
}

func (p puzzle) getBounds() (b boundingBox) {
	b.min.x, b.min.y, b.min.z = math.MaxInt, math.MaxInt, math.MaxInt
	b.max.x, b.max.y, b.max.z = math.MinInt, math.MinInt, math.MinInt
	for c := range p.grid {
		b.min.x, b.min.y, b.min.z = min(c.x, b.min.x), min(c.y, b.min.y), min(c.z, b.min.z)
		b.max.x, b.max.y, b.max.z = max(c.x, b.max.x), max(c.y, b.max.y), max(c.z, b.max.z)
	}
	return
}

func (b boundingBox) resize(size int) boundingBox {
	return boundingBox{
		min: coordinate{b.min.x - size, b.min.y - size, b.min.z - size},
		max: coordinate{b.max.x + size, b.max.y + size, b.max.z + size},
	}
}

func (c coordinate) inside(b boundingBox) bool {
	return c.x >= b.min.x && c.y >= b.min.y && c.z >= b.min.z &&
		c.x <= b.max.x && c.y <= b.max.y && c.z <= b.max.z
}

func (p *puzzle) fillGas(c coordinate, b boundingBox) {
	// If the coordinate is in the grid, but not already filled with gas or the shape, fill it with gas
	if c.inside(b) && !contains(p.gas, c) && !contains(p.grid, c) {
		p.gas[c] = true

		// Then recursively fill its neighbors
		for _, n := range c.neighbors() {
			p.fillGas(n, b)
		}
	}
}

func (c coordinate) neighbors() []coordinate {
	n := make([]coordinate, 0, 6)
	n = append(n, coordinate{c.x - 1, c.y, c.z})
	n = append(n, coordinate{c.x + 1, c.y, c.z})
	n = append(n, coordinate{c.x, c.y - 1, c.z})
	n = append(n, coordinate{c.x, c.y + 1, c.z})
	n = append(n, coordinate{c.x, c.y, c.z - 1})
	n = append(n, coordinate{c.x, c.y, c.z + 1})
	return n
}

func contains(grid map[coordinate]bool, c coordinate) bool {
	_, ok := grid[c]
	return ok
}

func (c coordinate) surfaceArea(grid map[coordinate]bool) int {
	adjacent := 0
	for _, n := range c.neighbors() {
		if contains(grid, n) {
			adjacent++
		}
	}
	return 6 - adjacent
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
