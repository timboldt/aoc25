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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Domain Structures ---

type Point struct {
	X, Y int
}

type Shape struct {
	ID     int
	Points []Point
	Width  int
	Height int
}

type Region struct {
	Width    int
	Height   int
	Presents []int // List of Shape IDs
}

type PresentItem struct {
	ID   int
	Area int
}

// --- Solver State Structures ---

type SolverFrame struct {
	ItemIndex int
	// Iterator state
	R, C, VarIdx int
	// Backtracking state
	Placed      bool
	PlacedShape Shape
	PlacedX     int
	PlacedY     int
}

type Solver struct {
	Region        Region
	Items         []PresentItem
	Grid          [][]int // -1 for empty, ShapeID for occupied
	Stack         []*SolverFrame
	AllVariations [][]Shape
	Solved        bool
	Failed        bool
}

func NewSolver(region Region, allVariations [][]Shape) *Solver {
	grid := make([][]int, region.Height)
	for i := range grid {
		grid[i] = make([]int, region.Width)
		for j := range grid[i] {
			grid[i][j] = -1
		}
	}

	var items []PresentItem
	totalArea := 0
	for _, pid := range region.Presents {
		if pid >= len(allVariations) || len(allVariations[pid]) == 0 {
			continue
		}
		area := len(allVariations[pid][0].Points)
		items = append(items, PresentItem{ID: pid, Area: area})
		totalArea += area
	}

	// Optimization: Sort items by area (descending)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Area != items[j].Area {
			return items[i].Area > items[j].Area
		}
		return items[i].ID < items[j].ID
	})

	// Initial Frame
	stack := []*SolverFrame{
		{ItemIndex: 0, R: 0, C: 0, VarIdx: 0},
	}

	// Quick fail check
	failed := false
	if totalArea > region.Width*region.Height {
		failed = true
		stack = nil
	}

	return &Solver{
		Region:        region,
		Items:         items,
		Grid:          grid,
		Stack:         stack,
		AllVariations: allVariations,
		Failed:        failed,
	}
}

func (s *Solver) Step(steps int) {
	if s.Solved || s.Failed || len(s.Stack) == 0 {
		return
	}

	for i := 0; i < steps; i++ {
		if len(s.Stack) == 0 {
			s.Failed = true
			return
		}

		frame := s.Stack[len(s.Stack)-1]

		// Success Check
		if frame.ItemIndex >= len(s.Items) {
			s.Solved = true
			return
		}

		// Backtrack if we returned from a child frame
		if frame.Placed {
			s.removeShape(frame.PlacedShape, frame.PlacedX, frame.PlacedY)
			frame.Placed = false
			// Advance to next option
			frame.VarIdx++
		}

		// Iteration Loop
		placed := false

		// Resume loops from frame state
		item := s.Items[frame.ItemIndex]
		vars := s.AllVariations[item.ID]

		for ; frame.R < s.Region.Height; frame.R++ {
			for ; frame.C < s.Region.Width; frame.C++ {
				// Optimization: If cell is occupied, skip
				if s.Grid[frame.R][frame.C] != -1 {
					continue
				}

				for ; frame.VarIdx < len(vars); frame.VarIdx++ {
					v := vars[frame.VarIdx]
					if s.canPlace(v, frame.C, frame.R) {
						s.placeShape(v, frame.C, frame.R)

						// Record what we did so we can undo it
						frame.Placed = true
						frame.PlacedShape = v
						frame.PlacedX = frame.C
						frame.PlacedY = frame.R

						// Push next frame
						nextFrame := &SolverFrame{
							ItemIndex: frame.ItemIndex + 1,
							R:         0,
							C:         0,
							VarIdx:    0,
						}
						s.Stack = append(s.Stack, nextFrame)
						placed = true
						goto StopStep
					}
				}
				frame.VarIdx = 0 // Reset variation for next cell
			}
			frame.C = 0 // Reset column for next row
		}

	StopStep:
		if !placed {
			// If we exit the loops without placing, this frame failed.
			// Pop it to return to parent.
			s.Stack = s.Stack[:len(s.Stack)-1]
		} else {
			// Verify if we are done immediately after push
			if s.Stack[len(s.Stack)-1].ItemIndex >= len(s.Items) {
				s.Solved = true
				return
			}
			// If we successfully placed, we break the inner step loop to let the UI update
			// or continue depending on 'steps' parameter.
			// For visualization smoothness, we might want to return here.
			// return
		}
	}
}

