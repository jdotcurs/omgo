package omgo_test

import (
	"sync"
	"testing"
	"time"

	"github.com/jdotcurs/omgo"
	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	cache := omgo.NewCache()
	require.NotNil(t, cache)
}

func TestCacheSetGet(t *testing.T) {
	cache := omgo.NewCache()
	key := "testKey"
	value := []byte("testValue")
	expiration := 1 * time.Second

	cache.Set(key, value, expiration)

	retrievedValue, found := cache.Get(key)
	require.True(t, found)
	require.Equal(t, value, retrievedValue)
}

func TestCacheExpiration(t *testing.T) {
	cache := omgo.NewCache()
	key := "testKey"
	value := []byte("testValue")
	expiration := 1 * time.Second

	cache.Set(key, value, expiration)

	time.Sleep(2 * time.Second)

	_, found := cache.Get(key)
	require.False(t, found)
}

func TestCacheConcurrency(t *testing.T) {
	cache := omgo.NewCache()
	key := "testKey"
	value := []byte("testValue")
	expiration := 5 * time.Second

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			cache.Set(key, value, expiration)
			_, _ = cache.Get(key)
		}()
	}
	wg.Wait()

	retrievedValue, found := cache.Get(key)
	require.True(t, found)
	require.Equal(t, value, retrievedValue)
}
