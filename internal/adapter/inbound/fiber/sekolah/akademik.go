package sekolah

import (
	"net/http"
	"prabogo/internal/domain/sekolah"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"

	"github.com/gofiber/fiber/v2"
)

type akademikHandler struct {
	service sekolah.AkademikDomain
}

func NewAkademikHandler(service sekolah.AkademikDomain) inbound_port.SekolahHttpPort {
	return &akademikHandler{
		service: service,
	}
}

// ------ Siswa Handler ------

func (h *akademikHandler) GetSiswaList(c *fiber.Ctx) error {
	// Get TenantID from context (set by middleware)
	tenantID := c.Locals("tenant_id").(string)

	siswaList, err := h.service.GetSiswaList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": siswaList,
	})
}

func (h *akademikHandler) CreateSiswa(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var siswa model.Siswa
	if err := c.BodyParser(&siswa); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	siswa.TenantID = tenantID

	if err := h.service.CreateSiswa(c.Context(), siswa); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Siswa created successfully"})
}

// ------ Guru Handler ------

func (h *akademikHandler) GetGuruList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	guruList, err := h.service.GetGuruList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": guruList})
}

func (h *akademikHandler) CreateGuru(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var guru model.Guru
	if err := c.BodyParser(&guru); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	guru.TenantID = tenantID

	if err := h.service.CreateGuru(c.Context(), guru); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Guru created successfully"})
}

// ------ Mapel Handler ------

func (h *akademikHandler) GetMapelList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	mapelList, err := h.service.GetMapelList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": mapelList})
}

// ------ Kelas Handler ------

func (h *akademikHandler) GetKelasList(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	kelasList, err := h.service.GetKelasList(c.Context(), tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": kelasList})
}

func (h *akademikHandler) CreateKelas(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var kelas model.Kelas
	if err := c.BodyParser(&kelas); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	kelas.TenantID = tenantID

	if err := h.service.CreateKelas(c.Context(), tenantID, &kelas); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Kelas created", "data": kelas})
}
