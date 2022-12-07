package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type day7 struct {
	s       *bufio.Scanner
	rootDir *Directory
	curDir  *Directory
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	d := day7{s, NewDirectory("/", nil), nil}
	d.parse()

	// Part 1
	totalSize := 0
	smallDirs := findDirs(make([]*Directory, 0), d.rootDir, func(d *Directory) bool { return d.Size() < 100000 })
	for _, smallDir := range smallDirs {
		totalSize += smallDir.Size()
	}
	fmt.Println("Small Directory Total Size", totalSize)

	// Part 2
	toDelete := 30000000 - (70000000 - d.rootDir.Size())
	largeDirs := findDirs(make([]*Directory, 0), d.rootDir, func(d *Directory) bool { return d.Size() > toDelete })
	smallest := d.rootDir
	for _, largeDir := range largeDirs {
		if largeDir.Size() < smallest.Size() {
			smallest = largeDir
		}
	}
	fmt.Println("The smallest directory to delete is", smallest.name, "with size", smallest.Size())
}

func (d *day7) parse() {
	for d.s.Scan() {
		fields := strings.Fields(d.s.Text())
		if fields[0] == "$" && fields[1] == "cd" {
			d.changeDirectory(fields[2])
		} else if fields[0] == "$" && fields[1] == "ls" {
			//fmt.Println("Listing contents of", d.curDir.Name())
		} else if len(fields) == 2 {
			if fields[0] == "dir" {
				d.curDir.AddSubdir(fields[1])
			} else {
				size, _ := strconv.Atoi(fields[0])
				d.curDir.AddFile(fields[1], size)
			}
		} else {
			fmt.Printf("Unhandled Command: %v\n", fields)
		}
	}
}

func (d *day7) changeDirectory(dir string) {
	if dir == "/" {
		d.curDir = d.rootDir
		return
	}
	if dir == ".." {
		d.curDir = d.curDir.parent
		return
	}
	if c, ok := d.curDir.subdirs[dir]; ok {
		d.curDir = c
		return
	}
	fmt.Println("Could not Change Directory", dir)
}

func findDirs(dirs []*Directory, currentDir *Directory, predicate func(*Directory) bool) []*Directory {
	if predicate(currentDir) {
		dirs = append(dirs, currentDir)
	}
	for _, d := range currentDir.subdirs {
		dirs = findDirs(dirs, d, predicate)
	}
	return dirs
}
