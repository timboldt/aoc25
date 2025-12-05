package main

import "testing"

const example = `987654321111111
811111111111119
234234234234278
818181911112111
`

func TestPart1(t *testing.T) {
	result := part1(example)
	expected := uint32(357)
	if result != expected {
		t.Errorf("part1() = %d, want %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(example)
	expected := uint64(3121910778619)
	if result != expected {
		t.Errorf("part2() = %d, want %d", result, expected)
	}
}
