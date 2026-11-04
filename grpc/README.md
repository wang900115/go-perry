# grpc

A unary gRPC CRUD service for "person" records, with a server that keeps data in memory and a client that exercises every method.

## What the code does

One module (`grpc/go.mod`, module `grpc`) with three packages:

- **`proto/person.proto`**: service `PersonService` with `Create`, `Read`, `Update`, `Delete` and messages for name/email/phoneNumber. `person.pb.go` and `person_grpc.pb.go` are generated (protoc-gen-go v1.36.6, protoc-gen-go-grpc v1.5.1, protoc v6.30.1); do not edit them.
- **`server/main.go`**: listens on `:8080`. Persons are stored in a package-level `map[int32]Person` with an incrementing `nextID` starting at 1.
  - `Create`: all three fields required, else error `fields missing`.
  - `Read` / `Delete`: error `not found` for unknown IDs.
  - `Update`: `not found` for unknown IDs, `fields missing` if any field is empty; otherwise replaces the record and returns `updated`.
- **`client/main.go`**: connects to `localhost:8080` with insecure credentials and a 15 s timeout, then in order: creates "John Wick", reads it back by ID, updates it to "Luke Skywalker", deletes it. It prints `Person created: ...`, `Person details ...`, `Update response: updated`, `Delete response: deleted`, and exits via `log.Fatalf` on the first error.

## Concepts you can learn

- **Contract-first APIs**: the `.proto` file defines messages and RPCs; code for both sides is generated from it.
- **Implementing a service**: embedding `UnimplementedPersonServiceServer` keeps the server compiling if new RPCs are added.
- **Client stubs and deadlines**: `NewPersonServiceClient` gives typed methods; `context.WithTimeout` sets one 15 s deadline shared by all four calls; every RPC carries it to the server.
- **Insecure transport**: `insecure.NewCredentials()` disables TLS, which is only acceptable locally.

## How to run

```bash
cd grpc
go run ./server     # terminal 1, listens on :8080
go run ./client     # terminal 2
```

Regenerating the Go code requires `protoc`, `protoc-gen-go` and `protoc-gen-go-grpc`; the exact command used is not recorded in the repo.

## Notes

- The map and `nextID` are not protected by a mutex; concurrent RPCs would race.
- Errors are plain `errors.New`, so clients see status code `Unknown` instead of `NotFound` / `InvalidArgument`.
- On error the server also returns a non-nil response message; gRPC ignores it.
- `go.mod` marks all dependencies `// indirect`; it builds anyway.
