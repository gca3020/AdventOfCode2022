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
	fmt.Println("Part 1 - The root monkey's value is", m.getValue("root"))
}

func part2(m monkeyMap) {
	defer duration(time.Now(), "Part 2")

	// Update the root monkey to have an equality op
	root := m["root"].(monkey)
	root.operation = '='
	m["root"] = root

	fmt.Println("Part 2 - The human value to achieve equality is", m.findHumanValue("root", 0))
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

func (m monkeyMap) isVariable(name string) bool {
	if name == "humn" {
		return true
	}

	switch val := m[name].(type) {
	case int:
		return false
	case monkey:
		return m.isVariable(val.m1) || m.isVariable(val.m2)
	default:
		panic(fmt.Sprint("Error looking up monkey:", name))
	}
}

func (m monkeyMap) findHumanValue(name string, exp int) int {
	if name == "humn" {
		return exp
	}

	v := m[name].(monkey)
	m1Var, m1Val, m1Str := m.isVariable(v.m1), 0, ""
	m2Var, m2Val, m2Str := m.isVariable(v.m2), 0, ""
	if m1Var && m2Var {
		panic("Both sides of the equation are variable!")
	}

	// Do some stuff so we can print some helpful debug
	if m1Var {
		m1Str = "x"
		m2Val = m.getValue(v.m2)
		m2Str = fmt.Sprint(m2Val)
	}
	if m2Var {
		m2Str = "x"
		m1Val = m.getValue(v.m1)
		m1Str = fmt.Sprint(m1Val)
	}
	fmt.Printf("[%s] Finding Human \"x\" such that (%s %c %s = %d)\n", name, m1Str, v.operation, m2Str, exp)

	switch v.operation {
	case '=':
		if m1Var {
			return m.findHumanValue(v.m1, m2Val)
		}
		return m.findHumanValue(v.m2, m1Val)
	case '+':
		if m1Var {
			return m.findHumanValue(v.m1, exp-m2Val)
		}
		return m.findHumanValue(v.m2, exp-m1Val)
	case '-':
		if m1Var {
			return m.findHumanValue(v.m1, exp+m2Val)
		}
		return m.findHumanValue(v.m2, m1Val-exp)
	case '*':
		if m1Var {
			return m.findHumanValue(v.m1, exp/m2Val)
		}
		return m.findHumanValue(v.m2, exp/m1Val)
	case '/':
		if m1Var {
			return m.findHumanValue(v.m1, exp*m2Val)
		}
		return m.findHumanValue(v.m2, m1Val/exp)
	default:
		panic(fmt.Sprintf("Unhandled Operation: \"%c\"!\n", v.operation))
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
