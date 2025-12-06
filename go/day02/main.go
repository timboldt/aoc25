package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	start uint64
	end   uint64
}

func parseRanges(input string) []Range {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", "")
	parts := strings.Split(input, ",")

	var ranges []Range
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		rangeParts := strings.Split(part, "-")
		start, _ := strconv.ParseUint(rangeParts[0], 10, 64)
		end, _ := strconv.ParseUint(rangeParts[1], 10, 64)
		ranges = append(ranges, Range{start: start, end: end})
	}

	return ranges
}

func isRepeatedPattern(n uint64) bool {
	s := strconv.FormatUint(n, 10)
	length := len(s)

	// Must have even length
	if length%2 != 0 {
		return false
	}

	// Split in half
	mid := length / 2
	firstHalf := s[:mid]
	secondHalf := s[mid:]

	// Check no leading zeros
	if strings.HasPrefix(firstHalf, "0") {
		return false
	}

	// Check if halves are equal
	return firstHalf == secondHalf
}

func part1(input string) uint64 {
	ranges := parseRanges(input)
	var sum uint64

	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isRepeatedPattern(id) {
				sum += id
			}
		}
	}

	return sum
}

func isRepeatedPatternV2(n uint64) bool {
	s := strconv.FormatUint(n, 10)
	length := len(s)

	// Try all possible pattern lengths from 1 to len/2
	for patternLen := 1; patternLen <= length/2; patternLen++ {
		// Check if len is divisible by patternLen (pattern must repeat evenly)
		if length%patternLen == 0 {
			pattern := s[:patternLen]

			// Check for leading zeros
			if strings.HasPrefix(pattern, "0") {
				continue
			}

			// Check if all chunks match the pattern
			allMatch := true
			for i := 0; i < length; i += patternLen {
				if s[i:i+patternLen] != pattern {
					allMatch = false
					break
				}
			}

			if allMatch {
				return true
			}
		}
	}

	return false
}

func part2(input string) uint64 {
	ranges := parseRanges(input)
	var sum uint64

	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isRepeatedPatternV2(id) {
				sum += id
			}
		}
	}

	return sum
}

func main() {
	data, err := os.ReadFile("../../inputs/day02.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
