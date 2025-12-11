package main

import (
	"fmt"
	"image/color"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Logic ---

type Machine struct {
	target        []bool
	buttons       [][]bool
	buttonIndices [][]int
	targets       []int64
}

func parseMachine(line string) Machine {
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == '[' || r == ']' || r == '{' || r == '}'
	})

	var validParts []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			validParts = append(validParts, p)
		}
	}

	targetStr := strings.TrimSpace(validParts[0])
	target := make([]bool, len(targetStr))
	for i, c := range targetStr {
		target[i] = c == '#'
	}

	buttonsStr := strings.TrimSpace(validParts[1])
	buttonSplits := strings.Split(buttonsStr, ")")
	var buttons [][]bool
	var buttonIndices [][]int

	for _, s := range buttonSplits {
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "(")
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		buttonMask := make([]bool, len(target))
		var indices []int

		nums := strings.Split(s, ",")
		for _, numStr := range nums {
			idx, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err == nil {
				if idx < len(target) {
					buttonMask[idx] = true
				}
				indices = append(indices, idx)
			}
		}
		buttons = append(buttons, buttonMask)
		buttonIndices = append(buttonIndices, indices)
	}

	var targets []int64
	if len(validParts) > 2 {
		targetNums := strings.Split(validParts[2], ",")
		for _, s := range targetNums {
			val, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err == nil {
				targets = append(targets, val)
			}
		}
	}

	return Machine{target, buttons, buttonIndices, targets}
}

func solveMachine(machine Machine) ([]int64, int64) {
	nLights := len(machine.target)
	nButtons := len(machine.buttons)

	if nButtons == 0 {
		allOff := true
		for _, b := range machine.target {
			if b {
				allOff = false
				break
			}
		}
		if allOff {
			return make([]int64, 0), 0
		}
		return nil, -1
	}

	matrix := make([][]bool, nLights)
	for i := range matrix {
		matrix[i] = make([]bool, nButtons+1)
		for j := 0; j < nButtons; j++ {
			matrix[i][j] = machine.buttons[j][i]
		}
		matrix[i][nButtons] = machine.target[i]
	}

	pivotCol := []int{}
	row := 0

	for col := 0; col < nButtons; col++ {
		pivotRow := -1
		for r := row; r < nLights; r++ {
			if matrix[r][col] {
				pivotRow = r
				break
			}
		}

		if pivotRow == -1 {
			continue
		}

		matrix[row], matrix[pivotRow] = matrix[pivotRow], matrix[row]
		pivotCol = append(pivotCol, col)

		for r := 0; r < nLights; r++ {
			if r != row && matrix[r][col] {
				for c := 0; c <= nButtons; c++ {
					if matrix[row][c] {
						matrix[r][c] = !matrix[r][c]
					}
				}
			}
		}
		row++
	}

	for r := row; r < nLights; r++ {
		if matrix[r][nButtons] {
			return nil, -1
		}
	}

	isPivot := make([]bool, nButtons)
	for _, col := range pivotCol {
		isPivot[col] = true
	}
	var freeVars []int
	for i := 0; i < nButtons; i++ {
		if !isPivot[i] {
			freeVars = append(freeVars, i)
		}
	}

	nFree := len(freeVars)
	minPresses := int64(-1)
	var bestSolution []int64

	for mask := 0; mask < (1 << nFree); mask++ {
		solution := make([]bool, nButtons)
		for i, v := range freeVars {
			if (mask & (1 << i)) != 0 {
				solution[v] = true
			}
		}

		for i := len(pivotCol) - 1; i >= 0; i-- {
			col := pivotCol[i]
			r := i
			val := matrix[r][nButtons]
			for b := 0; b < nButtons; b++ {
				if b != col && matrix[r][b] && solution[b] {
					val = !val
				}
			}
			solution[col] = val
		}

		presses := int64(0)
		currentSol := make([]int64, nButtons)
		for i, b := range solution {
			if b {
				presses++
				currentSol[i] = 1
			}
		}

		if minPresses == -1 || presses < minPresses {
			minPresses = presses
			bestSolution = currentSol
		}
	}

	return bestSolution, minPresses
}

