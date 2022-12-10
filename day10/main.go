package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	values []int
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

	p := Puzzle{make([]int, 0)}
	p.values = append(p.values, 1)
	p.run(lines)

	// Part 1 - Compute Signal Strength
	sum := 0
	for i := 20; i < len(p.values); i += 40 {
		sum += i * p.values[i-1]
	}
	fmt.Println("The sum of the signal strengths is", sum)

	// Debug
	/*
		for i, value := range p.values {
			fmt.Printf("%03d: %d\n", i, value)
		}
	*/

	// Part 2 - Draw the output
	p.draw()

	/*
		###..###....##..##..####..##...##..###..
		#..#.#..#....#.#..#....#.#..#.#..#.#..#.
		###..#..#....#.#..#...#..#....#..#.#..#.
		#..#.###.....#.####..#...#.##.####.###..
		#..#.#....#..#.#..#.#....#..#.#..#.#....
		###..#.....##..#..#.####..###.#..#.#....
	*/
}

func (p *Puzzle) run(instructions []string) {
	for _, instruction := range instructions {
		fields := strings.Fields(instruction)
		if fields[0] == "noop" {
			p.doNoop()
		} else {
			x, _ := strconv.Atoi(fields[1])
			p.doAddX(x)
		}
	}
}

// doNoop takes a single cycle and does nothing to the X value
func (p *Puzzle) doNoop() {
	p.values = append(p.values, p.values[len(p.values)-1])
}

// doAddX takes two cycles, and adds
func (p *Puzzle) doAddX(x int) {
	curX := p.values[len(p.values)-1]
	p.values = append(p.values, curX, curX+x)
}

func (p *Puzzle) draw() {
	for i := 0; i < len(p.values); i++ {
		pixelMid := p.values[i]
		drawPos := i % 40
		if drawPos == pixelMid-1 || drawPos == pixelMid || drawPos == pixelMid+1 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
		if i%40 == 39 {
			fmt.Print("\n")
		}
	}
}
