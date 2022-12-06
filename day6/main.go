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

	measureTime(somSearch, "somSearch4x1000", lines[0], 4, 1000)
	measureTime(somSearch, "somSearch14x1000", lines[0], 14, 1000)
	measureTime(somSearchBasic, "somSearchBasic4x1000", lines[0], 4, 1000)
	measureTime(somSearchBasic, "somSearchBasic14x1000", lines[0], 14, 1000)

	//somSearch4x1000 took 163.7761ms
	//somSearch14x1000 took 490.0602ms
	//somSearchBasic4x1000 took 207.9999ms
	//somSearchBasic14x1000 took 967.55ms

	startOfPacket := somSearch(lines[0], 4)
	startOfMessage := somSearch(lines[0], 14)
	sopBasic := somSearchBasic(lines[0], 4)
	somBasic := somSearchBasic(lines[0], 14)
	fmt.Printf("SOP=%d, SOM=%d, SOP(B)=%d, SOM(B)=%d\n", startOfPacket, startOfMessage, sopBasic, somBasic)
}

func measureTime(f func(string, int) int, name, line string, len, iterations int) {
	defer duration(time.Now(), name)
	for i := 0; i < iterations; i++ {
		f(line, len)
	}
}

func duration(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func onlyUnique(window string) bool {
	set := make(map[rune]interface{})
	for _, r := range window {
		if _, ok := set[r]; ok {
			return false
		}
		set[r] = true
	}
	return true
}

func somSearchBasic(line string, numUnique int) int {
	for i := numUnique; i < len(line); i++ {
		// Check if the slice only has unique characters
		if onlyUnique(line[i-numUnique : i]) {
			return i
		}
	}
	return 0
}

func somSearch(line string, numUnique int) int {
	window := make(map[rune]int, numUnique)
	for i, c := range line {
		// Pop the oldest character out of the set representing the sliding window
		if i >= numUnique {
			leavingWindow := rune(line[i-numUnique])
			window[leavingWindow]--
			if window[leavingWindow] <= 0 {
				delete(window, leavingWindow)
			}
		}

		// If this character is not already in the window, add it
		_, exists := window[c]
		if exists {
			window[c]++
		} else {
			window[c] = 1
		}

		// If the window has exactly numUnique keys, then we have a full SOM,
		// and the message starts on the next character
		if len(window) == numUnique {
			return i + 1
		}
	}
	return 0
}
