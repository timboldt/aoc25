const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Helper function to create and configure an executable
    inline for ([_]struct { name: []const u8, src: []const u8 }{
        .{ .name = "day01", .src = "src/day01.zig" },
        .{ .name = "day02", .src = "src/day02.zig" },
        .{ .name = "day03", .src = "src/day03.zig" },
    }) |day| {
        const exe = b.addExecutable(.{
            .name = day.name,
            .root_module = b.createModule(.{
                .root_source_file = b.path(day.src),
                .target = target,
                .optimize = optimize,
            }),
        });
        b.installArtifact(exe);
        const run_cmd = b.addRunArtifact(exe);
        const run_step = b.step(day.name, b.fmt("Run {s}", .{day.name}));
        run_step.dependOn(&run_cmd.step);

        // Add test for this day
        const tests = b.addTest(.{
            .root_module = b.createModule(.{
                .root_source_file = b.path(day.src),
                .target = target,
                .optimize = optimize,
            }),
        });
        const test_step = b.step(b.fmt("test-{s}", .{day.name}), b.fmt("Test {s}", .{day.name}));
        test_step.dependOn(&b.addRunArtifact(tests).step);
    }

    // Add a master test step
    const test_step = b.step("test", "Run all tests");
    for ([_][]const u8{ "day01", "day02", "day03" }) |day| {
        const tests = b.addTest(.{
            .root_module = b.createModule(.{
                .root_source_file = b.path(b.fmt("src/{s}.zig", .{day})),
                .target = target,
                .optimize = optimize,
            }),
        });
        test_step.dependOn(&b.addRunArtifact(tests).step);
    }
}
