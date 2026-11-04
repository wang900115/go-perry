# paseto

Issues and verifies PASETO v2 tokens in both modes: `v2.local` (symmetric encryption) and `v2.public` (Ed25519 signatures). It uses `github.com/o1egl/paseto`.

## What the code does

`main.go` runs two round trips:

1. **v2.local**: generates a random 32-byte key and builds a `paseto.JSONToken`. The token has issuer `perry`, subject = the user ID, audience `example-app`, a random hex `jti`, `iat`=now, `nbf`=now-5s, and `exp`=now+30m, plus custom claims `uid=1` and `role=admin`. It encrypts the token with footer `k=v2.local.demo` and prints it. It then decrypts the token, validates it with `paseto.ValidAt(now)`, and prints the footer and the payload as indented JSON.
2. **v2.public**: generates an Ed25519 key pair and signs a token with the same structure (`uid=2`, `role=manager`, footer `k=v2.public.demo`). It verifies the token with the public key, validates the time claims, and prints the footer and payload. The public token string itself isn't printed.

Any error panics (the `must` helper).

## Concepts you can learn

- **local vs public tokens**: `local` encrypts the payload, so only holders of the shared key can read or create tokens. `public` signs the payload, which stays readable, so anyone with the public key can verify the token but only the private-key holder can issue one.
- **No algorithm negotiation**: the version and purpose (`v2.local.` / `v2.public.`) fix the algorithm. JWT has the `alg` header pitfalls that PASETO avoids this way.
- **Footer**: unencrypted but authenticated data. Here it carries a key hint (`k=...`).
- **Time-claim validation is explicit**: decrypting or verifying a token doesn't check `exp` or `nbf`. You have to call `Validate(ValidAt(...))` yourself.
- **The `nbf` skew allowance**: setting `nbf` 5 seconds in the past tolerates small clock differences between servers.

## How to run

```bash
cd paseto
go run .
```

You don't need any external services. Expected output: the `v2.local.` token, `llocal Footer: k=v2.local.demo`, the local payload JSON, `public Footer: k=v2.public.demo`, and the public payload JSON.

## Notes

- The label `llocal Footer` has a typo in the source.
- Ed25519 comes from `golang.org/x/crypto/ed25519`, pinned to a 2018 version. The standard library's `crypto/ed25519` provides the same functionality today.
- The `Claims` struct feeds `sub` (from `UserID`) plus the custom `uid` and `role` claims. Its own JSON tags aren't used.
