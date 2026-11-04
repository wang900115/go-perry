# build_info

Embed version metadata in a Go binary and serve it over HTTP. Values injected with `-ldflags -X` take priority; `runtime/debug.ReadBuildInfo` fills in the rest.

## What the code does

`main.go`:

- Declares `Version = "dev"` and `Commit = "none"` as package variables that `-ldflags` can override.
- `versionInfo()` returns a `BuildInfo{Version, Commit, GoVersion, Modified}` struct:
  - `GoVersion` comes from `debug.ReadBuildInfo()`.
  - If `Version` is still `"dev"`, it is replaced by the main module version, unless that version is empty or `(devel)`.
  - If `Commit` is still `"none"`, it is replaced by the `vcs.revision` build setting, shortened to 7 characters.
  - `Modified` comes from the `vcs.modified` build setting.
- `main` prints `Build info: {...}`, registers `GET /version` (returns the struct as JSON) and listens on `:8080`.

## Concepts you can learn

- **`-ldflags "-X main.Var=value"`**: sets a string variable at link time. This is the usual way to stamp a release version and commit from CI without generating code.
- **`debug.ReadBuildInfo`**: reads metadata the toolchain embeds in every module-aware binary (Go version, module version, VCS settings), so a binary can describe itself at runtime.
- **VCS stamping**: `vcs.revision` and `vcs.modified` are recorded when the binary is built inside a VCS checkout (controlled by `-buildvcs`).
- **Precedence**: values set explicitly with ldflags always win, and the embedded metadata is only a fallback.
- **Version endpoints**: `/version` lets you see which build is deployed.

## How to run

```bash
cd build_info
go run .
# then: curl http://localhost:8080/version

go build -ldflags "-X main.Version=v1.2.3 -X main.Commit=abc1234" -o res .
```

## Notes

- Errors from `http.ListenAndServe` and from the JSON encoder are ignored. If port 8080 is already in use, the program prints the build info and exits silently.
- When no VCS information is embedded, `Modified` is an empty string and `Commit` stays `none`.
- `res`, `res.exe` and `res_id.exe` are committed Windows binaries. The code does not show which flags they were built with.
