// 代码生成时间: 2025-10-23 06:38:42
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "log"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// HashCalculatorHandler 定义哈希值计算工具的处理器
func HashCalculatorHandler(c *fiber.Ctx) error {
    // 获取请求中的文本
    text := c.Query("text", "")

    // 检查输入是否为空
    if text == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "text parameter is required",
        })
    }

    // 计算SHA-256哈希值
    hashBytes := sha256.Sum256([]byte(text))
    hash := hex.EncodeToString(hashBytes[:])

    // 返回哈希值
    return c.JSON(fiber.Map{
        "text": text,
        "hash": hash,
    })
}

func main() {
    // 创建Fiber实例
    app := fiber.New()

    // 设置路由，响应GET请求
    app.Get("/hash", HashCalculatorHandler)

    // 启动服务器
    log.Fatal(app.Listen(":3000"))
}
