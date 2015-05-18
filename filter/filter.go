// A thread-safe mechanism for swapping out existence keys.

package filter

import "sync"

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
