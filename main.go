package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Helo world")
	// })

	// menangkap parameter
	// GET http://localhost:8080/helloworld
	// app.Get("/:sapa", func(c *fiber.Ctx) error {
	// 	return c.SendString("Helo" + c.Params("sapa"))
	// })

	// // GET http://localhost:3000/john
	// opsional parameter
	// app.Get("/:sapa?", func(c *fiber.Ctx) error {
	// 	if c.Params("sapa") != "" {
	// 		return c.SendString("Helo " + c.Params("sapa"))
	// 	} else {
	// 		return c.SendString("where is json?")
	// 	}
	// })

	// Wildcards
	// // GET http://localhost:3000/sapa/user/john
	// app.Get("/:sapa/*", func(c *fiber.Ctx) error {
	// 	return c.SendString("Helo " + c.Params("*"))
	// })

	// Static files
	// app.Static("/", "./public")

	// new Error
	app.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(404, "error tolol")
	})

	app.Listen(":8080")
}
