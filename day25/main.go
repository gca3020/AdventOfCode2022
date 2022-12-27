package main

import (
	"bufio"
	"fmt"
	"os"
)

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

	part1(lines)
}

func part1(lines []string) {
	sum := 0
	for _, line := range lines {
		dec := toDecimal(line)
		sum += dec
		fmt.Printf("%s -> %d\n", line, toDecimal(line))
	}

	fmt.Printf("The total fuel required is %d, which is %s in SNAFU format", sum, toSnafu(sum))
}

func toDecimal(snafu string) int {
	mult := 1
	dec := 0
	for i := len(snafu) - 1; i >= 0; i-- {
		switch snafu[i] {
		case '=':
			dec += mult * -2
		case '-':
			dec += mult * -1
		case '0':
			dec += 0
		case '1':
			dec += 1 * mult
		case '2':
			dec += 2 * mult
		}
		mult *= 5
	}
	return dec
}

func toSnafu(decimal int) string {
	div := 5
	snafu := ""
	for decimal > 0 {
		remainder := decimal % div
		switch remainder {
		case 0:
			snafu = "0" + snafu
		case 1:
			snafu = "1" + snafu
		case 2:
			snafu = "2" + snafu
		case 3:
			snafu = "=" + snafu
			decimal += 5
		case 4:
			snafu = "-" + snafu
			decimal += 5
		}
		decimal /= div
	}
	return snafu
}
