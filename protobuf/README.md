# protobuf

Protocol Buffers in Go: a `Person` message defined in a `.proto` file, serialized to binary, written to disk, and exchanged over a small HTTP API.

## What the code does

- `models/person.proto` - proto3 schema: enum `PhoneType {MOBILE, HOME, WORK}`, message `PhoneNumber {number, type}`, and message `Person {name, id, email, repeated phones}`.
- `models/person.pb.go` - **generated** by `protoc` with `protoc-gen-go` (v1.36.6, protoc v6.30.1, per its header) from `models/person.proto`. Do not edit it by hand.
- `main.go` - HTTP server on `:8080` with an in-memory `map[int32]*Person`:
  - `POST /add`: the body is a binary-encoded `Person`. It is unmarshaled, stored by `Id`, and the server replies with `Person added <text form>`.
  - `GET /get?id=N`: marshals the stored person and returns raw bytes (`application/octet-stream`), or 404 `Person not founded`.
- `person_marshal.go` - builds `Person{Name: "John Wick", Id: 1234, Email: "wick@codehemi.io", Phones: [{"123-456-789", MOBILE}]}`, marshals it, and writes `tmp/person.bin`.
- `person_unmarshal.go` - reads `tmp/retrieved_person.bin`, unmarshals it, and prints `Deserialized Person: ...`.
- `tmp/person.bin`, `tmp/retrieved_person.bin` - committed sample files (47 bytes each).

The intended flow (inferred from the file names, not written down in the code): marshal to `person.bin`, POST it to `/add`, GET it back into `retrieved_person.bin`, and unmarshal that.

## Concepts you can learn

- **Schema-first, generated types**: the `.proto` file is the contract. Field numbers (`= 1`, `= 2`, ...) are what go on the wire, not field names, so the messages are compact and fields can be renamed safely.
- **Binary encoding with `proto.Marshal` / `proto.Unmarshal`**: smaller and faster than JSON, but the bytes cannot be read without the schema.
- **proto3 defaults**: enum value `0` (`MOBILE`) is the default and is not written to the wire when it has that value.
- **Protobuf over plain HTTP**: no gRPC needed. The body is just bytes, and the content type is set to `application/octet-stream`.

## How to run

Own `go.mod` (module `protobuf`). The three `.go` files in the module root are all `package main` with their own `func main`, so `go run .` / `go build .` will fail (`main` redeclared). Run one file at a time from inside `protobuf/` (the `tmp/` paths are relative to the current directory):

```bash
cd protobuf
go run main.go              # server on :8080
go run person_marshal.go    # writes tmp/person.bin
go run person_unmarshal.go  # reads tmp/retrieved_person.bin
```

Example round-trip against the server (with curl):

```bash
curl -X POST --data-binary @tmp/person.bin http://localhost:8080/add
curl -o tmp/retrieved_person.bin "http://localhost:8080/get?id=1234"
```

To regenerate `person.pb.go` you need `protoc` and `protoc-gen-go`. For example, from `protobuf/`: `protoc --go_out=. --go_opt=paths=source_relative models/person.proto` (not verified here).

## Notes

- The server's map has no mutex, and concurrent requests can race on it.
- Several status codes look wrong: a read or unmarshal failure and a missing `id` return `502 Bad Gateway` rather than `400`. A non-numeric `id` only prints to stdout and returns an empty `200`. The marshal-failure message says "unmarshal".
- The `/add` reply is text, while `/get` returns binary.
- Each file builds individually. `go run person_unmarshal.go` prints `Deserialized Person: name:"John Wick" id:1234 email:"wick@codehemi.io" phones:{number:"123-456-789"}`.
