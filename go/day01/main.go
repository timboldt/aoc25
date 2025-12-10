package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Original Logic Types ---

type Direction int

const (
	Left Direction = iota
	Right
)

type Rotation struct {
	direction Direction
	distance  int
}

func parseRotations(input string) []Rotation {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rotations := make([]Rotation, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var direction Direction
		if line[0] == 'L' {
			direction = Left
		} else if line[0] == 'R' {
			direction = Right
		} else {
			panic(fmt.Sprintf("Unknown direction: %c", line[0]))
		}

		distance, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(fmt.Sprintf("Failed to parse distance: %v", err))
		}

		rotations = append(rotations, Rotation{
			direction: direction,
			distance:  distance,
		})
	}

	return rotations
}

// mod implements modulo operation that handles negative numbers correctly
func mod(a, b int) int {
	return ((a % b) + b) % b
}

func part1(input string) int64 {
	rotations := parseRotations(input)
	position := 50
	count := int64(0)

	for _, rotation := range rotations {
		switch rotation.direction {
		case Left:
			position = mod(position-rotation.distance, 100)
		case Right:
			position = mod(position+rotation.distance, 100)
		}

		if position == 0 {
			count++
		}
	}

	return count
}

func part2(input string) int64 {
	rotations := parseRotations(input)
	position := 50
	count := int64(0)

	for _, rotation := range rotations {
		var clicksThroughZero int64
		switch rotation.direction {
		case Right:
			clicksThroughZero = int64((position+rotation.distance)/100 - position/100)
		case Left:
			if position == 0 {
				clicksThroughZero = int64(rotation.distance / 100)
			} else if rotation.distance >= position {
				clicksThroughZero = int64((rotation.distance-position)/100 + 1)
			} else {
				clicksThroughZero = 0
			}
		}
		count += clicksThroughZero
		switch rotation.direction {
		case Left:
			position = mod(position-rotation.distance, 100)
		case Right:
			position = mod(position+rotation.distance, 100)
		}
	}
	return count
}

// --- Ebitengine Visualization ---

type Game struct {
	rotations    []Rotation
	idx          int
	currPos      float64
	targetPos    float64
	p1Count      int64
	p2Count      int64
	state        string // "idle", "rotating", "finished"
	waitTimer    int
	windowWidth  int
	windowHeight int
}

func NewGame(input string) *Game {
	return &Game{
		rotations: parseRotations(input),
		currPos:   50.0, // Start at 50 as per problem description (implied by first move relative to it? No, problem usually starts at 0, but the code says 50. Wait, let's check the code provided.)
		// Original code: position := 50. So we start at 50.
		targetPos:    50.0,
		state:        "idle",
		windowWidth:  800,
		windowHeight: 600,
	}
}

