package warmcache

import (
	"sync"
	"time"
)

// FetchCacheFunc is a function for fetching origin data and then caching it.
type FetchCacheFunc func() (interface{}, error)

// CacheMgr caches origin data and updates it.
type CacheMgr struct {
	mu             sync.RWMutex
	cacheData      interface{}
	isStale        bool
	fetchCacheFunc FetchCacheFunc
	fetchInterval  time.Duration
	fetchErr       error
	fetchStopChan  chan struct{}
	elapsedChan    chan<- time.Duration
}

// NewCacheMgr creates cachemgr.
// If elapsedChan is nil, send no messages to elapsedChan.
func NewCacheMgr(cacheFunc FetchCacheFunc, interval time.Duration, elapsedChan chan<- time.Duration) (*CacheMgr, error) {
	cacheData, err := cacheFunc()
	if err != nil {
		return nil, err
	}
	c := &CacheMgr{
		cacheData:      cacheData,
		isStale:        false,
		fetchCacheFunc: cacheFunc,
		fetchInterval:  interval,
		fetchStopChan:  make(chan struct{}),
		elapsedChan:    elapsedChan,
	}
	go func() {
		<-time.After(c.fetchInterval)
		c.fetchLoop()
	}()
	return c, nil
}

// CacheData returns cached data.
// If the given FetchCacheFunc returned err, CacheData returns staled data and the error.
func (c *CacheMgr) CacheData() (interface{}, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cacheData, c.isStale, c.fetchErr
}

// Stop stops auto fetching cache.
func (c *CacheMgr) Stop() {
	close(c.fetchStopChan)
}

func (c *CacheMgr) fetchLoop() {
	for {
		timerChan := time.After(c.fetchInterval)
		data, isStale, err := c.fetch()

		c.mu.Lock()
		c.cacheData, c.isStale, c.fetchErr = data, isStale, err
		c.mu.Unlock()

		select {
		case <-timerChan:
		case <-c.fetchStopChan:
			close(c.elapsedChan)
			return
		}
	}
}

func (c *CacheMgr) fetch() (data interface{}, isStale bool, err error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		if c.fetchInterval < elapsed {
			isStale = true
		} else {
			isStale = false
		}
		if c.elapsedChan != nil {
			c.elapsedChan <- elapsed
		}
	}()
	data, err = c.fetchCacheFunc()
	if err != nil {
		// Still use stale cached data
		data = c.cacheData
	}
	return
}
