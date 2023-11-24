package middleware

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func TestBasictAuth(t *testing.T) {
	app := fiber.New()

	// Or extend your config for customization
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"john":  "doe",
			"admin": "123456",
		},
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			if user == "john" && pass == "doe" {
				return true
			}
			if user == "admin" && pass == "123456" {
				return true
			}
			return false
		},
		// jika kode di dalam Authorizer tidak ada yang memenuhi dia akan menjalankan Unauthorized
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendFile("./unauthorized.htm", true)
		},
		// kemudian kita masuakn pass dan user ke dalam contextPass dan User dengan variabel _pass dan _user
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))
	app.Use(cors.New(cors.Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Content-Type, Accept, Origin",
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// ketika kita ingin memanggil password di dalam basic auth kita gunakan
		// c.Locals("_pass")  _pass == ContextPassword
		if c.Locals("_user") == "admin" && c.Locals("_pass") == "123456" {
			return c.SendString("hello admin")
		}
		if c.Locals("_user") == "john" && c.Locals("_pass") == "doe" {
			return c.SendString("hello john")
		}
		fmt.Println(c.Locals("_pass"))
		return c.SendStatus(404)
	})

	app.Listen(":3000")
}
