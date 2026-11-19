// Command string_utf8 shows how Go strings, bytes, runes and UTF-8 fit together.
package main

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func main() {
	part1()
	part2()
	part3()
	part4()
	part5()
	part6()
}

// Part 1: a string is a read-only sequence of bytes; len counts bytes.
func part1() {
	fmt.Println("== Part 1: bytes vs runes ==")
	s := "Go語言"
	fmt.Printf("s=%q len=%d runes=%d\n", s, len(s), utf8.RuneCountInString(s))
	fmt.Printf("s[2]=%d (%#x) is a byte, not a character\n", s[2], s[2])
	fmt.Printf("bytes: % x\n", s)
	for i, r := range s { // range decodes UTF-8: i is a byte offset
		fmt.Printf("  offset=%d rune=%q U+%04X width=%d\n", i, r, r, utf8.RuneLen(r))
	}
	fmt.Println()
}

// Part 2: strings may hold invalid UTF-8; range yields U+FFFD for bad bytes.
func part2() {
	fmt.Println("== Part 2: invalid UTF-8 ==")
	bad := "a\xffb"
	fmt.Printf("bad=%q len=%d valid=%v\n", bad, len(bad), utf8.ValidString(bad))
	for i, r := range bad {
		fmt.Printf("  offset=%d rune=%q\n", i, r)
	}
	r, size := utf8.DecodeRuneInString(bad[1:])
	fmt.Printf("DecodeRuneInString: rune=%q size=%d (RuneError=%v)\n", r, size, r == utf8.RuneError)
	fmt.Printf("ToValidUTF8: %q\n", strings.ToValidUTF8(bad, "?"))
	fmt.Println()
}

// Part 3: conversions and what they cost.
func part3() {
	fmt.Println("== Part 3: conversions ==")
	s := "héllo"
	b := []byte(s) // copies: strings are immutable, []byte is not
	b[0] = 'H'
	fmt.Printf("s=%q b=%q (the string did not change)\n", s, b)
	rs := []rune(s) // decodes: one int32 per code point
	fmt.Printf("len([]byte)=%d len([]rune)=%d\n", len(b), len(rs))
	rs[1] = 'e'
	fmt.Printf("string(rs)=%q\n", string(rs))
	fmt.Printf("string(rune(65))=%q  string(rune(0x8A9E))=%q\n", string(rune(65)), string(rune(0x8A9E)))
	fmt.Printf("string(rune(-1))=%q (invalid code point becomes U+FFFD)\n", string(rune(-1)))
	fmt.Println()
}

// Part 4: slicing a string slices bytes, so it can cut a rune in half.
func part4() {
	fmt.Println("== Part 4: safe truncation ==")
	s := "Go語言很好"
	naive := s[:4]
	fmt.Printf("naive s[:4]=%q valid=%v\n", naive, utf8.ValidString(naive))
	fmt.Printf("truncateRunes(s, 3)=%q\n", truncateRunes(s, 3))
	fmt.Printf("truncateBytes(s, 4)=%q\n", truncateBytes(s, 4))
	fmt.Printf("truncateBytes(s, 8)=%q\n", truncateBytes(s, 8))
	fmt.Println()
}

// truncateRunes keeps at most n runes.
func truncateRunes(s string, n int) string {
	count := 0
	for i := range s {
		if count == n {
			return s[:i]
		}
		count++
	}
	return s
}

// truncateBytes keeps at most max bytes without splitting a rune.
func truncateBytes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	for max > 0 && !utf8.RuneStart(s[max]) {
		max--
	}
	return s[:max]
}

// Part 5: building strings. += copies every time; strings.Builder does not.
func part5() {
	fmt.Println("== Part 5: building strings ==")
	words := strings.Fields("the quick brown fox jumps over the lazy dog")

	concat := testing.AllocsPerRun(100, func() {
		out := ""
		for _, w := range words {
			out += w + " "
		}
		_ = out
	})
	builder := testing.AllocsPerRun(100, func() {
		var sb strings.Builder
		sb.Grow(64)
		for _, w := range words {
			sb.WriteString(w)
			sb.WriteByte(' ')
		}
		_ = sb.String()
	})
	fmt.Printf("words=%d allocs with +=: %.0f, with Builder+Grow: %.0f\n", len(words), concat, builder)

	var sb strings.Builder
	sb.WriteString("Go")
	sb.WriteRune('語')
	fmt.Fprintf(&sb, "%d", 1)
	fmt.Printf("builder=%q len=%d\n", sb.String(), sb.Len())
	fmt.Println()
}

// Part 6: one visible character can be several runes.
func part6() {
	fmt.Println("== Part 6: a rune is not a visible character ==")
	precomposed := "é" // é as one code point
	combining := "é"  // e + COMBINING ACUTE ACCENT
	fmt.Printf("%q vs %q equal=%v runes=%d,%d bytes=%d,%d\n",
		precomposed, combining, precomposed == combining,
		utf8.RuneCountInString(precomposed), utf8.RuneCountInString(combining),
		len(precomposed), len(combining))
	flag := "\U0001F1F9\U0001F1FC" // regional indicators T + W: one flag on screen
	fmt.Printf("flag runes=%d bytes=%d\n", utf8.RuneCountInString(flag), len(flag))
	fmt.Printf("EqualFold(\"Go\", \"GO\")=%v\n", strings.EqualFold("Go", "GO"))
}
