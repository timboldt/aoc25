package main

import "testing"

const example = `
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`

func TestPart1(t *testing.T) {
	result := part1(example)
	expected := 7
	if result != expected {
		t.Errorf("part1() = %d, want %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(example)
	expected := int64(33)
	if result != expected {
		t.Errorf("part2() = %d, want %d", result, expected)
	}
}
