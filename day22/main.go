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

	p1 := parse(lines)
	p1.run(false)

	p2 := parse(lines)
	p2.run(true)
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

	return Puzzle{b, steps, Player{0, b.firstCol(0), right}}
}

func (p *Puzzle) run(part2 bool) {
	// TODO: Run through the maze
	for _, step := range p.steps {
		switch s := step.(type) {
		case int:
			if part2 {
				p.player.moveCube(s, p.board)
			} else {
				p.player.moveFlat(s, p.board)
			}
		case Direction:
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

func (p Player) nextLocationFlat(b Board) (row, col int) {
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

func (p *Player) moveFlat(n int, b Board) {
	for i := 0; i < n; i++ {
		r, c := p.nextLocationFlat(b)
		if b[r][c] == '#' {
			// We hit a rock, break early
			break
		}
		p.row, p.col = r, c
	}
}

func getCubeFace(row, col int) int {
	if row >= 0 && row <= 49 && col >= 50 && col <= 99 {
		return 0
	}
	if row >= 0 && row <= 49 && col >= 100 && col <= 149 {
		return 1
	}
	if row >= 50 && row <= 99 && col >= 50 && col <= 99 {
		return 2
	}
	if row >= 100 && row <= 149 && col >= 0 && col <= 49 {
		return 3
	}
	if row >= 100 && row <= 149 && col >= 50 && col <= 99 {
		return 4
	}
	if row >= 150 && row <= 199 && col >= 0 && col <= 49 {
		return 5
	}
	return -1
}

func (p Player) nextLocationCube(b Board) (row, col int, dir Direction) {
	face := getCubeFace(p.row, p.col)
	row, col, dir = p.row, p.col, p.facing
	switch p.facing {
	case right:
		col++
		if col > b.lastCol(row) {
			switch face {
			case 1:
				row, col, dir = (49-p.row)+100, 99, left
			case 2:
				row, col, dir = 49, 50+p.row, up
			case 4:
				row, col, dir = 49-(p.row-100), 149, left
			case 5:
				row, col, dir = 149, p.row-100, up
			default:
				fmt.Println("Can't move right from face", face)
			}
		}
	case down:
		row++
		if row > b.lastRow(col) {
			switch face {
			case 1:
				row, col, dir = p.col-50, 99, left
			case 4:
				row, col, dir = 150+(p.col-50), 49, left
			case 5:
				row, col, dir = 0, p.col+100, down
			default:
				fmt.Println("Can't move down from face", face)
			}
		}
	case left:
		col--
		if col < b.firstCol(row) {
			switch face {
			case 0:
				row, col, dir = (49-p.row)+100, 0, right
			case 2:
				row, col, dir = 100, p.row-50, down
			case 3:
				row, col, dir = 49-(p.row-100), 50, right
			case 5:
				row, col, dir = 0, 50+(p.row-150), down
			default:
				fmt.Println("Can't move left from face", face)
			}
		}
	case up:
		row--
		if row < b.firstRow(col) {
			switch face {
			case 0:
				row, col, dir = 150+(p.col-50), 0, right
			case 1:
				row, col, dir = 199, col-100, up
			case 3:
				row, col, dir = col+50, 50, right
			default:
				fmt.Println("Can't move up from face", face)
			}
		}
	}

	newFace := getCubeFace(row, col)
	if face != newFace {
		fmt.Printf("Moving from %d-(%d,%d,%v) to %d-(%d,%d,%v)(%c)\n", face, p.row, p.col, p.facing, newFace, row, col, dir, b[row][col])
	}
	return
}

func (p *Player) moveCube(n int, b Board) {
	for i := 0; i < n; i++ {
		r, c, d := p.nextLocationCube(b)
		if b[r][c] == '#' {
			break
		}
		p.row, p.col, p.facing = r, c, d
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

func (d Direction) String() string {
	switch d {
	case right:
		return "right"
	case down:
		return "down"
	case left:
		return "left"
	case up:
		return "up"
	}
	return ""
}
