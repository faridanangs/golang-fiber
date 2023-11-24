package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
)

func TestEarlyData(t *testing.T) {
	app := fiber.New()
	// Initialize default config
	app.Use(earlydata.New())

	// Or extend your config for customization
	app.Use(earlydata.New(earlydata.Config{
		Error: fiber.ErrTooEarly,
		// ...
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
