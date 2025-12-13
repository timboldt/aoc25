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
    // Build list of shapes to place
    let mut shapes_to_place: Vec<usize> = Vec::new();
    for (shape_idx, &count) in region.required_shapes.iter().enumerate() {
        for _ in 0..count {
            shapes_to_place.push(shape_idx);
        }
    }

    if shapes_to_place.is_empty() {
        return true;
    }

    // Sort by size (largest first)
    shapes_to_place.sort_by_key(|&idx| std::cmp::Reverse(shapes[idx].cells.len()));

    let total_cells_needed: usize = shapes_to_place
        .iter()
        .map(|&idx| shapes[idx].cells.len())
        .sum();
    let total_cells = region.width * region.height;

    if total_cells_needed > total_cells {
        return false;
    }

    let mut grid = vec![vec![false; region.width]; region.height];
    backtrack_simple(&mut grid, &mut shapes_to_place, all_orientations)
}

fn backtrack_simple(
    grid: &mut [Vec<bool>],
    shapes_remaining: &mut Vec<usize>,
    all_orientations: &[Vec<Shape>],
) -> bool {
    if shapes_remaining.is_empty() {
        return true;
    }

    let height = grid.len();
    let width = grid[0].len();

    // Try placing the first shape (already sorted by size)
    let shape_idx = shapes_remaining[0];
    let orientations = &all_orientations[shape_idx];

    // Try each position in the grid
    // Only try positions where we can actually place the shape's top-left corner
    for row in 0..height {
        for col in 0..width {
            // Try each orientation of this shape at this position
            for orientation in orientations {
                if orientation.can_place(grid, row as i32, col as i32) {
                    orientation.place(grid, row as i32, col as i32);
                    shapes_remaining.remove(0);

                    if backtrack_simple(grid, shapes_remaining, all_orientations) {
                        shapes_remaining.insert(0, shape_idx);
                        orientation.remove(grid, row as i32, col as i32);
                        return true;
                    }

                    shapes_remaining.insert(0, shape_idx);
                    orientation.remove(grid, row as i32, col as i32);
                }
            }
        }
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

    #[test]
    fn test_shape4_orientations() {
        let shape4 = Shape::parse(&["###", "#..", "###"]);
        eprintln!("Shape 4: {:?}", shape4.cells);

        let orientations = shape4.all_orientations();
        eprintln!("Shape 4 has {} unique orientations", orientations.len());

        for (i, orient) in orientations.iter().enumerate() {
            eprintln!("  Orientation {}: {:?}", i, orient.cells);

            // Visualize
            let max_x = orient.cells.iter().map(|(x, _)| *x).max().unwrap_or(0);
            let max_y = orient.cells.iter().map(|(_, y)| *y).max().unwrap_or(0);

            for y in 0..=max_y {
                let mut row = String::new();
                for x in 0..=max_x {
                    if orient.cells.contains(&(x, y)) {
                        row.push('#');
                    } else {
                        row.push('.');
                    }
                }
                eprintln!("    {}", row);
            }
        }

        // Test if two shape 4's can fit in a 4x4 grid
        let all_orientations = vec![orientations];
        let shapes_vec = vec![shape4];
        let region = Region {
            width: 4,
            height: 4,
            required_shapes: vec![2], // 2 of shape 0 (which is shape 4)
        };

        let fits = can_fit_all_shapes(&region, &shapes_vec, &all_orientations);
        eprintln!("Two shape 4's fit in 4x4: {}", fits);
        assert!(fits, "Two shape 4's should fit in a 4x4 grid");
    }
}
