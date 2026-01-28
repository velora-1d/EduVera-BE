package sdm

import (
	"context"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

// SDMDomain interface
type SDMDomain interface {
	// Employee operations
	GetEmployees(ctx context.Context, tenantID string) ([]model.Employee, error)
	CreateEmployee(ctx context.Context, input *model.EmployeeInput) (*model.Employee, error)
	UpdateEmployee(ctx context.Context, id string, input *model.EmployeeInput) (*model.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error

	// Payroll operations
	GetPayrollByPeriod(ctx context.Context, tenantID, period string) ([]model.Payroll, error)
	GeneratePayroll(ctx context.Context, input *model.GeneratePayrollInput) ([]model.Payroll, error)
	MarkPayrollPaid(ctx context.Context, id string) error
	GetPaySlip(ctx context.Context, payrollID string) (*model.PaySlip, error)
	GetPayrollConfig(ctx context.Context, tenantID string) (*model.PayrollConfig, error)
	SavePayrollConfig(ctx context.Context, config *model.PayrollConfig) error

	// Attendance operations
	GetAttendanceByDate(ctx context.Context, tenantID, date string) ([]model.Attendance, error)
	RecordAttendance(ctx context.Context, input *model.AttendanceInput) (*model.Attendance, error)
	GetAttendanceSummary(ctx context.Context, tenantID, period string) (map[string]int, error)
}

type sdmDomain struct {
	db outbound_port.SDMDatabasePort
}

func NewSDMDomain(db outbound_port.SDMDatabasePort) SDMDomain {
	return &sdmDomain{db: db}
}

// ==========================================
// EMPLOYEE OPERATIONS
// ==========================================

func (d *sdmDomain) GetEmployees(ctx context.Context, tenantID string) ([]model.Employee, error) {
	return d.db.GetEmployeesByTenant(ctx, tenantID)
}

func (d *sdmDomain) CreateEmployee(ctx context.Context, input *model.EmployeeInput) (*model.Employee, error) {
	return d.db.CreateEmployee(ctx, input)
}

func (d *sdmDomain) UpdateEmployee(ctx context.Context, id string, input *model.EmployeeInput) (*model.Employee, error) {
	return d.db.UpdateEmployee(ctx, id, input)
}

func (d *sdmDomain) DeleteEmployee(ctx context.Context, id string) error {
	return d.db.DeleteEmployee(ctx, id)
}

// ==========================================
// PAYROLL OPERATIONS
// ==========================================

func (d *sdmDomain) GetPayrollByPeriod(ctx context.Context, tenantID, period string) ([]model.Payroll, error) {
	return d.db.GetPayrollByPeriod(ctx, tenantID, period)
}

// GeneratePayroll creates payroll records for all employees in a period
func (d *sdmDomain) GeneratePayroll(ctx context.Context, input *model.GeneratePayrollInput) ([]model.Payroll, error) {
	// Get all active employees
	employees, err := d.db.GetEmployeesByTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	// Get payroll config for this tenant
	config, err := d.db.GetPayrollConfig(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	var payrolls []model.Payroll

	for _, emp := range employees {
		// Build pay details based on config
		details := d.calculatePayDetails(emp, config.Components)

		payrollInput := &model.PayrollInput{
			TenantID:   input.TenantID,
			EmployeeID: emp.ID,
			Period:     input.Period,
			BaseSalary: emp.BaseSalary,
			Details:    details,
		}

		payroll, err := d.db.CreatePayroll(ctx, payrollInput)
		if err != nil {
			return nil, err
		}

		payroll.EmployeeName = emp.Name
		payroll.EmployeeNIP = emp.NIP
		payrolls = append(payrolls, *payroll)
	}

	return payrolls, nil
}

// calculatePayDetails applies payroll components to employee
func (d *sdmDomain) calculatePayDetails(emp model.Employee, components []model.PayComponent) []model.PayDetail {
	var details []model.PayDetail

	for _, comp := range components {
		// Check if component applies to this employee
		if !d.componentApplies(emp, comp) {
			continue
		}

		var amount int64
		if comp.Amount > 0 {
			amount = comp.Amount
		} else if comp.Percentage > 0 {
			amount = emp.BaseSalary * int64(comp.Percentage) / 100
		}

		if amount > 0 {
			details = append(details, model.PayDetail{
				Name:   comp.Name,
				Type:   comp.Type,
				Amount: amount,
			})
		}
	}

	return details
}

// componentApplies checks if a pay component applies to an employee
func (d *sdmDomain) componentApplies(emp model.Employee, comp model.PayComponent) bool {
	switch comp.AppliesTo {
	case "all":
		return true
	case "guru":
		return emp.Role == model.EmployeeRoleGuru
	case "staf":
		return emp.Role == model.EmployeeRoleStaf || emp.Role == model.EmployeeRoleTU
	default:
		return true
	}
}

func (d *sdmDomain) MarkPayrollPaid(ctx context.Context, id string) error {
	return d.db.UpdatePayrollStatus(ctx, id, model.PayrollStatusPaid)
}

func (d *sdmDomain) GetPaySlip(ctx context.Context, payrollID string) (*model.PaySlip, error) {
	payroll, err := d.db.GetPayrollByID(ctx, payrollID)
	if err != nil || payroll == nil {
		return nil, err
	}

	employee, err := d.db.GetEmployeeByID(ctx, payroll.EmployeeID)
	if err != nil || employee == nil {
		return nil, err
	}

	// In production, get school name from tenant
	return &model.PaySlip{
		Employee:   *employee,
		Payroll:    *payroll,
		SchoolName: "SMK Negeri 1 Example", // TODO: Get from tenant
	}, nil
}

func (d *sdmDomain) GetPayrollConfig(ctx context.Context, tenantID string) (*model.PayrollConfig, error) {
	return d.db.GetPayrollConfig(ctx, tenantID)
}

func (d *sdmDomain) SavePayrollConfig(ctx context.Context, config *model.PayrollConfig) error {
	return d.db.SavePayrollConfig(ctx, config)
}

// ==========================================
// ATTENDANCE OPERATIONS
// ==========================================

func (d *sdmDomain) GetAttendanceByDate(ctx context.Context, tenantID, date string) ([]model.Attendance, error) {
	return d.db.GetAttendanceByDate(ctx, tenantID, date)
}

func (d *sdmDomain) RecordAttendance(ctx context.Context, input *model.AttendanceInput) (*model.Attendance, error) {
	return d.db.RecordAttendance(ctx, input)
}

func (d *sdmDomain) GetAttendanceSummary(ctx context.Context, tenantID, period string) (map[string]int, error) {
	return d.db.GetAttendanceSummary(ctx, tenantID, period)
}
