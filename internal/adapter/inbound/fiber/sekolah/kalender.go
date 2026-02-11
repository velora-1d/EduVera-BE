package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (h *akademikHandler) GetKalenderEvents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetKalenderEvents(c.Context(), tenantID)
	if err != nil {
		if err != nil {
			return SendError(c, http.StatusInternalServerError, "Gagal mengambil data kalender", err)
		}
	}
	return SendSuccess(c, "Data kalender berhasil diambil", data)
}

func (h *akademikHandler) CreateKalenderEvent(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID, CreatedAt, UpdatedAt)
	var input struct {
		Title       string `json:"title"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&input); err != nil {
		if err := c.BodyParser(&input); err != nil {
			return SendError(c, http.StatusBadRequest, "Invalid request body", err)
		}
	}

	// Explicit mapping: DTO → DB Model
	m := model.KalenderEvent{
		TenantID:    tenantID, // From JWT, not user input
		Title:       input.Title,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Category:    input.Category,
		Description: input.Description,
	}

	if err := h.service.CreateKalenderEvent(c.Context(), tenantID, &m); err != nil {
		if err := h.service.CreateKalenderEvent(c.Context(), tenantID, &m); err != nil {
			return SendError(c, http.StatusInternalServerError, "Gagal membuat event", err)
		}
	}
	return SendCreated(c, "Event created", m)
}
