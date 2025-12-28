package tools

import (
	"sync"
	"time"
)

// item 表示缓存条目，包含值与过期时间。
// 这里使用 string 作为通用值类型，已满足当前 uuid->answer 的需求；
// 如需扩展可改为 any。
type item struct {
	value    string
	expireAt time.Time
}

// 内存缓存（并发安全）
var (
	cacheMu sync.RWMutex
	cache   = make(map[string]item)
)

// 启动周期清理过期键，避免长期占用内存。
func init() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cleanupExpired()
		}
	}()
}

// CacheSet 设置键值并指定 TTL。
func CacheSet(key, value string, ttl time.Duration) {
	cacheMu.Lock()
	cache[key] = item{value: value, expireAt: time.Now().Add(ttl)}
	cacheMu.Unlock()
}

// CacheGet 获取键值；若不存在或已过期则返回 ok=false，并清理该键。
func CacheGet(key string) (val string, ok bool) {
	cacheMu.RLock()
	it, exists := cache[key]
	cacheMu.RUnlock()
	if !exists {
		return "", false
	}
	if time.Now().After(it.expireAt) {
		CacheDelete(key)
		return "", false
	}
	return it.value, true
}

// CacheDelete 删除指定键。
func CacheDelete(key string) {
	cacheMu.Lock()
	delete(cache, key)
	cacheMu.Unlock()
}

// cleanupExpired 扫描清理过期键；由后台协程周期调用。
func cleanupExpired() {
	now := time.Now()
	cacheMu.Lock()
	for k, v := range cache {
		if now.After(v.expireAt) {
			delete(cache, k)
		}
	}
	cacheMu.Unlock()
}
