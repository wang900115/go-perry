// Regular expressions in Go: the regexp package and the RE2 engine.
//
// Run: go run .
package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Compile once at package level. MustCompile panics on a bad pattern,
// which is fine for patterns written in source code.
var (
	dateRE = regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
	wordRE = regexp.MustCompile(`\w+`)
)

// 1. Basic matching: Match* finds a match ANYWHERE unless you anchor it.
func basics() {
	fmt.Println("== 1. matching is unanchored ==")
	digits := regexp.MustCompile(`\d+`)
	fmt.Println("  \\d+ matches \"abc123\":     ", digits.MatchString("abc123"))
	full := regexp.MustCompile(`^\d+$`)
	fmt.Println("  ^\\d+$ matches \"abc123\":   ", full.MatchString("abc123"))
	fmt.Println("  FindString:               ", digits.FindString("id=42, n=7"))
	fmt.Println("  FindAllString(n=-1):      ", digits.FindAllString("id=42, n=7, x=1000", -1))
	fmt.Println("  FindAllString(n=2):       ", digits.FindAllString("id=42, n=7, x=1000", 2))
	fmt.Printf("  FindStringIndex:           %v\n", digits.FindStringIndex("id=42"))
	fmt.Printf("  no match -> %q and %v\n", digits.FindString("none"), digits.FindStringIndex("none") == nil)
}

// 2. Capture groups and named groups.
func groups() {
	fmt.Println("== 2. capture groups ==")
	m := dateRE.FindStringSubmatch("released on 2026-11-25.")
	fmt.Printf("  FindStringSubmatch: %q\n", m)
	fmt.Println("  SubexpNames:       ", dateRE.SubexpNames())
	fmt.Println("  month via SubexpIndex:", m[dateRE.SubexpIndex("month")])

	all := dateRE.FindAllStringSubmatch("2024-01-02 and 2025-03-04", -1)
	for _, g := range all {
		fmt.Println("  all ->", g[1], g[2], g[3])
	}
}

// 3. Replacement: $1 / ${name} expansion, and the $1x pitfall.
func replace() {
	fmt.Println("== 3. replacement ==")
	s := "2026-11-25"
	fmt.Println("  ${day}/${month}/${year}: ", dateRE.ReplaceAllString(s, "${day}/${month}/${year}"))

	// "$1x" means the group named "1x", which does not exist -> empty string.
	re := regexp.MustCompile(`(\w+)@example\.com`)
	fmt.Printf("  \"$1x\"   -> %q\n", re.ReplaceAllString("bob@example.com", "$1x"))
	fmt.Printf("  \"${1}x\" -> %q\n", re.ReplaceAllString("bob@example.com", "${1}x"))

	// Literal replacement: "$" is not special.
	price := regexp.MustCompile(`PRICE`)
	fmt.Printf("  Literal: %q\n", price.ReplaceAllLiteralString("cost: PRICE", "$5"))
	fmt.Printf("  Expand:  %q\n", price.ReplaceAllString("cost: PRICE", "$5"))

	// Func replacement: compute the new text.
	fmt.Println("  Func:   ", wordRE.ReplaceAllStringFunc("go is fun", strings.ToUpper))
}

// 4. Leftmost-first (Perl-like) vs leftmost-longest (POSIX).
func leftmost() {
	fmt.Println("== 4. leftmost-first vs leftmost-longest ==")
	re := regexp.MustCompile(`a|ab`)
	fmt.Println("  a|ab on \"ab\" (default):  ", re.FindString("ab"))
	re.Longest()
	fmt.Println("  a|ab on \"ab\" (Longest):  ", re.FindString("ab"))

	lazy := regexp.MustCompile(`<.+?>`)
	greedy := regexp.MustCompile(`<.+>`)
	html := "<b>bold</b>"
	fmt.Println("  lazy  <.+?>:", lazy.FindString(html))
	fmt.Println("  greedy <.+>:", greedy.FindString(html))
}

// 5. Flags and anchors: (?i), (?m), (?s), \A, \z.
func flags() {
	fmt.Println("== 5. flags ==")
	text := "first line\nsecond line"
	fmt.Println("  ^\\w+ (no m):  ", regexp.MustCompile(`^\w+`).FindAllString(text, -1))
	fmt.Println("  (?m)^\\w+:     ", regexp.MustCompile(`(?m)^\w+`).FindAllString(text, -1))
	fmt.Println("  (?i)GO:       ", regexp.MustCompile(`(?i)GO`).FindAllString("go Go GO", -1))
	fmt.Printf("  a.b on \"a\\nb\": %v, (?s)a.b: %v\n",
		regexp.MustCompile(`a.b`).MatchString("a\nb"),
		regexp.MustCompile(`(?s)a.b`).MatchString("a\nb"))
	fmt.Println("  . matches one rune: ", regexp.MustCompile(`^.$`).MatchString("é")) // "é" is 2 bytes, 1 rune
}

// 6. RE2 refuses features that need backtracking.
func unsupported() {
	fmt.Println("== 6. unsupported syntax (Compile returns an error) ==")
	for _, p := range []string{`(\w)\1`, `foo(?=bar)`, `(?<!x)y`, `[a-`} {
		_, err := regexp.Compile(p)
		fmt.Printf("  %-12s -> %v\n", p, err)
	}
}

// 7. Linear time: the classic "catastrophic backtracking" pattern is fast.
func linearTime() {
	fmt.Println("== 7. linear time guarantee ==")
	re := regexp.MustCompile(`^(a+)+$`)
	input := strings.Repeat("a", 50_000) + "!"
	start := time.Now()
	ok := re.MatchString(input)
	elapsed := time.Since(start)
	fmt.Println("  (a+)+$ on 50000 'a' + '!': matched =", ok, ", under 1s =", elapsed < time.Second)
}

// 8. User input: QuoteMeta escapes metacharacters.
func quoteMeta() {
	fmt.Println("== 8. QuoteMeta for user input ==")
	user := "1.5+2"
	fmt.Println("  QuoteMeta:", regexp.QuoteMeta(user))
	raw := regexp.MustCompile(user)
	safe := regexp.MustCompile(regexp.QuoteMeta(user))
	fmt.Println("  raw pattern matches \"1x55552\": ", raw.MatchString("1x55552"))
	fmt.Println("  quoted pattern matches \"1x55552\":", safe.MatchString("1x55552"))
}

// 9. A *Regexp is safe for concurrent use.
func concurrent() {
	fmt.Println("== 9. concurrent use ==")
	var wg sync.WaitGroup
	counts := make([]int, 4)
	for i := range counts {
		wg.Go(func() {
			for range 1000 {
				counts[i] += len(wordRE.FindAllString("a b c", -1))
			}
		})
	}
	wg.Wait()
	fmt.Println("  words counted per goroutine:", counts)
}

func main() {
	basics()
	groups()
	replace()
	leftmost()
	flags()
	unsupported()
	linearTime()
	quoteMeta()
	concurrent()
}
