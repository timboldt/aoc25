# Advent of Code 2025 - Go Solutions & Visualizations

This repository contains Go solutions for the Advent of Code 2025 challenges. The project is currently being enhanced to include interactive visualizations for each problem using the [Ebitengine](https://ebitengine.org/) 2D game library.

## Project Structure

*   **`dayXX/`**: Directories containing the source code (`main.go`) and tests (`main_test.go`) for each day.
    *   **`day01/` - `day03/`**: Currently feature fully implemented Ebitengine visualizations.
    *   **`day04/` onwards**: Currently contain standard CLI-based solutions.
*   **`go.mod` / `go.sum`**: Go module definitions and dependencies.
*   **`../inputs/`**: (External to this module) Expected location for puzzle input text files (e.g., `day01.txt`).

## Getting Started

### Prerequisites
*   **Go**: Version 1.24 or later.
*   **Dependencies**: Run `go mod tidy` to ensure all dependencies (specifically Ebitengine) are installed.

### Running a Solution/Visualization

**For Visualized Days (Day 1-3):**
These run as a GUI window.

```bash
cd day01
go run main.go
```

**For Standard Days:**
These run in the terminal and output the answers.

```bash
cd day04
go run main.go
```

### Running Tests

To run tests for all days:

```bash
go test ./...
```

To run tests for a specific day:

```bash
cd day01
go test -v
```

## Development & Visualization Guide

When converting a standard solution to a visualization:

1.  **Refactor Logic**: Move the core solving logic into a stateful structure that can be advanced incrementally.
2.  **Implement `ebiten.Game`**:
    *   **`Update()`**: Advance the simulation state (e.g., process one step of the algorithm).
    *   **`Draw(screen *ebiten.Image)`**: Render the current state using `ebitenutil` or `vector` packages.
    *   **`Layout()`**: Define the window dimensions.
3.  **Input Handling**: Ensure the input is read correctly from `../../inputs/dayXX.txt`.
4.  **Verification**: Maintain the original `part1` and `part2` functions to ensure the visualization logic still produces the correct results (verified via `go test`).

## Current Progress

*   **Day 01 (Safe Dial)**: Visualized (Rotating dial).
*   **Day 02 (Pattern Search)**: Visualized (Fast scanning with numbers).
*   **Day 03 (Battery Banks)**: Visualized (Stack-based greedy algorithm animation).
*   **Day 04 (Paper Rolls)**: Visualized (Grid-based removal simulation).
*   **Day 05 (Range Merging)**: Visualized (Range merging logic).
*   **Day 06 (Column Math)**: Visualized (Column-based parsing and solving).
*   **Day 07 (Tachyon Beams)**: Visualized (Beam propagation and splitting).
*   **Day 08 (Junction Boxes)**: Visualized (3D point clustering/MST animation).
*   **Day 09+**: Pending visualization.
