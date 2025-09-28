package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Context机制深度解析示例代码

func main() {
	fmt.Println("Context机制深度解析示例")

	// 基本Context操作
	basicContextOperations()

	// 取消传播演示
	cancellationPropagation()

	// 超时控制演示
	timeoutControl()

	// Context值存储演示
	valueStorage()

	// 并发模式演示
	concurrentPatterns()

	// 性能测试
	performanceTests()

	// Web服务示例
	webServiceExample()
}

// 基本Context操作
func basicContextOperations() {
	fmt.Println("\n=== 基本Context操作 ===")

	// 创建根Context
	rootCtx := context.Background()
	fmt.Printf("根Context: %v\n", rootCtx)

	// TODO Context
	todoCtx := context.TODO()
	fmt.Printf("TODO Context: %v\n", todoCtx)

	// 可取消Context
	ctx, cancel := context.WithCancel(rootCtx)
	fmt.Printf("可取消Context: %v\n", ctx)
	defer cancel()

	// 检查Context状态
	fmt.Printf("Context错误: %v\n", ctx.Err())
}

// 取消传播演示
func cancellationPropagation() {
	fmt.Println("\n=== 取消传播演示 ===")

	// 创建父Context
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// 创建子Context
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()

	// 启动worker
	done := make(chan struct{})
	go func() {
		defer close(done)
		worker(childCtx, "Worker 1")
	}()

	// 等待一段时间后取消父Context
	time.Sleep(100 * time.Millisecond)
	fmt.Println("取消父Context...")
	parentCancel()

	// 等待worker完成
	<-done
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: 收到取消信号: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s: 正在工作...\n", name)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// 超时控制演示
func timeoutControl() {
	fmt.Println("\n=== 超时控制演示 ===")

	// WithTimeout示例
	fmt.Println("WithTimeout示例:")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		fmt.Printf("超时: %v\n", timeoutCtx.Err())
	case <-time.After(500 * time.Millisecond):
		fmt.Println("不应该到达这里")
	}

	// WithDeadline示例
	fmt.Println("\nWithDeadline示例:")
	deadline := time.Now().Add(300 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	if d, ok := deadlineCtx.Deadline(); ok {
		fmt.Printf("截止时间: %v\n", d)
		fmt.Printf("剩余时间: %v\n", time.Until(d))
	}

	select {
	case <-deadlineCtx.Done():
		fmt.Printf("截止时间到达: %v\n", deadlineCtx.Err())
	}
}

// Context值存储演示
func valueStorage() {
	fmt.Println("\n=== Context值存储演示 ===")

	// 定义Context键
	type contextKey string
	userIDKey := contextKey("userID")
	requestIDKey := contextKey("requestID")

	// 创建带有值的Context
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user123")
	ctx = context.WithValue(ctx, requestIDKey, "req456")

	// 获取值
	if userID := ctx.Value(userIDKey); userID != nil {
		fmt.Printf("用户ID: %v\n", userID)
	}

	if requestID := ctx.Value(requestIDKey); requestID != nil {
		fmt.Printf("请求ID: %v\n", requestID)
	}

	// 链式存储
	ctx2 := context.WithValue(ctx, contextKey("traceID"), "trace789")
	fmt.Printf("链式存储 - 追踪ID: %v\n", ctx2.Value(contextKey("traceID")))
	fmt.Printf("链式存储 - 用户ID: %v\n", ctx2.Value(userIDKey)) // 仍然可以获取
}

// 并发模式演示
func concurrentPatterns() {
	fmt.Println("\n=== 并发模式演示 ===")

	// 并发任务组
	fmt.Println("并发任务组示例:")
	err := concurrentTaskGroup()
	if err != nil {
		fmt.Printf("任务组错误: %v\n", err)
	}

	// 工作池模式
	fmt.Println("\n工作池模式示例:")
	workPoolExample()
}

// 并发任务组
type TaskGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	err    error
	errMu  sync.Mutex
}

func NewTaskGroup(ctx context.Context) *TaskGroup {
	childCtx, cancel := context.WithCancel(ctx)
	return &TaskGroup{
		ctx:    childCtx,
		cancel: cancel,
	}
}

func (tg *TaskGroup) Go(fn func(context.Context) error) {
	tg.wg.Add(1)
	go func() {
		defer tg.wg.Done()

		if err := fn(tg.ctx); err != nil {
			tg.errMu.Lock()
			if tg.err == nil {
				tg.err = err
				tg.cancel() // 取消其他任务
			}
			tg.errMu.Unlock()
		}
	}()
}

func (tg *TaskGroup) Wait() error {
	tg.wg.Wait()
	return tg.err
}

func concurrentTaskGroup() error {
	ctx := context.Background()
	tg := NewTaskGroup(ctx)

	// 启动多个任务
	tg.Go(func(ctx context.Context) error {
		return simulateTask(ctx, "Task 1", 100*time.Millisecond)
	})

	tg.Go(func(ctx context.Context) error {
		return simulateTask(ctx, "Task 2", 150*time.Millisecond)
	})

	tg.Go(func(ctx context.Context) error {
		return simulateTask(ctx, "Task 3", 200*time.Millisecond)
	})

	return tg.Wait()
}

