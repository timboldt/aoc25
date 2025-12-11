package main

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	input := `aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out`

	graph := parseString(input)
	expected := 5
	if result := part1(graph); result != expected {
		t.Errorf("Part 1 expected %d, got %d", expected, result)
	}
}

func TestPart2(t *testing.T) {
	input := `svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`

	graph := parseString(input)
	expected := 2
	if result := part2(graph); result != expected {
		t.Errorf("Part 2 expected %d, got %d", expected, result)
	}
}

func parseString(input string) Graph {
	graph := make(Graph)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		node := parts[0]
		neighbors := strings.Fields(parts[1])
		graph[node] = neighbors
	}
	return graph
}