func solveJoltage(machine Machine) ([]int64, int64) {
	numCounters := len(machine.targets)
	numButtons := len(machine.buttonIndices)

	if numButtons == 0 {
		allZero := true
		for _, v := range machine.targets {
			if v != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return make([]int64, 0), 0
		}
		return nil, -1
	}

	matrix := make([][]*big.Rat, numCounters)
	for r := 0; r < numCounters; r++ {
		matrix[r] = make([]*big.Rat, numButtons)
		for c := 0; c < numButtons; c++ {
			matrix[r][c] = big.NewRat(0, 1)
		}
	}
	targetVec := make([]*big.Rat, numCounters)

	for r, t := range machine.targets {
		targetVec[r] = new(big.Rat).SetInt64(t)
		for c, indices := range machine.buttonIndices {
			for _, idx := range indices {
				if idx == r {
					matrix[r][c] = big.NewRat(1, 1)
					break
				}
			}
		}
	}

	var pivotCols []int
	var pivotRows []int
	row := 0

	for col := 0; col < numButtons; col++ {
		pivotRow := -1
		for r := row; r < numCounters; r++ {
			if matrix[r][col].Sign() != 0 {
				pivotRow = r
				break
			}
		}

		if pivotRow == -1 {
			continue
		}

		matrix[row], matrix[pivotRow] = matrix[pivotRow], matrix[row]
		targetVec[row], targetVec[pivotRow] = targetVec[pivotRow], targetVec[row]

		pivotCols = append(pivotCols, col)
		pivotRows = append(pivotRows, row)

		pivotVal := new(big.Rat).Set(matrix[row][col])
		invPivot := new(big.Rat).Inv(pivotVal)

		for c := col; c < numButtons; c++ {
			matrix[row][c].Mul(matrix[row][c], invPivot)
		}
		targetVec[row].Mul(targetVec[row], invPivot)

		pivotRowVals := make([]*big.Rat, numButtons)
		for c := 0; c < numButtons; c++ {
			pivotRowVals[c] = new(big.Rat).Set(matrix[row][c])
		}
		pivotTarget := new(big.Rat).Set(targetVec[row])

		for r := 0; r < numCounters; r++ {
			if r != row && matrix[r][col].Sign() != 0 {
				factor := new(big.Rat).Set(matrix[r][col])
				for c := col; c < numButtons; c++ {
					term := new(big.Rat).Mul(factor, pivotRowVals[c])
					matrix[r][c].Sub(matrix[r][c], term)
				}
				term := new(big.Rat).Mul(factor, pivotTarget)
				targetVec[r].Sub(targetVec[r], term)
			}
		}

		row++
	}

	for r := row; r < numCounters; r++ {
		if targetVec[r].Sign() != 0 {
			return nil, -1
		}
	}

	isPivot := make([]bool, numButtons)
	for _, col := range pivotCols {
		isPivot[col] = true
	}
	var freeVars []int
	for i := 0; i < numButtons; i++ {
		if !isPivot[i] {
			freeVars = append(freeVars, i)
		}
	}

	minTotal := int64(-1)
	var bestSolution []int64

	maxTarget := int64(0)
	for _, t := range machine.targets {
		if t > maxTarget {
			maxTarget = t
		}
	}
	limit := maxTarget + 2

	currentFreeVals := make([]int64, len(freeVars))

	var solveRecursive func(idx int)
	solveRecursive = func(idx int) {
		if idx == len(freeVars) {
			currentSolution := make([]int64, numButtons)
			sum := int64(0)

			for i, fv := range freeVars {
				currentSolution[fv] = currentFreeVals[i]
				sum += currentFreeVals[i]
			}

			valid := true
			for i, pc := range pivotCols {
				r := pivotRows[i]
				val := new(big.Rat).Set(targetVec[r])

				for j, fv := range freeVars {
					coeff := matrix[r][fv]
					if coeff.Sign() != 0 {
						term := new(big.Rat).SetInt64(currentFreeVals[j])
						term.Mul(term, coeff)
						val.Sub(val, term)
					}
				}

				if !val.IsInt() || val.Sign() < 0 {
					valid = false
					break
				}
				vInt := val.Num().Int64()
				currentSolution[pc] = vInt
				sum += vInt
			}

			if valid {
				if minTotal == -1 || sum < minTotal {
					minTotal = sum
					bestSolution = make([]int64, numButtons)
					copy(bestSolution, currentSolution)
				}
			}
			return
		}

		for val := int64(0); val <= limit; val++ {
			currentFreeVals[idx] = val
			solveRecursive(idx + 1)
		}
	}

	if len(freeVars) == 0 {
		valid := true
		sum := int64(0)
		solution := make([]int64, numButtons)
		for i, pc := range pivotCols {
			r := pivotRows[i]
			val := targetVec[r]
			if !val.IsInt() || val.Sign() < 0 {
				valid = false
				break
			}
			vInt := val.Num().Int64()
			solution[pc] = vInt
			sum += vInt
		}
		if valid {
			return solution, sum
		}
		return nil, -1
	}

	solveRecursive(0)
	return bestSolution, minTotal
}

func part1(input string) int {
	total := int64(0)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		m := parseMachine(line)
		_, res := solveMachine(m)
		if res != -1 {
			total += res
		}
	}
	return int(total)
}

