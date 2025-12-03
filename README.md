# Advent of Code 2025

Solutions for Advent of Code 2025 in Rust.

## Structure

- `src/bin/` - Daily solutions (day01.rs, day02.rs, etc.)
- `inputs/` - Puzzle inputs (day01.txt, day02.txt, etc.)

## Running Solutions

Run a specific day:
```bash
cargo run --bin day01
```

Run with release optimizations:
```bash
cargo run --release --bin day01
```

## Testing

Run tests for a specific day:
```bash
cargo test --bin day01
```

Run all tests:
```bash
cargo test
```

## Creating a New Day

Copy the template:
```bash
cp src/bin/day01.rs src/bin/dayXX.rs
```

Then update the input file reference and add your solution.
