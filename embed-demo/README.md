# embed-demo

A minimal `//go:embed` example that bakes a text file into the compiled binary at build time.

## What the code does

`main.go` declares `var content string` with the directive `//go:embed assets/hello.txt` and prints `content`. The file is read by the compiler when you build, not at runtime, so the binary prints the file's contents even if it later runs somewhere `assets/hello.txt` doesn't exist.

## Concepts you can learn

- **`//go:embed`**: a compiler directive that embeds files into the binary. It needs the `embed` package imported (a blank `_ "embed"` import is enough for `string`/`[]byte` targets).
- **Target types**: `string` and `[]byte` hold a single file, and `embed.FS` holds a tree of files. This demo uses `string`.
- **Single-file deployment**: templates, SQL, or static assets ship inside the executable, so there are no runtime path problems.
- **Build-time failure**: if the pattern matches no files, the build fails instead of the program failing at runtime.

## How to run

```bash
cd embed-demo
mkdir -p assets
echo "Hello, embed!" > assets/hello.txt   # any content
go run .
```

## Notes

- **`assets/hello.txt` isn't in the repository**, so the build currently fails with `main.go:8:12: pattern assets/hello.txt: no matching files found`. Create the file first, as shown above.
- The directory contains a committed `app.exe` (a prebuilt Windows binary). Its contents weren't inspected, and you don't need it to run the demo.
