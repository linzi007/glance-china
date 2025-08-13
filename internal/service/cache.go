package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CacheManager 缓存管理器
type CacheManager interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type     string        `yaml:"type"`     // memory, redis, disk
	TTL      time.Duration `yaml:"ttl"`      // 默认TTL
	MaxSize  int64         `yaml:"max-size"` // 最大缓存大小
	RedisURL string        `yaml:"redis-url"`
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data   map[string]*cacheItem
	mu     sync.RWMutex
	maxSize int64
	currentSize int64
}

type cacheItem struct {
	value     []byte
	expiredAt time.Time
}

func NewCacheManager(config CacheConfig) CacheManager {
	switch config.Type {
	case "redis":
		return NewRedisCache(config)
	case "disk":
		return NewDiskCache(config)
	default:
		return NewMemoryCache(config)
	}
}

func NewMemoryCache(config CacheConfig) *MemoryCache {
	cache := &MemoryCache{
		data:    make(map[string]*cacheItem),
		maxSize: config.MaxSize,
	}
	
	// 启动清理协程
	go cache.cleanup()
	
	return cache
}

func (m *MemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	m.mu.RLock()
	item, exists := m.data[key]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cache miss: %s", key)
	}
	
	if time.Now().After(item.expiredAt) {
		m.Delete(ctx, key)
		return fmt.Errorf("cache expired: %s", key)
	}
	
	return json.Unmarshal(item.value, dest)
}

func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查缓存大小限制
	if m.maxSize > 0 && m.currentSize+int64(len(data)) > m.maxSize {
		m.evictLRU()
	}
	
	m.data[key] = &cacheItem{
		value:     data,
		expiredAt: time.Now().Add(ttl),
	}
	m.currentSize += int64(len(data))
	
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if item, exists := m.data[key]; exists {
		m.currentSize -= int64(len(item.value))
		delete(m.data, key)
	}
	
	return nil
}

func (m *MemoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data = make(map[string]*cacheItem)
	m.currentSize = 0
	
	return nil
}

func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, item := range m.data {
			if now.After(item.expiredAt) {
				m.currentSize -= int64(len(item.value))
				delete(m.data, key)
			}
		}
		m.mu.Unlock()
	}
}

func (m *MemoryCache) evictLRU() {
	// 简单的LRU实现：删除最旧的条目
	var oldestKey string
	var oldestTime time.Time
	
	for key, item := range m.data {
		if oldestKey == "" || item.expiredAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiredAt
		}
	}
	
	if oldestKey != "" {
		m.currentSize -= int64(len(m.data[oldestKey].value))
		delete(m.data, oldestKey)
	}
}

// RedisCache Redis缓存实现（占位符）
type RedisCache struct {
	// Redis客户端实现
}

func NewRedisCache(config CacheConfig) *RedisCache {
	// TODO: 实现Redis缓存
	return &RedisCache{}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	return fmt.Errorf("redis cache not implemented")
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return fmt.Errorf("redis cache not implemented")
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return fmt.Errorf("redis cache not implemented")
}

func (r *RedisCache) Clear(ctx context.Context) error {
	return fmt.Errorf("redis cache not implemented")
}

// DiskCache 磁盘缓存实现（占位符）
type DiskCache struct {
	// 磁盘缓存实现
}

func NewDiskCache(config CacheConfig) *DiskCache {
	// TODO: 实现磁盘缓存
	return &DiskCache{}
}

func (d *DiskCache) Get(ctx context.Context, key string, dest interface{}) error {
	return fmt.Errorf("disk cache not implemented")
}

func (d *DiskCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return fmt.Errorf("disk cache not implemented")
}

func (d *DiskCache) Delete(ctx context.Context, key string) error {
	return fmt.Errorf("disk cache not implemented")
}

func (d *DiskCache) Clear(ctx context.Context) error {
	return fmt.Errorf("disk cache not implemented")
}
