package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func TestRecover(t *testing.T) {
	app := fiber.New()

	// supaya kodenya tidak berhenti jika terjadi panic kita tangkap menggunakan recover dari viber
	app.Use(recover.New())
	app.Get("/", func(c *fiber.Ctx) error {
		panic("i am error")
	})
	app.Listen(":3000")
}
