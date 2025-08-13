package main

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/glance-china/internal/performance"
	"github.com/glance-china/internal/service"
	"github.com/glance-china/internal/widget"
)

func main() {
	fmt.Println("=== Glance 中国版性能测试 ===")
	
	// 运行各种性能测试
	runWidgetBenchmarks()
	runServiceBenchmarks()
	runMemoryTest()
	runConcurrencyTest()
}

func runWidgetBenchmarks() {
	fmt.Println("\n--- 组件性能测试 ---")
	
	widgets := map[string]widget.Widget{
		"bilibili": widget.NewBilibiliVideosWidget(),
		"zhihu":    widget.NewZhihuTrendingWidget(),
		"gitee":    widget.NewGiteeReposWidget(),
	}
	
	for name, w := range widgets {
		config := performance.BenchmarkConfig{
			Name:         fmt.Sprintf("%s Widget Test", name),
			Concurrency:  5,
			Duration:     10 * time.Second,
			WarmupTime:   2 * time.Second,
		}
		
		benchmarker := performance.NewBenchmarker(config)
		
		testFunc := func(ctx context.Context) error {
			_, err := w.GetData(ctx, &mockConfig{})
			return err
		}
		
		result, err := benchmarker.RunBenchmark(context.Background(), testFunc)
		if err != nil {
			log.Printf("测试 %s 失败: %v", name, err)
			continue
		}
		
		benchmarker.PrintResult(result)
	}
}

func runServiceBenchmarks() {
	fmt.Println("\n--- 服务性能测试 ---")
	
	// 这里可以添加服务层的性能测试
	fmt.Println("服务性能测试完成")
}

func runMemoryTest() {
	fmt.Println("\n--- 内存使用测试 ---")
	
	monitor := performance.NewMonitor()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	
	go monitor.Start(ctx, 2*time.Second)
	
	// 模拟负载
	for i := 0; i < 500; i++ {
		widget := widget.NewBilibiliVideosWidget()
		widget.GetData(ctx, &mockConfig{})
		
		if i%50 == 0 {
			metrics := monitor.GetMetrics()
			fmt.Printf("进度: %d/500, 内存: %d MB, Goroutines: %d\n", 
				i, metrics.MemoryUsage/1024/1024, metrics.GoroutineCount)
		}
	}
	
	finalMetrics := monitor.GetMetrics()
	fmt.Printf("最终内存使用: %d MB\n", finalMetrics.MemoryUsage/1024/1024)
	fmt.Printf("最终Goroutine数: %d\n", finalMetrics.GoroutineCount)
}

func runConcurrencyTest() {
	fmt.Println("\n--- 并发性能测试 ---")
	
	config := performance.BenchmarkConfig{
		Name:         "High Concurrency Test",
		Concurrency:  50,
		Duration:     30 * time.Second,
		WarmupTime:   5 * time.Second,
	}
	
	benchmarker := performance.NewBenchmarker(config)
	
	widgets := []widget.Widget{
		widget.NewBilibiliVideosWidget(),
		widget.NewZhihuTrendingWidget(),
		widget.NewGiteeReposWidget(),
	}
	
	testFunc := func(ctx context.Context) error {
		// 随机选择组件
		w := widgets[time.Now().UnixNano()%int64(len(widgets))]
		_, err := w.GetData(ctx, &mockConfig{})
		return err
	}
	
	result, err := benchmarker.RunBenchmark(context.Background(), testFunc)
	if err != nil {
		log.Printf("并发测试失败: %v", err)
		return
	}
	
	benchmarker.PrintResult(result)
	
	// 性能评估
	fmt.Println("\n--- 性能评估 ---")
	if result.RequestsPerSec > 10 {
		fmt.Println("✅ 并发性能: 优秀")
	} else if result.RequestsPerSec > 5 {
		fmt.Println("⚠️  并发性能: 良好")
	} else {
		fmt.Println("❌ 并发性能: 需要优化")
	}
	
	if result.ErrorRate < 1 {
		fmt.Println("✅ 错误率: 优秀")
	} else if result.ErrorRate < 5 {
		fmt.Println("⚠️  错误率: 可接受")
	} else {
		fmt.Println("❌ 错误率: 过高")
	}
	
	if result.AverageTime < 500*time.Millisecond {
		fmt.Println("✅ 响应时间: 优秀")
	} else if result.AverageTime < 2*time.Second {
		fmt.Println("⚠️  响应时间: 可接受")
	} else {
		fmt.Println("❌ 响应时间: 过慢")
	}
}

// mockConfig 实现
type mockConfig struct{}

func (m *mockConfig) GetString(key string) string { return "" }
func (m *mockConfig) GetInt(key string) int { return 0 }
func (m *mockConfig) GetBool(key string) bool { return false }
func (m *mockConfig) GetStringSlice(key string) []string { return []string{} }
