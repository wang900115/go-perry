# linkname-demo

Use `//go:linkname` to call the unexported runtime function `runtime.nanotime`.

## What the code does

`main.go` declares `func nanotime() int64` without a body, binds it to `runtime.nanotime` with `//go:linkname`, and prints the result with `println`, which writes to stderr. The output is a single integer, such as `231070078600`.

## Concepts you can learn

- **`//go:linkname` pull**: a body-less declaration plus `import _ "unsafe"` (required for the directive) lets the linker resolve the local name to a symbol in another package.
- **Monotonic clock**: `runtime.nanotime` reads the runtime's monotonic clock. The value is not wall-clock time and has no fixed epoch, so it only makes sense for measuring intervals.
- **Linker allowlist**: since Go 1.23 the linker blocks linkname pulls of most runtime symbols. `runtime.nanotime` still links (verified on Go 1.27.1), because widely used packages depend on it and the runtime explicitly allows it. Compare `coroutine/`, where `runtime.newcoro` is rejected.

## How to run

```bash
cd linkname-demo
go run .
```

## Notes

- Linkname into the runtime depends on internals that can change between releases. Use `time.Now()` / `time.Since()` in real code, which already use the monotonic clock.
- The module is named `linkdemo-demo`, which does not match the directory name.
