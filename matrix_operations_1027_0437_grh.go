// 代码生成时间: 2025-10-27 04:37:05
package main

import (
    "fmt"
    "math"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// Matrix represents a 2D matrix of float64 values.
type Matrix [][][]float64

// MatrixOperations is a struct containing Matrix operations.
type MatrixOperations struct{}

// NewMatrixOperations creates a new MatrixOperations instance.
func NewMatrixOperations() *MatrixOperations {
    return &MatrixOperations{}
}

// Add adds two matrices together and returns the result as a new matrix.
func (m *MatrixOperations) Add(matrix1, matrix2 Matrix) (Matrix, error) {
    if len(matrix1) != len(matrix2) || len(matrix1[0]) != len(matrix2[0]) {
        return nil, fmt.Errorf("matrices have incompatible dimensions")
    }

    result := make(Matrix, len(matrix1))
    for i := range matrix1 {
        result[i] = make([]float64, len(matrix1[i]))
        for j := range matrix1[i] {
            result[i][j] = matrix1[i][j] + matrix2[i][j]
        }
    }
    return result, nil
}

// Subtract subtracts one matrix from another and returns the result as a new matrix.
func (m *MatrixOperations) Subtract(matrix1, matrix2 Matrix) (Matrix, error) {
    if len(matrix1) != len(matrix2) || len(matrix1[0]) != len(matrix2[0]) {
        return nil, fmt.Errorf("matrices have incompatible dimensions")
    }

    result := make(Matrix, len(matrix1))
    for i := range matrix1 {
        result[i] = make([]float64, len(matrix1[i]))
        for j := range matrix1[i] {
            result[i][j] = matrix1[i][j] - matrix2[i][j]
        }
    }
    return result, nil
}

// Multiply multiplies two matrices together and returns the result as a new matrix.
func (m *MatrixOperations) Multiply(matrix1, matrix2 Matrix) (Matrix, error) {
    if len(matrix1[0]) != len(matrix2) {
        return nil, fmt.Errorf("matrices have incompatible dimensions for multiplication")
    }

    result := make(Matrix, len(matrix1))
    for i := range matrix1 {
        result[i] = make([]float64, len(matrix2[0]))
        for j := range matrix2[0] {
            for k := range matrix1[0] {
                result[i][j] += matrix1[i][k] * matrix2[k][j]
            }
        }
    }
    return result, nil
}

// StartServer starts the Fiber HTTP server with matrix operations endpoints.
func StartServer() error {
    app := fiber.New()

    // Define the /add route for matrix addition.
    app.Post("/add", func(c *fiber.Ctx) error {
        var matrices struct {
            Matrix1 Matrix `json:"matrix1"`
            Matrix2 Matrix `json:"matrix2"`
        }
        if err := c.BodyParser(&matrices); err != nil {
            return err
        }

        result, err := NewMatrixOperations().Add(matrices.Matrix1, matrices.Matrix2)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.JSON(result)
    })

    // Define the /subtract route for matrix subtraction.
    app.Post("/subtract", func(c *fiber.Ctx) error {
        var matrices struct {
            Matrix1 Matrix `json:"matrix1"`
            Matrix2 Matrix `json:"matrix2"`
        }
        if err := c.BodyParser(&matrices); err != nil {
            return err
        }

        result, err := NewMatrixOperations().Subtract(matrices.Matrix1, matrices.Matrix2)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.JSON(result)
    })

    // Define the /multiply route for matrix multiplication.
    app.Post("/multiply", func(c *fiber.Ctx) error {
        var matrices struct {
            Matrix1 Matrix `json:"matrix1"`
            Matrix2 Matrix `json:"matrix2"`
        }
        if err := c.BodyParser(&matrices); err != nil {
            return err
        }

        result, err := NewMatrixOperations().Multiply(matrices.Matrix1, matrices.Matrix2)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.JSON(result)
    })

    return app.Listen(":3000")
}

func main() {
    if err := StartServer(); err != nil {
        fmt.Println("Server failed to start: ", err)
    }
}
