# string_utf8

How Go strings really work: a string is a read-only sequence of bytes, `len`
counts bytes, indexing returns a byte, and `for range` decodes UTF-8 into runes
(Unicode code points). Covers invalid UTF-8 and `U+FFFD`, the cost of
`[]byte(s)` / `[]rune(s)` conversions, truncating without cutting a rune in
half, `strings.Builder` vs `+=`, and why one visible character can be several
runes.

```bash
go run .
```

Standard library only. Verified on go1.26.3 and go1.24.7 linux/amd64 (identical
output); `go vet .` and `GOOS=windows go vet .` pass. The allocation counts in
Part 5 depend on compiler optimizations and may differ between Go versions.
On Windows, use a UTF-8 capable terminal (Windows Terminal) to see the CJK
characters correctly.
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: bytes vs runes ==
s="Go語言" len=8 runes=4
s[2]=232 (0xe8) is a byte, not a character
bytes: 47 6f e8 aa 9e e8 a8 80
  offset=0 rune='G' U+0047 width=1
  offset=1 rune='o' U+006F width=1
  offset=2 rune='語' U+8A9E width=3
  offset=5 rune='言' U+8A00 width=3

== Part 2: invalid UTF-8 ==
bad="a\xffb" len=3 valid=false
  offset=0 rune='a'
  offset=1 rune='�'
  offset=2 rune='b'
DecodeRuneInString: rune='�' size=1 (RuneError=true)
ToValidUTF8: "a?b"

== Part 3: conversions ==
s="héllo" b="Héllo" (the string did not change)
len([]byte)=6 len([]rune)=5
string(rs)="hello"
string(rune(65))="A"  string(rune(0x8A9E))="語"
string(rune(-1))="�" (invalid code point becomes U+FFFD)

== Part 4: safe truncation ==
naive s[:4]="Go\xe8\xaa" valid=false
truncateRunes(s, 3)="Go語"
truncateBytes(s, 4)="Go"
truncateBytes(s, 8)="Go語言"

== Part 5: building strings ==
words=9 allocs with +=: 9, with Builder+Grow: 1
builder="Go語1" len=6

== Part 6: a rune is not a visible character ==
"é" vs "é" equal=false runes=1,2 bytes=2,3
flag runes=2 bytes=8
EqualFold("Go", "GO")=true
```
