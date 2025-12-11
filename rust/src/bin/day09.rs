use std::collections::{HashMap, HashSet, VecDeque};
use std::fs;

fn parse_positions(input: &str) -> Vec<(i32, i32)> {
    input
        .lines()
        .map(|line| {
            let parts: Vec<i32> = line.split(',').map(|s| s.trim().parse().unwrap()).collect();
            (parts[0], parts[1])
        })
        .collect()
}

fn part1(input: &str) -> i64 {
    let positions = parse_positions(input.trim());
    let n = positions.len();

    let mut max_area = 0_i64;

    // Try all pairs of red tiles as opposite corners
    for i in 0..n {
        for j in (i + 1)..n {
            let (x1, y1) = positions[i];
            let (x2, y2) = positions[j];

            // Calculate rectangle dimensions (inclusive of corners)
            let width = ((x2 - x1).abs() + 1) as i64;
            let height = ((y2 - y1).abs() + 1) as i64;
            let area = width * height;

            max_area = max_area.max(area);
        }
    }

    max_area
}

fn part2(input: &str) -> i64 {
    let points = parse_positions(input.trim());
    let n = points.len();

    // 1. Coordinate Compression
    let mut xs: Vec<i32> = points.iter().map(|&(x, _)| x).collect();
    let mut ys: Vec<i32> = points.iter().map(|&(_, y)| y).collect();
    xs.sort();
    xs.dedup();
    ys.sort();
    ys.dedup();

    let x_map: HashMap<i32, usize> = xs.iter().enumerate().map(|(i, &x)| (x, i)).collect();
    let y_map: HashMap<i32, usize> = ys.iter().enumerate().map(|(i, &y)| (y, i)).collect();

    // 2. Create compressed grid (2*N+1 to represent lines and gaps)
    let height = 2 * ys.len() + 1;
    let width = 2 * xs.len() + 1;
    let mut grid = vec![vec![false; width]; height]; // false = outside, true = boundary/inside

    // 3. Draw boundaries on compressed grid
    for k in 0..n {
        let p1 = points[k];
        let p2 = points[(k + 1) % n];

        // Convert to grid indices (2*i + 1 for actual coordinates)
        let c1 = x_map[&p1.0] * 2 + 1;
        let r1 = y_map[&p1.1] * 2 + 1;
        let c2 = x_map[&p2.0] * 2 + 1;
        let r2 = y_map[&p2.1] * 2 + 1;

        let r_min = r1.min(r2);
        let r_max = r1.max(r2);
        let c_min = c1.min(c2);
        let c_max = c1.max(c2);

        for r in r_min..=r_max {
            for c in c_min..=c_max {
                grid[r][c] = true;
            }
        }
    }

    // 4. Flood fill from (0,0) to mark all "outside" cells
    let mut is_outside = HashSet::new();
    let mut queue = VecDeque::new();
    queue.push_back((0, 0));
    is_outside.insert((0, 0));

    while let Some((r, c)) = queue.pop_front() {
        for (dr, dc) in [(-1, 0), (1, 0), (0, -1), (0, 1)] {
            let nr = r as i32 + dr;
            let nc = c as i32 + dc;

            if nr >= 0 && nr < height as i32 && nc >= 0 && nc < width as i32 {
                let nr = nr as usize;
                let nc = nc as usize;

                if !is_outside.contains(&(nr, nc)) && !grid[nr][nc] {
                    is_outside.insert((nr, nc));
                    queue.push_back((nr, nc));
                }
            }
        }
    }

    // 5. Build 2D prefix sum of "outside" cells
    let mut prefix_sum = vec![vec![0i32; width]; height];
    for r in 0..height {
        for c in 0..width {
            let val = if is_outside.contains(&(r, c)) { 1 } else { 0 };
            let top = if r > 0 { prefix_sum[r - 1][c] } else { 0 };
            let left = if c > 0 { prefix_sum[r][c - 1] } else { 0 };
            let top_left = if r > 0 && c > 0 {
                prefix_sum[r - 1][c - 1]
            } else {
                0
            };
            prefix_sum[r][c] = val + top + left - top_left;
        }
    }

    // Helper function to check if a rectangle has no outside cells
    let region_is_clean = |r1: usize, c1: usize, r2: usize, c2: usize| -> bool {
        let mut total = prefix_sum[r2][c2];
        if r1 > 0 {
            total -= prefix_sum[r1 - 1][c2];
        }
        if c1 > 0 {
            total -= prefix_sum[r2][c1 - 1];
        }
        if r1 > 0 && c1 > 0 {
            total += prefix_sum[r1 - 1][c1 - 1];
        }
        total == 0
    };

    // 6. Check all pairs of red tiles
    let mut max_area = 0i64;
    let mapped: Vec<(usize, usize)> = points
        .iter()
        .map(|&(x, y)| (x_map[&x] * 2 + 1, y_map[&y] * 2 + 1))
        .collect();

    for i in 0..n {
        for j in (i + 1)..n {
            let (c1, r1) = mapped[i];
            let (c2, r2) = mapped[j];

            let r_min = r1.min(r2);
            let r_max = r1.max(r2);
            let c_min = c1.min(c2);
            let c_max = c1.max(c2);

            // Check if rectangle on compressed grid has no outside cells
            if region_is_clean(r_min, c_min, r_max, c_max) {
                // Calculate real area using original coordinates
                let width = ((points[i].0 - points[j].0).abs() + 1) as i64;
                let height = ((points[i].1 - points[j].1).abs() + 1) as i64;
                let area = width * height;
                max_area = max_area.max(area);
            }
        }
    }

    max_area
}

fn main() {
    let input = fs::read_to_string("../inputs/day09.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 50_i64);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 24);
    }
}
