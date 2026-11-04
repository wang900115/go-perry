# assembly

Notes on Go's Plan 9-style assembly (amd64). The directory contains only commented-out examples; there is no runnable program.

## What the code does

Every file is `package main` with no executable code:

- `branch.go`: a Go `if x > y {}` next to the assembly it lowers to: `CMPQ AX, BX` followed by a conditional jump `JLE end`.
- `cpu_register.go`: `MOVQ $3, AX`, `MOVQ $5, BX`, `ADDQ BX, AX`, `MOVQ AX, Z` (compute 3 + 5 in registers and store the result), plus a list of register names (`AX, BX, CX, DX, SP, BP, PC`).
- `memory.go`: the pseudo-registers `(SP) (BP) (PC)` and a hand-written `Add(a, b int) int` that reads its arguments from `a+0(FP)` and `b+8(FP)` and writes the result to `ret+16(FP)`.
- `cpu_instruction.go`: empty (only `package main`).

## Concepts you can learn

- **Plan 9 operand order**: the source comes first and the destination last (`MOVQ $3, AX` loads 3 into AX). Intel syntax uses the opposite order, so keep this in mind when reading `go tool objdump` output. `CMP` is the exception: `CMPQ AX, BX` compares AX to BX in written order, so `JLE` jumps when AX <= BX.
- **Instruction size suffixes**: the `Q` in `MOVQ`/`ADDQ`/`CMPQ` means a 64-bit quadword, which matches `int` on amd64.
- **Branches are compare + conditional jump**: an `if` becomes `CMPQ` plus a jump that skips the body when the condition is false (`JLE` for `x > y`).
- **The FP pseudo-register**: in hand-written assembly, arguments and return values are addressed as `name+offset(FP)`. With 8-byte ints, `a` is at offset 0, `b` at 8 and the return value at 16.

## How to run

There is nothing to run. The module has its own `go.mod`, but `go build` fails with `function main is undeclared in the main package`. To see real compiler output, run `go build -gcflags=-S .` in another topic, such as `bce/`.

## Notes

- The `memory.go` example uses the stack-based calling convention (ABI0) that hand-written `.s` files use. Since Go 1.17, compiled Go functions on amd64 pass arguments in registers (ABIInternal), so compiler output will look different.
- No `.s` file exists, so none of these snippets are assembled or tested.
