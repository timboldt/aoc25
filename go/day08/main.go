package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Original Logic ---

type Point struct {
	x, y, z int
	id      int
}

type Connection struct {
	p1, p2 int // indices in points slice
	distSq int
}

func parseInput(input string) []Point {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	points := make([]Point, 0, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points = append(points, Point{x, y, z, i})
	}
	return points
}

func distSq(a, b Point) int {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

// DSU implementation
type DSU struct {
	parent []int
	size   []int
	count  int // number of disjoint sets
}

func NewDSU(n int) *DSU {
	dsu := &DSU{
		parent: make([]int, n),
		size:   make([]int, n),
		count:  n,
	}
	for i := 0; i < n; i++ {
		dsu.parent[i] = i
		dsu.size[i] = 1
	}
	return dsu
}

func (d *DSU) Find(i int) int {
	if d.parent[i] != i {
		d.parent[i] = d.Find(d.parent[i])
	}
	return d.parent[i]
}

func (d *DSU) Union(i, j int) bool {
	rootI := d.Find(i)
	rootJ := d.Find(j)
	if rootI != rootJ {
		// Merge smaller into larger
		if d.size[rootI] < d.size[rootJ] {
			rootI, rootJ = rootJ, rootI
		}
		d.parent[rootJ] = rootI
		d.size[rootI] += d.size[rootJ]
		d.count--
		return true
	}
	return false
}

func solve(input string) (int, int) {
	points := parseInput(input)
	n := len(points)

	connections := make([]Connection, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := distSq(points[i], points[j])
			connections = append(connections, Connection{i, j, d})
		}
	}

	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distSq < connections[j].distSq
	})

	// Part 1
	dsu1 := NewDSU(n)
	limit := 1000
	if len(connections) < limit {
		limit = len(connections)
	}

	for i := 0; i < limit; i++ {
		dsu1.Union(connections[i].p1, connections[i].p2)
	}

	sizes := make(map[int]int)
	for i := 0; i < n; i++ {
		root := dsu1.Find(i)
		sizes[root] = dsu1.size[root]
	}

	sizeList := make([]int, 0, len(sizes))
	for _, s := range sizes {
		sizeList = append(sizeList, s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizeList)))

	part1 := 0
	if len(sizeList) >= 3 {
		part1 = sizeList[0] * sizeList[1] * sizeList[2]
	} else {
		part1 = 1
		for _, s := range sizeList {
			part1 *= s
		}
	}

	// Part 2
	dsu2 := NewDSU(n)
	part2 := 0

	for _, conn := range connections {
		if dsu2.Union(conn.p1, conn.p2) {
			if dsu2.count == 1 {
				p1 := points[conn.p1]
				p2 := points[conn.p2]
				part2 = p1.x * p2.x
				break
			}
		}
	}

	return part1, part2
}

// --- Ebitengine Visualization ---

const (
	screenWidth  = 1200
	screenHeight = 900
)

type Game struct {
	points      []Point
	connections []Connection
	dsu         *DSU

	// Visual State
	screenPoints [][2]float32
	colors       []color.RGBA

	// Animation State
	connIdx     int
	activeLines []Connection // Recently added lines for drawing

	// Results
	part1Result int
	part2Result int
	finished    bool
	paused      bool

	// Zoom/Pan (simple fitting for now)
	scale   float32
	offsetX float32
	offsetY float32
}

func NewGame(input string) *Game {
	points := parseInput(input)
	n := len(points)

	// Generate connections
	connections := make([]Connection, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := distSq(points[i], points[j])
			connections = append(connections, Connection{i, j, d})
		}
	}
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distSq < connections[j].distSq
	})

	// Calculate screen coordinates (Isometric projection-ish)
	// Find bounds first
	minX, maxX := 1000000, -1000000
	minY, maxY := 1000000, -1000000

	sp := make([][2]float32, n)
	for i, p := range points {
		// Project: x' = x - z*0.5, y' = y - z*0.5
		// Adjust for better spread?
		// Input range is likely 0-1000.

		// Simple Isometric rotation
		// x_iso = (x - y) * cos(30)
		// y_iso = (x + y) * sin(30) - z

		// Let's try simple cabinet projection
		// x_screen = x + 0.5 * z * cos(angle)
		// y_screen = y + 0.5 * z * sin(angle)

		isoX := float32(p.x) + float32(p.z)*0.5
		isoY := float32(p.y) + float32(p.z)*0.5

		sp[i] = [2]float32{isoX, isoY}

		if int(isoX) < minX {
			minX = int(isoX)
		}
		if int(isoX) > maxX {
			maxX = int(isoX)
		}
		if int(isoY) < minY {
			minY = int(isoY)
		}
		if int(isoY) > maxY {
			maxY = int(isoY)
		}
	}

	// Auto scale
	contentW := float32(maxX - minX)
	contentH := float32(maxY - minY)

	scale := float32(screenWidth-100) / contentW
	scaleY := float32(screenHeight-100) / contentH
	if scaleY < scale {
		scale = scaleY
	}

	offsetX := float32(50) - float32(minX)*scale
	offsetY := float32(50) - float32(minY)*scale

	// Initialize Colors (random per component)
	colors := make([]color.RGBA, n)
	for i := range colors {
		colors[i] = color.RGBA{
			R: uint8(rand.Intn(200) + 55),
			G: uint8(rand.Intn(200) + 55),
			B: uint8(rand.Intn(200) + 55),
			A: 255,
		}
	}

	return &Game{
		points:       points,
		connections:  connections,
		dsu:          NewDSU(n),
		screenPoints: sp,
		colors:       colors,
		scale:        scale,
		offsetX:      offsetX,
		offsetY:      offsetY,
		activeLines:  make([]Connection, 0),
	}
}

