package main

import (
	"bufio"
	"fmt"
	"os"
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

	fmt.Println("Part1:", computeTotal(lines, false), "Part2:", computeTotal(lines, true))
}

func computeTotal(rounds []string, p2 bool) int {
	totalScore := 0
	for _, round := range rounds {
		split := strings.Split(round, " ")
		opp, player := decrypt(split[0]), decrypt(split[1])

		if p2 {
			player = computePlayerMove(opp, split[1])
		}

		totalScore += shapeScore(player) + outcomeScore(opp, player)
	}
	return totalScore
}

func decrypt(shape string) string {
	switch shape {
	case "A", "X":
		return "R"
	case "B", "Y":
		return "P"
	case "C", "Z":
		return "S"
	default:
		return ""
	}
}

func shapeScore(shape string) int {
	switch shape {
	case "R":
		return 1
	case "P":
		return 2
	case "S":
		return 3
	default:
		return 0
	}
}

func outcomeScore(opp, player string) int {
	if (opp == "R" && player == "P") || (opp == "P" && player == "S") || (opp == "S" && player == "R") {
		return 6
	}
	if (opp == "R" && player == "R") || (opp == "P" && player == "P") || (opp == "S" && player == "S") {
		return 3
	}
	return 0
}

func computePlayerMove(opp, outcome string) string {
	if (opp == "R" && outcome == "Y") || (opp == "P" && outcome == "X") || (opp == "S" && outcome == "Z") {
		return "R"
	}
	if (opp == "R" && outcome == "Z") || (opp == "P" && outcome == "Y") || (opp == "S" && outcome == "X") {
		return "P"
	}
	return "S"
}
