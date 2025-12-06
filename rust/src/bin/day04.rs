use std::fs;

fn parse_input(input: &str) -> Vec<Vec<char>> {
    input.lines().map(|line| line.chars().collect()).collect()
}

fn count_adjacent_papers(grid: &[Vec<char>], row: usize, col: usize) -> usize {
    let rows = grid.len() as i32;
    let cols = grid[0].len() as i32;
    let directions = [
        (-1, -1), (-1, 0), (-1, 1),
        (0, -1),           (0, 1),
        (1, -1),  (1, 0),  (1, 1),
    ];

    let mut count = 0;
    for (dr, dc) in directions {
        let new_row = row as i32 + dr;
        let new_col = col as i32 + dc;

        if new_row >= 0 && new_row < rows && new_col >= 0 && new_col < cols {
            if grid[new_row as usize][new_col as usize] == '@' {
                count += 1;
            }
        }
    }
    count
}

fn part1(input: &str) -> u32 {
    let grid = parse_input(input.trim());
    let mut accessible = 0;

    for (row, line) in grid.iter().enumerate() {
        for (col, &cell) in line.iter().enumerate() {
            if cell == '@' {
                let adjacent_count = count_adjacent_papers(&grid, row, col);
                if adjacent_count < 4 {
                    accessible += 1;
                }
            }
        }
    }

    accessible
}

fn part2(input: &str) -> u32 {
    let mut grid = parse_input(input.trim());
    let mut total_removed = 0;

    loop {
        // Find all accessible rolls in this iteration
        let mut to_remove = Vec::new();

        for (row, line) in grid.iter().enumerate() {
            for (col, &cell) in line.iter().enumerate() {
                if cell == '@' {
                    let adjacent_count = count_adjacent_papers(&grid, row, col);
                    if adjacent_count < 4 {
                        to_remove.push((row, col));
                    }
                }
            }
        }

        // If no rolls can be removed, we're done
        if to_remove.is_empty() {
            break;
        }

        // Remove all accessible rolls
        for (row, col) in &to_remove {
            grid[*row][*col] = '.';
        }

        total_removed += to_remove.len() as u32;
    }

    total_removed
}

fn main() {
    let input = fs::read_to_string("../inputs/day04.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 13);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 43);
    }
}
