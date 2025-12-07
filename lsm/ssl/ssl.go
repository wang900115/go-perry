package ssl

import (
	"encoding/binary"
	"fmt"
	"io"
	"lsm/mem"
	"os"
	"sort"
)

// Format:
// body records
// each record: keyLen(uint32) | valLen(uint32) | key bytes | val bytes
// index: numEntries(uint32) [ for each: keyLen(uint32) | key bytes | offset(uint64) ]
// trailer: indexOffset(uint64)

type SSTable struct {
	Path  string
	index map[string]uint64 // key to offset
}

func WriteSSTable(path string, kv map[string]mem.Value) error {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	offsets := make(map[string]uint64)
	// write body records and remember offsets
	for _, k := range keys {
		offset, _ := f.Seek(0, io.SeekCurrent)
		offsets[k] = uint64(offset)
		v := kv[k]
		if err := binary.Write(f, binary.LittleEndian, uint32(len(k))); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, uint32(len(v))); err != nil {
			return err
		}
		if _, err := f.Write([]byte(k)); err != nil {
			return err
		}
		if _, err := f.Write(v); err != nil {
			return err
		}
	}
	// write index
	indexOffset, _ := f.Seek(0, io.SeekCurrent)
	if err := binary.Write(f, binary.LittleEndian, uint32(len(keys))); err != nil {
		return err
	}
	for _, k := range keys {
		if err := binary.Write(f, binary.LittleEndian, uint32(len(k))); err != nil {
			return err
		}
		if _, err := f.Write([]byte(k)); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, offsets[k]); err != nil {
			return err
		}
	}
	// write trailer
	if err := binary.Write(f, binary.LittleEndian, uint64(indexOffset)); err != nil {
		return err
	}
	return nil
}

func LoadSSTable(path string) (*SSTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read trailer
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() < 8 {
		return nil, fmt.Errorf("file too small to be a valid SSTable")
	}
	if _, err := f.Seek(-8, io.SeekEnd); err != nil {
		return nil, err
	}
	var indexOffset uint64
	if err := binary.Read(f, binary.LittleEndian, &indexOffset); err != nil {
		return nil, err
	}
	if _, err := f.Seek(int64(indexOffset), io.SeekStart); err != nil {
		return nil, err
	}
	var n uint32
	if err := binary.Read(f, binary.LittleEndian, &n); err != nil {
		return nil, err
	}
	index := make(map[string]uint64)
	for i := uint32(0); i < n; i++ {
		var keyLen uint32
		if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
			return nil, err
		}
		kbuf := make([]byte, keyLen)
		if _, err := io.ReadFull(f, kbuf); err != nil {
			return nil, err
		}
		var offset uint64
		if err := binary.Read(f, binary.LittleEndian, &offset); err != nil {
			return nil, err
		}
		index[string(kbuf)] = offset
	}
	return &SSTable{Path: path, index: index}, nil
}

func (sst *SSTable) Get(key string) (mem.Value, bool, error) {
	offset, ok := sst.index[key]
	if !ok {
		return nil, false, nil
	}
	f, err := os.Open(sst.Path)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	if _, err := f.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, false, err
	}
	var keyLen uint32
	if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
		return nil, false, err
	}
	var valLen uint32
	if err := binary.Read(f, binary.LittleEndian, &valLen); err != nil {
		return nil, false, err
	}
	kbuf := make([]byte, keyLen)
	if _, err := io.ReadFull(f, kbuf); err != nil {
		return nil, false, err
	}
	if string(kbuf) != key {
		return nil, false, fmt.Errorf("key mismatch at offset")
	}
	vbuf := make([]byte, valLen)
	if _, err := io.ReadFull(f, vbuf); err != nil {
		return nil, false, err
	}
	return vbuf, true, nil
}
