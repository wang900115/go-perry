# flatbuffer

Building and reading a FlatBuffers message in Go. FlatBuffers is a schema-based binary format whose fields can be read straight from the buffer without a decode step.

## What the code does

- `schema/user.fbs` - schema: `namespace example; table User { id:uint64; name:string; age:uint8; active:bool; }`, with `root_type User`.
- `generated/example/User.go` - **generated** by the FlatBuffers compiler (`flatc`) from `schema/user.fbs`. Do not edit it by hand.
- `main.go`:
  1. Creates a `flatbuffers.Builder`, creates the `name` string first, then fills a `User` table between `UserStart` and `UserEnd`, and calls `Finish`.
  2. Prints the buffer size and hex bytes.
  3. Reads it back with `example.GetRootAsUser(buf, 0)` and prints each field.

Actual output:

```
Binary size: 56 bytes
Binary data: 14000000000000000c0018000c000800070006000c0000000000011910000000e90300000000000000000000050000005065727279000000

User:
ID: 1001
Name: Perry
Age: 25
Active: true
```

## Concepts you can learn

- **Zero-copy access**: `GetRootAsUser` does not parse anything. Accessors such as `Id()` and `Name()` read through the vtable at the field offset in `buf`, so reads cost almost nothing and field reads do not allocate (only the small `*User` wrapper is allocated) (`Name()` returns a `[]byte` view).
- **Build order matters**: strings and other child objects must be created before `UserStart`, because a table cannot be nested inside another object while it is being built.
- **Vtables and optional fields**: each table has a vtable of field offsets, which supports schema evolution and default values. The cost is extra bytes: 56 bytes here for a small record.
- **Schema + code generation**: the `.fbs` file is the contract, and `flatc` generates typed accessors.

## How to run

Own `go.mod` (module `flatbuffer`, depends on `github.com/google/flatbuffers`):

```bash
cd flatbuffer
go run .
```

To regenerate the Go code after changing the schema, install `flatc` (not needed just to run the example). For example, from `flatbuffer/`:

```bash
flatc --go -o generated schema/user.fbs
```

This command is not verified here. Since the namespace is `example`, it should write `generated/example/User.go`.
