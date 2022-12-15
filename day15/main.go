package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type puzzle struct {
	sensors map[coordinate]sensor
	beacons map[coordinate]bool
}

type coordinate struct {
	x, y int
}

type intRange struct {
	min, max int
}

type sensor struct {
	location      coordinate
	closestBeacon coordinate
	curRange      int
}

func main() {
	//file, row, searchMax := "input_test", 10, 20
	file, row, searchMax := "input.txt", 2000000, 4000000

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

	// Parse the puzzle input
	p := puzzle{make(map[coordinate]sensor), make(map[coordinate]bool)}
	p.parse(lines)

	// Part 1 - The number of locations in the row that are covered by our scanners
	beaconFree := p.beaconFreePositions(row)
	fmt.Printf("There are %d locations that cannot contain beacons in row %d\n", beaconFree, row)

	// Part 2 - Find the beacon in the given range
	beaconLoc := p.findBeaconInArea(searchMax)
	fmt.Printf("There is a beacon at (%d,%d) with frequency %d", beaconLoc.x, beaconLoc.y, tuningFrequency(beaconLoc))
}

func (p *puzzle) parse(lines []string) {
	for _, line := range lines {
		sc, bc := coordinate{0, 0}, coordinate{0, 0}
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sc.x, &sc.y, &bc.x, &bc.y)
		s := NewSensor(sc, bc)
		p.sensors[s.location] = s
		p.beacons[s.closestBeacon] = true
	}
}

func (p puzzle) getCoveredRanges(row int) []intRange {
	ranges := make([]intRange, 0)
	for _, s := range p.sensors {
		if s.affectsRow(row) {
			affectedRange := s.getRangeInRow(row)
			ranges = append(ranges, affectedRange)
			//fmt.Printf("Sensor %v covers [%d, %d] in row %d\n", s, affectedRange.min, affectedRange.max, row)
		}
	}
	return combineRanges(ranges)
}

func (p puzzle) beaconFreePositions(row int) int {
	covered := 0
	for _, r := range p.getCoveredRanges(row) {
		covered += r.max - r.min + 1
	}

	// The row could have a beacon on it. If so, they need to be subtracted out
	beaconsOnRow := 0
	for knownBeacons := range p.beacons {
		if knownBeacons.y == row {
			beaconsOnRow++
		}
	}
	return covered - beaconsOnRow
}

func (p puzzle) findBeaconInArea(searchMax int) coordinate {
	found := coordinate{0, 0}
	for row := 0; row <= searchMax; row++ {
		for _, r := range p.getCoveredRanges(row) {
			if r.min <= 0 && r.max >= 0 && r.max < searchMax {
				found = coordinate{r.max + 1, row}
				fmt.Printf("Row %d has a range that ends at %d\n", row, r.max)
			}
		}
		//fmt.Printf("%02d: %v\n", row, p.getCoveredRanges(row))
	}
	return found
}

func NewSensor(location, closestBeacon coordinate) sensor {
	s := sensor{location, closestBeacon, manhattanDist(location, closestBeacon)}
	return s
}

func (s sensor) affectsRow(row int) bool {
	return absDiff(s.location.y, row) <= s.curRange
}

func (s sensor) getRangeInRow(row int) intRange {
	xDelta := s.curRange - absDiff(s.location.y, row)
	return intRange{s.location.x - xDelta, s.location.x + xDelta}
}

func (s sensor) String() string {
	return fmt.Sprintf("(%d,%d)[%d]", s.location.x, s.location.y, s.curRange)
}

func combineRanges(ranges []intRange) []intRange {
	combined := make([]intRange, 0)

	// Sort the ranges based on their minimum point
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].min < ranges[j].min
	})

	// Iterate over them, combining them if necessary
	open := ranges[0]
	for _, r := range ranges {
		if r.min <= open.max+1 {
			if r.max > open.max {
				open.max = r.max
			}
		} else {
			combined = append(combined, open)
			open = r
		}
	}
	combined = append(combined, open)
	return combined
}

func tuningFrequency(c coordinate) int {
	return c.x*4000000 + c.y
}

func manhattanDist(a, b coordinate) int {
	return absDiff(a.x, b.x) + absDiff(a.y, b.y)
}

func absDiff(x, y int) int {
	if x > y {
		return x - y
	}
	return y - x
}
