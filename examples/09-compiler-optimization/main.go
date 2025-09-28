package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// 编译器优化与逃逸分析示例代码

func main() {
	fmt.Println("编译器优化与逃逸分析示例")

	// 逃逸分析示例
	escapeAnalysisExample()

	// 栈分配与堆分配对比
	stackVsHeapAllocation()

	// 内联优化演示
	inliningExample()

	// 常量折叠与传播
	constantFoldingExample()

	// 循环优化
	loopOptimizationExample()

	// 内存池优化
	memoryPoolExample()

	// 对象池优化
	objectPoolExample()

	// 预分配优化
	preallocationExample()

	// 并发优化
	concurrentOptimizationExample()

	// 性能测试
	performanceTests()

	// 编译器指令演示
	compilerDirectivesExample()

	// 内存对齐优化
	memoryAlignmentExample()
}

// 逃逸分析示例
func escapeAnalysisExample() {
	fmt.Println("\n=== 逃逸分析示例 ===")

	// 不逃逸的变量
	noEscape1()

	// 逃逸的变量
	escape1()

	// 复杂的逃逸分析
	complexEscape()
}

// 不逃逸的示例
func noEscape1() {
	x := 42
	fmt.Printf("不逃逸的变量: %d\n", x)
}

// 逃逸的示例
func escape1() *int {
	x := 42 // x会逃逸到堆上
	return &x
}

// 复杂的逃逸分析
func complexEscape() {
	// 闭包中的变量
	x := 42
	func() {
		fmt.Printf("闭包中的变量: %d\n", x) // x可能会逃逸
	}()

	// slice中的指针
	var pointers []*int
	y := 42
	pointers = append(pointers, &y) // y会逃逸

	// channel中的指针
	ch := make(chan *int)
	z := 42
	ch <- &z // z会逃逸
}

// 栈分配与堆分配对比
func stackVsHeapAllocation() {
	fmt.Println("\n=== 栈分配与堆分配对比 ===")

	// 栈分配示例
	stackAllocated := stackAllocation()
	fmt.Printf("栈分配结果: %d\n", stackAllocated)

	// 堆分配示例
	heapAllocated := heapAllocation()
	fmt.Printf("堆分配结果: %d\n", *heapAllocated)
}

func stackAllocation() int {
	x := 42 // 栈分配
	return x
}

func heapAllocation() *int {
	x := 42 // 堆分配
	return &x
}

// 内联优化演示
func inliningExample() {
	fmt.Println("\n=== 内联优化演示 ===")

	// 简单函数，容易被内联
	result := simpleAdd(5, 3)
	fmt.Printf("简单加法: %d\n", result)

	// 更复杂的内联情况
	data := []int{1, 2, 3, 4, 5}
	sum := calculateSum(data)
	fmt.Printf("计算总和: %d\n", sum)
}

// 简单函数，容易被内联
func simpleAdd(a, b int) int {
	return a + b
}

