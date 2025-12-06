const std = @import("std");

fn parseBanks(allocator: std.mem.Allocator, input: []const u8) ![][]const u8 {
    var banks = try std.ArrayList([]const u8).initCapacity(allocator, 10);
    errdefer banks.deinit(allocator);

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;
        try banks.append(allocator, trimmed);
    }

    return banks.toOwnedSlice(allocator);
}

fn maxJoltage(bank: []const u8) u32 {
    var max: u32 = 0;

    // Try all pairs (i, j) where i < j
    var i: usize = 0;
    while (i < bank.len) : (i += 1) {
        var j = i + 1;
        while (j < bank.len) : (j += 1) {
            // Form two-digit number from bank[i] and bank[j]
            const first = bank[i] - '0';
            const second = bank[j] - '0';
            const joltage: u32 = @as(u32, first) * 10 + @as(u32, second);
            if (joltage > max) {
                max = joltage;
            }
        }
    }

    return max;
}

fn part1(input: []const u8) !u32 {
    const allocator = std.heap.page_allocator;
    const banks = try parseBanks(allocator, input);
    defer allocator.free(banks);

    var sum: u32 = 0;

    for (banks) |bank| {
        sum += maxJoltage(bank);
    }

    return sum;
}

fn maxJoltageN(allocator: std.mem.Allocator, bank: []const u8, n: usize) !u64 {
    const to_remove = bank.len - n;

    var stack = try std.ArrayList(u8).initCapacity(allocator, bank.len);
    defer stack.deinit(allocator);

    var removals_left = to_remove;

    // Greedy approach: remove smaller digits when a larger digit comes
    for (bank) |digit| {
        while (stack.items.len > 0 and removals_left > 0 and stack.items[stack.items.len - 1] < digit) {
            _ = stack.pop();
            removals_left -= 1;
        }
        try stack.append(allocator, digit);
    }

    // If we still have removals left, remove from the end
    while (removals_left > 0) : (removals_left -= 1) {
        _ = stack.pop();
    }

    // Take first n digits and form the number
    const result = stack.items[0..n];
    return try std.fmt.parseUnsigned(u64, result, 10);
}

fn part2(input: []const u8) !u64 {
    const allocator = std.heap.page_allocator;
    const banks = try parseBanks(allocator, input);
    defer allocator.free(banks);

    var sum: u64 = 0;

    for (banks) |bank| {
        sum += try maxJoltageN(allocator, bank, 12);
    }

    return sum;
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
