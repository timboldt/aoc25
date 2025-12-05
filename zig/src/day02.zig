const std = @import("std");

fn part1(input: []const u8) !u64 {
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

    const file = try std.fs.cwd().openFile("../inputs/day02.txt", .{});
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
        \\11-22,95-115,998-1012,1188511880-1188511890,222220-222224,
        \\1698522-1698528,446443-446449,38593856-38593862,565653-565659,
        \\824824821-824824827,2121212118-2121212124
    ;

    const result = try part1(example);
    try std.testing.expectEqual(@as(u64, 1227775554), result);
}

test "part2" {
    const example =
        \\11-22,95-115,998-1012,1188511880-1188511890,222220-222224,
        \\1698522-1698528,446443-446449,38593856-38593862,565653-565659,
        \\824824821-824824827,2121212118-2121212124
    ;

    const result = try part2(example);
    try std.testing.expectEqual(@as(u64, 4174379265), result);
}
