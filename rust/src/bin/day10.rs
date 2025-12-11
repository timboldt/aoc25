use std::fs;
use std::ops::{Add, Sub, Mul, Div, Neg, AddAssign, SubAssign};
use std::cmp::Ordering;

#[derive(Debug, Clone)]
struct Machine {
    target: Vec<bool>,
    buttons: Vec<Vec<bool>>,
    // For part 2: buttons as index lists
    button_indices: Vec<Vec<usize>>,
    // For part 2: target counter values
    targets: Vec<i64>,
}

#[derive(Clone, Copy, Debug, Eq)]
struct Rat {
    num: i128,
    den: i128,
}

impl Rat {
    fn new(num: i128, den: i128) -> Self {
        if den == 0 {
            panic!("Division by zero");
        }
        let g = gcd(num.abs(), den.abs());
        let mut r = Rat {
            num: num / g,
            den: den / g,
        };
        if r.den < 0 {
            r.num = -r.num;
            r.den = -r.den;
        }
        r
    }

    fn zero() -> Self {
        Rat { num: 0, den: 1 }
    }

    fn one() -> Self {
        Rat { num: 1, den: 1 }
    }
    
    fn from_i64(v: i64) -> Self {
        Rat::new(v as i128, 1)
    }

    fn is_integer(&self) -> bool {
        self.den == 1
    }

    fn to_i64(&self) -> Option<i64> {
        if self.is_integer() && self.num >= i64::MIN as i128 && self.num <= i64::MAX as i128 {
            Some(self.num as i64)
        } else {
            None
        }
    }
}

impl PartialEq for Rat {
    fn eq(&self, other: &Self) -> bool {
        self.num == other.num && self.den == other.den
    }
}

