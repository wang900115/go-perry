# 2026-11-15 · text/template & html/template: actions, pipelines, FuncMap, safe HTML

## Why learn this
- Generating text from data is everywhere: code generators (`go generate`, `stringer`-style tools), config files (Helm charts and many CLIs use Go templates), emails, reports, and server-rendered HTML.
- `go` tools expose templates directly, e.g. `go list -f '{{.ImportPath}}'` and `docker inspect --format '{{.State.Status}}'`.
- `html/template` is the standard library's defense against XSS. Using `text/template` for HTML by mistake is a real security bug.
- It builds on earlier lessons: `reflect` (templates look up fields and methods by reflection), `io_reader_writer` (`Execute` writes to an `io.Writer`), `error_wrapping` (execution errors).

## Core concepts
1. **Parse once, execute many times.** `template.New(name).Parse(src)` builds a parse tree; `Execute(w, data)` walks it and writes to `w`. `template.Must` panics on a parse error, which is fine for templates known at compile time (https://pkg.go.dev/text/template).
2. **Dot (`.`) is the current data.** `{{.Field}}`, `{{.Method}}`, `{{.Method "arg"}}`, and `{{.key}}` for maps. `range` and `with` change dot; `$` is always the root data.
3. **Actions**: `{{if}}…{{else}}…{{end}}`, `{{range $i, $v := .Items}}…{{else}}…{{end}}` (the `else` runs when the collection is empty), `{{with}}`, and since Go 1.24 `{{range 3}}` and range over iterator functions (https://go.dev/doc/go1.24#texttemplatepkgtexttemplate).
4. **Whitespace trimming**: `{{-` trims whitespace (including newlines) before the action, `-}}` trims after. The space between `-` and the content is required.
5. **Pipelines**: `{{.Name | upper | printf "%-10s"}}` passes the result of each command as the LAST argument of the next one.
6. **Functions**: built-ins include `len`, `index`, `printf`, `eq`/`ne`/`lt`/`gt`, `and`/`or`/`not`, `html`, `js`, `urlquery`. Your own functions go in a `template.FuncMap` and must be added with `.Funcs(...)` **before** `Parse`. A function (or method) may return `(value, error)`; a non-nil error stops execution.
7. **Associated templates**: `{{define "name"}}`, `{{template "name" .}}`, and `{{block "name" .}}default{{end}}` (= define + template). `Clone` lets you override a block in a copy without touching the original.
8. **Missing keys**: for maps, a missing key prints `<no value>` by default; `Option("missingkey=error")` makes it an error. A missing struct field is always an error.
9. **`html/template`** has the same API but escapes each `{{.}}` according to where it appears: HTML text, attribute, URL, JavaScript, CSS (https://pkg.go.dev/html/template). Typed strings such as `template.HTML` mark content as trusted and skip escaping.

### How it works (standard library)
Based on the go1.26 source and docs; internal details may change.
- `Parse` uses `text/template/parse` to lex and parse the source into a tree of nodes (text, action, if, range, pipe…). Function names are checked at parse time, which is why `Funcs` must come first.
- `Execute` walks the tree with a `state` that holds the writer, the current dot (a `reflect.Value`) and the variable stack. Field and method lookups use `reflect` on every execution. I have not checked whether lookups are cached across executions; as far as I know they are not.
- Output is written to the writer as execution goes. There is no buffering or rollback, so an error in the middle leaves partial output in `w`.
- Once parsing is done, a template can be executed by many goroutines at the same time. If they share a writer, the output can be interleaved (see the package overview at https://pkg.go.dev/text/template).
- `html/template` wraps a `text/template`. On the first `Execute` it runs an escaper that analyzes the context of each action (e.g. "inside a double-quoted href attribute") and rewrites the tree to add escaping functions to the pipelines. After that, the template can no longer be changed with `Parse` or cloned. I think the exact error text differs between methods; I have not checked each one.

## Full example
See `main.go`. Verified with `go run .` (2 runs, identical output) and `go vet` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. `range`/`else`/`if`, method call `{{.Total}}`, whitespace trimming.
2. `FuncMap` + pipelines, `len`/`gt`/`index`, `{{range 3}}`, a method with an argument, and the parse error you get when `Funcs` is missing.
3. `define`/`block`/`template`, overriding `title` in a `Clone`, listing templates (sorted, because `Templates()` comes from a map).
4. `<no value>` vs `missingkey=error`, partial output left in the buffer after an error, unknown struct field.
5. The same template with `text/template` and `html/template`: no escaping vs URL, JavaScript and HTML escaping in one line.

## Common pitfalls
1. **Partial output on error.** `Execute` streams to the writer. If you write straight to an `http.ResponseWriter` and an error happens halfway, the client gets half a page with a 200 status. Execute into a `bytes.Buffer` first, then copy it to the response only on success.
2. **`text/template` for HTML.** It does not escape anything (Part 5). Always use `html/template` for HTML. Do not wrap user input in `template.HTML`/`template.JS` to "fix" output; that turns escaping off.
3. **Silent `<no value>` and map ordering assumptions.** A typo in a map key prints `<no value>` instead of failing. Use `missingkey=error` for configs and generated code. (`range` over a map in a template is sorted by key, unlike a Go `for range`; but `Templates()`/`DefinedTemplates()` are not sorted.)

## Exercise (no answer given)
Write a small code generator: given `[]Field{{Name: "ID", Type: "int"}, {Name: "Email", Type: "string"}}` and a type name `User`, generate the Go source of a struct and a `String()` method that prints all fields. Requirements:
- Use a `FuncMap` with a `lower` function for the struct tags (`json:"id"`).
- Use `{{-`/`-}}` so the output has no blank lines inside the struct.
- Run the output through `go/format.Source` and print it; make `go/format` fail on purpose once to see the error.
- Use `missingkey=error` and pass a `map[string]any` once with a misspelled key to see the error.

Suggested directory: `text_template` (this one), or `codegen_template` if you keep the exercise separate.

## References
- https://pkg.go.dev/text/template
- https://pkg.go.dev/html/template
- https://go.dev/doc/go1.24#texttemplatepkgtexttemplate
