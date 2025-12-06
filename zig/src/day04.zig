const std = @import("std");

fn parseInput(allocator: std.mem.Allocator, input: []const u8) ![][]u8 {
    var grid = try std.ArrayList([]u8).initCapacity(allocator, 10);
    errdefer {
        for (grid.items) |row| allocator.free(row);
        grid.deinit(allocator);
    }

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        const row = try allocator.alloc(u8, trimmed.len);
        @memcpy(row, trimmed);
        try grid.append(allocator, row);
    }

    return grid.toOwnedSlice(allocator);
}

fn countAdjacentPapers(grid: []const []const u8, row: usize, col: usize) usize {
    const rows = @as(i32, @intCast(grid.len));
    const cols = @as(i32, @intCast(grid[0].len));
    const directions = [_][2]i32{
        .{ -1, -1 }, .{ -1, 0 }, .{ -1, 1 },
        .{ 0, -1 },              .{ 0, 1 },
        .{ 1, -1 },  .{ 1, 0 },  .{ 1, 1 },
    };

    var count: usize = 0;
    for (directions) |dir| {
        const new_row = @as(i32, @intCast(row)) + dir[0];
        const new_col = @as(i32, @intCast(col)) + dir[1];

        if (new_row >= 0 and new_row < rows and new_col >= 0 and new_col < cols) {
            const r = @as(usize, @intCast(new_row));
            const c = @as(usize, @intCast(new_col));
            if (grid[r][c] == '@') {
                count += 1;
            }
        }
    }
    return count;
}

fn part1(input: []const u8) !u32 {
    const allocator = std.heap.page_allocator;
    const grid = try parseInput(allocator, input);
    defer {
        for (grid) |row| allocator.free(row);
        allocator.free(grid);
    }

    var accessible: u32 = 0;

    for (grid, 0..) |row, r| {
        for (row, 0..) |cell, c| {
            if (cell == '@') {
                const adjacent_count = countAdjacentPapers(grid, r, c);
                if (adjacent_count < 4) {
                    accessible += 1;
                }
            }
        }
    }

    return accessible;
}

fn part2(input: []const u8) !u32 {
    const allocator = std.heap.page_allocator;
    var grid = try parseInput(allocator, input);
    defer {
        for (grid) |row| allocator.free(row);
        allocator.free(grid);
    }

    var total_removed: u32 = 0;

    while (true) {
        // Find all accessible rolls in this iteration
        var to_remove = try std.ArrayList([2]usize).initCapacity(allocator, 10);
        defer to_remove.deinit(allocator);

        for (grid, 0..) |row, r| {
            for (row, 0..) |cell, c| {
                if (cell == '@') {
                    const adjacent_count = countAdjacentPapers(grid, r, c);
                    if (adjacent_count < 4) {
                        try to_remove.append(allocator, .{ r, c });
                    }
                }
            }
        }

        // If no rolls can be removed, we're done
        if (to_remove.items.len == 0) {
            break;
        }

        // Remove all accessible rolls
        for (to_remove.items) |pos| {
            grid[pos[0]][pos[1]] = '.';
        }

        total_removed += @intCast(to_remove.items.len);
    }

    return total_removed;
}

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    const file = try std.fs.cwd().openFile("../inputs/day04.txt", .{});
    defer file.close();

    const input = try file.readToEndAlloc(allocator, 1024 * 1024);
    defer allocator.free(input);

    const result1 = try part1(input);
    const result2 = try part2(input);

    std.debug.print("Part 1: {}\n", .{result1});
    std.debug.print("Part 2: {}\n", .{result2});
}

test "part1" {
    const example =
        \\..@@.@@@@.
        \\@@@.@.@.@@
        \\@@@@@.@.@@
        \\@.@@@@..@.
        \\@@.@@@@.@@
        \\.@@@@@@@.@
        \\.@.@.@.@@@
        \\@.@@@.@@@@
        \\.@@@@@@@@.
        \\@.@.@@@.@.
    ;

    const result = try part1(example);
    try std.testing.expectEqual(@as(u32, 13), result);
}

test "part2" {
    const example =
        \\..@@.@@@@.
        \\@@@.@.@.@@
        \\@@@@@.@.@@
        \\@.@@@@..@.
        \\@@.@@@@.@@
        \\.@@@@@@@.@
        \\.@.@.@.@@@
        \\@.@@@.@@@@
        \\.@@@@@@@@.
        \\@.@.@@@.@.
    ;

    const result = try part2(example);
    try std.testing.expectEqual(@as(u32, 43), result);
}
