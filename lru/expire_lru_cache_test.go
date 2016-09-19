package lru

import (
	"fmt"
	"testing"
	"time"
)

func TestLruCache(t *testing.T) {
	cache := NewExpiringLRUCache(204800, 600)
	t1 := time.Now().UnixNano()
	for i := 0; i < 409600; i++ {
		cache.Put(int64(i), "world", 500)
	}
	t2 := time.Now().UnixNano()
	fmt.Println("Put 409600 cost time ", (t2-t1)/1000, " us")
	vs := make([]string, 409600, 409600)
	var defaultV string = ""
	for i := 0; i < 40960; i++ {
		v := cache.Get(int64(i), defaultV).(string)
		vs[i] = v
	}
	t3 := time.Now().UnixNano()
	fmt.Println("Get 40960 cost time ", (t3-t2)/1000, " us")
	keys := make([]interface{}, 40960, 40960)
	for i := 0; i < 40960; i++ {
		keys[i] = int64(i + 40960)
	}
	t4 := time.Now().UnixNano()
	v1s := cache.Mget(keys, defaultV, true)
	t5 := time.Now().UnixNano()
	v2s := cache.Mget(keys, defaultV, false)
	t6 := time.Now().UnixNano()
	fmt.Println("Mget 40960 strict cost ", (t5-t4)/1000, " us, and Mget 40960 no strict cost ", (t6-t5)/1000, " us, v1s len is ", len(v1s), ", and v2s len is ", len(v2s))
}
