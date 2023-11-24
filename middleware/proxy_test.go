package middleware

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func TestProxy(t *testing.T) {
	app := fiber.New()

	// Forward to url
	// app.Get("/gif", proxy.Forward("https://i.imgur.com/IWaBepg.gif"))

	// // If you want to forward with a specific domain. You have to use proxy.DomainForward.
	// app.Get("/payments", proxy.DomainForward("docs.gofiber.io", "http://localhost:8000"))

	// Forward to url with local custom client
	// app.Get("/gif", proxy.Forward("https://i.imgur.com/IWaBepg.gif", &fasthttp.Client{
	// 	NoDefaultUserAgentHeader: true,
	// 	DisablePathNormalizing: true,
	// }))

	// Make request within handler
	// app.Get("/:id", func(c *fiber.Ctx) error {
	// 	url := "https://i.imgur.com/" + c.Params("id") + ".gif"
	// 	if err := proxy.Do(c, url); err != nil {
	// 		return err
	// 	}
	// 	// remove server handler from server
	// 	c.Response().Header.Del(fiber.HeaderServer)
	// 	return c.SendStatus(200)
	// })

	// Make proxy requests while following redirects
	// app.Get("/proxy", func(c *fiber.Ctx) error {
	// 	// kita di arahkan ke halaman seperti ggogle tapi bukan goole hanya halamannya saja
	// 	proxy.DoRedirects(c, "https://yandex.com", 3)

	// 	// remove server handler from server
	// 	c.Response().Header.Del(fiber.HeaderServer)
	// 	return nil
	// })

	// Make proxy requests and wait up to 5 seconds before timing out
	// app.Get("/proxy", func(c *fiber.Ctx) error {
	// 	if err := proxy.DoTimeout(c, "https://yandex.com", time.Second*5); err != nil {
	// 		return err
	// 	}
	// 	// Remove Server header from response
	// 	c.Response().Header.Del(fiber.HeaderServer)
	// 	return nil
	// })

	// Or extend your balancer for customization
	app.Use(proxy.Balancer(proxy.Config{
		Servers: []string{
			"http://localhost:3001",
			"http://localhost:3002",
			"http://localhost:3003",
		},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
	}))

	// Make proxy requests, timeout a minute from now
	app.Get("/proxy", func(c *fiber.Ctx) error {
		if err := proxy.DoDeadline(c, "https://yandex.com", time.Now().Add(10*time.Second)); err != nil {
			return err
		}
		// Remove Server header from response
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	app.Listen(":3000")
}
