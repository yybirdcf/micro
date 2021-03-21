package algos

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Data map[interface{}]*Value

const MAX_TIMES = 30

var (
	CacheFilterIns = NewCacheFilter()
)

type Value struct {
	value interface{}
	n     int       //过期时间内，访问达到一定次数，才认为存在，否则仍然认为不存在，可以穿透访问缓存，数据库之类的
	ttl   time.Time //过期时间
}

//过滤器，目前用来防止缓存穿透等功能
type CacheFilter struct {
	cleanTime time.Time //上次清理时间
	lock      sync.RWMutex
	data      Data
}

func NewCacheFilter() *CacheFilter {
	return &CacheFilter{
		cleanTime: time.Now(),
		lock:      sync.RWMutex{},
		data:      make(Data),
	}
}

func (self *CacheFilter) Set(key interface{}, value interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if v, ok := self.data[key]; ok {
		v.n = v.n + 1
	} else {
		self.data[key] = &Value{
			value: value,
			n:     1,
			ttl:   time.Now().Add(time.Second * time.Duration(rand.Int63n(50)+10)),
		}
	}
}

func (self *CacheFilter) Exist(key interface{}) bool {
	if time.Now().Sub(self.cleanTime) >= time.Second*5 {
		go func() {
			self.checkExpire()
			self.cleanTime = time.Now()
		}()
	}

	self.lock.RLock()
	defer self.lock.RUnlock()

	v, ok := self.data[key]
	if !ok {
		return false
	}

	if v.ttl.Before(time.Now()) || v.n <= MAX_TIMES {
		return false
	}

	return true
}

func (self *CacheFilter) checkExpire() {
	self.lock.Lock()
	defer self.lock.Unlock()

	for k, v := range self.data {
		if v.ttl.Before(time.Now()) {
			fmt.Printf("delete key: %+v \n", k)
			delete(self.data, k)
		}
	}
}
