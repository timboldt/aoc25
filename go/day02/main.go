package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// --- Original Logic Types ---

type Range struct {
	start uint64
	end   uint64
}

func parseRanges(input string) []Range {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", "") // Remove newlines first
	parts := strings.Split(input, ",")

	var ranges []Range
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		rangeParts := strings.Split(part, "-")
		if len(rangeParts) != 2 {
			log.Printf("Warning: Malformed range part: %s, skipping", part)
			continue
		}
		start, err := strconv.ParseUint(rangeParts[0], 10, 64)
		if err != nil {
			log.Printf("Warning: Failed to parse start of range %s: %v, skipping", part, err)
			continue
		}
		end, err := strconv.ParseUint(rangeParts[1], 10, 64)
		if err != nil {
			log.Printf("Warning: Failed to parse end of range %s: %v, skipping", part, err)
			continue
		}
		ranges = append(ranges, Range{start: start, end: end})
	}

	return ranges
}

func isRepeatedPattern(n uint64) bool {
	s := strconv.FormatUint(n, 10)
	length := len(s)

	if length%2 != 0 {
		return false
	}

	mid := length / 2
	firstHalf := s[:mid]
	secondHalf := s[mid:]

	if strings.HasPrefix(firstHalf, "0") && len(firstHalf) > 1 { // Only check for leading zeros if more than one digit
		return false
	}

	return firstHalf == secondHalf
}

func part1(input string) uint64 {
	ranges := parseRanges(input)
	var sum uint64

	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isRepeatedPattern(id) {
				sum += id
			}
		}
	}

	return sum
}

func isRepeatedPatternV2(n uint64) bool {
	s := strconv.FormatUint(n, 10)
	length := len(s)

	for patternLen := 1; patternLen <= length/2; patternLen++ {
		if length%patternLen == 0 {
			pattern := s[:patternLen]

			if strings.HasPrefix(pattern, "0") && len(pattern) > 1 { // Only check for leading zeros if more than one digit
				continue
			}

			allMatch := true
			for i := 0; i < length; i += patternLen {
				if s[i:i+patternLen] != pattern {
					allMatch = false
					break
				}
			}

			if allMatch {
				return true
			}
		}
	}

	return false
}

func part2(input string) uint64 {
	ranges := parseRanges(input)
	var sum uint64

	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isRepeatedPatternV2(id) {
				sum += id
			}
		}
	}

	return sum
}

// --- Ebitengine Visualization ---

const (
	maxFoundNumbers = 20
)

type Game struct {
	ranges          []Range
	currentRangeIdx int
	currentID       uint64 // The number currently being processed

	part1Sum uint64
	part2Sum uint64

	foundPart1 []uint64
	foundPart2 []uint64

	state         string // "iterating", "finished"
	p1MatchResult bool   // Result of last check for display
	p2MatchResult bool   // Result of last check for display

	batchSize int // how many IDs to process per update

	windowWidth  int
	windowHeight int
}

func NewGame(input string) *Game {
	parsedRanges := parseRanges(input)
	if len(parsedRanges) == 0 {
		log.Fatal("No valid ranges found in input for Day 2.")
	}

	return &Game{
		ranges:          parsedRanges,
		currentRangeIdx: 0,
		currentID:       parsedRanges[0].start,
		state:           "iterating",
		batchSize:       1000, // Process 1000 IDs per frame
		windowWidth:     1000,
		windowHeight:    800,
	}
}

func (g *Game) Update() error {
	if g.state == "finished" {
		return nil
	}

	// Process multiple IDs per update for speed
	for i := 0; i < g.batchSize; i++ {
		currentRange := g.ranges[g.currentRangeIdx]

		if g.currentID > currentRange.end {
			// Move to next range or finish
			g.currentRangeIdx++
			if g.currentRangeIdx >= len(g.ranges) {
				g.state = "finished"
				return nil
			}
			g.currentID = g.ranges[g.currentRangeIdx].start
			continue // Process new range's start ID immediately
		}

		// Check Part 1
		g.p1MatchResult = isRepeatedPattern(g.currentID)
		if g.p1MatchResult {
			g.part1Sum += g.currentID
			g.foundPart1 = append(g.foundPart1, g.currentID)
			if len(g.foundPart1) > maxFoundNumbers {
				g.foundPart1 = g.foundPart1[1:] // Keep only last N
			}
		}

		// Check Part 2
		g.p2MatchResult = isRepeatedPatternV2(g.currentID)
		if g.p2MatchResult {
			g.part2Sum += g.currentID
			g.foundPart2 = append(g.foundPart2, g.currentID)
			if len(g.foundPart2) > maxFoundNumbers {
				g.foundPart2 = g.foundPart2[1:] // Keep only last N
			}
		}

		g.currentID++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255}) // Dark background

	xOffset := 20
	yOffset := 20
	lineHeight := 20

	// General Info
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Range: %d/%d", g.currentRangeIdx+1, len(g.ranges)), xOffset, yOffset)
	yOffset += lineHeight
	if g.currentRangeIdx < len(g.ranges) {
		currentRange := g.ranges[g.currentRangeIdx]
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Current Range: %d - %d", currentRange.start, currentRange.end), xOffset, yOffset)
	} else {
		ebitenutil.DebugPrintAt(screen, "All Ranges Processed", xOffset, yOffset)
	}
	yOffset += lineHeight
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Current ID: %d", g.currentID-1), xOffset, yOffset) // -1 because currentID was incremented at end of Update
	yOffset += lineHeight * 2

	// Part 1 Section
	ebitenutil.DebugPrintAt(screen, "--- Part 1 (Half-Half Pattern) ---", xOffset, yOffset)
	yOffset += lineHeight
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Sum: %d", g.part1Sum), xOffset, yOffset)
	yOffset += lineHeight
	if g.p1MatchResult {
		ebitenutil.DebugPrintAt(screen, "MATCH!", xOffset+100, yOffset-lineHeight)
	}
	ebitenutil.DebugPrintAt(screen, "Found Numbers:", xOffset, yOffset)
	yOffset += lineHeight
	for i := len(g.foundPart1) - 1; i >= 0; i-- { // Display newest first
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("- %d", g.foundPart1[i]), xOffset+20, yOffset)
		yOffset += lineHeight
		if yOffset > g.windowHeight-50 { // Don't overflow
			break
		}
	}

	// Part 2 Section (right side)
	xOffset2 := g.windowWidth / 2
	yOffset2 := 20

	ebitenutil.DebugPrintAt(screen, "--- Part 2 (Any Repeating Pattern) ---", xOffset2, yOffset2)
	yOffset2 += lineHeight
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Sum: %d", g.part2Sum), xOffset2, yOffset2)
	yOffset2 += lineHeight
	if g.p2MatchResult {
		ebitenutil.DebugPrintAt(screen, "MATCH!", xOffset2+100, yOffset2-lineHeight)
	}
	ebitenutil.DebugPrintAt(screen, "Found Numbers:", xOffset2, yOffset2)
	yOffset2 += lineHeight
	for i := len(g.foundPart2) - 1; i >= 0; i-- { // Display newest first
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("- %d", g.foundPart2[i]), xOffset2+20, yOffset2)
		yOffset2 += lineHeight
		if yOffset2 > g.windowHeight-50 { // Don't overflow
			break
		}
	}

	if g.state == "finished" {
		ebitenutil.DebugPrintAt(screen, "SIMULATION FINISHED!", g.windowWidth/2-100, g.windowHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day02.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Print calculated results first
	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 02 - Repeated Patterns")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
