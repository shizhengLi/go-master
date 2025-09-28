package main

import (
    "fmt"
    "runtime"
    "time"
)

// GC三色标记法示例代码

func main() {
    fmt.Println("GC三色标记法示例")

    // 基本GC测试
    basicGCTest()

    // 三色标记模拟
    simulateThreeColorMarking()

    // GC性能监控
    gcMonitoring()

    // 内存优化演示
    memoryOptimizationDemo()
}

// 基本GC测试
func basicGCTest() {
    fmt.Println("\n=== 基本GC测试 ===")

    // 获取GC统计
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)

    fmt.Printf("GC次数: %d\n", stats.NumGC)
    fmt.Printf("内存分配: %d bytes\n", stats.Alloc)

    // 创建内存压力
    createMemoryPressure()

    // 再次获取GC统计
    runtime.ReadMemStats(&stats)
    fmt.Printf("GC次数: %d\n", stats.NumGC)
    fmt.Printf("内存分配: %d bytes\n", stats.Alloc)
}

// 创建内存压力
func createMemoryPressure() {
    const iterations = 10000

    for i := 0; i < iterations; i++ {
        data := make([]byte, 1024)
        _ = data // 使用数据避免编译器优化
    }
}

// 三色标记模拟
type Color int

const (
    White Color = iota
    Gray
    Black
)

type GCObject struct {
    id       int
    color    Color
    children []*GCObject
}

func simulateThreeColorMarking() {
    fmt.Println("\n=== 三色标记模拟 ===")

    // 创建对象图
    root := &GCObject{id: 1, color: White}
    child1 := &GCObject{id: 2, color: White}
    child2 := &GCObject{id: 3, color: White}

    root.children = []*GCObject{child1, child2}

    // 模拟标记过程
    fmt.Println("开始三色标记:")
    markObject(root)

    // 显示结果
    fmt.Printf("根对象: %d (颜色: %v)\n", root.id, root.color)
    fmt.Printf("子对象1: %d (颜色: %v)\n", child1.id, child1.color)
    fmt.Printf("子对象2: %d (颜色: %v)\n", child2.id, child2.color)
}

func markObject(obj *GCObject) {
    if obj.color != White {
        return
    }

    // 变为灰色
    obj.color = Gray
    fmt.Printf("标记对象 %d 为灰色\n", obj.id)

    // 标记子对象
    for _, child := range obj.children {
        markObject(child)
    }

    // 变为黑色
    obj.color = Black
    fmt.Printf("标记对象 %d 为黑色\n", obj.id)
}

// GC监控
func gcMonitoring() {
    fmt.Println("\n=== GC监控 ===")

    // 设置GOGC
    oldGOGC := runtime.GOMAXPROCS(0)
    runtime.GOMAXPROCS(1) // 单核测试

    gogcValues := []int{50, 100, 200}

    for _, gogc := range gogcValues {
        fmt.Printf("\n测试 GOGC=%d:\n", gogc)

        // 设置GOGC
        debug.SetGCPercent(gogc)

        start := time.Now()
        createMemoryPressure()
        duration := time.Since(start)

        var stats runtime.MemStats
        runtime.ReadMemStats(&stats)

        fmt.Printf("耗时: %v\n", duration)
        fmt.Printf("GC次数: %d\n", stats.NumGC)
        fmt.Printf("暂停时间: %.2f ms\n", float64(stats.PauseTotalNs)/1e6)
    }

    // 恢复默认
    debug.SetGCPercent(oldGOGC)
}

// 内存优化演示
func memoryOptimizationDemo() {
    fmt.Println("\n=== 内存优化演示 ===")

    // 1. 对象池优化
    demonstrateObjectPool()

    // 2. 预分配优化
    demonstratePreAllocation()
}

// 对象池优化
func demonstrateObjectPool() {
    fmt.Println("对象池优化:")

    type LargeObject struct {
        data [1024]byte // 1KB
    }

    // 创建对象池
    pool := &sync.Pool{
        New: func() interface{} {
            return &LargeObject{}
        },
    }

    // 不使用池
    start := time.Now()
    for i := 0; i < 10000; i++ {
        obj := &LargeObject{}
        _ = obj
    }
    withoutPool := time.Since(start)

    // 使用池
    start = time.Now()
    for i := 0; i < 10000; i++ {
        obj := pool.Get().(*LargeObject)
        pool.Put(obj)
    }
    withPool := time.Since(start)

    fmt.Printf("不使用池: %v\n", withoutPool)
    fmt.Printf("使用池: %v\n", withPool)
    fmt.Printf("性能提升: %.2fx\n", float64(withoutPool)/float64(withPool))
}

// 预分配优化
func demonstratePreAllocation() {
    fmt.Println("预分配优化:")

    // 动态分配
    start := time.Now()
    var dynamicSlice []int
    for i := 0; i < 100000; i++ {
        dynamicSlice = append(dynamicSlice, i)
    }
    dynamicTime := time.Since(start)

    // 预分配
    start = time.Now()
    preallocated := make([]int, 0, 100000)
    for i := 0; i < 100000; i++ {
        preallocated = append(preallocated, i)
    }
    preallocatedTime := time.Since(start)

    fmt.Printf("动态分配: %v\n", dynamicTime)
    fmt.Printf("预分配: %v\n", preallocatedTime)
    fmt.Printf("性能提升: %.2fx\n", float64(dynamicTime)/float64(preallocatedTime))
}