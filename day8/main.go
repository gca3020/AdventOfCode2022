package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

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

	// Build the grid
	grid := buildGrid(lines)

	// Part 1, count visible
	n := numVisible(grid)
	fmt.Printf("Counted %d visible trees\n", n)

	// Part 2, scenic score
	scenicScore, row, col := bestScenicScore(grid)
	fmt.Printf("The best scenic score was %d and found at (%d, %d)\n", scenicScore, row, col)
}

func duration(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func buildGrid(lines []string) [][]int {
	defer duration(time.Now(), "buildGrid")
	grid := make([][]int, len(lines))
	for i, line := range lines {
		row := make([]int, len(line))
		for j, rune := range line {
			row[j] = int(rune - '0')
		}
		grid[i] = row
	}
	return grid
}

func numVisible(grid [][]int) int {
	defer duration(time.Now(), "numVisible")
	num := 0
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if isVisible(row, col, grid) {
				num++
			}
		}
	}
	return num
}

func isVisible(row, col int, grid [][]int) bool {
	height := grid[row][col]

	// From Left
	for i := 0; i <= col; i++ {
		if i == col {
			return true
		}
		if grid[row][i] >= height {
			break
		}
	}

	// From Right
	for i := len(grid[row]) - 1; i >= col; i-- {
		if i == col {
			return true
		}
		if grid[row][i] >= height {
			break
		}
	}

	// From Top
	for i := 0; i <= row; i++ {
		if i == row {
			return true
		}
		if grid[i][col] >= height {
			break
		}
	}

	// From Bottom
	for i := len(grid) - 1; i >= row; i-- {
		if i == row {
			return true
		}
		if grid[i][col] >= height {
			break
		}
	}
	return false
}

func bestScenicScore(grid [][]int) (topScore, topRow, topCol int) {
	defer duration(time.Now(), "bestScenicScore")
	topScore, topRow, topCol = 0, 0, 0
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if score := scenicScore(row, col, grid); score > topScore {
				topScore = score
				topRow = row
				topCol = col
			}
		}
	}

	return topScore, topRow, topCol
}

func scenicScore(row, col int, grid [][]int) int {
	height := grid[row][col]
	// Looking Left
	left := 0
	for i := col - 1; i >= 0; i-- {
		left++
		if grid[row][i] >= height {
			break
		}
	}

	// Looking Right
	right := 0
	for i := col + 1; i < len(grid[row]); i++ {
		right++
		if grid[row][i] >= height {
			break
		}
	}

	// Looking Up
	up := 0
	for i := row - 1; i >= 0; i-- {
		up++
		if grid[i][col] >= height {
			break
		}
	}

	// Looking Down
	down := 0
	for i := row + 1; i < len(grid); i++ {
		down++
		if grid[i][col] >= height {
			break
		}
	}

	// Calc the total
	total := left * right * up * down
	return total
}
