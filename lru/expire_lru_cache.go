package lru

import (
	"sync"
	"time"
)

type LRUItem struct {
	pos        int64
	expiration int64
	value      interface{}
}

type ExpiringLRUCache struct {
	size       int64
	max        int64
	hits       int64
	misses     int64
	lookups    int64
	evictions  int64
	hand       int64
	expiration int64
	mu         *sync.RWMutex
	data       map[interface{}]*LRUItem
	clockKeys  []interface{}
	clockRefs  []bool
}

func NewExpiringLRUCache(size int64, expiration int64) *ExpiringLRUCache {
	cache := &ExpiringLRUCache{
		size:       size,
		max:        size - 1,
		hits:       0,
		misses:     0,
		lookups:    0,
		evictions:  0,
		hand:       0,
		expiration: expiration,
		mu:         new(sync.RWMutex),
		data:       make(map[interface{}]*LRUItem),
		clockKeys:  make([]interface{}, size),
		clockRefs:  make([]bool, size),
	}
	return cache
}

func (p *ExpiringLRUCache) Put(key, value interface{}, expiration int64) {
	timeout := p.expiration
	if expiration != 0 {
		timeout = expiration
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if item, ok := p.data[key]; ok {
		item.value = value
		item.expiration = (int64(time.Now().Unix()) + timeout)
		p.clockRefs[item.pos] = true
		return
	}
	hand := p.hand
	count := 0
	maxCount := 107
	var ref bool
	var old interface{}
	for true {
		ref = p.clockRefs[hand]
		if ref {
			p.clockRefs[hand] = false
			hand++
			if hand > p.max {
				hand = 0
			}
			count++
			if count >= maxCount {
				p.clockRefs[hand] = false
			}
		} else {
			old = p.clockKeys[hand]
			delete(p.data, old)
			p.evictions++
			p.clockKeys[hand] = key
			p.clockRefs[hand] = true
			p.data[key] = &LRUItem{
				pos:        hand,
				expiration: (int64(time.Now().Unix()) + timeout),
				value:      value,
			}
			hand++
			if hand > p.max {
				hand = 0
			}
			p.hand = hand
			break
		}
	}
}

func (p *ExpiringLRUCache) Get(key interface{}, defaultV interface{}) interface{} {
	p.lookups++
	var item *LRUItem
	var ok bool
	p.mu.RLock()
	item, ok = p.data[key]
	p.mu.RUnlock()
	if ok {
		if item.expiration > int64(time.Now().Unix()) {
			p.hits++
			p.clockRefs[item.pos] = true
			return item.value
		} else {
			p.clockRefs[item.pos] = false
		}
	}
	p.misses++
	return defaultV
}

//strict=false,不检查是否失效,不修改LRU状态
//strict=true, 检查是否失效，修改LRU状态
func (p *ExpiringLRUCache) Mget(keys []interface{}, defaultV interface{}, strict bool) []interface{} {
	size := len(keys)
	values := make([]interface{}, size, size)
	if strict {
		for i := 0; i < size; i++ {
			values[i] = p.Get(keys[i], defaultV)
		}
	} else {
		p.mu.RLock()
		for i := 0; i < size; i++ {
			if item, ok := p.data[keys[i]]; !ok {
				values[i] = defaultV
			} else {
				values[i] = item.value
			}
		}
		p.mu.RUnlock()
	}
	return values
}
