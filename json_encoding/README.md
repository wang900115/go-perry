# json_encoding

`encoding/json` (v1) beyond the basics: struct tags (`omitempty`, `-`, `,string`),
Go 1.24 `omitzero`, custom `MarshalJSON`/`UnmarshalJSON`, `json.RawMessage` for
tagged unions, strict decoding with `DisallowUnknownFields`, `*json.UnmarshalTypeError`,
and why large integers break when decoded into `interface{}` (`UseNumber`).

```bash
go run .
```

Standard library only. This module declares `go 1.24` (`omitzero` needs Go 1.24).
Verified on go1.26 linux/amd64 (2 runs, identical output); `go vet .` and `GOOS=windows go vet .` pass.

## Expected output

```text
== Part 1: struct tags ==
{"id":1,"name":"perry","age":"30"}
err=<nil> name="Bob" password="" age=41
== Part 2: omitempty vs omitzero ==
{"name":"deploy","empty_at":"0001-01-01T00:00:00Z","zero_tags":[]}
{"name":"nil slices","empty_at":"0001-01-01T00:00:00Z"}
== Part 3: custom Marshaler / Unmarshaler ==
{"name":"backup","timeout":"1m30s"}
err=<nil> timeout=2h15m0s
err: duration must be a string: json: cannot unmarshal number into Go value of type string
== Part 4: RawMessage tagged union ==
main.Click {X:10 Y:20} err=<nil>
main.KeyPress {Key:Enter} err=<nil>
<nil> <nil> err=unknown event type "scroll"
== Part 5: strict decoding ==
lenient err: <nil>
strict err:  json: unknown field "admin"
type error: field=id value=string want=int
== Part 6: numbers in interface{} ==
default: float64 9007199254740992
UseNumber: json.Number 9007199254740993 -> int64 9007199254740993
```
