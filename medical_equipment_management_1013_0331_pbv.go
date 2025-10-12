// 代码生成时间: 2025-10-13 03:31:31
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// MedicalEquipment represents a medical equipment entity
type MedicalEquipment struct {
    gorm.Model
    Name       string `json:"name"`
    Type       string `json:"type"`
    SerialNumber string `json:"serial_number"`
    Manufacturer string `json:"manufacturer"`
    Status     string `json:"status"`
}

// Database is a global variable to hold the database connection
var Database *gorm.DB

func main() {
    // Initialize the SQLite database
    db, err := gorm.Open(sqlite.Open("medical_equipment.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    Database = db

    // Migrate the schema
    Database.AutoMigrate(&MedicalEquipment{})

    // Initialize Fiber
    app := fiber.New()

    // CRUD operations for Medical Equipment
    app.Get("/equipments", getAllEquipments)
    app.Post("/equipments", addEquipment)
    app.Get("/equipments/:id", getEquipment)
    app.Put("/equipments/:id", updateEquipment)
    app.Delete("/equipments/:id", deleteEquipment)

    // Start the server
    app.Listen(":3000")
}

// getAllEquipments retrieves all medical equipment records
func getAllEquipments(c *fiber.Ctx) error {
    var equipments []MedicalEquipment
    if err := Database.Find(&equipments).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(equipments)
}

// addEquipment adds a new medical equipment record
func addEquipment(c *fiber.Ctx) error {
    var equipment MedicalEquipment
    if err := c.BodyParser(&equipment); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    if err := Database.Create(&equipment).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "Equipment added successfully",
    })
}

// getEquipment retrieves a single medical equipment record by ID
func getEquipment(c *fiber.Ctx) error {
    id := c.Params("id")
    var equipment MedicalEquipment
    if err := Database.First(&equipment, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "status":  "error",
            "message": "Equipment not found",
        })
    }
    return c.Status(fiber.StatusOK).JSON(equipment)
}

// updateEquipment updates an existing medical equipment record
func updateEquipment(c *fiber.Ctx) error {
    id := c.Params("id")
    var equipment MedicalEquipment
    if err := c.BodyParser(&equipment); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    equipment.ID, _ = strconv.Atoi(id)
    if err := Database.Save(&equipment).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Equipment updated successfully",
    })
}

// deleteEquipment deletes a medical equipment record by ID
func deleteEquipment(c *fiber.Ctx) error {
    id := c.Params("id")
    var equipment MedicalEquipment
    equipment.ID, _ = strconv.Atoi(id)
    if err := Database.Delete(&equipment, id).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Equipment deleted successfully",
    })
}