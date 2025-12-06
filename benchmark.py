#!/usr/bin/env python3
import subprocess
import time
import os

def measure_time(command, cwd=None, runs=3):
    """Run command multiple times and return average time in milliseconds"""
    # Check if command exists (handling both absolute and relative paths)
    check_path = os.path.join(cwd, command[0]) if cwd else command[0]
    if not os.path.exists(check_path):
        return None

    times = []
    for _ in range(runs):
        start = time.perf_counter()
        try:
            subprocess.run(command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                         check=True, timeout=30, cwd=cwd)
            elapsed = (time.perf_counter() - start) * 1000  # Convert to ms
            times.append(elapsed)
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired, FileNotFoundError):
            return None

    return sum(times) / len(times) if times else None

# Results storage
results = []

print("Running benchmarks (3 runs each)...")
print()

for day in range(1, 6):
    day_str = f"{day:02d}"
    print(f"Day {day}...", end=" ", flush=True)

    rust_time = measure_time([f"target/release/day{day_str}"], cwd="rust")
    go_time = measure_time([f"../../target/day{day_str}-go"], cwd=f"go/day{day_str}")
    zig_time = measure_time([f"zig-out/bin/day{day_str}"], cwd="zig")

    results.append((day, rust_time, go_time, zig_time))
    print("done")

# Print results table
print("\n" + "="*60)
print("Benchmark Results (average of 3 runs)")
print("="*60)
print(f"{'Day':<6} {'Rust (ms)':<15} {'Go (ms)':<15} {'Zig (ms)':<15}")
print("-"*60)

for day, rust_time, go_time, zig_time in results:
    rust_str = f"{rust_time:.2f}" if rust_time is not None else "N/A"
    go_str = f"{go_time:.2f}" if go_time is not None else "N/A"
    zig_str = f"{zig_time:.2f}" if zig_time is not None else "N/A"
    print(f"{day:<6} {rust_str:<15} {go_str:<15} {zig_str:<15}")

print("="*60)

# Print markdown table
print("\nMarkdown Format:")
print("| Day | Rust (ms) | Go (ms) | Zig (ms) |")
print("|-----|-----------|---------|----------|")
for day, rust_time, go_time, zig_time in results:
    rust_str = f"{rust_time:.2f}" if rust_time is not None else "N/A"
    go_str = f"{go_time:.2f}" if go_time is not None else "N/A"
    zig_str = f"{zig_time:.2f}" if zig_time is not None else "N/A"
    print(f"| {day} | {rust_str} | {go_str} | {zig_str} |")
