// Package main demonstrates encoding/binary: byte order, fixed-size
// struct encoding, varints with zigzag, and length-prefixed framing.
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

// Header is a fixed-size wire header: every field has a fixed size,
// so encoding/binary can encode it with reflection.
type Header struct {
	Magic   [4]byte
	Version uint8
	Flags   uint8
	Length  uint32 // in memory the compiler inserts 2 bytes of padding before this field
	Seq     uint64
}

// BadHeader uses int, whose size depends on the platform, so it is not fixed-size.
type BadHeader struct {
	ID int
}

func section(title string) { fmt.Printf("\n== %s ==\n", title) }

func main() {
	// 1. Byte order: the same number, two different byte layouts.
	section("1. byte order")
	v := uint32(0x01020304)
	be := binary.BigEndian.AppendUint32(nil, v)
	le := binary.LittleEndian.AppendUint32(nil, v)
	fmt.Printf("BigEndian    % x\n", be)
	fmt.Printf("LittleEndian % x\n", le)
	fmt.Printf("read LE bytes as BE -> 0x%08x (wrong!)\n", binary.BigEndian.Uint32(le))

	// 2. Fixed-size structs: binary.Size ignores Go's memory padding.
	section("2. fixed-size struct")
	h := Header{Magic: [4]byte{'P', 'R', 'R', 'Y'}, Version: 1, Flags: 0x80, Length: 5, Seq: 42}
	fmt.Printf("unsafe.Sizeof=%d binary.Size=%d\n", unsafe.Sizeof(h), binary.Size(h))
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, h); err != nil {
		panic(err)
	}
	fmt.Printf("wire bytes   % x\n", buf.Bytes())
	var got Header
	if err := binary.Read(&buf, binary.BigEndian, &got); err != nil {
		panic(err)
	}
	fmt.Printf("decoded      %+v equal=%v\n", got, got == h)

	// int is not fixed-size: Write returns an error instead of guessing.
	err := binary.Write(io.Discard, binary.BigEndian, BadHeader{ID: 1})
	fmt.Println("int field    ->", err)
	fmt.Println("Size(int)    ->", binary.Size(BadHeader{}))

	// 3. Go 1.23 Encode / Decode / Append work on []byte directly.
	section("3. Encode/Decode/Append (Go 1.23+)")
	small := make([]byte, 10)
	_, err = binary.Encode(small, binary.LittleEndian, h)
	fmt.Println("Encode into 10 bytes ->", err)
	out, err := binary.Append([]byte{0xAA}, binary.LittleEndian, h)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Append len=%d first=% x\n", len(out), out[:5])
	var h2 Header
	n, err := binary.Decode(out[1:], binary.LittleEndian, &h2)
	fmt.Printf("Decode n=%d err=%v seq=%d\n", n, err, h2.Seq)

	// 4. Varints: small numbers use fewer bytes.
	section("4. uvarint sizes")
	for _, x := range []uint64{1, 127, 128, 300, 1 << 32} {
		b := binary.AppendUvarint(nil, x)
		fmt.Printf("%-10d -> %d byte(s) % x\n", x, len(b), b)
	}

	// Signed varints use zigzag: 0,-1,1,-2,2 -> 0,1,2,3,4.
	section("5. zigzag (Varint) vs casting to uint64")
	for _, x := range []int64{0, -1, 1, -2, 2, -64, 64} {
		zz := binary.AppendVarint(nil, x)
		fmt.Printf("%4d -> Varint % -6x (%d B) | Uvarint(uint64) %2d B\n",
			x, zz, len(zz), len(binary.AppendUvarint(nil, uint64(x))))
	}

	// 6. Uvarint's n tells you why decoding failed.
	section("6. Uvarint error signalling")
	full := binary.AppendUvarint(nil, 300) // ac 02
	for _, in := range [][]byte{full, full[:1], {}, bytes.Repeat([]byte{0xff}, 11)} {
		x, n := binary.Uvarint(in)
		fmt.Printf("in=% -33x -> x=%d n=%d (%s)\n", in, x, n, explain(n))
	}

	// 7. Length-prefixed framing: uvarint length + payload.
	section("7. length-prefixed frames")
	var stream []byte
	for _, msg := range []string{"hi", "", "gopher"} {
		stream = binary.AppendUvarint(stream, uint64(len(msg)))
		stream = append(stream, msg...)
	}
	fmt.Printf("stream % x\n", stream)
	r := bytes.NewReader(stream)
	for {
		msg, err := readFrame(r, 1<<20)
		if errors.Is(err, io.EOF) {
			fmt.Println("done (clean EOF)")
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("frame %q\n", msg)
	}
	_, err = readFrame(bytes.NewReader([]byte{0x05, 'a', 'b'}), 1<<20)
	fmt.Println("truncated frame ->", err)
}

func explain(n int) string {
	switch {
	case n > 0:
		return "ok"
	case n == 0:
		return "buffer too small"
	default:
		return "overflow"
	}
}

// readFrame reads one uvarint-length-prefixed message.
// maxLen protects against a hostile length allocating huge memory.
func readFrame(r *bytes.Reader, maxLen uint64) ([]byte, error) {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err // io.EOF only if no byte of the length was read
	}
	if size > maxLen {
		return nil, fmt.Errorf("frame too large: %d", size)
	}
	msg := make([]byte, size)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, fmt.Errorf("read payload: %w", err)
	}
	return msg, nil
}