func part2(input string) int64 {
	total := int64(0)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		m := parseMachine(line)
		_, res := solveJoltage(m)
		if res != -1 {
			total += res
		}
	}
	return total
}

// --- Visualization ---

type Game struct {
	machines   []Machine
	machineIdx int
	mode       int // 0: Part 1, 1: Part 2

	// Pre-calculated solutions
	solutionsP1 [][]int64
	minCostsP1  []int64
	solutionsP2 [][]int64
	minCostsP2  []int64

	// Current State
	currLights   []bool
	currCounters []int64
	userPresses  int64

	// Animation
	autoQueue []int // Queue of button indices to press
	autoTimer int
}

func NewGame(input string) *Game {
	g := &Game{
		mode: 0,
	}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := parseMachine(line)
		g.machines = append(g.machines, m)

		// Pre-calc P1
		sol1, cost1 := solveMachine(m)
		g.solutionsP1 = append(g.solutionsP1, sol1)
		g.minCostsP1 = append(g.minCostsP1, cost1)

		// Pre-calc P2
		sol2, cost2 := solveJoltage(m)
		g.solutionsP2 = append(g.solutionsP2, sol2)
		g.minCostsP2 = append(g.minCostsP2, cost2)
	}
	g.resetMachine()
	return g
}

func (g *Game) resetMachine() {
	m := g.machines[g.machineIdx]
	g.currLights = make([]bool, len(m.target))     // All off initially
	g.currCounters = make([]int64, len(m.targets)) // All zero initially
	g.userPresses = 0
	g.autoQueue = nil
}

func (g *Game) pressButton(idx int) {
	m := g.machines[g.machineIdx]
	g.userPresses++

	if g.mode == 0 {
		// Part 1: Toggle lights
		for i, affected := range m.buttons[idx] {
			if affected {
				g.currLights[i] = !g.currLights[i]
			}
		}
	} else {
		// Part 2: Increment counters
		for _, counterIdx := range m.buttonIndices[idx] {
			if counterIdx < len(g.currCounters) {
				g.currCounters[counterIdx]++
			}
		}
	}
}

