# memcache

Read-through caching of a MySQL-backed blog API with Memcached, based on "Episode #45: Caching in Go: How to Integrate Memcached" ([episode link](https://www.codeheim.io/courses/Episode-45-Caching-with-Memcached-and-Golang-668d72cd6aad903c4582b419)).

## What the code does

Module `go_memcached`, using `net/http`, GORM (MySQL driver) and `github.com/bradfitz/gomemcache`.

- **`main.go`**: connects to MySQL, runs `AutoMigrate` for `Blog`, and serves `GET /blogs/{id}` on `:8080`.
- **`models/setup.go`**: opens MySQL with DSN `net_http_blogs:tmp_pwd@tcp(127.0.0.1:3306)/net_http_blogs?charset=utf8&parseTime=true`; panics if it cannot connect.
- **`models/blog.go`**: `Blog` = `gorm.Model` + `Title` (varchar 255) + `Content` (text). `BlogsFind(id)` loads one row; `BlogsAll()` exists but is not used by any route.
- **`models/memcache.go`**: client for `localhost:11211`. `CacheData(key, ttl, fn)` returns the cached bytes on a hit; on a miss (or any `Get` error) it calls `fn` via reflection (it must return `[]byte`), stores the result with the TTL, and returns it.
- **`controllers/blogs_controller.go`**: parses the ID from the path (`400 Invalid blog ID` if not a number) and calls `CacheData("blog:<id>", 60, ...)`. The loader queries MySQL, marshals the blog to JSON and **sleeps 2 seconds** to simulate a slow query.

Result: the first request for an ID takes about 2 s, repeated requests within 60 s return immediately from Memcached.

## Concepts you can learn

- **Cache-aside / read-through**: check the cache, fall back to the source, then populate the cache so later reads are cheap.
- **TTL-based invalidation**: entries expire after 60 s, which bounds how stale the data can be.
- **Cache keys**: a namespaced key (`blog:<id>`) per resource.
- **Generic loaders via reflection**: `CacheData` accepts any function returning `[]byte` (an older alternative to generics).

## How to run

Prerequisites:

- MySQL on `127.0.0.1:3306` with database `net_http_blogs` and user `net_http_blogs` / password `tmp_pwd`. The `blogs` table is created by `AutoMigrate`; insert rows yourself.
- Memcached on `localhost:11211`.

```bash
cd memcache
go run .
curl http://localhost:8080/blogs/1   # ~2 s (miss)
curl http://localhost:8080/blogs/1   # fast (hit)
```

## Notes

- A missing ID is not treated as an error: the zero-value `Blog` is marshaled and cached for 60 s.
- If Memcached is down, every `Get` fails, the loader runs, then `cache.Set` **panics** (net/http recovers it per request, so the client gets no response body).
- Concurrent misses for the same key all run the slow loader (no request coalescing).
- `LICENSE` in this directory is Apache License 2.0.