impl PartialOrd for Rat {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for Rat {
    fn cmp(&self, other: &Self) -> Ordering {
        (self.num * other.den).cmp(&(other.num * self.den))
    }
}

impl Add for Rat {
    type Output = Self;
    fn add(self, rhs: Self) -> Self {
        Rat::new(self.num * rhs.den + rhs.num * self.den, self.den * rhs.den)
    }
}

impl Sub for Rat {
    type Output = Self;
    fn sub(self, rhs: Self) -> Self {
        Rat::new(self.num * rhs.den - rhs.num * self.den, self.den * rhs.den)
    }
}

impl Mul for Rat {
    type Output = Self;
    fn mul(self, rhs: Self) -> Self {
        Rat::new(self.num * rhs.num, self.den * rhs.den)
    }
}

impl Div for Rat {
    type Output = Self;
    fn div(self, rhs: Self) -> Self {
        Rat::new(self.num * rhs.den, self.den * rhs.num)
    }
}

impl Neg for Rat {
    type Output = Self;
    fn neg(self) -> Self {
        Rat { num: -self.num, den: self.den }
    }
}

impl AddAssign for Rat {
    fn add_assign(&mut self, rhs: Self) {
        *self = *self + rhs;
    }
}

impl SubAssign for Rat {
    fn sub_assign(&mut self, rhs: Self) {
        *self = *self - rhs;
    }
}

fn gcd(mut a: i128, mut b: i128) -> i128 {
    while b != 0 {
        let t = b;
        b = a % b;
        a = t;
    }
    a
}

fn parse_machine(line: &str) -> Machine {
    let parts: Vec<&str> = line.split(|c| c == '[' || c == ']' || c == '{' || c == '}')
        .filter(|s| !s.trim().is_empty())
        .collect();

    // Parse target state (indicator lights)
    let target_str = parts[0].trim();
    let target: Vec<bool> = target_str.chars().map(|c| c == '#').collect();

    // Parse buttons (everything between ] and {)
    let buttons_str = parts[1].trim();
    let buttons: Vec<Vec<bool>> = buttons_str
        .split(')')
        .filter_map(|s| {
            let s = s.trim().trim_start_matches('(').trim();
            if s.is_empty() {
                return None;
            }

            let num_lights = target.len();
            let mut button_mask = vec![false; num_lights];

            for num_str in s.split(',') {
                if let Ok(idx) = num_str.trim().parse::<usize>() {
                    if idx < num_lights {
                        button_mask[idx] = true;
                    }
                }
            }

            Some(button_mask)
        })
        .collect();

    // Parse button_indices for part 2
    let mut button_indices = Vec::new();
    let mut cursor = 0;
    while let Some(start) = line[cursor..].find('(') {
        let open = cursor + start;
        if let Some(end) = line[open..].find(')') {
            let close = open + end;
            let content = &line[open+1..close];
            let indices: Vec<usize> = content.split(',')
                .filter_map(|s| s.trim().parse().ok())
                .collect();
            button_indices.push(indices);
            cursor = close + 1;
        } else {
            break;
        }
    }

    // Parse joltage requirements / targets
    let targets: Vec<i64> = if parts.len() > 2 {
        parts[2]
            .trim()
            .split(',')
            .filter_map(|s| s.trim().parse().ok())
            .collect()
    } else {
        vec![]
    };

    Machine { target, buttons, button_indices, targets }
}

fn solve_machine(machine: &Machine) -> usize {
    let n_lights = machine.target.len();
    let n_buttons = machine.buttons.len();

    if n_buttons == 0 {
        return if machine.target.iter().all(|&b| !b) { 0 } else { usize::MAX };
    }

    // Create augmented matrix for Gaussian elimination over GF(2)
    let mut matrix = vec![vec![false; n_buttons + 1]; n_lights];

    for light in 0..n_lights {
        for button in 0..n_buttons {
            matrix[light][button] = machine.buttons[button][light];
        }
        matrix[light][n_buttons] = machine.target[light];
    }

    // Gaussian elimination
    let mut pivot_col = Vec::new();
    let mut row = 0;

    for col in 0..n_buttons {
        let mut pivot_row = None;
        for r in row..n_lights {
            if matrix[r][col] {
                pivot_row = Some(r);
                break;
            }
        }

        let Some(pivot) = pivot_row else {
            continue;
        };

        if pivot != row {
            matrix.swap(row, pivot);
        }

        pivot_col.push(col);

        for r in 0..n_lights {
            if r != row && matrix[r][col] {
                for c in 0..=n_buttons {
                    matrix[r][c] ^= matrix[row][c];
                }
            }
        }

        row += 1;
    }

    for r in row..n_lights {
        if matrix[r][n_buttons] {
            return usize::MAX;
        }
    }

    let mut is_pivot = vec![false; n_buttons];
    for &col in &pivot_col {
        is_pivot[col] = true;
    }

    let free_vars: Vec<usize> = (0..n_buttons)
        .filter(|&i| !is_pivot[i])
        .collect();

    let n_free = free_vars.len();
    let mut min_presses = usize::MAX;

    for mask in 0..(1 << n_free) {
        let mut solution = vec![false; n_buttons];

        for (i, &var) in free_vars.iter().enumerate() {
            solution[var] = (mask & (1 << i)) != 0;
        }

        for (row_idx, &col) in pivot_col.iter().enumerate().rev() {
            let mut val = matrix[row_idx][n_buttons];
            for button in 0..n_buttons {
                if button != col && matrix[row_idx][button] {
                    val ^= solution[button];
                }
            }
            solution[col] = val;
        }

        let presses = solution.iter().filter(|&&b| b).count();
        min_presses = min_presses.min(presses);
    }

    min_presses
}

fn solve_joltage(machine: &Machine) -> i64 {
    let num_counters = machine.targets.len();
    let num_buttons = machine.button_indices.len();

    if num_buttons == 0 {
        return if machine.targets.iter().all(|&v| v == 0) { 0 } else { i64::MAX };
    }
    if machine.targets.iter().all(|&v| v == 0) {
        return 0;
    }

    // Build the system Ax = b
    // A: num_counters x num_buttons
    let mut matrix: Vec<Vec<Rat>> = vec![vec![Rat::zero(); num_buttons]; num_counters];
    let mut target_vec: Vec<Rat> = Vec::with_capacity(num_counters);

    for (r, &t) in machine.targets.iter().enumerate() {
        target_vec.push(Rat::from_i64(t));
        for (c, indices) in machine.button_indices.iter().enumerate() {
            if indices.contains(&r) {
                matrix[r][c] = Rat::one();
            }
        }
    }

    // Gaussian Elimination with Rational Arithmetic
    let mut pivot_cols = Vec::new();
    let mut pivot_rows = Vec::new(); // map from col index to row index in RREF
    let mut row = 0;

    for col in 0..num_buttons {
        let mut pivot_row = None;
        for r in row..num_counters {
            if matrix[r][col] != Rat::zero() {
                pivot_row = Some(r);
                break;
            }
        }

        let Some(pivot) = pivot_row else {
            continue;
        };

        if pivot != row {
            matrix.swap(row, pivot);
            target_vec.swap(row, pivot);
        }

        pivot_cols.push(col);
        pivot_rows.push(row);

        // Normalize pivot
        let pivot_val = matrix[row][col];
        for c in col..num_buttons {
            matrix[row][c] = matrix[row][c] / pivot_val;
        }
        target_vec[row] = target_vec[row] / pivot_val;

        // Eliminate
        let pivot_row_vals: Vec<Rat> = matrix[row].clone();
        let pivot_target = target_vec[row];

        for r in 0..num_counters {
            if r != row && matrix[r][col] != Rat::zero() {
                let factor = matrix[r][col];
                for c in col..num_buttons {
                    matrix[r][c] -= factor * pivot_row_vals[c];
                }
                target_vec[r] -= factor * pivot_target;
            }
        }

        row += 1;
    }

    // Check inconsistency
    for r in row..num_counters {
        if target_vec[r] != Rat::zero() {
            return i64::MAX;
        }
    }

    // Identify free variables
    let mut is_pivot = vec![false; num_buttons];
    for &col in &pivot_cols {
        is_pivot[col] = true;
    }
    let free_vars: Vec<usize> = (0..num_buttons).filter(|&i| !is_pivot[i]).collect();

    // Minimize sum
    let mut min_total = i64::MAX;
    
    // Bounds for free variables
    let max_target = *machine.targets.iter().max().unwrap_or(&0);
    let limit = max_target + 2; 

    // Search free variables
    let mut current_free_vals = vec![0i64; free_vars.len()];
    
    // Optimization: if no free vars, just check the solution
    if free_vars.is_empty() {
        let mut valid = true;
        let mut sum = 0;
        for i in 0..num_buttons {
             if let Some(pos) = pivot_cols.iter().position(|&c| c == i) {
                 let r = pivot_rows[pos];
                 let val = target_vec[r];
                 if !val.is_integer() || val < Rat::zero() {
                     valid = false;
                     break;
                 }
                 sum += val.to_i64().unwrap();
             }
        }
        if valid {
            return sum;
        } else {
            return i64::MAX;
        }
    }

    solve_recursive(
        0,
        &free_vars,
        &pivot_cols,
        &pivot_rows,
        &matrix,
        &target_vec,
        limit,
        &mut current_free_vals,
        &mut min_total,
        num_buttons
    );

    min_total
}

fn solve_recursive(
    idx: usize,
    free_vars: &[usize],
    pivot_cols: &[usize],
    pivot_rows: &[usize],
    matrix: &[Vec<Rat>],
    target_vec: &[Rat],
    limit: i64,
    current_free_vals: &mut [i64],
    min_total: &mut i64,
    num_buttons: usize,
) {
    if idx == free_vars.len() {
        // Calculate pivot variables
        let mut current_solution = vec![0i64; num_buttons];
        let mut sum = 0;

        // Set free vars
        for (i, &fv) in free_vars.iter().enumerate() {
            current_solution[fv] = current_free_vals[i];
            sum += current_free_vals[i];
        }

        // Calculate pivots
        // x_pivot = target[row] - sum(matrix[row][free] * x_free)
        let mut valid = true;
        for (i, &pc) in pivot_cols.iter().enumerate() {
            let row = pivot_rows[i];
            let mut val = target_vec[row];
            for (j, &fv) in free_vars.iter().enumerate() {
                let coeff = matrix[row][fv];
                if coeff != Rat::zero() {
                    val -= coeff * Rat::from_i64(current_free_vals[j]);
                }
            }

            if !val.is_integer() || val < Rat::zero() {
                valid = false;
                break;
            }
            let v_int = val.to_i64().unwrap();
            current_solution[pc] = v_int;
            sum += v_int;
        }

        if valid {
            *min_total = (*min_total).min(sum);
        }
        return;
    }

    for val in 0..=limit {
        current_free_vals[idx] = val;
        solve_recursive(idx + 1, free_vars, pivot_cols, pivot_rows, matrix, target_vec, limit, current_free_vals, min_total, num_buttons);
    }
}

fn part1(input: &str) -> usize {
    input
        .lines()
        .map(|line| {
            let machine = parse_machine(line);
            solve_machine(&machine)
        })
        .sum()
}

fn part2(input: &str) -> i64 {
    input
        .lines()
        .map(|line| {
            let machine = parse_machine(line);
            solve_joltage(&machine)
        })
        .filter(|&result| result != i64::MAX)
        .sum()
}

fn main() {
    let input = fs::read_to_string("../inputs/day10.txt").expect("Failed to read input file");

    println!("Part 1: {}", part1(&input));
    println!("Part 2: {}", part2(&input));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}";

    #[test]
    fn test_parse() {
        let machine = parse_machine("[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}");
        assert_eq!(machine.target, vec![false, true, true, false]);
        assert_eq!(machine.buttons.len(), 6);
        assert_eq!(machine.buttons[0], vec![false, false, false, true]);
        assert_eq!(machine.buttons[1], vec![false, true, false, true]);
    }

    #[test]
    fn test_machine1() {
        let machine = parse_machine("[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}");
        assert_eq!(solve_machine(&machine), 2);
    }

    #[test]
    fn test_part1() {
        assert_eq!(part1(EXAMPLE), 7);
    }

    #[test]
    fn test_joltage1() {
        let machine = parse_machine("[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}");
        assert_eq!(solve_joltage(&machine), 10);
    }

    #[test]
    fn test_joltage2() {
        let machine = parse_machine("[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}");
        assert_eq!(solve_joltage(&machine), 12);
    }

    #[test]
    fn test_joltage3() {
        let machine = parse_machine("[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}");
        assert_eq!(solve_joltage(&machine), 11);
    }

    #[test]
    fn test_part2() {
        assert_eq!(part2(EXAMPLE), 33);
    }
}