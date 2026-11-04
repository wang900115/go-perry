# kafka

A coffee-order pipeline: an HTTP service publishes orders to Kafka and a worker consumes them, both using `github.com/IBM/sarama`.

## What the code does

There are two independent modules, each with its own `go.mod`:

- **`producer/`** (module `producer`): HTTP server on `:3000` with one route, `POST /order`. It decodes `{"customer_name": "...", "coffee_type": "..."}`, re-marshals it and sends it to topic `coffee_orders` on `localhost:9092` with a **sync producer** (`RequiredAcks = WaitForAll`, `Retry.Max = 5`, `Return.Successes = true`). It logs `Order is stored in topic(coffee_orders)/partition(N)/offset(M)` and replies `{"msg":"Order for <name> placed successfully!","success":true}`. Non-POST requests get `405`, bad JSON `400`, Kafka errors `500`.
- **`worker/`** (module `worker`): connects to `localhost:9092`, consumes **partition 0** of `coffee_orders` from `OffsetOldest`, and for each message prints `Received order count N: | Topic(coffee_orders) | Message(...)` and `Brewing coffee for order: <json>`. On SIGINT/SIGTERM it prints `Interrupt is detected`, then `Processed N messages`, and closes the consumer.

## Concepts you can learn

- **Sync producer and acks**: `WaitForAll` waits for all in-sync replicas, trading latency for durability; `SendMessage` returns the partition and offset the message landed on.
- **Partition consumer vs. consumer group**: the worker reads one partition directly with `ConsumePartition`; there is no group, so no committed offsets and no rebalancing.
- **Offsets**: `OffsetOldest` replays the whole partition every time the worker starts.
- **Decoupling via a queue**: the HTTP request returns once Kafka has the order; brewing happens asynchronously in the worker.
- **Graceful shutdown with signals**: `signal.Notify` + a done channel to stop the consume loop.

## How to run

Prerequisite: a Kafka broker reachable at `localhost:9092`. The topic `coffee_orders` must exist, or the broker must be configured to auto-create topics.

```bash
cd kafka/worker && go run .        # terminal 1
cd kafka/producer && go run .      # terminal 2

curl -X POST http://localhost:3000/order \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"Alice","coffee_type":"Latte"}'
```

## Notes

- The producer opens and closes a new Kafka connection **per HTTP request** (`PushOrderToQueue`), which is simple but expensive.
- The worker only reads partition 0. If the topic has more partitions, orders on other partitions are never seen.
- `worker/go.mod` lists sarama as `// indirect` even though it is imported directly; it still builds (`go mod tidy` would fix the marker).
