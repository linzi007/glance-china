package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	
	"github.com/glance-china/internal/performance"
	"github.com/glance-china/internal/widget"
)

// TestBilibiliWidgetBenchmark Bilibili组件基准测试
func TestBilibiliWidgetBenchmark(t *testing.T) {
	config := performance.BenchmarkConfig{
		Name:         "Bilibili Widget Load Test",
		Concurrency:  10,
		Duration:     30 * time.Second,
		RequestCount: 1000,
		WarmupTime:   5 * time.Second,
	}
	
	benchmarker := performance.NewBenchmarker(config)
	
	// 创建测试组件
	bilibiliWidget := widget.NewBilibiliVideosWidget()
	bilibiliWidget.UPMasters = []widget.BilibiliUPMaster{
		{UID: "123456", Name: "测试UP主"},
	}
	
	testFunc := func(ctx context.Context) error {
		_, err := bilibiliWidget.GetData(ctx, &mockConfig{})
		return err
	}
	
	result, err := benchmarker.RunBenchmark(context.Background(), testFunc)
	if err != nil {
		t.Fatalf("基准测试失败: %v", err)
	}
	
	benchmarker.PrintResult(result)
	
	// 验证性能指标
	if result.ErrorRate > 5.0 {
		t.Errorf("错误率过高: %.2f%%", result.ErrorRate)
	}
	
	if result.AverageTime > 2*time.Second {
		t.Errorf("平均响应时间过长: %v", result.AverageTime)
	}
	
	if result.RequestsPerSec < 1.0 {
		t.Errorf("每秒请求数过低: %.2f", result.RequestsPerSec)
	}
}

// TestZhihuWidgetBenchmark 知乎组件基准测试
func TestZhihuWidgetBenchmark(t *testing.T) {
	config := performance.BenchmarkConfig{
		Name:         "Zhihu Widget Load Test",
		Concurrency:  5,
		Duration:     20 * time.Second,
		RequestCount: 500,
		WarmupTime:   3 * time.Second,
	}
	
	benchmarker := performance.NewBenchmarker(config)
	
	zhihuWidget := widget.NewZhihuTrendingWidget()
	
	testFunc := func(ctx context.Context) error {
		_, err := zhihuWidget.GetData(ctx, &mockConfig{})
		return err
	}
	
	result, err := benchmarker.RunBenchmark(context.Background(), testFunc)
	if err != nil {
		t.Fatalf("基准测试失败: %v", err)
	}
	
	benchmarker.PrintResult(result)
	
	// 验证性能指标
	if result.ErrorRate > 10.0 {
		t.Errorf("错误率过高: %.2f%%", result.ErrorRate)
	}
}

// TestConcurrentWidgetLoad 并发组件加载测试
func TestConcurrentWidgetLoad(t *testing.T) {
	config := performance.BenchmarkConfig{
		Name:         "Concurrent Widget Load Test",
		Concurrency:  20,
		Duration:     60 * time.Second,
		WarmupTime:   10 * time.Second,
	}
	
	benchmarker := performance.NewBenchmarker(config)
	
	widgets := []widget.Widget{
		widget.NewBilibiliVideosWidget(),
		widget.NewZhihuTrendingWidget(),
		widget.NewGiteeReposWidget(),
	}
	
	testFunc := func(ctx context.Context) error {
		// 随机选择一个组件进行测试
		w := widgets[time.Now().UnixNano()%int64(len(widgets))]
		_, err := w.GetData(ctx, &mockConfig{})
		return err
	}
	
	result, err := benchmarker.RunBenchmark(context.Background(), testFunc)
	if err != nil {
		t.Fatalf("并发测试失败: %v", err)
	}
	
	benchmarker.PrintResult(result)
	
	// 验证并发性能
	if result.RequestsPerSec < 5.0 {
		t.Errorf("并发性能不足: %.2f requests/sec", result.RequestsPerSec)
	}
}

// TestMemoryUsage 内存使用测试
func TestMemoryUsage(t *testing.T) {
	monitor := performance.NewMonitor()
	
	// 启动监控
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	go monitor.Start(ctx, 1*time.Second)
	
	// 模拟高负载
	for i := 0; i < 1000; i++ {
		widget := widget.NewBilibiliVideosWidget()
		widget.GetData(ctx, &mockConfig{})
		
		if i%100 == 0 {
			metrics := monitor.GetMetrics()
			fmt.Printf("内存使用: %d MB, Goroutine数: %d\n", 
				metrics.MemoryUsage/1024/1024, metrics.GoroutineCount)
		}
	}
	
	// 检查最终指标
	finalMetrics := monitor.GetMetrics()
	
	if finalMetrics.MemoryUsage > 500*1024*1024 { // 500MB
		t.Errorf("内存使用过高: %d MB", finalMetrics.MemoryUsage/1024/1024)
	}
	
	if finalMetrics.GoroutineCount > 1000 {
		t.Errorf("Goroutine数量过多: %d", finalMetrics.GoroutineCount)
	}
}

// mockConfig 模拟配置
type mockConfig struct{}

func (m *mockConfig) GetString(key string) string {
	return ""
}

func (m *mockConfig) GetInt(key string) int {
	return 0
}

func (m *mockConfig) GetBool(key string) bool {
	return false
}

func (m *mockConfig) GetStringSlice(key string) []string {
	return []string{}
}
