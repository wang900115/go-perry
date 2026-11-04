# messagepack

Serializing a Go struct to MessagePack (a compact binary format similar to JSON) with `github.com/vmihailenco/msgpack/v5`, then decoding it back.

## What the code does

`main.go`:

- Defines `User{ID, Name, Email, Age}` with `msgpack:"..."` struct tags for the field names.
- Marshals `User{1001, "Perry", "perry@example.com", 25}`, prints the bytes in hex and the size, then unmarshals into a new `User` and prints it.

Actual output:

```
MessagePack bytes: 84a26964cd03e9a46e616d65a55065727279a5656d61696cb17065727279406578616d706c652e636f6da361676519
Size: 47 bytes
Decoded: {ID:1001 Name:Perry Email:perry@example.com Age:25}
```

## Concepts you can learn

- **Self-describing binary encoding**: the struct is encoded as a map (`84` = map of 4 entries) whose keys are the tag names (`a2 6964` = 2-byte string `"id"`). No schema is needed to decode it. The trade-off is that field names are repeated in every message.
- **Compact integer and string encodings**: `25` fits in one byte (`19`, positive fixint), `1001` uses `cd 03e9` (uint16), and short strings carry their length in the prefix byte. This is why MessagePack output is usually smaller than the equivalent JSON.
- **Struct tags control the wire names**, as with `encoding/json`.

## How to run

Own `go.mod` (module `messagepack`):

```bash
cd messagepack
go run .
```

## Notes

- `go.mod` marks `msgpack/v5` as `// indirect` even though `main.go` imports it directly. It still builds; `go mod tidy` would fix the annotation.
