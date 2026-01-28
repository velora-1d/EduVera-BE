package fiber_inbound_adapter

import (
	"eduvera/internal/domain"
	"eduvera/internal/model"

	"github.com/gofiber/fiber/v2"
)

type sdmAdapter struct {
	domain domain.Domain
}

func NewSDMAdapter(d domain.Domain) *sdmAdapter {
	return &sdmAdapter{domain: d}
}

// ==========================================
// EMPLOYEE HANDLERS
// ==========================================

func (h *sdmAdapter) GetEmployees(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	employees, err := h.domain.SDM().GetEmployees(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data pegawai",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   employees,
	})
}

func (h *sdmAdapter) CreateEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.EmployeeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak valid",
		})
	}
	input.TenantID = tenantID

	employee, err := h.domain.SDM().CreateEmployee(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat data pegawai",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Pegawai berhasil ditambahkan",
		"data":    employee,
	})
}

func (h *sdmAdapter) UpdateEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	employeeID := c.Params("id")

	var input model.EmployeeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak valid",
		})
	}
	input.TenantID = tenantID

	employee, err := h.domain.SDM().UpdateEmployee(ctx, employeeID, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengupdate data pegawai",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pegawai berhasil diupdate",
		"data":    employee,
	})
}

func (h *sdmAdapter) DeleteEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	employeeID := c.Params("id")

	if err := h.domain.SDM().DeleteEmployee(ctx, employeeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menghapus pegawai",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pegawai berhasil dihapus",
	})
}

// ==========================================
// PAYROLL HANDLERS
// ==========================================

func (h *sdmAdapter) GetPayrollByPeriod(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	period := c.Query("period", "2026-01") // Default current month

	payrolls, err := h.domain.SDM().GetPayrollByPeriod(ctx, tenantID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data gaji",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   payrolls,
	})
}

func (h *sdmAdapter) GeneratePayroll(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.GeneratePayrollInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak valid",
		})
	}
	input.TenantID = tenantID

	payrolls, err := h.domain.SDM().GeneratePayroll(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal generate gaji",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Gaji berhasil digenerate",
		"data":    payrolls,
	})
}

func (h *sdmAdapter) MarkPayrollPaid(c *fiber.Ctx) error {
	ctx := c.Context()
	payrollID := c.Params("id")

	if err := h.domain.SDM().MarkPayrollPaid(ctx, payrollID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menandai gaji sebagai dibayar",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Gaji berhasil ditandai sebagai dibayar",
	})
}

func (h *sdmAdapter) GetPaySlip(c *fiber.Ctx) error {
	ctx := c.Context()
	payrollID := c.Params("id")

	slip, err := h.domain.SDM().GetPaySlip(ctx, payrollID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil slip gaji",
		})
	}

	if slip == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Slip gaji tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   slip,
	})
}

func (h *sdmAdapter) GetPayrollConfig(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	config, err := h.domain.SDM().GetPayrollConfig(ctx, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil konfigurasi gaji",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   config,
	})
}

func (h *sdmAdapter) SavePayrollConfig(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var config model.PayrollConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak valid",
		})
	}
	config.TenantID = tenantID

	if err := h.domain.SDM().SavePayrollConfig(ctx, &config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan konfigurasi gaji",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Konfigurasi gaji berhasil disimpan",
	})
}

// ==========================================
// ATTENDANCE HANDLERS
// ==========================================

func (h *sdmAdapter) GetAttendance(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	date := c.Query("date")

	if date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Parameter date diperlukan",
		})
	}

	attendances, err := h.domain.SDM().GetAttendanceByDate(ctx, tenantID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data absensi",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   attendances,
	})
}

func (h *sdmAdapter) RecordAttendance(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.AttendanceInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak valid",
		})
	}
	input.TenantID = tenantID

	attendance, err := h.domain.SDM().RecordAttendance(ctx, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mencatat absensi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Absensi berhasil dicatat",
		"data":    attendance,
	})
}

func (h *sdmAdapter) GetAttendanceSummary(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	period := c.Query("period", "2026-01")

	summary, err := h.domain.SDM().GetAttendanceSummary(ctx, tenantID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil rekapitulasi absensi",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   summary,
	})
}
