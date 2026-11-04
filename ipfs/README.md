# ipfs

Talking to a local IPFS node (for example Kubo) through its HTTP RPC API with `github.com/ipfs/go-ipfs-api`: add content, read it back, download it and publish it to IPNS.

The original exercise description is in [ipfs-go-client/README.md](ipfs-go-client/README.md).

## What the code does

- **`ipfs-go-client/app/`** (own `go.mod`, module `test`, `go 1.18`, `go-ipfs-api v0.3.0`):
  - `util.go`: `performChecks` refuses to continue unless `YourLocalPath` and `YourPublicKey` are filled in and a node answers on `localhost:5001`; `separator` prints a dashed line.
  - `main.go`: connects to `localhost:5001` and runs five steps, stopping on the first error:
    1. `sh.Add` the text `Hello from Launchpad!` and print the resulting CID.
    2. `sh.Cat("/ipfs/<cid>")` and print the content.
    3. `sh.Get(cid, YourLocalPath)` to download it to disk.
    4. `sh.PublishWithDetails(cid, YourPublicKey, 50h lifetime, 1µs TTL, resolve=true)` to publish to IPNS.
    5. `sh.Resolve(YourPublicKey)` and print what the IPNS name points to.
- **`ipfs-go-client/solution/main.go`**: the same `main.go` (identical logic, different formatting). It has no `go.mod` and no `util.go`.

## Concepts you can learn

- **Content addressing (CID)**: the ID is derived from the content, so the same bytes always yield the same CID.
- **IPFS HTTP RPC client**: the Go program does not run IPFS itself; it drives a running node over its API port.
- **IPNS**: a mutable name (bound to a key) that points to an immutable CID, so you can update what a name refers to; lifetime and TTL control record validity and caching.

## How to run

Prerequisites: a running IPFS node with its API on `localhost:5001`. Edit `app/main.go` first: set `YourLocalPath` (download destination) and `YourPublicKey` (the key used for IPNS publish/resolve). With either left empty the program only prints an error and exits.

```bash
cd ipfs/ipfs-go-client/app
go run .
```

## Notes

- The exercise README says `app/main.go` has functions marked `TO BE IMPLEMENTED`; in this repo they are already implemented.
- `solution/` cannot be built where it is: without its own `go.mod` it falls into the root module (`arena`), which does not require `go-ipfs-api`, and it lacks `util.go` (`separator`, `performChecks`). Copy its functions into `app/` instead, as the exercise README describes.
- `go-ipfs-api v0.3.0` is old; compatibility with current Kubo API versions was not verified.
