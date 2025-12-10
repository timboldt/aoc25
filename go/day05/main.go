package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Original Logic Types ---

type Range struct {
	start uint64
	end   uint64
}

func (r Range) contains(id uint64) bool {
	return id >= r.start && id <= r.end
}

func parseInput(input string) ([]Range, []uint64) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")

	ranges := []Range{}
	for _, line := range strings.Split(parts[0], "\n") {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}
		start, _ := strconv.ParseUint(parts[0], 10, 64)
		end, _ := strconv.ParseUint(parts[1], 10, 64)
		ranges = append(ranges, Range{start: start, end: end})
	}

	ids := []uint64{}
	if len(parts) > 1 {
		for _, line := range strings.Split(parts[1], "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			id, _ := strconv.ParseUint(line, 10, 64)
			ids = append(ids, id)
		}
	}

	return ranges, ids
}

func isFresh(id uint64, ranges []Range) bool {
	for _, r := range ranges {
		if r.contains(id) {
			return true
		}
	}
	return false
}

func part1(input string) int {
	ranges, ids := parseInput(input)
	count := 0
	for _, id := range ids {
		if isFresh(id, ranges) {
			count++
		}
	}
	return count
}

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{}
	}

	// Sort by start position
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].start < sorted[j].start
	})

	merged := []Range{sorted[0]}

	for _, r := range sorted[1:] {
		last := &merged[len(merged)-1]

		// Check if ranges overlap or are adjacent
		if r.start <= last.end+1 {
			// Merge by extending the end
			if r.end > last.end {
				last.end = r.end
			}
		} else {
			// No overlap, add as new range
			merged = append(merged, r)
		}
	}

	return merged
}

func countIDsInRanges(ranges []Range) uint64 {
	merged := mergeRanges(ranges)
	total := uint64(0)
	for _, r := range merged {
		total += r.end - r.start + 1
	}
	return total
}

func part2(input string) uint64 {
	ranges, _ := parseInput(input)
	return countIDsInRanges(ranges)
}

// --- Ebitengine Visualization ---

const (
	maxVisibleRanges = 15
	lineHeight       = 30
)

type Game struct {
	// Data
	ids           []uint64
	initialRanges []Range // The raw input ranges
	mergedRanges  []Range // The ranges being built/merged in Part 2

	// Part 1 State
	p1IdIdx      int
	p1MatchCount int
	p1CurrentId  uint64
	p1ResultMsg  string
	p1Matched    bool

	// Part 2 State
	p2SortedInput []Range // Input ranges sorted for merging
	p2MergeIdx    int     // Index in p2SortedInput we are currently considering
	p2TotalCount  uint64  // Running total of merged size
	p2Current     *Range  // The "active" range we are trying to merge into (points to element in mergedRanges)

	// Mode
	mode string // "part1", "part2_sort", "part2_merge", "finished"

	// Animation Control
	lastUpdate time.Time
	stepDelay  time.Duration

	windowWidth  int
	windowHeight int
}

func NewGame(input string) *Game {
	ranges, ids := parseInput(input)

	g := &Game{
		ids:           ids,
		initialRanges: ranges,
		mergedRanges:  make([]Range, 0),
		mode:          "part1",
		stepDelay:     25 * time.Millisecond, // Fast speed for Part 1
		windowWidth:   1000,
		windowHeight:  800,
	}

	return g
}