// 复杂一点的函数，也可能被内联
func calculateSum(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

// 常量折叠与传播
func constantFoldingExample() {
	fmt.Println("\n=== 常量折叠与传播 ===")

	// 常量折叠
	x := 10 + 20 // 编译器会优化为 x := 30
	y := 5 * 4   // 编译器会优化为 y := 20
	fmt.Printf("常量折叠: x=%d, y=%d\n", x, y)

	// 常量传播
	const debug = true
	if debug {
		fmt.Println("调试模式开启") // 这段代码会被保留
	}

	const production = false
	if production {
		fmt.Println("生产模式开启") // 这段代码会被消除
	}
}

// 循环优化
func loopOptimizationExample() {
	fmt.Println("\n=== 循环优化 ===")

	// 原始循环
	data := []int{1, 2, 3, 4, 5}
	sum := optimizedLoop(data)
	fmt.Printf("优化循环总和: %d\n", sum)

	// 循环展开
	sum2 := unrolledLoop(data)
	fmt.Printf("循环展开总和: %d\n", sum2)

	// 边界检查消除
	sum3 := boundsCheckEliminated(data)
	fmt.Printf("边界检查消除总和: %d\n", sum3)
}

// 优化的循环
func optimizedLoop(data []int) int {
	sum := 0
	// 使用range消除边界检查
	for _, v := range data {
		sum += v
	}
	return sum
}

// 循环展开
func unrolledLoop(data []int) int {
	sum := 0
	n := len(data)

	// 手动展开4次
	i := 0
	for i <= n-4 {
		sum += data[i] + data[i+1] + data[i+2] + data[i+3]
		i += 4
	}

	// 处理剩余元素
	for i < n {
		sum += data[i]
		i++
	}

	return sum
}

// 边界检查消除
func boundsCheckEliminated(data []int) int {
	if len(data) == 0 {
		return 0
	}

	// 预先检查边界
	if len(data) > 0 {
		// 编译器可以消除这些访问的边界检查
		return data[0] + data[len(data)-1]
	}
	return 0
}

// 内存池优化
type MemoryPool struct {
	pool chan []byte
	size int
}

func NewMemoryPool(size, capacity int) *MemoryPool {
	pool := &MemoryPool{
		pool: make(chan []byte, capacity),
		size: size,
	}

	// 预填充池
	for i := 0; i < capacity; i++ {
		pool.pool <- make([]byte, size)
	}

	return pool
}

func (mp *MemoryPool) Get() []byte {
	select {
	case buf := <-mp.pool:
		return buf
	default:
		return make([]byte, mp.size)
	}
}

func (mp *MemoryPool) Put(buf []byte) {
	if cap(buf) >= mp.size {
		select {
		case mp.pool <- buf[:mp.size]:
		default:
			// 池已满，丢弃
		}
	}
}

func memoryPoolExample() {
	fmt.Println("\n=== 内存池优化 ===")

	pool := NewMemoryPool(1024, 10)

	// 使用内存池
	buf1 := pool.Get()
	copy(buf1, []byte("Hello, Memory Pool!"))
	fmt.Printf("从池中获取: %s\n", buf1)

	// 放回池中
	pool.Put(buf1)

	// 再次获取
	buf2 := pool.Get()
	fmt.Printf("再次获取: %s\n", buf2)
}

// 对象池优化
type LargeObject struct {
	data []int
}

func (lo *LargeObject) Reset() {
	for i := range lo.data {
		lo.data[i] = 0
	}
}

type ObjectPool struct {
	pool sync.Pool
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &LargeObject{
					data: make([]int, 1000),
				}
			},
		},
	}
}

func (op *ObjectPool) Get() *LargeObject {
	return op.pool.Get().(*LargeObject)
}

func (op *ObjectPool) Put(obj *LargeObject) {
	obj.Reset()
	op.pool.Put(obj)
}

func objectPoolExample() {
	fmt.Println("\n=== 对象池优化 ===")

	pool := NewObjectPool()

	// 使用对象池
	obj1 := pool.Get()
	obj1.data[0] = 42
	fmt.Printf("对象池对象数据[0]: %d\n", obj1.data[0])

	// 放回池中
	pool.Put(obj1)

	// 再次获取
	obj2 := pool.Get()
	fmt.Printf("重用对象数据[0]: %d\n", obj2.data[0])
}

// 预分配优化
func preallocationExample() {
	fmt.Println("\n=== 预分配优化 ===")

	// 动态分配vs预分配
	fmt.Println("Slice预分配测试:")
	benchmarkSliceAllocation()

	fmt.Println("\nMap预分配测试:")
	benchmarkMapAllocation()
}

func benchmarkSliceAllocation() {
	const size = 10000

	// 动态分配
	start := time.Now()
	var dynamicSlice []int
	for i := 0; i < size; i++ {
		dynamicSlice = append(dynamicSlice, i)
	}
	dynamicTime := time.Since(start)

	// 预分配
	start = time.Now()
	preallocatedSlice := make([]int, 0, size)
	for i := 0; i < size; i++ {
		preallocatedSlice = append(preallocatedSlice, i)
	}
	preallocatedTime := time.Since(start)

	fmt.Printf("动态分配: %v\n", dynamicTime)
	fmt.Printf("预分配: %v\n", preallocatedTime)
	fmt.Printf("性能提升: %.2fx\n", float64(dynamicTime)/float64(preallocatedTime))
}

