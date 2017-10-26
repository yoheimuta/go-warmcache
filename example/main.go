package main

import (
	"log"
	"sync/atomic"
	"time"

	warmcache "github.com/yoheimuta/go-warmcache"
)

func main() {
	var counter int32

	// prepare cache setting
	cacheFunc := func() (interface{}, error) {
		atomic.AddInt32(&counter, 1)
		return atomic.LoadInt32(&counter), nil
	}
	interval := 1 * time.Second

	// create CacheMgr
	c, err := warmcache.NewCacheMgr(cacheFunc, interval, nil)
	if err != nil {
		log.Fatal("cacheFunc is failed", err)
	}

	// fetch cache. no origin access. always warm.
	data, isStale, err := c.CacheData()
	if isStale {
		log.Fatal("cache is stale")
	}
	if err != nil {
		log.Fatal("cacheFunc is failed", err)
	}
	log.Println("cache data", data)

	// wait updating
	time.Sleep(interval + 100*time.Millisecond)

	// fetch new cache.
	data, isStale, err = c.CacheData()
	if isStale {
		log.Fatal("cache is stale")
	}
	if err != nil {
		log.Fatal("cacheFunc is failed", err)
	}
	log.Println("cache data", data)
}
