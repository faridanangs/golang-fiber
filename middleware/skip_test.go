package middleware

import (
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

func TestSkip(t *testing.T) {
	app := fiber.New()

	app.Use(skip.New(func(c *fiber.Ctx) error {
		return c.SendString("it was not a post req")
	},
		// jika di bawah ini mengembalikan false maka func di atas baru di jalankan naum jika
		// true func di atas tidak di jalankan
		func(c *fiber.Ctx) bool {
			return c.Method() == http.MethodPost
		}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	log.Fatal(app.Listen(":3000"))
}
