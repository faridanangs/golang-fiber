package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func TestRedirect(t *testing.T) {
	app := fiber.New()

	// jika kita masuk ke dalam path old walaupun di get pathny NEW MAKA KIATA akan
	// di redirect ke url new
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/old":   "/new",
			"/old/*": "/new/$1",
		},
		StatusCode: 301,
	}))

	app.Get("/new", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})
	app.Get("/new/*", func(c *fiber.Ctx) error {
		return c.SendString("hello " + c.Params("*"))
	})

	app.Listen(":3000")

}
