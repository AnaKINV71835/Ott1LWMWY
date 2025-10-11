// 代码生成时间: 2025-10-12 02:39:17
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/sourcegraph/syntaxhighlight"
)

// CodeHighlighterHandler handles requests for code highlighting
func CodeHighlighterHandler(c *fiber.Ctx) error {
    code := c.Query("code")
    language := c.Query("language")

    if code == "" || language == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Code and language are required",
        })
    }

    highlighted, err := syntaxhighlight.HTMLTag(code, language, true)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to highlight code",
        })
    }

    return c.SendString(highlighted)
}

func main() {
    app := fiber.New()

    // Register the code highlighter handler
    app.Get("/highlight", CodeHighlighterHandler)

    // Start the server
    if err := app.Listen(":3000"); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
