package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	input := `0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2`

	shapes, regions := parseInput(input)
	if len(shapes) != 6 {
		t.Fatalf("Expected 6 shapes, got %d", len(shapes))
	}
	if len(regions) != 3 {
		t.Fatalf("Expected 3 regions, got %d", len(regions))
	}

	maxID := 0
	for _, s := range shapes {
		if s.ID > maxID {
			maxID = s.ID
		}
	}
	variations := make([][]Shape, maxID+1)
	for _, s := range shapes {
		variations[s.ID] = generateVariations(s)
	}

	results := []bool{true, true, false}

	for i, region := range regions {
		fits := canFit(region, variations)
		if fits != results[i] {
			t.Errorf("Region %d: expected %v, got %v", i, results[i], fits)
		}
	}
}
