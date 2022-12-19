package main

import (
	"bufio"
	"fmt"
	"os"
)

type blueprint struct {
	id                                  int
	oreRobotOre                         int
	clayRobotOre                        int
	obsidianRobotOre, obsidianRobotClay int
	geodeRobotOre, geodeRobotObsidian   int
}

type puzzle struct {
	blueprints []blueprint
}

func main() {
	file := "input_test"
	//file := "input.txt"

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

	p := puzzle{make([]blueprint, 0)}
	p.parse(lines)

	p.part1()
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		bp := blueprint{0, 0, 0, 0, 0, 0, 0}
		fmt.Sscanf(line,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bp.id, &bp.oreRobotOre, &bp.clayRobotOre, &bp.obsidianRobotOre, &bp.obsidianRobotClay, &bp.geodeRobotOre, &bp.geodeRobotObsidian,
		)
		p.blueprints = append(p.blueprints, bp)
	}
}

func (p *puzzle) part1() {
	qlTotal := 0
	for _, bp := range p.blueprints {
		o := maxGeodes(bp, 1, 0, 0, 0, 0, 1, 0, 0, 0)
		ql := o * bp.id
		qlTotal += ql
		fmt.Printf("The best score for BP-%d was %d for a quality level of %d\n", bp.id, o, ql)
	}
	fmt.Printf("The total Quality Level was %d\n", qlTotal)
}

func maxGeodes(bp blueprint, turn int, o, c, ob, g int, oR, cR, obR, gR int) int {
	// If this is the last turn, return the number of geodes we have, plus the number that will be mined
	if turn == 24 {
		//fmt.Printf("Ore:%d(+%d) Clay:%d(+%d) Obsidian:%d(+%d) Geode:%d(+%d)\n",
		//	o, oR, c, cR, ob, obR, g, gR)
		return g + gR
	}

	// Each turn, we can potentially do a few things, so we basically recursively calculate each of these options
	// and maximize the amount of ore gained
	maxG := 0

	// 1. Do Nothing, let our ores accumulate
	maxG = max(maxG, maxGeodes(bp, turn+1, o+oR, c+cR, ob+obR, g+gR, oR, cR, obR, gR))

	// 2. Build an Ore Robot if we have enough ore
	if o >= bp.oreRobotOre && cR < 1 {
		maxG = max(maxG, maxGeodes(bp, turn+1, o+oR-bp.oreRobotOre, c+cR, ob+obR, g+gR, oR+1, cR, obR, gR))
	}

	// 3. Build a Clay Robot if we have enough ore
	if o >= bp.clayRobotOre && obR < 1 {
		maxG = max(maxG, maxGeodes(bp, turn+1, o+oR-bp.clayRobotOre, c+cR, ob+obR, g+gR, oR, cR+1, obR, gR))
	}

	// 4. Build an Obsidian Robot if we have enough ore and clay
	if o >= bp.obsidianRobotOre && c >= bp.obsidianRobotClay && gR < 1 {
		maxG = max(maxG, maxGeodes(bp, turn+1, o+oR-bp.obsidianRobotOre, c+cR-bp.obsidianRobotClay, ob+obR, g+gR, oR, cR, obR+1, gR))
	}

	//  5. Build a Geode Robot if we have enough ore and Obsidian
	if o >= bp.geodeRobotOre && ob >= bp.geodeRobotObsidian {
		maxG = max(maxG, maxGeodes(bp, turn+1, o+oR-bp.geodeRobotOre, c+cR, ob+obR-bp.geodeRobotObsidian, g+gR, oR, cR, obR, gR+1))
	}
	return maxG
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
