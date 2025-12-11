package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func main() {
	path := "../../inputs/day11.txt"
	graph := parseInput(path)

	fmt.Printf("Part 1: %d\n", part1(graph))
	fmt.Printf("Part 2: %d\n", part2(graph))

	game := NewGame(graph)
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("AoC 2025 Day 11 - Reactor Paths")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// --- Logic ---

type Graph map[string][]string

func parseInput(path string) Graph {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	graph := make(Graph)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		node := parts[0]
		neighbors := strings.Fields(parts[1])
		graph[node] = neighbors
	}
	return graph
}

func countPaths(g Graph, current, target string, memo map[string]int) int {
	if current == target {
		return 1
	}
	if val, ok := memo[current]; ok {
		return val
	}

	count := 0
	for _, neighbor := range g[current] {
		count += countPaths(g, neighbor, target, memo)
	}
	memo[current] = count
	return count
}

func part1(g Graph) int {
	memo := make(map[string]int)
	return countPaths(g, "you", "out", memo)
}

func part2(g Graph) int {
	count := func(start, end string) int {
		return countPaths(g, start, end, make(map[string]int))
	}

	sd := count("svr", "dac")
	df := count("dac", "fft")
	fo := count("fft", "out")
	path1 := sd * df * fo

	sf := count("svr", "fft")
	fd := count("fft", "dac")
	do := count("dac", "out")
	path2 := sf * fd * do

	return path1 + path2
}

// --- Visualization ---

type Node struct {
	Name  string
	X, Y  float32
	Layer int
}

type Particle struct {
	From  string
	To    string
	T     float32
	Speed float32
	Color color.RGBA
}

type Game struct {
	graph     Graph
	nodes     map[string]*Node
	particles []*Particle
	mode      int // 0 = Part 1, 1 = Part 2
	tick      int

	// Layout params
	minX, maxX float32
	minY, maxY float32
}

func NewGame(g Graph) *Game {
	game := &Game{
		graph: g,
		nodes: make(map[string]*Node),
		mode:  0,
	}
	game.computeLayout()
	return game
}

