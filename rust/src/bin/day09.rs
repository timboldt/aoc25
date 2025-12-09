use std::collections::HashSet;
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

fn is_inside_polygon(point: (i32, i32), polygon: &[(i32, i32)]) -> bool {
    // Ray casting algorithm: count how many times a ray crosses the polygon boundary
    let (x, y) = point;
    let n = polygon.len();
    let mut inside = false;
    let mut j = n - 1;

    for i in 0..n {
        let (xi, yi) = polygon[i];
        let (xj, yj) = polygon[j];

        if ((yi > y) != (yj > y)) && (x < (xj - xi) * (y - yi) / (yj - yi) + xi) {
            inside = !inside;
        }

        j = i;
    }

    inside
}

fn part2(input: &str) -> i64 {
    let red_tiles = parse_positions(input.trim());
    let n = red_tiles.len();

    // Create a set of red tiles for quick lookup
    let red_set: HashSet<_> = red_tiles.iter().cloned().collect();

    // Build set of boundary tiles (edges between consecutive red tiles)
    let mut boundary_tiles = HashSet::new();
    for i in 0..n {
        let next = (i + 1) % n;
        let (x1, y1) = red_tiles[i];
        let (x2, y2) = red_tiles[next];

        if y1 == y2 {
            // Same row - horizontal line
            let min_x = x1.min(x2);
            let max_x = x1.max(x2);
            for x in min_x..=max_x {
                boundary_tiles.insert((x, y1));
            }
        } else if x1 == x2 {
            // Same column - vertical line
            let min_y = y1.min(y2);
            let max_y = y1.max(y2);
            for y in min_y..=max_y {
                boundary_tiles.insert((x1, y));
            }
        }
    }

    // Check if a tile is valid (red, green boundary, or inside polygon)
    let is_valid_tile = |point: (i32, i32)| -> bool {
        red_set.contains(&point)
            || boundary_tiles.contains(&point)
            || is_inside_polygon(point, &red_tiles)
    };

    // Build list of all candidate rectangles with their areas
    let mut candidates = Vec::new();
    for i in 0..n {
        for j in (i + 1)..n {
            let (x1, y1) = red_tiles[i];
            let (x2, y2) = red_tiles[j];

            let min_x = x1.min(x2);
            let max_x = x1.max(x2);
            let min_y = y1.min(y2);
            let max_y = y1.max(y2);

            let width = (max_x - min_x + 1) as i64;
            let height = (max_y - min_y + 1) as i64;
            let area = width * height;

            candidates.push((area, min_x, max_x, min_y, max_y));
        }
    }

    // Sort by area descending
    candidates.sort_by(|a, b| b.0.cmp(&a.0));

    // Check rectangles from largest to smallest
    let mut checked_count = 0;
    for &(area, min_x, max_x, min_y, max_y) in &candidates {
        // Skip rectangles that are too large (would take forever to check)
        if area > 100_000_000 {
            continue;
        }

        // Check if all tiles in rectangle are valid
        let mut valid = true;
        'outer: for x in min_x..=max_x {
            for y in min_y..=max_y {
                if !is_valid_tile((x, y)) {
                    valid = false;
                    break 'outer;
                }
            }
        }

        if valid {
            eprintln!("Found valid rectangle with area {} after checking {} candidates", area, checked_count + 1);
            return area;
        }

        checked_count += 1;
        if checked_count % 1000 == 0 {
            eprintln!("Checked {} candidates, current area: {}", checked_count, area);
        }
    }

    0
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
