// Package main demonstrates how time.Time values really behave:
// layout strings, parsing and locations, the monotonic clock reading,
// comparing instants, calendar arithmetic, and daylight saving time.
package main

import (
	"fmt"
	"strings"
	"time"

	// Embed the IANA time zone database into the binary (~450 KB).
	// Without it, LoadLocation depends on the OS / GOROOT zoneinfo,
	// which may be missing on Windows machines without a Go install.
	_ "time/tzdata"
)

func main() {
	layouts()
	parsing()
	monotonic()
	comparing()
	calendarMath()
	daylightSaving()
}

// 1. Layouts are written with the reference time Mon Jan 2 15:04:05 MST 2006.
func layouts() {
	fmt.Println("== 1. layouts ==")
	t := time.Date(2026, time.November, 8, 9, 5, 7, 0, time.UTC)

	fmt.Println(t.Format("2006-01-02 15:04:05")) // correct
	fmt.Println(t.Format(time.DateTime))         // same thing (Go 1.20+)
	fmt.Println(t.Format(time.RFC3339))
	fmt.Println(t.Format("Jan 2, 2006 at 3:04pm (MST)"))

	// Pitfalls: these "look right" but are not layouts.
	fmt.Println(t.Format("YYYY-MM-DD")) // no magic tokens: printed literally
	// Swapped month/day. `go vet` flags this when the layout is a literal,
	// but not when it comes from a variable or config file.
	swapped := "2006-02-01"
	fmt.Println(t.Format(swapped))            // 2026-08-11
	fmt.Println(t.Format("2006-01-02 03:04")) // 03 = 12-hour clock
	fmt.Println()
}

// 2. Parse assumes UTC when the input has no zone; ParseInLocation does not.
func parsing() {
	fmt.Println("== 2. parsing ==")
	taipei, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	const s = "2026-11-28 09:00:00"

	a, _ := time.Parse(time.DateTime, s)
	b, _ := time.ParseInLocation(time.DateTime, s, taipei)
	fmt.Println("Parse:          ", a)
	fmt.Println("ParseInLocation:", b)
	fmt.Println("same instant?   ", a.Equal(b), "difference:", a.Sub(b))

	_, err = time.Parse(time.DateOnly, "2026-02-30")
	fmt.Println("invalid date:   ", err)
	fmt.Println()
}

// 3. time.Now carries a monotonic clock reading used by Sub/Since/Before/After.
func monotonic() {
	fmt.Println("== 3. monotonic clock ==")
	start := time.Now()
	stripped := start.Round(0) // Round(0) removes the monotonic reading

	// The String form contains "m=+..." only when a monotonic reading exists.
	hasM := func(t time.Time) bool { return strings.Contains(t.String(), " m=") }
	fmt.Println("Now has monotonic reading:     ", hasM(start))
	fmt.Println("Round(0) has monotonic reading:", hasM(stripped))
	fmt.Println("start == stripped:             ", start == stripped)
	fmt.Println("start.Equal(stripped):         ", start.Equal(stripped))

	// Elapsed time uses the monotonic clock, so it never goes negative
	// even if the wall clock is adjusted (NTP, user changes the clock).
	time.Sleep(20 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Println("elapsed >= 20ms:               ", elapsed >= 20*time.Millisecond)
	fmt.Println()
}

// 4. Use Equal, not ==, and never use time.Time as a map key without normalizing.
func comparing() {
	fmt.Println("== 4. comparing ==")
	utc := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	taipei := utc.In(time.FixedZone("UTC+8", 8*60*60))

	fmt.Println(utc, "|", taipei)
	fmt.Println("==    :", utc == taipei)
	fmt.Println("Equal :", utc.Equal(taipei))

	seen := map[time.Time]bool{utc: true}
	fmt.Println("map lookup with other location:", seen[taipei])
	fmt.Println("map lookup after .UTC():       ", seen[taipei.UTC()])

	var zero time.Time
	fmt.Println("zero value:", zero, "IsZero:", zero.IsZero())
	fmt.Println()
}

// 5. AddDate normalizes overflow; it does not clamp to the end of month.
func calendarMath() {
	fmt.Println("== 5. calendar math ==")
	jan31 := time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC)
	fmt.Println("Jan 31 + 1 month:", jan31.AddDate(0, 1, 0).Format(time.DateOnly))

	// End of next month, the safe way: first day of the month after, minus one day.
	firstOfMonthAfter := time.Date(jan31.Year(), jan31.Month()+2, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("end of next month:", firstOfMonthAfter.AddDate(0, 0, -1).Format(time.DateOnly))

	// Date() also normalizes: month 13 / day 0 are legal inputs.
	fmt.Println("Date(2026, 13, 0):", time.Date(2026, 13, 0, 0, 0, 0, 0, time.UTC).Format(time.DateOnly))

	d := 90*time.Minute + 30*time.Second + 250*time.Millisecond
	fmt.Println("duration:", d, "| hours:", d.Hours(), "| truncated:", d.Truncate(time.Minute))
	fmt.Println()
}

// 6. A "day" is not always 24 hours.
func daylightSaving() {
	fmt.Println("== 6. daylight saving time ==")
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	// DST starts 2026-03-08 at 02:00 in New York (clocks jump to 03:00).
	before := time.Date(2026, time.March, 7, 12, 0, 0, 0, ny)
	fmt.Println("start:           ", before)
	fmt.Println("Add(24h):        ", before.Add(24*time.Hour))
	fmt.Println("AddDate(0, 0, 1):", before.AddDate(0, 0, 1))

	// 02:30 does not exist that day; Date picks a valid time instead of failing.
	gap := time.Date(2026, time.March, 8, 2, 30, 0, 0, ny)
	fmt.Println("02:30 in the gap:", gap)
}
