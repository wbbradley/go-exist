// A thread-safe mechanism for swapping out existence keys.

package filter

import (
	"sync"

	"github.com/willf/bloom"
)

type ExistenceFilter interface {
	KeyExists(key string) bool
}

type ReplaceableExistenceFilter interface {
	ExistenceFilter
	ReplaceKeys(keys []string)
}

type mapFilter struct {
	rwmu sync.RWMutex
	keys map[string]bool
}

func (f *mapFilter) ReplaceKeys(keys []string) {
	filterKeys := make(map[string]bool)
	for _, key := range keys {
		filterKeys[key] = true
	}

	f.rwmu.Lock()
	defer f.rwmu.Unlock()

	f.keys = filterKeys
}

func (f *mapFilter) KeyExists(key string) bool {
	f.rwmu.RLock()
	defer f.rwmu.RUnlock()

	return f.keys[key]
}

func NewMapFilter(keys []string) ReplaceableExistenceFilter {
	f := new(mapFilter)
	f.ReplaceKeys(keys)
	return f
}

type bloomFilter struct {
	rwmu   sync.RWMutex
	bloomf *bloom.BloomFilter
}

func NewBloomFilter(keys []string) ReplaceableExistenceFilter {
	f := new(bloomFilter)
	if len(keys) > 0 {
		f.ReplaceKeys(keys)
	}
	return f
}

func (bf *bloomFilter) ReplaceKeys(keys []string) {
	// Load the new keys into a Bloom filter
	bloomf := bloom.NewWithEstimates(uint(len(keys)), 0.001)
	for _, key := range keys {
		bloomf.AddString(key)
	}

	data, _ := bloomf.MarshalJSON()

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
