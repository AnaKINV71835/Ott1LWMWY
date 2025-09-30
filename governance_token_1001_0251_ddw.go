// 代码生成时间: 2025-10-01 02:51:29
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
)

// Token represents a single governance token
type Token struct {
    ID    string `json:"id"`
    Owner string `json:"owner"`
    Amount int    `json:"amount"`
}

// TokenService handles logic related to governance tokens
type TokenService struct {
    // This could be a database connection or any other persistent storage
    storage map[string]Token
}

// NewTokenService creates a new TokenService instance
func NewTokenService() *TokenService {
    return &TokenService{
        storage: make(map[string]Token),
    }
}

// CreateToken adds a new token to the storage
func (s *TokenService) CreateToken(token Token) (string, error) {
    if _, exists := s.storage[token.ID]; exists {
        return "", fmt.Errorf("token with id %s already exists", token.ID)
    }
    s.storage[token.ID] = token
    return token.ID, nil
}

// GetToken retrieves a token by its ID
func (s *TokenService) GetToken(tokenID string) (Token, error) {
    token, exists := s.storage[tokenID]
    if !exists {
        return Token{}, fmt.Errorf("token with id %s not found", tokenID)
    }
    return token, nil
}

// UpdateToken updates an existing token in the storage
func (s *TokenService) UpdateToken(tokenID string, token Token) error {
    if _, exists := s.storage[tokenID]; !exists {
        return fmt.Errorf("token with id %s not found", tokenID)
    }
    s.storage[tokenID] = token
    return nil
}

// DeleteToken removes a token from the storage
func (s *TokenService) DeleteToken(tokenID string) error {
    if _, exists := s.storage[tokenID]; !exists {
        return fmt.Errorf("token with id %s not found", tokenID)
    }
    delete(s.storage, tokenID)
    return nil
}

// SetupRoutes configures the Fiber app with routes for governance tokens
func SetupRoutes(app *fiber.App, service *TokenService) {
    app.Post("/token", func(c *fiber.Ctx) error {
        var token Token
        if err := c.BodyParser(&token); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Unable to parse token data",
            })
        }
        tokenID, err := service.CreateToken(token)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.Status(fiber.StatusCreated).JSON(fiber.Map{
            "id": tokenID,
        })
    })

    app.Get("/token/:tokenID", func(c *fiber.Ctx) error {
        tokenID := c.Params("tokenID\)
        token, err := service.GetToken(tokenID)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(token)
    })

    app.Put("/token/:tokenID", func(c *fiber.Ctx) error {
        tokenID := c.Params("tokenID\)
        var token Token
        if err := c.BodyParser(&token); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Unable to parse token data",
            })
        }
        err := service.UpdateToken(tokenID, token)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.SendStatus(fiber.StatusOK)
    })

    app.Delete("/token/:tokenID", func(c *fiber.Ctx) error {
        tokenID := c.Params("tokenID\)
        err := service.DeleteToken(tokenID)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.SendStatus(fiber.StatusOK)
    })
}

func main() {
    app := fiber.New()
    service := NewTokenService()
    SetupRoutes(app, service)
    app.Listen(":3000")
}