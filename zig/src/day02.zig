const std = @import("std");

const Range = struct {
    start: u64,
    end: u64,
};

fn parseRanges(allocator: std.mem.Allocator, input: []const u8) ![]Range {
    var input_clean = try std.ArrayList(u8).initCapacity(allocator, input.len);
    defer input_clean.deinit(allocator);

    // Remove newlines and trim
    for (input) |c| {
        if (c != '\n' and c != '\r') {
            try input_clean.append(allocator, c);
        }
    }

    const cleaned = std.mem.trim(u8, input_clean.items, " \t");

    var ranges = try std.ArrayList(Range).initCapacity(allocator, 10);
    errdefer ranges.deinit(allocator);

    var parts = std.mem.splitScalar(u8, cleaned, ',');
    while (parts.next()) |part| {
        const trimmed = std.mem.trim(u8, part, " \t");
        if (trimmed.len == 0) continue;

        var range_parts = std.mem.splitScalar(u8, trimmed, '-');
        const start_str = range_parts.next() orelse continue;
        const end_str = range_parts.next() orelse continue;

        const start = try std.fmt.parseUnsigned(u64, start_str, 10);
        const end = try std.fmt.parseUnsigned(u64, end_str, 10);

        try ranges.append(allocator, Range{ .start = start, .end = end });
    }

    return ranges.toOwnedSlice(allocator);
}

fn isRepeatedPattern(n: u64) bool {
    var buf: [32]u8 = undefined;
    const s = std.fmt.bufPrint(&buf, "{}", .{n}) catch return false;
    const length = s.len;

    // Must have even length
    if (length % 2 != 0) {
        return false;
    }

    // Split in half
    const mid = length / 2;
    const first_half = s[0..mid];
    const second_half = s[mid..];

    // Check no leading zeros
    if (first_half[0] == '0') {
        return false;
    }

    // Check if halves are equal
    return std.mem.eql(u8, first_half, second_half);
}

fn part1(input: []const u8) !u64 {
    const allocator = std.heap.page_allocator;
    const ranges = try parseRanges(allocator, input);
    defer allocator.free(ranges);

    var sum: u64 = 0;

    for (ranges) |r| {
        var id = r.start;
        while (id <= r.end) : (id += 1) {
            if (isRepeatedPattern(id)) {
                sum += id;
            }
        }
    }

    return sum;
}

fn isRepeatedPatternV2(n: u64) bool {
    var buf: [32]u8 = undefined;
    const s = std.fmt.bufPrint(&buf, "{}", .{n}) catch return false;
    const length = s.len;

    // Try all possible pattern lengths from 1 to len/2
    var pattern_len: usize = 1;
    while (pattern_len <= length / 2) : (pattern_len += 1) {
        // Check if len is divisible by pattern_len (pattern must repeat evenly)
        if (length % pattern_len == 0) {
            const pattern = s[0..pattern_len];

            // Check for leading zeros
            if (pattern[0] == '0') {
                continue;
            }

            // Check if all chunks match the pattern
            var all_match = true;
            var i: usize = 0;
            while (i < length) : (i += pattern_len) {
                if (!std.mem.eql(u8, s[i .. i + pattern_len], pattern)) {
                    all_match = false;
                    break;
                }
            }

            if (all_match) {
                return true;
            }
        }
    }

    return false;
}

fn part2(input: []const u8) !u64 {
    const allocator = std.heap.page_allocator;
    const ranges = try parseRanges(allocator, input);
    defer allocator.free(ranges);

    var sum: u64 = 0;

    for (ranges) |r| {
        var id = r.start;
        while (id <= r.end) : (id += 1) {
            if (isRepeatedPatternV2(id)) {
                sum += id;
            }
        }
    }

    return sum;
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
