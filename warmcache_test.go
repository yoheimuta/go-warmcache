package warmcache

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheMgr(t *testing.T) {
	var counter int32

	tests := []struct {
		inputInterval    time.Duration
		inputCacheFunc   FetchCacheFunc
		inputElapsedChan chan time.Duration
		wantElapsedMsg   string
		wantCacheData    interface{}
		wantCacheData2   interface{}
		wantIsStale      bool
		wantIsErr        bool
		wantIsErr2       bool
	}{
		{
			inputInterval: 1 * time.Second,
			inputCacheFunc: func() (interface{}, error) {
				time.Sleep(500 * time.Millisecond)
				atomic.AddInt32(&counter, 1)
				return atomic.LoadInt32(&counter), nil
			},
			inputElapsedChan: make(chan time.Duration),
			wantElapsedMsg:   `50.*ms`,
			wantCacheData:    int32(1),
			wantCacheData2:   int32(2),
			wantIsStale:      false,
			wantIsErr:        false,
			wantIsErr2:       false,
		},
		{
			inputInterval: 1 * time.Second,
			inputCacheFunc: func() (interface{}, error) {
				if atomic.LoadInt32(&counter) == 0 {
					atomic.AddInt32(&counter, 1)
					return atomic.LoadInt32(&counter), nil
				}
				return nil, errors.New("failed to fetch a cache")
			},
			inputElapsedChan: make(chan time.Duration),
			wantElapsedMsg:   `.*µs`,
			wantCacheData:    int32(1),
			wantCacheData2:   int32(1),
			wantIsStale:      false,
			wantIsErr:        false,
			wantIsErr2:       true,
		},
		{
			inputInterval: 1 * time.Second,
			inputCacheFunc: func() (interface{}, error) {
				time.Sleep(500 * time.Millisecond)
				atomic.AddInt32(&counter, 1)
				return []string{"hoge", fmt.Sprintf("fuga%d", atomic.LoadInt32(&counter))}, nil
			},
			inputElapsedChan: make(chan time.Duration),
			wantElapsedMsg:   `50.*ms`,
			wantCacheData:    []string{"hoge", "fuga1"},
			wantCacheData2:   []string{"hoge", "fuga2"},
			wantIsStale:      false,
			wantIsErr:        false,
			wantIsErr2:       false,
		},
	}

	for _, test := range tests {
		atomic.StoreInt32(&counter, 0)

		c, err := NewCacheMgr(test.inputCacheFunc, test.inputInterval, test.inputElapsedChan)
		if err != nil {
			t.Fatalf(`failed to NewCacheMgr: err="%v"`, err)
		}

		// first get
		data, isStale, err := c.CacheData()
		if isStale != test.wantIsStale {
			t.Errorf(`got %v, but want %v`, isStale, test.wantIsStale)
		}
		if test.wantIsErr && err == nil {
			t.Errorf(`should be err`)
		}
		if !test.wantIsErr && err != nil {
			t.Errorf(`should be not err, but got %s`, err)
		}
		if !reflect.DeepEqual(data, test.wantCacheData) {
			t.Errorf(`got %v, but want %v`, data, test.wantCacheData)
		}

		// second get
		msg := fmt.Sprintf("%v", <-test.inputElapsedChan)
		if !regexp.MustCompile(test.wantElapsedMsg).Match([]byte(msg)) {
			t.Errorf(`got %s, but want %s`, msg, test.wantElapsedMsg)
		}

		time.Sleep(test.inputInterval * 2)
		data, isStale, err = c.CacheData()
		if isStale != test.wantIsStale {
			t.Errorf(`should be not stale`)
		}
		if test.wantIsErr2 && err == nil {
			t.Errorf(`should be err`)
		}
		if !test.wantIsErr2 && err != nil {
			t.Errorf(`should be not err, but got %s`, err)
		}
		if !reflect.DeepEqual(data, test.wantCacheData2) {
			t.Errorf(`got %v, but want %v`, data, test.wantCacheData2)
		}

		c.Stop()
	}
}
