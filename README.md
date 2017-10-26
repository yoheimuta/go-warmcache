# go-warmcache [![Build Status](https://travis-ci.org/yoheimuta/go-warmcache.svg?branch=master)](https://travis-ci.org/yoheimuta/go-warmcache) [![GoDoc](https://godoc.org/github.com/yoheimuta/go-warmcache?status.svg)](https://godoc.org/github.com/yoheimuta/go-warmcache)

go-warmcache is a thin go package which manages an in-memory warm cache.
It provides thread safety.

You can use this package...

- To avoid delay of response time at first access
- To avoid thundering hurd

### Installation

```
go get github.com/yoheimuta/go-warmcache
```

### Usage

See example/main.go and example_test.go in detail.

```go
// create CacheMgr. call cacheFunc to create a warm cache in this method. refresh cache automatically every interval.
c, err := warmcache.NewCacheMgr(cacheFunc, interval, elapsedChan)
if err != nil {
    log.Fatal("cacheFunc is failed", err)
}

// fetch cache. no origin access. always warm.
data, isStale, err := c.CacheData()
if isStale {
    log.Fatal("cache data is stale. cacheFunc is slower than the interval")
}
if err != nil {
    log.Fatal("cacheFunc is failed. cache data is also stale, not empty", err)
}
log.Println("cache data", data)

// after the interval...
// fetch new cache.
data, isStale, err = c.CacheData()
```
