package main

import "fmt"

func debug() {
	example := `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

	polygon := parsePositions(example)
	redSet := make(map[Point]bool)
	for _, p := range polygon {
		redSet[p] = true
	}

	fmt.Println("Polygon vertices:", polygon)
	fmt.Println()

	// Test green ranges for y=3,4,5
	for y := 3; y <= 5; y++ {
		ranges := computeGreenRanges(y, polygon, redSet)
		fmt.Printf("Y=%d, Green ranges: %v\n", y, ranges)

		// Check if [2,9] is covered
		covered := isRangeCovered(2, 9, ranges)
		fmt.Printf("  Is [2,9] covered? %v\n", covered)
	}

	fmt.Println()
	result := part2(example)
	fmt.Printf("Part2 result: %d (expected 24)\n", result)
}
