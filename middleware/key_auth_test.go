package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

var (
	apiKey        = "hello world"
	protectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/authenticated$"),
		regexp.MustCompile("^/auth2$"),
	}
)

func validateApiKey(c *fiber.Ctx, key string) (bool, error) {
	hasedApiKey := sha256.Sum256([]byte(apiKey))
	hasedKey := sha256.Sum256([]byte(key))
	if subtle.ConstantTimeCompare(hasedApiKey[:], hasedKey[:]) == 1 {
		return true, nil
	} else {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}
}

func TestKeyAuth(t *testing.T) {
	app := fiber.New()

	app.Use(keyauth.New(keyauth.Config{
		// jika kita ingin mengecek kode ini kita gunakan
		// curl http://localhost:3000  ini dia error karna tidak ada access tokenya
		// missing or malformed API Key

		//curl --cookie "access_token=hello world" http://localhost:3000
		// successfully  ini berhasil karna kita beri acces token yang kemudian akan masuk ke dalam validator
		// menjadi key di parameter ke dua

		KeyLookup: "cookie:access_token",
		Validator: validateApiKey,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("successfully")
	})

	app.Listen(":3000")
}

// Authenticate only certain endpoints
// If you want to authenticate only certain endpoints, you can use the Config of keyauth
// and apply a filter function (eg. authFilter) like so

func authFilter(c *fiber.Ctx) bool {
	originalUrl := strings.ToLower(c.OriginalURL())

	fmt.Println(originalUrl)
	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalUrl) {
			// jika kita mereturn false maka dia akan lanjut menjalankan validator
			// sehingga di validator dia mengecek apakah key cookienya sama dengan key yang kita tentukan
			return false

			// namun jika kita mereturn true dia akan langsung lanjut mengeksekusi app.Get tanpa menjalankan validator
			// return true
		}
	}
	return true
}
func TestKeyAuthEndpoint(t *testing.T) {
	app := fiber.New()

	app.Use(keyauth.New(keyauth.Config{
		Next:      authFilter,
		KeyLookup: "cookie:access_token",
		Validator: validateApiKey,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("welcome")
	})
	app.Get("/authenticated", func(c *fiber.Ctx) error {
		return c.SendString("successfully")
	})
	app.Get("/auth2", func(c *fiber.Ctx) error {
		return c.SendString("successfully auth 2")
	})

	app.Listen(":3000")
}

// Specifying middleware in the handler
const apiKeySpesific = "super-key"

func TestSpecifyingMiddleware(t *testing.T) {
	app := fiber.New()

	authMiddleKey := keyauth.New(keyauth.Config{
		// jika kita tidak menggunakan keylookup dia secara otomatis akan menggunakan
		// "header:Authorization" sehingga cara kita mengirim reuest kita harus menggunakan
		// beanr sebagai nilai defaultnya contoh requestnya seperti ini
		//  curl --header "Authorization: super-key"  http://localhost:3000/allowed

		// jika kita menyetel nilai dari lokupya kita harus memanggilnya seperti ini
		// curl --cookie "access=super-key"  http://localhost:3000/allowed
		KeyLookup: "cookie:access",
		Validator: func(c *fiber.Ctx, s string) (bool, error) {
			hasedApiKey := sha256.Sum256([]byte(apiKeySpesific))
			hasedApi := sha256.Sum256([]byte(s))

			if subtle.ConstantTimeCompare(hasedApiKey[:], hasedApi[:]) == 1 {
				return true, nil
			} else {
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}
		},
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})
	app.Get("/allowed", authMiddleKey, func(c *fiber.Ctx) error {
		return c.SendString("successfuly")
	})

	app.Listen(":3000")
}

// Default Config
// var ConfigDefault = Config{
//     SuccessHandler: func(c *fiber.Ctx) error {
//         return c.Next()
//     },
//     ErrorHandler: func(c *fiber.Ctx, err error) error {
//         if err == ErrMissingOrMalformedAPIKey {
//             return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
//         }
//         return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired API Key")
//     },
//     KeyLookup:  "header:" + fiber.HeaderAuthorization,
//     AuthScheme: "Bearer",
//     ContextKey: "token",
// }