func (s *Solver) canPlace(shp Shape, x, y int) bool {
	if x+shp.Width > s.Region.Width || y+shp.Height > s.Region.Height {
		return false
	}
	for _, p := range shp.Points {
		px, py := x+p.X, y+p.Y
		if px < 0 || px >= s.Region.Width || py < 0 || py >= s.Region.Height {
			return false
		}
		if s.Grid[py][px] != -1 {
			return false
		}
	}
	return true
}

func (s *Solver) placeShape(shp Shape, x, y int) {
	for _, p := range shp.Points {
		s.Grid[y+p.Y][x+p.X] = shp.ID
	}
}

func (s *Solver) removeShape(shp Shape, x, y int) {
	for _, p := range shp.Points {
		s.Grid[y+p.Y][x+p.X] = -1
	}
}

// --- Game ---

type Game struct {
	Shapes     []Shape
	Regions    []Region
	Variations [][]Shape

	CurrentRegionIdx int
	Solver           *Solver

	SolvedCount int

	// Speed control
	StepsPerTick int
	AutoPlay     bool
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.AutoPlay = !g.AutoPlay
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.StepsPerTick *= 2
		if g.StepsPerTick == 0 {
			g.StepsPerTick = 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.StepsPerTick /= 2
		if g.StepsPerTick < 1 {
			g.StepsPerTick = 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		// Force next region
		g.nextRegion()
	}

	if g.Solver == nil {
		if g.CurrentRegionIdx < len(g.Regions) {
			g.Solver = NewSolver(g.Regions[g.CurrentRegionIdx], g.Variations)
		}
	}

	if g.Solver != nil {
		if !g.Solver.Solved && !g.Solver.Failed {
			g.Solver.Step(g.StepsPerTick)
		} else {
			// Move to next region automatically after a short delay or if fast fwd?
			// For now, just count and move on immediately if autoplay is on
			if g.Solver.Solved {
				g.SolvedCount++
			}
			if g.AutoPlay {
				g.nextRegion()
			}
		}
	}

	return nil
}

func (g *Game) nextRegion() {
	g.CurrentRegionIdx++
	if g.CurrentRegionIdx < len(g.Regions) {
		g.Solver = NewSolver(g.Regions[g.CurrentRegionIdx], g.Variations)
	} else {
		g.Solver = nil // Done with all
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Layout info
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	if g.Solver == nil {
		msg := fmt.Sprintf("All Done!\nTotal Solved: %d", g.SolvedCount)
		ebitenutil.DebugPrint(screen, msg)
		return
	}

	// Sidebar Info
	status := "Running"
	if g.Solver.Solved {
		status = "SOLVED!"
	}
	if g.Solver.Failed {
		status = "Impossible"
	}

	info := fmt.Sprintf(`Region: %d / %d
Size: %dx%d
Solved Count: %d
Status: %s
Speed: %d (Left/Right)
AutoPlay: %v (Space)
(N) Next Region`,
		g.CurrentRegionIdx+1, len(g.Regions),
		g.Solver.Region.Width, g.Solver.Region.Height,
		g.SolvedCount,
		status,
		g.StepsPerTick,
		g.AutoPlay,
	)

	ebitenutil.DebugPrint(screen, info)

	// Draw Grid
	gridW, gridH := g.Solver.Region.Width, g.Solver.Region.Height
	availW := float32(sw - 200)
	availH := float32(sh - 40)

	cellSize := min(availW/float32(gridW), availH/float32(gridH))
	if cellSize > 50 {
		cellSize = 50
	}

	startX := float32(200)
	startY := float32(20)

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			cx, cy := startX+float32(x)*cellSize, startY+float32(y)*cellSize

			// Draw cell background
			vector.DrawFilledRect(screen, cx, cy, cellSize-1, cellSize-1, color.RGBA{40, 40, 50, 255}, false)

			val := g.Solver.Grid[y][x]
			if val != -1 {
				// Shape color
				c := getColor(val)
				vector.DrawFilledRect(screen, cx, cy, cellSize-1, cellSize-1, c, false)
			}
		}
	}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func getColor(id int) color.RGBA {
	h := float64(id) * 137.508 // Golden angle approximation
	r, g, b := hslToRgb(h, 0.7, 0.6)
	return color.RGBA{r, g, b, 255}
}

// Simple HSL to RGB helper
func hslToRgb(h, s, l float64) (uint8, uint8, uint8) {
	c := (1 - abs(2*l-1)) * s
	x := c * (1 - abs(mod(h/60, 2)-1))
	m := l - c/2
	var r, g, b float64

	h = mod(h, 360)
	if h < 60 {
		r, g, b = c, x, 0
	} else if h < 120 {
		r, g, b = x, c, 0
	} else if h < 180 {
		r, g, b = 0, c, x
	} else if h < 240 {
		r, g, b = 0, x, c
	} else if h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
func mod(a, b float64) float64 { return a - float64(int(a/b))*b }

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func canFit(region Region, allVariations [][]Shape) bool {
	solver := NewSolver(region, allVariations)
	// Run until solved or failed
	for !solver.Solved && !solver.Failed {
		solver.Step(1000) // Large step count to speed up tests
	}
	return solver.Solved
}

func main() {
	inputPath := "../../inputs/day12.txt"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	contentBytes, err := os.ReadFile(inputPath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return
	}
	content := string(contentBytes)

	shapes, regions := parseInput(content)

	// Precompute variations
	maxID := 0
	for _, s := range shapes {
		if s.ID > maxID {
			maxID = s.ID
		}
	}
	variations := make([][]Shape, maxID+1)
	for _, s := range shapes {
		variations[s.ID] = generateVariations(s)
	}

	game := &Game{
		Shapes:       shapes,
		Regions:      regions,
		Variations:   variations,
		StepsPerTick: 10, // Fast by default
		AutoPlay:     true,
	}

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Advent of Code 2025 - Day 12 Visualization")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// --- Original Parsing Logic ---

func parseInput(content string) ([]Shape, []Region) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	blocks := strings.Split(strings.TrimSpace(content), "\n\n")
	var shapes []Shape
	var regions []Region

	for _, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		firstLine := lines[0]

		isRegionBlock := false
		if strings.Contains(firstLine, ":") {
			parts := strings.Split(firstLine, ":")
			if strings.Contains(parts[0], "x") {
				isRegionBlock = true
			}
		}

		if isRegionBlock {
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				regions = append(regions, parseRegion(line))
			}
		} else if strings.HasSuffix(strings.TrimSpace(firstLine), ":") {
			idStr := strings.TrimSuffix(strings.TrimSpace(firstLine), ":")
			id, _ := strconv.Atoi(idStr)
			var points []Point
			for y, line := range lines[1:] {
				for x, char := range line {
					if char == '#' {
						points = append(points, Point{x, y})
					}
				}
			}
			shapes = append(shapes, normalizeShape(id, points))
		}
	}

	return shapes, regions
}

