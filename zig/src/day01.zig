const std = @import("std");

const Direction = enum {
    left,
    right,
};

const Rotation = struct {
    direction: Direction,
    distance: i32,
};

fn parseRotations(allocator: std.mem.Allocator, input: []const u8) ![]Rotation {
    var rotations = try std.ArrayList(Rotation).initCapacity(allocator, 10);
    errdefer rotations.deinit(allocator);

    var lines = std.mem.splitScalar(u8, input, '\n');
    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        const direction: Direction = switch (trimmed[0]) {
            'L' => .left,
            'R' => .right,
            else => return error.InvalidDirection,
        };

        const distance = try std.fmt.parseInt(i32, trimmed[1..], 10);

        try rotations.append(allocator, .{
            .direction = direction,
            .distance = distance,
        });
    }

    return rotations.toOwnedSlice(allocator);
}

fn mod(a: i32, b: i32) i32 {
    return @mod(@mod(a, b) + b, b);
}

fn part1(input: []const u8) !i64 {
    const allocator = std.heap.page_allocator;
    const rotations = try parseRotations(allocator, input);
    defer allocator.free(rotations);

    var position: i32 = 50;
    var count: i64 = 0;

    for (rotations) |rotation| {
        switch (rotation.direction) {
            .left => {
                position = mod(position - rotation.distance, 100);
            },
            .right => {
                position = mod(position + rotation.distance, 100);
            },
        }

        if (position == 0) {
            count += 1;
        }
    }

    return count;
}

fn part2(input: []const u8) !i64 {
    const allocator = std.heap.page_allocator;
    const rotations = try parseRotations(allocator, input);
    defer allocator.free(rotations);

    var position: i32 = 50;
    var count: i64 = 0;

    for (rotations) |rotation| {
        // Count how many times we pass through 0 during this rotation
        const clicks_through_zero: i64 = switch (rotation.direction) {
            .right => blk: {
                // Going right from position P by distance D
                // We cross 0 at positions where (P + k) % 100 == 0 for k in [1, D]
                // This happens floor((P + D) / 100) - floor(P / 100) times
                const after = @divFloor(position + rotation.distance, 100);
                const before = @divFloor(position, 100);
                break :blk @as(i64, after - before);
            },
            .left => blk: {
                // Going left from position P by distance D
                // We cross 0 at positions where (P - k) % 100 == 0 for k in [1, D]
                // This happens when k = P, P+100, P+200, ... (all ≤ D)
                if (position == 0) {
                    // Starting at 0, we cross it again at k = 100, 200, ...
                    break :blk @as(i64, @divFloor(rotation.distance, 100));
                } else if (rotation.distance >= position) {
                    // We cross 0 at k = P, then every 100 clicks after
                    break :blk @as(i64, @divFloor(rotation.distance - position, 100) + 1);
                } else {
                    // We don't reach 0
                    break :blk 0;
                }
            },
        };

        count += clicks_through_zero;

        // Update position
        switch (rotation.direction) {
            .left => {
                position = mod(position - rotation.distance, 100);
            },
            .right => {
                position = mod(position + rotation.distance, 100);
            },
        }
    }

    return count;
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
