# 2026-11-08 · os.Root: traversal-resistant file access

## Why learn this
- **Any time a file name comes from outside**: HTTP upload/download handlers, static file servers, "export to path" options, template loaders, plugin directories.
- **Archive extraction**: a zip/tar entry named `../../.bashrc` ("Zip Slip") or a symlink entry pointing outside the target directory.
- **Path traversal is a top CWE** (CWE-22). Checking strings is easy to get wrong; `os.Root` lets the OS enforce the boundary.

This follows the `sync_*` series: after a few concurrency primitives, a standard-library security primitive that is small, recent (Go 1.24/1.25), and runs on Windows with no external services.

## Core concepts
1. **`os.OpenRoot(dir)`** returns a `*os.Root`. Every method (`Open`, `Create`, `OpenFile`, `ReadFile`, `WriteFile`, `Mkdir`, `MkdirAll`, `Remove`, `Rename`, `Stat`, `Lstat`, `Symlink`, …) takes a name **relative to the root** and fails if the name resolves outside it (https://pkg.go.dev/os#Root).
2. **What is blocked**: `..` that climbs above the root, absolute paths, and **symlinks whose target is outside the root** — including absolute symlinks. `..` that stays inside (`docs/../hello.txt`) is allowed.
3. **Symlinks inside the root are followed** (`link_in` works).
4. **Shortcut**: `os.OpenInRoot(dir, name)` opens a single file without keeping a `Root` around.
5. **`root.FS()`** returns an `fs.FS`, so it plugs into `fs.WalkDir`, `fs.ReadFile`, `http.FileServerFS`, `template.ParseFS`. Note `fs.FS` additionally requires `fs.ValidPath` names (no `..`, no leading `/`), so `fs.ReadFile(fsys, "../x")` fails with `fs.ErrInvalid` before the OS is even asked.
6. **Version timeline**: Go 1.24 introduced `os.Root` with the core methods (https://go.dev/doc/go1.24#directory-limited-filesystem-access). Go 1.25 added `Chmod`, `Chown`, `Chtimes`, `Lchown`, `Link`, `MkdirAll`, `ReadFile`, `Readlink`, `RemoveAll`, `Rename`, `Symlink`, `WriteFile` (https://go.dev/doc/go1.25#os). Background: https://go.dev/blog/osroot.

### How it works
- **Why string checks are not enough**: `filepath.IsLocal` / `filepath.Clean` only look at the text. They cannot know that `link_out` is a symlink to `../secret.txt`, and there is a TOCTOU window: an attacker can swap a directory for a symlink between your check and your `open`.
- **Unix** (simplified): `OpenRoot` opens the directory and keeps the file descriptor. Each operation walks the path one component at a time with `openat(dirfd, component, O_NOFOLLOW…)`, relative to the already-opened directory fd. When a component is a symlink, the runtime reads it with `readlinkat` and resolves it itself, tracking depth so it never climbs above the root. Because everything is relative to the fd, renaming or replacing the root path afterward does not redirect it (moving the root directory keeps the `Root` pointed at the moved directory).
- **Windows**: the same idea with a directory handle and `NtCreateFile` relative to a parent handle (`RootDirectory` in `OBJECT_ATTRIBUTES`). Reserved device names such as `NUL` and `COM1` are rejected. **Uncertain**: I have not checked the Windows implementation line by line; the doc comment on `os.Root` is the authoritative statement of guarantees.
- **The error** is a `*fs.PathError` whose inner message is `path escapes from parent`. The sentinel is unexported (to my knowledge in 1.26), so match on the behavior (non-nil error) rather than on the exact text.

## Full example
See `main.go` in this directory (`go run .` and `go vet` on go1.26.3 linux/amd64; `GOOS=windows go build` also succeeds). Requires Go 1.25.

Expected output: see [README.md](README.md). On Windows without Developer Mode the symlink rows are skipped and a notice is printed instead.

## Common pitfalls
1. **Thinking `filepath.Clean` / `IsLocal` / `strings.HasPrefix(path, base)` is enough.** Lexical checks miss symlinks and race with concurrent filesystem changes. `HasPrefix("/srv/data-evil", "/srv/data")` is also true.
2. **Mixing `os.*` and `root.*` calls.** Once you have a root, every access must go through it. One stray `os.Open(filepath.Join(dir, name))` reopens the hole. Likewise, a `*os.File` returned by the root is fine, but re-deriving a path from `f.Name()` and calling `os.Open` on it is not.
3. **Over-trusting the boundary.** `os.Root` does not stop traversal of mount points, Linux bind mounts, `/proc` special files, or Unix device files; and on Unix `Chmod`/`Chown`/`Chtimes` have a documented race. It confines *names*, it is not a sandbox. Also, only names are relative — passing an absolute path is an error, not "relative to the root".

## Exercise
Write `SafeExtract(zipPath, destDir string) error` that extracts a zip archive using `archive/zip` and an `os.Root` opened on `destDir`:
- Every entry must be written via the root (`MkdirAll`, `OpenFile`/`Create`); no `os.*` path calls on entry names.
- Build a test archive in code containing a normal file, a nested file `a/b/c.txt`, and an evil entry `../../evil.txt`. Your function must extract the good entries and return an error (or skip with a log line) for the evil one.
- Bonus: limit total uncompressed size to guard against zip bombs.
