# rabbitmq

A Gin HTTP endpoint publishes text messages to a RabbitMQ queue and a separate consumer prints them, using `github.com/rabbitmq/amqp091-go`.

## What the code does

One module (`rabbitmq/go.mod`) with two programs; both use queue `ServicelQueue` on `amqp://localhost:5672`.

- **`producer/main.go`**: connects, opens a channel, declares `ServicelQueue` as **durable**, then starts Gin on `:8080`. `GET /send?msg=<text>` publishes `<text>` (`text/plain`) to the default exchange with routing key `ServicelQueue` and returns `{"message":"<text>","status":"success"}`. Missing `msg` returns `400 {"error":"Mesasage is required"}`.
- **`consumer/main.go`**: connects, opens a channel and consumes `ServicelQueue` with **auto-ack**, logging `Message: <body>` for each delivery. On SIGINT/SIGTERM it logs `Interrupt detected` and exits.

## Concepts you can learn

- **Connection vs. channel**: one TCP connection, lightweight channels on top; all operations go through a channel.
- **Queue declaration and durability**: `QueueDeclare(..., durable=true, ...)` makes the queue survive a broker restart (messages themselves are not marked persistent here).
- **Default exchange**: publishing with exchange `""` and the queue name as routing key delivers directly to that queue.
- **Auto-ack**: messages are acknowledged on delivery, so a consumer crash after receipt loses them; manual ack is the alternative.

## How to run

Prerequisite: RabbitMQ listening on `localhost:5672`. The URL contains no credentials, so the client library's defaults are used.

```bash
cd rabbitmq
go run ./producer        # terminal 1 - start first, it declares the queue
go run ./consumer        # terminal 2

curl "http://localhost:8080/send?msg=hello"
```

## Notes

- Only the producer declares the queue. Starting the consumer first on a fresh broker fails because `ServicelQueue` does not exist yet.
- The `/send` handler assigns to the outer `err` variable from `main`, so concurrent requests write to a shared variable (a data race).
- If the broker closes the connection, the `messages` channel is closed and the consumer's `select` keeps receiving zero-value deliveries, logging empty messages in a loop.
- `go.mod` marks all dependencies `// indirect`; it builds anyway.
