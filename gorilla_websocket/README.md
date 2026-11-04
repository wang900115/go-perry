# gorilla_websocket

A WebSocket echo server built with `github.com/gorilla/websocket`, plus a plain HTML page to talk to it.

## What the code does

- **`main.go`**: serves `/ws` on `:8080` and prints `WebSocket server is running on :8080/ws`. Each request is upgraded with an `Upgrader` (1 KB read/write buffers, `CheckOrigin` always `true`). For every message it prints `Received message: <data>` and writes the same bytes back with the same message type (text or binary). Read/write errors (including client disconnects) are printed and end that connection's loop.
- **`client.html`**: opens `ws://localhost:8080/ws`, sends the input box text on button click, and appends `Received: <data>` to the page for each echoed message.

## Concepts you can learn

- **HTTP upgrade**: a WebSocket starts as an HTTP request and is switched to a persistent, bidirectional connection by `Upgrader.Upgrade`.
- **Per-connection read loop**: `net/http` runs each handler in its own goroutine, so one blocking `ReadMessage` loop per client is the natural pattern.
- **Message types**: echoing with the received `messageType` preserves text vs. binary frames.
- **Origin checking**: `CheckOrigin` is the CSRF-style guard for browsers; returning `true` disables it.

## How to run

```bash
cd gorilla_websocket
go run .
```

Then open `client.html` in a browser (as a local file is fine) and send messages.

## Notes

- `CheckOrigin: true` accepts connections from any website; fine for a demo, not for production.
- The client inserts echoed text with `innerHTML`, so HTML in a message is rendered, not escaped.
- `go.mod` marks gorilla/websocket as `// indirect`; it builds anyway.
