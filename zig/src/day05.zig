const std = @import("std");

const Range = struct {
    start: u64,
    end: u64,

    fn contains(self: Range, id: u64) bool {
        return id >= self.start and id <= self.end;
    }
};

fn parseInput(allocator: std.mem.Allocator, input: []const u8) !struct { ranges: []Range, ids: []u64 } {
    var parts = std.mem.splitSequence(u8, std.mem.trim(u8, input, " \t\r\n"), "\n\n");

    const ranges_part = parts.next() orelse return error.InvalidInput;
    const ids_part = parts.next() orelse return error.InvalidInput;

    // Parse ranges
    var ranges = try std.ArrayList(Range).initCapacity(allocator, 10);
    errdefer ranges.deinit(allocator);

    var range_lines = std.mem.splitScalar(u8, ranges_part, '\n');
    while (range_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        var parts_iter = std.mem.splitScalar(u8, trimmed, '-');
        const start_str = parts_iter.next() orelse return error.InvalidRange;
        const end_str = parts_iter.next() orelse return error.InvalidRange;

        const start = try std.fmt.parseUnsigned(u64, start_str, 10);
        const end = try std.fmt.parseUnsigned(u64, end_str, 10);

        try ranges.append(allocator, .{ .start = start, .end = end });
    }

    // Parse IDs
    var ids = try std.ArrayList(u64).initCapacity(allocator, 10);
    errdefer ids.deinit(allocator);

    var id_lines = std.mem.splitScalar(u8, ids_part, '\n');
    while (id_lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        const id = try std.fmt.parseUnsigned(u64, trimmed, 10);
        try ids.append(allocator, id);
    }

    return .{
        .ranges = try ranges.toOwnedSlice(allocator),
        .ids = try ids.toOwnedSlice(allocator),
    };
}

fn isFresh(id: u64, ranges: []const Range) bool {
    for (ranges) |r| {
        if (r.contains(id)) {
            return true;
        }
    }
    return false;
}

fn part1(input: []const u8) !usize {
    const allocator = std.heap.page_allocator;
    const parsed = try parseInput(allocator, input);
    defer allocator.free(parsed.ranges);
    defer allocator.free(parsed.ids);

    var count: usize = 0;
    for (parsed.ids) |id| {
        if (isFresh(id, parsed.ranges)) {
            count += 1;
        }
    }

    return count;
}

fn lessThan(_: void, a: Range, b: Range) bool {
    return a.start < b.start;
}

fn mergeRanges(allocator: std.mem.Allocator, ranges: []const Range) ![]Range {
    if (ranges.len == 0) {
        return try allocator.alloc(Range, 0);
    }

    // Sort ranges by start position
    const sorted = try allocator.alloc(Range, ranges.len);
    @memcpy(sorted, ranges);
    std.mem.sort(Range, sorted, {}, lessThan);

    var merged = try std.ArrayList(Range).initCapacity(allocator, 10);
    errdefer merged.deinit(allocator);

    try merged.append(allocator, sorted[0]);

    for (sorted[1..]) |r| {
        var last = &merged.items[merged.items.len - 1];

        // Check if ranges overlap or are adjacent
        if (r.start <= last.end + 1) {
            // Merge by extending the end
            if (r.end > last.end) {
                last.end = r.end;
            }
        } else {
            // No overlap, add as new range
            try merged.append(allocator, r);
        }
    }

    allocator.free(sorted);
    return merged.toOwnedSlice(allocator);
}

fn countIDsInRanges(allocator: std.mem.Allocator, ranges: []const Range) !u64 {
    const merged = try mergeRanges(allocator, ranges);
    defer allocator.free(merged);

    var total: u64 = 0;
    for (merged) |r| {
        total += r.end - r.start + 1;
    }
    return total;
}

fn part2(input: []const u8) !u64 {
    const allocator = std.heap.page_allocator;
    const parsed = try parseInput(allocator, input);
    defer allocator.free(parsed.ranges);
    defer allocator.free(parsed.ids);

    return try countIDsInRanges(allocator, parsed.ranges);
}

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    const file = try std.fs.cwd().openFile("../inputs/day05.txt", .{});
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
        \\3-5
        \\10-14
        \\16-20
        \\12-18
        \\
        \\1
        \\5
        \\8
        \\11
        \\17
        \\32
    ;

    const result = try part1(example);
    try std.testing.expectEqual(@as(usize, 3), result);
}

test "part2" {
    const example =
        \\3-5
        \\10-14
        \\16-20
        \\12-18
        \\
        \\1
        \\5
        \\8
        \\11
        \\17
        \\32
    ;

    const result = try part2(example);
    try std.testing.expectEqual(@as(u64, 14), result);
}
