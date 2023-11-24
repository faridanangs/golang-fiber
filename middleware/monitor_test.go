package middleware

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func TestMonitor(t *testing.T) {
	app := fiber.New()

	// Initialize default config (Assign the middleware to /metrics)
	// app.Get("/metrics", monitor.New())

	// Or extend your config for customization
	// Assign the middleware to /metrics
	// and change the Title to `MyService Metrics Page`
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:   "MyService Metrics Page",
		Refresh: 1 * time.Nanosecond,
		// jika kita gunakan apOnly maka dia tidak mngebalikan kode krapik
		// tap dia mengembalikan format json
		// APIOnly: true,
	}))

	app.Listen(":3000")
}
