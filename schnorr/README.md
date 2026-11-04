# schnorr

A toy implementation of the interactive Schnorr identification protocol over a tiny group (`p = 23`, `g = 5`), using `math/big`. A prover shows that it knows a private key `x` without revealing it.

## What the code does

`main.go` runs one round of the protocol:

1. `NewProver` picks a private key `x` in `[1, p-2]` and computes the public key `y = g^x mod p`.
2. **Commit**: the prover picks a random `r` in `[0, p-2]` and sends `t = g^r mod p`.
3. **Challenge**: the verifier picks a random `c` in `[0, p-2]`.
4. **Response**: the prover sends `s = (r + c·x) mod (p-1)`.
5. **Verify**: the verifier checks `g^s ≡ t · y^c (mod p)`.

The program prints `y`, `t`, `c`, `s`, and `Verify result = true`. The numbers change on every run.

## Concepts you can learn

- **The commit, challenge, response structure (a Sigma protocol)**: this is the basis of zero-knowledge proofs of knowledge. The prover must commit to `t` before it sees `c`.
- **Why verification works**: `g^s = g^(r + c·x) = g^r · (g^x)^c = t · y^c`. Exponents are reduced modulo the group order (`p-1`), not modulo `p`.
- **The discrete-log assumption**: anyone can compute `y` from `x`, but recovering `x` from `y` is assumed to be hard in a large group.
- **`crypto/rand` + `math/big`** for modular exponentiation (`Exp`) and modular arithmetic.

## How to run

```bash
cd schnorr
go run .
```

You don't need any external services.

## Notes

- **This is for learning only and isn't secure.** With `p = 23`, `x` can be brute-forced instantly. Real deployments use large prime-order groups or elliptic curves.
- The protocol here is interactive identification, not a signature scheme. Turning it into a Schnorr *signature* would derive `c` from a hash of `t` and the message (Fiat-Shamir), and the code doesn't do that. The comment at the top of the file calls it "Schnorr signatures".
- `c` can be `0`. Then `s = r` and the round proves nothing about `x`. A real verifier would draw `c` from a large range.
- The errors from `rand.Int` are ignored.
