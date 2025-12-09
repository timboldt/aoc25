const std = @import("std");

fn parseAndSolve(allocator: std.mem.Allocator, input: []const u8) !i64 {
    var lines = try std.ArrayList([]const u8).initCapacity(allocator, 10);
    defer lines.deinit(allocator);

    var line_iter = std.mem.splitScalar(u8, input, '\n');
    while (line_iter.next()) |line| {
        try lines.append(allocator, line);
    }

    if (lines.items.len == 0) {
        return 0;
    }

    // Find max line length
    var max_len: usize = 0;
    for (lines.items) |line| {
        if (line.len > max_len) {
            max_len = line.len;
        }
    }

    // Parse character by character into columns
    var columns = try allocator.alloc(std.ArrayList(u8), max_len);
    defer {
        for (columns) |*col| {
            col.deinit(allocator);
        }
        allocator.free(columns);
    }

    for (columns) |*col| {
        col.* = try std.ArrayList(u8).initCapacity(allocator, 10);
    }

    for (lines.items) |line| {
        for (line, 0..) |ch, col_idx| {
            try columns[col_idx].append(allocator, ch);
        }
        // Pad shorter lines with spaces
        for (line.len..max_len) |col_idx| {
            try columns[col_idx].append(allocator, ' ');
        }
    }

    // Group consecutive non-space columns into problems
    var problems = try std.ArrayList([]usize).initCapacity(allocator, 10);
    defer {
        for (problems.items) |problem| {
            allocator.free(problem);
        }
        problems.deinit(allocator);
    }

    var current_problem = try std.ArrayList(usize).initCapacity(allocator, 10);
    defer current_problem.deinit(allocator);

    for (columns, 0..) |column, col_idx| {
        // Check if this is a separator column (all spaces)
        var all_spaces = true;
        for (column.items) |ch| {
            if (ch != ' ') {
                all_spaces = false;
                break;
            }
        }

        if (all_spaces) {
            if (current_problem.items.len > 0) {
                try problems.append(allocator, try current_problem.toOwnedSlice(allocator));
                current_problem.clearRetainingCapacity();
            }
        } else {
            try current_problem.append(allocator, col_idx);
        }
    }

    // Don't forget the last problem
    if (current_problem.items.len > 0) {
        try problems.append(allocator, try current_problem.toOwnedSlice(allocator));
    }

    // Solve each problem
    var grand_total: i64 = 0;
    for (problems.items) |problem_cols| {
        if (try solveProblem(allocator, columns, problem_cols)) |result| {
            grand_total += result;
        }
    }

    return grand_total;
}

fn solveProblem(allocator: std.mem.Allocator, columns: []std.ArrayList(u8), problem_cols: []usize) !?i64 {
    if (problem_cols.len == 0) {
        return null;
    }

    const num_rows = columns[problem_cols[0]].items.len;

    // Extract the operation from the last row
    var operation_str = try std.ArrayList(u8).initCapacity(allocator, 10);
    defer operation_str.deinit(allocator);

    for (problem_cols) |col_idx| {
        if (col_idx < columns.len and num_rows > 0) {
            try operation_str.append(allocator, columns[col_idx].items[num_rows - 1]);
        }
    }

    const operation = std.mem.trim(u8, operation_str.items, " ");

    if (!std.mem.eql(u8, operation, "*") and !std.mem.eql(u8, operation, "+")) {
        return null;
    }

    // Extract numbers from rows above the operation
    var numbers = try std.ArrayList(i64).initCapacity(allocator, 10);
    defer numbers.deinit(allocator);

    for (0..num_rows - 1) |row_idx| {
        var row_str = try std.ArrayList(u8).initCapacity(allocator, 10);
        defer row_str.deinit(allocator);

        for (problem_cols) |col_idx| {
            if (col_idx < columns.len and row_idx < columns[col_idx].items.len) {
                try row_str.append(allocator, columns[col_idx].items[row_idx]);
            }
        }

        const trimmed = std.mem.trim(u8, row_str.items, " ");
        if (trimmed.len > 0) {
            if (std.fmt.parseInt(i64, trimmed, 10)) |num| {
                try numbers.append(allocator, num);
            } else |_| {}
        }
    }

    if (numbers.items.len == 0) {
        return null;
    }

    // Calculate result
    var result: i64 = if (std.mem.eql(u8, operation, "*")) 1 else 0;

    if (std.mem.eql(u8, operation, "*")) {
        for (numbers.items) |num| {
            result *= num;
        }
    } else {
        for (numbers.items) |num| {
            result += num;
        }
    }

    return result;
}

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    const file = try std.fs.cwd().openFile("../inputs/day06.txt", .{});
    defer file.close();

    const input = try file.readToEndAlloc(allocator, 1024 * 1024);
    defer allocator.free(input);

    const result = try parseAndSolve(allocator, input);
    std.debug.print("Part 1: {d}\n", .{result});
}

test "example" {
    const input =
        \\123 328  51 64
        \\ 45 64  387 23
        \\  6 98  215 314
        \\*   +   *   +
    ;

    const allocator = std.testing.allocator;
    const result = try parseAndSolve(allocator, input);
    try std.testing.expectEqual(@as(i64, 4277556), result);
}
