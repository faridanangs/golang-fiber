package middleware

import (
	"expvar"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	expvarmw "github.com/gofiber/fiber/v2/middleware/expvar"
)

func TestExpVar(t *testing.T) {
	app := fiber.New()

	var count = expvar.NewInt("count")

	app.Use(expvarmw.New())
	app.Get("/", func(c *fiber.Ctx) error {
		count.Add(1)
		return c.SendString(fmt.Sprintf("hello expvar count %d", count.Value()))
	})
	// Visit path /debug/vars to see all vars and use query r=key to filter exposed variables.
	// curl 127.0.0.1:3000
	// hello expvar count 1
	
	// curl 127.0.0.1:3000/debug/vars
	// {
	// 	"cmdline": ["xxx"],
	// 	"count": 1,
	// 	"expvarHandlerCalls": 33,
	// 	"expvarRegexpErrors": 0,
	// 	"memstats": {...}
	// }

	// curl 127.0.0.1:3000/debug/vars?r=c
	// {
	// 	"cmdline": ["xxx"],
	// 	"count": 1
	// }
	app.Listen(":3000")
}
