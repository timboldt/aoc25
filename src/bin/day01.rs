use std::fs;

fn part1(input: &str) -> i64 {
    // TODO: Implement part 1
    0
}

fn part2(input: &str) -> i64 {
    // TODO: Implement part 2
    0
}

fn main() {
    let input = fs::read_to_string("inputs/day01.txt")
        .expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
TODO: Add example input
";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 0);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 0);
    }
}
