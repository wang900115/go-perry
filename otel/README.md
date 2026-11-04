# otel

OpenTelemetry tracing across three Go services (HTTP -> gRPC -> gRPC). As written, the trace context is **not** propagated between processes (see Notes), so each service starts its own trace.

## What the code does

One module (`otel/go.mod`) with three programs and two shared packages:

- **`telemetry/telemetry.go`**: `Init(serviceName)` builds a `TracerProvider` with a `service.name` resource and a **batching** span processor. The exporter is a hand-written `stdoutExporter` that prints one line per span:
  `[TRACE] service=<name> span=<span name> trace_id=<id> span_id=<id> parent=<parent span id>`.
  No collector or Jaeger is involved.
- **`main.go`** (`main-service`): HTTP server on `:8080`, wrapped in `otelhttp.NewHandler`. `/orders` decodes `{"user_id","sku","quantity"}`, calls `OrderService.CreateOrder` on `localhost:50051`, then `InventoryService.ReserveStock` on `localhost:50052` (both gRPC clients use the `otelgrpc` client stats handler), and returns `{"order_id","order_status","inventory_reserved","remaining_stock"}`.
- **`order/main.go`** (`order-service`, `:50051`): validates the request (`InvalidArgument` on empty user/sku or quantity <= 0), starts a manual span `order.create`, logs the order and always returns `order_id: "order-001", status: "created"`.
- **`inventory/main.go`** (`inventory-service`, `:50052`): in-memory stock `iphone: 10, macbook: 5` behind a mutex. Starts span `inventory.reserve`; unknown SKU -> `NotFound`; insufficient stock -> `success: false` with the current amount; otherwise decrements and returns the remainder.
- **`proto/`**: `order.proto`, `inventory.proto` and the generated `*.pb.go` / `*_grpc.pb.go` (protoc-gen-go v1.36.12). Do not edit the generated files.

For one request, spans inside one process share a `trace_id` and are linked by `parent` (e.g. main-service: `POST /orders` -> two gRPC client spans; order-service: server span -> `order.create`). Each service prints a **different** `trace_id`, and the gRPC server spans have `parent=0000000000000000`, because no propagator is configured.

## Concepts you can learn

- **Trace ID vs. span ID vs. parent**: one trace per request, one span per unit of work, and parent IDs to rebuild the call tree.
- **Context propagation across processes**: `otelhttp`/`otelgrpc` inject and extract trace context using the global `TextMapPropagator`, which is a no-op until you call `otel.SetTextMapPropagator(...)`. This code omits that call, which is exactly why traces break at process boundaries.
- **Auto vs. manual instrumentation**: middleware/stats handlers create spans for you; `tracer.Start(ctx, ...)` adds business-level child spans.
- **Resources**: `service.name` identifies which process emitted a span.
- **Custom `SpanExporter`**: implementing `ExportSpans`/`Shutdown` shows what an exporter actually receives.
- **Batch span processor**: spans are buffered and exported in batches, not one by one.

## How to run

No external services are required. Start each in its own terminal from the `otel` directory:

```bash
cd otel
go run ./order        # :50051
go run ./inventory    # :50052
go run .              # :8080

curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":"u1","sku":"iphone","quantity":2}'
```

Regenerating `proto/` requires `protoc` with `protoc-gen-go` and `protoc-gen-go-grpc`; the exact command used is not recorded in the repo.

## Notes

- `[TRACE]` lines appear with a delay because of the batcher.
- All three `main`s exit via `log.Fatal` or a kill signal, which skips the deferred `shutdown`, so the last unflushed batch of spans can be lost.
- `/orders` does not check the HTTP method, and the order is "created" before stock is reserved; a failed reservation does not undo it.
- No propagator is set, so the three services do not share a trace (verified: one `curl` produced three different `trace_id`s). Adding `otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))` in `telemetry.Init` would make them join one trace.
- `telemetry.Tracer` is defined but unused (services call `otel.Tracer` directly).
