package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func parsePositions(input string) []Point {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	positions := make([]Point, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic(fmt.Sprintf("Invalid line: %s", line))
		}

		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse x: %v", err))
		}

		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse y: %v", err))
		}

		positions = append(positions, Point{x, y})
	}

	return positions
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func part1(input string) int64 {
	positions := parsePositions(input)
	n := len(positions)

	maxArea := int64(0)

	// Try all pairs of red tiles as opposite corners
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := positions[i]
			p2 := positions[j]

			// Calculate rectangle dimensions (inclusive of corners)
			width := int64(abs(p2.x-p1.x) + 1)
			height := int64(abs(p2.y-p1.y) + 1)
			area := width * height

			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

type Range struct {
	min, max int
}

// isPointOnSegment checks if point (px, py) is on the line segment from (x1,y1) to (x2,y2)
func isPointOnSegment(px, py, x1, y1, x2, y2 int) bool {
	// Check if point is on horizontal segment
	if y1 == y2 && py == y1 {
		return (px >= min(x1, x2)) && (px <= max(x1, x2))
	}
	// Check if point is on vertical segment
	if x1 == x2 && px == x1 {
		return (py >= min(y1, y2)) && (py <= max(y1, y2))
	}
	return false
}

// isPointInPolygon uses ray casting algorithm to check if a point is inside a polygon
func isPointInPolygon(px, py int, polygon []Point) bool {
	n := len(polygon)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := polygon[i].x, polygon[i].y
		xj, yj := polygon[j].x, polygon[j].y

		// Check if ray crosses this edge
		if ((yi > py) != (yj > py)) && (px < (xj-xi)*(py-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// isGreenTile checks if a tile is green (on edge or inside polygon)
func isGreenTile(p Point, polygon []Point) bool {
	n := len(polygon)

	// Check if on any edge
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		p1 := polygon[i]
		p2 := polygon[next]

		if isPointOnSegment(p.x, p.y, p1.x, p1.y, p2.x, p2.y) {
			return true
		}
	}

	// Check if inside polygon
	return isPointInPolygon(p.x, p.y, polygon)
}

// computeGreenRanges computes which x-ranges are green (red or on edge or inside) for a given y
func computeGreenRanges(y int, polygon []Point, redSet map[Point]bool) []Range {
	n := len(polygon)
	var ranges []Range

	// Find horizontal edges at this y
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		p1 := polygon[i]
		p2 := polygon[next]

		if p1.y == y && p2.y == y {
			minX := min(p1.x, p2.x)
			maxX := max(p1.x, p2.x)
			ranges = append(ranges, Range{minX, maxX})
		}
	}

	// Find vertical edges that cross this y and use scanline algorithm
	var crossings []int
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		p1 := polygon[i]
		p2 := polygon[next]

		if p1.x == p2.x {
			minY := min(p1.y, p2.y)
			maxY := max(p1.y, p2.y)
			// Use half-open interval [minY, maxY) to avoid double-counting at vertices
			if y >= minY && y < maxY {
				crossings = append(crossings, p1.x)
			}
		}
	}

	sort.Ints(crossings)

	// Remove duplicates
	if len(crossings) > 0 {
		unique := []int{crossings[0]}
		for j := 1; j < len(crossings); j++ {
			if crossings[j] != unique[len(unique)-1] {
				unique = append(unique, crossings[j])
			}
		}
		crossings = unique
	}

	// Pair up crossings to get inside ranges
	for j := 0; j < len(crossings); j += 2 {
		if j+1 < len(crossings) {
			ranges = append(ranges, Range{crossings[j], crossings[j+1]})
		}
	}

	// Merge overlapping ranges and return
	if len(ranges) == 0 {
		return ranges
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].min < ranges[j].min
	})

	merged := []Range{ranges[0]}
	for i := 1; i < len(ranges); i++ {
		last := &merged[len(merged)-1]
		if ranges[i].min <= last.max+1 {
			last.max = max(last.max, ranges[i].max)
		} else {
			merged = append(merged, ranges[i])
		}
	}

	return merged
}

// isRangeCovered checks if [xMin, xMax] is fully covered by the given ranges
func isRangeCovered(xMin, xMax int, ranges []Range) bool {
	for _, r := range ranges {
		if xMin >= r.min && xMax <= r.max {
			return true
		}
	}
	return false
}

func part2(input string) int64 {
	redTiles := parsePositions(input)
	n := len(redTiles)

	// Create a set of red tiles for quick lookup
	redSet := make(map[Point]bool)
	for _, p := range redTiles {
		redSet[p] = true
	}

	// Build list of all candidate rectangles with their areas
	type Candidate struct {
		area       int64
		minX, maxX int
		minY, maxY int
	}

	var candidates []Candidate
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := redTiles[i]
			p2 := redTiles[j]

			minX := min(p1.x, p2.x)
			maxX := max(p1.x, p2.x)
			minY := min(p1.y, p2.y)
			maxY := max(p1.y, p2.y)

			width := int64(maxX - minX + 1)
			height := int64(maxY - minY + 1)
			area := width * height

			candidates = append(candidates, Candidate{area, minX, maxX, minY, maxY})
		}
	}

	// Sort by area descending
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].area > candidates[j].area
	})

	// Cache for green ranges by y-coordinate
	greenRangeCache := make(map[int][]Range)

	// Check rectangles from largest to smallest
	for _, c := range candidates {
		// Check if all rows in the rectangle have the x-range fully covered
		valid := true
		for y := c.minY; y <= c.maxY && valid; y++ {
			// Get or compute green ranges for this y
			ranges, exists := greenRangeCache[y]
			if !exists {
				ranges = computeGreenRanges(y, redTiles, redSet)
				greenRangeCache[y] = ranges
			}

			if !isRangeCovered(c.minX, c.maxX, ranges) {
				valid = false
			}
		}

		if valid {
			return c.area
		}
	}

	return 0
}

func main() {
	data, err := os.ReadFile("../../inputs/day09.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
