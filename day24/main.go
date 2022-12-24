package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coordinate struct {
	row, col int
}

type Puzzle struct {
	wind []string

	start Coordinate
	end   Coordinate
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

	p := parse(lines)
	p.part1()
	p.part2()
}

func parse(lines []string) Puzzle {
	p := Puzzle{lines, Coordinate{0, 0}, Coordinate{0, 0}}
	for row, line := range lines {
		for col, r := range line {
			if row == 0 && r == '.' {
				p.start.row, p.start.col = row, col
			}
			if row == len(lines)-1 && r == '.' {
				p.end.row, p.end.col = row, col
			}
		}
	}
	return p
}

func (p Puzzle) part1() {
	if turns, ok := p.shortestTurns(p.start, p.end, 0); ok {
		fmt.Println("Part 1 - The shortest number of turns to the goal is", turns)
		return
	}
}

func (p Puzzle) part2() {
	turns, ok := p.shortestTurns(p.start, p.end, 0)
	if !ok {
		fmt.Println("No valid paths!")
	}

	turns, ok = p.shortestTurns(p.end, p.start, turns)
	if !ok {
		fmt.Println("No valid paths!")
	}

	turns, ok = p.shortestTurns(p.start, p.end, turns)
	if !ok {
		fmt.Println("No valid paths!")
	}

	fmt.Println("Part 2 - The shortest number of turns to the goal (back and forth) is", turns)
}

func (p *Puzzle) shortestTurns(start, end Coordinate, startTurn int) (turns int, ok bool) {
	thisTurnLocations := make(map[Coordinate]bool)

	turn := startTurn
	thisTurnLocations[start] = true

	// Iterate over all of the possible locations we could be on this turn, adding the places we could possibly go
	for {
		nextTurnLocations := make(map[Coordinate]bool)

		for loc := range thisTurnLocations {
			next := getNext(loc)
			for _, n := range next {
				// Return if we've found our end conditions
				if n == end {
					return turn + 1, true
				}

				// If this cell is open next turn, add it to our list of possible places to be
				if p.isOpen(n.row, n.col, turn+1) {
					nextTurnLocations[n] = true
				}
			}
		}

		thisTurnLocations = nextTurnLocations
		turn++

		if turn > 2000 {
			return 0, false
		}
	}
}

func getNext(c Coordinate) []Coordinate {
	next := make([]Coordinate, 0, 5)
	next = append(next, c)
	next = append(next, Coordinate{c.row - 1, c.col}, Coordinate{c.row + 1, c.col})
	next = append(next, Coordinate{c.row, c.col - 1}, Coordinate{c.row, c.col + 1})
	return next
}

func (p Puzzle) isOpen(row, col int, turn int) bool {
	numRows := len(p.wind) - 2
	numCols := len(p.wind[0]) - 2

	if row == p.start.row && col == p.start.col {
		return true
	}
	if row == p.end.row && col == p.end.col {
		return true
	}

	if row < 1 || row > numRows {
		return false
	}
	if col < 1 || col > numCols {
		return false
	}

	rightCol := wrapAround(col, turn*-1, 1, numCols)
	leftCol := wrapAround(col, turn, 1, numCols)
	downRow := wrapAround(row, turn*-1, 1, numRows)
	upRow := wrapAround(row, turn, 1, numRows)

	if p.wind[row][rightCol] == '>' {
		return false
	}
	if p.wind[row][leftCol] == '<' {
		return false
	}
	if p.wind[downRow][col] == 'v' {
		return false
	}
	if p.wind[upRow][col] == '^' {
		return false
	}
	return true
}

func wrapAround(val, delta int, min, max int) int {
	mod := max + 1 - min
	val += delta - min
	val += (1 - val/mod) * mod
	return val%mod + min
}

/*
func (p Puzzle) draw(turn int) {
	fmt.Println("Drawing Turn", turn)
	for row := 0; row < len(p.wind); row++ {
		for col := 0; col < len(p.wind[0]); col++ {
			if row == 0 || row == len(p.wind)-1 || col == 0 || col == len(p.wind[0])-1 {
				fmt.Print("#")
			} else {
				if p.isOpen(row, col, turn) {
					fmt.Print(".")
				} else {
					fmt.Print("x")
				}
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}
*/
