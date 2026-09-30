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
| `sync_map` | `sync.Map` examples. |
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
| `direct-io` | Direct I/O experiments bypassing the OS page cache. |

### 🌐 Networking

| Directory | Description |
| --- | --- |
| `gorilla_websocket` | WebSocket server and client examples using gorilla/websocket. |
| `grpc` | gRPC client/server examples with Protocol Buffers. |
| `graphql` | GraphQL server examples. |
| `multipartload` | Multipart file upload client/server examples. |

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
| `version-vector` | Version Vector experiments, including causal history tracking, Dotted Version Vectors, and causal/concurrent event detection in distributed systems. |

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

### 📦 Serialization

| Directory | Description |
| --- | --- |
| `protobuf` | Protocol Buffers serialization and code generation examples. |
| `messagepack` | Binary serialization and deserialization using MessagePack. |
| `flatbuffer` | Schema-based binary serialization with direct buffer access. |
| `capProto` | Cap'n Proto schema-based binary serialization and generated Go code examples. |
| `cbor` | CBOR and Canonical CBOR serialization examples, including deterministic encoding and hashing. |
| `convert` | Type conversion and data format conversion examples. |

### 🔐 Security

| Directory | Description |
| --- | --- |
| `jwt` | JWT authentication and authorization examples. |
| `paseto` | PASETO authentication and security examples. |
| `schnorr` | Schnorr signatures and cryptography-related examples. |
| `x509` | X.509 certificate parsing and inspection examples. |
| `pqc` | Post-quantum cryptography examples using ML-KEM for shared-secret establishment and ML-DSA for digital signatures. |

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