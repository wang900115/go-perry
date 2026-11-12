// Package main demonstrates the io.Reader / io.Writer contracts and how small
// readers and writers compose into pipelines with the standard library.
package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ---------- Part 1: the Read contract ----------

// eofWithDataReader returns its last bytes together with io.EOF in the same
// call. This is legal: Read may return n > 0 and err != nil at the same time.
type eofWithDataReader struct {
	data []byte
	pos  int
}

func (r *eofWithDataReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	if r.pos >= len(r.data) {
		return n, io.EOF // data AND EOF together
	}
	return n, nil
}

// wrongLoop drops the bytes that arrive together with io.EOF.
func wrongLoop(r io.Reader) string {
	var sb strings.Builder
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		if err != nil {
			break // BUG: buf[:n] is lost
		}
		sb.Write(buf[:n])
	}
	return sb.String()
}

// rightLoop always consumes buf[:n] first, then looks at err.
func rightLoop(r io.Reader) (string, error) {
	var sb strings.Builder
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		sb.Write(buf[:n])
		if errors.Is(err, io.EOF) {
			return sb.String(), nil
		}
		if err != nil {
			return sb.String(), err
		}
	}
}

// ---------- Part 2: writing your own Reader and Writer ----------

// upperReader wraps another Reader and upper-cases ASCII letters on the fly.
type upperReader struct{ src io.Reader }

func (u upperReader) Read(p []byte) (int, error) {
	n, err := u.src.Read(p)
	for i := range p[:n] {
		if 'a' <= p[i] && p[i] <= 'z' {
			p[i] -= 'a' - 'A'
		}
	}
	return n, err
}

// countingWriter counts bytes and Write calls, then forwards to dst.
type countingWriter struct {
	dst    io.Writer
	bytes  int64
	writes int
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.dst.Write(p)
	c.bytes += int64(n)
	c.writes++
	return n, err
}

// ---------- Part 3: composition ----------

// onlyReader hides every method except Read, so io.Copy cannot use the
// WriterTo / ReaderFrom fast paths of the wrapped value.
type onlyReader struct{ io.Reader }

func main() {
	fmt.Println("== Part 1: the Read contract ==")
	fmt.Printf("wrong loop: %q\n", wrongLoop(&eofWithDataReader{data: []byte("hello, gopher")}))
	s, err := rightLoop(&eofWithDataReader{data: []byte("hello, gopher")})
	fmt.Printf("right loop: %q err=%v\n", s, err)
	all, err := io.ReadAll(&eofWithDataReader{data: []byte("hello, gopher")})
	fmt.Printf("io.ReadAll: %q err=%v (ReadAll never returns io.EOF)\n", all, err)

	fmt.Println("== Part 2: custom Reader and Writer ==")
	cw := &countingWriter{dst: os.Stdout}
	n, err := io.Copy(cw, onlyReader{upperReader{strings.NewReader("small interfaces compose\n")}})
	fmt.Printf("copied %d bytes, err=%v, cw.bytes=%d\n", n, err, cw.bytes)

	fmt.Println("== Part 3: TeeReader + MultiWriter + LimitReader ==")
	src := strings.NewReader("line one\nline two\nline three\n")
	hash := sha256.New() // hash.Hash is an io.Writer
	var copyBuf strings.Builder
	tee := io.TeeReader(src, hash) // everything read from tee is also written to hash
	n, err = io.Copy(io.MultiWriter(&copyBuf, io.Discard), tee)
	fmt.Printf("copied %d bytes, err=%v\n", n, err)
	fmt.Printf("sha256 = %x...\n", hash.Sum(nil)[:8])
	limited, _ := io.ReadAll(io.LimitReader(strings.NewReader(copyBuf.String()), 8))
	fmt.Printf("first 8 bytes via LimitReader: %q\n", limited)

	fmt.Println("== Part 4: io.Copy fast paths ==")
	// strings.Reader implements io.WriterTo, so io.Copy calls WriteTo once.
	cw1 := &countingWriter{dst: io.Discard}
	io.Copy(cw1, strings.NewReader(strings.Repeat("x", 100_000)))
	// Hidden behind onlyReader, io.Copy falls back to its 32 KiB buffer loop.
	cw2 := &countingWriter{dst: io.Discard}
	io.Copy(cw2, onlyReader{strings.NewReader(strings.Repeat("x", 100_000))})
	fmt.Printf("with WriterTo:    %d bytes in %d Write call(s)\n", cw1.bytes, cw1.writes)
	fmt.Printf("without WriterTo: %d bytes in %d Write call(s)\n", cw2.bytes, cw2.writes)

	fmt.Println("== Part 5: io.Pipe connects a writer goroutine to a reader ==")
	pr, pw := io.Pipe()
	go func() {
		bw := bufio.NewWriter(pw)
		for i := 1; i <= 3; i++ {
			fmt.Fprintf(bw, "event %d\n", i)
		}
		// Flush before closing, or buffered bytes are lost.
		if err := bw.Flush(); err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.CloseWithError(errors.New("producer finished early")) // reader sees this error
	}()
	sc := bufio.NewScanner(pr)
	for sc.Scan() {
		fmt.Println("got:", sc.Text())
	}
	fmt.Println("scanner err:", sc.Err())
}
