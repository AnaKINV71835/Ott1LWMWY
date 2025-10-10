// 代码生成时间: 2025-10-11 02:05:24
// data_encryption_tool.go
// 一个使用GOLANG和FIBER框架的数据加密传输工具

package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fiber"
    "log"
    "net/http"
)

// EncryptionConfig 配置结构体
type EncryptionConfig struct {
    Key string
}

// NewEncryptionConfig 创建一个新的配置实例
func NewEncryptionConfig(key string) *EncryptionConfig {
    return &EncryptionConfig{Key: key}
}

// Encrypt 加密数据
func (config *EncryptionConfig) Encrypt(plaintext []byte) (string, error) {
    block, err := aes.NewCipher([]byte(config.Key))
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := rand.Read(iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
    
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据
func (config *EncryptionConfig) Decrypt(ciphertext string) ([]byte, error) {
    decoded, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher([]byte(config.Key))
    if err != nil {
        return nil, err
    }
    
    if len(decoded) < aes.BlockSize {
        return nil, err
    }
    
    iv := decoded[:aes.BlockSize]
    decoded = decoded[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decoded, decoded)
    
    return decoded, nil
}

// App app实例
var App = fiber.New()

func main() {
    config := NewEncryptionConfig("your-encryption-key") // 替换为你的密钥
    
    App.Post("/encrypt", func(c *fiber.Ctx) error {
        data := c.Get("plaintext")
        encryptedData, err := config.Encrypt([]byte(data))
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
                "error": "Encryption failed",
            })
        }
        return c.Status(http.StatusOK).JSON(fiber.Map{
            "ciphertext": encryptedData,
        })
    })
    
    App.Post("/decrypt", func(c *fiber.Ctx) error {
        ciphertext := c.Get("ciphertext")
        decryptedData, err := config.Decrypt(ciphertext)
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
                "error": "Decryption failed",
            })
        }
        return c.Status(http.StatusOK).JSON(fiber.Map{
            "plaintext": string(decryptedData),
        })
    })
    
    log.Fatal(App.Listen(":3000"))
}