func (g *Game) Update() error {
	if g.finished || g.paused {
		return nil
	}

	// Process multiple connections per frame to speed up
	speed := 4
	// if g.connIdx > 1000 {
	// 	speed = 5 // Faster for Part 2
	// }

	for k := 0; k < speed; k++ {
		if g.connIdx >= len(g.connections) {
			g.finished = true
			break
		}

		conn := g.connections[g.connIdx]

		// Part 1 Logic: Capture state at 1000
		if g.connIdx == 1000 {
			// Calculate Part 1 result
			sizes := make(map[int]int)
			for i := 0; i < len(g.points); i++ {
				root := g.dsu.Find(i)
				sizes[root] = g.dsu.size[root]
			}
			sizeList := make([]int, 0, len(sizes))
			for _, s := range sizes {
				sizeList = append(sizeList, s)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(sizeList)))
			if len(sizeList) >= 3 {
				g.part1Result = sizeList[0] * sizeList[1] * sizeList[2]
			}

			// Optional pause
			// g.paused = true
		}

		// Attempt Union
		if g.dsu.Union(conn.p1, conn.p2) {
			// Successful merge
			// Update colors?
			// Propagate color of larger component to smaller one?
			// The DSU Union merges smaller (J) into larger (I).
			// So parent[rootJ] = rootI.
			// We can lazily handle colors in Draw by checking Find().

			// Keep track of active line for drawing
			g.activeLines = append(g.activeLines, conn)
			if len(g.activeLines) > 500 {
				g.activeLines = g.activeLines[1:] // Keep last 500 lines active
			}

			// Check Part 2 condition
			if g.dsu.count == 1 {
				p1 := g.points[conn.p1]
				p2 := g.points[conn.p2]
				g.part2Result = p1.x * p2.x
				g.finished = true
				break
			}
		}

		g.connIdx++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 15, 255})

	// Draw Connections
	// Drawing ALL connections is too heavy and messy.
	// Draw the MST connections (those that caused a merge).
	// We stored them in g.activeLines? No, activeLines are just recent ones.
	// We probably want to draw all "tree" edges found so far.
	// But storing them all might be good.
	// For visualization, drawing the *recent* activity + points is mostly enough.

	// Draw Lines
	for _, conn := range g.activeLines {
		p1 := g.screenPoints[conn.p1]
		p2 := g.screenPoints[conn.p2]

		x1 := p1[0]*g.scale + g.offsetX
		y1 := p1[1]*g.scale + g.offsetY
		x2 := p2[0]*g.scale + g.offsetX
		y2 := p2[1]*g.scale + g.offsetY

		col := color.RGBA{100, 100, 100, 100}
		// Highlight the latest one
		if conn == g.activeLines[len(g.activeLines)-1] {
			col = color.RGBA{255, 255, 255, 255}
		}

		vector.StrokeLine(screen, x1, y1, x2, y2, 1, col, false)
	}

	// Draw Points
	for i, sp := range g.screenPoints {
		x := sp[0]*g.scale + g.offsetX
		y := sp[1]*g.scale + g.offsetY

		// Get color based on component
		root := g.dsu.Find(i)
		col := g.colors[root] // Initial random color of the root

		// Draw
		vector.DrawFilledCircle(screen, x, y, 2, col, false)
	}

	// UI
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Processed Pairs: %d / %d", g.connIdx, len(g.connections)), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Components: %d", g.dsu.count), 10, 30)

	if g.part1Result > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 1 Result: %d", g.part1Result), 10, 50)
	}
	if g.finished {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Part 2 Result: %d", g.part2Result), 10, 70)
		ebitenutil.DebugPrintAt(screen, "FINISHED", 10, 90)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	data, err := os.ReadFile("../../inputs/day08.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)

	// Pre-calc
	p1, p2 := solve(input)
	fmt.Printf("Calculated Part 1: %d\n", p1)
	fmt.Printf("Calculated Part 2: %d\n", p2)

	game := NewGame(input)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AoC 2025 Day 08 - Junction Boxes")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
