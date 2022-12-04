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

	fullOverlapCount := 0
	overlapCount := 0
	for _, line := range lines {
		l1, h1, l2, h2 := parseLine(line)
		if fullyOverlaps(l1, h1, l2, h2) {
			fullOverlapCount++
		}
		if overlapsAtAll(l1, h1, l2, h2) {
			overlapCount++
		}
	}
	fmt.Println("Fully Overlapping Ranges", fullOverlapCount)
	fmt.Println("Any Overlapping Ranges", overlapCount)
}

func parseLine(line string) (l1, h1, l2, h2 int) {
	elves := strings.Split(line, ",")
	e1 := strings.Split(elves[0], "-")
	e2 := strings.Split(elves[1], "-")
	l1, _ = strconv.Atoi(e1[0])
	h1, _ = strconv.Atoi(e1[1])
	l2, _ = strconv.Atoi(e2[0])
	h2, _ = strconv.Atoi(e2[1])

	return l1, h1, l2, h2
}

func fullyOverlaps(l1, h1 int, l2, h2 int) bool {
	if ((l1 >= l2) && (h1 <= h2)) || ((l2 >= l1) && (h2 <= h1)) {
		return true
	}
	return false
}

func overlapsAtAll(l1, h1, l2, h2 int) bool {
	if ((l1 <= l2) && (h1 >= l2)) || ((l2 <= l1) && (h2 >= l1)) {
		return true
	}
	return false
}
