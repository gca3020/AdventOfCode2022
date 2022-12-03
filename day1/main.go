package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	calories := make([]string, 0)
	for s.Scan() {
		calories = append(calories, s.Text())
	}

	part1(calories)
	part2(calories)
}

func part1(calorieCounts []string) {
	var currentCalorieCount, highestCalorieCount = 0, 0
	for _, calories := range calorieCounts {
		i, err := strconv.Atoi(calories)
		if err != nil {
			if currentCalorieCount > highestCalorieCount {
				highestCalorieCount = currentCalorieCount
			}
			currentCalorieCount = 0
			continue
		}

		currentCalorieCount += i
	}

	fmt.Println("Highest Calorie Count:", highestCalorieCount)
}

func part2(calorieCounts []string) {
	calorieSums := make([]int, 0)
	var currentCalorieCount = 0
	for _, calories := range calorieCounts {
		i, err := strconv.Atoi(calories)
		if err != nil {
			calorieSums = append(calorieSums, currentCalorieCount)
			currentCalorieCount = 0
			continue
		}

		currentCalorieCount += i
	}

	sort.Slice(calorieSums, func(i, j int) bool { return calorieSums[i] > calorieSums[j] })

	fmt.Println("1", calorieSums[0], "2", calorieSums[1], "3", calorieSums[2])
	fmt.Println("total", calorieSums[0]+calorieSums[1]+calorieSums[2])
}
