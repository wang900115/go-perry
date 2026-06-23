package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type SimpleCache[K comparable, V any] interface {
	Add(key K, value V) error
	AddWithTTL(key K, value V, ttl time.Duration) error
	Get(key K) (V, error)
	GetWithLoader(key K, loader func(key K) (V, time.Duration, error)) (V, error)
	Remove(key K)
	Keys() []K
	Values() []V
	Len() int
	OnExpiration(fn func(K, V)) func()
}

type Options struct {
	SizeLimit  uint64
	DefaultTTL time.Duration
}

const evictionTimeout = 1 * time.Hour

func NewSimpleCache[K comparable, V any](options ...Options) SimpleCache[K, V] {
	opts := []ttlcache.Option[K, V]{
		ttlcache.WithDisableTouchOnHit[K, V](),
	}
	if len(options) > 0 {
		o := options[0]
		if o.SizeLimit > 0 {
			opts = append(opts, ttlcache.WithCapacity[K, V](o.SizeLimit))
		}
		if o.DefaultTTL > 0 {
			opts = append(opts, ttlcache.WithTTL[K, V](o.DefaultTTL))
		}
	}
	c := ttlcache.New[K, V](opts...)
	cache := &simpleCache[K, V]{
		data: c,
	}
	go cache.data.Start()

	// Register a cleanup function to stop the cache when the program exits
	// This ensures that the cache is properly stopped and resources are released when the program terminates
	runtime.AddCleanup(cache, func(ttlCache *ttlcache.Cache[K, V]) {
		ttlCache.Stop()
	}, cache.data)

	return cache
}

type simpleCache[K comparable, V any] struct {
	data            *ttlcache.Cache[K, V]
	evictionDeadlin atomic.Pointer[time.Time]
}

func (c *simpleCache[K, V]) Add(key K, value V) error {
	c.evictExpired()
	return c.AddWithTTL(key, value, ttlcache.DefaultTTL)
}

func (c *simpleCache[K, V]) AddWithTTL(key K, value V, ttl time.Duration) error {
	c.evictExpired()
	item := c.data.Set(key, value, ttl)
	if item == nil {
		return errors.New("failed to add item to cache")
	}
	return nil
}

func (c *simpleCache[K, V]) Remove(key K) {
	c.data.Delete(key)
}

func (c *simpleCache[K, V]) Get(key K) (V, error) {
	item := c.data.Get(key)
	if item == nil {
		var zero V
		return zero, errors.New("item not found")
	}
	return item.Value(), nil
}

func (c *simpleCache[K, V]) GetWithLoader(key K, loader func(key K) (V, time.Duration, error)) (V, error) {
	var err error

	// Wrap the loader function to handle eviction of expired items before loading a new item
	// example: db query, file read, etc.
	loaderWrapper := ttlcache.LoaderFunc[K, V](
		func(t *ttlcache.Cache[K, V], key K) *ttlcache.Item[K, V] {
			c.evictExpired()
			var value V
			var ttl time.Duration
			value, ttl, err = loader(key)
			if err != nil {
				return nil
			}
			return t.Set(key, value, ttl)
		},
	)
	item := c.data.Get(key, ttlcache.WithLoader[K, V](loaderWrapper))
	if item == nil {
		var zero V
		if err != nil {
			return zero, fmt.Errorf("cache error: loader returned %w", err)
		}
		return zero, errors.New("item not found")
	}
	return item.Value(), nil
}

func (c *simpleCache[K, V]) evictExpired() {
	if c.evictionDeadlin.Load() == nil || c.evictionDeadlin.Load().Before(time.Now()) {
		c.data.DeleteExpired()
		c.evictionDeadlin.Store(new(time.Now().Add(evictionTimeout)))
	}
}

func (c *simpleCache[K, V]) Keys() []K {
	res := make([]K, 0, c.data.Len())
	c.data.Range(func(item *ttlcache.Item[K, V]) bool {
		if !item.IsExpired() {
			res = append(res, item.Key())
		}
		return true
	})
	return res
}

func (c *simpleCache[K, V]) Values() []V {
	res := make([]V, 0, c.data.Len())
	c.data.Range(func(item *ttlcache.Item[K, V]) bool {
		if !item.IsExpired() {
			res = append(res, item.Value())
		}
		return true
	})
	return res
}

func (c *simpleCache[K, V]) Len() int { return c.data.Len() }

// returned function can be used to unregister the expiration callback if needed
// if you want to stop receiving expiration notifications, simply call the returned function to unregister the callback
func (c *simpleCache[K, V]) OnExpiration(fn func(K, V)) func() {
	return c.data.OnEviction(func(_ context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[K, V]) {
		if reason == ttlcache.EvictionReasonExpired {
			fn(item.Key(), item.Value())
		}
	})
}
