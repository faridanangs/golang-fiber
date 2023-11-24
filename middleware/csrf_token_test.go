package middleware

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

// CSRF middleware for Fiber that provides Cross-site request forgery protection by passing a csrf token via cookies. This cookie value will be used to compare against the client csrf token on requests, other than those defined as "safe" by RFC7231 (GET, HEAD, OPTIONS, or TRACE). When the csrf token is invalid, this middleware will return the fiber.ErrForbidden error.
// CSRF Tokens are generated on GET requests. You can retrieve the CSRF token with c.Locals(contextKey), where contextKey is the string you set in the config (see Custom Config below).
// When no csrf_ cookie is set, or the token has expired, a new token will be generated and csrf_ cookie set.
// NOTE
// This middleware uses our Storage package to support various databases through a single interface. The default configuration for this middleware saves data to memory, see the examples below for other databases.

func TestCsrfToken(t *testing.T) {
	app := fiber.New()

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
		CookieHTTPOnly: true,
	}))

	app.Get("/user", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "user",
			Value:    "you passed in auth ",
			MaxAge:   -1,
			SameSite: "Lax",
			HTTPOnly: true,
		})
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
