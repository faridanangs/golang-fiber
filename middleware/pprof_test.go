package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func TestPprof(t *testing.T) {
	app := fiber.New()

	// untuk menampilkan datany kita taruh debug/pprof/ di akhiran pathnya
	app.Use(pprof.New(pprof.Config{
		Prefix: "/prefix",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
