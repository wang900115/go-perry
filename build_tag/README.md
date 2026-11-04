# build_tag

Use `//go:build` constraints to compile different files for different environments (custom tags) and operating systems.

## What the code does

- `main.go` declares the package variables `env` and `os` and prints them:
  ```
  Environment is :  <env>
  Running on: <os>
  ```
- `main_production.go` (`//go:build production`) sets `env = "production"` in `init`.
- `main_staging.go` (`//go:build staging`) sets `env = "staging"`.
- `main_linux.go` (`//go:build linux`) sets `os = "Linux"`.
- `main_window.go` (`//go:build windows`) sets `os = "Windows"`.

Example output on Windows with `-tags production`: `Environment is :  production` / `Running on: Windows`. Without any tag, `env` is empty.

## Concepts you can learn

- **Custom build tags**: `-tags production` includes or excludes whole files at compile time, so environment-specific code never ships in the other build.
- **OS constraints**: `//go:build linux` / `//go:build windows` select per-platform code. The `_linux.go` filename suffix applies the same constraint implicitly.
- **`init` for wiring**: each optional file only sets a variable in `init`, so `main.go` does not need to know which files were compiled in.

## How to run

```bash
cd build_tag
go run .                    # env empty
go run -tags production .
go run -tags staging .
```

## Notes

- `main_window.go` is named `window`, not `windows`, so the filename suffix does not apply. Only the `//go:build windows` line restricts it.
- On operating systems other than Linux and Windows (for example macOS), `os` stays empty.
- Passing both `production` and `staging` compiles both `init` functions. They run in file-name order, so `staging` would win; nothing prevents combining them.
- The package variable named `os` shadows the name of the standard `os` package inside this package.
