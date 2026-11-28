<p align="center">
  <img src="https://lh3.googleusercontent.com/d/1PyOdY6ynrZSK88chCXBWVUUpoVymLGTm" alt="go-perry Logo" width="220"/>
</p>

# go-perry

A collection of Golang experiments, and topic-based examples.
Each directory generally represents a technical topic, implementation exercise, or experiment focused on a specific Go mechanism.

## Topics

### 🧠 Runtime Internals

| Directory | Description |
| --- | --- |
| `assembly` | Experiments involving assembly, CPU registers, instructions, and memory models. |
| `bce` | Bounds Check Elimination experiments. |
| `build_tag` | Go build tags and conditional compilation examples. |
| `dump` | Program dump generation and debugging experiments. |
| `escape` | Escape analysis and memory allocation experiments. |
| `linkname-demo` | Experimental examples using `//go:linkname`. |
| `nosplit-demo` | Experiments with `//go:nosplit`. |
| `pin_thread` | Experiments involving goroutines and OS thread binding. |
| `reflect` | Go reflection experiments. |
| `runtime` | Go runtime concepts and runtime behavior experiments. |
| `unsafe` | `unsafe` package and low-level memory experiments. |
|`coroutine` | Asymmetric coroutines and single-thread register-swapping experiments via raw runtime primitives. |
| `swisstable` | Go 1.24 version's hashmap mechanism |
| `unique` | Experiments with the `unique` package for canonicalization (interning) and fast O(1) comparisons. |
| `weak_pointer` | Go 1.24 `weak` package: weak pointers that do not keep objects alive, combined with `runtime.AddCleanup` for a GC-aware cache. |
| `interface_values` | Interface values as (type, data) word pairs: `eface`/`iface` and itabs, boxing copies and allocations, type assertions and type switches, run-time panics when comparing uncomparable dynamic types, the typed-nil trap, and compile-time satisfaction checks. |
| `build_info` | Go build metadata experiments, including debug.ReadBuildInfo, VCS revision, vcs.modified, ldflags, reproducible builds, and build provenance. |

### 🧵 Concurrency Patterns

| Directory | Description |
| --- | --- |
| `async` | Asynchronous processing patterns, including producers, consumers, and background tasks. |
| `barrier` | Goroutine synchronization using barriers. |
| `channel` | Basic Go channel examples. |
| `channel_buffered` | Buffered channel examples. |
| `channel_unbuffered` | Unbuffered channel examples. |
| `context` | Cancellation, timeouts, and value propagation using Go context. |
| `errgroup` | Managing goroutines and collecting errors with `errgroup`. |
| `generator` | Generator patterns and goroutine-based generators. |
| `goroutine_confinement` | Goroutine confinement patterns for controlling access to shared data. |
| `goroutine_mutex` | Goroutine synchronization using mutexes. |
| `orDone&T` | Pipeline patterns such as `orDone` and `tee`. |
| `pipeline` | Pipeline-based concurrent data processing. |
| `pool` | Resource pool and worker pool patterns. |
| `singleflight` | Request deduplication using `singleflight`. |
| `sync_cond` | Condition variables with `sync.Cond`: `Wait`/`Signal`/`Broadcast` and a bounded queue built on two conds. |
| `sync_map` | `sync.Map` examples. |
| `sync_once` | `sync.Once` and Go 1.21 `OnceFunc`/`OnceValue`/`OnceValues`: one-time lazy init, its atomic fast path, and panic/error caching behavior. |
| `sync_atomic` | `sync/atomic` typed atomics: lock-free counters, CAS loops, `atomic.Pointer` copy-on-write config hot-swap, and Go 1.23 `Or`/`And` flags. |
| `time_timer` | `time.Timer`/`Ticker`/`AfterFunc` and the Go 1.23 timer semantics: synchronous timer channels, no stale values after `Stop`/`Reset`, dropped ticks. |
| `worker_pool` | Worker pool pattern examples. |
| `workerpool` | Additional worker pool implementations and experiments. |

### ⚡ Performance

| Directory | Description |
| --- | --- |
| `benchmark` | Go benchmarks, CPU profiling, and memory profiling experiments. |
| `cache` | Cache implementations and cached HTTP client examples. |
| `falseshare` | False sharing and CPU cache performance experiments. |
| `ratelimit` | Rate limiting implementations and performance-related examples. |
| `memcache` | Memcache usage and caching experiments. |
| `zero-copy` | Linux zero-copy mechanisms and kernel data paths |
| `directio` | Direct I/O experiments bypassing the OS page cache. |

### 🌐 Networking

| Directory | Description |
| --- | --- |
| `gorilla_websocket` | WebSocket server and client examples using gorilla/websocket. |
| `grpc` | gRPC client/server examples with Protocol Buffers. |
| `graphql` | GraphQL server examples. |
| `multipartload` | Multipart file upload client/server examples. |
| `http_servemux` | Go 1.22+ `net/http.ServeMux` routing: method patterns, `{id}`/`{path...}`/`{$}` wildcards, `PathValue`, most-specific-wins precedence, automatic 405 with `Allow`, redirects, registration-time conflict panics, and Go 1.23 `r.Pattern`. |

