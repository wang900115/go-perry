// Package main demonstrates the enhanced routing patterns of net/http.ServeMux
// (Go 1.22+): methods, wildcards, {$}, {name...}, precedence, 405 responses,
// redirects, and registration-time conflict detection.
//
// Every request is served in memory with httptest.NewRecorder, so no port is opened.
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// do sends one in-memory request through h and prints the result.
func do(h http.Handler, method, target string) {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body := strings.TrimSpace(rec.Body.String())
	extra := ""
	if loc := rec.Header().Get("Location"); loc != "" {
		extra = " Location=" + loc
	}
	if allow := rec.Header().Get("Allow"); allow != "" {
		extra += " Allow=" + allow
	}
	if method == http.MethodHead {
		body = "(no body for HEAD)"
	}
	fmt.Printf("%-6s %-28s -> %d %s%s\n", method, target, rec.Code, body, extra)
}

// text returns a handler that writes a label plus the matched pattern.
// r.Pattern was added in Go 1.23.
func text(label string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s [pattern %q]", label, r.Pattern)
	}
}

func part1() {
	fmt.Println("== Part 1: methods, wildcards, PathValue ==")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "get user id=%s", r.PathValue("id"))
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "create user")
	})
	// {path...} must be last and matches the rest of the path, slashes included.
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "file path=%q", r.PathValue("path"))
	})
	// {$} matches only the exact path "/", not every path.
	mux.HandleFunc("GET /{$}", text("home"))

	do(mux, "GET", "/users/42")
	do(mux, "POST", "/users")
	do(mux, "GET", "/files/docs/2026/notes.txt")
	do(mux, "GET", "/files/")
	do(mux, "GET", "/")
	do(mux, "GET", "/nope")
	// A wildcard matches one whole segment; it never matches an empty segment
	// or more than one segment.
	do(mux, "GET", "/users/42/posts")
	// Wildcard values are unescaped: %2F becomes "/", but routing still treats
	// it as one segment.
	do(mux, "GET", "/users/a%2Fb")
}

func part2() {
	fmt.Println("== Part 2: precedence, 405, HEAD ==")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", text("user by id"))
	mux.HandleFunc("GET /users/me", text("current user")) // more specific: wins for /users/me
	mux.HandleFunc("/users/{id}", text("any method"))     // less specific: no method
	mux.HandleFunc("GET /orders/{id}", text("order"))
	mux.HandleFunc("DELETE /orders/{id}", text("delete order"))

	do(mux, "GET", "/users/me")
	do(mux, "GET", "/users/7")
	do(mux, "PUT", "/users/7")   // falls back to the method-less pattern
	do(mux, "PUT", "/orders/1")  // path matches, method does not -> 405 + Allow
	do(mux, "HEAD", "/orders/1") // "GET" patterns also match HEAD
}

func part3() {
	fmt.Println("== Part 3: trailing slash and path cleaning ==")
	mux := http.NewServeMux()
	// A pattern ending in "/" matches the whole subtree.
	mux.HandleFunc("GET /static/", text("static subtree"))
	mux.HandleFunc("GET /docs/{$}", text("docs index only"))

	do(mux, "GET", "/static/css/site.css")
	do(mux, "GET", "/static")                // redirected to the subtree root
	do(mux, "GET", "/docs/")                 // {$}: only exact "/docs/"
	do(mux, "GET", "/docs/intro")            // not matched by /docs/{$}
	do(mux, "GET", "/static/../static/a.js") // unclean path -> redirect to clean path
}

func part4() {
	fmt.Println("== Part 4: conflicts, mux.Handler, SetPathValue ==")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts/{id}", text("post"))

	// Neither pattern is more specific than the other, and both match
	// "/posts/latest", so registration panics.
	func() {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprint(r)
				// The full message also names the file:line of both
				// registrations; cut it to keep the output machine-independent.
				if i := strings.Index(msg, " (registered"); i >= 0 {
					msg = msg[:i] + " conflicts with \"GET /posts/{id}\" ..."
				}
				fmt.Println("panic:", msg)
			}
		}()
		mux.HandleFunc("GET /{kind}/latest", text("latest"))
	}()

	// mux.Handler reports which pattern would handle a request, without serving it.
	for _, target := range []string{"/posts/9", "/missing"} {
		_, pattern := mux.Handler(httptest.NewRequest("GET", target, nil))
		fmt.Printf("Handler(%s) -> pattern %q\n", target, pattern)
	}

	// SetPathValue lets a unit test call a handler directly, without a mux.
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "post id=%s", r.PathValue("id"))
	}
	req := httptest.NewRequest("GET", "/whatever", nil)
	req.SetPathValue("id", "123")
	rec := httptest.NewRecorder()
	h(rec, req)
	fmt.Println("direct call:", rec.Body.String())
	// PathValue returns "" for an unknown name; it does not panic.
	fmt.Printf("unknown name: %q\n", req.PathValue("nope"))
}

func main() {
	part1()
	part2()
	part3()
	part4()
}
