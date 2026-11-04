# p2p

Placeholder for peer-to-peer networking experiments. **There is no Go source code in this directory yet.**

## What the code does

Nothing to run. The directory only contains:

- `go.mod` (module `p2p`, `go 1.24.4`) requiring `github.com/libp2p/go-libp2p v0.42.1` and `github.com/ipfs/go-ipfs-api v0.7.0`.
- `go.sum` for those dependencies.

The dependency list suggests the intended direction (libp2p hosts/streams and an IPFS HTTP API client), but no code exercises either library.

## How to run

```bash
cd p2p
go build ./...   # prints: go: warning: "./..." matched no packages
```
