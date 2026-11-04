# breakcircut

A minimal circuit breaker state machine (`Open` / `Closed` / `HalfOpen`) protected by a mutex. The directory name is misspelled on purpose; it is the real directory name.

## What the code does

`main.go` (the only file) defines:

- `Breaker` with a failure threshold, a success threshold and a reset timeout.
- `Allow()`: in `Open`, rejects with an error until `resetTimeout` has passed since the last failure, then moves to `HalfOpen` and lets the call through. In any other state it allows the call.
- `Success()`: in `HalfOpen`, counts successes and moves to `Closed` once `successThreshold` is reached. In `Closed`, resets the failure counter (so only consecutive failures count).
- `Failure()`: in `HalfOpen`, goes straight back to `Open`. In `Closed`, counts failures and opens after `failThreshold`.
- `unreliable(pFail)`: a fake service that fails with probability `pFail`, otherwise sleeps 0-49 ms.

There is **no `func main`**, so nothing runs and there is no output. `unreliable` and `NewBreaker` are never called.

## Concepts you can learn

- **Circuit breaker states**: stop calling a failing dependency (`Open`), probe it carefully (`HalfOpen`), resume normal traffic (`Closed`). This prevents a slow or broken service from dragging down its callers.
- **Consecutive-failure counting**: resetting `failureCount` on success means sporadic errors do not trip the breaker; only a run of failures does.
- **Time-based recovery**: `lastFailTime` + `resetTimeout` decides when a probe is allowed.
- **Mutex-guarded state transitions**: every method locks, so concurrent callers see consistent state.

## How to run

`breakcircut/` has no `go.mod`; it belongs to the root module (`module arena`). From the repo root:

```bash
go vet breakcircut/main.go   # type-checks
```

`go run breakcircut/main.go` fails with `function main is undeclared in the main package`. To try it, add a `main` that loops: `Allow()` -> `unreliable(p)` -> `Success()`/`Failure()`, printing the state.

## Notes

- `Open` is the zero value of `State`, and `NewBreaker` starts in `Open`. Because `lastFailTime` is zero, the first `Allow()` immediately moves to `HalfOpen`. Conventional breakers start `Closed`.
- `HalfOpen` lets every caller through (no limit on concurrent probe requests).
- `State` has no `String()` method, so printing it shows `0/1/2`.