func (g *Game) computeLayout() {
	// Simple layer-based layout
	// Calculate in-degrees to find roots, but since we want layers from "start" (left),
	// we can try to do a BFS/longest path from 'svr' and 'you'.

	// 1. Build Reverse Graph to find depths? Or just topological sort.

	// Better approach for visualization:
	// Topological sort to assign X layers.
	// Since it's a DAG, this is valid.

	inDegree := make(map[string]int)
	for u, neighbors := range g.graph {
		if _, ok := inDegree[u]; !ok {
			inDegree[u] = 0 // Ensure existence
		}
		for _, v := range neighbors {
			inDegree[v]++
		}
	}

	// Kahn's algorithm for layers
	queue := []string{}
	for u := range inDegree {
		if inDegree[u] == 0 {
			queue = append(queue, u)
		}
	}

	// Sort queue for deterministic layout
	sort.Strings(queue)

	maxLayer := 0

	// Process by layers
	// Note: Kahn's usually gives a linear order. We want parallel layers.
	// Using "Longest Path Layering" (rank = max rank of parents + 1) is better for DAG drawing.
	// To do that, we need to process nodes in topological order.

	// Re-do with topo order first to ensure we process parents before children
	topoOrder := []string{}
	tempInDegree := make(map[string]int)
	for k, v := range inDegree {
		tempInDegree[k] = v
	}
	q := append([]string{}, queue...)

	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		topoOrder = append(topoOrder, u)

		for _, v := range g.graph[u] {
			tempInDegree[v]--
			if tempInDegree[v] == 0 {
				q = append(q, v)
			}
		}
	}

	// Now assign ranks based on parents
	// We need parents for this, or just push rank to children.
	nodeRank := make(map[string]int)
	for _, u := range topoOrder {
		r := nodeRank[u]
		if r > maxLayer {
			maxLayer = r
		}
		for _, v := range g.graph[u] {
			if nodeRank[v] < r+1 {
				nodeRank[v] = r + 1
			}
		}
	}

	// Assign nodes to ranks
	nodesByRank := make(map[int][]string)
	for _, u := range topoOrder {
		r := nodeRank[u]
		nodesByRank[r] = append(nodesByRank[r], u)
	}

	// Create Node objects
	// Screen space: 1200 x 700 essentially
	layerWidth := 1100.0 / float32(maxLayer+1)

	for r := 0; r <= maxLayer; r++ {
		nodes := nodesByRank[r]
		sort.Strings(nodes) // Deterministic Y

		layerHeight := float32(650.0)
		spacingY := layerHeight / float32(len(nodes)+1)

		for i, name := range nodes {
			g.nodes[name] = &Node{
				Name:  name,
				X:     50 + float32(r)*layerWidth,
				Y:     50 + float32(i+1)*spacingY,
				Layer: r,
			}
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.mode = 1 - g.mode
		g.particles = nil // Clear particles on mode switch
	}

	g.tick++

	// Spawn particles
	spawnRate := 5 // Every 5 ticks
	if g.tick%spawnRate == 0 {
		startNode := "you"
		targetColor := color.RGBA{100, 200, 255, 200}
		if g.mode == 1 {
			startNode = "svr"
			targetColor = color.RGBA{255, 100, 100, 200}
		}

		if _, ok := g.nodes[startNode]; ok {
			g.particles = append(g.particles, &Particle{
				From:  startNode,
				To:    startNode, // Initially at start
				T:     1.0,       // Ready to move
				Speed: 0.02 + rand.Float32()*0.02,
				Color: targetColor,
			})
		}
	}

	// Update particles
	activeParticles := g.particles[:0]
	for _, p := range g.particles {
		p.T += p.Speed
		if p.T >= 1.0 {
			// Reached destination, pick next
			p.From = p.To
			p.T = 0

			// Logic for picking next node
			// If Part 2, strictly following valid paths is hard without pre-calc.
			// Let's just flow randomly and let them die if they hit a dead end or 'out'.
			// Maybe visually distinguish "valid" particles later? For now, just flow.

			neighbors := g.graph[p.From]
			if len(neighbors) == 0 {
				// Reached end or dead end
				continue // Die
			}

			p.To = neighbors[rand.Intn(len(neighbors))]
		}
		activeParticles = append(activeParticles, p)
	}
	g.particles = activeParticles

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw Edges
	for u, neighbors := range g.graph {
		n1, ok1 := g.nodes[u]
		if !ok1 {
			continue
		}

		for _, v := range neighbors {
			n2, ok2 := g.nodes[v]
			if !ok2 {
				continue
			}

			col := color.RGBA{60, 60, 80, 255}
			vector.StrokeLine(screen, n1.X, n1.Y, n2.X, n2.Y, 1, col, true)
		}
	}

	// Draw Particles
	for _, p := range g.particles {
		n1 := g.nodes[p.From]
		n2 := g.nodes[p.To]
		if n1 == nil || n2 == nil {
			continue
		}

		x := n1.X + (n2.X-n1.X)*p.T
		y := n1.Y + (n2.Y-n1.Y)*p.T

		vector.DrawFilledCircle(screen, x, y, 3, p.Color, true)
	}

	// Draw Nodes
	for name, node := range g.nodes {
		col := color.RGBA{100, 100, 100, 255}
		radius := float32(4.0)

		// Highlight special nodes
		if name == "you" || name == "svr" {
			col = color.RGBA{100, 255, 100, 255}
			radius = 6.0
		} else if name == "out" {
			col = color.RGBA{255, 100, 100, 255}
			radius = 6.0
		} else if name == "dac" || name == "fft" {
			col = color.RGBA{255, 255, 100, 255}
			radius = 5.0
		}

		vector.DrawFilledCircle(screen, node.X, node.Y, radius, col, true)

		// Optional: Draw names on hover or just for special nodes
		if radius > 4.0 {
			ebitenutil.DebugPrintAt(screen, name, int(node.X)-10, int(node.Y)-20)
		}
	}

	// UI
	modeStr := "Part 1 (Start: you)"
	if g.mode == 1 {
		modeStr = "Part 2 (Start: svr)"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: %s (Space to toggle)", modeStr), 10, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 720
}
