// Package main demonstrates structured logging with log/slog:
// handlers, attributes, groups, dynamic levels, lazy values (LogValuer),
// a context-aware custom Handler, and Go 1.26 slog.NewMultiHandler.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// noTime removes the time attribute so the output is deterministic.
// ReplaceAttr is called for every attribute; groups is the current group path.
func noTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{} // an empty Attr is dropped
	}
	return a
}

func textLogger(level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: noTime,
	}))
}

func jsonLogger(level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: noTime,
	}))
}

// ---------- Part 3 / 4: LogValuer ----------

// Token is a secret. Its LogValue hides the value whenever it is logged.
type Token string

func (Token) LogValue() slog.Value { return slog.StringValue("REDACTED") }

// expensive counts how often it is actually resolved.
type expensive struct{ calls *int }

func (e expensive) LogValue() slog.Value {
	*e.calls++
	return slog.GroupValue(slog.Int("rows", 42), slog.String("table", "users"))
}

// ---------- Part 5: a custom Handler ----------

type ctxKey struct{}

// requestIDHandler wraps another Handler and adds request_id from the context.
type requestIDHandler struct{ next slog.Handler }

func (h requestIDHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.next.Enabled(ctx, l)
}

func (h requestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		r = r.Clone() // do not mutate a Record we might share
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.next.Handle(ctx, r)
}

// WithAttrs / WithGroup must return a wrapped handler too,
// otherwise logger.With(...) would silently drop our behavior.
func (h requestIDHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return requestIDHandler{h.next.WithAttrs(as)}
}

func (h requestIDHandler) WithGroup(name string) slog.Handler {
	return requestIDHandler{h.next.WithGroup(name)}
}

func main() {
	fmt.Println("== Part 1: key-value pairs vs typed Attrs ==")
	log := textLogger(slog.LevelInfo)
	log.Info("user login", "user", "perry", "attempt", 3)
	log.Info("user login", slog.String("user", "perry"), slog.Int("attempt", 3))
	// A literal log.Info("odd", "user") is caught by go vet; args built at
	// runtime are not, and the dangling key becomes !BADKEY.
	odd := []any{"user"}
	log.Info("odd args", odd...)
	log.Debug("not printed: below Info")

	fmt.Println("== Part 2: With and WithGroup (JSON) ==")
	base := jsonLogger(slog.LevelInfo)
	reqLog := base.With("service", "api").WithGroup("req")
	reqLog.Info("handled", "method", "GET", "path", "/users", "status", 200)
	base.Info("nested", slog.Group("db", slog.String("op", "select"), slog.Int("ms", 7)))

	fmt.Println("== Part 3: dynamic level and lazy values ==")
	var level slog.LevelVar // zero value is Info
	dyn := textLogger(&level)
	calls := 0
	dyn.Debug("query", "detail", expensive{&calls})
	fmt.Println("LogValue calls while level=INFO: ", calls)
	level.Set(slog.LevelDebug)
	dyn.Debug("query", "detail", expensive{&calls})
	fmt.Println("LogValue calls while level=DEBUG:", calls)
	fmt.Println("Enabled(Debug) =", dyn.Enabled(context.Background(), slog.LevelDebug))

	fmt.Println("== Part 4: LogValuer for redaction ==")
	log.Info("auth", "user", "perry", "token", Token("s3cr3t-abc"))

	fmt.Println("== Part 5: custom Handler reading the context ==")
	h := requestIDHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: noTime})}
	clog := slog.New(h).With("component", "billing")
	ctx := context.WithValue(context.Background(), ctxKey{}, "req-7f3a")
	clog.InfoContext(ctx, "charge created", "amount", 1200)
	clog.Info("no context, no request_id")

	fmt.Println("== Part 6: Go 1.26 NewMultiHandler ==")
	multi := slog.New(slog.NewMultiHandler(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: noTime}),
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: noTime}),
	))
	multi.Info("only text handler")
	multi.Warn("both handlers", "disk_free_pct", 4)
}
