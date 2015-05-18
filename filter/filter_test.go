package filter

import "testing"

func TestMapFilter(t *testing.T) {
	filter := NewMapFilter([]string{"a", "b"})
	if !filter.KeyExists("a") {
		t.Error("failed to find key")
	}
	if filter.KeyExists("c") {
		t.Error("found non-existent key")
	}
}

func TestBloomFilter(t *testing.T) {
	filter := NewBloomFilter([]string{"a", "b", "c"})
	if !filter.KeyExists("a") {
		t.Error("failed to find key")
	}
	if !filter.KeyExists("b") {
		t.Error("failed to find key")
	}
	if !filter.KeyExists("c") {
		t.Error("failed to find key")
	}
}

func TestBloomFilterEmpty(t *testing.T) {
	filter := NewBloomFilter([]string{})
	if filter.KeyExists("a") {
		t.Error("no key should exist")
	}
	if filter.KeyExists("b") {
		t.Error("no key should exist")
	}
	if filter.KeyExists("c") {
		t.Error("no key should exist")
	}
}
