use std::fs;

// Union-Find (Disjoint Set Union) data structure
struct UnionFind {
    parent: Vec<usize>,
    rank: Vec<usize>,
}

impl UnionFind {
    fn new(size: usize) -> Self {
        UnionFind {
            parent: (0..size).collect(),
            rank: vec![0; size],
        }
    }

    fn find(&mut self, x: usize) -> usize {
        if self.parent[x] != x {
            self.parent[x] = self.find(self.parent[x]); // Path compression
        }
        self.parent[x]
    }

    fn union(&mut self, x: usize, y: usize) -> bool {
        let root_x = self.find(x);
        let root_y = self.find(y);

        if root_x == root_y {
            return false; // Already in same set
        }

        // Union by rank
        if self.rank[root_x] < self.rank[root_y] {
            self.parent[root_x] = root_y;
        } else if self.rank[root_x] > self.rank[root_y] {
            self.parent[root_y] = root_x;
        } else {
            self.parent[root_y] = root_x;
            self.rank[root_x] += 1;
        }

        true // Successfully merged
    }

    fn component_sizes(&mut self) -> Vec<usize> {
        let n = self.parent.len();
        let mut sizes = vec![0; n];

        for i in 0..n {
            let root = self.find(i);
            sizes[root] += 1;
        }

        sizes.into_iter().filter(|&s| s > 0).collect()
    }
}

fn parse_positions(input: &str) -> Vec<(i32, i32, i32)> {
    input
        .lines()
        .map(|line| {
            let parts: Vec<i32> = line.split(',').map(|s| s.parse().unwrap()).collect();
            (parts[0], parts[1], parts[2])
        })
        .collect()
}

fn distance_squared(a: (i32, i32, i32), b: (i32, i32, i32)) -> i64 {
    let dx = (a.0 - b.0) as i64;
    let dy = (a.1 - b.1) as i64;
    let dz = (a.2 - b.2) as i64;
    dx * dx + dy * dy + dz * dz
}

fn solve(input: &str, num_connections: usize) -> usize {
    let positions = parse_positions(input.trim());
    let n = positions.len();

    // Generate all edges with distances
    let mut edges: Vec<(i64, usize, usize)> = Vec::new();
    for i in 0..n {
        for j in (i + 1)..n {
            let dist_sq = distance_squared(positions[i], positions[j]);
            edges.push((dist_sq, i, j));
        }
    }

    // Sort by distance
    edges.sort_by_key(|&(dist, _, _)| dist);

    // Try to connect the num_connections closest pairs using Union-Find
    let mut uf = UnionFind::new(n);

    for &(_dist, i, j) in edges.iter().take(num_connections) {
        uf.union(i, j); // Try to connect, even if already connected
    }

    // Get component sizes
    let mut sizes = uf.component_sizes();
    sizes.sort_by(|a, b| b.cmp(a)); // Sort descending

    // Multiply the three largest
    sizes.iter().take(3).product()
}

fn part1(input: &str) -> usize {
    solve(input, 1000)
}

fn part2(input: &str) -> i64 {
    let positions = parse_positions(input.trim());
    let n = positions.len();

    // Generate all edges with distances
    let mut edges: Vec<(i64, usize, usize)> = Vec::new();
    for i in 0..n {
        for j in (i + 1)..n {
            let dist_sq = distance_squared(positions[i], positions[j]);
            edges.push((dist_sq, i, j));
        }
    }

    // Sort by distance
    edges.sort_by_key(|&(dist, _, _)| dist);

    // Connect pairs until all in one component
    let mut uf = UnionFind::new(n);
    let mut num_components = n;
    let mut last_connection = (0, 0);

    for &(_dist, i, j) in &edges {
        if uf.union(i, j) {
            num_components -= 1;
            last_connection = (i, j);

            if num_components == 1 {
                break;
            }
        }
    }

    // Return product of X coordinates
    let (i, j) = last_connection;
    (positions[i].0 as i64) * (positions[j].0 as i64)
}

fn main() {
    let input = fs::read_to_string("../inputs/day08.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689";

    #[test]
    fn test_part1() {
        assert_eq!(solve(EXAMPLE, 10), 40);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 25272_i64);
    }
}
