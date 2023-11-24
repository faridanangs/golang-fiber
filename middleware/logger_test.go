package middleware

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLogger(t *testing.T) {
	app := fiber.New()

	// Initialize default config
	// app.Use(logger.New())

	// Or extend your config for customization
	// Logging remote IP and Port
	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	// }))

	// Logging Request ID
	// app.Use(requestid.New())
	// app.Use(logger.New(logger.Config{
	// 	// For more options, see the Config section
	// 	Format: "${pid} - ${locals:requestid} ${status} - ${method} ${path}​\n",
	// }))

	// Changing TimeZone & TimeFormat
	// app.Use(logger.New(logger.Config{
	// 	Format:     "${pid} ${status} - ${method} ${path}\n",
	// 	TimeFormat: "02-Jan-2006",
	// 	TimeZone:   "America/New_York",
	// }))

	// Custom File Writer
	// file, err := os.OpenFile("./123.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer file.Close()
	// app.Use(logger.New(logger.Config{
	// 	Output: file,
	// }))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}
