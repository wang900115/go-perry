#!/usr/bin/env bash
set -euo pipefail
mkdir -p bench_output

go test -bench=. -benchmem ./... -run=^$ -cpuprofile=bench_output/cpu.prof -memprofile=bench_output/mem.prof

echo "Benchmarks completed. Profiles in bench_output/"