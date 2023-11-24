package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
)

// Favicon
// Favicon middleware for Fiber that ignores favicon requests or caches a provided
// icon in memory to improve performance by skipping disk access.
// User agents request favicon.ico frequently and indiscriminately,
// so you may wish to exclude these requests from your logs by using this middleware before your logger middleware.
// NOTE
// This middleware is exclusively for serving the default, implicit favicon, which is GET /favicon.ico or custom favicon URL.


func TestFavicon(t *testing.T) {
	app := fiber.New()

	// jika kita ingin mengakses faviconya kita lebih harus masukke dalam urlnya
	app.Use(favicon.New(favicon.Config{
		File: "../image/05.png",
		URL:  "/anangs",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}
