package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type list struct {
	orig []*node
	head *node
	tail *node
}

type node struct {
	next  *node
	prev  *node
	value int64
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

func part1(l *list) {
	defer duration(time.Now(), "Part 1")
	// Mix the list once
	l.mix()
	fmt.Printf("Part 1 - The code is %d\n", l.findCode())
}

func part2(l *list) {
	defer duration(time.Now(), "Part 2")

	// Multiply all numbers in the list by the decryption key
	decryptionKey := int64(811589153)
	for _, n := range l.orig {
		n.value *= decryptionKey
	}

	// Mix 10 times, then print the code
	for i := 0; i < 10; i++ {
		l.mix()
	}

	fmt.Printf("Part 2 - The code is %d\n", l.findCode())
}

func (l *list) findCode() int64 {
	sum := int64(0)
	n := l.find(0)
	for i := 1; i <= 3000; i++ {
		n = n.next
		if i%1000 == 0 {
			sum += n.value
			fmt.Printf("%d: %d\n", i, n.value)
		}
	}
	return sum
}

func parse(lines []string) *list {
	l := &list{make([]*node, 0), nil, nil}
	for _, line := range lines {
		if val, err := strconv.Atoi(line); err == nil {
			l.pushBack(&node{nil, nil, int64(val)})
		} else {
			fmt.Println("Parser Error:", err)
		}
	}

	// After we're done parsing, link the head and tail to make the list circular
	l.tail.next = l.head
	l.head.prev = l.tail
	return l
}

func (l *list) pushBack(n *node) {
	if len(l.orig) == 0 {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		n.prev = l.tail
		l.tail = n
	}
	l.orig = append(l.orig, n)
}

func (n *node) move(l *list, pos int64) {
	// We mod by one smaller than the total size of the list, since this item has technically been removed already
	modPos := int(pos % int64(len(l.orig)-1))
	// Shortcut if this element is not being moved. Do nothing
	if modPos == 0 {
		return
	}

	// First, pull this node out of the list and fill the gap
	n.prev.next, n.next.prev = n.next, n.prev

	// Then find the new location for this node
	newPrev, newNext := n.newLocation(modPos)

	// Then insert this element back into the list
	newPrev.next, newNext.prev = n, n
	n.prev, n.next = newPrev, newNext
}

func (n *node) newLocation(pos int) (newPrev, newNext *node) {
	//fmt.Printf("Finding new location %d cells away\n", pos)
	newPrev, newNext = n.prev, n.next
	if pos > 0 {
		for i := 0; i < pos; i++ {
			newPrev, newNext = newPrev.next, newNext.next
		}
	} else if pos < 0 {
		for i := 0; i > pos; i-- {
			newPrev, newNext = newPrev.prev, newNext.prev
		}
	}
	return newPrev, newNext
}

func (l *list) mix() {
	fmt.Println("Mixing List")
	for _, n := range l.orig {
		n.move(l, n.value)
	}
}

func (l *list) find(val int64) *node {
	for _, n := range l.orig {
		if n.value == val {
			return n
		}
	}
	return nil
}
