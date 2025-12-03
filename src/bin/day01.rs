use nom::{
    IResult,
    branch::alt,
    character::complete::{char, i32 as nom_i32, line_ending},
    combinator::{map, value},
    multi::separated_list1,
    sequence::tuple,
};
use std::fs;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Direction {
    Left,
    Right,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Rotation {
    direction: Direction,
    distance: i32,
}

fn parse_direction(input: &str) -> IResult<&str, Direction> {
    alt((
        value(Direction::Left, char('L')),
        value(Direction::Right, char('R')),
    ))(input)
}

fn parse_rotation(input: &str) -> IResult<&str, Rotation> {
    map(
        tuple((parse_direction, nom_i32)),
        |(direction, distance)| Rotation {
            direction,
            distance,
        },
    )(input)
}

fn parse_rotations(input: &str) -> Vec<Rotation> {
    let result: IResult<&str, Vec<Rotation>> = separated_list1(line_ending, parse_rotation)(input);
    match result {
        Ok((_, rotations)) => rotations,
        Err(e) => panic!("Failed to parse input: {}", e),
    }
}

fn part1(input: &str) -> i64 {
    let rotations = parse_rotations(input);
    let mut position = 50;
    let mut count = 0;

    for rotation in rotations {
        match rotation.direction {
            Direction::Left => {
                position = (position - rotation.distance).rem_euclid(100);
            }
            Direction::Right => {
                position = (position + rotation.distance).rem_euclid(100);
            }
        }

        if position == 0 {
            count += 1;
        }
    }

    count
}

fn part2(input: &str) -> i64 {
    let rotations = parse_rotations(input);
    let mut position = 50;
    let mut count: i64 = 0;

    for rotation in rotations {
        // Count how many times we pass through 0 during this rotation
        let clicks_through_zero = match rotation.direction {
            Direction::Right => {
                // Going right from position P by distance D
                // We cross 0 at positions where (P + k) % 100 == 0 for k in [1, D]
                // This happens floor((P + D) / 100) - floor(P / 100) times
                (((position + rotation.distance) / 100) - (position / 100)) as i64
            }
            Direction::Left => {
                // Going left from position P by distance D
                // We cross 0 at positions where (P - k) % 100 == 0 for k in [1, D]
                // This happens when k = P, P+100, P+200, ... (all ≤ D)
                if position == 0 {
                    // Starting at 0, we cross it again at k = 100, 200, ...
                    (rotation.distance / 100) as i64
                } else if rotation.distance >= position {
                    // We cross 0 at k = P, then every 100 clicks after
                    ((rotation.distance - position) / 100 + 1) as i64
                } else {
                    // We don't reach 0
                    0
                }
            }
        };

        count += clicks_through_zero;

        // Update position
        match rotation.direction {
            Direction::Left => {
                position = (position - rotation.distance).rem_euclid(100);
            }
            Direction::Right => {
                position = (position + rotation.distance).rem_euclid(100);
            }
        }
    }

    count
}

fn main() {
    let input = fs::read_to_string("inputs/day01.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 3);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 6);
    }
}
