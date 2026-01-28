package sekolah

import (
	"net/http"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

// Extending the same handler struct or creating a new one?
// Since `akademikHandler` implements `SekolahHttpPort`, and we added methods to `SekolahHttpPort` (wait, did we?),
// we should probably check if `SekolahHttpPort` has these methods.
// Actually, `akademikHandler` is used for general sekolah stuff. Let's see if we should add to it.
// The `akademikHandler` struct is in `akademik.go`.
// Go allows methods to be in different files if they are in the same package.
// So we can use `func (h *akademikHandler) ...` here.

// But first I need to check if `SekolahHttpPort` has been updated. I forgot to update `SekolahHttpPort` in `inbound/sekolah.go`.
// Let me update the port first.
// I will write this content to a temporary buffer or just wait?
// I will update the port first.

// Wait, I can't cancel a tool call once thought process starts, but I can chain them.
// I will write the handler file assuming I will update the port.

func (h *akademikHandler) GetPelanggaranAturanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPelanggaranAturanList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *akademikHandler) CreatePelanggaranAturan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var m model.PelanggaranAturan
	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.service.CreatePelanggaranAturan(c.Context(), tenantID, &m); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Aturan created", "data": m})
}

func (h *akademikHandler) GetPelanggaranSiswaList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPelanggaranSiswaList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *akademikHandler) CreatePelanggaranSiswa(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var m model.PelanggaranSiswa
	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.service.CreatePelanggaranSiswa(c.Context(), tenantID, &m); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Pelanggaran recorded", "data": m})
}

func (h *akademikHandler) GetPerizinanList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	data, err := h.service.GetPerizinanList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (h *akademikHandler) CreatePerizinan(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var m model.Perizinan
	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.service.CreatePerizinan(c.Context(), tenantID, &m); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Perizinan created", "data": m})
}
