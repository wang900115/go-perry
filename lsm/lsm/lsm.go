package lsm

import (
	"encoding/binary"
	"fmt"
	"io"
	"lsm/mem"
	"lsm/ssl"
	"lsm/wal"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type LSM struct {
	dir      string
	mem      *mem.MemTable
	wal      *wal.WAL
	mu       sync.RWMutex
	sstFiles []string
}

func OpenLSM(dir string) (*LSM, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	mt := mem.NewMemTable()
	w, err := wal.OpenWAL(dir)
	if err != nil {
		return nil, err
	}
	// replay
	if err := w.ReplayInto(mt); err != nil {
		return nil, err
	}
	lsm := &LSM{
		dir: dir,
		mem: mt,
		wal: w,
	}
	// load sst files
	files, _ := filepath.Glob(filepath.Join(dir, "sst-*.db"))
	// sort files by higher timestamp
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})
	lsm.sstFiles = files
	return lsm, nil
}

func (lsm *LSM) Put(key string, value mem.Value) error {
	if err := lsm.wal.Append(key, value); err != nil {
		return err
	}
	lsm.mem.Put(key, value)
	return nil
}

func (lsm *LSM) Get(key string) (mem.Value, bool, error) {
	// 1. memtable
	if v, ok := lsm.mem.Get(key); ok {
		return v, true, nil
	}
	// 2. sst files (from newest to oldest)
	lsm.mu.RLock()
	files := append([]string{}, lsm.sstFiles...)
	lsm.mu.RUnlock()
	for _, f := range files {
		sst, err := ssl.LoadSSTable(f)
		if err != nil {
			return nil, false, err
		}
		v, ok, err := sst.Get(key)
		if err != nil {
			return nil, false, err
		}
		if ok {
			return v, true, nil
		}
	}
	return nil, false, nil
}

func (lsm *LSM) FlushMemToSST() error {
	kv := lsm.mem.Flush()
	if len(kv) == 0 {
		return nil
	}
	// write sst
	name := filepath.Join(lsm.dir, fmt.Sprintf("sst-%d.db", time.Now().UnixNano()))
	if err := ssl.WriteSSTable(name, kv); err != nil {
		return err
	}
	// truncate wal
	if err := lsm.wal.Truncate(); err != nil {
		return err
	}
	// prepend sst file as newest
	lsm.mu.Lock()
	lsm.sstFiles = append([]string{name}, lsm.sstFiles...)
	lsm.mu.Unlock()
	return nil
}

func (lsm *LSM) CompactTwoNewest() error {
	lsm.mu.Lock()
	if len(lsm.sstFiles) < 2 {
		lsm.mu.Unlock()
		return nil
	}
	a := lsm.sstFiles[0]
	b := lsm.sstFiles[1]

	lsm.sstFiles = lsm.sstFiles[2:]
	lsm.mu.Unlock()

	sa, err := loadKVFromSST(a)
	if err != nil {
		return err
	}
	sb, err := loadKVFromSST(b)
	if err != nil {
		return err
	}
	merged := make(map[string]mem.Value)
	for k, v := range sa {
		merged[k] = v
	}
	for k, v := range sb {
		merged[k] = v
	}
	outName := filepath.Join(lsm.dir, fmt.Sprintf("sst-%d.db", time.Now().UnixNano()))
	if err := ssl.WriteSSTable(outName, merged); err != nil {
		return err
	}
	_ = os.Remove(a)
	_ = os.Remove(b)
	lsm.mu.Lock()
	lsm.sstFiles = append([]string{outName}, lsm.sstFiles...)
	lsm.mu.Unlock()
	return nil
}

func loadKVFromSST(path string) (map[string]mem.Value, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, _ := f.Stat()
	if fi.Size() < 8 {
		return nil, fmt.Errorf("invalid sst file")
	}
	if _, err := f.Seek(-8, os.SEEK_END); err != nil {
		return nil, err
	}
	var indexOffset int64
	if err := binary.Read(f, binary.LittleEndian, &indexOffset); err != nil {
		return nil, err
	}
	if _, err := f.Seek(indexOffset, io.SeekStart); err != nil {
		return nil, err
	}
	var n uint32
	if err := binary.Read(f, binary.LittleEndian, &n); err != nil {
		return nil, err
	}

	index := make([]struct {
		key string
		off uint64
	}, 0, n)
	for i := uint32(0); i < n; i++ {
		var keyLen uint32
		if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
			return nil, err
		}
		keyBytes := make([]byte, keyLen)
		if _, err := io.ReadFull(f, keyBytes); err != nil {
			return nil, err
		}
		var offset uint64
		if err := binary.Read(f, binary.LittleEndian, &offset); err != nil {
			return nil, err
		}
		index = append(index, struct {
			key string
			off uint64
		}{key: string(keyBytes), off: offset})
	}

	out := make(map[string]mem.Value, len(index))
	for _, entry := range index {
		if _, err := f.Seek(int64(entry.off), io.SeekStart); err != nil {
			return nil, err
		}
		var keyLen uint32
		if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
			return nil, err
		}
		var valLen uint32
		if err := binary.Read(f, binary.LittleEndian, &valLen); err != nil {
			return nil, err
		}
		keyBuf := make([]byte, keyLen)
		if _, err := io.ReadFull(f, keyBuf); err != nil {
			return nil, err
		}
		valBuf := make([]byte, valLen)
		if _, err := io.ReadFull(f, valBuf); err != nil {
			return nil, err
		}
		out[string(keyBuf)] = valBuf
	}
	return out, nil
}