func (g *Game) Update() error {
	switch g.state {
	case "idle":
		if g.idx >= len(g.rotations) {
			g.state = "finished"
			return nil
		}

		// Prepare next move
		rot := g.rotations[g.idx]
		if rot.direction == Right {
			g.targetPos = g.currPos + float64(rot.distance)
		} else {
			g.targetPos = g.currPos - float64(rot.distance)
		}
		g.state = "rotating"

	case "rotating":
		speed := 0.5 // Base speed
		diff := g.targetPos - g.currPos

		// Accelerate based on distance remaining
		speed = math.Max(math.Abs(diff)*0.1, 0.5)

		if math.Abs(diff) < speed {
			// Arrived
			oldPosFloor := int64(math.Floor(g.currPos / 100.0))
			g.currPos = g.targetPos
			newPosFloor := int64(math.Floor(g.currPos / 100.0))

			// Check Part 2 (crossing 0)
			// Wait, if we jump the last step, we might miss a crossing if the step is large?
			// But here the step is small (< speed), so we are fine.
			// Actually, we need to accumulate crossings from the jump.
			g.p2Count += int64(math.Abs(float64(newPosFloor - oldPosFloor)))

			// Check Part 1 (land on 0)
			normalizedPos := mod(int(math.Round(g.currPos)), 100)
			if normalizedPos == 0 {
				g.p1Count++
			}

			g.idx++
			g.state = "idle"
			// Add a small pause between moves for visual clarity
			// We can implement a timer state if needed, but instant is fine for now
		} else {
			oldPosFloor := int64(math.Floor(g.currPos / 100.0))
			if diff > 0 {
				g.currPos += speed
			} else {
				g.currPos -= speed
			}
			newPosFloor := int64(math.Floor(g.currPos / 100.0))
			g.p2Count += int64(math.Abs(float64(newPosFloor - oldPosFloor)))
		}

	case "finished":
		// Do nothing
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	cx, cy := float32(g.windowWidth)/2, float32(g.windowHeight)/2
	radius := float32(200.0)

	// Draw Dial
	vector.StrokeCircle(screen, cx, cy, radius, 5, color.White, true)
	vector.StrokeCircle(screen, cx, cy, radius-20, 2, color.Gray{100}, true)

	// Draw Ticks and Numbers
	// Position 0 is at Top (12 o'clock).
	// Values increase Clockwise (Right).
	// So angle for value V: -90 degrees + (V/100 * 360)
	// But the dial rotates.
	// If the Indicator is fixed at Top, and we are at Position P.
	// Then the number P is at Top.
	// So the number 0 is at angle: -90 - (P/100 * 360).
	// Let's verify: If P=0, 0 is at -90 (Top). Correct.
	// If P=25 (Right), 25 is at Top. So 0 should be at Top - 90 = Left.
	// -90 - (25/100*360) = -90 - 90 = -180 (Left). Correct.

	dialRotation := -float64(g.currPos) / 100.0 * 2 * math.Pi

	for i := 0; i < 100; i++ {
		angle := float64(i)/100.0*2*math.Pi + dialRotation - math.Pi/2
		x1 := cx + float32(math.Cos(angle))*radius
		y1 := cy + float32(math.Sin(angle))*radius
		x2 := cx + float32(math.Cos(angle))*(radius-10)
		y2 := cy + float32(math.Sin(angle))*(radius-10)

		col := color.RGBA{100, 100, 100, 255}
		if i%10 == 0 {
			col = color.RGBA{255, 255, 255, 255}
			x2 = cx + float32(math.Cos(angle))*(radius-20)
			y2 = cy + float32(math.Sin(angle))*(radius-20)

			// Draw number
			numStr := fmt.Sprintf("%d", i)
			// Position text slightly inside the tick
			textRadius := float64(radius - 35)
			tx := float64(cx) + math.Cos(angle)*textRadius
			ty := float64(cy) + math.Sin(angle)*textRadius
			// Center text (approximate)
			ebitenutil.DebugPrintAt(screen, numStr, int(tx)-4*len(numStr), int(ty)-8)
		}

		vector.StrokeLine(screen, x1, y1, x2, y2, 2, col, true)
	}

	// Draw Indicator (Fixed at Top)
	// Make it more prominent, pointing down to the current number
	vector.StrokeLine(screen, cx, cy-radius-15, cx, cy-radius+15, 4, color.RGBA{255, 50, 50, 255}, true)
	// Add a small triangle pointing down
	vector.DrawFilledCircle(screen, cx, cy-radius-15, 5, color.RGBA{255, 50, 50, 255}, true)

	// Draw Center Knob
	vector.DrawFilledCircle(screen, cx, cy, 40, color.Gray{50}, true)
	vector.StrokeCircle(screen, cx, cy, 40, 2, color.White, true)

	// Draw Text Info
	normalizedPos := mod(int(math.Round(g.currPos)), 100)
	msg := fmt.Sprintf("Pos: %d", normalizedPos)
	ebitenutil.DebugPrintAt(screen, msg, int(cx)-20, int(cy)-10)

	status := fmt.Sprintf("Part 1 (Land on 0): %d\nPart 2 (Pass 0): %d", g.p1Count, g.p2Count)
	ebitenutil.DebugPrintAt(screen, status, 10, 10)

	if g.state == "rotating" && g.idx < len(g.rotations) {
		rot := g.rotations[g.idx]
		dirStr := "L"
		if rot.direction == Right {
			dirStr = "R"
		}
		cmd := fmt.Sprintf("Executing: %s%d", dirStr, rot.distance)
		ebitenutil.DebugPrintAt(screen, cmd, 10, 50)
	} else if g.state == "finished" {
		ebitenutil.DebugPrintAt(screen, "FINISHED", 10, 50)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.windowWidth, g.windowHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day01.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Print calculated results first
	fmt.Printf("Calculated Part 1: %d\n", part1(input))
	fmt.Printf("Calculated Part 2: %d\n", part2(input))

	game := NewGame(input)
	ebiten.SetWindowSize(game.windowWidth, game.windowHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 01 - Safe Dial")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
