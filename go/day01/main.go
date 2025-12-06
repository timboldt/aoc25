package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	Left Direction = iota
	Right
)

type Rotation struct {
	direction Direction
	distance  int
}

func parseRotations(input string) []Rotation {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rotations := make([]Rotation, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var direction Direction
		if line[0] == 'L' {
			direction = Left
		} else if line[0] == 'R' {
			direction = Right
		} else {
			panic(fmt.Sprintf("Unknown direction: %c", line[0]))
		}

		distance, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(fmt.Sprintf("Failed to parse distance: %v", err))
		}

		rotations = append(rotations, Rotation{
			direction: direction,
			distance:  distance,
		})
	}

	return rotations
}

// mod implements modulo operation that handles negative numbers correctly
func mod(a, b int) int {
	return ((a % b) + b) % b
}

func part1(input string) int64 {
	rotations := parseRotations(input)
	position := 50
	count := int64(0)

	for _, rotation := range rotations {
		switch rotation.direction {
		case Left:
			position = mod(position-rotation.distance, 100)
		case Right:
			position = mod(position+rotation.distance, 100)
		}

		if position == 0 {
			count++
		}
	}

	return count
}

func part2(input string) int64 {
	rotations := parseRotations(input)
	position := 50
	count := int64(0)

	for _, rotation := range rotations {
		// Count how many times we pass through 0 during this rotation
		var clicksThroughZero int64

		switch rotation.direction {
		case Right:
			// Going right from position P by distance D
			// We cross 0 at positions where (P + k) % 100 == 0 for k in [1, D]
			// This happens floor((P + D) / 100) - floor(P / 100) times
			clicksThroughZero = int64((position+rotation.distance)/100 - position/100)

		case Left:
			// Going left from position P by distance D
			// We cross 0 at positions where (P - k) % 100 == 0 for k in [1, D]
			// This happens when k = P, P+100, P+200, ... (all ≤ D)
			if position == 0 {
				// Starting at 0, we cross it again at k = 100, 200, ...
				clicksThroughZero = int64(rotation.distance / 100)
			} else if rotation.distance >= position {
				// We cross 0 at k = P, then every 100 clicks after
				clicksThroughZero = int64((rotation.distance-position)/100 + 1)
			} else {
				// We don't reach 0
				clicksThroughZero = 0
			}
		}

		count += clicksThroughZero

		// Update position
		switch rotation.direction {
		case Left:
			position = mod(position-rotation.distance, 100)
		case Right:
			position = mod(position+rotation.distance, 100)
		}
	}

	return count
}

func main() {
	data, err := os.ReadFile("../../inputs/day01.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
