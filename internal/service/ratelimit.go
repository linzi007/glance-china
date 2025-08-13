package service

import (
	"sync"
	"time"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(service string) bool
	Reset(service string)
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	DefaultLimit int           `yaml:"default-limit"` // 默认每分钟请求数
	Services     map[string]int `yaml:"services"`      // 各服务的限流配置
	Window       time.Duration `yaml:"window"`        // 时间窗口
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	buckets map[string]*tokenBucket
	config  RateLimitConfig
	mu      sync.RWMutex
}

type tokenBucket struct {
	tokens    int
	capacity  int
	refillRate int
	lastRefill time.Time
	mu        sync.Mutex
}

func NewRateLimiter(config RateLimitConfig) RateLimiter {
	if config.Window == 0 {
		config.Window = time.Minute
	}
	if config.DefaultLimit == 0 {
		config.DefaultLimit = 60
	}
	
	return &TokenBucketLimiter{
		buckets: make(map[string]*tokenBucket),
		config:  config,
	}
}

func (t *TokenBucketLimiter) Allow(service string) bool {
	t.mu.RLock()
	bucket, exists := t.buckets[service]
	t.mu.RUnlock()
	
	if !exists {
		t.mu.Lock()
		// 双重检查
		if bucket, exists = t.buckets[service]; !exists {
			limit := t.config.DefaultLimit
			if serviceLimit, ok := t.config.Services[service]; ok {
				limit = serviceLimit
			}
			
			bucket = &tokenBucket{
				tokens:     limit,
				capacity:   limit,
				refillRate: limit,
				lastRefill: time.Now(),
			}
			t.buckets[service] = bucket
		}
		t.mu.Unlock()
	}
	
	return bucket.consume()
}

func (t *TokenBucketLimiter) Reset(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if bucket, exists := t.buckets[service]; exists {
		bucket.mu.Lock()
		bucket.tokens = bucket.capacity
		bucket.lastRefill = time.Now()
		bucket.mu.Unlock()
	}
}

func (b *tokenBucket) consume() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	tokensToAdd := int(elapsed.Minutes()) * b.refillRate
	
	if tokensToAdd > 0 {
		b.tokens = min(b.capacity, b.tokens+tokensToAdd)
		b.lastRefill = now
	}
	
	// 消费令牌
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
