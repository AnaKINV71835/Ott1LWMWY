// 代码生成时间: 2025-10-07 01:32:48
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/hex"
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "github.com/gofiber/fiber/v2"
)

// FileEncryptionDecryption provides functionality for encrypting and decrypting files.
type FileEncryptionDecryption struct {
    key []byte
}

// NewFileEncryptionDecryption initializes the encryption/decryption tool with a key.
func NewFileEncryptionDecryption(key []byte) *FileEncryptionDecryption {
    return &FileEncryptionDecryption{key: key}
}

// EncryptFile encrypts a file using AES-256-GCM.
func (f *FileEncryptionDecryption) EncryptFile(filename string) error {
    fileContent, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    block, err := aes.NewCipher(f.key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(ioutil.Discard, nonce); err != nil {
        return err
    }

    encrypted := gcm.Seal(nonce, nonce, fileContent, nil)
    return ioutil.WriteFile(filename+".enc", encrypted, 0644)
}

// DecryptFile decrypts a file using AES-256-GCM.
func (f *FileEncryptionDecryption) DecryptFile(encryptedFilename string) error {
    encryptedFileContent, err := ioutil.ReadFile(encryptedFilename)
    if err != nil {
        return err
    }

    nonceSize := aes.BlockSize
    if len(encryptedFileContent) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := encryptedFileContent[:nonceSize], encryptedFileContent[nonceSize:]

    block, err := aes.NewCipher(f.key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    if _, err := gcm.Open(nil, nonce, ciphertext, nil); err != nil {
        return err
    }

    return ioutil.WriteFile(encryptedFilename[:len(encryptedFilename)-4], ciphertext, 0644)
}

// StartServer starts a Fiber server that handles file encryption and decryption.
func StartServer(key []byte) {
    app := fiber.New()

    // Route for encrypting a file.
    app.Post("/encrypt", func(c *fiber.Ctx) error {
        filename := c.Query("filename")
        if filename == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "filename parameter is required",
            })
        }

        tool := NewFileEncryptionDecryption(key)
        if err := tool.EncryptFile(filename); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "message": "File encrypted successfully",
        })
    })

    // Route for decrypting a file.
    app.Post("/decrypt", func(c *fiber.Ctx) error {
        encryptedFilename := c.Query("filename")
        if encryptedFilename == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "filename parameter is required",
            })
        }

        tool := NewFileEncryptionDecryption(key)
        if err := tool.DecryptFile(encryptedFilename); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "message": "File decrypted successfully",
        })
    })

    // Start the server.
    if err := app.Listen(":3000"); err != nil {
        fmt.Println("Error starting server: ", err)
    }}