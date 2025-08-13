package performance

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Monitor 性能监控器
type Monitor struct {
	metrics    *Metrics
	collectors []Collector
	mu         sync.RWMutex
	running    bool
	stopCh     chan struct{}
}

// Metrics 性能指标
type Metrics struct {
	// HTTP 指标
	HTTPRequests      int64         `json:"http_requests"`
	HTTPErrors        int64         `json:"http_errors"`
	HTTPResponseTime  time.Duration `json:"http_response_time"`
	
	// API 指标
	APIRequests       map[string]int64 `json:"api_requests"`
	APIErrors         map[string]int64 `json:"api_errors"`
	APIResponseTimes  map[string]time.Duration `json:"api_response_times"`
	
	// 缓存指标
	CacheHits         int64 `json:"cache_hits"`
	CacheMisses       int64 `json:"cache_misses"`
	CacheSize         int64 `json:"cache_size"`
	
	// 系统指标
	MemoryUsage       uint64 `json:"memory_usage"`
	GoroutineCount    int    `json:"goroutine_count"`
	GCPauseTime       time.Duration `json:"gc_pause_time"`
	
	// 组件指标
	WidgetLoadTimes   map[string]time.Duration `json:"widget_load_times"`
	WidgetErrors      map[string]int64 `json:"widget_errors"`
	
	// 时间戳
	LastUpdated       time.Time `json:"last_updated"`
	
	mu sync.RWMutex
}

// Collector 指标收集器接口
type Collector interface {
	Collect(ctx context.Context) error
	GetName() string
}

// NewMonitor 创建性能监控器
func NewMonitor() *Monitor {
	return &Monitor{
		metrics: &Metrics{
			APIRequests:      make(map[string]int64),
			APIErrors:        make(map[string]int64),
			APIResponseTimes: make(map[string]time.Duration),
			WidgetLoadTimes:  make(map[string]time.Duration),
			WidgetErrors:     make(map[string]int64),
		},
		collectors: []Collector{
			NewSystemCollector(),
			NewCacheCollector(),
			NewAPICollector(),
		},
		stopCh: make(chan struct{}),
	}
}

// Start 启动监控
func (m *Monitor) Start(ctx context.Context, interval time.Duration) {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.collect(ctx)
		}
	}
}

// Stop 停止监控
func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.running {
		close(m.stopCh)
		m.running = false
	}
}

// collect 收集指标
func (m *Monitor) collect(ctx context.Context) {
	for _, collector := range m.collectors {
		if err := collector.Collect(ctx); err != nil {
			// 记录错误但继续收集其他指标
			continue
		}
	}
	
	m.metrics.mu.Lock()
	m.metrics.LastUpdated = time.Now()
	m.metrics.mu.Unlock()
}

// GetMetrics 获取当前指标
func (m *Monitor) GetMetrics() *Metrics {
	m.metrics.mu.RLock()
	defer m.metrics.mu.RUnlock()
	
	// 深拷贝指标
	metrics := &Metrics{
		HTTPRequests:     m.metrics.HTTPRequests,
		HTTPErrors:       m.metrics.HTTPErrors,
		HTTPResponseTime: m.metrics.HTTPResponseTime,
		CacheHits:        m.metrics.CacheHits,
		CacheMisses:      m.metrics.CacheMisses,
		CacheSize:        m.metrics.CacheSize,
		MemoryUsage:      m.metrics.MemoryUsage,
		GoroutineCount:   m.metrics.GoroutineCount,
		GCPauseTime:      m.metrics.GCPauseTime,
		LastUpdated:      m.metrics.LastUpdated,
		APIRequests:      make(map[string]int64),
		APIErrors:        make(map[string]int64),
		APIResponseTimes: make(map[string]time.Duration),
		WidgetLoadTimes:  make(map[string]time.Duration),
		WidgetErrors:     make(map[string]int64),
	}
	
	// 复制映射
	for k, v := range m.metrics.APIRequests {
		metrics.APIRequests[k] = v
	}
	for k, v := range m.metrics.APIErrors {
		metrics.APIErrors[k] = v
	}
	for k, v := range m.metrics.APIResponseTimes {
		metrics.APIResponseTimes[k] = v
	}
	for k, v := range m.metrics.WidgetLoadTimes {
		metrics.WidgetLoadTimes[k] = v
	}
	for k, v := range m.metrics.WidgetErrors {
		metrics.WidgetErrors[k] = v
	}
	
	return metrics
}

