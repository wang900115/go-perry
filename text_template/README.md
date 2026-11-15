# text_template

`text/template` and `html/template`: actions (`if`/`range`/`else`), whitespace trimming
with `{{-` and `-}}`, method calls, pipelines, `FuncMap`, built-in functions, Go 1.24
range-over-int, `define`/`template`/`block` with `Clone`, `missingkey=error`,
partial output on execution errors, and context-aware auto-escaping in `html/template`.

```bash
go run .
```

Standard library only. This module declares `go 1.24` (`{{range 3}}` needs Go 1.24).
Verified on go1.26 linux/amd64 (2 runs, identical output); `go vet .` and `GOOS=windows go vet .` pass.
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: actions, range, if, whitespace trim ==
Order #7 for Perry
  0. keyboard $120
  1. mouse $40
Total: $160 [UNPAID]
Order #8 for Ann
  (no items)
Total: $0 [PAID]
== Part 2: FuncMap, pipelines, built-ins ==
KEYBOARD   tags=usb,rgb n=2 (expensive)
MOUSE      tags=wireless n=1
first tag of first item: usb
range over int (Go 1.24+): 0 1 2
discount: 16
parse without Funcs: template: bad:1: function "upper" not defined
== Part 3: define / template / block ==
<Default title> base body for perry
<Custom title> base body for perry
defined: [base body page title]
== Part 4: missing keys and errors ==
hello perry, role=<no value>
partial output="hello perry, role="
err=template: m:1:24: executing "m" at <.role>: map has no entry for key "role"
partial output="total=50 discount="
err=template: d:1:28: executing "d" at <.Discount>: error calling Discount: invalid coupon BAD
unknown field err: template: f:1:2: executing "f" at <.Nope>: can't evaluate field Nope in type main.Order
== Part 5: text/template vs html/template ==
text: <a href="/u?q=<b>"Tom" & 'Jerry'</b>" onclick="greet('<b>"Tom" & 'Jerry'</b>')"><b>"Tom" & 'Jerry'</b></a>
html: <a href="/u?q=%3cb%3e%22Tom%22%20%26%20%27Jerry%27%3c%2fb%3e" onclick="greet('\u003cb\u003e\u0022Tom\u0022 \u0026 \u0027Jerry\u0027\u003c\/b\u003e')">&lt;b&gt;&#34;Tom&#34; &amp; &#39;Jerry&#39;&lt;/b&gt;</a>
trusted: <b>bold</b>
```
