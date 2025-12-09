package main

import (
	"testing"
)

func TestPart1(t *testing.T) {
	input := `123 328  51 64
 45 64  387 23
  6 98  215 314
*   +   *   +  `

	expected := int64(4277556)
	result := parseAndSolve(input)

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
