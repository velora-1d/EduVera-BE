package fiber_inbound_adapter

import (
	"context"

	"prabogo/internal/domain"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

type eraporAdapter struct {
	domain domain.Domain
}

func NewERaporAdapter(domain domain.Domain) *eraporAdapter {
	return &eraporAdapter{domain: domain}
}

// GET /api/v1/sekolah/erapor/subjects
func (h *eraporAdapter) GetSubjects(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)

	subjects, err := h.domain.ERapor().GetSubjectsByTenant(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data mata pelajaran",
		})
	}

	return c.JSON(fiber.Map{
		"data": subjects,
	})
}

// POST /api/v1/sekolah/erapor/subjects
func (h *eraporAdapter) CreateSubject(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)

	var input model.SubjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	input.TenantID = tenantID

	subject, err := h.domain.ERapor().CreateSubject(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat mata pelajaran: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Mata pelajaran berhasil dibuat",
		"data":    subject,
	})
}

// PUT /api/v1/sekolah/erapor/subjects/:id
func (h *eraporAdapter) UpdateSubject(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var input model.SubjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	subject, err := h.domain.ERapor().UpdateSubject(ctx, id, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengupdate mata pelajaran",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mata pelajaran berhasil diupdate",
		"data":    subject,
	})
}

// DELETE /api/v1/sekolah/erapor/subjects/:id
func (h *eraporAdapter) DeleteSubject(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	err := h.domain.ERapor().DeleteSubject(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus mata pelajaran",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mata pelajaran berhasil dihapus",
	})
}

// POST /api/v1/sekolah/erapor/grades
func (h *eraporAdapter) SaveGrade(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)

	var input model.StudentGradeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	input.TenantID = tenantID

	grade, err := h.domain.ERapor().SaveGrade(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan nilai: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Nilai berhasil disimpan",
		"data":    grade,
	})
}

// POST /api/v1/sekolah/erapor/grades/batch
func (h *eraporAdapter) BatchSaveGrades(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)

	var input model.BatchGradeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	input.TenantID = tenantID

	grades, err := h.domain.ERapor().BatchSaveGrades(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan nilai batch: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Nilai batch berhasil disimpan",
		"count":   len(grades),
		"data":    grades,
	})
}

// GET /api/v1/sekolah/erapor/grades/student/:student_id
func (h *eraporAdapter) GetStudentGrades(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Params("student_id")
	semesterID := c.Query("semester", "")

	if semesterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parameter semester diperlukan",
		})
	}

	grades, err := h.domain.ERapor().GetGradesByStudent(ctx, studentID, semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil nilai siswa",
		})
	}

	return c.JSON(fiber.Map{
		"data": grades,
	})
}

// GET /api/v1/sekolah/erapor/grades/subject/:subject_id
func (h *eraporAdapter) GetSubjectGrades(c *fiber.Ctx) error {
	ctx := context.Background()
	subjectID := c.Params("subject_id")
	semesterID := c.Query("semester", "")

	if semesterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parameter semester diperlukan",
		})
	}

	grades, err := h.domain.ERapor().GetGradesBySubject(ctx, subjectID, semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil nilai mata pelajaran",
		})
	}

	return c.JSON(fiber.Map{
		"data": grades,
	})
}

// GET /api/v1/sekolah/erapor/rapor/:student_id/:semester
func (h *eraporAdapter) GetStudentRapor(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Params("student_id")
	semesterID := c.Params("semester")

	rapor, err := h.domain.ERapor().GetStudentRapor(ctx, studentID, semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data rapor",
		})
	}

	return c.JSON(fiber.Map{
		"data": rapor,
	})
}

// POST /api/v1/sekolah/erapor/generate
func (h *eraporAdapter) GenerateRapor(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)

	var input struct {
		StudentID   string `json:"student_id"`
		SemesterID  string `json:"semester_id"`
		CatatanWali string `json:"catatan_wali"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	rapor, err := h.domain.ERapor().GenerateRapor(ctx, tenantID, input.StudentID, input.SemesterID, input.CatatanWali)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal generate rapor: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Rapor berhasil digenerate",
		"data":    rapor,
	})
}

// GET /api/v1/sekolah/erapor/stats
func (h *eraporAdapter) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()
	tenantID := c.Locals("tenant_id").(string)
	semesterID := c.Query("semester", "")

	if semesterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parameter semester diperlukan",
		})
	}

	stats, err := h.domain.ERapor().GetGradeStats(ctx, tenantID, semesterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil statistik nilai",
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}

// GET /api/v1/sekolah/erapor/curriculum
func (h *eraporAdapter) GetCurriculum(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	profil, err := h.domain.Sekolah().GetProfil(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data profil",
		})
	}

	curriculum := "K13" // Default
	if profil != nil && profil.Curriculum != "" {
		curriculum = profil.Curriculum
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"curriculum": curriculum,
		},
	})
}

// PUT /api/v1/sekolah/erapor/curriculum
func (h *eraporAdapter) SetCurriculum(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var input struct {
		Curriculum string `json:"curriculum"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate curriculum value
	validCurriculums := map[string]bool{"K13": true, "MERDEKA": true, "PESANTREN": true}
	if !validCurriculums[input.Curriculum] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Curriculum harus K13, MERDEKA, atau PESANTREN",
		})
	}

	// Update profil
	update := &model.ProfilUpdate{
		Curriculum: input.Curriculum,
	}
	if err := h.domain.Sekolah().UpdateProfil(c.Context(), tenantID, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan setting curriculum",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Curriculum berhasil diupdate",
		"data": fiber.Map{
			"curriculum": input.Curriculum,
		},
	})
}

// GET /api/v1/sekolah/erapor/rapor/history
func (h *eraporAdapter) GetRaporHistory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	rapors, err := h.domain.Sekolah().GetRaporList(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil riwayat rapor",
		})
	}

	return c.JSON(fiber.Map{
		"data": rapors,
	})
}
