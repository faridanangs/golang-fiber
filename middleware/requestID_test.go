package middleware

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
)

func TestRequestID(t *testing.T) {
	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		Header: "X-Header-Key",
		Generator: utils.UUIDv4,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
