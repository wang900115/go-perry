#!/usr/bin/env bash
set -euo pipefail
mkdir -p bench_output
/go/bin/go test -bench=. -benchmem ./... -run=^$ -cpuprofile=bench_output/cpu.prof -memprofile=bench_output/mem.prof
# Optional: generate flamegraph if flamegraph tools available
# go tool pprof -svg bench_output/cpu.prof > bench_output/cpu.svg

echo "Benchmarks completed. Profiles in bench_output/"
