package pokecache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache(10 * time.Second)

	key := "test_key"
	value := []byte("test_value")

	cache.Add(key, value)

	retrieved, found := cache.Get(key)
	if !found {
		t.Errorf("Expected key %s in cache, none found", key)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected value %s, but got %s", value, retrieved)
	}
}
