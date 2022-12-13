package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type part1 struct {
	pairs []packetPair
}

type part2 struct {
	packets []interface{}
}

type packetPair struct {
	left, right []interface{}
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

	p1 := &part1{}
	p1.parse(lines)
	fmt.Println("Sum of Correct Indices:", p1.sumOfCorrectIndices())

	p2 := &part2{}
	p2.parse(lines)
	fmt.Println("The decoder key is:", p2.findDecoderKey())
}

func (p *part1) parse(lines []string) {
	p.pairs = make([]packetPair, 0)
	for i := 0; i < len(lines); i += 3 {
		pair := packetPair{}
		json.Unmarshal([]byte(lines[i]), &pair.left)
		json.Unmarshal([]byte(lines[i+1]), &pair.right)
		p.pairs = append(p.pairs, pair)
	}
}

func (p *part1) sumOfCorrectIndices() int {
	sumOfCorrectIndices := 0
	for i, pair := range p.pairs {
		//fmt.Println("")
		if compareAny(pair.left, pair.right) > 0 {
			sumOfCorrectIndices += i + 1
		}
	}
	return sumOfCorrectIndices
}

func (p *part2) parse(lines []string) {
	// Append the two divider packets
	lines = append(lines, "[[2]]", "[[6]]")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var packet []interface{}
		json.Unmarshal([]byte(line), &packet)
		p.packets = append(p.packets, packet)
	}
}

func (p *part2) findDecoderKey() int {
	// Sort the slice using the comparison function
	sort.Slice(p.packets, func(i, j int) bool {
		return compareAny(p.packets[i], p.packets[j]) > 0
	})

	first, second := 0, 0
	for i, packet := range p.packets {
		packetStr, _ := json.Marshal(packet)
		if string(packetStr) == "[[2]]" {
			first = i + 1
		}
		if string(packetStr) == "[[6]]" {
			second = i + 1
		}
	}
	return first * second
}

func compareAny(left, right any) int {
	l, lok := left.(float64)
	r, rok := right.(float64)
	result := 0
	if lok && rok {
		//fmt.Printf("(int|int)\n")
		result = compareInts(int(l), int(r))
	} else if lok && !rok {
		//fmt.Printf("(int|list)\n")
		result = compareAny(makeArray(l), right)
	} else if !lok && rok {
		//fmt.Printf("(list|int)\n")
		result = compareAny(left, makeArray(r))
	} else {
		//fmt.Printf("(list|list)\n")
		result = compareArrays(left.([]interface{}), right.([]interface{}))
	}
	return result
}

func makeArray(val float64) []interface{} {
	i := make([]interface{}, 0)
	i = append(i, val)
	return i
}

func compareInts(left, right int) int {
	return right - left
}

func compareArrays(left, right []interface{}) int {
	for i := 0; i < max(len(left), len(right)); i++ {
		if len(left)-1 < i {
			//fmt.Println("Left Ran Out", 1)
			return 1
		} else if len(right)-1 < i {
			//fmt.Println("Right Ran Out", -1)
			return -1
		} else {
			if comparison := compareAny(left[i], right[i]); comparison != 0 {
				//fmt.Println("Comparison Result", comparison)
				return comparison
			}
		}
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
