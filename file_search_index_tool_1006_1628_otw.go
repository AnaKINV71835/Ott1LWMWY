// 代码生成时间: 2025-10-06 16:28:01
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "github.com/gofiber/fiber/v2"
)

// FileSearchIndexTool 结构体，用于保存文件搜索和索引信息
type FileSearchIndexTool struct {
    RootPath     string
    Index       map[string][]string // key: 文件名，value: 文件路径列表
    IgnoreCase bool
}

// NewFileSearchIndexTool 构造函数，用于创建 FileSearchIndexTool 实例
func NewFileSearchIndexTool(rootPath string, ignoreCase bool) *FileSearchIndexTool {
    return &FileSearchIndexTool{
        RootPath:     rootPath,
        Index:        make(map[string][]string),
        IgnoreCase:   ignoreCase,
    }
}

// BuildIndex 构建文件索引
func (f *FileSearchIndexTool) BuildIndex() error {
    err := filepath.WalkDir(f.RootPath, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        fileName := d.Name()
        if f.IgnoreCase {
            fileName = strings.ToLower(fileName)
        }
        f.Index[fileName] = append(f.Index[fileName], path)
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

// SearchByFileName 根据文件名搜索文件
func (f *FileSearchIndexTool) SearchByFileName(fileName string) ([]string, error) {
    if f.IgnoreCase {
        fileName = strings.ToLower(fileName)
    }
    if paths, exists := f.Index[fileName]; exists {
        return paths, nil
    }
    return nil, fmt.Errorf("file not found: %s", fileName)
}

// StartServer 启动 Fiber 服务器，提供文件搜索和索引的 RESTful API
func (f *FileSearchIndexTool) StartServer(port int) error {
    app := fiber.New()
    app.Get("/search", func(c *fiber.Ctx) error {
        query := c.Query("q")
        paths, err := f.SearchByFileName(query)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(fiber.Map{
            "query": query,
            "paths": paths,
        })
    })
    return app.Listen(fmt.Sprintf(":%d", port))
}

func main() {
    rootPath := "./" // 搜索的根目录
    port := 3000       // 服务器端口
    tool := NewFileSearchIndexTool(rootPath, true) // 创建工具实例
    err := tool.BuildIndex()
    if err != nil {
        log.Fatalf("Failed to build index: %v", err)
    }
    err = tool.StartServer(port)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
