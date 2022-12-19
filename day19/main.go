package main

import (
	"bufio"
	"fmt"
	"os"
)

type blueprint struct {
	id int

	oreRobotOre                         int
	clayRobotOre                        int
	obsidianRobotOre, obsidianRobotClay int
	geodeRobotOre, geodeRobotObsidian   int

	oRMax, cRMax, obRMax int
}

type puzzle struct {
	blueprints []blueprint
	totalTurns int
	calls      int64
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

	p := puzzle{make([]blueprint, 0), 0, 0}
	p.parse(lines)

	p.part1()
	p.part2()
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		bp := blueprint{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		fmt.Sscanf(line,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bp.id, &bp.oreRobotOre, &bp.clayRobotOre, &bp.obsidianRobotOre, &bp.obsidianRobotClay, &bp.geodeRobotOre, &bp.geodeRobotObsidian,
		)
		bp.oRMax = max(max(bp.oreRobotOre, bp.clayRobotOre), max(bp.obsidianRobotOre, bp.geodeRobotOre))
		bp.cRMax = bp.obsidianRobotClay
		bp.obRMax = bp.geodeRobotObsidian
		p.blueprints = append(p.blueprints, bp)
	}
}

func (p *puzzle) part1() {
	qlTotal := 0
	for _, bp := range p.blueprints {
		p.calls = 0
		p.totalTurns = 24
		o := p.maxGeodes(bp, 1, 0, 0, 0, 0, 1, 0, 0, 0, false)
		ql := o * bp.id
		qlTotal += ql
		fmt.Printf("The best score for BP-%d was %d for a quality level of %d (%d calls)\n", bp.id, o, ql, p.calls)
	}
	fmt.Printf("The total Quality Level was %d\n", qlTotal)
}

func (p *puzzle) part2() {
	product := 1
	for _, bp := range p.blueprints {
		if bp.id > 3 {
			continue
		}
		p.totalTurns = 32
		o := p.maxGeodes(bp, 1, 0, 0, 0, 0, 1, 0, 0, 0, false)
		product *= o
		fmt.Printf("The best score for BP-%d was %d (%d calls)\n", bp.id, o, p.calls)
	}
	fmt.Printf("The product of max geodes was %d\n", product)
}

func (p *puzzle) maxGeodes(bp blueprint, turn int, o, c, ob, g int, oR, cR, obR, gR int, idle bool) int {
	// If this is the last turn, return the number of geodes we have, plus the number that will be mined
	p.calls++
	if turn == p.totalTurns {
		//fmt.Printf("Ore:%d(+%d) Clay:%d(+%d) Obsidian:%d(+%d) Geode:%d(+%d)\n",
		//	o, oR, c, cR, ob, obR, g, gR)
		return g + gR
	}

	// Each turn, we can potentially do a few things, so we basically recursively calculate each of these options
	// and maximize the amount of ore gained
	maxG := 0

	// If we have enough to build a geode robot, then that is always our best move, don't bother branching to others
	if o >= bp.geodeRobotOre && ob >= bp.geodeRobotObsidian {
		return p.maxGeodes(bp, turn+1, o+oR-bp.geodeRobotOre, c+cR, ob+obR-bp.geodeRobotObsidian, g+gR, oR, cR, obR, gR+1, false)
	}

	// Build an Ore Robot if:
	//  - we have enough ore,
	//  - we have below the max useful number of ore robots
	//  - we didn't have enough to build one last turn OR we did, but built something else instead
	if o >= bp.oreRobotOre && oR < bp.oRMax && (o-oR < bp.oreRobotOre || !idle) {
		maxG = max(maxG, p.maxGeodes(bp, turn+1, o+oR-bp.oreRobotOre, c+cR, ob+obR, g+gR, oR+1, cR, obR, gR, false))
	}

	// Build a Clay Robot if:
	//  - we have enough ore
	//  - we're below our max of useful clay robots, and below our max of useful obsidian robots
	//  - we didn't have enough to build one last turn OR we did, but built something else instead
	if o >= bp.clayRobotOre && cR < bp.cRMax && obR < bp.obRMax && (o-oR < bp.clayRobotOre || !idle) {
		maxG = max(maxG, p.maxGeodes(bp, turn+1, o+oR-bp.clayRobotOre, c+cR, ob+obR, g+gR, oR, cR+1, obR, gR, false))
	}

	// Build an Obsidian Robot if:
	//  - we have enough ore and clay
	//  - we're below our useful max of obsidian robots
	//  - we didn't have enough to build one last turn OR we did, but built something else instead
	if o >= bp.obsidianRobotOre && c >= bp.obsidianRobotClay && obR < bp.obRMax && (o-oR < bp.obsidianRobotOre || c-cR < bp.obsidianRobotClay || !idle) {
		maxG = max(maxG, p.maxGeodes(bp, turn+1, o+oR-bp.obsidianRobotOre, c+cR-bp.obsidianRobotClay, ob+obR, g+gR, oR, cR, obR+1, gR, false))
	}

	// Stay idle and let our ores accumulate
	maxG = max(maxG, p.maxGeodes(bp, turn+1, o+oR, c+cR, ob+obR, g+gR, oR, cR, obR, gR, true))

	return maxG
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
