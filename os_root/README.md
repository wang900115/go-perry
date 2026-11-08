# os_root

`os.Root` (Go 1.24, more methods in Go 1.25): traversal-resistant file access confined to a single directory tree. Compares a naive `filepath.Join`, a string-only `filepath.IsLocal` check, and `os.Root` against `../` and symlink escapes.

## Run

```
go run .
```

Requires Go 1.25 (`Root.ReadFile`, `Root.WriteFile`, `Root.MkdirAll`). On Windows, creating symlinks needs Developer Mode or admin rights; without it the program prints a notice and skips the symlink cases.

## What the code does

- **Part 1 – naive `filepath.Join`**: `../secret.txt` and the `link_out` symlink both leak the secret.
- **Part 2 – `filepath.IsLocal`**: blocks `../secret.txt` lexically, but the symlink still escapes.
- **Part 3 – `os.Root`**: blocks both; the in-tree symlink `link_in` still works.
- **Part 4 – writes**: `Root.WriteFile("../pwned.txt")` fails; `MkdirAll` + `WriteFile` inside the root succeed.
- **Part 5 – `Root.FS()`**: use the root as an `fs.FS` (e.g. `fs.WalkDir`, `http.FileServerFS`).

Expected output (observed on Linux with go1.26.3; `GOOS=windows go build` also succeeds):

```
== 1. naive filepath.Join ==
  hello.txt      -> "hello"
  docs/a.txt     -> "doc A"
  ../secret.txt  -> "TOP SECRET"
  link_in        -> "doc A"
  link_out       -> "TOP SECRET"
== 2. filepath.IsLocal check (string-only) ==
  hello.txt      -> "hello"
  docs/a.txt     -> "doc A"
  ../secret.txt  -> ERROR (rejected: not a local path)
  link_in        -> "doc A"
  link_out       -> "TOP SECRET"
== 3. os.Root ==
  hello.txt      -> "hello"
  docs/a.txt     -> "doc A"
  ../secret.txt  -> ERROR (path escapes from parent)
  link_in        -> "doc A"
  link_out       -> ERROR (path escapes from parent)
== 4. writes are confined too ==
WriteFile ../pwned.txt -> ERROR (path escapes from parent)
WriteFile uploads/2026/ok.txt -> ok
== 5. Root.FS() as an fs.FS ==
   docs/a.txt
   hello.txt
   link_in (symlink)
   link_out (symlink)
   uploads/2026/ok.txt
fs.ReadFile ../secret.txt -> ERROR (invalid path)
```

See [LESSON.md](LESSON.md) for the full lesson.
