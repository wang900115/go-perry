# 2026-11-14 · encoding/json: tags, omitzero, custom (un)marshalers, RawMessage

## Why learn this
- Every HTTP API, config file, webhook and message payload you touch in Go goes through `encoding/json`.
- Most production JSON bugs are silent: a field that is never sent because it is unexported, a `time.Time` that shows up as `0001-01-01T00:00:00Z`, a typo in a request body that is quietly ignored, or an order ID that changes its last digit.
- It builds on earlier lessons: `error_wrapping` (`errors.As` with `*json.UnmarshalTypeError`), `io_reader_writer` (`json.NewDecoder` takes an `io.Reader`), `reflect` (encoding/json is built on reflection).

## Core concepts
1. **Only exported fields are visible.** The field name is the key unless a tag says otherwise: `json:"name,opt1,opt2"`. Options: `omitempty`, `omitzero` (Go 1.24+), `string`, and `-` to skip the field (https://pkg.go.dev/encoding/json#Marshal).
2. **Unmarshal matches keys case-insensitively** and ignores unknown keys by default (https://pkg.go.dev/encoding/json#Unmarshal).
3. **`omitempty` vs `omitzero`**: `omitempty` drops `false`, `0`, `nil` pointer/interface, and empty array/slice/map/string. A struct is never "empty", so `time.Time` is always written. `omitzero` drops the zero value of any type, and if the type has an `IsZero() bool` method it uses that (https://go.dev/doc/go1.24#encodingjsonpkgencodingjson).
4. **Custom encoding**: implement `json.Marshaler` (`MarshalJSON() ([]byte, error)`) and `json.Unmarshaler` (`UnmarshalJSON([]byte) error`). Put `UnmarshalJSON` on the pointer receiver because it must modify the value.
5. **`json.RawMessage`** keeps a value as raw bytes, so you can read a `type` field first and decode `data` into the right struct later (https://pkg.go.dev/encoding/json#RawMessage).
6. **Strict decoding**: `Decoder.DisallowUnknownFields()` turns unknown keys into an error (https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields).
7. **Numbers in `interface{}` become `float64`**, which has 53 bits of mantissa; integers above 2^53 lose precision. `Decoder.UseNumber()` gives you `json.Number` (a string) instead (https://pkg.go.dev/encoding/json#Decoder.UseNumber).

### How it works (standard library)
Based on `src/encoding/json/encode.go` and `decode.go` in go1.26; internal details may change.
- The first time a type is encoded, the package walks its fields with `reflect`, parses tags, and builds an encoder function. The result is cached per type in a `sync.Map`, so later calls skip the tag parsing.
- For each value, it first checks whether the type (or its pointer, if addressable) implements `json.Marshaler` or `encoding.TextMarshaler`; if so, it calls that method instead of the default reflection path. The output of `MarshalJSON` is validated and compacted.
- On decode, `Unmarshal` first scans the whole input to check it is valid JSON, and only then fills the value. That is why a syntax error leaves the target untouched, but a type error can leave it partly filled (it keeps going and returns the first such error).
- There is an experimental `encoding/json/v2` behind `GOEXPERIMENT=jsonv2` (Go 1.25+). This lesson does not use it; I am not sure when or whether it becomes the default.

## Full example
See `main.go`. Verified with `go run .` (2 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. Struct tags: `-` hides `Password`, `,string` writes `Age` as `"30"`, case-insensitive keys on decode.
2. `omitempty` keeps a zero `time.Time`; `omitzero` drops it. `omitzero` keeps `[]string{}` but drops `nil`.
3. A `Duration` type that encodes as `"1m30s"` with `MarshalJSON`/`UnmarshalJSON`.
4. A tagged union decoded with `json.RawMessage`.
5. `DisallowUnknownFields` and `errors.As` with `*json.UnmarshalTypeError`.
6. `9007199254740993` becomes `...992` as `float64`; `UseNumber` keeps it exact.

## Common pitfalls
1. **Unexported or untagged fields**: lowercase fields are silently skipped in both directions, no error. Also, `json:"-"` only hides a field from JSON; it does not make the struct safe to log with `%+v`.
2. **`omitempty` on structs and `time.Time`**: it never omits them. Use `omitzero` (Go 1.24+) or a pointer (`*time.Time`).
3. **Infinite recursion in `MarshalJSON`**: calling `json.Marshal(v)` inside `func (v T) MarshalJSON()` calls itself forever. Convert to a type without the method first (`type alias T; json.Marshal(alias(v))`) or marshal a different type as in Part 3.
4. **Decoding into `map[string]any`**: every number is a `float64`. Large IDs (Twitter/Snowflake-style) silently change. Use a typed struct field (`int64`), `UseNumber`, or `,string`.

## Exercise (no answer given)
Write a config loader for this JSON:

```json
{"name":"api","port":8080,"timeout":"30s","started_at":"2026-11-14T08:00:00Z","features":{"type":"cache","data":{"ttl":"5m","size":128}}}
```

Requirements:
- Reject unknown fields (try adding `"prot":8080`).
- `timeout` and `ttl` use your `Duration` type.
- `started_at` is optional; when you marshal the config back, it must disappear when unset (do it once with `omitzero` and once with `*time.Time`; compare).
- `features` is a tagged union with at least two types (`cache`, `ratelimit`).
- Marshal the config back and check that `Unmarshal(Marshal(cfg))` gives an equal value (you may need `reflect.DeepEqual`; think about why `time.Time` may need `Equal` instead).

## Suggested directory
`json_encoding`
