package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
)

type BoomFilter struct {
	bitset []byte
	m      uint // bit array size
	k      uint // number of hash functions
}

// Optimal m = -(n * ln(p)) / (ln 2)^2
func optimalM(n uint, p float64) uint {
	return uint(math.Ceil(-(float64(n) * math.Log(p)) / (math.Ln2 * math.Ln2)))
}

// Optimal k = (m/n) * ln(2)
func optimalK(m, n uint) uint {
	return uint(math.Ceil((float64(m) / float64(n)) * math.Ln2))
}

func NewBloom(n uint, p float64) *BoomFilter {
	m := optimalM(n, p)
	k := optimalK(m, n)

	bits := make([]byte, (m+7)/8)

	return &BoomFilter{
		bitset: bits,
		m:      m,
		k:      k,
	}
}

func (b *BoomFilter) Add(item []byte) {
	for i := uint(0); i < b.k; i++ {
		pos := b.hash(item, i) % b.m
		b.setBit(pos)
	}
}

func (b *BoomFilter) Exists(item []byte) bool {
	for i := uint(0); i < b.k; i++ {
		pos := b.hash(item, i) % b.m
		if !b.getBit(pos) {
			return false
		}
	}
	return true
}

func (b *BoomFilter) setBit(pos uint) {
	byteIndex := pos / 8
	bitIndex := pos % 8
	b.bitset[byteIndex] |= 1 << bitIndex
}

func (b *BoomFilter) getBit(pos uint) bool {
	byteIndex := pos / 8
	bitIndex := pos % 8
	return (b.bitset[byteIndex] & (1 << bitIndex)) != 0
}

func (b *BoomFilter) hash(item []byte, seed uint) uint {
	h := fnv.New64a()
	_, _ = h.Write(item)

	// add extra seed
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, uint64(seed))
	_, _ = h.Write(seedBytes)
	return uint(h.Sum64())
}

func main() {
	// 預期 100000 筆資料，false positive = 1%
	bf := NewBloom(10000, 0.01)

	bf.Add([]byte("alice"))
	bf.Add([]byte("bob"))

	fmt.Println(bf.Exists([]byte("alice"))) // true
	fmt.Println(bf.Exists([]byte("bob")))   // true
	fmt.Println(bf.Exists([]byte("carol"))) // false (or maybe true: false positive)
}
