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
	mu   sync.Mutex
	keys map[string]bool
}

func (f *mapFilter) ReplaceKeys(keys []string) {
	filterKeys := make(map[string]bool)
	for _, key := range keys {
		filterKeys[key] = true
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.keys = filterKeys
}

func (f *mapFilter) KeyExists(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.keys[key]
}

func NewMapFilter(keys []string) ExistenceFilter {
	f := new(mapFilter)
	f.ReplaceKeys(keys)
	return f
}
