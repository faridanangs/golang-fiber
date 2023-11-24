package middleware

import (
	"embed"
	_ "embed"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// FileSystem
// Filesystem middleware for Fiber that enables you to serve files from a directory.
// CAUTION
// :params & :optionals? within the prefix path are not supported!
// To handle paths with spaces (or other url encoded values) make sure to set fiber.Config{ UnescapePath: true }
func TestFileSystem(T *testing.T) {
	app := fiber.New()

	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("../public"),
		Browse: true,
		MaxAge: 3600,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// If your environment (Go 1.16+) supports it, we recommend using Go Embed instead of the other solutions listed as this one is native to Go and the easiest to use.
// embed
// Embed is the native method to embed files in a Golang excecutable. Introduced in Go 1.16.

// Embed a single file
//
//go:embed unauthorized.htm
var unauth embed.FS

//go:embed public/*
var direc embed.FS

func TestEmbedFunc(t *testing.T) {
	app := fiber.New()

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(unauth),
	}))

	// Access file "image.png" under `static/` directory via URL: `http://<server>/static/image.png`.
	// Without `PathPrefix`, you have to access it via URL:
	// `http://<server>/static/static/image.png`.
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(direc),
		PathPrefix: "static",
		Browse:     true,
	}))

	log.Fatal(app.Listen(":3000"))

}

// Utils
// SendFile
// Serves a file from an HTTP file system at the specified path
func TestSendFile(*testing.T) {
	app := fiber.New()

	// Serve static files from the "build" directory using Fiber's built-in middleware.
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(unauth), // Specify the root directory for static files.
		PathPrefix: "build",         // Define the path prefix where static files are served.
	}))

	// For all other routes (wildcard "*"), serve the "index.html" file from the "build" directory.
	app.Use("*", func(ctx *fiber.Ctx) error {
		return filesystem.SendFile(ctx, http.FS(unauth), "build/index.html")
	})

	// Define a route to serve a specific file
	app.Get("/download", func(c *fiber.Ctx) error {
		// Serve the file using SendFile function
		err := filesystem.SendFile(c, http.Dir("your/filesystem/root"), "path/to/your/file.txt")
		if err != nil {
			// Handle the error, e.g., return a 404 Not Found response
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		return nil
	})
	app.Listen(":3000")
}
