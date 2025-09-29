// 代码生成时间: 2025-09-29 21:21:35
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "flag"
    "fmt"
    "io"
    "io/fs"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/gofiber/fiber/v2"
)

// FileIntegrityChecker 结构体用于存储文件路径和预期的哈希值
type FileIntegrityChecker struct {
    Path    string
    ExpectedHash string
}

// NewFileIntegrityChecker 创建一个新的 FileIntegrityChecker 实例
func NewFileIntegrityChecker(path string, expectedHash string) *FileIntegrityChecker {
    return &FileIntegrityChecker{
        Path:    path,
        ExpectedHash: expectedHash,
    }
}

// CheckIntegrity 检查文件的完整性是否与预期的哈希值匹配
func (fic *FileIntegrityChecker) CheckIntegrity() error {
    // 打开文件
    file, err := os.Open(fic.Path)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // 计算文件的哈希值
    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return fmt.Errorf("failed to calculate hash: %w", err)
    }
    actualHash := hex.EncodeToString(hash.Sum(nil))

    // 比较实际哈希值和预期哈希值
    if actualHash != fic.ExpectedHash {
        return fmt.Errorf("file integrity check failed, expected hash: %s, actual hash: %s", fic.ExpectedHash, actualHash)
    }

    return nil
}

// fileIntegrityCheckHandler 处理文件完整性检查请求的 handler
func fileIntegrityCheckHandler(c *fiber.Ctx) error {
    // 解析请求参数
    path := c.Query("path")
    expectedHash := c.Query("expectedHash")

    // 创建文件完整性检查器实例
    fic := NewFileIntegrityChecker(path, expectedHash)

    // 执行检查
    if err := fic.CheckIntegrity(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // 返回成功响应
    return c.JSON(fiber.Map{
        "message": "File integrity check passed",
    })
}

func main() {
    app := fiber.New()

    // 设置路由和 handler
    app.Get("/check", fileIntegrityCheckHandler)

    // 启动服务
    if err := app.Listen(":3000"); err != nil {
        fmt.Println("Error starting server: ", err)
        os.Exit(1)
    }
}