func benchmarkMapAllocation() {
	const size = 10000

	// 动态分配
	start := time.Now()
	dynamicMap := make(map[int]string)
	for i := 0; i < size; i++ {
		dynamicMap[i] = fmt.Sprintf("value%d", i)
	}
	dynamicTime := time.Since(start)

	// 预分配
	start = time.Now()
	preallocatedMap := make(map[int]string, size)
	for i := 0; i < size; i++ {
		preallocatedMap[i] = fmt.Sprintf("value%d", i)
	}
	preallocatedTime := time.Since(start)

	fmt.Printf("动态分配: %v\n", dynamicTime)
	fmt.Printf("预分配: %v\n", preallocatedTime)
	fmt.Printf("性能提升: %.2fx\n", float64(dynamicTime)/float64(preallocatedTime))
}

// 并发优化
type OptimizedWorkerPool struct {
	tasks   chan Task
	results chan Result
	workers int
	wg      sync.WaitGroup
}

type Task struct {
	ID   int
	Data int
}

type Result struct {
	TaskID int
	Value  int
}

func NewOptimizedWorkerPool(workers int) *OptimizedWorkerPool {
	pool := &OptimizedWorkerPool{
		tasks:   make(chan Task, workers*2),
		results: make(chan Result, workers*2),
		workers: workers,
	}

	// 启动worker
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func (wp *OptimizedWorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		result := wp.processTask(task)
		wp.results <- result
	}
}

func (wp *OptimizedWorkerPool) processTask(task Task) Result {
	// 模拟处理
	time.Sleep(time.Millisecond)
	return Result{
		TaskID: task.ID,
		Value:  task.Data * 2,
	}
}

func (wp *OptimizedWorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *OptimizedWorkerPool) Stop() {
	close(wp.tasks)
	wp.wg.Wait()
	close(wp.results)
}

func concurrentOptimizationExample() {
	fmt.Println("\n=== 并发优化 ===")

	pool := NewOptimizedWorkerPool(4)

	// 提交任务
	for i := 0; i < 10; i++ {
		task := Task{
			ID:   i,
			Data: i,
		}
		pool.Submit(task)
	}

	// 收集结果
	go func() {
		pool.Stop()
	}()

	for result := range pool.results {
		fmt.Printf("任务 %d 结果: %d\n", result.TaskID, result.Value)
	}
}

// 性能测试
func performanceTests() {
	fmt.Println("\n=== 性能测试 ===")

	// 内存分配测试
	benchmarkMemoryAllocation()

	// 逃逸分析测试
	benchmarkEscapeAnalysis()

	// 内联优化测试
	benchmarkInlining()
}

func benchmarkMemoryAllocation() {
	fmt.Println("内存分配基准测试:")

	const iterations = 10000

	// 栈分配测试
	start := time.Now()
	for i := 0; i < iterations; i++ {
		stackAllocation()
	}
	stackTime := time.Since(start)

	// 堆分配测试
	start = time.Now()
	for i := 0; i < iterations; i++ {
		heapAllocation()
	}
	heapTime := time.Since(start)

	fmt.Printf("栈分配: %v\n", stackTime)
	fmt.Printf("堆分配: %v\n", heapTime)
	fmt.Printf("性能差异: %.2fx\n", float64(heapTime)/float64(stackTime))
}

func benchmarkEscapeAnalysis() {
	fmt.Println("逃逸分析基准测试:")

	const iterations = 100000

	// 不逃逸的函数调用
	start := time.Now()
	for i := 0; i < iterations; i++ {
		noEscape1()
	}
	noEscapeTime := time.Since(start)

	// 逃逸的函数调用
	start = time.Now()
	for i := 0; i < iterations; i++ {
		escape1()
	}
	escapeTime := time.Since(start)

	fmt.Printf("不逃逸调用: %v\n", noEscapeTime)
	fmt.Printf("逃逸调用: %v\n", escapeTime)
	fmt.Printf("性能差异: %.2fx\n", float64(escapeTime)/float64(noEscapeTime))
}