func (g *Game) Update() error {
	if g.mode == "finished" {
		return nil
	}

	if time.Since(g.lastUpdate) < g.stepDelay {
		return nil
	}
	g.lastUpdate = time.Now()

	switch g.mode {
	case "part1":
		if g.p1IdIdx >= len(g.ids) {
			// Finished Part 1
			g.mode = "part2_sort"
			g.stepDelay = 500 * time.Millisecond // Slower for sort/merge visualization

			// Prepare Part 2 data
			// 1. Sort a copy of ranges
			g.p2SortedInput = make([]Range, len(g.initialRanges))
			copy(g.p2SortedInput, g.initialRanges)
			sort.Slice(g.p2SortedInput, func(i, j int) bool {
				return g.p2SortedInput[i].start < g.p2SortedInput[j].start
			})

			return nil
		}

		g.p1CurrentId = g.ids[g.p1IdIdx]
		g.p1Matched = isFresh(g.p1CurrentId, g.initialRanges)

		if g.p1Matched {
			g.p1MatchCount++
			g.p1ResultMsg = "MATCH"
		} else {
			g.p1ResultMsg = "NO MATCH"
		}

		g.p1IdIdx++

	case "part2_sort":
		// Just a visual pause to show "Sorting..."
		// Then initialize the merge process
		if len(g.p2SortedInput) > 0 {
			g.mergedRanges = append(g.mergedRanges, g.p2SortedInput[0])
			g.p2Current = &g.mergedRanges[0]
			g.p2MergeIdx = 1 // Start checking from second element
		}
		g.mode = "part2_merge"
		g.stepDelay = 200 * time.Millisecond

	case "part2_merge":
		if g.p2MergeIdx >= len(g.p2SortedInput) {
			// Calculate final count
			for _, r := range g.mergedRanges {
				g.p2TotalCount += r.end - r.start + 1
			}
			g.mode = "finished"
			return nil
		}

		nextRange := g.p2SortedInput[g.p2MergeIdx]
		// Get pointer to the last element in mergedRanges (our "current")
		lastIdx := len(g.mergedRanges) - 1
		current := &g.mergedRanges[lastIdx]

		if nextRange.start <= current.end+1 {
			// Overlap/Adjacent -> Merge
			// We just update the end of current if needed
			if nextRange.end > current.end {
				current.end = nextRange.end
			}
			// Visually, we consumed nextRange.
		} else {
			// No overlap -> Start new current
			g.mergedRanges = append(g.mergedRanges, nextRange)
			// g.p2Current will point to new last in next frame implicitly or we can ignore pointer
		}
		g.p2MergeIdx++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Layout
	// Top: Part 1
	// Bottom: Part 2

	// --- Header ---
	ebitenutil.DebugPrintAt(screen, "AoC 2025 Day 05 - Range Merging", 10, 10)

	// --- Part 1 Area ---
	y := 40
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("--- Part 1: ID Checking (%d/%d) ---", g.p1IdIdx, len(g.ids)), 10, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Current ID: %d", g.p1CurrentId), 10, y)
	y += 20

	resultCol := color.RGBA{150, 150, 150, 255}
	if g.p1ResultMsg == "MATCH" {
		resultCol = color.RGBA{50, 255, 50, 255}
	} else if g.p1ResultMsg == "NO MATCH" {
		resultCol = color.RGBA{255, 50, 50, 255}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Result: %s", g.p1ResultMsg), 10, y)
	// Draw a colored box for result
	vector.DrawFilledRect(screen, 200, float32(y), 20, 20, resultCol, false)

	y += 30
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Total Matches: %d", g.p1MatchCount), 10, y)

	// --- Part 2 Area ---
	y = 200
	ebitenutil.DebugPrintAt(screen, "---", 10, y)
	y += 30

	if g.mode == "part1" {
		ebitenutil.DebugPrintAt(screen, "Waiting for Part 1...", 10, y)
		return
	}

	// Draw merged ranges (stack growing downwards)
	// We only show the last few merged ranges to keep it on screen

	startIndex := 0
	if len(g.mergedRanges) > maxVisibleRanges {
		startIndex = len(g.mergedRanges) - maxVisibleRanges
	}

	ebitenutil.DebugPrintAt(screen, "Merged/Active Ranges:", 10, y)
	y += 20

	for i := startIndex; i < len(g.mergedRanges); i++ {
		r := g.mergedRanges[i]
		// Highlight the very last one as "Active"
		col := color.RGBA{100, 100, 255, 255} // Blueish for settled
		prefix := "  "
		if i == len(g.mergedRanges)-1 {
			col = color.RGBA{100, 255, 100, 255} // Green for active/merging
			prefix = "> "
		}

		str := fmt.Sprintf("%s[%d - %d]", prefix, r.start, r.end)
		ebitenutil.DebugPrintAt(screen, str, 10, y)

		// Draw a representative bar
		// Since coords are huge, we just draw a fixed width bar
		barWidth := float32(400)
		vector.DrawFilledRect(screen, 300, float32(y), barWidth, 15, col, false)

		y += 20
	}

	// Show the "Next" range being considered (from sorted input)
	if g.mode == "part2_merge" && g.p2MergeIdx < len(g.p2SortedInput) {
		nextR := g.p2SortedInput[g.p2MergeIdx]
		y += 20
		ebitenutil.DebugPrintAt(screen, "Next Input Range:", 10, y)
		y += 20
		str := fmt.Sprintf("? [%d - %d]", nextR.start, nextR.end)
		ebitenutil.DebugPrintAt(screen, str, 10, y)

		// Draw "Next" bar in yellow
		vector.DrawFilledRect(screen, 300, float32(y), 200, 15, color.RGBA{255, 255, 50, 255}, false)
	}

	// Final Result
	if g.mode == "finished" {
		y += 60
		ebitenutil.DebugPrintAt(screen, "---", 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 2 Total Coverage: %d", g.p2TotalCount), 10, y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day05.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Pre-calc for CLI consistency check
	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 05 - Range Logic")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
