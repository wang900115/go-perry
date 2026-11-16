# 2026-11-16 · net/http ServeMux routing patterns (Go 1.22+)

## Why learn this
- Almost every Go backend has an HTTP router. Since Go 1.22 the standard `http.ServeMux` supports methods and path wildcards (`GET /users/{id}`), so many services no longer need a third-party router such as gorilla/mux or chi just for routing.
- You will read these patterns in REST APIs, health/metrics endpoints, admin tools, and the servers in this repo (`ratelimit`, `otel`, `multipartload`).
- The rules (precedence, 405, redirects, conflicts) decide which handler runs. Getting them wrong gives silent 404s or the wrong handler, not compile errors.
- It builds on earlier lessons: `context` (every `*http.Request` carries one), `io_reader_writer` (`ResponseWriter` is an `io.Writer`), `slog_structured` and `json_encoding` (what handlers usually do), `defer_panic_recover` (registration conflicts panic).

## Core concepts
1. **Pattern syntax**: `[METHOD ][HOST]/[PATH]`, for example `GET /users/{id}`, `POST /users`, `api.example.com/`. Exactly one space between method and path (https://pkg.go.dev/net/http#hdr-Patterns-ServeMux).
2. **Methods**: a pattern with a method only matches that method. `GET` also matches `HEAD`. A pattern without a method matches every method.
3. **Wildcards**: `{name}` matches exactly one non-empty path segment. `r.PathValue("name")` returns its value (unescaped, so `%2F` becomes `/` but still counts as one segment). `PathValue` of an unknown name returns `""`.
4. **Rest wildcard**: `{name...}` must be the last segment and matches the rest of the path, slashes included, possibly empty (`/files/` gives `""`).
5. **Trailing slash**: a pattern ending in `/` matches the whole subtree (`/static/` matches `/static/css/site.css`). `/` therefore matches everything. `{$}` anchors the end: `/{$}` matches only `/`, `/docs/{$}` only `/docs/`.
6. **Precedence: most specific wins.** If two patterns match a request, the one that matches a strict subset of requests wins: `GET /users/me` beats `GET /users/{id}`, which beats `/users/{id}`. A host pattern beats a host-less one.
7. **Conflicts panic at registration.** If two patterns overlap and neither is more specific (`GET /posts/{id}` vs `GET /{kind}/latest`, both match `/posts/latest`), `Handle`/`HandleFunc` panics. This happens at startup, not at request time.
8. **405 and Allow**: if a path matches some pattern but the method does not, the mux replies `405 Method Not Allowed` with an `Allow` header listing the methods that would match.
9. **Redirects**: a request for `/static` when only `/static/` is registered gets a `301` to `/static/`. Unclean paths (`..`, `.`, double slashes) get a `301` to the cleaned path.
10. **Inspection helpers**: `mux.Handler(r)` returns the handler and pattern without serving. `r.Pattern` (Go 1.23, https://go.dev/doc/go1.23#nethttppkgnethttp) holds the matched pattern inside the handler. `r.SetPathValue` lets tests call a handler directly with path values.

Background: https://go.dev/blog/routing-enhancements and https://go.dev/doc/go1.22#enhanced_routing_patterns.

### How it works (standard library)
Based on the go1.26 source (`net/http/server.go`, `pattern.go`, `routing_tree.go`, `routing_index.go`); internal details may change.
- `Handle` parses the string into a `pattern` (method, host, list of segments). Then it checks the new pattern against existing ones. A `routingIndex` narrows that check to patterns that could possibly conflict, so registration does not compare against every pattern. If two patterns overlap without one being more specific, it panics with a message naming both patterns and where they were registered.
- Patterns are stored in a tree of `routingNode`s: first by host, then by method, then segment by segment. Literal segments are looked up before wildcards, which is how "more specific" wins at match time.
- `ServeHTTP` → `findHandler`: it cleans the path, tries to match, and if needed builds a redirect (missing trailing slash, unclean path). If nothing matches, it collects the methods that would match the path; if there are any, it returns a 405 handler with `Allow`, otherwise `NotFound`.
- On a match, the wildcard values are stored on the request, so `PathValue` is a lookup over the matched pattern's wildcard names, not a re-parse.
- The old (Go 1.21) mux is still in the source as `serveMux121`. It is used when `GODEBUG=httpmuxgo121=1`, and that is the default when the main module's `go` directive is older than 1.22.

## Full example
See `main.go`. Verified with `go run .` (2 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. Method patterns, `{id}`, `{path...}`, `/{$}`, 404 for extra segments, `%2F` in a wildcard value.
2. Precedence (`/users/me` vs `/users/{id}` vs method-less), 405 with `Allow: DELETE, GET, HEAD`, `HEAD` handled by a `GET` pattern.
3. Subtree patterns, the `/static` → `/static/` redirect, `/docs/{$}`, cleaning `/static/../static/a.js`.
4. A conflicting registration recovered from a panic, `mux.Handler`, `SetPathValue` in a direct handler call.

## Common pitfalls
1. **Old `go` directive = old mux, silently.** If `go.mod` says `go 1.21` (or older), `GET /users/{id}` is not a method pattern. It is registered without error and every request returns 404. Check the `go` line in `go.mod`, and check that `GODEBUG=httpmuxgo121=1` is not set.
2. **`/` catches everything.** `mux.HandleFunc("/", home)` is a subtree pattern, so every unknown path goes to `home` with status 200 instead of 404. Use `GET /{$}` for "only the root".
3. **Expecting first-registered-wins.** Order does not matter; specificity does. Two overlapping patterns with no clear winner panic at startup. If you build routes from config or plugins, a conflict can crash the program on boot.

## Exercise (no answer given)
Build an in-memory todo API in a new `main.go` (or as tests with `httptest`):
- `GET /todos`, `POST /todos`, `GET /todos/{id}`, `PUT /todos/{id}`, `DELETE /todos/{id}`, and `GET /todos/{id}/files/{name...}`.
- Store todos in a map protected by a `sync.Mutex`; return JSON.
- Write checks for: a `PATCH /todos/1` request (what status and `Allow` header do you get?), `GET /todos/abc` (validate that `id` is a number and return 400), and `GET /todos/` (which pattern, if any, matches?).
- Add a logging middleware that wraps the whole mux and logs `r.Pattern`. Explain why the value is empty in the middleware and find a way to log the matched pattern anyway.
- Try to register `GET /todos/{id}` and `GET /{resource}/{id}` together. Does it panic? Then try `GET /todos/{id}` and `GET /{resource}/export`. Explain the difference.

## Suggested directory
`http_servemux`
