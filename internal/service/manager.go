package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/glance-china/internal/performance"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	clients   map[string]APIClient
	cache     CacheManager
	limiter   RateLimiter
	config    *Config
	monitor   *performance.Monitor
	optimizer *performance.Optimizer
	mu        sync.RWMutex
}

// APIClient 通用API客户端接口
type APIClient interface {
	GetName() string
	GetBaseURL() string
	IsHealthy(ctx context.Context) bool
	SetRateLimit(requests int, duration time.Duration)
	Request(ctx context.Context, req *APIRequest) (*APIResponse, error)
}

// APIRequest 统一请求结构
type APIRequest struct {
	Method   string
	Path     string
	Headers  map[string]string
	Params   map[string]interface{}
	Body     interface{}
	Timeout  time.Duration
}

// APIResponse 统一响应结构
type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Duration   time.Duration
}

// Config 服务配置
type Config struct {
	Region      string                    `yaml:"region"`
	APISources  map[string]APISourceConfig `yaml:"api-sources"`
	Cache       CacheConfig              `yaml:"cache"`
	RateLimit   RateLimitConfig          `yaml:"rate-limit"`
	Performance PerformanceConfig        `yaml:"performance"`
}

type APISourceConfig struct {
	BaseURL     string            `yaml:"base-url"`
	RateLimit   int              `yaml:"rate-limit"`
	Timeout     time.Duration    `yaml:"timeout"`
	Headers     map[string]string `yaml:"headers"`
	Token       string           `yaml:"token"`
	Fallbacks   []string         `yaml:"fallbacks"`
}

type PerformanceConfig struct {
	MaxWorkers        int           `yaml:"max-workers"`
	WorkerIdleTimeout time.Duration `yaml:"worker-idle-timeout"`
	QueueSize         int           `yaml:"queue-size"`
	GCPercent         int           `yaml:"gc-percent"`
	MaxMemoryMB       int           `yaml:"max-memory-mb"`
}

// NewServiceManager 创建服务管理器
func NewServiceManager(config *Config) *ServiceManager {
	sm := &ServiceManager{
		clients:   make(map[string]APIClient),
		config:    config,
		cache:     NewCacheManager(config.Cache),
		limiter:   NewRateLimiter(config.RateLimit),
		monitor:   performance.NewMonitor(),
		optimizer: performance.NewOptimizer(performance.OptimizerConfig{
			MaxWorkers:        config.Performance.MaxWorkers,
			WorkerIdleTimeout: config.Performance.WorkerIdleTimeout,
			QueueSize:         config.Performance.QueueSize,
			GCPercent:         config.Performance.GCPercent,
			MaxMemoryMB:       config.Performance.MaxMemoryMB,
		}),
	}
	
	// 启动性能监控
	go sm.monitor.Start(context.Background(), 30*time.Second)
	
	// 初始化各种服务客户端
	sm.initializeClients()
	
	return sm
}

func (sm *ServiceManager) initializeClients() {
	// 初始化 Bilibili 客户端
	if config, exists := sm.config.APISources["bilibili"]; exists {
		sm.clients["bilibili"] = NewBilibiliClient(config)
	}
	
	// 初始化知乎客户端
	if config, exists := sm.config.APISources["zhihu"]; exists {
		sm.clients["zhihu"] = NewZhihuClient(config)
	}
	
	// 初始化 Gitee 客户端
	if config, exists := sm.config.APISources["gitee"]; exists {
		sm.clients["gitee"] = NewGiteeClient(config)
	}
	
	// 初始化微博客户端
	if config, exists := sm.config.APISources["weibo"]; exists {
		sm.clients["weibo"] = NewWeiboClient(config)
	}
	
	// 初始化斗鱼客户端
	if config, exists := sm.config.APISources["douyu"]; exists {
		sm.clients["douyu"] = NewDouyuClient(config)
	}
}

// GetClient 获取指定服务的客户端
func (sm *ServiceManager) GetClient(serviceName string) (APIClient, error) {
	sm.mu.RLock()
	client, exists := sm.clients[serviceName]
	sm.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("service client not found: %s", serviceName)
	}
	
	return client, nil
}

// RequestWithFallback 带容错的请求
func (sm *ServiceManager) RequestWithFallback(ctx context.Context, serviceName string, req *APIRequest) (*APIResponse, error) {
	start := time.Now()
	var err error
	
	defer func() {
		duration := time.Since(start)
		sm.monitor.RecordAPIRequest(serviceName, duration, err != nil)
	}()
	
	// 尝试主要服务
	client, err := sm.GetClient(serviceName)
	if err != nil {
		return nil, err
	}
	
	// 检查限流
	if !sm.limiter.Allow(serviceName) {
		err = fmt.Errorf("rate limit exceeded for service: %s", serviceName)
		return nil, err
	}
	
	// 使用工作池执行请求
	pool := sm.optimizer.GetOrCreatePool(serviceName)
	resultCh := make(chan error, 1)
	
	job := performance.Job{
		ID:      fmt.Sprintf("%s-%d", serviceName, time.Now().UnixNano()),
		Context: ctx,
		Result:  resultCh,
		Function: func(ctx context.Context) error {
			resp, reqErr := client.Request(ctx, req)
			if reqErr == nil && resp.StatusCode < 400 {
				return nil
			}
			return reqErr
		},
	}
	
	if err := pool.Submit(job); err != nil {
		return nil, err
	}
	
	// 等待结果
	select {
	case err = <-resultCh:
		if err == nil {
			resp, _ := client.Request(ctx, req)
			return resp, nil
		}
	case <-ctx.Done():
		err = ctx.Err()
		return nil, err
	}
	
	// 尝试备用服务
	if config, exists := sm.config.APISources[serviceName]; exists {
		for _, fallback := range config.Fallbacks {
			if fallbackClient, fallbackErr := sm.GetClient(fallback); fallbackErr == nil {
				if resp, fallbackErr := fallbackClient.Request(ctx, req); fallbackErr == nil {
					return resp, nil
				}
			}
		}
	}
	
	return nil, fmt.Errorf("all services failed for: %s", serviceName)
}

// GetMetrics 获取性能指标
func (sm *ServiceManager) GetMetrics() *performance.Metrics {
	return sm.monitor.GetMetrics()
}
