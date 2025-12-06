package main

import "testing"

const example = `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

func TestPart1(t *testing.T) {
	result := part1(example)
	expected := 13
	if result != expected {
		t.Errorf("part1() = %d, expected %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(example)
	expected := 43
	if result != expected {
		t.Errorf("part2() = %d, expected %d", result, expected)
	}
}
