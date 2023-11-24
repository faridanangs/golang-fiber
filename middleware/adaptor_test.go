package middleware

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// net/http to Fiber
func TestAdaptor(t *testing.T) {
	// New fiber app
	app := fiber.New()

	// http.Handler -> fiber.Handler
	app.Get("/", adaptor.HTTPHandler(handler(greet)))

	// http.HandlerFunc -> fiber.Handler
	app.Get("/func", adaptor.HTTPHandlerFunc(greet))

	// Listen on port 3000
	app.Listen(":3000")
}
func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

// net/http middleware to Fiber
func TestMiddleware(t *testing.T) {
	// New fiber app
	app := fiber.New()

	// http middleware -> fiber.Handler
	app.Use(adaptor.HTTPMiddleware(logMiddleware))

	// Listen on port 3000
	app.Listen(":3000")
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log middleware")
		next.ServeHTTP(w, r)
	})
}

// Fiber Handler to net/http
func TestFiberToHttp(t *testing.T) {
	http.Handle("/", adaptor.FiberHandler(greetFiber))

	http.HandleFunc("/func", adaptor.FiberHandlerFunc(greetFiber))

	http.ListenAndServe(":3000", nil)
}

func greetFiber(c *fiber.Ctx) error {
	return c.SendString("hello world")
}

// Fiber App to net/http
func TestFiberAppToHttp(t *testing.T) {
	app := fiber.New()
	app.Get("/", greetApp)

	http.ListenAndServe(":3000", adaptor.FiberApp(app))
}

func greetApp(c *fiber.Ctx) error {
	return c.SendString("hello world")
}

// Fiber Context to (net/http).Request
func TestFiberContextToHttp(t *testing.T) {
	app := fiber.New()
	app.Get("/greet", greetWithHttpReq)

	http.ListenAndServe(":3000", adaptor.FiberApp(app))
}

func greetWithHttpReq(c *fiber.Ctx) error {
	req, err := adaptor.ConvertRequest(c, false)
	if err != nil {
		return err
	}
	fmt.Println("Req url: " + req.URL.String())
	fmt.Println(req.Body)
	fmt.Println(req.Form)
	fmt.Println(req.Header)
	fmt.Println(req.Host)
	fmt.Println(req.Method)
	fmt.Println(req.MultipartForm)
	fmt.Println(req.PostForm)
	fmt.Println(req.Proto)
	fmt.Println(req.ProtoMajor)
	fmt.Println(req.ProtoMinor)
	fmt.Println(req.RemoteAddr)
	fmt.Println(req.RequestURI)
	fmt.Println(req.Response)
	fmt.Println(req.TLS)
	fmt.Println(req.Trailer)
	fmt.Println(req.TransferEncoding)
	fmt.Println(req.URL)
	return c.SendStatus(200)
}
