package main

import "fmt"

// Go语言内存模型与指针哲学示例代码
func main() {
    fmt.Println("Go语言内存模型与指针哲学示例")

    // 基本指针操作
    x := 42
    p := &x
    fmt.Printf("x = %d, p = %p, *p = %d\n", x, p, *p)

    // 指针接收器 vs 值接收器
    type Counter struct {
        count int
    }

    func (c Counter) IncrementValue() {
        c.count++
    }

    func (c *Counter) IncrementPointer() {
        c.count++
    }

    counter := Counter{count: 0}
    counter.IncrementValue()
    fmt.Printf("值接收器后: %d\n", counter.count) // 0

    counter.IncrementPointer()
    fmt.Printf("指针接收器后: %d\n", counter.count) // 1
}