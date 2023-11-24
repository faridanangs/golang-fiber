package middleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func TestSession(t *testing.T) {
	app := fiber.New()
	store := session.New(session.Config{
		Expiration:     time.Second * 10,
		CookieHTTPOnly: true,
		CookiePath:     "/",
		KeyLookup:      "header:apiKey",
		// CookieSecure:   true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		// set key and value
		sess.Set("name", "farid anang samudra")
		sess.Set("id", "2121312")

		// get value
		name := sess.Get("name")
		fmt.Println(name)

		// get all keys
		keys := sess.Keys()
		fmt.Println("before", keys)

		// Destroy session
		// menghapus semua data di dalam session
		// if err := sess.Destroy(); err != nil {
		// 	panic(err)
		// }

		// Delete key
		sess.Delete("name")

		keys = sess.Keys()
		fmt.Println("after", keys)

		// Sets a specific expiration for this session
		sess.SetExpiry(time.Second * 10)

		if err := sess.Save(); err != nil {
			panic(err)
		}

		return c.SendString(fmt.Sprintf("hello %v", name))
	})
	app.Listen(":3000")

}
