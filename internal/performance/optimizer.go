package performance

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Optimizer 性能优化器
type Optimizer struct {
	config OptimizerConfig
	pools  map[string]*WorkerPool
	mu     sync.RWMutex
}

// OptimizerConfig 优化器配置
type OptimizerConfig struct {
	MaxWorkers        int           `yaml:"max-workers"`
	WorkerIdleTimeout time.Duration `yaml:"worker-idle-timeout"`
	QueueSize         int           `yaml:"queue-size"`
	GCPercent         int           `yaml:"gc-percent"`
	MaxMemoryMB       int           `yaml:"max-memory-mb"`
}

// WorkerPool 工作池
type WorkerPool struct {
	name       string
	workers    chan chan Job
	jobQueue   chan Job
	quit       chan bool
	maxWorkers int
	mu         sync.RWMutex
}

// Job 工作任务
type Job struct {
	ID       string
	Function func(ctx context.Context) error
	Context  context.Context
	Result   chan error
}

// NewOptimizer 创建性能优化器
func NewOptimizer(config OptimizerConfig) *Optimizer {
	// 设置默认值
	if config.MaxWorkers == 0 {
		config.MaxWorkers = runtime.NumCPU() * 2
	}
	if config.WorkerIdleTimeout == 0 {
		config.WorkerIdleTimeout = 30 * time.Second
	}
	if config.QueueSize == 0 {
		config.QueueSize = 1000
	}
	if config.GCPercent == 0 {
		config.GCPercent = 100
	}
	
	optimizer := &Optimizer{
		config: config,
		pools:  make(map[string]*WorkerPool),
	}
	
	// 应用系统优化
	optimizer.applySystemOptimizations()
	
	return optimizer
}

// applySystemOptimizations 应用系统级优化
func (o *Optimizer) applySystemOptimizations() {
	// 设置GC百分比
	if o.config.GCPercent > 0 {
		runtime.SetGCPercent(o.config.GCPercent)
	}
	
	// 设置最大CPU使用数
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	// 启动内存监控
	go o.monitorMemory()
}

// monitorMemory 监控内存使用
func (o *Optimizer) monitorMemory() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// 检查内存使用是否超过限制
		if o.config.MaxMemoryMB > 0 {
			currentMB := int(m.Alloc / 1024 / 1024)
			if currentMB > o.config.MaxMemoryMB {
				// 强制GC
				runtime.GC()
				runtime.GC() // 连续两次GC确保清理
			}
		}
	}
}

// GetOrCreatePool 获取或创建工作池
func (o *Optimizer) GetOrCreatePool(name string) *WorkerPool {
	o.mu.RLock()
	pool, exists := o.pools[name]
	o.mu.RUnlock()
	
	if exists {
		return pool
	}
	
	o.mu.Lock()
	defer o.mu.Unlock()
	
	// 双重检查
	if pool, exists := o.pools[name]; exists {
		return pool
	}
	
	pool = NewWorkerPool(name, o.config.MaxWorkers, o.config.QueueSize)
	pool.Start()
	o.pools[name] = pool
	
	return pool
}

// NewWorkerPool 创建工作池
func NewWorkerPool(name string, maxWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		name:       name,
		workers:    make(chan chan Job, maxWorkers),
		jobQueue:   make(chan Job, queueSize),
		quit:       make(chan bool),
		maxWorkers: maxWorkers,
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		worker := NewWorker(wp.workers, wp.quit)
		worker.Start()
	}
	
	go wp.dispatch()
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	close(wp.quit)
}

// Submit 提交任务
func (wp *WorkerPool) Submit(job Job) error {
	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// dispatch 分发任务
func (wp *WorkerPool) dispatch() {
	for {
		select {
		case job := <-wp.jobQueue:
			// 获取可用的worker
			select {
			case jobChannel := <-wp.workers:
				jobChannel <- job
			case <-wp.quit:
				return
			}
		case <-wp.quit:
			return
		}
	}
}

// Worker 工作者
type Worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

// NewWorker 创建工作者
func NewWorker(workerPool chan chan Job, quit chan bool) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       quit,
	}
}

// Start 启动工作者
func (w *Worker) Start() {
	go func() {
		for {
			// 将worker的job channel注册到worker pool
			w.workerPool <- w.jobChannel
			
			select {
			case job := <-w.jobChannel:
				// 执行任务
				err := job.Function(job.Context)
				if job.Result != nil {
					job.Result <- err
				}
			case <-w.quit:
				return
			}
		}
	}()
}

// 错误定义
var (
	ErrQueueFull = fmt.Errorf("job queue is full")
)

// ConnectionPool 连接池
type ConnectionPool struct {
	connections chan interface{}
	factory     func() (interface{}, error)
	close       func(interface{}) error
	maxSize     int
	mu          sync.Mutex
}

// NewConnectionPool 创建连接池
func NewConnectionPool(maxSize int, factory func() (interface{}, error), closeFunc func(interface{}) error) *ConnectionPool {
	return &ConnectionPool{
		connections: make(chan interface{}, maxSize),
		factory:     factory,
		close:       closeFunc,
		maxSize:     maxSize,
	}
}

// Get 获取连接
func (cp *ConnectionPool) Get() (interface{}, error) {
	select {
	case conn := <-cp.connections:
		return conn, nil
	default:
		return cp.factory()
	}
}

// Put 归还连接
func (cp *ConnectionPool) Put(conn interface{}) error {
	select {
	case cp.connections <- conn:
		return nil
	default:
		// 池已满，关闭连接
		return cp.close(conn)
	}
}

// Close 关闭连接池
func (cp *ConnectionPool) Close() error {
	close(cp.connections)
	
	for conn := range cp.connections {
		if err := cp.close(conn); err != nil {
			return err
		}
	}
	
	return nil
}
