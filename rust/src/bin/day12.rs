use std::collections::HashSet;
use std::fs;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
struct Shape {
    cells: Vec<(i32, i32)>, // Normalized coordinates where shape has '#'
}

impl Shape {
    fn parse(lines: &[&str]) -> Self {
        let mut cells = Vec::new();
        for (y, line) in lines.iter().enumerate() {
            for (x, ch) in line.chars().enumerate() {
                if ch == '#' {
                    cells.push((x as i32, y as i32));
                }
            }
        }
        Self::normalize(cells)
    }

    fn normalize(cells: Vec<(i32, i32)>) -> Self {
        if cells.is_empty() {
            return Shape { cells };
        }

        let min_x = cells.iter().map(|(x, _)| *x).min().unwrap();
        let min_y = cells.iter().map(|(_, y)| *y).min().unwrap();

        let cells: Vec<_> = cells.iter().map(|(x, y)| (x - min_x, y - min_y)).collect();

        Shape { cells }
    }

    fn rotate_90(&self) -> Self {
        let cells = self.cells.iter().map(|(x, y)| (-*y, *x)).collect();
        Self::normalize(cells)
    }

    fn flip_horizontal(&self) -> Self {
        let cells = self.cells.iter().map(|(x, y)| (-*x, *y)).collect();
        Self::normalize(cells)
    }

    fn all_orientations(&self) -> Vec<Shape> {
        let mut orientations = HashSet::new();
        let mut current = self.clone();

        // 4 rotations
        for _ in 0..4 {
            orientations.insert(current.clone());
            current = current.rotate_90();
        }

        // Flip and 4 more rotations
        current = self.flip_horizontal();
        for _ in 0..4 {
            orientations.insert(current.clone());
            current = current.rotate_90();
        }

        orientations.into_iter().collect()
    }

    fn can_place(&self, grid: &[Vec<bool>], row: i32, col: i32) -> bool {
        let height = grid.len() as i32;
        let width = grid[0].len() as i32;

        for (dx, dy) in &self.cells {
            let x = col + dx;
            let y = row + dy;

            if x < 0 || x >= width || y < 0 || y >= height {
                return false;
            }

            if grid[y as usize][x as usize] {
                return false;
            }
        }

        true
    }

    fn place(&self, grid: &mut [Vec<bool>], row: i32, col: i32) {
        for (dx, dy) in &self.cells {
            let x = (col + dx) as usize;
            let y = (row + dy) as usize;
            grid[y][x] = true;
        }
    }

    fn remove(&self, grid: &mut [Vec<bool>], row: i32, col: i32) {
        for (dx, dy) in &self.cells {
            let x = (col + dx) as usize;
            let y = (row + dy) as usize;
            grid[y][x] = false;
        }
    }

    fn first_cell(&self) -> (i32, i32) {
        // Find cell with min_y, then min_x
        let mut min = self.cells[0];
        for &cell in &self.cells[1..] {
            if cell.1 < min.1 || (cell.1 == min.1 && cell.0 < min.0) {
                min = cell;
            }
        }
        min
    }
}

#[derive(Debug)]
struct Region {
    width: usize,
    height: usize,
    required_shapes: Vec<usize>, // Count of each shape index
}

fn parse_input(input: &str) -> (Vec<Shape>, Vec<Region>) {
    // Split into lines and find the boundary between shapes and regions
    let lines: Vec<&str> = input.lines().collect();

    let mut shapes = Vec::new();
    let mut regions = Vec::new();

    let mut i = 0;

    // Parse shapes
    while i < lines.len() {
        let line = lines[i].trim();

        // Check if this is a shape definition (starts with number and colon)
        if line.is_empty() {
            i += 1;
            continue;
        }

        // Check if this is a region line (contains 'x' like "4x4:")
        if line.contains('x') && line.contains(':') {
            // This is a region, parse all remaining regions
            while i < lines.len() {
                let region_line = lines[i].trim();
                if !region_line.is_empty() {
                    let parts: Vec<&str> = region_line.split(':').collect();
                    if parts.len() == 2 {
                        let dims: Vec<&str> = parts[0].trim().split('x').collect();
                        let width = dims[0].parse().unwrap();
                        let height = dims[1].parse().unwrap();

                        let counts: Vec<usize> = parts[1]
                            .split_whitespace()
                            .map(|s| s.parse().unwrap())
                            .collect();

                        regions.push(Region {
                            width,
                            height,
                            required_shapes: counts,
                        });
                    }
                }
                i += 1;
            }
            break;
        }

        // This should be a shape definition
        if line.contains(':') {
            // Parse this shape
            i += 1; // Move to first line of shape pattern
            let mut shape_lines = Vec::new();

            while i < lines.len() {
                let shape_line = lines[i];
                if shape_line.trim().is_empty() {
                    break; // End of this shape
                }
                if shape_line.contains(':') {
                    break; // Start of next shape
                }
                if shape_line.contains('x') && shape_line.contains(':') {
                    break; // Start of regions
                }
                shape_lines.push(shape_line);
                i += 1;
            }

            if !shape_lines.is_empty() {
                shapes.push(Shape::parse(&shape_lines));
            }
        } else {
            i += 1;
        }
    }

    (shapes, regions)
}

