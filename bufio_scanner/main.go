// Package main demonstrates bufio: Scanner (lines, words, custom SplitFunc,
// token size limits), Reader (ReadString, Peek), and Writer (buffering, Flush).
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// countingWriter counts how many Write calls reach the underlying writer.
type countingWriter struct {
	calls int
	buf   bytes.Buffer
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.calls++
	return c.buf.Write(p)
}

// chunkReader returns at most n bytes per Read, like a slow network or pipe.
type chunkReader struct {
	r io.Reader
	n int
}

func (c *chunkReader) Read(p []byte) (int, error) {
	if len(p) > c.n {
		p = p[:c.n]
	}
	return c.r.Read(p)
}

func part1Lines() {
	fmt.Println("== Part 1: Scanner lines (CRLF and missing final newline) ==")
	// Windows-style line endings, and the last line has no newline at all.
	input := "alpha\r\nbeta\r\n\r\ngamma"
	sc := bufio.NewScanner(strings.NewReader(input))
	n := 0
	for sc.Scan() {
		n++
		fmt.Printf("line %d: %q\n", n, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("lines:", n, "(\\r is dropped, empty line kept, last line returned)")
}

func part2Words() {
	fmt.Println("\n== Part 2: Scanner words and a 3-byte-per-Read source ==")
	src := &chunkReader{r: strings.NewReader("  the quick\tbrown\n\nfox  "), n: 3}
	sc := bufio.NewScanner(src)
	sc.Split(bufio.ScanWords)
	var words []string
	for sc.Scan() {
		words = append(words, sc.Text())
	}
	fmt.Printf("words: %q\n", words)
	fmt.Println("tokens are reassembled even though each Read returned at most 3 bytes")
}

// commaSplit is a custom SplitFunc: fields separated by ',' with spaces trimmed.
// A field equal to "END" stops scanning via bufio.ErrFinalToken.
func commaSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil // no more data: stop
	}
	if i := bytes.IndexByte(data, ','); i >= 0 {
		field := bytes.TrimSpace(data[:i])
		if string(field) == "END" {
			return i + 1, nil, bufio.ErrFinalToken // stop, no token
		}
		return i + 1, field, nil
	}
	if atEOF {
		return len(data), bytes.TrimSpace(data), nil // last field without ','
	}
	return 0, nil, nil // need more data: ask the Scanner to read more
}

func part3CustomSplit() {
	fmt.Println("\n== Part 3: custom SplitFunc + ErrFinalToken ==")
	sc := bufio.NewScanner(strings.NewReader("red, green ,blue,END,ignored,also-ignored"))
	sc.Split(commaSplit)
	for sc.Scan() {
		fmt.Printf("field: %q\n", sc.Text())
	}
	fmt.Println("err after ErrFinalToken:", sc.Err())

	sc = bufio.NewScanner(strings.NewReader("a,b,  last  "))
	sc.Split(commaSplit)
	var got []string
	for sc.Scan() {
		got = append(got, sc.Text())
	}
	fmt.Printf("no END: %q err=%v\n", got, sc.Err())
}

func part4TooLong() {
	fmt.Println("\n== Part 4: token too long and Scanner.Buffer ==")
	long := strings.Repeat("x", 70*1024) // 70 KiB line > default 64 KiB limit
	input := "short\n" + long + "\nafter\n"

	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		fmt.Printf("default: got line of %d bytes\n", len(sc.Bytes()))
	}
	err := sc.Err()
	fmt.Println("default: err =", err, "| is ErrTooLong:", errors.Is(err, bufio.ErrTooLong))
	fmt.Println("default MaxScanTokenSize =", bufio.MaxScanTokenSize)

	sc = bufio.NewScanner(strings.NewReader(input))
	sc.Buffer(make([]byte, 0, 4096), 1024*1024) // start at 4 KiB, allow growth to 1 MiB
	for sc.Scan() {
		fmt.Printf("1MiB max: got line of %d bytes\n", len(sc.Bytes()))
	}
	fmt.Println("1MiB max: err =", sc.Err())

	// Bytes() aliases the internal buffer: it is overwritten by the next Scan.
	// A tiny 8-byte buffer forces the Scanner to move data around.
	sc = bufio.NewScanner(strings.NewReader("one\ntwo\nsix\n"))
	sc.Buffer(make([]byte, 0, 8), 8)
	var kept, copied [][]byte
	for sc.Scan() {
		kept = append(kept, sc.Bytes())                  // BUG: aliases the buffer
		copied = append(copied, bytes.Clone(sc.Bytes())) // safe copy
	}
	fmt.Printf("kept without copy: %q\n", kept)
	fmt.Printf("kept with copy:    %q\n", copied)
}

func part5Reader() {
	fmt.Println("\n== Part 5: bufio.Reader ReadString and Peek ==")
	r := bufio.NewReader(strings.NewReader("HTTP/1.1 200 OK\nlast line no newline"))
	head, _ := r.Peek(4) // look ahead without consuming
	fmt.Printf("peek: %q buffered=%d\n", head, r.Buffered())
	for {
		line, err := r.ReadString('\n')
		fmt.Printf("ReadString: %q err=%v\n", line, err)
		if err != nil {
			break // io.EOF still came with the final data
		}
	}
}

func part6Writer() {
	fmt.Println("\n== Part 6: bufio.Writer buffering and Flush ==")
	cw := &countingWriter{}
	w := bufio.NewWriterSize(cw, 16)
	for i := 0; i < 5; i++ {
		fmt.Fprintf(w, "row%d\n", i) // 5 bytes each
	}
	fmt.Printf("before Flush: underlying calls=%d, buffered=%d, written=%q\n",
		cw.calls, w.Buffered(), cw.buf.String())
	if err := w.Flush(); err != nil {
		fmt.Println("flush error:", err)
	}
	fmt.Printf("after Flush:  underlying calls=%d, buffered=%d, total bytes=%d\n",
		cw.calls, w.Buffered(), cw.buf.Len())

	// A write larger than the buffer, while the buffer is empty, goes straight through.
	cw2 := &countingWriter{}
	w2 := bufio.NewWriterSize(cw2, 16)
	w2.Write(bytes.Repeat([]byte("z"), 40))
	fmt.Printf("large write on empty buffer: calls=%d buffered=%d\n", cw2.calls, w2.Buffered())

	// Forgetting Flush loses data.
	cw3 := &countingWriter{}
	w3 := bufio.NewWriter(cw3)
	fmt.Fprintln(w3, "this line is never flushed")
	fmt.Printf("no Flush: underlying got %d bytes (data still in buffer: %d)\n",
		cw3.buf.Len(), w3.Buffered())
}

func main() {
	part1Lines()
	part2Words()
	part3CustomSplit()
	part4TooLong()
	part5Reader()
	part6Writer()
}
