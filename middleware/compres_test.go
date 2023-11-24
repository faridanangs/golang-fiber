package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func TestCompress(t *testing.T) {
	app := fiber.New()

	// /Compression middleware for Fiber that will compress the response using gzip, deflate and brotli compression depending on the Accept-Encoding header.
	// After you initiate your Fiber app, you can use the following possibilities:

	// Initialize default config
	app.Use(compress.New())

	// Or extend your config for customization
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// Skip middleware for specific routes
	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/dont_compress"
		},
		Level: compress.LevelBestSpeed, // 1
	}))

	// 	Default Config
	// var ConfigDefault = Config{
	//     Next:  nil,
	//     Level: LevelDefault,
	// }

	// Constants
	// // Compression levels
	// const (
	//     LevelDisabled        = -1
	//     LevelDefault         = 0
	//     LevelBestSpeed       = 1
	//     LevelBestCompression = 2
	// )

	app.Get("/dont_compress", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3000")
}
