package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// Goroutine调度器底层原理示例代码

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d 开始\n", id)
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Worker %d 完成\n", id)
}

func demonstrateGMP() {
    fmt.Println("Goroutine调度器演示")

    // 设置GOMAXPROCS
    runtime.GOMAXPROCS(2)
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

    // 创建多个goroutine
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("所有worker完成")
}

// 抢占式调度演示
func preemptiveDemo() {
    fmt.Println("\n抢占式调度演示")

    runtime.GOMAXPROCS(1) // 单核观察抢占效果

    var wg sync.WaitGroup

    // 长时间运行的goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            if i%100 == 0 {
                fmt.Printf("长时间运行Goroutine %d\n", i)
            }
        }
    }()

    // 短时间运行的goroutine
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("短时间Goroutine %d\n", id)
            time.Sleep(50 * time.Millisecond)
        }(i)
    }

    wg.Wait()
}

func main() {
    demonstrateGMP()
    preemptiveDemo()

    fmt.Printf("\n最终goroutine数量: %d\n", runtime.NumGoroutine())
}