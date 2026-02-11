package fiber_inbound_adapter

import (
	"fmt"

	"prabogo/internal/domain"
	"prabogo/internal/model"

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
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data pegawai", err)
	}

	return SendSuccess(c, "Data pegawai berhasil diambil", employees)
}

func (h *sdmAdapter) CreateEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.EmployeeInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid", err)
	}
	input.TenantID = tenantID

	employee, err := h.domain.SDM().CreateEmployee(ctx, &input)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat data pegawai", err)
	}

	return SendSuccess(c, "Pegawai berhasil ditambahkan", employee)
}

func (h *sdmAdapter) UpdateEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	employeeID := c.Params("id")

	var input model.EmployeeInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid", err)
	}
	input.TenantID = tenantID

	employee, err := h.domain.SDM().UpdateEmployee(ctx, employeeID, &input)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengupdate data pegawai", err)
	}

	return SendSuccess(c, "Pegawai berhasil diupdate", employee)
}

func (h *sdmAdapter) DeleteEmployee(c *fiber.Ctx) error {
	ctx := c.Context()
	employeeID := c.Params("id")

	if err := h.domain.SDM().DeleteEmployee(ctx, employeeID); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menghapus pegawai", err)
	}

	return SendSuccess(c, "Pegawai berhasil dihapus", nil)
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
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data gaji", err)
	}

	return SendSuccess(c, "Data gaji berhasil diambil", payrolls)
}

func (h *sdmAdapter) GeneratePayroll(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.GeneratePayrollInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid", err)
	}
	input.TenantID = tenantID

	payrolls, err := h.domain.SDM().GeneratePayroll(ctx, &input)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal generate gaji", err)
	}

	return SendSuccess(c, "Gaji berhasil digenerate", payrolls)
}

func (h *sdmAdapter) MarkPayrollPaid(c *fiber.Ctx) error {
	ctx := c.Context()
	payrollID := c.Params("id")

	if err := h.domain.SDM().MarkPayrollPaid(ctx, payrollID); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menandai gaji sebagai dibayar", err)
	}

	return SendSuccess(c, "Gaji berhasil ditandai sebagai dibayar", nil)
}

func (h *sdmAdapter) GetPaySlip(c *fiber.Ctx) error {
	ctx := c.Context()
	payrollID := c.Params("id")

	slip, err := h.domain.SDM().GetPaySlip(ctx, payrollID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil slip gaji", err)
	}

	if slip == nil {
		return SendError(c, fiber.StatusNotFound, "Slip gaji tidak ditemukan", nil)
	}

	return SendSuccess(c, "Slip gaji berhasil diambil", slip)
}

func (h *sdmAdapter) DownloadPaySlip(c *fiber.Ctx) error {
	ctx := c.Context()
	payrollID := c.Params("id")

	pdfBytes, err := h.domain.SDM().DownloadPaySlip(ctx, payrollID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal generate PDF slip gaji", err)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"slip-gaji-%s.pdf\"", payrollID))
	return c.Send(pdfBytes)
}

func (h *sdmAdapter) GetPayrollConfig(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	config, err := h.domain.SDM().GetPayrollConfig(ctx, tenantID)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil konfigurasi gaji", err)
	}

	return SendSuccess(c, "Konfigurasi gaji berhasil diambil", config)
}

func (h *sdmAdapter) SavePayrollConfig(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	// DTO: Only allow fillable fields (no ID)
	var input struct {
		Components []model.PayComponent `json:"components"`
	}
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid", err)
	}

	// Explicit mapping: DTO → DB Model
	config := model.PayrollConfig{
		TenantID:   tenantID, // From JWT, not user input
		Components: input.Components,
	}

	if err := h.domain.SDM().SavePayrollConfig(ctx, &config); err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan konfigurasi gaji", err)
	}

	return SendSuccess(c, "Konfigurasi gaji berhasil disimpan", nil)
}

// ==========================================
// ATTENDANCE HANDLERS
// ==========================================

func (h *sdmAdapter) GetAttendance(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	date := c.Query("date")

	if date == "" {
		return SendError(c, fiber.StatusBadRequest, "Parameter date diperlukan", nil)
	}

	attendances, err := h.domain.SDM().GetAttendanceByDate(ctx, tenantID, date)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data absensi", err)
	}

	return SendSuccess(c, "Data absensi berhasil diambil", attendances)
}

func (h *sdmAdapter) RecordAttendance(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	var input model.AttendanceInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid", err)
	}
	input.TenantID = tenantID

	attendance, err := h.domain.SDM().RecordAttendance(ctx, &input)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mencatat absensi", err)
	}

	return SendSuccess(c, "Absensi berhasil dicatat", attendance)
}

func (h *sdmAdapter) GetAttendanceSummary(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	period := c.Query("period", "2026-01")

	summary, err := h.domain.SDM().GetAttendanceSummary(ctx, tenantID, period)
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil rekapitulasi absensi", err)
	}

	return SendSuccess(c, "Rekapitulasi absensi berhasil diambil", summary)
}
