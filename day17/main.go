package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type coordinate struct {
	x, y int
}

type rock []coordinate

type puzzle struct {
	jet          string
	jetIdx       int
	chamber      map[coordinate]bool
	rockCount    int64
	repeatHeight int64
	repeatCount  int64
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

	p := puzzle{
		jet:          lines[0],
		jetIdx:       0,
		chamber:      make(map[coordinate]bool),
		rockCount:    0,
		repeatHeight: 0,
		repeatCount:  0,
	}
	fmt.Printf("The jet has %d entries\n", len(p.jet))
	//p.part1()
	p.part2()
}

func duration(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func (p *puzzle) part1() {
	defer duration(time.Now(), "part1")
	for p.rockCount < 2022 {
		// Add a rock and simulate its fall
		r, ok := p.newRock(), true
		for {
			r = p.pushRock(r)
			r, ok = p.fallRock(r)
			if !ok {
				break
			}
		}

		// Add the resting place of the rock to the chamber
		for _, coord := range r {
			p.chamber[coord] = true
		}

		// Increment the total number of rocks dropped
		p.rockCount++
	}
	fmt.Printf("The tower height is %d after %d rocks\n", p.getHeight(), p.rockCount)
}

func (p *puzzle) part2() {
	simulatedHeight := int64(0)
	simulatedRocks := int64(0)
	rocksToSimulate := int64(1000000000000)
	h5, r5, h10, r10 := int64(0), int64(0), int64(0), int64(0)
	for p.rockCount+simulatedRocks < rocksToSimulate {
		// Add a rock and simulate its fall
		r, ok := p.newRock(), true
		for {
			r = p.pushRock(r)
			r, ok = p.fallRock(r)

			if p.jetIdx%len(p.jet) == 0 && p.jetIdx/len(p.jet) == 5 {
				h5, r5 = int64(p.getHeight()), p.rockCount
				fmt.Printf("After 5 repeats, the height is %d and the rock count is %d\n", h5, r5)
			}
			if p.jetIdx%len(p.jet) == 0 && p.jetIdx/len(p.jet) == 10 {
				p.repeatHeight = int64(p.getHeight()) - h5
				p.repeatCount = int64(p.rockCount) - r5
				fmt.Printf("After 10 repeats, the height is %d (%d) and the rock count is %d (%d)\n", h10, p.repeatHeight, r10, p.repeatCount)
				simulatedRepeats := (rocksToSimulate - p.rockCount) / p.repeatCount
				simulatedHeight = simulatedRepeats * p.repeatHeight
				simulatedRocks = simulatedRepeats * p.repeatCount
				fmt.Printf("Simulating %d repeats adds %d extra height and %d extra rocks\n", simulatedRepeats, simulatedHeight, simulatedRocks)
			}

			if !ok {
				break
			}
		}

		// Add the resting place of the rock to the chamber
		for _, coord := range r {
			p.chamber[coord] = true
		}

		// Increment the total number of rocks dropped
		p.rockCount++
	}
	fmt.Printf("The tower height is %d (%d simulated) after %d rocks (%d simulated)\n",
		int64(p.getHeight())+simulatedHeight, simulatedHeight,
		p.rockCount+simulatedRocks, simulatedRocks,
	)
}

func (p *puzzle) newRock() rock {
	var r rock

	h := p.getHeight() + 4
	switch p.rockCount % 5 {
	case 0:
		//fmt.Printf("Adding a new h-line\n")
		r = rock{coordinate{2, h}, coordinate{3, h}, coordinate{4, h}, coordinate{5, h}}
	case 1:
		//fmt.Printf("Adding a new plus\n")
		r = rock{coordinate{2, h + 1}, coordinate{3, h + 2}, coordinate{3, h + 1}, coordinate{3, h}, coordinate{4, h + 1}}
	case 2:
		//fmt.Printf("Adding a new L\n")
		r = rock{coordinate{2, h}, coordinate{3, h}, coordinate{4, h}, coordinate{4, h + 1}, coordinate{4, h + 2}}
	case 3:
		//fmt.Printf("Adding a new v-line\n")
		r = rock{coordinate{2, h}, coordinate{2, h + 1}, coordinate{2, h + 2}, coordinate{2, h + 3}}
	case 4:
		//fmt.Printf("Adding a new square\n")
		r = rock{coordinate{2, h}, coordinate{2, h + 1}, coordinate{3, h}, coordinate{3, h + 1}}
	}
	return r
}

func (p *puzzle) pushRock(r rock) rock {
	// Determine the direction to move, and increment the tracking index
	direction := p.jet[p.jetIdx%len(p.jet)]
	p.jetIdx++

	// Determine L/R shift based on jet direction
	offset := -1
	if direction == '>' {
		offset = 1
	}

	if r1, ok := p.translateRock(r, offset, 0); ok {
		return r1
	}
	return r
}

func (p *puzzle) fallRock(r rock) (rock, bool) {
	if r1, ok := p.translateRock(r, 0, -1); ok {
		return r1, true
	}
	return r, false
}

func (p puzzle) translateRock(r rock, x, y int) (rock, bool) {
	r1 := make(rock, 0, len(r))
	for _, c := range r {
		c.x, c.y = c.x+x, c.y+y
		r1 = append(r1, c)
		// Check for any out-of-bounds conditions
		if c.x < 0 || c.x >= 7 || c.y <= 0 {
			return nil, false
		}
		// Check for overlaps
		if _, ok := p.chamber[c]; ok {
			return nil, false
		}
	}

	return r1, true
}

func (p puzzle) getHeight() int {
	maxY := 0
	for coord := range p.chamber {
		if coord.y > maxY {
			maxY = coord.y
		}
	}
	return maxY
}

func (p puzzle) draw() {
	for y := p.getHeight() + 3; y > 0; y-- {
		fmt.Print("|")
		for x := 0; x < 7; x++ {
			if _, ok := p.chamber[coordinate{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Println("+-------+")
}
