package logger

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func TestLogger(t *testing.T) {
	/////// Note: The method of calling the Fatal level will interrupt the program running after printing the log, please use it with caution. Directly print logs of different levels, which will be entered into messageKey, the default is msg.
	// log.Info("Hello, World!")
	// log.Debug("Are you OK?")
	// log.Info("42 is the answer to life, the universe, and everything")
	// log.Warn("We are under attack!")
	// log.Error("Houston, we have a problem.")
	// log.Fatal("So Long, and Thanks for All the Fislog.")
	// log.Panic("The system is down.")

	///// Format and print logs of different levels, all methods end with f
	// log.Debugf("error: at Testlogerr, message:%s", "erororko tololll")
	// log.Infof("%d is the answer to life, the universe, and everything", 233)
	// log.Warnf("We are under attack %s!", "boss")
	// log.Errorf("%s, we have a problem.", "Master Shifu")
	// log.Fatalf("So Long, and Thanks for All the %s.", "banana")

	/////// Print a message with the key and value, or KEYVALS UNPAIRED if the key and value are not a pair.
	log.Debugw("", "error", "boy", "message", "ini terjadi error di testlogger")
	log.Infow("", "number", 233)
	log.Warnw("", "job", "boss")
	log.Errorw("", "name", "Master Shifu")
	log.Fatalw("", "fruit", "banana")

	//// set level
	// log.SetLevel sets the level of logs below which logs will not be output. The default logger is LevelTrace.
	// Note that this method is not concurrent-safe.
	log.SetLevel(log.LevelPanic)

	// // set output
	// 	var logger log.AllLogger = &defaultLogger{
	// 		stdlog: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
	// 		depth:  4,
	// 	}
	// 	Set the output destination to the file.
	// // Output to ./test.log file
	// f, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	//     return
	// }
	// log.SetOutput(f)

	////// Bind context

}

func loggerMiddleware(ctx *fiber.Ctx) error {
	commonLogger := log.WithContext(ctx.Context())
	ctx.Locals("logger", commonLogger)
	return ctx.Next()
}

func TestLoggerWithContext(t *testing.T) {
	app := fiber.New()
	app.Use(loggerMiddleware)
	app.Get("/", func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(log.Logger)
		logger.Info("hello")
		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
