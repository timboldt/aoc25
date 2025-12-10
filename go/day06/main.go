package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Helper Wrapper for Tests ---
func parseAndSolve(input string) int64 {
	cols, _ := parseColumns(input)
	probs := identifyProblems(cols)
	total := int64(0)
	for _, p := range probs {
		total += p.result
	}
	return total
}

// --- Original Logic Refactored for State ---

type ProblemInfo struct {
	colIndices []int
	numbers    []int64
	operator   string
	result     int64
}

// Global parsed data
var (
	gColumns [][]rune
	gRows    int
)

func parseColumns(input string) ([][]rune, int) {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil, 0
	}
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	columns := make([][]rune, maxLen)
	for i := range columns {
		columns[i] = make([]rune, 0)
	}
	for _, line := range lines {
		runes := []rune(line)
		for colIdx, ch := range runes {
			columns[colIdx] = append(columns[colIdx], ch)
		}
		for colIdx := len(runes); colIdx < maxLen; colIdx++ {
			columns[colIdx] = append(columns[colIdx], ' ')
		}
	}
	return columns, len(lines)
}

func identifyProblems(columns [][]rune) []ProblemInfo {
	var problems []ProblemInfo
	var currentCols []int

	for colIdx, column := range columns {
		allSpaces := true
		for _, ch := range column {
			if ch != ' ' {
				allSpaces = false
				break
			}
		}

		if allSpaces {
			if len(currentCols) > 0 {
				if p := processProblemCols(columns, currentCols); p != nil {
					problems = append(problems, *p)
				}
				currentCols = currentCols[:0]
			}
		} else {
			currentCols = append(currentCols, colIdx)
		}
	}
	if len(currentCols) > 0 {
		if p := processProblemCols(columns, currentCols); p != nil {
			problems = append(problems, *p)
		}
	}
	return problems
}

func processProblemCols(columns [][]rune, cols []int) *ProblemInfo {
	if len(cols) == 0 {
		return nil
	}
	numRows := len(columns[cols[0]])

	// Op
	opStr := ""
	for _, c := range cols {
		if c < len(columns) && numRows > 0 {
			opStr += string(columns[c][numRows-1])
		}
	}
	op := strings.TrimSpace(opStr)
	if op != "*" && op != "+" {
		return nil
	}

	// Numbers
	var nums []int64
	for r := 0; r < numRows-1; r++ {
		rowStr := ""
		for _, c := range cols {
			if c < len(columns) && r < len(columns[c]) {
				rowStr += string(columns[c][r])
			}
		}
		trimmed := strings.TrimSpace(rowStr)
		if trimmed != "" {
			if n, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
				nums = append(nums, n)
			}
		}
	}
	if len(nums) == 0 {
		return nil
	}

	// Calc
	var res int64
	if op == "*" {
		res = 1
		for _, n := range nums {
			res *= n
		}
	} else {
		res = 0
		for _, n := range nums {
			res += n
		}
	}

	return &ProblemInfo{
		colIndices: append([]int(nil), cols...), // copy
		numbers:    nums,
		operator:   op,
		result:     res,
	}
}

// --- Ebitengine Visualization ---

type Game struct {
	columns  [][]rune
	problems []ProblemInfo

	// State
	currentProbIdx int
	grandTotal     int64

	// Viewport (Grid is huge, so we only show the current problem area)
	// We'll center the view on the current problem's columns

	// Animation
	state      string // "scrolling", "solving"
	lastUpdate time.Time
	stepDelay  time.Duration

	// Visuals
	flashTimer int
}

func NewGame(input string) *Game {
	cols, rows := parseColumns(input)
	probs := identifyProblems(cols)

	// Global access needed? no, passing around.
	gColumns = cols
	gRows = rows

	return &Game{
		columns:        cols,
		problems:       probs,
		state:          "scrolling",
		stepDelay:      100 * time.Millisecond,
		currentProbIdx: 0,
	}
}

func (g *Game) Update() error {
	if g.currentProbIdx >= len(g.problems) {
		g.state = "finished"
		return nil
	}

	if time.Since(g.lastUpdate) < g.stepDelay {
		return nil
	}
	g.lastUpdate = time.Now()

	switch g.state {
	case "scrolling":
		// Fast forward through "scrolling" by just waiting one frame to center
		// In a real smooth scroll we'd interpolate coordinates, but jumping is fine for now
		g.state = "solving"
		g.flashTimer = 10 // Frames to show highlight

	case "solving":
		if g.flashTimer > 0 {
			g.flashTimer--
		} else {
			// Add result to total
			g.grandTotal += g.problems[g.currentProbIdx].result

			// Move next
			g.currentProbIdx++
			g.state = "scrolling"
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Info
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Problem %d / %d", g.currentProbIdx+1, len(g.problems)), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Grand Total: %d", g.grandTotal), 10, 30)

	if g.state == "finished" {
		ebitenutil.DebugPrintAt(screen, "FINISHED!", 10, 60)
		return
	}

	// Draw relevant columns centered
	prob := g.problems[g.currentProbIdx]

	// Calculate center column of the problem
	probStartCol := prob.colIndices[0]
	probEndCol := prob.colIndices[len(prob.colIndices)-1]
	probCenter := (probStartCol + probEndCol) / 2

	// Screen center
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	cx, cy := sw/2, sh/2

	// Char size
	cw, ch := 12, 20

	// Range of columns to draw to fill screen
	colsVisible := sw / cw
	startCol := probCenter - colsVisible/2
	if startCol < 0 {
		startCol = 0
	}

	for c := startCol; c < startCol+colsVisible && c < len(g.columns); c++ {
		colData := g.columns[c]
		x := cx + (c-probCenter)*cw

		// Highlight if part of current problem
		isCurrent := false
		for _, pc := range prob.colIndices {
			if c == pc {
				isCurrent = true
				break
			}
		}

		if isCurrent && g.state == "solving" {
			// Background highlight
			vector.DrawFilledRect(screen, float32(x), float32(cy-ch*len(colData)/2), float32(cw), float32(ch*len(colData)), color.RGBA{40, 40, 60, 255}, false)
		}

				for r, char := range colData {
					y := cy - (len(colData)/2 * ch) + r*ch
					
					// Highlight operator specially with a small box
					if isCurrent && r == len(colData)-1 && char != ' ' {
						vector.DrawFilledRect(screen, float32(x), float32(y), float32(cw), float32(ch), color.RGBA{255, 100, 100, 100}, false)
					}
					
					ebitenutil.DebugPrintAt(screen, string(char), x, y)
				}	}

	// Draw solved result overlay
	if g.state == "solving" {
		resStr := fmt.Sprintf("= %d", prob.result)
		ebitenutil.DebugPrintAt(screen, resStr, cx, cy+(len(g.columns[0])/2*ch)+20)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	data, err := os.ReadFile("../../inputs/day06.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// CLI-like output for consistency
	// We can't easily call parseAndSolve because it's replaced, but we can rely on viz
	// or re-implement simplistic Part 1 print if needed.
	// For now, let's just run the game.

	game := NewGame(input)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("AoC 2025 Day 06 - Column Math")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
