package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Direction int
type Board map[Coordinate]Coordinate

const (
	North Direction = iota
	South
	West
	East
	All
)

type Coordinate struct {
	x, y int
}

func (d Direction) String() string {
	return [...]string{"North", "South", "West", "East", "All"}[d]
}

func (d Direction) next() Direction {
	d++
	if d >= All {
		d = North
	}
	return d
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

	part1(parse(lines))
	part2(parse(lines))
}

func parse(lines []string) Board {
	b := make(Board)
	for row, line := range lines {
		for col, r := range line {
			if r == '#' {
				b[Coordinate{col, row}] = Coordinate{col, row}
			}
		}
	}
	return b
}

func part1(b Board) {
	d := North
	for i := 0; i < 10; i++ {
		b = b.doRound(d)
		d = d.next()
	}
	fmt.Printf("Part 1 - After 10 rounds, the board contains %d elves and %d empty tiles\n", len(b), b.getEmptyTiles())
}

func part2(b Board) {
	d := North
	for i := 0; ; i++ {
		next := b.doRound(d)
		if next.equals(b) {
			fmt.Printf("Part 2 - No elves moved on turn %d\n", i+1)
			break
		}

		b = next
		d = d.next()
	}
}

func (b Board) doRound(firstDir Direction) Board {
	nextB := make(Board)

	// Iterate over each elf's current coordinate
	for c := range b {
		// Figure out where this elf would like to move based on the current board state
		nextC := c.getProposed(b, firstDir)

		// Check if that proposed location is already occupied
		if prev, ok := nextB[nextC]; ok {
			// If it was, move the elf there back to its previous location and free up that spot. No one moves there
			nextB[prev] = prev
			nextB[c] = c
			delete(nextB, nextC)
		} else {
			// The spot is free, so take it (but record the previous position in case it needs to move back)
			nextB[nextC] = c
		}
	}

	return nextB
}

func (b Board) equals(other Board) bool {
	if len(b) != len(other) {
		return false
	}
	for k := range b {
		if !other.contains(k.x, k.y) {
			return false
		}
	}
	return true
}

func (c Coordinate) getProposed(b Board, firstDir Direction) Coordinate {
	// If the 8 cells around this one are empty, this one stays put
	if !b.contains(c.x-1, c.y-1) && !b.contains(c.x, c.y-1) && !b.contains(c.x+1, c.y-1) &&
		!b.contains(c.x-1, c.y) && !b.contains(c.x+1, c.y) &&
		!b.contains(c.x-1, c.y+1) && !b.contains(c.x, c.y+1) && !b.contains(c.x+1, c.y+1) {
		return c
	}

	// Figure out where this elf would like to move based on the current board state
	d, done, nextC := firstDir, false, c
	for i := 0; i < int(All) && !done; i++ {
		switch d {
		case North:
			if !b.contains(c.x-1, c.y-1) && !b.contains(c.x, c.y-1) && !b.contains(c.x+1, c.y-1) {
				nextC = Coordinate{c.x, c.y - 1}
				done = true
			}
		case South:
			if !b.contains(c.x-1, c.y+1) && !b.contains(c.x, c.y+1) && !b.contains(c.x+1, c.y+1) {
				nextC = Coordinate{c.x, c.y + 1}
				done = true
			}
		case West:
			if !b.contains(c.x-1, c.y-1) && !b.contains(c.x-1, c.y) && !b.contains(c.x-1, c.y+1) {
				nextC = Coordinate{c.x - 1, c.y}
				done = true
			}
		case East:
			if !b.contains(c.x+1, c.y-1) && !b.contains(c.x+1, c.y) && !b.contains(c.x+1, c.y+1) {
				nextC = Coordinate{c.x + 1, c.y}
				done = true
			}
		}
		d = d.next()
	}
	return nextC
}

func (b Board) contains(x, y int) bool {
	if _, ok := b[Coordinate{x, y}]; ok {
		return true
	}
	return false
}

func (b Board) getEmptyTiles() int {
	minX, minY, maxX, maxY := b.getDimensions()
	return (1+maxX-minX)*(1+maxY-minY) - len(b)
}

func (b Board) getDimensions() (minX, minY, maxX, maxY int) {
	minX, minY, maxX, maxY = math.MaxInt, math.MaxInt, math.MinInt, math.MinInt
	for c := range b {
		minX, minY, maxX, maxY = min(c.x, minX), min(c.y, minY), max(c.x, maxX), max(c.y, maxY)
	}
	return
}

/*
func (b Board) print() {
	minX, minY, maxX, maxY := b.getDimensions()
	for row := minY; row <= maxY; row++ {
		for col := minX; col <= maxX; col++ {
			if _, ok := b[Coordinate{col, row}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}
*/

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
