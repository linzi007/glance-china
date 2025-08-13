package performance

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
	Name            string        `json:"name"`
	TotalRequests   int           `json:"total_requests"`
	SuccessRequests int           `json:"success_requests"`
	FailedRequests  int           `json:"failed_requests"`
	TotalTime       time.Duration `json:"total_time"`
	AverageTime     time.Duration `json:"average_time"`
	MinTime         time.Duration `json:"min_time"`
	MaxTime         time.Duration `json:"max_time"`
	RequestsPerSec  float64       `json:"requests_per_sec"`
	ErrorRate       float64       `json:"error_rate"`
}

// BenchmarkConfig 基准测试配置
type BenchmarkConfig struct {
	Name         string        `json:"name"`
	Concurrency  int           `json:"concurrency"`
	Duration     time.Duration `json:"duration"`
	RequestCount int           `json:"request_count"`
	WarmupTime   time.Duration `json:"warmup_time"`
}

// Benchmarker 基准测试器
type Benchmarker struct {
	config  BenchmarkConfig
	results []BenchmarkResult
	mu      sync.Mutex
}

// TestFunc 测试函数类型
type TestFunc func(ctx context.Context) error

// NewBenchmarker 创建基准测试器
func NewBenchmarker(config BenchmarkConfig) *Benchmarker {
	return &Benchmarker{
		config:  config,
		results: make([]BenchmarkResult, 0),
	}
}

// RunBenchmark 运行基准测试
func (b *Benchmarker) RunBenchmark(ctx context.Context, testFunc TestFunc) (*BenchmarkResult, error) {
	// 预热
	if b.config.WarmupTime > 0 {
		fmt.Printf("预热中... (%v)\n", b.config.WarmupTime)
		warmupCtx, cancel := context.WithTimeout(ctx, b.config.WarmupTime)
		b.runWarmup(warmupCtx, testFunc)
		cancel()
		time.Sleep(time.Second) // 预热后等待1秒
	}
	
	fmt.Printf("开始基准测试: %s\n", b.config.Name)
	fmt.Printf("并发数: %d, 持续时间: %v\n", b.config.Concurrency, b.config.Duration)
	
	result := &BenchmarkResult{
		Name:    b.config.Name,
		MinTime: time.Hour, // 初始化为很大的值
	}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	var times []time.Duration
	
	startTime := time.Now()
	endTime := startTime.Add(b.config.Duration)
	
	// 启动并发测试
	for i := 0; i < b.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for time.Now().Before(endTime) {
				requestStart := time.Now()
				err := testFunc(ctx)
				requestDuration := time.Since(requestStart)
				
				mu.Lock()
				result.TotalRequests++
				times = append(times, requestDuration)
				
				if err != nil {
					result.FailedRequests++
				} else {
					result.SuccessRequests++
				}
				
				if requestDuration < result.MinTime {
					result.MinTime = requestDuration
				}
				if requestDuration > result.MaxTime {
					result.MaxTime = requestDuration
				}
				mu.Unlock()
				
				// 检查是否达到请求数限制
				if b.config.RequestCount > 0 && result.TotalRequests >= b.config.RequestCount {
					break
				}
			}
		}()
	}
	
	wg.Wait()
	
	// 计算统计信息
	result.TotalTime = time.Since(startTime)
	
	if len(times) > 0 {
		var totalTime time.Duration
		for _, t := range times {
			totalTime += t
		}
		result.AverageTime = totalTime / time.Duration(len(times))
	}
	
	if result.TotalTime > 0 {
		result.RequestsPerSec = float64(result.TotalRequests) / result.TotalTime.Seconds()
	}
	
	if result.TotalRequests > 0 {
		result.ErrorRate = float64(result.FailedRequests) / float64(result.TotalRequests) * 100
	}
	
	b.mu.Lock()
	b.results = append(b.results, *result)
	b.mu.Unlock()
	
	return result, nil
}

// runWarmup 运行预热
func (b *Benchmarker) runWarmup(ctx context.Context, testFunc TestFunc) {
	var wg sync.WaitGroup
	
	for i := 0; i < b.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					return
				default:
					testFunc(ctx)
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}
	
	wg.Wait()
}

// GetResults 获取所有测试结果
func (b *Benchmarker) GetResults() []BenchmarkResult {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	results := make([]BenchmarkResult, len(b.results))
	copy(results, b.results)
	return results
}

// PrintResult 打印测试结果
func (b *Benchmarker) PrintResult(result *BenchmarkResult) {
	fmt.Printf("\n=== 基准测试结果: %s ===\n", result.Name)
	fmt.Printf("总请求数: %d\n", result.TotalRequests)
	fmt.Printf("成功请求: %d\n", result.SuccessRequests)
	fmt.Printf("失败请求: %d\n", result.FailedRequests)
	fmt.Printf("总耗时: %v\n", result.TotalTime)
	fmt.Printf("平均响应时间: %v\n", result.AverageTime)
	fmt.Printf("最小响应时间: %v\n", result.MinTime)
	fmt.Printf("最大响应时间: %v\n", result.MaxTime)
	fmt.Printf("每秒请求数: %.2f\n", result.RequestsPerSec)
	fmt.Printf("错误率: %.2f%%\n", result.ErrorRate)
	fmt.Printf("================================\n")
}