func (g *Game) Update() error {
	// Machine Navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.machineIdx = (g.machineIdx + 1) % len(g.machines)
		g.resetMachine()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.machineIdx = (g.machineIdx - 1 + len(g.machines)) % len(g.machines)
		g.resetMachine()
	}

	// Mode Toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.mode = 1 - g.mode
		g.resetMachine()
	}

	// Reset
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetMachine()
	}

	// Auto Solve
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resetMachine()
		if g.mode == 0 {
			if g.solutionsP1[g.machineIdx] != nil {
				for bIdx, count := range g.solutionsP1[g.machineIdx] {
					if count > 0 {
						g.autoQueue = append(g.autoQueue, bIdx)
					}
				}
			}
		} else {
			if g.solutionsP2[g.machineIdx] != nil {
				for bIdx, count := range g.solutionsP2[g.machineIdx] {
					for i := int64(0); i < count; i++ {
						g.autoQueue = append(g.autoQueue, bIdx)
					}
				}
			}
		}
	}

	// Process Animation
	if len(g.autoQueue) > 0 {
		g.autoTimer++
		// Speed up if queue is large
		speed := 10
		if len(g.autoQueue) > 100 {
			speed = 1
		}

		if g.autoTimer >= speed {
			g.autoTimer = 0
			// Process up to N items to keep up
			batch := 1
			if len(g.autoQueue) > 1000 {
				batch = 10
			}

			for i := 0; i < batch && len(g.autoQueue) > 0; i++ {
				idx := g.autoQueue[0]
				g.autoQueue = g.autoQueue[1:]
				g.pressButton(idx)
			}
		}
	}

	// Mouse Interaction
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		m := g.machines[g.machineIdx]
		// Check buttons
		// Buttons are drawn at bottom
		// Let's define layout in Draw and reuse calculations or hardcode here
		startX := 50
		startY := 400
		btnW, btnH := 40, 40
		gap := 10

		for i := 0; i < len(m.buttons); i++ {
			x := startX + (i%10)*(btnW+gap)
			y := startY + (i/10)*(btnH+gap)
			if mx >= x && mx <= x+btnW && my >= y && my <= y+btnH {
				g.pressButton(i)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	m := g.machines[g.machineIdx]

	// Header
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Machine %d/%d  [Left/Right] Select  [M] Mode: %s  [Space] Auto-Solve  [R] Reset",
		g.machineIdx+1, len(g.machines), map[int]string{0: "Part 1 (Lights)", 1: "Part 2 (Joltage)"}[g.mode]), 10, 10)

	// Status
	minCost := g.minCostsP1[g.machineIdx]
	if g.mode == 1 {
		minCost = g.minCostsP2[g.machineIdx]
	}

	statusColor := color.RGBA{200, 200, 200, 255}
	if minCost != -1 && g.userPresses == minCost {
		// Check if actually solved
		solved := true
		if g.mode == 0 {
			for i, t := range m.target {
				if g.currLights[i] != t {
					solved = false
					break
				}
			}
		} else {
			for i, t := range m.targets {
				if g.currCounters[i] != t {
					solved = false
					break
				}
			}
		}
		if solved {
			statusColor = color.RGBA{0, 255, 0, 255}
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Presses: %d / Min: %d", g.userPresses, minCost), 30, 30)
	vector.DrawFilledCircle(screen, 15, 36, 5, statusColor, true)

	// Visualization Area
	// cx, cy := 400, 200 // Unused

	if g.mode == 0 {
		// Draw Lights
		// Target Row
		ebitenutil.DebugPrintAt(screen, "Target:", 50, 100)
		for i, t := range m.target {
			col := color.RGBA{50, 50, 50, 255}
			if t {
				col = color.RGBA{255, 255, 0, 255}
			}
			vector.DrawFilledCircle(screen, float32(150+i*40), 110, 15, col, true)
		}

		// Current Row
		ebitenutil.DebugPrintAt(screen, "Current:", 50, 160)
		for i, l := range g.currLights {
			col := color.RGBA{50, 50, 50, 255}
			if l {
				col = color.RGBA{0, 255, 0, 255}
			}
			vector.DrawFilledCircle(screen, float32(150+i*40), 170, 15, col, true)
		}

	} else {
		// Draw Counters
		ebitenutil.DebugPrintAt(screen, "Joltage Counters (Current / Target):", 50, 100)
		for i, t := range m.targets {
			curr := g.currCounters[i]
			var col color.Color = color.White
			if curr == t {
				col = color.RGBA{0, 255, 0, 255}
			}
			if curr > t {
				col = color.RGBA{255, 0, 0, 255}
			}

			msg := fmt.Sprintf("%d / %d", curr, t)
			x, y := 50+(i%5)*150, 130+(i/5)*30

			// Draw colored indicator
			vector.DrawFilledRect(screen, float32(x-5), float32(y), 3, 12, col, true)
			ebitenutil.DebugPrintAt(screen, msg, x, y)
		}
	}

	// Draw Buttons
	startX := 50
	startY := 400
	btnW, btnH := 40, 40
	gap := 10

	ebitenutil.DebugPrintAt(screen, "Buttons (Click to Toggle/Add):", startX, startY-20)

	mx, my := ebiten.CursorPosition()

	for i := 0; i < len(m.buttons); i++ {
		x := startX + (i%10)*(btnW+gap)
		y := startY + (i/10)*(btnH+gap)

		hover := mx >= x && mx <= x+btnW && my >= y && my <= y+btnH
		col := color.RGBA{100, 100, 200, 255}
		if hover {
			col = color.RGBA{150, 150, 255, 255}
		}

		// Check if this button is part of optimal solution (Cheating hint!)
		/*
			sol := g.solutionsP1[g.machineIdx]
			if g.mode == 1 { sol = g.solutionsP2[g.machineIdx] }
			if sol != nil && sol[i] > 0 {
				vector.StrokeRect(screen, float32(x-2), float32(y-2), float32(btnW+4), float32(btnH+4), 2, color.RGBA{255, 215, 0, 255}, true)
			}
		*/

		vector.DrawFilledRect(screen, float32(x), float32(y), float32(btnW), float32(btnH), col, true)

		// Draw label
		label := fmt.Sprintf("%d", i)
		ebitenutil.DebugPrintAt(screen, label, x+10, y+10)

		// Draw Hover Effects (Lines to affected lights/counters)
		if hover {
			if g.mode == 0 {
				for lIdx, affected := range m.buttons[i] {
					if affected {
						targetX := float32(150 + lIdx*40)
						targetY := float32(170)
						vector.StrokeLine(screen, float32(x+btnW/2), float32(y), targetX, targetY+15, 1, color.RGBA{255, 255, 255, 100}, true)
					}
				}
			} else {
				// Part 2 lines
				for _, cIdx := range m.buttonIndices[i] {
					targetX := float32(50 + (cIdx%5)*150 + 20)
					targetY := float32(130 + (cIdx/5)*30 + 10)
					vector.StrokeLine(screen, float32(x+btnW/2), float32(y), targetX, targetY, 1, color.RGBA{255, 255, 255, 100}, true)
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	data, err := os.ReadFile("../../inputs/day10.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Print Part 1 and Part 2 results to CLI
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("AoC 2025 Day 10 - Factory")
	game := NewGame(input)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
