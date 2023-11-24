package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestApp(t *testing.T) {
	// Static
	app := fiber.New()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Ini method get")
	// })
	// app.Post("/api", func(c *fiber.Ctx) error {
	// 	return c.SendString("Ini method post")
	// })

	// // // Route Handlers
	// // Match any request jika cocok lanjutkan
	// app.Use("/", func(c *fiber.Ctx) error {
	// 	return c.Next()
	// })

	// // Match with first /api request jika cocok lanjutkan
	// app.Use("/api", func(c *fiber.Ctx) error {
	// 	return c.Next()
	// })

	// Match requests starting with /api or /home (multiple-prefix support)
	// app.Use([]string{"/api", "/home"}, func(c *fiber.Ctx) error {
	// 	return c.Next()
	// })

	// Attach multiple handlers
	app.Use("/api", func(c *fiber.Ctx) error {
		c.Set("X-Custom-Header", "hello loli")
		return c.Next()
	}, func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Listen(":3000")
}

func TestMount(t *testing.T) {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/jhon", micro)

	micro.Get("/doe", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":3000"))
}

func TestGroup(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	v1.Get("/list", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	v1.Get("/user", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// result api/v2/list
	// result api/v2/user
	v2 := api.Group("/v2", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	v2.Get("/list", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	v2.Get("/user", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":3000"))
}

func TestRouter(t *testing.T) {
	app := fiber.New()

	// untuk bisa mendapatkan nilai return dari 200 kita harus masuke ke dalam url atau
	// path api/foo, jika kita hanya masuk ke dalam url /api saja maka dia tidak akan mereturn data 200 nya
	app.Route("/api", func(router fiber.Router) {
		router.Get("/foo", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})
	}, "test")
	log.Fatal(app.Listen(":3000"))
}

// dia awal berhenti jika awal awal program di jalankan dan tidak di gunakan
func TestShutdown(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3000")
	app.ShutdownWithTimeout(1 * time.Minute)
}

var handler = func(c *fiber.Ctx) error { return nil }

// dia akan menamilkan semua stack dari methodnya
func TestStack(t *testing.T) {
	app := fiber.New()

	app.Get("/jhon/:age", handler)
	app.Post("/register", handler)

	res, _ := json.MarshalIndent(app.Stack(), "", " ")
	fmt.Println(string(res))
	app.Listen(":3000")
}

func TestName(t *testing.T) {

	// untuk memberikan nama di dalamnya
	app := fiber.New()

	app.Get("/", handler).Name("index")

	app.Get("/doe", handler).Name("home")

	app.Trace("/tracer", handler).Name("tracert")

	app.Delete("/delete", handler).Name("delete")
	// // contoh
	// {
	// 	"method": "DELETE",
	// 	"name": "delete",
	// 	"path": "/delete",
	// 	"params": null
	// }

	a := app.Group("/a").Name("fd.")

	a.Get("/test", handler).Name("test")

	data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	fmt.Print(string(data))

	app.Listen(":3000")
}

func TestGetRoute(t *testing.T) {
	app := fiber.New()

	app.Get("/", handler).Name("index")

	// jika kita ingin menggunakan app.GetRoute maka .Name harus sama dengan yang kita panggil
	// di dalam getroute supaya tidak terjadi error
	data, _ := json.MarshalIndent(app.GetRoute("indx"), "", " ")

	fmt.Print(string(data))
	app.Listen(":3000")
}

func TestGetRoutes(t *testing.T) {
	// Metode ini mendapatkan semua rute.
	// func (app *App) GetRoutes(filterUseOption ...bool) []Route
	// Ketika filterUseOption sama dengan true, itu akan memfilter rute yang didaftarkan oleh middleware.

	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	}).Name("index")

	data, _ := json.MarshalIndent(app.GetRoutes(true), "", " ")
	fmt.Print(string(data))

	app.Listen(":3000")
}

func TestConfig(t *testing.T) {
	// Config mengembalikan konfigurasi aplikasi sebagai nilai ( read-only ) .
	// func (app *App) Config() Config

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}).Name("index")
	data, _ := json.MarshalIndent(app.Config(), "", " ")
	fmt.Print(string(data))

	app.Listen(":3000")
	// INI contoh hasil outputnya
	{
		/*

			// === RUN   TestConfig
			// {
			//  "prefork": false,
			//  "server_header": "",
			//  "strict_routing": false,
			//  "case_sensitive": false,
			//  "immutable": false,
			//  "unescape_path": false,
			//  "etag": false,
			//  "body_limit": 4194304,
			//  "concurrency": 262144,
			//  "views_layout": "",
			//  "pass_locals_to_views": false,
			//  "read_timeout": 0,
			//  "write_timeout": 0,
			//  "idle_timeout": 0,
			//  "read_buffer_size": 4096,
			//  "write_buffer_size": 4096,
			//  "compressed_file_suffix": ".fiber.gz",
			//  "proxy_header": "",
			//  "get_only": false,
			//  "disable_keepalive": false,
			//  "disable_default_date": false,
			//  "disable_default_content_type": false,
			//  "disable_header_normalizing": false,
			//  "disable_startup_message": false,
			//  "app_name": "",
			//  "StreamRequestBody": false,
			//  "DisablePreParseMultipartForm": false,
			//  "reduce_memory_usage": false,
			//  "Network": "tcp4",
			//  "enable_trusted_proxy_check": false,
			//  "trusted_proxies": null,
			//  "enable_ip_validation": false,
			//  "enable_print_routes": false,
			//  "color_scheme": {
			//   "Black": "\u001b[90m",
			//   "Red": "\u001b[91m",
			//   "Green": "\u001b[92m",
			//   "Yellow": "\u001b[93m",
			//   "Blue": "\u001b[94m",
			//   "Magenta": "\u001b[95m",
			//   "Cyan": "\u001b[96m",
			//   "White": "\u001b[97m",
			//   "Reset": "\u001b[0m"
			//  },
			//  "RequestMethods": [
			//   "GET",
			//   "HEAD",
			//   "POST",
			//   "PUT",
			//   "DELETE",
			//   "CONNECT",
			//   "OPTIONS",
			//   "TRACE",
			//   "PATCH"
			//  ],
			//  "enable_splitting_on_parsers": false
			// }
		*/
	}
}

func TestHandler(t *testing.T) {
	// Handler mengembalikan handler server yang dapat digunakan untuk melayani permintaan kustom * fasthttp.RequestCtx.
	// func (app *App) Handler() fasthttp.RequestHandler

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		data, _ := json.MarshalIndent(app.Handler(), "", " ")
		fmt.Print(string(data))
		return c.SendStatus(200)
	}).Name("index")

	app.Listen(":3000")
}

func TestRequest(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.BaseURL())
		fmt.Println(c.Get("X-Sapa"))
		return c.SendStatus(200)
	}).Name("home")

	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000", nil)
	req.Header.Set("X-Sapa", "hello")

	// kemudian kita kirim responsenya ke dalam fiber test
	resp, _ := app.Test(req)

	if resp.StatusCode == fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

func TestLatihan(t *testing.T) {
	app := fiber.New()
	micro := fiber.New()
	app.Mount("/api", micro)
	micro.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})
	micro.Get("/signin", func(c *fiber.Ctx) error {
		return c.SendString("Signin")
	})

	app.Listen(":3000")

}
