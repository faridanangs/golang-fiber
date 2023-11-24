package context

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestApplication(t *testing.T) {
	app := fiber.New()

	app.Get("/stack", func(c *fiber.Ctx) error {
		fmt.Println(c.JSON(c.App().Stack()))
		return c.SendStatus(200)
	})
	app.Listen(":8080")
}
func TestAppend(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Append("Link", "https://google.com")
		c.Append("Link", "test")
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})
	app.Listen(":8080")
}

func TestAtachment(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// untuk mengirim data img ke header dan mengunduhnya kita bisa menggunakan atachmen dan tempat file imgnya
		c.Type("image/png")
		c.Attachment("./image/logo1.png")
		fmt.Println(c.GetRespHeaders())
		fmt.Println(c.BaseURL())
		return c.SendStatus(200)
	})
	app.Listen(":8080")
}
func TestBind(t *testing.T) {
	// Tambahkan vars ke tampilan default pengikatan peta var ke
	// mesin templat. Variabel dibaca dengan metode Render dan dapat
	// ditimpa func (c *Ctx) Bind(vars Map) error

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Bind(fiber.Map{
			"Title": "Hello, World!",
		})
		return nil
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("xxx.tmpl", fiber.Map{}) // Render will use Title variable
	})

	app.Listen(":8080")
}

func TestBodyRaw(t *testing.T) {
	// Di dalam aplikasi Fiber, ketika ada permintaan POST pada rute '/', pengontrol (handler) yang diberikan akan dijalankan.
	// Di dalam handler, c.BodyRaw() digunakan untuk mendapatkan data body
	// dari permintaan dalam bentuk byte slice. Kemudian, data ini dikirimkan
	// kembali sebagai respons menggunakan c.Send(). Dalam kasus ini, data yang
	// dikirimkan kembali adalah []byte("user=john").

	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		c.Attachment("./image/logo1.png")
		c.Append("Link", "Hello, World")
		return c.Send(c.BodyRaw())
	})
	app.Listen(":8080")

}
func TestBody(t *testing.T) {
	app := fiber.New()

	// // ini adalah kode perintah untuk menjalankan kode body supaya bisa di kompres dan kode ini husus powershalel
	// $body = 'user=john'
	// PS D:\Pemrograman\Golang\src\golang-fiber\context> $body | Out-File -FilePath body.txt
	// >> Compress-Archive -Path .\body.txt -DestinationPath .\body.zip
	// >> Invoke-RestMethod -Uri "http://localhost:8080" -Method POST -InFile .\body.zip -Headers @{"Content-Encoding"="gzip"}

	app.Post("/", func(c *fiber.Ctx) error {
		// Mendekompresi body dari permintaan POST berdasarkan Content-Encoding
		// dan mengembalikan konten mentahnya:
		c.Attachment("./image/logo1.png")
		c.Append("Link", "Hello, World")
		return c.Send(c.Body())
	})
	app.Listen(":8080")

}

