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

type pattern struct {
	rockShape int
	jetIdx    int
}
type cycleState struct {
	height, rockCount int
}

type rock []coordinate

type puzzle struct {
	jet           string
	jetIdx        int
	chamber       map[coordinate]bool
	rockCount     int64
	cycleDetector map[pattern]cycleState
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

	p1 := NewPuzzle(lines[0])
	p1.simulate(2022, false)

	p1c := NewPuzzle(lines[0])
	p1c.simulate(2022, true)

	p2 := NewPuzzle(lines[0])
	p2.simulate(1000000000000, true)
}

func duration(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("  -- %s took %s\n\n", name, elapsed)
}

func NewPuzzle(line string) puzzle {
	return puzzle{
		jet:           line,
		jetIdx:        0,
		chamber:       make(map[coordinate]bool),
		rockCount:     0,
		cycleDetector: make(map[pattern]cycleState),
	}
}

func (p *puzzle) simulate(totalRocks int64, useCycleDetection bool) {
	defer duration(time.Now(), fmt.Sprintf("simulate %d (%v)", totalRocks, useCycleDetection))
	simulatedHeight := int64(0)
	simulatedRocks := int64(0)
	for p.rockCount+simulatedRocks < totalRocks {
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

		// Optional cycle detection
		if useCycleDetection {
			// Cycle detection. Don't start until after we've looped through the jet stream at least once to reach a steady state
			if p.jetIdx/len(p.jet) > 0 && simulatedHeight == 0 {
				detected, cycleHeight, cycleRocks := p.detectCycle(p.getHeight())
				if detected {
					simulatedRepeats := (totalRocks - p.rockCount) / cycleRocks
					simulatedHeight = simulatedRepeats * cycleHeight
					simulatedRocks = simulatedRepeats * cycleRocks
					fmt.Printf("Found a cycle of %d rocks adding %d height. Simulating %d repeats\n", cycleRocks, cycleHeight, simulatedRepeats)
				}
			}
		}
	}
	fmt.Printf("The tower height is %d (%d sim) after %d rocks (%d sim)\n",
		p.getHeight()+int(simulatedHeight), simulatedHeight,
		p.rockCount+simulatedRocks, simulatedRocks,
	)
}

func (p *puzzle) detectCycle(height int) (bool, int64, int64) {
	patt := pattern{int(p.rockCount % 5), p.jetIdx % len(p.jet)}
	new := cycleState{height, int(p.rockCount)}
	if val, ok := p.cycleDetector[patt]; ok {
		return true, int64(new.height - val.height), int64(new.rockCount - val.rockCount)
	}
	p.cycleDetector[patt] = new
	return false, 0, 0
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
