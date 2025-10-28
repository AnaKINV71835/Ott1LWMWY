// 代码生成时间: 2025-10-28 10:44:41
// api_test_tool.go
package main

import (
    "fmt"
    "net/http"
    "github.com/gofiber/fiber/v2"
)

// main 函数是程序的入口点
func main() {
    // 创建一个新的 Fiber 实例
    app := fiber.New()

    // 测试 API 路由
    app.Get("/test", func(c *fiber.Ctx) error {
        // 返回一个简单的响应
        return c.SendString("Hello, API Test Tool!")
    })

    // 处理错误，并启动服务器
    if err := app.Listen(":3000"); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
