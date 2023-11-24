package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"gorm.io/gorm"
)

// There exist two distinct implementations of timeout middleware Fiber.
// New
// Wraps a fiber.Handler with a timeout. If the handler takes longer than the given duration to return, the timeout error is set and forwarded to the centralized ErrorHandler.
// CAUTION
// This has been deprecated since it raises race conditions.
// NewWithContext
// As a fiber.Handler wrapper, it creates a context with context.WithTimeout and pass it in UserContext.
// If the context passed executions (eg. DB ops, Http calls) takes longer than the given duration to return, the timeout error is set and forwarded to the centralized ErrorHandler.
// It does not cancel long running executions. Underlying executions must handle timeout by using context.Context parameter.

func TestTimeout(t *testing.T) {
	app := fiber.New()

	h := func(c *fiber.Ctx) error {
		sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
		if err := sleepWithContext(c.UserContext(), sleepTime); err != nil {
			return fmt.Errorf("%w: execution error", err)
		}

		return nil
	}

	// timeout.New is deprecated: This implementation contains data race issues. Use NewWithContext instead. Find documentation and sample usage on
	app.Get("/foo/:sleepTime", timeout.New(h, 5*time.Second))

	app.Listen(":3000")

}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return context.DeadlineExceeded
	case <-timer.C:
	}

	return nil
}

var errorFooTimeout = errors.New("foo context canceled")

// Use with custom error:
func TestTimeoutCustomError(t *testing.T) {
	app := fiber.New()

	h := func(c *fiber.Ctx) error {
		sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
		if err := sleepWithContext(c.UserContext(), sleepTime); err != nil {
			return fmt.Errorf("%w: execution error", err)
		}

		return nil
	}

	// timeout.New is deprecated: This implementation contains data race issues. Use NewWithContext instead. Find documentation and sample usage on
	app.Get("/foo/:sleepTime", timeout.NewWithContext(h, 5*time.Second, errorFooTimeout), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3000")

}

func sleepWithContextCutromError(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return errorFooTimeout
	case <-timer.C:
	}

	return nil
}

// Sample usage with a DB call:
func TestTimeoutDbCall(t *testing.T) {
	app := fiber.New()
	db, _ := gorm.Open(postgres.Open("postgres://localhost/foodb"), &gorm.Config{})

	handler := func(ctx *fiber.Ctx) error {
		tran := db.WithContext(ctx.UserContext()).Begin()

		if tran = tran.Exec("SELECT pg_sleep(50)"); tran.Error != nil {
			return tran.Error
		}

		if tran = tran.Commit(); tran.Error != nil {
			return tran.Error
		}

		return nil
	}

	app.Get("/foo", timeout.NewWithContext(handler, 10*time.Second))
	log.Fatal(app.Listen(":3000"))
}
