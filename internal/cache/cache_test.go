package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	
	cache.Set("test-key", "test-value")
	
	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find value in cache")
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}
}

func TestCache_GetNotFound(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	
	_, found := cache.Get("non-existent-key")
	if found {
		t.Error("Expected not to find value in cache")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	
	cache.Set("test-key", "test-value")
	cache.Delete("test-key")
	
	_, found := cache.Get("test-key")
	if found {
		t.Error("Expected value to be deleted from cache")
	}
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Clear()
	
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	
	if found1 || found2 {
		t.Error("Expected cache to be cleared")
	}
}

func TestCache_Size(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	
	if cache.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cache.Size())
	}
	
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	
	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	
	cache.Set("test-key", "test-value")
	
	// Value should be found immediately
	_, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find value immediately after set")
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Value should not be found after expiration
	_, found = cache.Get("test-key")
	if found {
		t.Error("Expected value to be expired")
	}
}
