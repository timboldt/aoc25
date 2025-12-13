use std::collections::HashSet;
use std::fs;

fn parse_grid(input: &str) -> Vec<Vec<char>> {
    input.lines().map(|line| line.chars().collect()).collect()
}

fn find_start(grid: &[Vec<char>]) -> Option<(usize, usize)> {
    for (row, line) in grid.iter().enumerate() {
        for (col, &ch) in line.iter().enumerate() {
            if ch == 'S' {
                return Some((row, col));
            }
        }
    }
    None
}

fn simulate_beams(grid: &[Vec<char>]) -> usize {
    let (start_row, start_col) = find_start(grid).expect("No start position found");
    let height = grid.len();
    let width = if height > 0 { grid[0].len() } else { 0 };

    let mut split_count = 0;
    let mut current_beams: HashSet<usize> = HashSet::new();
    current_beams.insert(start_col);

    // Simulate beam movement row by row
    for row_data in grid.iter().skip(start_row + 1).take(height - start_row - 1) {
        let mut next_beams: HashSet<usize> = HashSet::new();

        for &col in &current_beams {
            if col >= width {
                continue;
            }

            if row_data[col] == '^' {
                // Beam hits a splitter
                split_count += 1;

                // Create two new beams to the left and right
                if col > 0 {
                    next_beams.insert(col - 1);
                }
                if col + 1 < width {
                    next_beams.insert(col + 1);
                }
            } else {
                // Beam continues downward
                next_beams.insert(col);
            }
        }

        current_beams = next_beams;

        // If no beams remain, we're done
        if current_beams.is_empty() {
            break;
        }
    }

    split_count
}

fn part1(input: &str) -> usize {
    let grid = parse_grid(input.trim());
    simulate_beams(&grid)
}

fn count_timelines_recursive(
    grid: &[Vec<char>],
    row: usize,
    col: usize,
    memo: &mut std::collections::HashMap<(usize, usize), usize>,
) -> usize {
    let height = grid.len();
    let width = if height > 0 { grid[0].len() } else { 0 };

    // Check if we've exited the grid (successfully completed a timeline)
    if row >= height {
        return 1;
    }

    // Check if out of bounds horizontally (timeline terminates)
    if col >= width {
        return 0;
    }

    // Check memoization cache
    if let Some(&count) = memo.get(&(row, col)) {
        return count;
    }

    let count = if grid[row][col] == '^' {
        // Particle splits into two timelines at the splitter
        let left = if col > 0 {
            count_timelines_recursive(grid, row + 1, col - 1, memo)
        } else {
            0
        };
        let right = if col + 1 < width {
            count_timelines_recursive(grid, row + 1, col + 1, memo)
        } else {
            0
        };
        left + right
    } else {
        // Continue downward in the same timeline
        count_timelines_recursive(grid, row + 1, col, memo)
    };

    memo.insert((row, col), count);
    count
}

fn count_timelines(grid: &[Vec<char>]) -> usize {
    let (start_row, start_col) = find_start(grid).expect("No start position found");
    let mut memo = std::collections::HashMap::new();
    count_timelines_recursive(grid, start_row + 1, start_col, &mut memo)
}

fn part2(input: &str) -> usize {
    let grid = parse_grid(input.trim());
    count_timelines(&grid)
}

fn main() {
    let input = fs::read_to_string("../inputs/day07.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 21);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 40);
    }
}
