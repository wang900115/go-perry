# unsafe

Raw memory access with `unsafe`: reinterpreting pointers, inspecting struct layout and padding, pointer arithmetic, and zero-copy string/byte conversions.

## What the code does

`main` only calls `Dump` on `User{A: 1, B: 0x11223344, C: 2}` (fields `int8, int32, int8`). `Dump` views the struct's memory as a `[]byte` with `unsafe.Slice` and prints it 8 bytes per row. A run on amd64/Windows printed:

```
Raw memory (12 bytes):

0x00: 01 7F 00 00 44 33 22 11
0x08: 02 0E 00 00
```

So `A` sits at offset 0, followed by 3 padding bytes. `B` is at offset 4 in little-endian order (`44 33 22 11`). `C` is at offset 8, followed by 3 padding bytes, for a total size of 12.

Other functions are defined but **never called** by `main`:

- `voidPointer`: reads an `int` through `unsafe.Pointer` as `*int` and then as `*float64`, reinterpreting the same bits.
- `memoryAlign`: prints each field's offset, size and alignment with `reflect`, plus the struct size.
- `pointerOffset`: walks an `[5]int` array by adding `i*Sizeof(elem)` to a `uintptr` base.
- `zeroCopyByteSliceToString` / `zeroCopyStringToByteSlice`: convert without copying using `unsafe.String`/`unsafe.SliceData` and `unsafe.Slice`/`unsafe.StringData`.

## Concepts you can learn

- **Alignment and padding**: fields are aligned to their size, so `int8, int32, int8` takes 12 bytes. Ordering fields from largest to smallest (`int32, int8, int8`) would reduce this to 8. This matters for large slices of structs.
- **Endianness**: the dump shows amd64 storing the least significant byte first.
- **`unsafe.Pointer` vs `uintptr`**: the GC tracks an `unsafe.Pointer`. A `uintptr` is just a number, so arithmetic must convert back to `unsafe.Pointer` in the same expression, as `pointerOffset` does.
- **Type punning**: casting `*int` to `*float64` reads the same 8 bytes under a different interpretation, without any conversion.
- **Zero-copy conversions**: these avoid allocation when converting between `[]byte` and `string`, at the cost of giving up Go's immutability guarantees.

## How to run

```bash
cd unsafe
go run .
```

To try the other functions, call them from `main`.

## Notes

- Padding bytes are not guaranteed to be zero and vary between runs (one run showed `7F`/`0E`, another `09`).
- Writing to the `[]byte` returned by `zeroCopyStringToByteSlice` is undefined behavior, because string data may be in read-only memory. Likewise, the string from `zeroCopyByteSliceToString` changes if the source slice is mutated.
- The module is named `un`. Comments are in Chinese.
