package swiss_table

import "math/bits"

const (
	groupSize = 8

	ctrlEmpty   uint8 = 0x80 // represent this slot not have been used.
	ctrlDeleted uint8 = 0xFE // represent this slot have been deleted.
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

// Table is a simplified educational Swiss Table.
//
// It demonstrates:
//   - open addressing
//   - control bytes
//   - fingerprints
//   - group probing
//   - tombstones
//
// It is intentionally simplified and is NOT the Go runtime map.
type Table[K comparable, V any] struct {
	ctrl    []uint8       // slot's metadata
	entries []entry[K, V] // key / value
	mask    uint64        // index sequence
	size    int           // how many data in table
}

// New creates a table with at least capacity slots.
//
// Capacity is rounded to a power of two.
func New[K comparable, V any](capacity int) *Table[K, V] {
	if capacity < groupSize {
		capacity = groupSize
	}

	capacity = nextPowerOfTwo(capacity)

	ctrl := make([]uint8, capacity)

	for i := range ctrl {
		ctrl[i] = ctrlEmpty
	}

	return &Table[K, V]{
		ctrl:    ctrl,
		entries: make([]entry[K, V], capacity),
		mask:    uint64(capacity - 1),
	}
}

func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}

	return 1 << bits.Len(uint(n-1))
}

// Len returns the number of live entries.
func (t *Table[K, V]) Len() int {
	return t.size
}

// Capacity returns the number of slots.
func (t *Table[K, V]) Capacity() int {
	return len(t.ctrl)
}

// LoadFactor returns entries / capacity.
func (t *Table[K, V]) LoadFactor() float64 {
	if len(t.ctrl) == 0 {
		return 0
	}

	return float64(t.size) / float64(len(t.ctrl))
}

// hashKey is a deterministic educational hash.
//
// The implementation intentionally uses a simple integer hash so that
// the mechanics remain visible.
//
// A production implementation would use an appropriate keyed hash.
func hashKey[K comparable](key K) uint64 {
	// We cannot generically convert arbitrary comparable values into a
	// stable integer without reflection.
	//
	// For this educational implementation we use the address-independent
	// textual representation through a type switch for common benchmark
	// types and fall back to a generic representation.
	//
	// The benchmarks use uint64, which takes the fast path.
	switch v := any(key).(type) {
	case uint64:
		return mix64(v)
	case uint32:
		return mix64(uint64(v))
	case uint:
		return mix64(uint64(v))
	case int:
		return mix64(uint64(v))
	case int64:
		return mix64(uint64(v))
	case int32:
		return mix64(uint64(v))
	case string:
		return hashString(v)
	default:
		return hashStringFallback(key)
	}
}

func mix64(x uint64) uint64 {
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9

	x ^= x >> 27
	x *= 0x94d049bb133111eb

	x ^= x >> 31

	return x
}

func hashString(s string) uint64 {
	var h uint64 = 14695981039346656037

	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}

	return mix64(h)
}

func hashStringFallback[K comparable](key K) uint64 {
	// This fallback is deliberately simple.
	//
	// For production use, use a proper hash function appropriate for
	// the key type.
	return hashString(stringify(key))
}

func stringify[K comparable](key K) string {
	return formatValue(key)
}

// formatValue avoids reflection-heavy code in the hot path for the
// benchmarked primitive types.
func formatValue[K comparable](key K) string {
	switch v := any(key).(type) {
	case string:
		return v
	case uint64:
		var buf [8]byte
		for i := range buf {
			buf[i] = byte(v >> (8 * i))
		}
		return string(buf[:])
	case uint:
		return formatValue(uint64(v))
	case uint32:
		return formatValue(uint64(v))
	case int:
		return formatValue(uint64(v))
	case int64:
		return formatValue(uint64(v))
	default:
		return ""
	}
}

// splitHash separates the hash into:
//
// H1 -> initial group/slot
// H2 -> small fingerprint
func splitHash(hash uint64, mask uint64) (uint64, uint8) {
	h1 := hash & mask

	// Keep H2 away from EMPTY and DELETED markers.
	h2 := uint8((hash >> 57) & 0x7F)

	if h2 == 0 {
		h2 = 1
	}

	return h1, h2
}

// Get returns a value associated with key.
func (t *Table[K, V]) Get(key K) (V, bool) {
	hash := hashKey(key)

	start, fingerprint := splitHash(hash, t.mask)

	capacity := uint64(len(t.ctrl))

	for probe := uint64(0); probe < capacity; probe += groupSize {
		groupStart := (start + probe) & t.mask

		for offset := uint64(0); offset < groupSize; offset++ {
			index := (groupStart + offset) & t.mask
			ctrl := t.ctrl[index]

			if ctrl == ctrlEmpty {
				var zero V
				return zero, false
			}

			if ctrl != fingerprint {
				continue
			}

			if t.entries[index].key == key {
				return t.entries[index].value, true
			}
		}
	}

	var zero V
	return zero, false
}

// Set inserts or updates a value.
func (t *Table[K, V]) Set(key K, value V) {
	// Keep the table away from complete saturation.
	if float64(t.size+1)/float64(len(t.ctrl)) > 0.75 {
		t.resize(len(t.ctrl) * 2)
	}

	hash := hashKey(key)

	start, fingerprint := splitHash(hash, t.mask)

	var firstDeleted uint64
	hasDeleted := false

	capacity := uint64(len(t.ctrl))

	for probe := uint64(0); probe < capacity; probe += groupSize {
		groupStart := (start + probe) & t.mask

		for offset := uint64(0); offset < groupSize; offset++ {
			index := (groupStart + offset) & t.mask
			ctrl := t.ctrl[index]

			switch ctrl {
			case ctrlEmpty:
				if hasDeleted {
					index = firstDeleted
				}

				t.ctrl[index] = fingerprint
				t.entries[index] = entry[K, V]{
					key:   key,
					value: value,
				}

				t.size++
				return

			case ctrlDeleted:
				if !hasDeleted {
					firstDeleted = index
					hasDeleted = true
				}

			default:
				if ctrl == fingerprint && t.entries[index].key == key {
					t.entries[index].value = value
					return
				}
			}
		}
	}

	// Defensive fallback. Normally resizing above prevents this.
	t.resize(len(t.ctrl) * 2)
	t.Set(key, value)
}

// Delete removes a key.
func (t *Table[K, V]) Delete(key K) bool {
	hash := hashKey(key)

	start, fingerprint := splitHash(hash, t.mask)

	capacity := uint64(len(t.ctrl))

	for probe := uint64(0); probe < capacity; probe += groupSize {
		groupStart := (start + probe) & t.mask

		for offset := uint64(0); offset < groupSize; offset++ {
			index := (groupStart + offset) & t.mask
			ctrl := t.ctrl[index]

			if ctrl == ctrlEmpty {
				return false
			}

			if ctrl != fingerprint {
				continue
			}

			if t.entries[index].key != key {
				continue
			}

			t.ctrl[index] = ctrlDeleted

			var zero entry[K, V]
			t.entries[index] = zero

			t.size--

			return true
		}
	}

	return false
}

func (t *Table[K, V]) resize(newCapacity int) {
	oldCtrl := t.ctrl
	oldEntries := t.entries

	nt := New[K, V](newCapacity)

	for i, ctrl := range oldCtrl {
		if ctrl == ctrlEmpty || ctrl == ctrlDeleted {
			continue
		}

		nt.Set(oldEntries[i].key, oldEntries[i].value)
	}

	*t = *nt
}
