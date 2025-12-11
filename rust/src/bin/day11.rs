use std::collections::HashMap;
use std::fs;

type Graph<'a> = HashMap<&'a str, Vec<&'a str>>;

fn parse_graph(input: &str) -> Graph<'_> {
    let mut graph = HashMap::new();

    for line in input.lines() {
        if line.trim().is_empty() {
            continue;
        }

        let parts: Vec<&str> = line.split(':').collect();
        if parts.len() != 2 {
            continue;
        }

        let device = parts[0].trim();
        let outputs: Vec<&str> = parts[1]
            .split_whitespace()
            .map(|s| s.trim())
            .filter(|s| !s.is_empty())
            .collect();

        graph.insert(device, outputs);
    }

    graph
}

// Memoized path counting
// Key: (current_node, visited_required_bitmask)
// Value: number of paths from current to target
type MemoKey<'a> = (&'a str, u8);

fn count_paths_memo<'a>(
    graph: &Graph<'a>,
    current: &'a str,
    target: &'a str,
    required: &[&'a str],
    visited_mask: u8,
    memo: &mut HashMap<MemoKey<'a>, usize>,
) -> usize {
    // Base case: reached the target
    if current == target {
        // Check if all required nodes were visited
        let all_required_mask = (1 << required.len()) - 1;
        return if visited_mask == all_required_mask {
            1
        } else {
            0
        };
    }

    // Check memo
    let key = (current, visited_mask);
    if let Some(&cached) = memo.get(&key) {
        return cached;
    }

    let mut total_paths = 0;

    // Explore all outputs from current device
    if let Some(outputs) = graph.get(current) {
        for &next in outputs {
            // Update visited mask if next is a required node
            let mut new_mask = visited_mask;
            for (i, &req) in required.iter().enumerate() {
                if next == req {
                    new_mask |= 1 << i;
                }
            }

            total_paths += count_paths_memo(graph, next, target, required, new_mask, memo);
        }
    }

    memo.insert(key, total_paths);
    total_paths
}

fn part1(input: &str) -> usize {
    let graph = parse_graph(input);
    let required: Vec<&str> = vec![];
    let mut memo = HashMap::new();
    count_paths_memo(&graph, "you", "out", &required, 0, &mut memo)
}

fn part2(input: &str) -> usize {
    let graph = parse_graph(input);
    let required = vec!["dac", "fft"];
    let mut memo = HashMap::new();
    count_paths_memo(&graph, "svr", "out", &required, 0, &mut memo)
}

fn main() {
    let input = fs::read_to_string("../inputs/day11.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out";

    #[test]
    fn test_parse() {
        let graph = parse_graph(EXAMPLE);
        assert_eq!(graph.get("you"), Some(&vec!["bbb", "ccc"]));
        assert_eq!(graph.get("bbb"), Some(&vec!["ddd", "eee"]));
    }

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 5);
    }

    const EXAMPLE2: &str = "\
svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out";

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE2), 2);
    }
}
