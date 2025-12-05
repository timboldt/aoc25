use nom::{
    IResult,
    bytes::complete::tag,
    character::complete::{char, multispace0, u64 as nom_u64},
    combinator::map,
    multi::separated_list1,
    sequence::{delimited, separated_pair},
};
use std::fs;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Range {
    start: u64,
    end: u64,
}

fn parse_range(input: &str) -> IResult<&str, Range> {
    map(
        separated_pair(nom_u64, char('-'), nom_u64),
        |(start, end)| Range { start, end },
    )(input)
}

fn parse_ranges(input: &str) -> Vec<Range> {
    let input = input.trim();
    // Parse comma with optional whitespace around it
    let result: IResult<&str, Vec<Range>> =
        separated_list1(delimited(multispace0, tag(","), multispace0), parse_range)(input);
    match result {
        Ok((_, ranges)) => ranges,
        Err(e) => panic!("Failed to parse input: {}", e),
    }
}

fn is_repeated_pattern(n: u64) -> bool {
    let s = n.to_string();
    let len = s.len();

    // Must have even length
    if !len.is_multiple_of(2) {
        return false;
    }

    // Split in half
    let mid = len / 2;
    let first_half = &s[..mid];
    let second_half = &s[mid..];

    // Check no leading zeros
    if first_half.starts_with('0') {
        return false;
    }

    // Check if halves are equal
    first_half == second_half
}

fn part1(input: &str) -> u64 {
    let ranges = parse_ranges(input);
    let mut sum = 0;

    for range in ranges {
        for id in range.start..=range.end {
            if is_repeated_pattern(id) {
                sum += id;
            }
        }
    }

    sum
}

fn is_repeated_pattern_v2(n: u64) -> bool {
    let s = n.to_string();
    let len = s.len();

    // Try all possible pattern lengths from 1 to len/2
    for pattern_len in 1..=(len / 2) {
        // Check if len is divisible by pattern_len (pattern must repeat evenly)
        if len.is_multiple_of(pattern_len) {
            let pattern = &s[..pattern_len];

            // Check for leading zeros
            if pattern.starts_with('0') {
                continue;
            }

            // Check if all chunks match the pattern
            if s.as_bytes()
                .chunks(pattern_len)
                .all(|chunk| chunk == pattern.as_bytes())
            {
                return true;
            }
        }
    }

    false
}

fn part2(input: &str) -> u64 {
    let ranges = parse_ranges(input);
    let mut sum = 0;

    for range in ranges {
        for id in range.start..=range.end {
            if is_repeated_pattern_v2(id) {
                sum += id;
            }
        }
    }

    sum
}

fn main() {
    let input = fs::read_to_string("../inputs/day02.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
11-22,95-115,998-1012,1188511880-1188511890,222220-222224,
1698522-1698528,446443-446449,38593856-38593862,565653-565659,
824824821-824824827,2121212118-2121212124
";

    #[test]
    fn test_is_repeated_pattern() {
        assert!(is_repeated_pattern(11));
        assert!(is_repeated_pattern(22));
        assert!(is_repeated_pattern(99));
        assert!(is_repeated_pattern(1010));
        assert!(is_repeated_pattern(6464));
        assert!(is_repeated_pattern(123123));
        assert!(is_repeated_pattern(222222));
        assert!(is_repeated_pattern(446446));
        assert!(is_repeated_pattern(1188511885));
        assert!(is_repeated_pattern(38593859));

        assert!(!is_repeated_pattern(101));
        assert!(!is_repeated_pattern(1698522));
        assert!(!is_repeated_pattern(0101)); // This is actually 101, not valid
    }

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 1227775554);
    }

    #[test]
    fn test_is_repeated_pattern_v2() {
        // Original patterns (2 repetitions)
        assert!(is_repeated_pattern_v2(11));
        assert!(is_repeated_pattern_v2(99));
        assert!(is_repeated_pattern_v2(1010));
        assert!(is_repeated_pattern_v2(222222));
        assert!(is_repeated_pattern_v2(446446));
        assert!(is_repeated_pattern_v2(1188511885));
        assert!(is_repeated_pattern_v2(38593859));

        // New patterns (3+ repetitions)
        assert!(is_repeated_pattern_v2(111)); // "1" x3
        assert!(is_repeated_pattern_v2(999)); // "9" x3
        assert!(is_repeated_pattern_v2(565656)); // "56" x3
        assert!(is_repeated_pattern_v2(824824824)); // "824" x3
        assert!(is_repeated_pattern_v2(2121212121)); // "21" x5
        assert!(is_repeated_pattern_v2(123123123)); // "123" x3
        assert!(is_repeated_pattern_v2(1212121212)); // "12" x5
        assert!(is_repeated_pattern_v2(1111111)); // "1" x7

        // Not repeated patterns
        assert!(!is_repeated_pattern_v2(101));
        assert!(!is_repeated_pattern_v2(1698522));
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 4174379265);
    }
}
