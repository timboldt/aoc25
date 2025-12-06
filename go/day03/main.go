package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseBanks(input string) []string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var banks []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			banks = append(banks, line)
		}
	}
	return banks
}

func maxJoltage(bank string) uint32 {
	digits := []rune(bank)
	var max uint32 = 0

	// Try all pairs (i, j) where i < j
	for i := 0; i < len(digits); i++ {
		for j := i + 1; j < len(digits); j++ {
			// Form two-digit number from digits[i] and digits[j]
			joltageStr := fmt.Sprintf("%c%c", digits[i], digits[j])
			joltage, _ := strconv.ParseUint(joltageStr, 10, 32)
			if uint32(joltage) > max {
				max = uint32(joltage)
			}
		}
	}

	return max
}

func part1(input string) uint32 {
	banks := parseBanks(input)
	var sum uint32 = 0

	for _, bank := range banks {
		sum += maxJoltage(bank)
	}

	return sum
}

func maxJoltageN(bank string, n int) uint64 {
	digits := []rune(bank)
	toRemove := len(digits) - n

	var stack []rune
	removalsLeft := toRemove

	// Greedy approach: remove smaller digits when a larger digit comes
	for _, digit := range digits {
		for len(stack) > 0 && removalsLeft > 0 && stack[len(stack)-1] < digit {
			stack = stack[:len(stack)-1]
			removalsLeft--
		}
		stack = append(stack, digit)
	}

	// If we still have removals left, remove from the end
	for removalsLeft > 0 {
		stack = stack[:len(stack)-1]
		removalsLeft--
	}

	// Take first n digits and form the number
	result := string(stack[:n])
	num, _ := strconv.ParseUint(result, 10, 64)
	return num
}

func part2(input string) uint64 {
	banks := parseBanks(input)
	var sum uint64 = 0

	for _, bank := range banks {
		sum += maxJoltageN(bank, 12)
	}

	return sum
}

func main() {
	data, err := os.ReadFile("../../inputs/day03.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
