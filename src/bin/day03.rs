use nom::{
    IResult,
    character::complete::{digit1, line_ending},
    multi::separated_list1,
};
use std::fs;

fn parse_banks(input: &str) -> IResult<&str, Vec<&str>> {
    separated_list1(line_ending, digit1)(input)
}

fn max_joltage(bank: &str) -> u32 {
    let digits: Vec<char> = bank.chars().collect();
    let mut max = 0;

    // Try all pairs (i, j) where i < j
    for i in 0..digits.len() {
        for j in (i + 1)..digits.len() {
            // Form two-digit number from digits[i] and digits[j]
            let joltage = format!("{}{}", digits[i], digits[j])
                .parse::<u32>()
                .unwrap();
            max = max.max(joltage);
        }
    }

    max
}

fn part1(input: &str) -> u32 {
    let (_, banks) = parse_banks(input.trim()).expect("Failed to parse input");
    banks.iter().map(|bank| max_joltage(bank)).sum()
}

fn max_joltage_n(bank: &str, n: usize) -> u64 {
    let digits: Vec<char> = bank.chars().collect();
    let to_remove = digits.len() - n;

    let mut stack = Vec::new();
    let mut removals_left = to_remove;

    // Greedy approach: remove smaller digits when a larger digit comes
    for &digit in &digits {
        while !stack.is_empty() && removals_left > 0 && stack.last().unwrap() < &digit {
            stack.pop();
            removals_left -= 1;
        }
        stack.push(digit);
    }

    // If we still have removals left, remove from the end
    while removals_left > 0 {
        stack.pop();
        removals_left -= 1;
    }

    // Take first n digits and form the number
    let result: String = stack.into_iter().take(n).collect();
    result.parse().unwrap()
}

fn part2(input: &str) -> u64 {
    let (_, banks) = parse_banks(input.trim()).expect("Failed to parse input");
    banks.iter().map(|bank| max_joltage_n(bank, 12)).sum()
}

fn main() {
    let input = fs::read_to_string("inputs/day03.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
987654321111111
811111111111119
234234234234278
818181911112111
";

    #[test]
    fn test_max_joltage() {
        assert_eq!(max_joltage("987654321111111"), 98);
        assert_eq!(max_joltage("811111111111119"), 89);
        assert_eq!(max_joltage("234234234234278"), 78);
        assert_eq!(max_joltage("818181911112111"), 92);
    }

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 357);
    }

    #[test]
    fn test_max_joltage_n() {
        assert_eq!(max_joltage_n("987654321111111", 12), 987654321111);
        assert_eq!(max_joltage_n("811111111111119", 12), 811111111119);
        assert_eq!(max_joltage_n("234234234234278", 12), 434234234278);
        assert_eq!(max_joltage_n("818181911112111", 12), 888911112111);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 3121910778619);
    }
}
