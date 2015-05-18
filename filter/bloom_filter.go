package filter

import (
	"log"
	"sync"

	"github.com/willf/bloom"
)

type bloomFilter struct {
	rwmu   sync.RWMutex
	bloomf *bloom.BloomFilter
}

func NewBloomFilter() MutableExistenceFilter {
	return new(bloomFilter)
}

// keyCount should be as close as possible to the real number of keys
func (bf *bloomFilter) ImportKeys(keyCount uint, keys <-chan string) {
	// Load the new keys into a Bloom filter
	log.Println("creating a bloom filter with", keyCount, "keys")
	bloomf := bloom.NewWithEstimates(keyCount, 0.001)
	for key := range keys {
		log.Println("adding key", key, "to bloom filter")
		bloomf.AddString(key)
	}

	// Swap the new Bloom filter with the old one
	bf.rwmu.Lock()
	defer bf.rwmu.Unlock()

	bf.bloomf = bloomf
}

func (bf *bloomFilter) KeyExists(key string) bool {
	bf.rwmu.RLock()
	defer bf.rwmu.RUnlock()

	if bf.bloomf != nil {
		return bf.bloomf.TestString(key)
	} else {
		return false
	}
}
