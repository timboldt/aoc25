package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	x, y int
}

func parsePositions(input string) []Point {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	positions := make([]Point, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		positions = append(positions, Point{x, y})
	}
	return positions
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func part1(positions []Point) int64 {
	n := len(positions)
	maxArea := int64(0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := positions[i]
			p2 := positions[j]
			width := int64(abs(p2.x-p1.x) + 1)
			height := int64(abs(p2.y-p1.y) + 1)
			area := width * height
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

// Coordinate Compression logic
func compressCoordinates(points []Point) ([]int, []int) {
	mux := make(map[int]bool)
	mux := make(map[int]bool)
	for _, p := range points {
		mux[p.x] = true
		mux[p.y] = true
	}
	mux := make([]int, 0, len(tmux))
	mux := make([]int, 0, len(tmux))
	for x := range tmux {
		mux = append(tmux, x)
	}
	for y := range tmux {
		mux = append(tmux, y)
	}
	sort.Ints(tmux)
	sort.Ints(tmux)
	return tmux, tmux
}

func getIndex(val int, arr []int) int {
	// Binary search
	idx := sort.SearchInts(arr, val)
	if idx < len(arr) && arr[idx] == val {
		return idx
	}
	return -1
}

// Build 2D grid where grid[i][j] is 1 if the rectangle [xs[i], xs[i+1]] x [ys[j], ys[j+1]] is inside
func buildGrid(points []Point, xs, ys []int) [][]int {
	w := len(xs) - 1
	h := len(ys) - 1
	if w <= 0 || h <= 0 {
		return make([][]int, 0)
	}
	grid := make([][]int, w)
	for i := range grid {
		grid[i] = make([]int, h)
	}

	// Scanline on compressed Y ranges
	for j := 0; j < h; j++ {
		// Sample Y for this strip
		midY := (float64(ys[j]) + float64(ys[j+1])) / 2.0

		// Find intersections with polygon edges
		var nodes []float64
		n := len(points)
		for k := 0; k < n; k++ {
			p1 := points[k]
			p2 := points[(k+1)%n]

			// Check vertical edge crossing midY
			// p1.y <= midY < p2.y or p2.y <= midY < p1.y
			if (float64(p1.y) < midY && float64(p2.y) >= midY) ||
				(float64(p2.y) < midY && float64(p1.y) >= midY) {
				// Vertical line x = p1.x
				nodes = append(nodes, float64(p1.x))
			}
		}
		sort.Float64s(nodes)

		// Fill intervals between pairs
		for k := 0; k < len(nodes); k += 2 {
			if k+1 >= len(nodes) {
				break
			}
			xStart := nodes[k]
			xEnd := nodes[k+1]

			// Find corresponding xs indices
			// xs[i] >= xStart, xs[i+1] <= xEnd
			for i := 0; i < w; i++ {
				midX := (float64(xs[i]) + float64(xs[i+1])) / 2.0
				if midX > xStart && midX < xEnd {
					grid[i][j] = 1
				}
			}
		}
	}
	return grid
}

func buildPrefixSum(grid [][]int) [][]int {
	w := len(grid)
	if w == 0 {
		return nil
	}
	h := len(grid[0])
	ps := make([][]int, w+1)
	for i := range ps {
		ps[i] = make([]int, h+1)
	}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			ps[i+1][j+1] = ps[i][j+1] + ps[i+1][j] - ps[i][j] + grid[i][j]
		}
	}
	return ps
}

func querySum(ps [][]int, x1, y1, x2, y2 int) int {
	if ps == nil || x1 >= x2 || y1 >= y2 {
		return 0
	}
	return ps[x2][y2] - ps[x1][y2] - ps[x2][y1] + ps[x1][y1]
}

// Checks if a point is on the boundary or inside (raycasting)
func isPointValid(p Point, polygon []Point) bool {
	// 1. Check boundary
	n := len(polygon)
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]
		
		if (p.y == p1.y && p.y == p2.y) && ((p.x >= min(p1.x, p2.x)) && (p.x <= max(p1.x, p2.x))) {
			return true // Point on horizontal segment
		}
		if (p.x == p1.x && p.x == p2.x) && ((p.y >= min(p1.y, p2.y)) && (p.y <= max(p1.y, p2.y))) {
			return true // Point on vertical segment
		}
	}

	// 2. Raycast for interior
	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		// Standard raycasting
		if ((polygon[i].y > p.y) != (polygon[j].y > p.y)) &&
			(p.x < (polygon[j].x-polygon[i].x)*(p.y-polygon[i].y)/(polygon[j].y-polygon[i].y)+polygon[i].x) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func part2(positions []Point) int64 {
	mux, ys := compressCoordinates(positions)
	grid := buildGrid(positions, xs, ys)
	ps := buildPrefixSum(grid)

	n := len(positions)
	maxArea := int64(0)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := positions[i]
			p2 := positions[j]

			xMin := min(p1.x, p2.x)
			xMax := max(p1.x, p2.x)
			yMin := min(p1.y, p2.y)
			yMax := max(p1.y, p2.y)

			// Map to grid indices
			ix1 := getIndex(xMin, xs)
			ix2 := getIndex(xMax, xs)
			iy1 := getIndex(yMin, ys)
			iy2 := getIndex(yMax, ys)

			isValid := false

			if ix1 == ix2 || iy1 == iy2 {
				// Degenerate rectangle (line)
				midX := (xMin + xMax) / 2
				midY := (yMin + yMax) / 2
				
				if isPointValid(Point{xMin, yMin}, positions) && isPointValid(Point{xMax, yMax}, positions) &&
					isPointValid(Point{midX, midY}, positions) {
					isValid = true
				}
			} else {
				// Check full area using prefix sum
				expected := (ix2 - ix1) * (iy2 - iy1)
				actual := querySum(ps, ix1, iy1, ix2, iy2)
				if actual == expected {
					isValid = true
				}
			}

			if isValid {
				area := int64(xMax-xMin+1) * int64(yMax-yMin+1)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}
	return maxArea
}

