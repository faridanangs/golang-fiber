package middleware

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// ETag middleware for Fiber that lets caches be more efficient and save
// bandwidth, as a web server does not need to resend a full response if the content has not changed.

// ETag adalah header HTTP yang digunakan untuk memverifikasi apakah
// konten dari sumber daya telah berubah sejak terakhir kali diambil oleh
//
//	klien. Ini adalah mekanisme cache kontrol yang memungkinkan server
//	mengirimkan header ETag bersama dengan respons. Kemudian, klien dapat
//
// menyimpan ETag dan menggunakannya dalam permintaan berikutnya untuk
// memeriksa apakah konten telah berubah sejak terakhir kali diambil.
// Saat klien mengirim permintaan, ia dapat menyertakan header
// If-None-Match dengan nilai ETag yang telah ia terima sebelumnya.
// Jika konten masih sama dengan nilai ETag yang dikirim oleh klien,
// server dapat mengirim respons 304 Not Modified tanpa mengirim ulang seluruh konten, sehingga menghemat bandwidth.
func TestETag(t *testing.T) {
	app := fiber.New()

	// sisi server
	app.Get("/", func(c *fiber.Ctx) error {
		etag := fmt.Sprintf("%x", md5.Sum([]byte("heo world")))
		c.Set("ETag", etag)
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func TestEtagClient(t *testing.T) {
	app := fiber.New()
	app.Get("/etag", func(c *fiber.Ctx) error {
		// Kirim permintaan dengan ETag yang diterima sebelumnya
		c.Set("If-None-Match", "5eb63bbbe01eeed093cb22bb8f5acdc3")

		// melakukan permintan ke server
		resp, err := http.Get("http://localhost:3000")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// / Periksa apakah status respons adalah 304 Not Modified
		if resp.StatusCode == fiber.StatusNotModified {
			return c.SendString("konten tidak berubah")
		}

		// Baca dan kirim konten dari respons
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.SendString(string(body))
	})

	app.Listen(":2000")
}
