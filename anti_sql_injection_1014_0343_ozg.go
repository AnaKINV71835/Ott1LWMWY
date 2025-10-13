// 代码生成时间: 2025-10-14 03:43:21
package main

import (
    "fmt"
    "github.com/go-fiber/fiber/v2"
    "github.com/go-fiber/fiber/v2/middleware/cors"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// Database struct to hold our SQLite connection
type Database struct {
    DB *gorm.DB
}

// initDB initializes the SQLite database connection
func initDB() *Database {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    // Migrations
    db.AutoMigrate(&User{})
    return &Database{DB: db}
}

// User model for demonstration
type User struct {
    gorm.Model
    Name string
    Email string `gorm:"uniqueIndex"`
}

func main() {
    // Initialize Fiber with a custom configuration
    app := fiber.New(fiber.Config{
        Prefork:       false,
        ServerHeader: "gofiber",
        CaseySensitive: true,
    })

    // Apply CORS middleware
    app.Use(cors.New())

    // Initialize the database
    db := initDB()

    // Define a route to create a new user, demonstrating SQL injection prevention
    app.Post("/users", func(c *fiber.Ctx) error {
        name := c.Query("name")
        email := c.Query("email")
        if name == "" || email == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Name and Email are required fields",
            })
        }

        // Use GORM to prevent SQL injection by parameterized queries
        result := db.DB.Create(&User{Name: name, Email: email})
        if result.Error != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": result.Error.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "message": "User created successfully",
            "user": result.Values,
        })
    })

    // Start the Fiber server
    app.Listen(":3000")
}
