package main

import "testing"

const example = `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

func TestPart1(t *testing.T) {
	result := part1(parsePositions(example))
	expected := int64(50)
	if result != expected {
		t.Errorf("part1() = %d, want %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(parsePositions(example))
	expected := int64(24)
	if result != expected {
		t.Errorf("part2() = %d, want %d", result, expected)
	}
}
