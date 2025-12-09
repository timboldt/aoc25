package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Problem struct {
	numbers   []int64
	operation string
}

func parseAndSolve(input string) int64 {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return 0
	}

	// Find max line length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Parse character by character into columns
	columns := make([][]rune, maxLen)
	for i := range columns {
		columns[i] = make([]rune, 0)
	}

	for _, line := range lines {
		runes := []rune(line)
		for colIdx, ch := range runes {
			columns[colIdx] = append(columns[colIdx], ch)
		}
		// Pad shorter lines with spaces
		for colIdx := len(runes); colIdx < maxLen; colIdx++ {
			columns[colIdx] = append(columns[colIdx], ' ')
		}
	}

	// Group consecutive non-space columns into problems
	problems := make([][]int, 0)
	currentProblem := make([]int, 0)

	for colIdx, column := range columns {
		// Check if this is a separator column (all spaces)
		allSpaces := true
		for _, ch := range column {
			if ch != ' ' {
				allSpaces = false
				break
			}
		}

		if allSpaces {
			if len(currentProblem) > 0 {
				problems = append(problems, append([]int(nil), currentProblem...))
				currentProblem = currentProblem[:0]
			}
		} else {
			currentProblem = append(currentProblem, colIdx)
		}
	}

	// Don't forget the last problem
	if len(currentProblem) > 0 {
		problems = append(problems, currentProblem)
	}

	// Solve each problem
	grandTotal := int64(0)
	for _, problemCols := range problems {
		if result := solveProblem(columns, problemCols); result != nil {
			grandTotal += *result
		}
	}

	return grandTotal
}

func solveProblem(columns [][]rune, problemCols []int) *int64 {
	if len(problemCols) == 0 {
		return nil
	}

	numRows := len(columns[problemCols[0]])

	// Extract the operation from the last row
	operationStr := ""
	for _, colIdx := range problemCols {
		if colIdx < len(columns) && numRows > 0 {
			operationStr += string(columns[colIdx][numRows-1])
		}
	}
	operation := strings.TrimSpace(operationStr)

	if operation != "*" && operation != "+" {
		return nil
	}

	// Extract numbers from rows above the operation
	numbers := make([]int64, 0)
	for rowIdx := 0; rowIdx < numRows-1; rowIdx++ {
		rowStr := ""
		for _, colIdx := range problemCols {
			if colIdx < len(columns) && rowIdx < len(columns[colIdx]) {
				rowStr += string(columns[colIdx][rowIdx])
			}
		}
		trimmed := strings.TrimSpace(rowStr)
		if trimmed != "" {
			if num, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
				numbers = append(numbers, num)
			}
		}
	}

	if len(numbers) == 0 {
		return nil
	}

	// Calculate result
	var result int64
	if operation == "*" {
		result = 1
		for _, num := range numbers {
			result *= num
		}
	} else {
		result = 0
		for _, num := range numbers {
			result += num
		}
	}

	return &result
}

func main() {
	data, err := os.ReadFile("../../inputs/day06.txt")
	if err != nil {
		panic(err)
	}

	input := string(data)
	fmt.Printf("Part 1: %d\n", parseAndSolve(input))
}
