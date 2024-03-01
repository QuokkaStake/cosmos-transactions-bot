package cache_test

import (
	cachePkg "main/pkg/cache"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCacheSet(t *testing.T) {
	t.Parallel()

	cache := cachePkg.NewCache()
	cache.Set("key", "value")

	entry, found := cache.Get("key")
	require.Equal(t, "value", entry)
	require.True(t, found)
}

func TestCacheGetNotExists(t *testing.T) {
	t.Parallel()

	cache := cachePkg.NewCache()
	_, found := cache.Get("key")
	require.False(t, found)
}

func TestCacheGetExpired(t *testing.T) {
	t.Parallel()

	cache := cachePkg.NewCache()
	cache.Entries["key"] = cachePkg.CacheEntry{
		Value: "test", StoredAt: time.Now().Add(-24 * time.Hour),
	}
	_, found := cache.Get("key")
	require.False(t, found)
}
