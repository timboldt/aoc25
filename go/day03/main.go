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
)

// --- Original Logic ---

func parseBanks(input string) []string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var banks []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			banks = append(banks, line)
		}
	}
	return banks
}

// Part 1: Max 2-digit number
func maxJoltage(bank string) uint32 {
	digits := []rune(bank)
	var max uint32 = 0
	for i := 0; i < len(digits); i++ {
		for j := i + 1; j < len(digits); j++ {
			joltageStr := fmt.Sprintf("%c%c", digits[i], digits[j])
			joltage, _ := strconv.ParseUint(joltageStr, 10, 32)
			if uint32(joltage) > max {
				max = uint32(joltage)
			}
		}
	}
	return max
}

func part1(input string) uint32 {
	banks := parseBanks(input)
	var sum uint32 = 0
	for _, bank := range banks {
		sum += maxJoltage(bank)
	}
	return sum
}

// Part 2: Max N-digit number (Greedy approach)
func maxJoltageN(bank string, n int) uint64 {
	digits := []rune(bank)
	toRemove := len(digits) - n
	var stack []rune
	removalsLeft := toRemove
	for _, digit := range digits {
		for len(stack) > 0 && removalsLeft > 0 && stack[len(stack)-1] < digit {
			stack = stack[:len(stack)-1]
			removalsLeft--
		}
		stack = append(stack, digit)
	}
	for removalsLeft > 0 {
		stack = stack[:len(stack)-1]
		removalsLeft--
	}
	result := string(stack[:n])
	num, _ := strconv.ParseUint(result, 10, 64)
	return num
}

func part2(input string) uint64 {
	banks := parseBanks(input)
	var sum uint64 = 0
	for _, bank := range banks {
		sum += maxJoltageN(bank, 12)
	}
	return sum
}

// --- Ebitengine Visualization ---

type VisualizationState struct {
	bankIdx int
	banks   []string

	// Part 2 Animation State (Stack based)
	currentBank  string
	stack        []rune
	digitIdx     int // Index of digit being considered from input
	removalsLeft int

	finishedBank bool

	p1Total uint32
	p2Total uint64

	p2BankResult uint64 // Result for the current bank in Part 2

	lastUpdate time.Time
	stepDelay  time.Duration
}

type Game struct {
	state        *VisualizationState
	windowWidth  int
	windowHeight int
}

func NewGame(input string) *Game {
	banks := parseBanks(input)

	// Initialize with first bank ready to process
	initialState := &VisualizationState{
		banks:        banks,
		bankIdx:      0,
		currentBank:  banks[0],
		stack:        make([]rune, 0),
		digitIdx:     0,
		removalsLeft: len(banks[0]) - 12,    // Target N=12
		stepDelay:    50 * time.Millisecond, // Speed of animation
	}

	return &Game{
		state:        initialState,
		windowWidth:  800,
		windowHeight: 600,
	}
}

func (g *Game) Update() error {
	s := g.state

	// Check if all finished
	if s.bankIdx >= len(s.banks) {
		return nil
	}

	// Simple timer for animation steps
	if time.Since(s.lastUpdate) < s.stepDelay {
		return nil
	}
	s.lastUpdate = time.Now()

	// Logic for Part 2 visualization (Greedy Stack)
	// We animate the stack operations for the current bank

	digits := []rune(s.currentBank)

	if s.digitIdx < len(digits) {
		digit := digits[s.digitIdx]

		// While loop condition from logic:
		// for len(stack) > 0 && removalsLeft > 0 && stack[len(stack)-1] < digit
		if len(s.stack) > 0 && s.removalsLeft > 0 && s.stack[len(s.stack)-1] < digit {
			// Pop from stack
			s.stack = s.stack[:len(s.stack)-1]
			s.removalsLeft--
			// Don't advance digitIdx yet, we re-evaluate with new stack top
		} else {
			// Push to stack
			s.stack = append(s.stack, digit)
			s.digitIdx++
		}
	} else if s.removalsLeft > 0 {
		// Post-processing: remove from end if still need to remove
		s.stack = s.stack[:len(s.stack)-1]
		s.removalsLeft--
	} else {
		// Finished this bank
		// Calculate Part 1 for this bank (instantly) and add to total
		s.p1Total += maxJoltage(s.currentBank)

		// Calculate Part 2 result for this bank from stack
		resStr := string(s.stack[:12]) // Assuming 12 is always valid
		val, _ := strconv.ParseUint(resStr, 10, 64)
		s.p2Total += val
		s.p2BankResult = val

		// Move to next bank
		s.bankIdx++
		if s.bankIdx < len(s.banks) {
			s.currentBank = s.banks[s.bankIdx]
			s.stack = make([]rune, 0)
			s.digitIdx = 0
			s.removalsLeft = len(s.currentBank) - 12
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 40, 255})
	s := g.state

	// Header
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bank %d / %d", s.bankIdx+1, len(s.banks)), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 1 Total: %d", s.p1Total), 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 2 Total: %d", s.p2Total), 10, 50)

	if s.bankIdx >= len(s.banks) {
		ebitenutil.DebugPrintAt(screen, "FINISHED", 350, 250)
		return
	}

	// Visualization of Stack for Part 2
	y := 100
	ebitenutil.DebugPrintAt(screen, "Current Bank (Part 2 Optimization):", 10, y)
	y += 30

	// Show input string with current cursor
	inputStr := s.currentBank
	// Draw full input string
	ebitenutil.DebugPrintAt(screen, inputStr, 10, y)

	// Draw cursor under current digit
	// cursorX calculation removed as unused
	// Actually DebugPrint is small. Let's assume standard debug font width.
	// We'll just draw a caret
	vectorStr := strings.Repeat(" ", s.digitIdx) + "^"
	ebitenutil.DebugPrintAt(screen, vectorStr, 10, y+15)

	y += 50
	ebitenutil.DebugPrintAt(screen, "Stack (Building Max Number):", 10, y)
	y += 30

	stackStr := string(s.stack)
	ebitenutil.DebugPrintAt(screen, stackStr, 10, y)

	y += 50
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Removals Left: %d", s.removalsLeft), 10, y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day03.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Pre-calculate to show expected
	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 03 - Battery Banks")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
