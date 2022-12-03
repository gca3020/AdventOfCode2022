package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
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

	totalPriority := 0
	for _, contents := range lines {
		totalPriority += itemPriority(findSharedItem(contents))
	}
	fmt.Println("Total Rucksack Priority:", totalPriority)

	badgePriority := 0
	for i := 0; i < len(lines); i += 3 {
		badgePriority += itemPriority(findBadge(lines[i], lines[i+1], lines[i+2]))
	}
	fmt.Println("Total Badge Priority:", badgePriority)
}

func findSharedItem(contents string) rune {
	common := findCommonItems(contents[:len(contents)/2], contents[len(contents)/2:])
	fmt.Println("Common Items in Rucksack: ", common)
	return rune(common[0])
}

func findCommonItems(s1, s2 string) string {
	common := ""
	for _, i := range s1 {
		if strings.ContainsRune(s2, i) {
			common += string(i)
		}
	}
	return common
}

func findBadge(elf1, elf2, elf3 string) rune {
	common := findCommonItems(elf1, elf2)
	common = findCommonItems(common, elf3)
	fmt.Println("Common Items In Elves", common)
	return rune(common[0])
}

func itemPriority(item rune) int {
	if unicode.IsUpper(item) {
		return int(item-'A') + 27
	}
	return int(item-'a') + 1
}
