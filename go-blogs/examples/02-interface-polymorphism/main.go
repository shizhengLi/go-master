package main

import "fmt"

// 接口与多态的魔法实现示例代码

// 定义接口
type Writer interface {
    Write([]byte) (int, error)
}

// 实现接口
type StringWriter struct {
    data string
}

func (s *StringWriter) Write(data []byte) (int, error) {
    s.data += string(data)
    return len(data), nil
}

func main() {
    fmt.Println("接口与多态的魔法实现示例")

    // 隐式实现接口
    var w Writer = &StringWriter{data: ""}
    w.Write([]byte("Hello, World!"))

    // 类型断言
    if sw, ok := w.(*StringWriter); ok {
        fmt.Printf("类型断言成功: %s\n", sw.data)
    }

    // 空接口
    var i interface{} = 42
    switch v := i.(type) {
    case int:
        fmt.Printf("i是整数: %d\n", v)
    case string:
        fmt.Printf("i是字符串: %s\n", v)
    default:
        fmt.Printf("i是未知类型: %T\n", v)
    }
}