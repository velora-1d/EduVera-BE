package outbound_port

import (
	"context"

	"prabogo/internal/model"
)

// SDMDatabasePort defines the interface for SDM database operations
type SDMDatabasePort interface {
	// Employee operations
	CreateEmployee(ctx context.Context, input *model.EmployeeInput) (*model.Employee, error)
	UpdateEmployee(ctx context.Context, id string, input *model.EmployeeInput) (*model.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*model.Employee, error)
	GetEmployeesByTenant(ctx context.Context, tenantID string) ([]model.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error

	// Payroll operations
	CreatePayroll(ctx context.Context, input *model.PayrollInput) (*model.Payroll, error)
	GetPayrollByID(ctx context.Context, id string) (*model.Payroll, error)
	GetPayrollByPeriod(ctx context.Context, tenantID, period string) ([]model.Payroll, error)
	GetPayrollByEmployee(ctx context.Context, employeeID, period string) (*model.Payroll, error)
	UpdatePayrollStatus(ctx context.Context, id, status string) error
	GetPayrollConfig(ctx context.Context, tenantID string) (*model.PayrollConfig, error)
	SavePayrollConfig(ctx context.Context, config *model.PayrollConfig) error

	// Attendance operations
	RecordAttendance(ctx context.Context, input *model.AttendanceInput) (*model.Attendance, error)
	GetAttendanceByDate(ctx context.Context, tenantID string, date string) ([]model.Attendance, error)
	GetAttendanceByEmployee(ctx context.Context, employeeID string, startDate, endDate string) ([]model.Attendance, error)
	GetAttendanceSummary(ctx context.Context, tenantID, period string) (map[string]int, error)
}
