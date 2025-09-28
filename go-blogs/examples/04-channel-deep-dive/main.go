package main

import (
    "fmt"
    "sync"
    "time"
)

// Channel深度剖析示例代码

func main() {
    fmt.Println("Channel深度剖析示例")

    // 基本Channel操作
    basicChannelOperations()

    // 生产者消费者模式
    producerConsumerDemo()

    // Select机制
    selectDemo()

    // Channel性能测试
    channelPerformanceTest()
}

// 基本Channel操作
func basicChannelOperations() {
    fmt.Println("\n=== 基本Channel操作 ===")

    // 无缓冲Channel
    unbuffered := make(chan int)
    go func() {
        unbuffered <- 42
    }()
    value := <-unbuffered
    fmt.Printf("无缓冲Channel接收到: %d\n", value)

    // 有缓冲Channel
    buffered := make(chan int, 3)
    buffered <- 1
    buffered <- 2
    buffered <- 3
    fmt.Printf("有缓冲Channel长度: %d\n", len(buffered))

    // 只发送和只接收Channel
    var sendOnly chan<- int = buffered
    var recvOnly <-chan int = buffered

    sendOnly <- 4
    received := <-recvOnly
    fmt.Printf("类型化Channel接收到: %d\n", received)
}

// 生产者消费者模式
func producerConsumerDemo() {
    fmt.Println("\n=== 生产者消费者模式 ===")

    tasks := make(chan int, 10)
    results := make(chan int, 10)

    // 启动消费者
    const numWorkers = 3
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for task := range tasks {
                result := task * 2
                results <- result
                fmt.Printf("Worker %d 处理任务 %d -> %d\n", workerID, task, result)
                time.Sleep(100 * time.Millisecond)
            }
        }(i)
    }

    // 启动生产者
    go func() {
        for i := 1; i <= 5; i++ {
            tasks <- i
            fmt.Printf("生产任务 %d\n", i)
        }
        close(tasks)
    }()

    // 收集结果
    go func() {
        wg.Wait()
        close(results)
    }()

    // 处理结果
    for result := range results {
        fmt.Printf("收到结果: %d\n", result)
    }
}

// Select机制
func selectDemo() {
    fmt.Println("\n=== Select机制 ===")

    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from ch1"
    }()

    go func() {
        time.Sleep(50 * time.Millisecond)
        ch2 <- "from ch2"
    }()

    select {
    case msg1 := <-ch1:
        fmt.Printf("从ch1接收到: %s\n", msg1)
    case msg2 := <-ch2:
        fmt.Printf("从ch2接收到: %s\n", msg2)
    case <-time.After(200 * time.Millisecond):
        fmt.Println("超时")
    }

    // Select with for loop
    fmt.Println("\nSelect with for loop:")
    done := make(chan bool)
    messages := make(chan string, 2)

    go func() {
        messages <- "message 1"
        messages <- "message 2"
        time.Sleep(100 * time.Millisecond)
        done <- true
    }()

    for {
        select {
        case msg := <-messages:
            fmt.Printf("处理消息: %s\n", msg)
        case <-done:
            fmt.Println("收到完成信号")
            return
        default:
            fmt.Println("等待消息...")
            time.Sleep(50 * time.Millisecond)
        }
    }
}

// Channel性能测试
func channelPerformanceTest() {
    fmt.Println("\n=== Channel性能测试 ===")

    iterations := 10000

    // 测试无缓冲Channel
    start := time.Now()
    testUnbufferedChannel(iterations)
    unbufferedTime := time.Since(start)

    // 测试有缓冲Channel
    start = time.Now()
    testBufferedChannel(iterations)
    bufferedTime := time.Since(start)

    fmt.Printf("无缓冲Channel: %v\n", unbufferedTime)
    fmt.Printf("有缓冲Channel: %v\n", bufferedTime)
    fmt.Printf("性能差异: %.2fx\n", float64(unbufferedTime)/float64(bufferedTime))
}

func testUnbufferedChannel(iterations int) {
    ch := make(chan int)
    done := make(chan bool)

    go func() {
        for i := 0; i < iterations; i++ {
            ch <- i
        }
        close(ch)
    }()

    go func() {
        for range ch {
            // 接收数据
        }
        done <- true
    }()

    <-done
}

func testBufferedChannel(iterations int) {
    ch := make(chan int, 100)
    done := make(chan bool)

    go func() {
        for i := 0; i < iterations; i++ {
            ch <- i
        }
        close(ch)
    }()

    go func() {
        for range ch {
            // 接收数据
        }
        done <- true
    }()

    <-done
}