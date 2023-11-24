package middleware

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// middleware yang digunakan untuk membatasi jumlah permintaan yang dapat
// diterima oleh server dalam interval waktu tertentu. Hal ini berguna untuk
// mencegah serangan penolakan layanan (DDoS) atau penggunaan berlebihan dari sumber daya server.
// Berikut adalah beberapa fungsi utama dari Limiter di Fiber:

// Pembatasan Lalu Lintas (Rate Limiting): Memungkinkan Anda untuk
// menetapkan batasan berapa banyak permintaan yang dapat diterima oleh
// server dalam interval waktu tertentu. Misalnya, Anda dapat membatasi
// server Anda untuk menerima hanya 100 permintaan per menit.

// Mengatur Kode Status Terkait: Limiter dapat dikonfigurasi untuk
// memberikan kode status berbeda (seperti 429 Too Many Requests)
// ketika batasan lalu lintas terlampaui.

// Manajemen Lebih Lanjut: Selain pembatasan lalu lintas,
// Limiter dalam Fiber juga dapat dikonfigurasi untuk memperhatikan header
// khusus atau melakukan tindakan khusus berdasarkan kondisi tertentu

func TestLimitter(t *testing.T) {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        2,
		Expiration: 2 * time.Minute,
		// jika request user lebih dari max dia akan mengirimkan limitedreached ke user
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("request mentok")
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3000")
}
