# boomfilter

A Bloom filter (the type is spelled `BoomFilter` in the code): a probabilistic set that can say "definitely not present" or "probably present" using a fixed-size bit array.

## What the code does

`main.go`:

- `optimalM(n, p)` computes the bit-array size `m = ceil(-(n * ln p) / (ln 2)^2)`.
- `optimalK(m, n)` computes the hash count `k = ceil((m / n) * ln 2)`.
- `NewBloom(n, p)` allocates `(m+7)/8` bytes. For `n = 10000` and `p = 0.01` this gives `m = 95851` bits (about 12 KB) and `k = 7`.
- `hash(item, seed)` derives the i-th hash by running FNV-1a 64-bit over the item followed by the 8-byte little-endian seed `i`.
- `Add` sets bit `hash(item, i) % m` for each `i` in `0..k-1`. `Exists` returns `false` as soon as one of those bits is unset.
- `main` adds `alice` and `bob`, then checks `alice`, `bob`, and `carol`.

Actual output:

```
true
true
false
```

With these fixed inputs the output is deterministic (FNV has no randomness), so `carol` always prints `false`. A different non-member could still hit all k bits (a false positive), but with two items in this filter that is very unlikely.

## Concepts you can learn

- **No false negatives, possible false positives**: an added item always has all k bits set. A non-member can also find all k bits set by coincidence. This is why a Bloom filter suits "skip the expensive lookup if definitely absent" checks.
- **Sizing from n and p**: the formulas trade memory for accuracy. A lower target false-positive rate needs more bits per element (about 9.6 bits per element for 1%).
- **Choosing k**: too few hashes leaves too little signal, and too many fill the array faster. `(m/n) ln 2` minimizes the false-positive rate for a given m and n.
- **Bit packing**: `pos/8` selects the byte and `pos%8` the bit, so the filter uses m bits rather than m bytes.
- **Deriving k hashes from one hash function** by seeding it differently. This is simple, but each Add or Exists call runs k full FNV passes over the item.

## How to run

Own `go.mod` (module `boom`, standard library only):

```bash
cd boomfilter
go run .
```

## Notes

- The comment in `main` says "100000" items, but the code passes `10000`.
- The filter does not support removal (clearing a bit could affect other items) and does not track how many items were added, so there is no guard against overfilling it past `n`.
- The directory and type names use "boom", not "bloom". The directory name is the real one, so do not rename it.
