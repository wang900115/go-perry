package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type WAL struct {
	f   *os.File
	mu  sync.Mutex
	dir string
}

func OpenWAL(dir string) (*WAL, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	fpath := dir + "/wal.log"
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{
		f:   f,
		dir: dir,
	}, nil
}

func (w *WAL) Append(key string, val Value) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(key))); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(val))); err != nil {
		return err
	}
	buf.WriteString(key)
	buf.Write(val)
	_, err := w.f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return w.f.Sync()
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.f.Close()
}

func (w *WAL) Path() string { return filepath.Join(w.dir, "wal.log") }

func (w *WAL) ReplayInto(mem *MemTable) error {
	f, err := os.Open(w.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	for {
		var klen uint32
		if err := binary.Read(f, binary.LittleEndian, &klen); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		var vlen uint32
		if err := binary.Read(f, binary.LittleEndian, &vlen); err != nil {
			return err
		}
		kbuf := make([]byte, klen)
		if _, err := io.ReadFull(f, kbuf); err != nil {
			return err
		}
		vbuf := make([]byte, vlen)
		if _, err := io.ReadFull(f, vbuf); err != nil {
			return err
		}
		mem.Put(string(kbuf), vbuf)
	}
}

func (w *WAL) Truncate() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.f.Close(); err != nil {
		return err
	}
	path := w.Path()
	if err := os.Remove(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.f = f
	return nil
}
