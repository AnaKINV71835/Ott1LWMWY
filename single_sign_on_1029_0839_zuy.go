// 代码生成时间: 2025-10-29 08:39:10
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

// 用户模型
type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"createdAt"`
}

// 用户存储，用于演示
var users = []User{
    {ID: "1", Username: "user1", CreatedAt: time.Now()},
    {ID: "2", Username: "user2", CreatedAt: time.Now()},
}

// Session管理结构
type SessionManager struct {
    Sessions map[string]string // 存储用户ID和Session ID的映射
}

// 单点登录系统
func main() {
    app := fiber.New()

    // 配置CORS
    app.Use(cors.New())
    // 配置Recover中间件
    app.Use(recover.New())
    // 配置Logger中间件
    app.Use(logger.New())

    // 初始化Session管理器
    sessionManager := &SessionManager{
        Sessions: make(map[string]string),
    }

    // 登录路由
    app.Post("/login", func(c *fiber.Ctx) error {
        // 从请求体中提取用户名和密码
        var loginData struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := c.BodyParser(&loginData); err != nil {
            return c.Status(http.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid request body",
            })
        }

        // 验证用户名和密码
        for _, user := range users {
            if user.Username == loginData.Username && user.Username == loginData.Password {
                // 创建Session ID
                sessionID := fmt.Sprintf("session_%d", time.Now().Unix())
                // 存储Session ID和用户ID的映射
                sessionManager.Sessions[sessionID] = user.ID
                // 返回Session ID
                return c.JSON(fiber.Map{
                    "session_id": sessionID,
                })
            }
        }

        // 用户名或密码错误
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid username or password",
        })
    })

    // 检查Session路由
    app.Get("/check_session", func(c *fiber.Ctx) error {
        sessionID := c.Get("session_id")
        if _, exists := sessionManager.Sessions[sessionID]; !exists {
            return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid session",
            })
        }
        // 验证Session有效，返回用户信息
        userID := sessionManager.Sessions[sessionID]
        for _, user := range users {
            if user.ID == userID {
                return c.JSON(user)
            }
        }
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "User not found",
        })
    })

    // 启动服务器
    app.Listen(":3000")
}