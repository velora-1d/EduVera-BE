package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain/student"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
)

type studentAdapter struct {
	domain student.StudentDomain
}

func NewStudentAdapter(domain student.StudentDomain) inbound_port.StudentHttpPort {
	return &studentAdapter{
		domain: domain,
	}
}

// GET /api/v1/students
func (h *studentAdapter) List(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	tenantID := c.Locals("tenant_id").(string)

	filter := model.StudentFilter{
		TenantID: tenantID,
		Type:     c.Query("type"), // siswa, santri, both
		Jenjang:  c.Query("jenjang"),
		KelasID:  c.Query("kelas_id"),
		KamarID:  c.Query("kamar_id"),
		Status:   c.Query("status"),
		Search:   c.Query("search"),
	}

	if c.Query("is_mukim") != "" {
		isMukim := c.Query("is_mukim") == "true"
		filter.IsMukim = &isMukim
	}

	students, err := h.domain.List(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data siswa/santri",
		})
	}

	return c.JSON(fiber.Map{
		"students": students,
		"count":    len(students),
	})
}

// GET /api/v1/students/:id
func (h *studentAdapter) Get(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID siswa/santri diperlukan",
		})
	}

	student, err := h.domain.FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Siswa/santri tidak ditemukan",
		})
	}

	return c.JSON(student)
}

// POST /api/v1/students
func (h *studentAdapter) Create(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	tenantID := c.Locals("tenant_id").(string)

	var input model.StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid",
		})
	}

	// Validate required fields
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama siswa/santri wajib diisi",
		})
	}

	if input.Type == "" {
		input.Type = model.StudentTypeSiswa // Default to siswa
	}

	// Validate type
	if input.Type != model.StudentTypeSiswa &&
		input.Type != model.StudentTypeSantri &&
		input.Type != model.StudentTypeBoth {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tipe harus: siswa, santri, atau both",
		})
	}

	student, err := h.domain.Create(ctx, tenantID, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menambahkan siswa/santri: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Siswa/santri berhasil ditambahkan",
		"student": student,
	})
}

// PUT /api/v1/students/:id
func (h *studentAdapter) Update(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID siswa/santri diperlukan",
		})
	}

	var input model.StudentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data tidak valid",
		})
	}

	student, err := h.domain.Update(ctx, id, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui data siswa/santri: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data siswa/santri berhasil diperbarui",
		"student": student,
	})
}

// DELETE /api/v1/students/:id
func (h *studentAdapter) Delete(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID siswa/santri diperlukan",
		})
	}

	err := h.domain.Delete(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus siswa/santri",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Siswa/santri berhasil dihapus",
	})
}

// GET /api/v1/students/count
func (h *studentAdapter) Count(a any) error {
	c := a.(*fiber.Ctx)
	ctx := context.Background()

	tenantID := c.Locals("tenant_id").(string)

	count, err := h.domain.CountByTenant(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghitung data",
		})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}
