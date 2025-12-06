package main

import "testing"

const example = `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

func TestPart1(t *testing.T) {
	result := part1(example)
	expected := 3
	if result != expected {
		t.Errorf("part1() = %d, expected %d", result, expected)
	}
}

func TestPart2(t *testing.T) {
	result := part2(example)
	expected := uint64(14)
	if result != expected {
		t.Errorf("part2() = %d, expected %d", result, expected)
	}
}
