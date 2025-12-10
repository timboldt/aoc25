package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Original Logic ---

func parseInput(input string) ([]string, int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	startCol := -1
	for col, char := range lines[0] {
		if char == 'S' {
			startCol = col
			break
		}
	}
	return lines, startCol
}

func part1(input string) int {
	grid, startCol := parseInput(input)
	if startCol == -1 {
		return 0
	}
	activeCols := make(map[int]bool)
	activeCols[startCol] = true
	splitCount := 0

	for r := 0; r < len(grid)-1; r++ {
		nextActiveCols := make(map[int]bool)
		for c := range activeCols {
			if c < 0 || c >= len(grid[r+1]) {
				continue
			}
			cell := grid[r+1][c]
			if cell == '^' {
				splitCount++
				nextActiveCols[c-1] = true
				nextActiveCols[c+1] = true
			} else {
				nextActiveCols[c] = true
			}
		}
		activeCols = nextActiveCols
	}
	return splitCount
}

func part2(input string) int {
	grid, startCol := parseInput(input)
	if startCol == -1 {
		return 0
	}
	activeCounts := make(map[int]int)
	activeCounts[startCol] = 1

	for r := 0; r < len(grid)-1; r++ {
		nextActiveCounts := make(map[int]int)
		for c, count := range activeCounts {
			if c < 0 || c >= len(grid[r+1]) {
				continue
			}
			cell := grid[r+1][c]
			if cell == '^' {
				nextActiveCounts[c-1] += count
				nextActiveCounts[c+1] += count
			} else {
				nextActiveCounts[c] += count
			}
		}
		activeCounts = nextActiveCounts
	}

	total := 0
	for _, count := range activeCounts {
		total += count
	}
	return total
}

// --- Ebitengine Visualization ---

const (
	cellSize   = 8 // Small cells to fit more
	viewWidth  = 1200
	viewHeight = 900
)

type Game struct {
	grid     []string
	startCol int

	// Simulation State
	currentRow  int
	activeBeams map[int]int // Map col -> count (used for both parts logic visually)

	// Stats
	p1SplitCount int
	p2Total      int

	// Visuals
	history [][]int // Store history of counts for drawing the trail [row][col] = count

	// Control
	finished   bool
	lastUpdate time.Time
	stepDelay  time.Duration
}

func NewGame(input string) *Game {
	grid, startCol := parseInput(input)

	g := &Game{
		grid:        grid,
		startCol:    startCol,
		currentRow:  0,
		activeBeams: make(map[int]int),
		stepDelay:   50 * time.Millisecond,
		history:     make([][]int, len(grid)),
	}

	// Initialize history rows
	for i := range g.history {
		g.history[i] = make([]int, len(grid[0]))
	}

	// Start
	g.activeBeams[startCol] = 1
	g.history[0][startCol] = 1

	return g
}

func (g *Game) Update() error {
	if g.finished {
		return nil
	}

	// Speed up if key pressed?
	// For now just auto run
	// Actually, let's run fast because 140 rows at 50ms is 7 seconds. That's fine.

	if time.Since(g.lastUpdate) < g.stepDelay {
		return nil
	}
	g.lastUpdate = time.Now()

	// Advance one row
	if g.currentRow >= len(g.grid)-1 {
		g.finished = true
		// Calculate final sum for P2 display
		sum := 0
		for _, c := range g.activeBeams {
			sum += c
		}
		g.p2Total = sum
		return nil
	}

	nextActiveBeams := make(map[int]int)
	r := g.currentRow

	for c, count := range g.activeBeams {
		// Bounds check
		if c < 0 || c >= len(g.grid[r+1]) {
			continue
		}

		cell := g.grid[r+1][c]
		if cell == '^' {
			// Split
			g.p1SplitCount++ // Count splits (Part 1 logic is simply number of split events encountered)

			nextActiveBeams[c-1] += count
			nextActiveBeams[c+1] += count
		} else {
			// Continue
			nextActiveBeams[c] += count
		}
	}

	g.activeBeams = nextActiveBeams
	g.currentRow++

	// Record history for drawing
	for c, count := range g.activeBeams {
		if c >= 0 && c < len(g.history[g.currentRow]) {
			g.history[g.currentRow][c] = count
		}
	}

	// Update running total for P2
	sum := 0
	for _, c := range g.activeBeams {
		sum += c
	}
	g.p2Total = sum

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 20, 255})

	// Calculate Viewport offset to keep current row in middle/bottom
	// Grid is usually ~140 wide, ~140 high.
	// 140 * 8 = 1120. Fits in 1200 width.
	// 140 * 8 = 1120. Fits in 900 height with scrolling.

	offsetY := 0
	if g.currentRow*cellSize > viewHeight/2 {
		offsetY = g.currentRow*cellSize - viewHeight/2
	}
	// Clamp offset
	maxOffset := len(g.grid)*cellSize - viewHeight + 50
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offsetY > maxOffset {
		offsetY = maxOffset
	}

	// Draw Grid (only visible rows)
	startRow := offsetY / cellSize
	endRow := (offsetY+viewHeight)/cellSize + 1
	if startRow < 0 {
		startRow = 0
	}
	if endRow > len(g.grid) {
		endRow = len(g.grid)
	}

	for r := startRow; r < endRow; r++ {
		rowStr := g.grid[r]
		for c, char := range rowStr {
			x := float32(c * cellSize)
			y := float32(r*cellSize - offsetY)

			// Draw static grid elements
			if char == '^' {
				vector.DrawFilledRect(screen, x+1, y+1, cellSize-2, cellSize-2, color.RGBA{100, 100, 100, 255}, false)
			} else if char == 'S' {
				vector.DrawFilledRect(screen, x+1, y+1, cellSize-2, cellSize-2, color.RGBA{255, 255, 0, 255}, false)
			}

			// Draw Beams from history
			count := g.history[r][c]
			if count > 0 {
				// Intensity based on count (logarithmic or capped for visibility?)
				// Counts get HUGE in Part 2.
				// Just use a gradient: Blue -> Cyan -> White

				if count > 1000 { // Very high -> White/Yellow
					vector.DrawFilledRect(screen, x+2, y+2, cellSize-4, cellSize-4, color.RGBA{255, 255, 200, 255}, false)
				} else if count > 100 {
					vector.DrawFilledRect(screen, x+2, y+2, cellSize-4, cellSize-4, color.RGBA{100, 255, 255, 255}, false)
				} else if count > 1 {
					vector.DrawFilledRect(screen, x+2, y+2, cellSize-4, cellSize-4, color.RGBA{0, 200, 255, 255}, false)
				} else {
					vector.DrawFilledRect(screen, x+2, y+2, cellSize-4, cellSize-4, color.RGBA{0, 100, 200, 255}, false)
				}
			}
		}
	}

	// Overlay Info
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Row: %d / %d", g.currentRow, len(g.grid)), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 1 (Splits): %d", g.p1SplitCount), 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 2 (Timelines): %d", g.p2Total), 10, 50)
	if g.finished {
		ebitenutil.DebugPrintAt(screen, "FINISHED", 10, 70)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return viewWidth, viewHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day07.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Pre-calc
	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(viewWidth, viewHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 07 - Tachyon Beams")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