// RecordHTTPRequest 记录HTTP请求
func (m *Monitor) RecordHTTPRequest(duration time.Duration, isError bool) {
	m.metrics.mu.Lock()
	defer m.metrics.mu.Unlock()
	
	m.metrics.HTTPRequests++
	if isError {
		m.metrics.HTTPErrors++
	}
	
	// 计算平均响应时间
	if m.metrics.HTTPRequests == 1 {
		m.metrics.HTTPResponseTime = duration
	} else {
		m.metrics.HTTPResponseTime = (m.metrics.HTTPResponseTime + duration) / 2
	}
}

// RecordAPIRequest 记录API请求
func (m *Monitor) RecordAPIRequest(service string, duration time.Duration, isError bool) {
	m.metrics.mu.Lock()
	defer m.metrics.mu.Unlock()
	
	if m.metrics.APIRequests[service] == 0 {
		m.metrics.APIRequests[service] = 0
		m.metrics.APIErrors[service] = 0
		m.metrics.APIResponseTimes[service] = 0
	}
	
	m.metrics.APIRequests[service]++
	if isError {
		m.metrics.APIErrors[service]++
	}
	
	// 计算平均响应时间
	if m.metrics.APIRequests[service] == 1 {
		m.metrics.APIResponseTimes[service] = duration
	} else {
		m.metrics.APIResponseTimes[service] = (m.metrics.APIResponseTimes[service] + duration) / 2
	}
}

// RecordWidgetLoad 记录组件加载
func (m *Monitor) RecordWidgetLoad(widgetType string, duration time.Duration, isError bool) {
	m.metrics.mu.Lock()
	defer m.metrics.mu.Unlock()
	
	if m.metrics.WidgetLoadTimes[widgetType] == 0 {
		m.metrics.WidgetLoadTimes[widgetType] = 0
		m.metrics.WidgetErrors[widgetType] = 0
	}
	
	if isError {
		m.metrics.WidgetErrors[widgetType]++
	} else {
		// 计算平均加载时间
		count := m.metrics.APIRequests[widgetType] - m.metrics.WidgetErrors[widgetType]
		if count == 1 {
			m.metrics.WidgetLoadTimes[widgetType] = duration
		} else if count > 1 {
			m.metrics.WidgetLoadTimes[widgetType] = (m.metrics.WidgetLoadTimes[widgetType] + duration) / 2
		}
	}
}

// SystemCollector 系统指标收集器
type SystemCollector struct{}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (s *SystemCollector) GetName() string {
	return "system"
}

func (s *SystemCollector) Collect(ctx context.Context) error {
	// 这里应该从全局监控器获取指标，简化实现
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// 更新系统指标（需要访问全局监控器实例）
	// 这里简化处理，实际应该通过依赖注入
	
	return nil
}

// CacheCollector 缓存指标收集器
type CacheCollector struct{}

func NewCacheCollector() *CacheCollector {
	return &CacheCollector{}
}

func (c *CacheCollector) GetName() string {
	return "cache"
}

func (c *CacheCollector) Collect(ctx context.Context) error {
	// 收集缓存指标
	return nil
}

// APICollector API指标收集器
type APICollector struct{}

func NewAPICollector() *APICollector {
	return &APICollector{}
}

func (a *APICollector) GetName() string {
	return "api"
}

func (a *APICollector) Collect(ctx context.Context) error {
	// 收集API指标
	return nil
}
