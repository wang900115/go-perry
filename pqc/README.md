# pqc

Post-quantum cryptography with Go's standard library: ML-DSA signatures (`crypto/mldsa`) and ML-KEM key encapsulation (`crypto/mlkem`). The directory has two independent programs.

## What the code does

- `digitalsignature/main.go`:
  1. Generates an ML-DSA-44 key pair with `mldsa.GenerateKey(mldsa.MLDSA44())`.
  2. Signs a fixed message with empty `mldsa.Options{}`.
  3. Prints the signature size and its first 32 bytes, then verifies the signature with the public key. The size observed in a run is 2420 bytes.
  4. Prints `Signature verification succeeded!`.

  A commented-out block shows how to pass a signing `Context` string.
- `sharedsecret/main.go`: an Alice/Bob ML-KEM-768 exchange.
  1. Alice generates a decapsulation key and serializes her encapsulation (public) key (1184 bytes).
  2. Bob parses Alice's key and calls `Encapsulate()`, which returns a 32-byte shared key and a 1088-byte ciphertext.
  3. Alice calls `Decapsulate(ciphertext)` to get the shared key.
  4. The program checks that both keys are equal and prints a success message.

## Concepts you can learn

- **KEM vs key exchange**: with ML-KEM, one side *encapsulates* a random secret to the other's public key, instead of both sides contributing a share as in Diffie-Hellman. The result is a 32-byte symmetric key for use with an AEAD cipher.
- **Post-quantum signatures**: ML-DSA (FIPS 204) replaces RSA and ECDSA where you need resistance to quantum attacks. Its signatures are much larger: 2420 bytes here, compared with 64 bytes for Ed25519.
- **Key sizes and serialization**: `EncapsulationKey().Bytes()` and `NewEncapsulationKey768(bytes)` show how a public key goes over the wire.
- **Implicit rejection**: `Decapsulate` only errors on a ciphertext of the wrong length. A tampered ciphertext of the right length returns a different, pseudo-random shared key without any error (FIPS 203 implicit rejection). Tampering only shows up later, when the keys don't match or AEAD decryption fails. (The comment in `sharedsecret/main.go` claiming otherwise is wrong.)

## How to run

`pqc/` has no `go.mod`, so it belongs to the root module (`go 1.26.3`). Run the programs from the repository root:

```bash
go run ./pqc/digitalsignature
go run ./pqc/sharedsecret
```

Both programs ran successfully with the go1.27.1 toolchain. You need a Go version whose standard library includes `crypto/mldsa` and `crypto/mlkem`. I didn't check the earliest version that has both.

## Notes

- `digitalsignature` passes empty `Options{}` both when signing and when verifying. If you enable a `Context`, you must use the same context on both sides.
