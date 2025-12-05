package main

import "testing"

const example = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

func TestPart1(t *testing.T) {
	result := part1(example)
	expected := int64(3)
	if result != expected {
		t.Errorf("part1() = %d, want %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(example)
	expected := int64(6)
	if result != expected {
		t.Errorf("part2() = %d, want %d", result, expected)
	}
}