func simulateTask(ctx context.Context, name string, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: 被取消: %w", name, ctx.Err())
	case <-time.After(duration):
		fmt.Printf("%s: 完成\n", name)
		return nil
	}
}

// 工作池示例
func workPoolExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tasks := make(chan int, 10)
	results := make(chan int, 10)

	// 创建工作池
	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		go workerPool(ctx, i, tasks, results)
	}

	// 发送任务
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case tasks <- i:
				fmt.Printf("发送任务 %d\n", i)
			case <-ctx.Done():
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		close(tasks)
	}()

	// 收集结果
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case result := <-results:
				fmt.Printf("收到结果: %d\n", result)
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(3 * time.Second)
}

func workerPool(ctx context.Context, id int, tasks <-chan int, results chan<- int) {
	for task := range tasks {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: 停止工作\n", id)
			return
		default:
			result := task * 2
			fmt.Printf("Worker %d: 处理任务 %d -> %d\n", id, task, result)
			results <- result
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// 性能测试
func performanceTests() {
	fmt.Println("\n=== 性能测试 ===")

	// Context创建性能
	fmt.Println("Context创建性能:")
	benchmarkContextCreation()

	// Context传播性能
	fmt.Println("\nContext传播性能:")
	benchmarkContextPropagation()
}

func benchmarkContextCreation() {
	const iterations = 10000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		ctx := context.Background()
		_ = ctx
	}
	backgroundTime := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_ = ctx
	}
	cancelTime := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		ctx := context.WithValue(context.Background(), "key", "value")
		_ = ctx
	}
	valueTime := time.Since(start)

	fmt.Printf("Background: %v\n", backgroundTime)
	fmt.Printf("WithCancel: %v\n", cancelTime)
	fmt.Printf("WithValue: %v\n", valueTime)
}

func benchmarkContextPropagation() {
	const iterations = 100000

	ctx := context.WithValue(context.Background(), "key", "value")

	start := time.Now()
	for i := 0; i < iterations; i++ {
		value := ctx.Value("key")
		_ = value
	}
	propagationTime := time.Since(start)

	fmt.Printf("Value传播: %v\n", propagationTime)
}

// Web服务示例
func webServiceExample() {
	fmt.Println("\n=== Web服务示例 ===")

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: createHandler(),
	}

	// 启动服务器
	go func() {
		fmt.Println("启动HTTP服务器 :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器错误: %v\n", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 发送测试请求
	sendTestRequest()

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("关闭服务器...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭错误: %v\n", err)
	}

	fmt.Println("服务器已关闭")
}

func createHandler() http.Handler {
	mux := http.NewServeMux()

	// 注册中间件
	handler := withMiddleware(mux)

	// 注册路由
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/api", handleAPI)
	mux.HandleFunc("/timeout", handleTimeout)

	return handler
}

// 中间件
func withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建Context
		ctx := r.Context()

		// 添加请求ID
		requestID := generateRequestID()
		ctx = context.WithValue(ctx, "requestID", requestID)

		// 添加超时
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// 添加日志
		logger := log.New(os.Stdout, "", log.LstdFlags)
		ctx = context.WithValue(ctx, "logger", logger)

		// 更新请求
		r = r.WithContext(ctx)

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		// 记录请求时间
		duration := time.Since(start)
		logger.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value("requestID").(string)
	logger := ctx.Value("logger").(*log.Logger)

	logger.Printf("处理首页请求 - RequestID: %s", requestID)
	fmt.Fprintf(w, "欢迎访问首页! RequestID: %s", requestID)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value("requestID").(string)

	// 模拟API处理
	select {
	case <-ctx.Done():
		http.Error(w, "请求超时", http.StatusRequestTimeout)
		return
	case <-time.After(100 * time.Millisecond):
		// 处理完成
	}

	response := map[string]interface{}{
		"status":    "success",
		"requestID": requestID,
		"data":      "API响应数据",
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success","requestID":"%s","data":"API响应数据"}`, requestID)
}

func handleTimeout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 模拟长时间操作
	select {
	case <-ctx.Done():
		http.Error(w, "操作超时", http.StatusRequestTimeout)
		return
	case <-time.After(3 * time.Second):
		// 不应该到达这里
	}

	fmt.Fprint(w, "操作完成")
}

func sendTestRequest() {
	client := &http.Client{Timeout: 10 * time.Second}

	// 测试首页
	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		log.Printf("请求首页失败: %v", err)
		return
	}
	resp.Body.Close()

	// 测试API
	resp, err = client.Get("http://localhost:8080/api")
	if err != nil {
		log.Printf("请求API失败: %v", err)
		return
	}
	resp.Body.Close()

	// 测试超时
	resp, err = client.Get("http://localhost:8080/timeout")
	if err != nil {
		log.Printf("请求超时测试失败: %v", err)
		return
	}
	resp.Body.Close()

	time.Sleep(500 * time.Millisecond)
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}