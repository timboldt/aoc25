const std = @import("std");

fn part1(input: []const u8) !i64 {
    _ = input;
    // TODO: Implement part 1
    return 0;
}

fn part2(input: []const u8) !i64 {
    _ = input;
    // TODO: Implement part 2
    return 0;
}

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    const file = try std.fs.cwd().openFile("../inputs/day01.txt", .{});
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
        \\L68
        \\L30
        \\R48
        \\L5
        \\R60
        \\L55
        \\L1
        \\L99
        \\R14
        \\L82
    ;

    const result = try part1(example);
    try std.testing.expectEqual(@as(i64, 3), result);
}

test "part2" {
    const example =
        \\L68
        \\L30
        \\R48
        \\L5
        \\R60
        \\L55
        \\L1
        \\L99
        \\R14
        \\L82
    ;

    const result = try part2(example);
    try std.testing.expectEqual(@as(i64, 6), result);
}
