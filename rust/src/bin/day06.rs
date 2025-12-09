use std::fs;

fn parse_and_solve(input: &str) -> i64 {
    let lines: Vec<&str> = input.lines().collect();
    if lines.is_empty() {
        return 0;
    }

    // Find the maximum line length
    let max_len = lines.iter().map(|l| l.len()).max().unwrap_or(0);

    // Parse character by character into columns
    let mut columns: Vec<Vec<char>> = vec![Vec::new(); max_len];
    for line in &lines {
        let chars: Vec<char> = line.chars().collect();
        for (col_idx, &ch) in chars.iter().enumerate() {
            columns[col_idx].push(ch);
        }
        // Pad shorter lines with spaces
        for col_idx in chars.len()..max_len {
            columns[col_idx].push(' ');
        }
    }

    // Group consecutive non-space columns into problems
    let mut problems: Vec<Vec<Vec<char>>> = Vec::new();
    let mut current_problem: Vec<Vec<char>> = Vec::new();

    for column in columns {
        // Check if this is a separator column (all spaces)
        if column.iter().all(|&ch| ch == ' ') {
            if !current_problem.is_empty() {
                problems.push(current_problem.clone());
                current_problem.clear();
            }
        } else {
            current_problem.push(column);
        }
    }

    // Don't forget the last problem
    if !current_problem.is_empty() {
        problems.push(current_problem);
    }

    // Solve each problem
    let mut grand_total = 0i64;
    for problem in problems {
        if let Some(result) = solve_problem(&problem) {
            grand_total += result;
        }
    }

    grand_total
}

fn solve_problem(problem: &[Vec<char>]) -> Option<i64> {
    if problem.is_empty() {
        return None;
    }

    let num_rows = problem[0].len();

    // Extract the operation from the last row
    let operation_str: String = problem
        .iter()
        .map(|col| col.last().copied().unwrap_or(' '))
        .collect();
    let operation = operation_str.trim();

    if operation != "*" && operation != "+" {
        return None;
    }

    // Extract numbers from rows above the operation
    let mut numbers = Vec::new();
    for row_idx in 0..num_rows - 1 {
        let row_str: String = problem
            .iter()
            .map(|col| col.get(row_idx).copied().unwrap_or(' '))
            .collect();
        let trimmed = row_str.trim();
        if !trimmed.is_empty() {
            if let Ok(num) = trimmed.parse::<i64>() {
                numbers.push(num);
            }
        }
    }

    if numbers.is_empty() {
        return None;
    }

    // Calculate result
    let result = if operation == "*" {
        numbers.iter().product()
    } else {
        numbers.iter().sum()
    };

    Some(result)
}

fn solve_problem_part2(problem: &[Vec<char>]) -> Option<i64> {
    if problem.is_empty() {
        return None;
    }

    let num_rows = problem[0].len();

    // Find the operator (should be in the last row)
    let operator = problem
        .iter()
        .map(|col| col.last().copied().unwrap_or(' '))
        .find(|&ch| ch == '*' || ch == '+')?;

    // Extract numbers by reading columns right-to-left
    // Each column represents one number (digits stacked vertically)
    let mut numbers = Vec::new();

    for col in problem.iter().rev() {
        // Build number from this column (top to bottom, excluding operator row)
        let mut digit_str = String::new();
        for row_idx in 0..num_rows - 1 {
            let ch = col.get(row_idx).copied().unwrap_or(' ');
            if ch != ' ' {
                digit_str.push(ch);
            }
        }

        if !digit_str.is_empty() {
            if let Ok(num) = digit_str.parse::<i64>() {
                numbers.push(num);
            }
        }
    }

    if numbers.is_empty() {
        return None;
    }

    // Calculate result
    let result = if operator == '*' {
        numbers.iter().product()
    } else {
        numbers.iter().sum()
    };

    Some(result)
}

fn part1(input: &str) -> i64 {
    parse_and_solve(input)
}

fn part2(input: &str) -> i64 {
    let lines: Vec<&str> = input.lines().collect();
    if lines.is_empty() {
        return 0;
    }

    // Find the maximum line length
    let max_len = lines.iter().map(|l| l.len()).max().unwrap_or(0);

    // Parse character by character into columns
    let mut columns: Vec<Vec<char>> = vec![Vec::new(); max_len];
    for line in &lines {
        let chars: Vec<char> = line.chars().collect();
        for (col_idx, &ch) in chars.iter().enumerate() {
            columns[col_idx].push(ch);
        }
        // Pad shorter lines with spaces
        for col_idx in chars.len()..max_len {
            columns[col_idx].push(' ');
        }
    }

    // Group consecutive non-space columns into problems
    let mut problems: Vec<Vec<Vec<char>>> = Vec::new();
    let mut current_problem: Vec<Vec<char>> = Vec::new();

    for column in columns {
        // Check if this is a separator column (all spaces)
        if column.iter().all(|&ch| ch == ' ') {
            if !current_problem.is_empty() {
                problems.push(current_problem.clone());
                current_problem.clear();
            }
        } else {
            current_problem.push(column);
        }
    }

    // Don't forget the last problem
    if !current_problem.is_empty() {
        problems.push(current_problem);
    }

    // Solve each problem using part2 logic
    let mut grand_total = 0i64;
    for problem in problems {
        if let Some(result) = solve_problem_part2(&problem) {
            grand_total += result;
        }
    }

    grand_total
}

fn main() {
    let input = fs::read_to_string("../inputs/day06.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
123 328  51 64
 45 64  387 23
  6 98  215 314
*   +   *   +  ";

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 4277556);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 3263827);
    }
}
