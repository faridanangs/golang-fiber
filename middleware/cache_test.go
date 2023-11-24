package middleware

import (
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"
)

func TestCache(t *testing.T) {
	app := fiber.New()
	// Initialize default config
	app.Use(cache.New())

	// Or extend your config for customization
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))
	// Or you can custom key and expire time like this:

	app.Use(cache.New(cache.Config{
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			newCacheTime, _ := strconv.Atoi(c.GetRespHeader("Cache-Time", "600"))
			return time.Second * time.Duration(newCacheTime)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", "6000")
		return c.SendString("hi")
	})

	app.Listen(":3000")
}
