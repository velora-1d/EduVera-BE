package sekolah

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// SendError sends a standardized error response
func SendError(c *fiber.Ctx, status int, clientMsg string, actualErr error) error {
	if actualErr != nil {
		logrus.WithFields(logrus.Fields{
			"path":      c.Path(),
			"method":    c.Method(),
			"status":    status,
			"client_ip": c.IP(),
			"error":     actualErr.Error(),
			"pkg":       "sekolah",
		}).Error("Handled error response")
	}

	if status == fiber.StatusInternalServerError && os.Getenv("APP_MODE") == "release" {
		clientMsg = "Internal Server Error"
	}

	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": clientMsg,
		"error":   clientMsg, // Backward compatibility
	})
}

// SendSuccess sends a standardized success response
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// SendCreated sends a standardized created response
func SendCreated(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}
