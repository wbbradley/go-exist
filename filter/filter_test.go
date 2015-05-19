package filter

import (
	"os"
	"testing"
)

func TestMapFilter(t *testing.T) {
	filter := NewMapFilter()
	ReadKeysIntoFilter(filter, []string{"a", "b"})

	if !filter.KeyExists("a") {
		t.Error("failed to find key")
	}
	if filter.KeyExists("c") {
		t.Error("found non-existent key")
	}
}

func TestBloomFilterCheckRealPositive(t *testing.T) {
	filter := NewBloomFilter(func(key string) bool {
		return true
	})
	ReadKeysIntoFilter(filter, []string{"a"})

	if !filter.KeyExists("a") {
		t.Error("failed to find a real positive key")
	}
}

func TestBloomFilterCheckNoMatchWithCheck(t *testing.T) {
	filter := NewBloomFilter(func(key string) bool {
		return true
	})
	ReadKeysIntoFilter(filter, []string{"a"})

	if filter.KeyExists("b") {
		t.Error("found a non-existent key")
	}
}

func TestBloomFilterCheckFalsePositive(t *testing.T) {
	filter := NewBloomFilter(func(key string) bool {
		return false
	})
	ReadKeysIntoFilter(filter, []string{"a"})

	if filter.KeyExists("a") {
		t.Error("found a false positive key")
	}
}

func TestBloomFilter(t *testing.T) {
	filter := NewBloomFilter(nil)
	ReadKeysIntoFilter(filter, []string{"a", "b", "c"})

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
	filter := NewBloomFilter(nil)
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

func TestReadingAFileIntoAFilter(t *testing.T) {
	testFilename := "../tests/email_data.txt"
	var f *os.File
	var err error

	f, err = os.Open(testFilename)
	if err != nil {
		t.Error("failed to open", testFilename)
	}

	filter := NewBloomFilter(nil)
	ReadStreamIntoFilter(filter, 100, f)
	if !filter.KeyExists("fake1@fakeplace.net") {
		t.Error("failed to load file correctly into Bloom filter")
	}
}
