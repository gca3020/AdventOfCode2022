package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Direction int
type Board []string

const (
	right Direction = iota
	down
	left
	up
)

type Player struct {
	row, col int
	facing   Direction
}

type Puzzle struct {
	board  Board
	steps  []interface{}
	player Player
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

	puzzle := parse(lines)
	puzzle.part1()
}

func parse(lines []string) Puzzle {
	b := Board(lines[:len(lines)-2])

	// Parse the line with directions -- This is ugly, but it works
	steps := make([]interface{}, 0)
	stepLine := lines[len(lines)-1]
	for i := 0; i < len(stepLine); {
		nextChar := strings.IndexAny(stepLine[i:], "RL")
		if nextChar == 0 {
			d := left
			if stepLine[i] == 'R' {
				d = right
			}
			steps = append(steps, d)
			i++
		} else {
			toParse := ""
			if nextChar > 0 {
				toParse = stepLine[i : nextChar+i]
			} else {
				toParse = stepLine[i:]
				nextChar = 10
			}
			if num, err := strconv.Atoi(toParse); err == nil {
				steps = append(steps, num)
				i += nextChar
			}
		}
	}
	fmt.Printf("%#v\n", steps)

	return Puzzle{b, steps, Player{0, b.firstCol(0), right}}
}

func (p *Puzzle) part1() {
	// TODO: Run through the maze
	for _, step := range p.steps {
		fmt.Println("The player is", p.player)
		switch s := step.(type) {
		case int:
			fmt.Printf("Moving %d\n", s)
			p.player.move(s, p.board)
		case Direction:
			fmt.Printf("Turning %c\n", s)
			p.player.turn(s)
		}
	}

	fmt.Println("Part 1 - The player score is", p.player.getScore())
}

func (b Board) firstRow(col int) int {
	for row, line := range b {
		// If the line is not long enough to cover this column, just continue
		if len(line) <= col || line[col] == ' ' {
			continue
		}
		return row
	}
	return -1
}

func (b Board) lastRow(col int) int {
	for i := len(b) - 1; i >= 0; i-- {
		if len(b[i]) <= col || b[i][col] == ' ' {
			continue
		}
		return i
	}
	return -1
}

func (b Board) firstCol(row int) int {
	for i, r := range b[row] {
		if r != ' ' {
			return i
		}
	}
	return -1
}

func (b Board) lastCol(row int) int {
	return len(b[row]) - 1
}

func (p Player) getScore() int {
	return ((p.row + 1) * 1000) + ((p.col + 1) * 4) + int(p.facing)
}

func (p Player) nextLocation(b Board) (row, col int) {
	row, col = p.row, p.col
	switch p.facing {
	case right:
		col++
		if col > b.lastCol(row) {
			col = b.firstCol(row)
		}
	case down:
		row++
		if row > b.lastRow(col) {
			row = b.firstRow(col)
		}
	case left:
		col--
		if col < b.firstCol(row) {
			col = b.lastCol(row)
		}
	case up:
		row--
		if row < b.firstRow(col) {
			row = b.lastRow(col)
		}
	}
	return
}

func (p *Player) move(n int, b Board) {
	for i := 0; i < n; i++ {
		r, c := p.nextLocation(b)
		if b[r][c] == '#' {
			// We hit a rock, break early
			break
		}
		p.row, p.col = r, c
	}
}

func (p *Player) turn(d Direction) {
	if d == left {
		p.facing--
		if p.facing < 0 {
			p.facing = up
		}
	} else {
		p.facing++
		if p.facing > up {
			p.facing = right
		}
	}
}
