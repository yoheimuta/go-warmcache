package warmcache

import (
	"fmt"
	"sync/atomic"
	"time"
)

func Example() {
	var counter int32

	// prepare cache setting
	cacheFunc := func() (interface{}, error) {
		atomic.AddInt32(&counter, 1)
		return atomic.LoadInt32(&counter), nil
	}
	interval := 1 * time.Second

	// create CacheMgr
	c, err := NewCacheMgr(cacheFunc, interval, nil)
	if err != nil {
		fmt.Println("cacheFunc is failed", err)
		return
	}

	// fetch cache. no origin access. always warm.
	data, isStale, err := c.CacheData()
	if isStale {
		fmt.Println("cache is stale")
		return
	}
	if err != nil {
		fmt.Println("cacheFunc is failed", err)
		return
	}
	fmt.Println(data)

	// wait updating
	time.Sleep(interval + 100*time.Millisecond)

	// fetch new cache.
	data, isStale, err = c.CacheData()
	if isStale {
		fmt.Println("cache is stale")
		return
	}
	if err != nil {
		fmt.Println("cacheFunc is failed", err)
		return
	}
	fmt.Println(data)

	// Output:
	// 1
	// 2
}
