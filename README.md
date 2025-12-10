# Advent of Code 2025

Solutions for Advent of Code 2025 in Rust (plus a few in Go and Zig)

Every year, I try a different approach in order to, well, learn new approaches.

This year, I'm using (okay, cheating with) AI in order to understand what modern AI can (and cannot) do.

P.S. If you want to learn, don't cheat by letting an LLM write a solution for you. Or at least try first and then if you get stuck, talk your ideas through with an LLM, and then go back and implement the new ideas. The reason I'm using an LLM this year is not to solve the puzzles but to learn about LLMs - this is probably not the reason you are participating in AOC, so don't do what I did.

## AI Performance Comparison

| Day | Claude Sonnet 4.5<br>Implementation<br>In Rust | Gemini 3 Pro<br>Visualization<br>In Go |
|-----|-------------------|--------------|
| 1   | ✅ ✅     | ✅       |
| 2   | ✅ ✅     | ✅       |
| 3   | ✅ ✅     | ✅       |
| 4   | ✅ ✅     | ✅       |
| 5   | ✅ ✅     | ✅       |
| 6   | ✅ ✅     | ✅       |
| 7   | ✅ ✅     | ⬜       |
| 8   | ✅ ✅     | ⬜       |
| 9   | ✅ ⚠️ (1) | ⬜       |
| 10  | ⬜ ⬜     | ⬜.      |
| 11  | ⬜ ⬜     | ⬜       |
| 12  | ⬜ ⬜     | ⬜       |

(1) Claude came up with a correct answer in Rust, but it took forever to run. I had it rewrite it in Go and it worked after a couple of attempts. Then I had Gemini look at the problem and it proposed a lightning fast solution in Python. I then gave Gemini's Python code to Claude and it was able to make the Rust implementation much faster, using Gemini's solution.

## Structure

```
.
├── rust/          # Rust implementations
│   ├── Cargo.toml
│   └── src/bin/   # Day solutions
├── zig/           # Zig implementations
│   ├── build.zig
│   └── src/       # Day solutions
├── go/            # Go implementations
│   ├── go.mod
│   ├── day01/
│   ├── day02/
│   └── day03/
└── inputs/        # Shared puzzle inputs
```

## Running Solutions

### Rust

Run a specific day:
```bash
cd rust
cargo run --bin day01
```

Run with release optimizations:
```bash
cd rust
cargo run --release --bin day01
```

Run tests:
```bash
cd rust
cargo test --bin day01  # Test a specific day
cargo test              # Test all days
```

### Zig

Run a specific day:
```bash
cd zig
zig build day01
```

Run with optimizations:
```bash
cd zig
zig build day01 -Doptimize=ReleaseFast
```

Run tests:
```bash
cd zig
zig build test
```

### Go

Run a specific day (from the `go/` directory):
```bash
cd go
go run ./day01
```

Or from within a day directory:
```bash
cd go/day01
go run .
```

Run tests:
```bash
cd go
go test ./day01    # Test a specific day
go test ./...      # Test all days
```

## Development

Each language follows its own idioms and best practices:
- **Rust**: Uses nom for parsing, bins for executables
- **Zig**: Simple build system with test support
- **Go**: One package per day with standard testing
