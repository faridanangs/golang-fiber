package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/envvar"
)

func TestEnvVar(t *testing.T) {
	app := fiber.New()

	// Initialize default config
	// app.Use("/expose/envvars", envvar.New())

	// Or extend your config for customization
	app.Use("/expose/envvar", envvar.New(envvar.Config{
		ExportVars:  map[string]string{"testKey": "", "testDefaultKey": "testDefaultVal"},
		ExcludeVars: map[string]string{"excludeKey": ""},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}
