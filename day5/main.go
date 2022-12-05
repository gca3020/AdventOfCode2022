package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	fmt.Println("Parsing Stacks")
	stacks := parseStacks(lines)
	printStacks(stacks)

	fmt.Println("Running algorithm")
	for _, line := range lines {
		move, from, to := 0, 0, 0
		if _, err := fmt.Sscanf(line, "move %d from %d to %d\n", &move, &from, &to); err != nil {
			continue
		}
		runStep(move, from, to, stacks)
	}
	printStacks(stacks)
}

/*
Part1 - Left for posterity
func runStep(move, from, to int, stacks map[int][]rune) {
	for i := 0; i < move; i++ {
		// Get the top of the "from" stack
		top := len(stacks[from]) - 1
		c := stacks[from][top]
		// Pop the item off the "from" stack
		stacks[from] = stacks[from][:top]
		// Add the item to the top of the "to" stack
		stacks[to] = append(stacks[to], c)
	}
}
*/

func runStep(move, from, to int, stacks map[int][]rune) {
	// Get the top "move" items off the from stack
	top := len(stacks[from])
	crates := stacks[from][top-move:]
	// Pop the items off the "From" stack
	stacks[from] = stacks[from][:top-move]
	// Append the items to the "To" stack
	stacks[to] = append(stacks[to], crates...)
}

func parseStacks(lines []string) (stacks map[int][]rune) {
	// Find the container IDs
	stacks = make(map[int][]rune)
	var idLine int
	for i, line := range lines {
		if line == "" {
			idLine = i - 1
			break
		}
	}

	// Create the stacks
	stackIds := strings.Fields(lines[idLine])
	for _, id := range stackIds {
		stackId, _ := strconv.Atoi(id)
		stacks[stackId] = make([]rune, 0)
	}

	// Fill the stacks
	for lineNum := idLine - 1; lineNum >= 0; lineNum-- {
		for i := 0; i < len(stacks); i++ {
			r := rune(lines[lineNum][(4*i)+1])
			if r != ' ' {
				stacks[i+1] = append(stacks[i+1], r)
			}
		}
	}

	return stacks
}

func printStacks(stacks map[int][]rune) {
	for i := 0; i < len(stacks); i++ {
		fmt.Printf("%d: %c\n", i+1, stacks[i+1])
	}
}
