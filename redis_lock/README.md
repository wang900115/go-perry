# redis_lock

The simplest Redis-based mutual exclusion: `SET key value NX` with an expiry, using `github.com/redis/go-redis/v9`.

## What the code does

`main.go` connects to Redis at `localhost:6379` and tries to acquire the key `my_lock` with `SetNX(ctx, "my_lock", "1", 20s)`.

- If the key was set: prints `Lock acquired succeessfully`, sleeps 20 seconds (the "work"), prints `Work done!`, deletes the key and prints `Lock released!`.
- If the key already exists: prints `Failed to acquire lock, Resource is already locked`.
- If Redis returns an error: prints `Error acquiring lock: ...` and treats it as not acquired.

## Concepts you can learn

- **`SETNX` as a lock**: only one client can create the key, so only one process enters the critical section.
- **Lock TTL**: the expiry ensures a crashed holder does not keep the lock forever.
- **Release by deletion**: removing the key lets the next client acquire it.

## How to run

Prerequisite: Redis on `localhost:6379` (no password).

```bash
cd redis_lock
go run .     # terminal 1: acquires the lock and holds it for 20s
go run .     # terminal 2, within 20s: "Failed to acquire lock..."
```

## Notes

- The work (20 s sleep) is as long as the TTL (20 s), so the lock can expire just before it is released, and another client could get in.
- Release is a plain `DEL` with no ownership check (the value is always `"1"`). If the lock expired and someone else took it, this process deletes **their** lock. Safer variants store a unique token and delete via a compare-and-delete Lua script.
- No retry/wait: a client that fails to acquire just exits.
