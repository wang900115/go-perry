package models

import (
	"reflect"

	"github.com/bradfitz/gomemcache/memcache"
)

var cache = memcache.New("localhost:11211")

func CacheData(cacheKey string, ttl int32, fn interface{}) []byte {
	var retValue []byte
	item, err := cache.Get(cacheKey)

	if err != nil { // Cache miss or any other error.

		// Reflect to get the function value
		fnValue := reflect.ValueOf(fn)

		// Create a slice of reflect.Value for the arguments
		args := make([]reflect.Value, fnValue.Type().NumIn())

		// Check if the function returns a slice of bytes
		if fnValue.Type().NumOut() != 1 || fnValue.Type().Out(0).Kind() != reflect.Slice || fnValue.Type().Out(0).Elem().Kind() != reflect.Uint8 {
			panic("function must return a slice of bytes")
		}

		// Call the function
		results := fnValue.Call(args)

		// Convert result to a slice of bytes
		retValue = results[0].Interface().([]byte)

		// Cache the result
		memcacheItem := memcache.Item{Key: cacheKey, Expiration: ttl, Value: retValue}
		if err := cache.Set(&memcacheItem); err != nil {
			panic(err)
		}
	} else { // cache hit
		retValue = item.Value
	}
	return retValue
}
