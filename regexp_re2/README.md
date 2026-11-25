# regexp_re2

Regular expressions with Go's `regexp` package (RE2 engine).

- `MatchString`/`Find*` match anywhere unless you anchor with `^`/`$` (or `\A`/`\z`).
- `Find*` naming: `All` = every match (`n=-1`), `String` = string input, `Submatch` = groups, `Index` = positions.
- Named groups `(?P<name>...)` and `SubexpIndex`; replacements with `${name}` / `${1}` (beware `$1x`).
- `ReplaceAllLiteralString` and `ReplaceAllStringFunc` when `$` must not expand.
- Leftmost-first (Perl) by default; `Longest()` switches to leftmost-longest (POSIX).
- Flags `(?i)`, `(?m)`, `(?s)`; `.` matches one UTF-8 rune.
- No backreferences or lookaround: RE2 guarantees linear-time matching, so `(a+)+$` is safe.
- `regexp.QuoteMeta` for user input; a `*Regexp` is safe for concurrent use.

```bash
go run .
```

See [LESSON.md](LESSON.md) for the full lesson.