func benchmarkInlining() {
	fmt.Println("内联优化基准测试:")

	const iterations = 1000000

	// 简单函数调用
	start := time.Now()
	for i := 0; i < iterations; i++ {
		simpleAdd(i, i+1)
	}
	simpleTime := time.Since(start)

	// 复杂函数调用
	start = time.Now()
	for i := 0; i < iterations; i++ {
		complexFunction(i)
	}
	complexTime := time.Since(start)

	fmt.Printf("简单函数: %v\n", simpleTime)
	fmt.Printf("复杂函数: %v\n", complexTime)
	fmt.Printf("性能差异: %.2fx\n", float64(complexTime)/float64(simpleTime))
}

func complexFunction(x int) int {
	// 更复杂的计算，不太可能被内联
	result := 0
	for i := 0; i < 10; i++ {
		result += x * i
	}
	return result
}

// 编译器指令演示
//go:noinline
func noInlineFunction() int {
	return 42
}

func compilerDirectivesExample() {
	fmt.Println("\n=== 编译器指令演示 ===")

	// 调用noinline函数
	result := noInlineFunction()
	fmt.Printf("noinline函数结果: %d\n", result)

	// 内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("内存统计:\n")
	fmt.Printf("  分配内存: %d bytes\n", m.Alloc)
	fmt.Printf("  总分配: %d bytes\n", m.TotalAlloc)
	fmt.Printf("  系统内存: %d bytes\n", m.Sys)
	fmt.Printf("  GC次数: %d\n", m.NumGC)
}

// 内存对齐优化
type AlignedStruct struct {
	// 按照内存对齐规则排列
	a int64   // 8 bytes
	b int32   // 4 bytes
	c int16   // 2 bytes
	d int8    // 1 byte
	_ [5]byte // 填充到16字节对齐
}

type UnalignedStruct struct {
	a int8  // 1 byte
	b int64 // 8 bytes
	c int16 // 2 bytes
	d int32 // 4 bytes
}

func memoryAlignmentExample() {
	fmt.Println("\n=== 内存对齐优化 ===")

	aligned := AlignedStruct{
		a: 1,
		b: 2,
		c: 3,
		d: 4,
	}

	unaligned := UnalignedStruct{
		a: 1,
		b: 2,
		c: 3,
		d: 4,
	}

	fmt.Printf("对齐结构体大小: %d bytes\n", unsafe.Sizeof(aligned))
	fmt.Printf("非对齐结构体大小: %d bytes\n", unsafe.Sizeof(unaligned))
	fmt.Printf("内存节省: %d bytes\n", unsafe.Sizeof(unaligned)-unsafe.Sizeof(aligned))

	// 位域优化示例
	fmt.Println("\n位域优化示例:")
	bitFieldExample()
}

// 位域优化
type BitField struct {
	value uint32
}

func (bf *BitField) SetField(offset, width, value uint32) {
	mask := uint32(1<<width-1) << offset
	bf.value = (bf.value & ^mask) | ((value << offset) & mask)
}

func (bf *BitField) GetField(offset, width uint32) uint32 {
	mask := uint32(1<<width-1) << offset
	return (bf.value & mask) >> offset
}

func bitFieldExample() {
	bf := BitField{}

	// 存储多个小字段
	bf.SetField(0, 8, 0x12)   // 字段0: 8位
	bf.SetField(8, 8, 0x34)   // 字段1: 8位
	bf.SetField(16, 16, 0x5678) // 字段2: 16位

	fmt.Printf("位域值: %08x\n", bf.value)
	fmt.Printf("字段0: %02x\n", bf.GetField(0, 8))
	fmt.Printf("字段1: %02x\n", bf.GetField(8, 8))
	fmt.Printf("字段2: %04x\n", bf.GetField(16, 16))

	// 传统方式需要3个uint32，现在只需要1个
	fmt.Printf("内存节省: %d bytes\n", 8) // 3*4 - 1*4 = 8 bytes
}