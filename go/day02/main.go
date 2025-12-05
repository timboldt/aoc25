package main

import (
	"fmt"
	"os"
)

func part1(input string) uint64 {
	// TODO: Implement part 1
	return 0
}

func part2(input string) uint64 {
	// TODO: Implement part 2
	return 0
}

func main() {
	data, err := os.ReadFile("../inputs/day02.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
