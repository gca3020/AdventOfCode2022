package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type monkeyMap map[string]interface{}

type monkey struct {
	operation rune
	m1, m2    string
}

func duration(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
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

	part1(parse(lines))
	part2(parse(lines))
}

func part1(m monkeyMap) {
	defer duration(time.Now(), "Part 1")
	fmt.Println("The root monkey's value is", m.getValue("root"))
}

func part2(m monkeyMap) {
	defer duration(time.Now(), "Part 2")
}

func parse(lines []string) monkeyMap {
	monkeys := make(monkeyMap)
	for _, line := range lines {
		name, op, _ := strings.Cut(line, ": ")
		fields := strings.Fields(op)
		if len(fields) == 1 {
			i, _ := strconv.Atoi(fields[0])
			monkeys[name] = i
		} else {
			m := monkey{rune(fields[1][0]), fields[0], fields[2]}
			monkeys[name] = m
		}
	}
	return monkeys
}

func (m monkeyMap) getValue(name string) int {
	switch val := m[name].(type) {
	case int:
		return val
	case monkey:
		return doOperation(val.operation, m.getValue(val.m1), m.getValue(val.m2))
	default:
		panic(fmt.Sprint("Error looking up monkey:", name))
	}
}

func doOperation(op rune, a, b int) int {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	default:
		panic(fmt.Sprintf("Invalid Operation: %c", op))
	}
}
