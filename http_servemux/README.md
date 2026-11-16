# http_servemux

Enhanced routing patterns in `net/http.ServeMux` (Go 1.22+): method matching
(`GET /users/{id}`), single-segment wildcards with `r.PathValue`, `{name...}`
rest wildcards, `{$}` exact matches, "most specific pattern wins" precedence,
automatic `405 Method Not Allowed` with an `Allow` header, `GET` also matching
`HEAD`, trailing-slash and path-cleaning redirects, registration-time conflict
panics, `mux.Handler`, Go 1.23 `r.Pattern`, and `SetPathValue` for unit tests.

```bash
go run .
```

Standard library only. Every request is served in memory with `httptest.NewRecorder`,
so no port is opened and no network is needed. This module declares `go 1.23`
(`Request.Pattern` needs Go 1.23; the new patterns need a `go` directive of 1.22 or later).
Verified on go1.26 linux/amd64 (2 runs, identical output); `go vet .` and `GOOS=windows go vet .` pass.
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: methods, wildcards, PathValue ==
GET    /users/42                    -> 200 get user id=42
POST   /users                       -> 200 create user
GET    /files/docs/2026/notes.txt   -> 200 file path="docs/2026/notes.txt"
GET    /files/                      -> 200 file path=""
GET    /                            -> 200 home [pattern "GET /{$}"]
GET    /nope                        -> 404 404 page not found
GET    /users/42/posts              -> 404 404 page not found
GET    /users/a%2Fb                 -> 200 get user id=a/b
== Part 2: precedence, 405, HEAD ==
GET    /users/me                    -> 200 current user [pattern "GET /users/me"]
GET    /users/7                     -> 200 user by id [pattern "GET /users/{id}"]
PUT    /users/7                     -> 200 any method [pattern "/users/{id}"]
PUT    /orders/1                    -> 405 Method Not Allowed Allow=DELETE, GET, HEAD
HEAD   /orders/1                    -> 200 (no body for HEAD)
== Part 3: trailing slash and path cleaning ==
GET    /static/css/site.css         -> 200 static subtree [pattern "GET /static/"]
GET    /static                      -> 301 <a href="/static/">Moved Permanently</a>. Location=/static/
GET    /docs/                       -> 200 docs index only [pattern "GET /docs/{$}"]
GET    /docs/intro                  -> 404 404 page not found
GET    /static/../static/a.js       -> 301 <a href="/static/a.js">Moved Permanently</a>. Location=/static/a.js
== Part 4: conflicts, mux.Handler, SetPathValue ==
panic: pattern "GET /{kind}/latest" conflicts with "GET /posts/{id}" ...
Handler(/posts/9) -> pattern "GET /posts/{id}"
Handler(/missing) -> pattern ""
direct call: post id=123
unknown name: ""
```
