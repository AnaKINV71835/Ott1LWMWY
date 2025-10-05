// 代码生成时间: 2025-10-06 02:51:31
package main

import (
    "image"
    "image/color"
    "image/png"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// WatermarkConfig contains configuration for watermarking
type WatermarkConfig struct {
    Text         string
    FontPath    string
    FontSize    int
    Color       color.RGBA
    opacity     float64
}

// WatermarkService provides functionality to add watermarks to images
type WatermarkService struct {
    config WatermarkConfig
}

// NewWatermarkService initializes a new WatermarkService with default settings
func NewWatermarkService(config WatermarkConfig) *WatermarkService {
    return &WatermarkService{
        config: config,
    }
}

// AddWatermark adds a watermark to an image
func (s *WatermarkService) AddWatermark(img image.Image, width, height int) image.Image {
    // Create a new RGBA image
    rgba := image.NewRGBA(image.Rect(0, 0, width, height))
    // Draw the original image onto the new image
    rgba.SetColorSpace(image.UniformSpace)
    draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

    // Define the font and prepare to draw the watermark text
    font, err := freetype.ParseFont(s.config.FontPath)
    if err != nil {
        log.Fatalf("Error parsing font: %s", err)
    }
    c := freetype.NewContext()
    c.SetDPI(72)
    c.SetFont(font)
    c.SetFontSize(float64(s.config.FontSize))
    pt := freetype.Pt(width/2, height/2)
    _, h := c.FontSize()
    _, baseline := c.Kern(c.Font, s.config.Text[0])
    baseline = baseline.Ceil()
    d := &font.DrawerContext{
        Dst: rgba,
        Src: image.NewUniform(s.config.Color),
        Face: c.Face,
        Dx:  pt.X - c.StringWidth(s.config.Text)/2,
        Dy:  pt.Y + (h-baseline)/2 + baseline.Ceil()/2,
    }

    // Draw the watermark text
    d.DrawString(s.config.Text)

    // Apply the opacity
    img = rgba.SubImage(rgba.Bounds()).(image.Image)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r, g, b, a := rgba.At(x, y).RGBA()
            ra := uint8(float64(r) * (1 - s.config.opacity))
            ga := uint8(float64(g) * (1 - s.config.opacity))
            ba := uint8(float64(b) * (1 - s.config.opacity))
            rgba.Set(x, y, color.RGBA{ra, ga, ba, a})
        }
    }

    return rgba
}

// WatermarkHandler handles the HTTP request to add a watermark
func WatermarkHandler(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "No file uploaded",
        })
    }
    defer file.Close()

    imgBytes, err := ioutil.ReadAll(file)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Error reading file",
        })
    }

    img, _, err := image.Decode(bytes.NewReader(imgBytes))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid image format",
        })
    }

    width := img.Bounds().Dx()
    height := img.Bounds().Dy()

    // Initialize the watermark service with default config
    watermarkService := NewWatermarkService(WatermarkConfig{
        Text:         "Watermark",
        FontPath:    "path/to/font.ttf",
        FontSize:    24,
        Color:       color.RGBA{255, 255, 255, 255},
        opacity:     0.5,
    })

    // Add watermark to the image
    watermarkedImg := watermarkService.AddWatermark(img, width, height)

    // Save the watermarked image to a buffer
    var buf bytes.Buffer
    png.Encode(&buf, watermarkedImg)

    // Send the watermarked image as a response
    return c.Send(buf.Bytes(), "image/png")
}

func main() {
    app := fiber.New()

    app.Post("/add-watermark", WatermarkHandler)

    log.Fatal(app.Listen(":3000"))
}
