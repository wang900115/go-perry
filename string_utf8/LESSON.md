# 2026-11-19 · Strings, bytes, runes and UTF-8

## Why learn this
- Every program handles text: HTTP headers, JSON, log lines, user names, file paths. Bugs appear as soon as the text is not ASCII (Chinese names, emoji, accented letters).
- Typical real bugs: a length check that counts bytes instead of characters, a "preview" that cuts `s[:20]` in the middle of a character and produces broken output, a loop that does `s += x` and becomes slow, two strings that look equal but are not.
- It builds on `slice_internals`: a string is like a read-only `[]byte` (pointer + length, no capacity), and it explains why `bufio.Scanner.Text()` allocates while `Bytes()` does not (`bufio_scanner`).

## Core concepts
1. **A string is an immutable sequence of bytes** (https://go.dev/ref/spec#String_types). Its header is two words: pointer and length. `len(s)` is the number of bytes; `s[i]` is a `byte`. Strings do not have to contain valid UTF-8.
2. **Go source is UTF-8, so string literals are UTF-8** (https://go.dev/blog/strings). `"語"` is 3 bytes: `e8 aa 9e`.
3. **A rune is an `int32` holding a Unicode code point.** `for i, r := range s` decodes one UTF-8 sequence per iteration; `i` is the byte offset, not the character index (https://go.dev/ref/spec#For_range).
4. **Invalid bytes become `utf8.RuneError` (U+FFFD) with width 1** in `range`, `utf8.DecodeRuneInString` and `[]rune(s)`. Use `utf8.ValidString` to check and `strings.ToValidUTF8` to clean (https://pkg.go.dev/unicode/utf8).
5. **Conversions:** `[]byte(s)` and `string(b)` copy the data (strings are immutable, slices are not). `[]rune(s)` decodes into 4 bytes per rune. `string(r)` encodes a rune; an invalid code point becomes `"\uFFFD"` (https://go.dev/ref/spec#Conversions). `string(65)` on an `int` gives `"A"`, not `"65"`; `go vet` reports it (stringintconv). Use `strconv.Itoa` for numbers.
6. **Slicing is by bytes.** `s[:n]` may cut a rune in half. Truncate by runes (`range`) or move back to a rune start with `utf8.RuneStart`.
7. **Building strings:** `s += x` creates a new string every time (O(n²) in a loop). `strings.Builder` appends into a growing `[]byte`; `Grow` preallocates (https://pkg.go.dev/strings#Builder).
8. **A rune is not a visible character.** `"é"` can be one code point (U+00E9) or two (`e` + U+0301); they look the same but `==` is false. A flag emoji is two runes. Normalization needs `golang.org/x/text/unicode/norm` (not in the standard library); grapheme clusters need a third-party package.

### How it works (runtime)
Based on go1.26 `src/runtime/string.go` and `src/runtime/utf8.go`; implementation details, not spec, may change.
- `range` over a string is compiled into a loop that handles ASCII (`< 0x80`) inline and calls `runtime.decoderune` only for multi-byte sequences.
- `[]byte(s)` calls `stringtoslicebyte`, `string(b)` calls `slicebytetostring`; both allocate and copy. When the result is small and does not escape, the compiler uses a 32-byte stack buffer instead of the heap.
- The compiler skips the copy in some cases where it can prove it is safe, e.g. `m[string(b)]` map lookups, `string(b) == "x"` comparisons, and `for i, c := range []byte(s)`. I am not certain of the complete list; check with `go build -gcflags=-m` or `testing.AllocsPerRun`.
- `a + b + c` is one call to `concatstrings` (one allocation for the whole expression), but each loop iteration still allocates a new string and copies everything before it.
- `strings.Builder.String()` returns a string that shares the builder's buffer through `unsafe.String`, without copying. That is why a `Builder` must not be copied after first use (it panics).

## Full example
See `main.go`. Verified with `go run .` (go1.26.3 and go1.24.7, identical output) and `go vet .` on linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. `len` vs rune count, indexing returns a byte, `range` gives byte offsets and runes.
2. A string with an invalid byte, `U+FFFD`, `ValidString`, `ToValidUTF8`.
3. `[]byte` copy, `[]rune` decode, `string(rune)`, invalid code points.
4. Naive `s[:4]` cuts a rune; `truncateRunes` and `truncateBytes` fix it.
5. Allocations of `+=` vs `strings.Builder` with `Grow` (via `testing.AllocsPerRun`).
6. Precomposed vs combining `é`, a two-rune flag, `strings.EqualFold`.

## Common pitfalls
1. **Using `len(s)` as a character count.** A "max 10 characters" check with `len` rejects a 4-character Chinese name (12 bytes). Use `utf8.RuneCountInString`, and decide whether you really mean bytes (storage limits), runes, or visible characters.
2. **Cutting strings by byte index.** `s[:100]` for a preview or log truncation can leave half a rune, which shows as `\uFFFD` or breaks JSON/DB validation. Truncate on rune boundaries.
3. **`string(i)` with an integer.** `string(65)` is `"A"`, and `string(-1)` is `"\uFFFD"`. Use `strconv.Itoa` / `fmt.Sprint` for numbers; `go vet` catches the `int` case.

## Exercise (no answer given)
Write a package `textutil` with tests:
- `Reverse(s string) string` that reverses by runes. Test it with `"Go語言"` and with a string that contains an invalid byte. What should happen to the invalid byte? Write down your decision.
- `Ellipsis(s string, maxBytes int) string` that returns `s` if it fits, otherwise cuts on a rune boundary and appends `"…"` so that the result (including `…`, which is 3 bytes) is at most `maxBytes` bytes.
- `Join(parts []string, sep string) string` with `strings.Builder` and a correct `Grow` size. Compare it with a `+=` version in a benchmark (`b.ReportAllocs()`) for 10, 100 and 1000 parts.
- Bonus: find a case where `Reverse` gives a visibly wrong result even though every rune is valid (hint: Part 6). Explain why.

## Suggested directory
`string_utf8`
