const std = @import("std");

fn part1(input: []const u8) !u32 {
    _ = input;
    // TODO: Implement part 1
    return 0;
}

fn part2(input: []const u8) !u64 {
    _ = input;
    // TODO: Implement part 2
    return 0;
}

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    const file = try std.fs.cwd().openFile("../inputs/day03.txt", .{});
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
        \\987654321111111
        \\811111111111119
        \\234234234234278
        \\818181911112111
    ;

    const result = try part1(example);
    try std.testing.expectEqual(@as(u32, 357), result);
}

test "part2" {
    const example =
        \\987654321111111
        \\811111111111119
        \\234234234234278
        \\818181911112111
    ;

    const result = try part2(example);
    try std.testing.expectEqual(@as(u64, 3121910778619), result);
}
