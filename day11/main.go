package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Puzzle struct {
	monkeys  []Monkey
	worryMod int
}

type Monkey struct {
	items        []int
	operation    func(old int) int
	test         func(old int) bool
	onTrue       int
	onFalse      int
	numInspected int
}

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

	// Build the puzzle and parse the input
	part1 := Puzzle{make([]Monkey, 0), 0}
	part1.parse(lines)

	for i := 0; i < 20; i++ {
		part1.doRound()
	}
	fmt.Println("Monkey Business (Part 1):", part1.monkeyBusiness())

	// Build the puzzle and parse the input
	part2 := Puzzle{make([]Monkey, 0), 1}
	part2.parse(lines)

	for i := 0; i < 10000; i++ {
		part2.doRound()
	}
	fmt.Println("Monkey Business (Part 2):", part2.monkeyBusiness())
}

func (p *Puzzle) parse(lines []string) {
	cur := 0
	for _, line := range lines {
		if strings.Contains(line, "Monkey") {
			p.monkeys = append(p.monkeys, Monkey{nil, nil, nil, 0, 0, 0})
			cur = len(p.monkeys) - 1
			//fmt.Println("Building Monkey", cur)
		}

		if strings.Contains(line, "Starting items:") {
			p.monkeys[cur].items = make([]int, 0)
			ints := strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), ", ")
			for _, intStr := range ints {
				if val, err := strconv.Atoi(intStr); err == nil {
					p.monkeys[cur].items = append(p.monkeys[cur].items, val)
				}
			}
			//fmt.Println("Starting Items:", p.monkeys[cur].items)
		}

		if strings.Contains(line, "Operation:") {
			fields := strings.Fields(line)

			// Parse the operation (Add and Multiply only)
			//fmt.Print("Operation: ")
			var op func(x, y int) int
			switch fields[4] {
			case "+":
				//fmt.Print("+")
				op = func(x, y int) int { return x + y }
			case "*":
				//fmt.Print("*")
				op = func(x, y int) int { return x * y }
			default:
				//fmt.Println("Could not parse operation")
			}

			// Now parse the value (with special case to handle "old * old")
			if val, err := strconv.Atoi(fields[5]); err != nil {
				//fmt.Println("old")
				p.monkeys[cur].operation = func(old int) int {
					return op(old, old)
				}
			} else {
				//fmt.Println(val)
				p.monkeys[cur].operation = func(old int) int {
					return op(old, val)
				}
			}
		}

		if strings.Contains(line, "Test:") {
			val, _ := strconv.Atoi(strings.Fields(line)[3])
			p.monkeys[cur].test = func(old int) bool {
				return old%val == 0
			}
			p.worryMod *= val
			//fmt.Println("Test:", val)
		}

		if strings.Contains(line, "If true:") {
			p.monkeys[cur].onTrue, _ = strconv.Atoi(strings.Fields(line)[5])
			//fmt.Println("On True:", p.monkeys[cur].onTrue)
		}

		if strings.Contains(line, "If false:") {
			p.monkeys[cur].onFalse, _ = strconv.Atoi(strings.Fields(line)[5])
			//fmt.Println("On False:", p.monkeys[cur].onFalse)
		}
	}
}

func (p *Puzzle) doRound() {
	for i := 0; i < len(p.monkeys); i++ {
		p.monkeys[i].takeTurn(p.monkeys, p.worryMod)
	}
}

func (p *Puzzle) monkeyBusiness() int {
	inspections := make([]int, 0)
	for _, m := range p.monkeys {
		inspections = append(inspections, m.numInspected)
	}

	sort.Slice(inspections, func(i, j int) bool {
		return inspections[i] > inspections[j]
	})

	return inspections[0] * inspections[1]
}

func (m *Monkey) takeTurn(monkeys []Monkey, worryMod int) {
	for _, item := range m.items {
		m.numInspected++         // Record the inspection
		item = m.operation(item) // Do the operation

		// Relief - Parts 1 & 2
		if worryMod == 0 {
			item /= 3 // If this is 0, use the Part 1 logic
		} else {
			item %= worryMod // Modulo by the LCM of all the monkey's divisors
		}

		if m.test(item) { // Test the item's value
			monkeys[m.onTrue].catch(item) // Throw to the "if true" neighbor
		} else {
			monkeys[m.onFalse].catch(item) // Throw to the "if false" neighbor
		}
	}
	// Reset the slice now that we're done processing all of this monkey's items
	m.items = make([]int, 0)
}

func (m *Monkey) catch(item int) {
	m.items = append(m.items, item)
}
