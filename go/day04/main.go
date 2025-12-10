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

func parseInput(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

func countAdjacentPapers(grid [][]rune, row, col int) int {
	rows := len(grid)
	cols := len(grid[0])
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	count := 0
	for _, dir := range directions {
		newRow := row + dir[0]
		newCol := col + dir[1]

		if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols {
			if grid[newRow][newCol] == '@' {
				count++
			}
		}
	}
	return count
}

func part1(input string) int {
	grid := parseInput(input)
	accessible := 0

	for row, line := range grid {
		for col, cell := range line {
			if cell == '@' {
				adjacentCount := countAdjacentPapers(grid, row, col)
				if adjacentCount < 4 {
					accessible++
				}
			}
		}
	}

	return accessible
}

// deepCopyGrid creates a deep copy of a 2D rune slice.
func deepCopyGrid(grid [][]rune) [][]rune {
	duplicate := make([][]rune, len(grid))
	for i := range grid {
		duplicate[i] = make([]rune, len(grid[i]))
		copy(duplicate[i], grid[i])
	}
	return duplicate
}

func part2(input string) int {
	grid := parseInput(input)
	totalRemoved := 0

	for {
		toRemove := [][2]int{}

		for row, line := range grid {
			for col, cell := range line {
				if cell == '@' {
					adjacentCount := countAdjacentPapers(grid, row, col)
					if adjacentCount < 4 {
						toRemove = append(toRemove, [2]int{row, col})
					}
				}
			}
		}

		if len(toRemove) == 0 {
			break
		}

		for _, pos := range toRemove {
			grid[pos[0]][pos[1]] = '.'
		}

		totalRemoved += len(toRemove)
	}

	return totalRemoved
}

// --- Ebitengine Visualization ---

const (
	cellSize    = 5
	gridPadding = 20
)

type Game struct {
	initialGrid [][]rune
	grid        [][]rune // Grid for part 2 simulation
	gridRows    int
	gridCols    int

	// Part 1 state
	p1AccessibleCount int
	p1AccessibleCells [][2]int // Coordinates of accessible cells for part 1

	// Part 2 state
	p2TotalRemoved      int
	p2CellsToRemove     [][2]int // Cells identified for removal in current iteration
	p2CurrentIteration  int
	p2IterationFinished bool // Flag for one pass of finding+removing
	p2Finished          bool // Flag for entire part 2 simulation

	// Mode control
	mode string // "part1_display", "part2_identifying", "part2_removing", "finished"

	lastUpdate time.Time
	stepDelay  time.Duration // Delay between animation steps (identifying -> removing -> next iter)

	windowWidth  int
	windowHeight int
}

func NewGame(input string) *Game {
	initialGrid := parseInput(input)
	if len(initialGrid) == 0 || len(initialGrid[0]) == 0 {
		log.Fatal("Invalid input grid for Day 4.")
	}

	rows := len(initialGrid)
	cols := len(initialGrid[0])

	game := &Game{
		initialGrid: initialGrid,
		grid:        deepCopyGrid(initialGrid), // Start Part 2 with a fresh copy
		gridRows:    rows,
		gridCols:    cols,
		stepDelay:   100 * time.Millisecond,
		mode:        "part1_display",
	}

	// Calculate Part 1 data upfront
	for r, line := range initialGrid {
		for c, cell := range line {
			if cell == '@' {
				adjacentCount := countAdjacentPapers(initialGrid, r, c)
				if adjacentCount < 4 {
					game.p1AccessibleCount++
					game.p1AccessibleCells = append(game.p1AccessibleCells, [2]int{r, c})
				}
			}
		}
	}

	game.windowWidth = cols*cellSize + 2*gridPadding + 200 // Add space for text
	game.windowHeight = rows*cellSize + 2*gridPadding + 100

	return game
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
	case "part1_display":
		// Display for a bit, then switch to part 2
		g.mode = "part2_identifying"
		g.p2CurrentIteration = 1

	case "part2_identifying":
		g.p2CellsToRemove = nil // Clear from previous iteration
		foundAny := false
		for r, line := range g.grid {
			for c, cell := range line {
				if cell == '@' {
					adjacentCount := countAdjacentPapers(g.grid, r, c)
					if adjacentCount < 4 {
						g.p2CellsToRemove = append(g.p2CellsToRemove, [2]int{r, c})
						foundAny = true
					}
				}
			}
		}

		if !foundAny {
			g.p2Finished = true
			g.mode = "finished"
		} else {
			g.mode = "part2_removing"
		}

	case "part2_removing":
		for _, pos := range g.p2CellsToRemove {
			g.grid[pos[0]][pos[1]] = '.'
		}
		g.p2TotalRemoved += len(g.p2CellsToRemove)
		g.p2CurrentIteration++
		g.mode = "part2_identifying" // Go back to identifying for next iteration
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw grid
	for r := 0; r < g.gridRows; r++ {
		for c := 0; c < g.gridCols; c++ {
			x := float32(gridPadding + c*cellSize)
			y := float32(gridPadding + r*cellSize)

			cellColor := color.RGBA{50, 50, 60, 255} // Default empty cell
			if g.grid[r][c] == '@' {
				cellColor = color.RGBA{150, 150, 150, 255} // Paper roll
			}

			// Highlight accessible cells for Part 1
			if g.mode == "part1_display" {
				for _, p := range g.p1AccessibleCells {
					if p[0] == r && p[1] == c {
						cellColor = color.RGBA{50, 200, 50, 255} // Green for accessible
						break
					}
				}
			}

			// Highlight cells to be removed for Part 2
			if g.mode == "part2_removing" || g.mode == "part2_identifying" {
				for _, p := range g.p2CellsToRemove {
					if p[0] == r && p[1] == c {
						cellColor = color.RGBA{200, 200, 50, 255} // Yellow for to-be-removed
						break
					}
				}
			}

			vector.DrawFilledRect(screen, x, y, cellSize, cellSize, cellColor, false)
			vector.StrokeRect(screen, x, y, cellSize, cellSize, 1, color.RGBA{80, 80, 80, 255}, false) // Cell border
		}
	}

	// Display Info
	textX := float32(g.gridCols*cellSize + 2*gridPadding + 10)
	textY := float32(gridPadding)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("AoC Day 04: Paper Rolls"), int(textX), int(textY))
	textY += 20

	if g.mode == "part1_display" {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: Part 1 Analysis"), int(textX), int(textY))
		textY += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Accessible Count: %d", g.p1AccessibleCount), int(textX), int(textY))
	} else if g.mode == "part2_identifying" || g.mode == "part2_removing" {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: Part 2 Simulation"), int(textX), int(textY))
		textY += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Iteration: %d", g.p2CurrentIteration), int(textX), int(textY))
		textY += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Total Removed: %d", g.p2TotalRemoved), int(textX), int(textY))
		textY += 20
		if g.mode == "part2_identifying" {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Identifying rolls to remove..."), int(textX), int(textY))
		} else if g.mode == "part2_removing" {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Removing %d rolls...", len(g.p2CellsToRemove)), int(textX), int(textY))
		}
	} else if g.mode == "finished" {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SIMULATION FINISHED!"), int(textX), int(textY))
		textY += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Final Removed: %d", g.p2TotalRemoved), int(textX), int(textY))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day04.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 04 - Paper Rolls")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
