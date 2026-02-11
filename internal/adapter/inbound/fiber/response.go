package fiber_inbound_adapter

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// SendError logs the actual error and returns a sanitized JSON response to the client.
//
// Parameters:
//   - c: Fiber context
//   - status: HTTP status code
//   - clientMsg: Message to show to the user
//   - actualErr: Actual error object (optional, for logging)
//
// In production mode (APP_MODE=release), 500 errors will always return "Internal Server Error"
// regardless of clientMsg, to prevent leaking sensitive details.
func SendError(c *fiber.Ctx, status int, clientMsg string, actualErr error) error {
	// Log the actual error if present
	if actualErr != nil {
		logrus.WithFields(logrus.Fields{
			"path":      c.Path(),
			"method":    c.Method(),
			"status":    status,
			"client_ip": c.IP(),
			"error":     actualErr.Error(),
		}).Error("Handled error response")
	}

	// Sanitize message for 500 errors in production
	if status == fiber.StatusInternalServerError && os.Getenv("APP_MODE") == "release" {
		clientMsg = "Internal Server Error"
	}

	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": clientMsg,
		"error":   clientMsg, // Backward compatibility for frontends expecting 'error' key
	})
}

// SendSuccess returns a standardized success response
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}
