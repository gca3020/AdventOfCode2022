package main

import (
	"bufio"
	"fmt"
	"os"
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

	startOfPacket := somSearch(lines[0], 4)
	fmt.Println("Start of Packet is at", startOfPacket)
	startOfMessage := somSearch(lines[0], 14)
	fmt.Println("Start of Message is at", startOfMessage)
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