### 🔗 Distributed Systems

| Directory | Description |
| --- | --- |
| `breakcircut` | Circuit breaker pattern experiments for distributed services. |
| `kafka` | Kafka producer, consumer, and worker examples. |
| `rabbitmq` | RabbitMQ producer and consumer examples. |
| `redis_lock` | Redis distributed lock examples. |
| `p2p` | Peer-to-peer networking experiments. |
| `ipfs` | IPFS and decentralized storage experiments. |
| `otel` | OpenTelemetry distributed tracing across services, including TraceID, SpanID, parent-child relationships, and context propagation. |
| `logical_clock`| Including distributed common clock with lamport and vector and interval tree. |
| `detector` | Version Vector experiments, including causal history tracking, Dotted Version Vectors, and causal/concurrent event detection in distributed systems. |

### 💽 Storage

| Directory | Description |
| --- | --- |
| `lsm` | LSM Tree storage engine experiments. |
| `gorm` | GORM examples, including models, migrations, and database operations. |
| `memcache` | Memcache usage examples. |
| `cache` | Cache implementations and cached HTTP client examples. |

### 🧮 Data Structures

| Directory | Description |
| --- | --- |
| `alogrithm` | Algorithms and data structures, including search, sort, and structural topics. |
| `sort` | Sorting algorithms and sorting examples. |
| `boomfilter` | Bloom Filter implementation examples. |
| `iter` | Iterator patterns and iterator-style programming. |
| `generic_hash` | Generic hash table implementation and hash-based data structure experiments. |
| `slice_internals` | How slices work: header (ptr/len/cap), shared backing arrays, `append` aliasing and the `s[lo:hi:max]` / `slices.Clip` fix, `runtime.growslice` growth and size classes, Go 1.22 `slices.Delete` tail zeroing, and memory retention by sub-slices. |
| `string_utf8` | Strings as read-only bytes: `len` vs rune count, `range` decoding UTF-8, invalid bytes and `U+FFFD`, `[]byte`/`[]rune` conversion costs, rune-safe truncation, `strings.Builder` vs `+=`, and combining characters. |
| `regexp_re2` | Regular expressions with the RE2 engine: unanchored matching, the `Find(All)?(String)?(Submatch)?(Index)?` family, named groups, `${name}` replacements and the `$1x` pitfall, leftmost-first vs `Longest()`, flags, no backreferences/lookaround, linear-time matching, `QuoteMeta`, and concurrent use. |
| `generics_constraints` | Generics: type parameters and inference, constraints as type sets (`any`, `comparable`, `cmp.Ordered`, unions, `~`), the `S ~[]E` pattern, generic types whose methods cannot add type parameters, inference limits, `comparable` run-time panics with interfaces, and Go 1.26 self-referential constraints. |
| `map_semantics` | Map language semantics: `nil` map reads vs write panics, pointer-like sharing and `maps.Clone`, comma-ok lookups, randomized iteration order and `slices.Sorted(maps.Keys(m))`, non-addressable elements, mutation during `range`, NaN keys and `clear`, the set idiom, and the unrecoverable concurrent-write fatal error. |
| `container_heap` | Priority queues with `container/heap`: the `heap.Interface` contract (your `Push`/`Pop` only grow/shrink the slice), `Init`/`Push`/`Pop` costs, index tracking for `heap.Fix`/`heap.Remove`, the direct `h.Push` pitfall, unstable ordering and sequence-number tie-breaks, and a bounded top-K min-heap. |

### 📦 Serialization

| Directory | Description |
| --- | --- |
| `protobuf` | Protocol Buffers serialization and code generation examples. |
| `messagepack` | Binary serialization and deserialization using MessagePack. |
| `flatbuffer` | Schema-based binary serialization with direct buffer access. |
| `capProto` | Cap'n Proto schema-based binary serialization and generated Go code examples. |
| `cbor` | CBOR and Canonical CBOR serialization examples, including deterministic encoding and hashing. |
| `convert` | Type conversion and data format conversion examples. |
| `encoding_binary` | `encoding/binary`: big- vs little-endian `PutUint32`/`AppendUint32`, fixed-size struct `Write`/`Read` and why `int` is rejected, `binary.Size` vs `unsafe.Sizeof` padding, Go 1.23 `Encode`/`Decode`/`Append`, uvarint and zigzag varint sizes, `Uvarint` `n` error signalling, and capped length-prefixed framing. |

### 🔐 Security