type Person struct {
	Name string `json:"name" xml:"name" form:"name"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

func TestBodyParse(t *testing.T) {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		p := new(Person)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		log.Println(p.Name)
		log.Println(p.Pass)

		return c.SendStatus(200)
	})

	app.Listen(":3000")

}

func TestDelCookie(t *testing.T) {
	app := fiber.New()

	app.Get("/set", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:        "token",
			Value:       "random lol",
			Expires:     time.Now().Add(24 * time.Hour),
			HTTPOnly:    true,
			SameSite:    "lax",
			SessionOnly: true, // ketika user menutup aplikasi cookie di hapus
		})
		return c.SendStatus(200)

	})

	app.Get("/dedl", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			HTTPOnly: true,
			Expires:  time.Now().Add(-(time.Hour * 2)),
			SameSite: "lax",
		})
		return c.SendStatus(200)
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})
	app.Listen(":3000")
}

func TestGetCookies(t *testing.T) {
	app := fiber.New()

	app.Get("/set", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:        "token",
			Value:       "hello world",
			Expires:     time.Now().Add(time.Hour * 24),
			HTTPOnly:    true,
			SameSite:    "lax",
			SessionOnly: true,
		})
		return c.SendStatus(200)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// untuk menambil nilai yang ada di dalam cookie kita menggunakan cookies
		name := c.Cookies("token")
		// untuk mengasih nilai default ke dalam cookie kita berika nilai setelah keynya
		empty := c.Cookies("data", "kosong sayang")
		return c.SendString(name + " " + empty)
	})
	app.Listen(":3000")
}

// // download
// Mentransfer file dari jalur sebagai attachment.
// Biasanya, browser akan meminta pengguna untuk mengunduh. Secara
// default, parameter header Content-Dispositionfilename= adalah jalur
// file ( biasanya muncul di dialog browser ) .

func TestDowload(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// return c.Download("./image/05.png", "05.png")
		return c.Download("./image/loli.txt", "loli.txt")
	})
	app.Listen(":3000")
}

// Mengembalikan header permintaan HTTP yang ditentukan oleh kolom.
// dan mengambil data dari dalam header http
func TestGetHeaderHttp(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Get("Content-Type"))      // "text/plain"
		fmt.Println(c.Get("CoNtEnT-TypE"))      // "text/plain"
		fmt.Println(c.Get("something", "john")) // "john"
		c.Set("loli", "hello loli[pop]")

		// response header
		fmt.Println(c.GetRespHeaders())

		// request header
		fmt.Println(c.GetReqHeaders())

		// mengambil data yang sesuai dengan key yang di kirim
		fmt.Println(c.GetRespHeader("Content-Type"))
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func TestGetRouteUrl(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("helo")
	}).Name("home")

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("id"))
	}).Name("user.show")

	app.Get("/test", func(c *fiber.Ctx) error {
		// nilai dari routeName di dalam GetRouteURL harus sama dengan nama url yang di tuju
		// contoh user.show dan dia mengembalikan url path dari url yang di tuju dan memasukan nilai ke dalam
		// url yang di tuju menggunakan fiber.Map{"pathnya": nilainya}
		value, _ := c.GetRouteURL("user.show", fiber.Map{"id": 1})
		return c.SendString(value)
	})

	app.Listen(":3000")
}

// Mengembalikan nama host yang berasal dari header HTTP Host
func TestHostName(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Hostname())
		return c.SendStatus(200)
	}).Name("home")

	app.Listen(":3000")
}

// Mengembalikan alamat IP jarak jauh dari permintaan.
func TestIp(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.IP())
		return c.SendStatus(200)
	}).Name("home")

	app.Listen(":3000")
}

func TestIps(t *testing.T) {
	// Saat mendaftarkan header permintaan proxy di aplikasi fiber, alamat ip header dikembalikan (konfigurasi Fiber)
	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
	})
	// Mengembalikan array alamat IP yang ditentukan dalam header permintaan X-Forwarded-For .
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.IPs())
		return c.SendStatus(200)
	}).Name("home")

	app.Listen(":3000")
}

// Mengembalikan tipe konten yang cocok, jika bidang header HTTP Tipe Konten permintaan masuk cocok dengan tipe MIME yang ditentukan oleh parameter tipe.
func TestIs(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Is("html"))
		fmt.Println(c.Is(".html"))
		fmt.Println(c.Is("json"))
		fmt.Println(c.Is("css"))
		return c.SendStatus(200)
	}).Name("home")

	app.Listen(":3000")
}

// Mengembalikan nilai benar jika permintaan datang dari localhost
func TestIsFromLocal(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.IsFromLocal())
		return c.SendStatus(200)
	}).Name("home")

	app.Listen(":3000")
}

type User struct {
	Username string
	Age      int
}

// Mengonversi antarmuka atau string apa pun menjadi JSON menggunakan paket pengkodean/json .
func TestJson(t *testing.T) {
	app := fiber.New()

	app.Get("/json", func(c *fiber.Ctx) error {
		// data := &User{
		// 	Username: "farid anang samudra",
		// 	Age:      17,
		// }
		// c.JSON(data)

		// selain menggunkan struct kita juga bisa menggunakan map bawaan dari fiber
		c.JSON(fiber.Map{
			"name": "farid anang samudra",
			"Age":  10,
		})
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// jika kita menggunakan json p dia memiliki nama di awalnya seperti ini
// callback({"Username":"farid anang samudra","Age":17}); jika kita tidak memberi nama
// maka defaultnya callback
func TestJsonP(t *testing.T) {
	app := fiber.New()

	app.Get("/json", func(c *fiber.Ctx) error {
		data := &User{
			Username: "farid anang samudra",
			Age:      17,
		}
		// c.JSONP(data) callback({"Username":"farid anang samudra","Age":17});
		c.JSONP(data, "anangs") // anangs({"Username":"farid anang samudra","Age":17});
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
func TestLink(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Links(
			// link ini bisa menerima lebih dari satu dan cara mengirimkan datanya urlnya, relnya
			// ini adalah hasil dari nilai yang kita buat
			// Link:<http://google.com>; rel="google",....
			"http://google.com", "google",
			"http://anangs.com", "anangs",
		)
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func TestLocals(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// untuk memasukan nilai ke dalam locals kita beri key dan nilainya
		c.Locals("user", "admin")
		return c.Next()
	})

	app.Get("/admin", func(c *fiber.Ctx) error {
		// untuk memanggil nilai yang ada di dalam locals kita masukan saja keynya
		if c.Locals("user") == "admin" {
			return c.Status(fiber.StatusOK).SendString("welcome admin")
		}
		return c.SendStatus(fiber.StatusForbidden)
	})

	app.Listen(":3000")

}

// Menyetel header HTTP Lokasi respons ke parameter jalur yang ditentukan.
func TestLocation(t *testing.T) {
	app := fiber.New()
	app.Get("/foo", func(c *fiber.Ctx) error {
		c.Location("https://google.com")
		c.Location("/foo/bar")
		return c.SendStatus(200)
	})

	app.Listen(":3000")

}

// Mengembalikan string yang sesuai dengan metode HTTP permintaan: GET, POST, PUT, dan seterusnya.
// Secara opsional, Anda dapat mengganti metode dengan meneruskan sebuah string.

// func (c *Ctx) Method(override ...string) string

func TestMethod(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendStatus(200)
	})
	app.Post("/user", func(c *fiber.Ctx) error {
		c.Method("GET")
		fmt.Println(c.Method())
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func TestMultipartForm(t *testing.T) {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		// Parse the multipart form:
		if form, err := c.MultipartForm(); err == nil {
			// => *multipart.Form

			if token := form.Value["token"]; len(token) > 0 {
				// Get key value:
				fmt.Println(token[0])
			}

			// Get all files from "documents" key:
			files := form.File["documents"]
			// => []*multipart.FileHeader

			// Loop through files:
			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				// => "tutorial.pdf" 360641 "application/pdf"

				// Save the files to disk:
				if err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename)); err != nil {
					return err
				}
			}
		}

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// // Next()
// Ketika Next dipanggil, ia mengeksekusi metode berikutnya dalam tumpukan yang
// cocok dengan rute saat ini. Anda dapat meneruskan struct kesalahan dalam
// metode yang akan mengakhiri rangkaian dan memanggil penangan kesalahan .

func TestNext(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		fmt.Println("ke satu 1")
		return c.Next()
	})
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("ke dua 2")
		time.Sleep(5 * time.Second)
		return c.Next()
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString("Hello World")
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Mengembalikan URL permintaan
// func (c *Ctx) OriginalURL() string
func TestOriginalURL(t *testing.T) {
	app := fiber.New()

	// GET http://example.com/search?q=something
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.OriginalURL()) // "/search?q=something"
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}

// Metode dapat digunakan untuk mendapatkan parameter rute, Anda dapat memberikan
// nilai default opsional yang akan dikembalikan jika kunci param tidak ada.
// INFO
// Defaultnya adalah string kosong ( "") , jika paramnya tidak ada.
// func (c *Ctx) Params(key string, defaultValue ...string) string

func TestParams(t *testing.T) {
	app := fiber.New()
	// GET http://example.com/user/1
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("id")) // "1"
		return c.SendString(c.Params("id"))
	})

	app.Get("/user/name/*", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("*"))

		return c.SendString(c.Params("*"))
	})

	app.Get("/v1/*/user/*", func(c *fiber.Ctx) error {
		// ROUTE: /v1/*/shop/*
		// GET:   /v1/brand/4/user/blue/xs
		// kita gunakan *1 dan *2 .. untuk menampung url yang kirim di url * dan berakhir di * berikutnya contoh
		//  GET:   /v1/(*1)brand/4/user/(*2)blue/xs
		fmt.Println(c.Params("*1")) // "brand/4"
		fmt.Println(c.Params("*2")) // blue/xs

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Metode dapat digunakan untuk mendapatkan bilangan bulat dari parameter
// rute. Harap dicatat jika parameter tersebut tidak ada dalam permintaan,
// nol akan dikembalikan. Jika parameter BUKAN angka, nol dan kesalahan akan dikembalikan

// Defaultnya adalah bilangan bulat nol ( 0) , jika paramnya tidak ada.
// func (c *Ctx) ParamsInt(key string) (int, error)

func TestParamsInt(t *testing.T) {
	app := fiber.New()
	// jika kita memasukan string dia errornya seperti inii
	// 0 failed to convert: strconv.Atoi: parsing "dda": invalid syntax
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		fmt.Println(c.ParamsInt("id")) // "1"
		return c.SendString("ok")
	})

	app.Listen(":3000")
}

// Metode ini mirip dengan BodyParser, tetapi untuk parameter jalur.
// Penting untuk menggunakan tag struct "params". Misalnya, jika Anda ingin
// mengurai parameter jalur dengan bidang bernama Pass, Anda akan
// menggunakan bidang struct params:"pass"
// func (c *Ctx) ParamsParser(out interface{}) error

func TestParamsParse(t *testing.T) {
	app := fiber.New()
	// GET http://example.com/user/111
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		data := struct {
			Id int `params:"id"`
		}{}
		c.ParamsParser(&data) // "{"id": 111}"
		return c.SendStatus(data.Id)
	})

	app.Listen(":3000")
}

// Berisi bagian jalur dari URL permintaan. Secara opsional,
// Anda dapat mengganti jalur dengan meneruskan string.
// Untuk pengalihan internal, Anda mungkin ingin memanggil RestartRouting daripada Next .
// func (c *Ctx) Path(override ...string) string

func TestPath(t *testing.T) {
	app := fiber.New()
	app.Get("/user", func(c *fiber.Ctx) error {
		c.Path("/json")
		// c.Path() /sjon
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Berisi string protokol permintaan: httpatau httpsuntuk permintaan TLS
func TestProtocol(t *testing.T) {
	app := fiber.New()
	app.Get("/user", func(c *fiber.Ctx) error {
		fmt.Println(c.Protocol()) // http
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// kita menggunakan queries untuk mengambil data yang ada di dalam url dengan cara
// mengambil nama url yang memiliki = contoh key=value&key2=value2&...
func TestQueries(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// GET http://exce.id/?first=lolii&second=udin&id=
		q := c.Queries()
		first := q["first"]   //"lolii"
		second := q["second"] // "udin"
		id := q["id"]         // " "
		fmt.Println(first, second, id)
		return c.SendStatus(200)
	})

	// GET http://example.com/?list_a=1&list_a=2&list_a=3&list_b[]=1&list_b[]=2&list_b[]=3&list_c=1,2,3
	// jika key urlnya sama maka dia akan mengambil url yang terakhir
	app.Get("/", func(c *fiber.Ctx) error {
		m := c.Queries()
		one := m["list_a"]   // "3"
		two := m["list_b[]"] // "3"
		three := m["list_c"] // "1,2,3"
		fmt.Println(one, two, three)
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Properti ini adalah objek yang berisi properti untuk setiap parameter string kueri di rute, Anda dapat meneruskan nilai default opsional yang akan dikembalikan jika kunci kueri tidak ada.
// Jika tidak ada string kueri, ia mengembalikan string kosong .
// func (c *Ctx) Query(key string, defaultValue ...string) string

func TestQuery(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Query("age", "17"))
		fmt.Println(c.Query("name", "anangs"))
		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}

// Properti ini adalah objek yang berisi properti untuk setiap parameter boolean kueri di rute, Anda dapat meneruskan nilai default opsional yang akan dikembalikan jika kunci kueri tidak ada.
// PERINGATAN
// Harap dicatat jika parameter tersebut tidak ada dalam permintaan, false akan dikembalikan. Jika parameternya bukan boolean, parameter tersebut tetap dicoba untuk dikonversi dan biasanya dikembalikan sebagai false.
// func (c *Ctx) QueryBool(key string, defaultValue ...bool) bool

func TestQueryBool(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// http://localhost:3000/?name=true&age=17
		fmt.Println(c.QueryBool("age"))
		fmt.Println(c.QueryBool("name"))
		fmt.Println(c.QueryBool("udin", true))

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}

// Properti ini adalah objek yang berisi properti untuk setiap parameter kueri float64 di rute, Anda dapat meneruskan nilai default opsional yang akan dikembalikan jika kunci kueri tidak ada.
// PERINGATAN
// Harap dicatat jika parameter tersebut tidak ada dalam permintaan, nol akan dikembalikan. Jika parameternya bukan angka, tetap dicoba dikonversi dan biasanya dikembalikan sebagai 1.
// INFO
// Defaultnya adalah float64 zero ( 0) , jika paramnya tidak ada.
// func (c *Ctx) QueryFloat(key string, defaultValue ...float64) float64

func TestQueryFloat(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// http://localhost:3000/?name=true&age=17.12123
		fmt.Println(c.QueryFloat("age"))       // 17.12123
		fmt.Println(c.QueryFloat("name"))      // 0
		fmt.Println(c.QueryFloat("count", 10)) // 0

		// Properti ini adalah objek yang berisi properti untuk setiap parameter bilangan bulat kueri di rute, Anda dapat meneruskan nilai default opsional yang akan dikembalikan jika kunci kueri tidak ada.
		fmt.Println(c.QueryInt("age"))       // 0 karna ini float bukan int
		fmt.Println(c.QueryInt("name"))      // 0
		fmt.Println(c.QueryInt("count", 10)) // 0

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}

// Metode ini mirip dengan BodyParser , tetapi untuk parameter kueri. Penting untuk menggunakan tag struct "query". Misalnya, jika Anda ingin mengurai parameter kueri dengan bidang bernama Pass, Anda akan menggunakan bidang struct sebesar query:"pass".
// func (c *Ctx) QueryParser(out interface{}) error
func TestQueryParser(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// http://localhost:3000/?name=farid&age=17&status=single

		// untuk menangkap nilai dati query yang akan di parse kita gunaka query;"nama_pathnya"
		p := struct {
			Name   string `query:"name"`
			Age    int    `query:"age"`
			Status string `query:"status"`
		}{}
		// kemudian kita parser data strucnya supaya bisa di kenali jika terjadi permintaan di user dan jika sama
		// dengan nama query kita di struct maka dia akan dimasukan ke dalamnya
		if err := c.QueryParser(&p); err != nil {
			panic(err)
		}

		log.Println(p.Name)
		log.Println(p.Age)
		log.Println(p.Status)

		// output
		// 2023/10/09 22:09:01 farid
		// 2023/10/09 22:09:01 17
		// 2023/10/09 22:09:01 single
		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}

// Sebuah struct yang berisi tipe dan potongan rentang akan dikembalikan.
// Tanda tangan
// func (c *Ctx) Range(size int) (Range, error)

func TestRange(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		b, _ := c.Range(1000)
		if b.Type == "bytes" {
			for r := range b.Ranges {
				fmt.Println(r)
				// [500, 700]
			}
		}
		return nil
	})

	app.Listen(":3000")
}

// Pengalihan ke URL yang berasal dari jalur tertentu, dengan status tertentu, bilangan bulat positif yang sesuai dengan kode status HTTP.
// Jika tidak ditentukan, status defaultnya adalah 302 Found .
// func (c *Ctx) Redirect(location string, status ...int) error

func TestRedirect(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/teapot")
	})
	app.Get("/teapot", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3000")
}

// yang ditentukan, bilangan bulat positif yang sesuai dengan kode status HTTP.
// INFO
// Jika tidak ditentukan, status defaultnya adalah 302 Found .
// func (c *Ctx) RedirectBack(fallback string, status ...int) error
func TestRedirectBack(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Home page")
	})
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(`<a href="/back">Back</a>`)
	})

	app.Get("/back", func(c *fiber.Ctx) error {
		fmt.Println("Sebelum redirect")
		err := c.RedirectBack("/")
		fmt.Println("Setelah redirect")
		return err
	})

	app.Listen(":3000")
}

// Minta pengembalian penunjuk * fasthttp.Request
// func (c *Ctx) Request() *fasthttp.Request

func TestRequest(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(string(c.Request().Header.Method()))
		return c.SendString("Home page")
	})

	app.Listen(":3000")
}

// Metode ini mirip dengan BodyParser , tetapi untuk header permintaan.
// Penting untuk menggunakan tag struct "reqHeader".
//
//	Misalnya, jika Anda ingin mengurai header permintaan dengan bidang bernama Pass, Anda akan menggunakan bidang struct sebesar reqHeader:"pass".
//
// Tanda tangan
// func (c *Ctx) ReqHeaderParser(out interface{}) error
func TestReqHeaderParser(t *testing.T) {
	app := fiber.New()
	p := struct {
		Name     string   `reqHeader:"name"`
		Pass     string   `reqHeader:"pass"`
		Products []string `reqHeader:"products"`
	}{}
	app.Get("/", func(c *fiber.Ctx) error {
		// cara mengirimkan nilai ke dalamnya
		// curl "http://localhost:3000/" -H "name: john" -H "pass: doe" -H "products: shoe,hat"
		if err := c.ReqHeaderParser(&p); err != nil {
			return err
		}
		fmt.Println(p.Name)
		fmt.Println(p.Products)
		fmt.Println(p.Pass)
		fmt.Println(c.GetRespHeaders())

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Response return the *fasthttp.Response pointer
// func (c *Ctx) Response() *fasthttp.Response

func TestResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().BodyWriter().Write([]byte("hello world response"))
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Daripada mengeksekusi metode selanjutnya saat memanggil Next ,
// RestartRouting memulai ulang eksekusi dari metode pertama yang cocok
// dengan rute saat ini. Ini mungkin berguna setelah mengganti jalur,
// yaitu pengalihan internal. Perhatikan bahwa penangan mungkin dieksekusi
// lagi yang dapat menghasilkan perulangan tak terbatas.

// Tanda tangan
// func (c *Ctx) RestartRouting() error
func TestRestartRouting(t *testing.T) {
	app := fiber.New()

	app.Get("/new", func(c *fiber.Ctx) error {
		return c.SendString("from /new")
	})
	app.Get("/udin", func(c *fiber.Ctx) error {
		return c.SendString("from /udin ngangak")
	})
	app.Get("/old", func(c *fiber.Ctx) error {
		// ketika kita masuk ke url old dia akan menjalankan data yang ada di dalam
		// /new namun urlnya tetap berada di old
		// jika kita tidak menggunakan RestartRouting maka yang akan di render data yang kita return bukan data
		// yang ada di dalam url /new
		c.Path("/new")
		return c.RestartRouting()
	})
	app.Get("/human", func(c *fiber.Ctx) error {
		c.Path("/udin")

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Mengembalikan struct Route yang cocok .

// Tanda tangan
// func (c *Ctx) Route() *Route

func TestRouter(t *testing.T) {
	app := fiber.New()

	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		r := c.Route()
		fmt.Println(r.Method, r.Path, r.Params, r.Handlers)
		return c.SendString("ook")
	})
	//// Jangan mengandalkan c.Route()middleware sebelum menelepon c.Next()- c.Route()mengembalikan rute yang terakhir dieksekusi .
	// func() fiber.Handler {
	// 	return func(c *fiber.Ctx) error {
	// 		before := c.Route().Path
	// 		err := c.Next()
	// 		after := c.Route().Path
	// 		return err
	// 	}
	// }()

	app.Listen(":3000")
}

// Metode ini digunakan untuk menyimpan file multi-bagian ke disk.
// Tanda tangan
// func (c *Ctx) SaveFile(fh *multipart.FileHeader, path string) error
func TestSaveFile(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["context"]

			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

				// Save the files to disk:
				if err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename)); err != nil {
					return err
				}
			}
		}
		return nil
	})
	app.Listen(":3000")
}

// Metode ini digunakan untuk menyimpan file multi-bagian ke sistem penyimpanan eksternal.
// Tanda tangan
// func (c *Ctx) SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage Storage) error
func TestSaveFileToStorage(t *testing.T) {
	app := fiber.New()
	// gak tau ini memory.New()nya gak ada
	// storage := memory.New()

	app.Post("/", func(c *fiber.Ctx) error {
		// Parse the multipart form:
		if form, err := c.MultipartForm(); err == nil {
			// => *multipart.Form

			// Get all files from "documents" key:
			files := form.File["documents"]
			// => []*multipart.FileHeader

			// Loop through files:
			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				// => "tutorial.pdf" 360641 "application/pdf"

				// Save the files to storage:
				// if err := c.SaveFileToStorage(file, fmt.Sprintf("./%s", file.Filename), storage); err != nil {
				// 	return err
				// }
			}
			return err
		}
		return nil
	})
	app.Listen(":3000")
}

// Menyetel isi respons HTTP.
// Tanda tangan
// func (c *Ctx) Send(body []byte) error
func TestSend(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("hello"))
	})

	// Fiber juga menyediakan SendStringmetode SendStreamuntuk input mentah.
	// TIP
	// Gunakan ini jika Anda tidak memerlukan pernyataan tipe, disarankan untuk kinerja yang lebih cepat .
	// func (c *Ctx) SendString(body string) error
	// func (c *Ctx) SendStream(stream io.Reader, size ...int) error
	app.Get("/lol", func(c *fiber.Ctx) error {
		return c.SendStream(bytes.NewReader([]byte("hello lolii stream")))
	})

	app.Listen(":3000")
}

// Mentransfer file dari jalur yang diberikan. Menyetel bidang header HTTP respons Tipe Konten berdasarkan ekstensi nama file .
// PERINGATAN
// Metode tidak menggunakan gzipping secara default, setel ke true untuk mengaktifkan.
// Tanda tangan
// func (c *Ctx) SendFile(file string, compress ...bool) error

func TestSendFile(t *testing.T) {
	app := fiber.New()
	app.Get("/not-found", func(c *fiber.Ctx) error {
		return c.SendFile("../public/404.htm", true)
		// return c.SendFile("../public/hello.htm")
	})

	app.Listen(":3000")
}

func TestSet(t *testing.T) {
	app := fiber.New()
	app.Get("/set-header", func(c *fiber.Ctx) error {
		c.Set("content-type", "text/html")
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}

func TestStatus(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// return c.Status(fiber.StatusBadRequest).SendString("hello world")
		return c.Status(fiber.StatusBadRequest).SendFile("../image/logo1.png")
	})
	app.Listen(":3000")
}

// Mengembalikan sepotong string subdomain dalam nama domain permintaan.
// Offset subdomain properti aplikasi, yang defaultnya adalah 2, digunakan untuk menentukan awal segmen subdomain.
// Tanda tangan
// func (c *Ctx) Subdomains(offset ...int) []string
func TestSubDomains(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// Host: "tobi.ferrets.example.com"
		c.Subdomains()  // ["ferrets", "tobi"]
		c.Subdomains(1) // ["tobi"]
		return c.SendStatus(200)
	})
}

// Menyetel header HTTP Tipe Konten ke tipe MIME yang tercantum di sini ditentukan oleh ekstensi file .
// Tanda tangan
// func (c *Ctx) Type(ext string, charset ...string) *Ctx
func TestType(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// dia akan mengambil tipe yang terakhir
		// c.Type(".html")
		// c.Type("html")
		c.Type("png")
		// c.Type("json", "utf-8") // "application/json; charset=utf-8"
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)
	})
	app.Listen(":3000")
}

// Menambahkan kolom header tertentu ke header respons Vary . Ini akan menambahkan header, jika belum terdaftar, jika tidak, biarkan header terdaftar di lokasi saat ini.
// INFO
// Beberapa bidang diperbolehkan .
// Tanda tangan
// func (c *Ctx) Vary(fields ...string)

func TestVary(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Vary("Origin")     // => Vary: Origin
		c.Vary("User-Agent") // => Vary: Origin, User-Agent

		// No duplicates
		c.Vary("Origin") // => Vary: Origin, User-Agent

		c.Vary("Accept-Encoding", "Accept")
		// => Vary: Origin, User-Agent, Accept-Encoding, Accept
		fmt.Println(c.GetRespHeaders())
		return c.SendStatus(200)

	})
	app.Listen(":3000")
}

func TestWrite(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// write ini akan di jalankan tidak peduli dengan apa yang di return
		// c.Write([]byte("hello world"))
		// // naum jika kita mengirim string dia akan mengutamakan string, jika kita kitim
		// // send status dia akan mengutanakab string

		// // kita juga bisa menggunakakn writef untuk mELAKUKAN  format
		// hello := "farid"
		// c.Writef("hello %s ", hello)

		// selain itu dia juga bisa menjalankan write string untuk mengadopri string
		c.WriteString("hello user anang s")
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

// Properti Boolean, yaitu true, jika bidang header X-Requested-With permintaan adalah
// XMLHttpRequest , yang menunjukkan bahwa permintaan tersebut dikeluarkan
// oleh perpustakaan klien ( seperti jQuery ) .
// Tanda tangan
// func (c *Ctx) XHR() bool
func TestXHR(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		// X-Requested-With: XMLHttpRequest
		fmt.Println(c.XHR()) // true
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}

func TestXML(t *testing.T) {
	app := fiber.New()
	someStruct := struct {
		XMLName xml.Name `xml:"Golang loli"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
	}{
		Name: "farid anang s",
		Age:  17,
	}
	app.Get("/", func(c *fiber.Ctx) error {

		return c.XML(someStruct)
	})

	app.Listen(":3000")
}
