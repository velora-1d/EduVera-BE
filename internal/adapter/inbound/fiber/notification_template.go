package fiber_inbound_adapter

import (
	"context"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/gofiber/fiber/v2"
)

type notificationTemplateAdapter struct {
	dbPort outbound_port.NotificationTemplateDatabasePort
}

func NewNotificationTemplateAdapter(dbPort outbound_port.DatabasePort) *notificationTemplateAdapter {
	return &notificationTemplateAdapter{
		dbPort: dbPort.NotificationTemplate(),
	}
}

// GET /api/v1/owner/notification-templates
func (a *notificationTemplateAdapter) List(c *fiber.Ctx) error {
	ctx := context.Background()

	templates, err := a.dbPort.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   templates,
	})
}

// GET /api/v1/owner/notification-templates/:id
func (a *notificationTemplateAdapter) Get(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	template, err := a.dbPort.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	if template == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Template tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   template,
	})
}

// PUT /api/v1/owner/notification-templates/:id
func (a *notificationTemplateAdapter) Update(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		TemplateName    string `json:"template_name"`
		TemplateContent string `json:"template_content"`
		Variables       string `json:"variables"`
		IsActive        *bool  `json:"is_active"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	// Get existing
	template, err := a.dbPort.GetByID(ctx, id)
	if err != nil || template == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Template tidak ditemukan",
		})
	}

	// Update fields
	if input.TemplateName != "" {
		template.TemplateName = input.TemplateName
	}
	if input.TemplateContent != "" {
		template.TemplateContent = input.TemplateContent
	}
	if input.Variables != "" {
		template.Variables = input.Variables
	}
	if input.IsActive != nil {
		template.IsActive = *input.IsActive
	}

	if err := a.dbPort.Update(ctx, template); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Template berhasil diupdate",
		"data":    template,
	})
}

// POST /api/v1/owner/notification-templates
func (a *notificationTemplateAdapter) Create(c *fiber.Ctx) error {
	ctx := context.Background()

	var input model.NotificationTemplate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	if input.EventType == "" || input.Channel == "" || input.TemplateContent == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "event_type, channel, dan template_content wajib diisi",
		})
	}

	input.IsActive = true

	if err := a.dbPort.Save(ctx, &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Template berhasil dibuat",
		"data":    input,
	})
}

// DELETE /api/v1/owner/notification-templates/:id
func (a *notificationTemplateAdapter) Delete(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	if err := a.dbPort.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Template berhasil dihapus",
	})
}

// POST /api/v1/owner/notification-templates/:id/test
func (a *notificationTemplateAdapter) TestSend(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input struct {
		Phone     string            `json:"phone"`
		Variables map[string]string `json:"variables"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	template, err := a.dbPort.GetByID(ctx, id)
	if err != nil || template == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Template tidak ditemukan",
		})
	}

	// Replace variables in template
	message := template.TemplateContent
	for key, val := range input.Variables {
		message = replaceVariable(message, key, val)
	}

	// Return the parsed message for preview
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Preview pesan",
		"data": fiber.Map{
			"original": template.TemplateContent,
			"parsed":   message,
			"channel":  template.Channel,
		},
	})
}

func replaceVariable(content, key, value string) string {
	placeholder := "{{" + key + "}}"
	return stringReplace(content, placeholder, value)
}

func stringReplace(s, old, new string) string {
	result := s
	for {
		i := indexOfString(result, old)
		if i < 0 {
			break
		}
		result = result[:i] + new + result[i+len(old):]
	}
	return result
}

func indexOfString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
