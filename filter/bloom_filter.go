package filter

import (
	"log"
	"sync"

	"github.com/willf/bloom"
)

type bloomFilter struct {
	rwmu        sync.RWMutex
	bloomf      *bloom.BloomFilter
	filterCheck FilterCheck
}

type FilterCheck func(string) bool

func NewBloomFilter(filterCheck FilterCheck) MutableExistenceFilter {
	filter := new(bloomFilter)
	filter.filterCheck = filterCheck
	return filter
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

func (bf *bloomFilter) keyIsInBloomFilter(key string) bool {
	bf.rwmu.RLock()
	defer bf.rwmu.RUnlock()

	return bf.bloomf != nil && bf.bloomf.TestString(key)
}

func (bf *bloomFilter) KeyExists(key string) bool {
	if bf.keyIsInBloomFilter(key) {
		if bf.filterCheck != nil {
			return bf.filterCheck(key)
		} else {
			return true
		}
	} else {
		return false
	}
}
