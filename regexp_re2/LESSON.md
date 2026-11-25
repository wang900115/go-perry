# Regular Expressions with RE2

> Go version used: `go1.26.3`.

## Why learn this

Regular expressions show up in everyday Go code:

- Validating input (IDs, dates, slugs) in HTTP handlers and CLIs
- Parsing log lines and extracting fields
- Rewriting text: renaming, redacting secrets, templating
- Config and routing rules where users provide the pattern

Go's `regexp` is different from PCRE (Perl, Python, JavaScript, Java). It uses the **RE2** design:
matching time is linear in the input size, so a user-provided pattern cannot freeze your server
(no "ReDoS"). The price is that some syntax (backreferences, lookaround) does not exist.
([regexp package docs](https://pkg.go.dev/regexp))

## Core concepts

### 1. Compile once, use many times

```go
var dateRE = regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
```

- `regexp.Compile` returns `(*Regexp, error)`: use it for patterns from users or config.
- `regexp.MustCompile` panics on error: use it for constant patterns at package level, so a bad pattern
  fails at program start, not on the first request.
- Use raw string literals (backquotes) so `\d` does not need to be written as `"\\d"`.
- Compiling is expensive compared with matching. Do not call `MustCompile` inside a hot loop or handler.
- A `*Regexp` is safe for concurrent use by multiple goroutines.

### 2. The `Find` method family

There are 16 `Find` methods; their names follow a pattern
([regexp docs](https://pkg.go.dev/regexp#pkg-overview)):

```
Find(All)?(String)?(Submatch)?(Index)?
```

- `All`: return all non-overlapping matches; the extra `n` argument limits the count (`-1` = no limit)
- `String`: input/output are `string` instead of `[]byte`
- `Submatch`: also return capture groups (index 0 is the whole match)
- `Index`: return byte positions instead of text

"No match" is `nil` (or `""` for `FindString`, which cannot be told apart from an empty match;
use `FindStringIndex` when that matters).

### 3. Matching is unanchored

`MatchString` reports whether the pattern matches **anywhere** in the input.
For validation you almost always want `^...$`. Without `(?m)`, `^`/`$` mean the start/end of the text;
`\A` and `\z` always do.

### 4. How RE2 works (why it is linear)

A backtracking engine (PCRE) tries one alternative, and if it fails, goes back and tries the next.
For patterns like `(a+)+$`, the number of ways to split the input grows exponentially.

Go compiles the pattern into a program for a small virtual machine and simulates a
**Thompson NFA**: it keeps the *set* of all possible states and advances them together, one input
rune at a time. Each rune costs at most O(pattern size), so the total is O(pattern × input).
The implementation picks the fastest of several strategies (one-pass, bounded backtracker
with a visited bitmap for small inputs, or the full NFA in `regexp/backtrack.go` and `regexp/exec.go`),
all of which keep the linear bound. See Russ Cox's article
[Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html).

Consequences:

- No backreferences (`\1`) and no lookahead/lookbehind (`(?=`, `(?!`, `(?<=`, `(?<!`): these cannot be
  done in linear time. `Compile` returns an error.
- Since Go 1.22, `(?<name>...)` is accepted as well as `(?P<name>...)`. That is why `(?<!x)y` reports
  "invalid named capture" instead of "lookbehind not supported".
- Go's `regexp` is sometimes slower than PCRE on simple patterns, but it never explodes.

### 5. Leftmost-first vs leftmost-longest

By default Go uses Perl semantics: among matches starting at the leftmost position, it picks the one
a backtracking engine would find first. So `a|ab` on `"ab"` matches `"a"`.
`re.Longest()` (or `regexp.CompilePOSIX`) switches to leftmost-longest: `"ab"`.
Lazy quantifiers `*?`, `+?`, `??` are supported.

### 6. Replacement templates

In `ReplaceAllString` the replacement is expanded with `Expand`:
`$1`, `${1}`, `$name`, `${name}`, and `$$` for a literal `$`
([Regexp.Expand](https://pkg.go.dev/regexp#Regexp.Expand)).
The name after `$` is taken as long as possible: `$1x` means the group named `1x`,
which does not exist, so it expands to an empty string. Write `${1}x`.

Use `ReplaceAllLiteralString` when the replacement must not be expanded, and
`ReplaceAllStringFunc` when it must be computed.

### 7. Flags and UTF-8

- `(?i)` case-insensitive, `(?m)` multi-line `^`/`$`, `(?s)` lets `.` match `\n`, `(?U)` swaps greedy/lazy.
- Flags can be scoped: `(?i:abc)def`.
- Matching is on UTF-8: `.` matches one rune, not one byte. `\w`, `\d`, `\b` are ASCII-only;
  use `\p{L}`, `\p{Han}` etc. for Unicode classes. Syntax reference:
  [regexp/syntax](https://pkg.go.dev/regexp/syntax).

## Full example

See [main.go](main.go). Run with `go run .`.

Expected output (Go 1.26.3):

```
== 1. matching is unanchored ==
  \d+ matches "abc123":      true
  ^\d+$ matches "abc123":    false
  FindString:                42
  FindAllString(n=-1):       [42 7 1000]
  FindAllString(n=2):        [42 7]
  FindStringIndex:           [3 5]
  no match -> "" and true
== 2. capture groups ==
  FindStringSubmatch: ["2026-11-25" "2026" "11" "25"]
  SubexpNames:        [ year month day]
  month via SubexpIndex: 11
  all -> 2024 01 02
  all -> 2025 03 04
== 3. replacement ==
  ${day}/${month}/${year}:  25/11/2026
  "$1x"   -> ""
  "${1}x" -> "bobx"
  Literal: "cost: $5"
  Expand:  "cost: "
  Func:    GO IS FUN
== 4. leftmost-first vs leftmost-longest ==
  a|ab on "ab" (default):   a
  a|ab on "ab" (Longest):   ab
  lazy  <.+?>: <b>
  greedy <.+>: <b>bold</b>
== 5. flags ==
  ^\w+ (no m):   [first]
  (?m)^\w+:      [first second]
  (?i)GO:        [go Go GO]
  a.b on "a\nb": false, (?s)a.b: true
  . matches one rune:  true
== 6. unsupported syntax (Compile returns an error) ==
  (\w)\1       -> error parsing regexp: invalid escape sequence: `\1`
  foo(?=bar)   -> error parsing regexp: invalid or unsupported Perl syntax: `(?=`
  (?<!x)y      -> error parsing regexp: invalid named capture: `(?<!x)y`
  [a-          -> error parsing regexp: missing closing ]: `[a-`
== 7. linear time guarantee ==
  (a+)+$ on 50000 'a' + '!': matched = false , under 1s = true
== 8. QuoteMeta for user input ==
  QuoteMeta: 1\.5\+2
  raw pattern matches "1x55552":  true
  quoted pattern matches "1x55552": false
== 9. concurrent use ==
  words counted per goroutine: [3000 3000 3000 3000]
```

## Common pitfalls

### Pitfall 1: Forgetting anchors in validation

```go
regexp.MustCompile(`\d{4}`).MatchString("abc12345xyz") // true!
```

Use `^\d{4}$`. Also remember `$` without `(?m)` means end of text, so a trailing `"\n"` is not allowed;
this differs from some other engines.

### Pitfall 2: `$1x` and `$` in replacements

`ReplaceAllString(s, "$1x")` silently produces an empty string for that part. Write `${1}x`.
If the replacement comes from user data (prices like `"$5"`), use `ReplaceAllLiteralString`.

### Pitfall 3: Compiling in a loop, or passing user input unescaped

`regexp.MustCompile(userInput)` both panics on bad input and treats `.`, `+`, `(` as operators.
Use `regexp.Compile` and handle the error, or `regexp.QuoteMeta(userInput)` if the user meant a literal string.
Compile once and reuse the `*Regexp`; for simple cases, `strings.Contains`/`strings.Cut` are much faster.

## Exercise

Write a small log parser for lines like:

```
2026-11-25T20:13:15Z level=ERROR user=alice msg="disk full" latency=120ms
```

1. Use one compiled regexp with **named groups** to extract `time`, `level`, `user`, `msg` and `latency`
   into a `map[string]string`. Lines that do not match must be reported, not silently skipped.
2. Redact every `user=<name>` to `user=***` with a replacement template. Then do it again but keep the
   first letter (`user=a***`) using `ReplaceAllStringFunc` or `${...}`.
3. Accept a search keyword from `os.Args` and highlight it as `[keyword]` in the output.
   Make sure a keyword like `a+b` or `(x)` works literally.
4. Generate a 1 MB input and compare the time of your regexp with a version built on
   `strings.Fields`/`strings.Cut`. Explain the difference.

## References

- regexp package: https://pkg.go.dev/regexp
- Syntax: https://pkg.go.dev/regexp/syntax
- RE2 syntax wiki: https://github.com/google/re2/wiki/Syntax
- Russ Cox, Regular Expression Matching Can Be Simple And Fast: https://swtch.com/~rsc/regexp/regexp1.html
- Russ Cox, Regular Expression Matching: the Virtual Machine Approach: https://swtch.com/~rsc/regexp/regexp2.html