fn can_fit_all_shapes(region: &Region, shapes: &[Shape], all_orientations: &[Vec<Shape>]) -> bool {
    let mut total_cells_needed = 0;
    let mut counts = vec![0; shapes.len()];

    for (idx, &count) in region.required_shapes.iter().enumerate() {
        if count > 0 {
            if idx >= shapes.len() {
                return false;
            }
            total_cells_needed += shapes[idx].cells.len() * count;
            counts[idx] = count;
        }
    }

    let total_cells = region.width * region.height;
    if total_cells_needed > total_cells {
        return false;
    }

    let gaps_allowed = total_cells - total_cells_needed;
    let mut grid = vec![vec![false; region.width]; region.height];

    // Sort shape indices by size (largest first) to try filling with big chunks first
    let mut shape_indices: Vec<usize> = (0..shapes.len()).collect();
    shape_indices.sort_by_key(|&i| std::cmp::Reverse(shapes[i].cells.len()));

    solve(
        &mut grid,
        &mut counts,
        &shape_indices,
        all_orientations,
        gaps_allowed,
    )
}

fn solve(
    grid: &mut [Vec<bool>],
    remaining_counts: &mut [usize],
    shape_indices: &[usize],
    all_orientations: &[Vec<Shape>],
    gaps_remaining: usize,
) -> bool {
    // Find first empty cell
    let mut target_r = 0;
    let mut target_c = 0;
    let mut found = false;

    'outer: for (r, row) in grid.iter().enumerate() {
        for (c, &cell) in row.iter().enumerate() {
            if !cell {
                target_r = r;
                target_c = c;
                found = true;
                break 'outer;
            }
        }
    }

    if !found {
        return true;
    }

    // Try to cover (target_r, target_c) with a shape
    for &shape_idx in shape_indices {
        if remaining_counts[shape_idx] > 0 {
            for orient in &all_orientations[shape_idx] {
                // Optimization: use first_cell to anchor the shape
                let (fx, fy) = orient.first_cell();

                // Calculate required origin to place `first_cell` at `(target_c, target_r)`
                let origin_col = target_c as i32 - fx;
                let origin_row = target_r as i32 - fy;

                if orient.can_place(grid, origin_row, origin_col) {
                    orient.place(grid, origin_row, origin_col);
                    remaining_counts[shape_idx] -= 1;

                    if solve(
                        grid,
                        remaining_counts,
                        shape_indices,
                        all_orientations,
                        gaps_remaining,
                    ) {
                        return true;
                    }

                    remaining_counts[shape_idx] += 1;
                    orient.remove(grid, origin_row, origin_col);
                }
            }
        }
    }

    // Try to treat as a gap
    if gaps_remaining > 0 {
        grid[target_r][target_c] = true; // Fill with "gap"
        if solve(
            grid,
            remaining_counts,
            shape_indices,
            all_orientations,
            gaps_remaining - 1,
        ) {
            return true;
        }
        grid[target_r][target_c] = false; // Backtrack
    }

    false
}

fn part1(input: &str) -> usize {
    let (shapes, regions) = parse_input(input);

    // Precompute all orientations
    let all_orientations: Vec<Vec<Shape>> = shapes
        .iter()
        .map(|shape| shape.all_orientations())
        .collect();

    regions
        .iter()
        .filter(|region| can_fit_all_shapes(region, &shapes, &all_orientations))
        .count()
}

fn part2(_input: &str) -> usize {
    0 // Part 2 not yet available
}

fn main() {
    let input = fs::read_to_string("../inputs/day12.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2";

    #[test]
    fn test_part1() {
        let result = part1(EXAMPLE);
        eprintln!("Part 1 result: {}", result);
        assert_eq!(result, 2);
    }

    // #[test]
    // fn test_shape4_orientations() {
    //     let shape4 = Shape::parse(&["###", "#..", "###"]);
    //     eprintln!("Shape 4: {:?}", shape4.cells);

    //     let orientations = shape4.all_orientations();
    //     eprintln!("Shape 4 has {} unique orientations", orientations.len());

    //     for (i, orient) in orientations.iter().enumerate() {
    //         eprintln!("  Orientation {}: {:?}", i, orient.cells);

    //         // Visualize
    //         let max_x = orient.cells.iter().map(|(x, _)| *x).max().unwrap_or(0);
    //         let max_y = orient.cells.iter().map(|(_, y)| *y).max().unwrap_or(0);

    //         for y in 0..=max_y {
    //             let mut row = String::new();
    //             for x in 0..=max_x {
    //                 if orient.cells.contains(&(x, y)) {
    //                     row.push('#');
    //                 } else {
    //                     row.push('.');
    //                 }
    //             }
    //             eprintln!("    {}", row);
    //         }
    //     }

    //     // Test if two shape 4's can fit in a 4x4 grid
    //     let all_orientations = vec![orientations];
    //     let shapes_vec = vec![shape4];
    //     let region = Region {
    //         width: 4,
    //         height: 4,
    //         required_shapes: vec![2], // 2 of shape 0 (which is shape 4)
    //     };

    //     let fits = can_fit_all_shapes(&region, &shapes_vec, &all_orientations);
    //     eprintln!("Two shape 4's fit in 4x4: {}", fits);
    //     assert!(fits, "Two shape 4's should fit in a 4x4 grid");
    // }
}