// --- Ebitengine Visualization ---

const (
	screenWidth  = 1000
	screenHeight = 800
)

type Game struct {
	points []Point
	// Part 2 pre-calc data
	mux, ys []int
	grid   [][]int
	ps     [][]int // Prefix Sum grid

	// State
	bestRectP1 struct {
		x, y, w, h int
		area       int64
	}
	bestRectP2 struct {
		x, y, w, h int
		area       int64
	}

	// Animation
	currentIndexI int
	currentIndexJ int
	mode          int // 0: Part 1, 1: Part 2

	scale   float64
	offsetX float64
	offsetY float64

	// Results
	p1Result int64
	p2Result int64
}

func NewGame(input string) *Game {
	points := parsePositions(input)

	// Calc bounds for scaling
	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y
	for _, p := range points {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	// Add padding
	rangeX := maxX - minX
	rangeY := maxY - minY

	scaleX := float64(screenWidth-100) / float64(rangeX)
	scaleY := float64(screenHeight-100) / float64(rangeY)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	offsetX := 50.0 - float64(minX)*scale
	offsetY := 50.0 - float64(minY)*scale

	// Pre-build Grid for Part 2 visualization context
	mux, ys := compressCoordinates(points)
	grid := buildGrid(points, xs, ys)
	ps := buildPrefixSum(grid)

	return &Game{
		points:        points,
		mux:           xs,
		ys:            ys,
		grid:          grid,
		ps:            ps,
		scale:         scale,
		offsetX:       offsetX,
		offsetY:       offsetY,
		mode:          0,
		currentIndexI: 0,
		currentIndexJ: 1,
	}
}

func (g *Game) Update() error {
	iterations := 1000
	if g.mode == 1 {
		iterations = 500
	}

	for k := 0; k < iterations; k++ {
		if g.currentIndexI >= len(g.points) {
			if g.mode == 0 {
				g.mode = 1
				g.currentIndexI = 0
				g.currentIndexJ = 1
				g.bestRectP2.area = 0 // Reset for Part 2 search
			} else {
				return nil
			}
		}

		p1 := g.points[g.currentIndexI]
		p2 := g.points[g.currentIndexJ]

		xMin := min(p1.x, p2.x)
		xMax := max(p1.x, p2.x)
		yMin := min(p1.y, p2.y)
		yMax := max(p1.y, p2.y)

		width := int64(xMax - xMin + 1)
		height := int64(yMax - yMin + 1)
		area := width * height

		if g.mode == 0 {
			if area > g.bestRectP1.area {
				g.bestRectP1.area = area
				g.bestRectP1.x = xMin
				g.bestRectP1.y = yMin
				g.bestRectP1.w = int(width)
				g.bestRectP1.h = int(height)
			}
		} else {
			// Part 2 Check inside visual loop
			isValid := false
			
			ix1 := getIndex(xMin, g.xs)
			ix2 := getIndex(xMax, g.xs)
			iy1 := getIndex(yMin, g.ys)
			iy2 := getIndex(yMax, g.ys)

			if ix1 != -1 && ix2 != -1 && iy1 != -1 && iy2 != -1 {
				if ix1 == ix2 || iy1 == iy2 {
					midX := (xMin + xMax) / 2
					midY := (yMin + yMax) / 2
					if isPointValid(Point{xMin, yMin}, g.points) && isPointValid(Point{xMax, yMax}, g.points) &&
						isPointValid(Point{midX, midY}, g.points) {
						isValid = true
					}
				} else {
					expected := (ix2 - ix1) * (iy2 - iy1)
					actual := querySum(g.ps, ix1, iy1, ix2, iy2)
					if actual == expected {
						isValid = true
					}
				}
			}

			if isValid {
				if area > g.bestRectP2.area {
					g.bestRectP2.area = area
					g.bestRectP2.x = xMin
					g.bestRectP2.y = yMin
					g.bestRectP2.w = int(width)
					g.bestRectP2.h = int(height)
				}
			}
		}
		}

		g.currentIndexJ++
		if g.currentIndexJ >= len(g.points) {
			g.currentIndexI++
			g.currentIndexJ = g.currentIndexI + 1
		}
	}
	
	// Sync final results
	g.p1Result = g.bestRectP1.area
	g.p2Result = g.bestRectP2.area

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 15, 255})

	// Draw Green area
	for i := 0; i < len(g.xs)-1; i++ {
		for j := 0; j < len(g.ys)-1; j++ {
			if i < len(g.grid) && j < len(g.grid[0]) && g.grid[i][j] == 1 {
				x := float64(g.xs[i])*g.scale + g.offsetX
				y := float64(g.ys[j])*g.scale + g.offsetY
				w := float64(g.xs[i+1]-g.xs[i]) * g.scale
				h := float64(g.ys[j+1]-g.ys[j]) * g.scale
				if w < 1 { w = 1 }
				if h < 1 { h = 1 }
				vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{0, 50, 0, 255}, false)
			}
		}
	}

	// Draw Edges
	for k := 0; k < len(g.points); k++ {
		p1 := g.points[k]
		p2 := g.points[(k+1)%len(g.points)]
		x1 := float32(float64(p1.x)*g.scale + g.offsetX)
		y1 := float32(float64(p1.y)*g.scale + g.offsetY)
		x2 := float32(float64(p2.x)*g.scale + g.offsetX)
		y2 := float32(float64(p2.y)*g.scale + g.offsetY)
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, color.RGBA{0, 255, 0, 255}, false)
	}

	// Draw Red tiles
	for _, p := range g.points {
		x := float32(float64(p.x)*g.scale + g.offsetX)
		y := float32(float64(p.y)*g.scale + g.offsetY)
		vector.DrawFilledCircle(screen, x, y, 3, color.RGBA{255, 0, 0, 255}, false)
	}

	// Draw Best Rect P1
	if g.bestRectP1.area > 0 {
		rx := float32(float64(g.bestRectP1.x)*g.scale + g.offsetX)
		ry := float32(float64(g.bestRectP1.y)*g.scale + g.offsetY)
		rw := float32(float64(g.bestRectP1.w) * g.scale)
		rh := float32(float64(g.bestRectP1.h) * g.scale)
		vector.StrokeRect(screen, rx, ry, rw, rh, 2, color.RGBA{255, 100, 100, 150}, false)
	}
	
	if g.bestRectP2.area > 0 && g.mode == 1 {
		rx := float32(float64(g.bestRectP2.x)*g.scale + g.offsetX)
		ry := float32(float64(g.bestRectP2.y)*g.scale + g.offsetY)
		rw := float32(float64(g.bestRectP2.w) * g.scale)
		rh := float32(float64(g.bestRectP2.h) * g.scale)
		vector.StrokeRect(screen, rx, ry, rw, rh, 3, color.RGBA{100, 100, 255, 255}, false)
	}

	// Text
	modeStr := "Part 1 (Any Rect)"
	if g.mode == 1 { modeStr = "Part 2 (Green Only)" }
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: %s", modeStr), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("P1 Max Area: %d", g.p1Result), 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("P2 Max Area: %d", g.p2Result), 10, 50)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day09.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Pre-calculate
	p1 := part1(parsePositions(input))
	p2 := part2(parsePositions(input))
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)

	game := NewGame(input)
	game.p1Result = 0
	game.p2Result = p2 // Just for show if needed, but animation updates it
	
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 09 - Largest Rectangle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
