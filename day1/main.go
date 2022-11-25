package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input_test")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		fmt.Println(s.Text())
	}
}
