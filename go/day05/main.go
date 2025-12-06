package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start uint64
	end   uint64
}

func (r Range) contains(id uint64) bool {
	return id >= r.start && id <= r.end
}

func parseInput(input string) ([]Range, []uint64) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")

	ranges := []Range{}
	for _, line := range strings.Split(parts[0], "\n") {
		parts := strings.Split(line, "-")
		start, _ := strconv.ParseUint(parts[0], 10, 64)
		end, _ := strconv.ParseUint(parts[1], 10, 64)
		ranges = append(ranges, Range{start: start, end: end})
	}

	ids := []uint64{}
	for _, line := range strings.Split(parts[1], "\n") {
		id, _ := strconv.ParseUint(line, 10, 64)
		ids = append(ids, id)
	}

	return ranges, ids
}

func isFresh(id uint64, ranges []Range) bool {
	for _, r := range ranges {
		if r.contains(id) {
			return true
		}
	}
	return false
}

func part1(input string) int {
	ranges, ids := parseInput(input)

	count := 0
	for _, id := range ids {
		if isFresh(id, ranges) {
			count++
		}
	}

	return count
}

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{}
	}

	// Sort by start position
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].start < sorted[j].start
	})

	merged := []Range{sorted[0]}

	for _, r := range sorted[1:] {
		last := &merged[len(merged)-1]

		// Check if ranges overlap or are adjacent
		if r.start <= last.end+1 {
			// Merge by extending the end
			if r.end > last.end {
				last.end = r.end
			}
		} else {
			// No overlap, add as new range
			merged = append(merged, r)
		}
	}

	return merged
}

func countIDsInRanges(ranges []Range) uint64 {
	merged := mergeRanges(ranges)
	total := uint64(0)
	for _, r := range merged {
		total += r.end - r.start + 1
	}
	return total
}

func part2(input string) uint64 {
	ranges, _ := parseInput(input)
	return countIDsInRanges(ranges)
}

func main() {
	data, err := os.ReadFile("../../inputs/day05.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
