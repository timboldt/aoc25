# Advent of Code 2025

Solutions for Advent of Code 2025 in Rust, Zig, and Go.

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
