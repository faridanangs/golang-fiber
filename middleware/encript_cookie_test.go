package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

// Encrypt middleware for Fiber which encrypts cookie values. Note: this middleware does not encrypt cookie names.
func TestEncriptCookie(t *testing.T) {
	// Intitializes the middleware
	// func New(config ...Config) fiber.Handler

	// Returns a random 32 character long string
	// func GenerateKey() string

	app := fiber.New()

	// Provide a minimal config
	// `Key` must be a 32 character string. It's used to encrypt the values, so make sure it is random and keep it secret.
	// You can run `openssl rand -base64 32` or call `encryptcookie.GenerateKey()` to create a random key for you.
	// Make sure not to set `Key` to `encryptcookie.GenerateKey()` because that will create a new key every run.
	// app.Use(encryptcookie.New(encryptcookie.Config{
	// 	Key: encryptcookie.GenerateKey(),
	// }))

	// 	Usage of CSRF and Encryptcookie Middlewares with Custom Cookie Names
	// Normally, encryptcookie middleware skips csrf_ cookies. However, it won't work when you use custom cookie names for CSRF. You should update Except config to avoid this problem. For example:
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    encryptcookie.GenerateKey(),
		Except: []string{"csrf_1"}, // exclude CSRF cookie
	}))
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:test",
		CookieName:     "csrf_1",
		CookieHTTPOnly: true,
	}))

	// Get / reading out the encrypted cookie
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("value:" + c.Cookies("test"))
	})
	// get /post create the encrypted cookie
	app.Get("/post", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:  "test",
			Value: "something",
		})
		return c.Redirect("/", 200)
	})

	app.Listen(":3000")
}
