use std::fs;

#[derive(Debug, Clone, Copy)]
struct Range {
    start: u64,
    end: u64,
}

impl Range {
    fn contains(&self, id: u64) -> bool {
        id >= self.start && id <= self.end
    }
}

fn parse_input(input: &str) -> (Vec<Range>, Vec<u64>) {
    let parts: Vec<&str> = input.split("\n\n").collect();

    let ranges: Vec<Range> = parts[0]
        .lines()
        .map(|line| {
            let parts: Vec<&str> = line.split('-').collect();
            Range {
                start: parts[0].parse().unwrap(),
                end: parts[1].parse().unwrap(),
            }
        })
        .collect();

    let ids: Vec<u64> = parts[1]
        .lines()
        .map(|line| line.parse().unwrap())
        .collect();

    (ranges, ids)
}

fn is_fresh(id: u64, ranges: &[Range]) -> bool {
    ranges.iter().any(|range| range.contains(id))
}

fn part1(input: &str) -> usize {
    let (ranges, ids) = parse_input(input.trim());

    ids.iter()
        .filter(|&&id| is_fresh(id, &ranges))
        .count()
}

fn merge_ranges(ranges: &[Range]) -> Vec<Range> {
    if ranges.is_empty() {
        return Vec::new();
    }

    let mut sorted = ranges.to_vec();
    sorted.sort_by_key(|r| r.start);

    let mut merged = vec![sorted[0]];

    for range in sorted.iter().skip(1) {
        let last = merged.last_mut().unwrap();

        // Check if ranges overlap or are adjacent
        if range.start <= last.end + 1 {
            // Merge by extending the end
            last.end = last.end.max(range.end);
        } else {
            // No overlap, add as new range
            merged.push(*range);
        }
    }

    merged
}

fn count_ids_in_ranges(ranges: &[Range]) -> u64 {
    let merged = merge_ranges(ranges);
    merged.iter().map(|r| r.end - r.start + 1).sum()
}

fn part2(input: &str) -> u64 {
    let (ranges, _ids) = parse_input(input.trim());
    count_ids_in_ranges(&ranges)
}

fn main() {
    let input = fs::read_to_string("../inputs/day05.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
3-5
10-14
16-20
12-18

1
5
8
11
17
32";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 3);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 14);
    }
}
