package filter

import "sync"

type mapFilter struct {
	rwmu sync.RWMutex
	keys map[string]bool
}

func NewMapFilter() MutableExistenceFilter {
	return new(mapFilter)
}

func (f *mapFilter) ImportKeys(keyCount uint, keys <-chan string) {
	filterKeys := make(map[string]bool)
	for key := range keys {
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
