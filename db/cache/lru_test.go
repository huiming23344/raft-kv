package cache

import (
	"testing"
)

func TestGet(t *testing.T) {
	lru := NewLRUCache(128)
	lru.Set("key1", "1234")
	if v, _ := lru.Get("key1"); v != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); !ok {
		t.Fatalf("cache miss key2 failed")
	}
}
