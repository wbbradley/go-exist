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