| Directory | Description |
| --- | --- |
| `jwt` | JWT authentication and authorization examples. |
| `paseto` | PASETO authentication and security examples. |
| `schnorr` | Schnorr signatures and cryptography-related examples. |
| `cert` | X.509 certificate parsing and inspection examples. |
| `pqc` | Post-quantum cryptography examples using ML-KEM for shared-secret establishment and ML-DSA for digital signatures. |
| `os_root` | Go 1.24/1.25 `os.Root`: traversal-resistant file access that blocks `../` and symlink escapes, compared with naive `filepath.Join` and `filepath.IsLocal`. |

### ⚙️ System Programming

| Directory | Description |
| --- | --- |
| `file_tracker` | Tracking filesystem events through kernel notifications. |
| `ebpf` | eBPF examples for tracing process execution and filtering network packets at the kernel level. |
| `iouring` | Linux `io_uring` experiments using `liburing`, including asynchronous TCP I/O. |

### 🏗️ Application Patterns

| Directory | Description |
| --- | --- |
| `decorator` | Decorator design pattern examples. |
| `dependency` | Dependency injection and dependency design examples. |
| `background_faktory` | Distributed background job examples using Faktory. |
| `background_gocraft` | Distributed background job examples using gocraft/work. |
| `gocron` | Scheduled job examples using gocron. |
| `logger` | Logging and file output examples. |
| `rss` | RSS generation and parsing examples. |
| `test` | Go testing examples and experiments. |
| `embed-demo` | Examples using Go's `embed` feature. |
| `error_wrapping` | Error wrapping with `%w`, `errors.Is`/`As`, Go 1.26 `errors.AsType`, custom `Is`/`Unwrap`, multi-error trees via `errors.Join`, and the typed-nil pitfall. |
| `defer_panic_recover` | `defer` order and argument evaluation, `panic`/`recover` rules, converting panics to errors with named results, `panic(nil)` (`*runtime.PanicNilError`), `runtime.Goexit`, and recovering inside goroutines. |
| `io_reader_writer` | `io.Reader`/`io.Writer` contracts (data returned with `io.EOF`), custom readers and writers, `TeeReader`/`MultiWriter`/`LimitReader`, the `WriterTo` fast path in `io.Copy`, and `io.Pipe` with `CloseWithError`. |
| `bufio_scanner` | Buffered I/O with `bufio`: `Scanner` lines/words (CRLF, missing final newline), custom `SplitFunc` with `ErrFinalToken`, the 64 KiB `ErrTooLong` limit and `Scanner.Buffer`, `Bytes()` aliasing, `Reader.Peek`/`ReadString`, and `Writer` buffering and `Flush`. |
| `slog_structured` | Structured logging with `log/slog`: Text/JSON handlers, typed attrs, `With`/`WithGroup`, `LevelVar` dynamic levels, lazy and redacted values via `LogValuer`, a context-aware custom `Handler`, and Go 1.26 `NewMultiHandler`. |
| `json_encoding` | `encoding/json` struct tags, Go 1.24 `omitzero` vs `omitempty`, custom `MarshalJSON`/`UnmarshalJSON`, `json.RawMessage` tagged unions, `DisallowUnknownFields`, and large-number precision with `UseNumber`. |
| `text_template` | `text/template` and `html/template`: actions, whitespace trimming, pipelines, `FuncMap`, `define`/`block` with `Clone`, `missingkey=error`, partial output on errors, and context-aware auto-escaping. |
| `struct_embedding` | Struct embedding as composition: field/method promotion, method sets when embedding `T` vs `*T`, shadowing is not overriding, the depth rule and ambiguous selectors, embedding an interface as a decorator, the `sync.Mutex` API leak, and JSON flattening. |
| `closures_loopvar` | Closures capture variables, not values: heap-allocated closure state, capture-by-value optimization, Go 1.22 per-iteration loop variables and 3-clause copy semantics, goroutines in loops, method values and `defer` evaluating early, and the `go.mod` version gate. |
| `time_values` | `time.Time` values: reference-time layouts and silent layout mistakes, `Parse` vs `ParseInLocation`, the monotonic clock reading and `Round(0)`, `Equal` vs `==` and map keys, `AddDate` normalization, DST-aware `Add` vs `AddDate`, and embedding `time/tzdata` for Windows. |

The topic list is continuously evolving. New topics are welcome.

## Contributing

Contributions are welcome.

You are not limited to improving existing examples. You are encouraged to create your own concept, experiment, or technical topic.

The recommended workflow is to fork then create a feature branch:

```bash
git checkout -b feature/<your-concept>
```

For example:

```bash
git checkout -b feature/tracing
```

After implementing and documenting your concept, open a Pull Request.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the contribution workflow and guidelines.

## Authors

Contributors can add their name and ID to [AUTHORS.md](AUTHORS.md).

The author list is intended to give credit to everyone who contributes to the project.

## Philosophy

The goal of `go-perry` is not simply to collect finished code.

It is a place to:

- Experiment
- Build
- Document
- Share

A small idea can become a new topic.

**Create your own concept. Implement it. Document it. Share it.**