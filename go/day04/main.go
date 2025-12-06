package main

import (
	"fmt"
	"os"
	"strings"
)

func parseInput(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

func countAdjacentPapers(grid [][]rune, row, col int) int {
	rows := len(grid)
	cols := len(grid[0])
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	count := 0
	for _, dir := range directions {
		newRow := row + dir[0]
		newCol := col + dir[1]

		if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols {
			if grid[newRow][newCol] == '@' {
				count++
			}
		}
	}
	return count
}

func part1(input string) int {
	grid := parseInput(input)
	accessible := 0

	for row, line := range grid {
		for col, cell := range line {
			if cell == '@' {
				adjacentCount := countAdjacentPapers(grid, row, col)
				if adjacentCount < 4 {
					accessible++
				}
			}
		}
	}

	return accessible
}

func part2(input string) int {
	grid := parseInput(input)
	totalRemoved := 0

	for {
		// Find all accessible rolls in this iteration
		toRemove := [][2]int{}

		for row, line := range grid {
			for col, cell := range line {
				if cell == '@' {
					adjacentCount := countAdjacentPapers(grid, row, col)
					if adjacentCount < 4 {
						toRemove = append(toRemove, [2]int{row, col})
					}
				}
			}
		}

		// If no rolls can be removed, we're done
		if len(toRemove) == 0 {
			break
		}

		// Remove all accessible rolls
		for _, pos := range toRemove {
			grid[pos[0]][pos[1]] = '.'
		}

		totalRemoved += len(toRemove)
	}

	return totalRemoved
}

func main() {
	data, err := os.ReadFile("../../inputs/day04.txt")
	if err != nil {
		panic(err)
	}
	input := string(data)

	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
