# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the Rust implementation directory for Advent of Code 2025 solutions. It's part of a multi-language repo that includes Rust, Zig, and Go implementations. Each day's puzzle has two parts (part1 and part2).

## Project Structure

- `src/bin/`: Each day's solution is a separate binary (day01.rs, day02.rs, etc.)
- `../inputs/`: Shared puzzle inputs for all languages (day01.txt, day02.txt, etc.)
- Input files are read from `../inputs/dayXX.txt` relative to the Rust directory
- Each binary is self-contained with parsing, part1, part2, and tests

## Running Solutions

Run a specific day:
```bash
cargo run --bin day01
```

Run with release optimizations (important for performance-sensitive puzzles):
```bash
cargo run --release --bin day01
```

Run tests for a specific day:
```bash
cargo test --bin day01
```

Run all tests:
```bash
cargo test
```

## Solution Structure Pattern

Each day follows this pattern:
- Parsing functions at the top (often using nom combinators)
- `part1(input: &str)` function for the first puzzle
- `part2(input: &str)` function for the second puzzle
- `main()` reads from `../inputs/dayXX.txt` and prints both parts
- Tests module with example inputs and expected outputs

## Dependencies

- **nom 7.1**: Parser combinator library used for parsing puzzle inputs (see day01.rs for examples)
- **rayon 1.10**: Data parallelism library (available but not used in all solutions)

## Common Patterns

### Input Reading
All solutions read input files from the parent `inputs/` directory:
```rust
fs::read_to_string("../inputs/dayXX.txt")
```

### Parsing with nom
Many solutions use nom for parsing (see day01.rs):
- `alt()` for alternatives
- `value()` for mapping to constants
- `separated_list1()` for lists
- `tuple()` for sequences
- Custom parsers composed from nom primitives

### Testing
Each solution includes a test module with the example from the puzzle description:
```rust
#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "...";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), expected);
    }
}
```

## Performance Notes

- Some puzzles (like day09) include performance-critical code with optimized algorithms
- Use `--release` flag for accurate performance testing
- Solutions may include diagnostic `eprintln!()` statements for debugging
