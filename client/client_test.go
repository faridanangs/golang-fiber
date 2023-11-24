package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getSomething(c *fiber.Ctx) (err error) {
	agent := fiber.Get("<URL>")
	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"err": errs,
			})
	}

	var something fiber.Map
	err = json.Unmarshal(body, &something)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	return c.Status(statusCode).JSON(something)
}

// Post something
func createSomething(c *fiber.Ctx) (err error) {
	agent := fiber.Post("<URL>")
	agent.Body(c.Body()) // set body received by request
	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	// pass status code and body received by the proxy
	return c.Status(statusCode).Send(body)
}

// Parse menginisialisasi HostClient.
func TestClient(t *testing.T) {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	req.SetRequestURI("http://anangs.com")
	req.Header.SetMethod(fiber.MethodGet)

	if err := agent.Parse(); err != nil {
		panic(err)
	}

	///// Set
	// Set menyetel key: valueheader yang diberikan.
	// Tanda tangan
	// func (a *Agent) Set(k, v string) *Agent
	// func (a *Agent) SetBytesK(k []byte, v string) *Agent
	// func (a *Agent) SetBytesV(k string, v []byte) *Agent
	// func (a *Agent) SetBytesKV(k []byte, v []byte) *Agent

	// Contoh
	// agent.Set("k1", "v1").
	//     SetBytesK([]byte("k1"), "v1").
	//     SetBytesV("k1", []byte("v1")).
	//     SetBytesKV([]byte("k2"), []byte("v2"))
	agent.Set("k1", "v1").SetBytesV("name", []byte("v1"))

	///// Add
	// 	Add menambahkan key: valueheader yang diberikan.
	// 	Beberapa header dengan kunci yang sama dapat ditambahkan dengan fungsi ini.
	// Tanda tangan
	// func (a *Agent) Add(k, v string) *Agent
	// func (a *Agent) AddBytesK(k []byte, v string) *Agent
	// func (a *Agent) AddBytesV(k string, v []byte) *Agent
	// func (a *Agent) AddBytesKV(k []byte, v []byte) *Agent

	// Contoh
	// agent.Add("k1", "v1").
	//     AddBytesK([]byte("k1"), "v1").
	//     AddBytesV("k1", []byte("v1")).
	//     AddBytesKV([]byte("k2"), []byte("v2"))
	// // Headers:
	// // K1: v1
	// // K1: v1
	// // K1: v1
	// // K2: v2
	agent.Add("k1", "v1").AddBytesV("k1", []byte("v1"))

	////// ConnectionClose
	// 	ConnectionClose menambahkan Connection: closeheader.
	// Tanda tangan
	// func (a *Agent) ConnectionClose() *Agent

	// Contoh
	// agent.ConnectionClose()

	////// UserAgent
	// 	UserAgent menetapkan User-Agentnilai header.
	// Tanda tangan
	// func (a *Agent) UserAgent(userAgent string) *Agent
	// func (a *Agent) UserAgentBytes(userAgent []byte) *Agent

	// Contoh
	agent.UserAgent("fiber")

	////// Cookie
	// 	Cookie menetapkan cookie dalam key: valuebentuk. Cookiesdapat digunakan untuk mengatur beberapa cookie.
	// Tanda tangan
	// func (a *Agent) Cookie(key, value string) *Agent
	// func (a *Agent) CookieBytesK(key []byte, value string) *Agent
	// func (a *Agent) CookieBytesKV(key, value []byte) *Agent
	// func (a *Agent) Cookies(kv ...string) *Agent
	// func (a *Agent) CookiesBytesKV(kv ...[]byte) *Agent

	// Contoh
	agent.Cookie("k", "v")
	agent.Cookies("k1", "v1", "k2", "v2")

	////// Referer
	// 	Referer menetapkan nilai header Referer. di gunakan untuk mengetahui
	// dari mana user masuk ke web kita di contoh ini user masuk melalui doc fiber
	// Tanda tangan
	// func (a *Agent) Referer(referer string) *Agent
	// func (a *Agent) RefererBytes(referer []byte) *Agent

	// Contoh
	agent.Referer("https://docs.gofiber.io")

	////// ContentType
	// ContentType menetapkan nilai header Tipe Konten.
	// Tanda tangan
	// func (a *Agent) ContentType(contentType string) *Agent
	// func (a *Agent) ContentTypeBytes(contentType []byte) *Agent

	// Contoh
	agent.ContentType(fiber.MIMEApplicationJSON)

	////// Host
	// Header Host menunjukkan tuan rumah atau domain yang diminta.
	// Misalnya, jika Anda membuat permintaan HTTP ke "http://example.com/page", maka header Host akan diatur ke "example.com".
	// 	Host menyetel header Host.
	// Tanda tangan
	// func (a *Agent) Host(host string) *Agent
	// func (a *Agent) HostBytes(host []byte) *Agent

	// Contoh
	agent.Host("example.com")
	/////// BasicAuth
	// 	BasicAuth menetapkan nama pengguna dan kata sandi URI menggunakan HTTP Basic Auth.

	// Tanda tangan
	// func (a *Agent) BasicAuth(username, password string) *Agent
	// func (a *Agent) BasicAuthBytes(username, password []byte) *Agent
	agent.BasicAuth("farid", "anodsauogsdaiy")

	///// Body
	// 	func (a *Agent) BodyString(bodyString string) *Agent
	// func (a *Agent) Body(body []byte) *Agent
	// BodyStream sets request body stream and, optionally body size.
	// If bodySize is >= 0, then the bodyStream must provide exactly bodySize bytes
	// before returning io.EOF.
	// If bodySize < 0, then bodyStream is read until io.EOF.
	// bodyStream.Close() is called after finishing reading all body data
	// if it implements io.Closer.

	// Note that GET and HEAD requests cannot have body.
	// func (a *Agent) BodyStream(bodyStream io.Reader, bodySize int) *Agent
	agent.BodyString("foo=bar")
	agent.Body([]byte("bar=baz"))
	agent.BodyStream(strings.NewReader("body=stream"), -1)

	///// json
	// JSON mengirimkan permintaan JSON dengan mengatur header Content-Type ke application/json.
	// Tanda tangan
	// func (a *Agent) JSON(v interface{}) *Agent
	// Contoh
	// agent.JSON(fiber.Map{"success": true})

	////// xml
	// XML mengirimkan permintaan XML dengan menyetel header Tipe Konten ke application/xml.
	// Tanda tangan
	// func (a *Agent) XML(v interface{}) *Agent
	// Contoh
	// agent.XML(fiber.Map{"success": true})

	//// form
	// Formulir mengirimkan permintaan formulir dengan mengatur header Tipe Konten ke application/x-www-form-urlencoded.

	// Tanda tangan
	// Form sends form request with body if args is non-nil.
	// It is recommended obtaining args via AcquireArgs and release it
	// manually in performance-critical code.
	// func (a *Agent) Form(args *Args) *Agent

	// Contoh
	args := fiber.AcquireArgs()
	args.Set("foo", "bar")

	agent.Form(args)
	fiber.ReleaseArgs(args)

	//// multipartform
	// 	MultipartForm mengirimkan permintaan formulir multipart dengan menyetel header Tipe Konten ke multipart/form-data. Permintaan ini dapat mencakup nilai kunci dan file.

	// Tanda tangan
	// // MultipartForm sends multipart form request with k-v and files.
	// // It is recommended to obtain args via AcquireArgs and release it
	// // manually in performance-critical code.
	// func (a *Agent) MultipartForm(args *Args) *Agent
	// Contoh
	args.Set("foo", "bar")
	agent.MultipartForm(args)
	fiber.ReleaseArgs(args)
	// Fiber menyediakan beberapa metode untuk mengirim file. Perhatikan bahwa mereka harus dipanggil sebelumnya MultipartForm.

	/////// sendfile
	// 	SendFile membaca file dan menambahkannya ke permintaan formulir multi-bagian. Sendfiles dapat digunakan untuk menambahkan banyak file.

	// Tanda tangan
	// func (a *Agent) SendFile(filename string, fieldname ...string) *Agent
	// func (a *Agent) SendFiles(filenamesAndFieldnames ...string) *Agent

	// Contoh
	agent.SendFile("f", "field name").SendFiles("f1", "field name1", "f2").MultipartForm(nil)

	///// timeout
	// Batas waktu menetapkan durasi batas waktu permintaan.
	// Tanda tangan
	// func (a *Agent) Timeout(timeout time.Duration) *Agent
	// Contoh
	agent.Timeout(time.Second)

	//////Reuse
	// 	Penggunaan kembali memungkinkan instans Agen untuk digunakan kembali setelah satu permintaan. Jika agen dapat digunakan kembali, maka harus dilepaskan secara manual ketika sudah tidak digunakan lagi.
	// Tanda tangan
	// func (a *Agent) Reuse() *Agent
	// Contoh
	agent.Reuse()

	///// InsecureSkipVerify
	// 	InsecureSkipVerify mengontrol apakah Agen memverifikasi rantai sertifikat server dan nama host.
	// Tanda tangan
	// func (a *Agent) InsecureSkipVerify() *Agent
	// Contoh
	agent.InsecureSkipVerify()

	///// TLSConfig
	// 	TLSConfig menyetel konfigurasi tls.

	// Tanda tangan
	// func (a *Agent) TLSConfig(config *tls.Config) *Agent

	// Contoh
	// Create tls certificate
	cer, _ := tls.LoadX509KeyPair("pem", "key")

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
	}

	agent.TLSConfig(config)

	////////MaxRedirectsCount
	// 	MaxRedirectsCount menetapkan jumlah pengalihan maksimal untuk GET dan HEAD.
	// Tanda tangan
	// func (a *Agent) MaxRedirectsCount(count int) *Agent
	// Contoh
	agent.MaxRedirectsCount(7)

	////// JSONEncoder
	// 	JSONEncoder menyetel pembuat enkode json khusus.

	// Tanda tangan
	// func (a *Agent) JSONEncoder(jsonEncoder utils.JSONMarshal) *Agent

	// Contoh
	agent.JSONEncoder(json.Marshal)

	//// JSONDecoder
	// 	JSONDecoder menyetel dekoder json khusus.
	// Tanda tangan
	// func (a *Agent) JSONDecoder(jsonDecoder utils.JSONUnmarshal) *Agent
	// Contoh
	agent.JSONDecoder(json.Unmarshal)

	///// request
	// 	Permintaan mengembalikan contoh permintaan Agen.
	// Tanda tangan
	// func (a *Agent) Request() *Request
	// Contoh
	// reque := agent.Request()
	// fmt.Println(reque)

	//////SetResponse
	// 	SetResponse menetapkan respons khusus untuk instans Agen. Disarankan untuk mendapatkan respons khusus melalui AcquireResponse dan merilisnya secara manual dalam kode yang penting bagi kinerja.
	// Tanda tangan
	// func (a *Agent) SetResponse(customResp *Response) *Agent
	resp := fiber.AcquireResponse()
	agent.SetResponse(resp)
	fiber.ReleaseResponse(resp)

	///// Dest
	// 	Tujuan menetapkan tujuan khusus. Isi dest akan digantikan oleh badan respons, jika dest terlalu kecil, potongan baru akan dialokasikan.
	// Tanda tangan
	// func (a *Agent) Dest(dest []byte) *Agent {
	// Contoh
	agent.Dest(nil)

	///// Bytes
	// 	Bytes mengembalikan kode status, isi byte, dan kesalahan url.
	// Tanda tangan
	// func (a *Agent) Bytes() (code int, body []byte, errs []error)

	// Contoh
	// code, body, errs := agent.Bytes()
	//  fmt.Println("code ", code)
	// fmt.Println("body ", body)
	// fmt.Println("error ", errs)

	///// String
	// 	String mengembalikan kode status, isi string, dan kesalahan url.
	// Tanda tangan
	// func (a *Agent) String() (int, string, []error)
	// code, body, errs := agent.String()
	// fmt.Println("code ", code)
	// fmt.Println("body ", body)
	// fmt.Println("error ", errs)

	////// Struct
	// 	Struct mengembalikan kode status, isi byte, dan kesalahan url. Dan isi byte tidak akan di-unmarshall ke v yang diberikan.
	// Tanda tangan
	// func (a *Agent) Struct(v interface{}) (code int, body []byte, errs []error)
	var d interface{}
	code, body, errs := agent.Struct(&d)
	fmt.Println("code ", code)
	fmt.Println("body ", body)
	fmt.Println("error ", errs)
	//// RetryIf
	// 	RetryIf mengontrol apakah percobaan ulang harus dilakukan setelah terjadi kesalahan. Secara default, akan menggunakan fungsi isIdempotent dari fasthttp
	// Tanda tangan
	// func (a *Agent) RetryIf(retryIf RetryIfFunc) *Agent

	// code, _, _ = agent.Get("http://localhost:8080")

	// heder := agent.Request().Header.String()
	// fmt.Println(heder)

}