func parseRegion(line string) Region {
	parts := strings.Split(line, ":")
	dims := strings.Split(parts[0], "x")
	w, _ := strconv.Atoi(dims[0])
	h, _ := strconv.Atoi(dims[1])

	countsStr := strings.Fields(parts[1])
	var presents []int
	for id, s := range countsStr {
		count, _ := strconv.Atoi(s)
		for k := 0; k < count; k++ {
			presents = append(presents, id)
		}
	}
	return Region{Width: w, Height: h, Presents: presents}
}

func normalizeShape(id int, points []Point) Shape {
	if len(points) == 0 {
		return Shape{ID: id, Points: nil, Width: 0, Height: 0}
	}
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	var newPoints []Point
	for _, p := range points {
		newPoints = append(newPoints, Point{p.X - minX, p.Y - minY})
	}

	// Sort points
	sort.Slice(newPoints, func(i, j int) bool {
		if newPoints[i].Y == newPoints[j].Y {
			return newPoints[i].X < newPoints[j].X
		}
		return newPoints[i].Y < newPoints[j].Y
	})

	return Shape{ID: id, Points: newPoints, Width: (maxX - minX) + 1, Height: (maxY - minY) + 1}
}

func generateVariations(base Shape) []Shape {
	unique := make(map[string]Shape)

	add := func(s Shape) {
		norm := normalizeShape(s.ID, s.Points)
		key := fmt.Sprintf("%v", norm.Points)
		unique[key] = norm
	}

	temp := base
	for i := 0; i < 4; i++ {
		add(temp) // Rotation i

		// Flip Horizontal
		var flipped []Point
		for _, p := range temp.Points {
			flipped = append(flipped, Point{-p.X, p.Y})
		}
		add(normalizeShape(base.ID, flipped))

		// Rotate 90 deg clockwise for next iteration
		var rotated []Point
		for _, p := range temp.Points {
			rotated = append(rotated, Point{-p.Y, p.X})
		}
		temp = normalizeShape(base.ID, rotated)
	}

	var res []Shape
	for _, s := range unique {
		res = append(res, s)
	}
	return res
}
