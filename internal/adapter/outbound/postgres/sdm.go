package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type sdmAdapter struct {
	db *goqu.Database
}

func NewSDMAdapter(sqlDB *sql.DB) *sdmAdapter {
	return &sdmAdapter{db: goqu.New("postgres", sqlDB)}
}

// ==========================================
// EMPLOYEE OPERATIONS
// ==========================================

func (a *sdmAdapter) CreateEmployee(ctx context.Context, input *model.EmployeeInput) (*model.Employee, error) {
	id := uuid.New().String()
	now := time.Now()

	joinDate, _ := time.Parse("2006-01-02", input.JoinDate)

	_, err := a.db.Insert("employees").Rows(
		goqu.Record{
			"id":            id,
			"tenant_id":     input.TenantID,
			"nip":           input.NIP,
			"name":          input.Name,
			"email":         input.Email,
			"phone":         input.Phone,
			"role":          input.Role,
			"employee_type": input.EmployeeType,
			"department":    input.Department,
			"join_date":     joinDate,
			"base_salary":   input.BaseSalary,
			"is_active":     true,
			"created_at":    now,
			"updated_at":    now,
		},
	).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Employee{
		ID:           id,
		TenantID:     input.TenantID,
		NIP:          input.NIP,
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		Role:         input.Role,
		EmployeeType: input.EmployeeType,
		Department:   input.Department,
		JoinDate:     joinDate,
		BaseSalary:   input.BaseSalary,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (a *sdmAdapter) UpdateEmployee(ctx context.Context, id string, input *model.EmployeeInput) (*model.Employee, error) {
	now := time.Now()
	joinDate, _ := time.Parse("2006-01-02", input.JoinDate)

	_, err := a.db.Update("employees").Set(
		goqu.Record{
			"nip":           input.NIP,
			"name":          input.Name,
			"email":         input.Email,
			"phone":         input.Phone,
			"role":          input.Role,
			"employee_type": input.EmployeeType,
			"department":    input.Department,
			"join_date":     joinDate,
			"base_salary":   input.BaseSalary,
			"updated_at":    now,
		},
	).Where(goqu.C("id").Eq(id)).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return a.GetEmployeeByID(ctx, id)
}

func (a *sdmAdapter) GetEmployeeByID(ctx context.Context, id string) (*model.Employee, error) {
	var emp model.Employee
	found, err := a.db.From("employees").Where(goqu.C("id").Eq(id)).ScanStructContext(ctx, &emp)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &emp, nil
}

func (a *sdmAdapter) GetEmployeesByTenant(ctx context.Context, tenantID string) ([]model.Employee, error) {
	var employees []model.Employee
	err := a.db.From("employees").
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Where(goqu.C("is_active").Eq(true)).
		Order(goqu.C("name").Asc()).
		ScanStructsContext(ctx, &employees)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (a *sdmAdapter) DeleteEmployee(ctx context.Context, id string) error {
	// Soft delete
	_, err := a.db.Update("employees").Set(
		goqu.Record{"is_active": false, "updated_at": time.Now()},
	).Where(goqu.C("id").Eq(id)).Executor().ExecContext(ctx)
	return err
}

// ==========================================
// PAYROLL OPERATIONS
// ==========================================

func (a *sdmAdapter) CreatePayroll(ctx context.Context, input *model.PayrollInput) (*model.Payroll, error) {
	id := uuid.New().String()
	now := time.Now()

	// Calculate totals
	var allowances, deductions int64
	for _, d := range input.Details {
		if d.Type == model.PayComponentAllowance {
			allowances += d.Amount
		} else {
			deductions += d.Amount
		}
	}
	netSalary := input.BaseSalary + allowances - deductions

	detailsJSON, _ := json.Marshal(input.Details)

	_, err := a.db.Insert("payrolls").Rows(
		goqu.Record{
			"id":          id,
			"tenant_id":   input.TenantID,
			"employee_id": input.EmployeeID,
			"period":      input.Period,
			"base_salary": input.BaseSalary,
			"allowances":  allowances,
			"deductions":  deductions,
			"net_salary":  netSalary,
			"details":     detailsJSON,
			"status":      model.PayrollStatusDraft,
			"created_at":  now,
			"updated_at":  now,
		},
	).OnConflict(
		goqu.DoUpdate("employee_id, period",
			goqu.Record{
				"base_salary": input.BaseSalary,
				"allowances":  allowances,
				"deductions":  deductions,
				"net_salary":  netSalary,
				"details":     detailsJSON,
				"updated_at":  now,
			}),
	).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Payroll{
		ID:         id,
		TenantID:   input.TenantID,
		EmployeeID: input.EmployeeID,
		Period:     input.Period,
		BaseSalary: input.BaseSalary,
		Allowances: allowances,
		Deductions: deductions,
		NetSalary:  netSalary,
		Details:    input.Details,
		Status:     model.PayrollStatusDraft,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (a *sdmAdapter) GetPayrollByID(ctx context.Context, id string) (*model.Payroll, error) {
	var row struct {
		ID         string     `db:"id"`
		TenantID   string     `db:"tenant_id"`
		EmployeeID string     `db:"employee_id"`
		Period     string     `db:"period"`
		BaseSalary int64      `db:"base_salary"`
		Allowances int64      `db:"allowances"`
		Deductions int64      `db:"deductions"`
		NetSalary  int64      `db:"net_salary"`
		Details    []byte     `db:"details"`
		Status     string     `db:"status"`
		PaidAt     *time.Time `db:"paid_at"`
		CreatedAt  time.Time  `db:"created_at"`
		UpdatedAt  time.Time  `db:"updated_at"`
	}

	found, err := a.db.From("payrolls").Where(goqu.C("id").Eq(id)).ScanStructContext(ctx, &row)
	if err != nil || !found {
		return nil, err
	}

	var details []model.PayDetail
	json.Unmarshal(row.Details, &details)

	return &model.Payroll{
		ID:         row.ID,
		TenantID:   row.TenantID,
		EmployeeID: row.EmployeeID,
		Period:     row.Period,
		BaseSalary: row.BaseSalary,
		Allowances: row.Allowances,
		Deductions: row.Deductions,
		NetSalary:  row.NetSalary,
		Details:    details,
		Status:     row.Status,
		PaidAt:     row.PaidAt,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func (a *sdmAdapter) GetPayrollByPeriod(ctx context.Context, tenantID, period string) ([]model.Payroll, error) {
	var rows []struct {
		ID           string     `db:"id"`
		TenantID     string     `db:"tenant_id"`
		EmployeeID   string     `db:"employee_id"`
		EmployeeName string     `db:"employee_name"`
		EmployeeNIP  string     `db:"employee_nip"`
		Period       string     `db:"period"`
		BaseSalary   int64      `db:"base_salary"`
		Allowances   int64      `db:"allowances"`
		Deductions   int64      `db:"deductions"`
		NetSalary    int64      `db:"net_salary"`
		Details      []byte     `db:"details"`
		Status       string     `db:"status"`
		PaidAt       *time.Time `db:"paid_at"`
		CreatedAt    time.Time  `db:"created_at"`
		UpdatedAt    time.Time  `db:"updated_at"`
	}

	err := a.db.From("payrolls").
		Select("payrolls.*", goqu.I("employees.name").As("employee_name"), goqu.I("employees.nip").As("employee_nip")).
		Join(goqu.T("employees"), goqu.On(goqu.I("payrolls.employee_id").Eq(goqu.I("employees.id")))).
		Where(goqu.I("payrolls.tenant_id").Eq(tenantID)).
		Where(goqu.I("payrolls.period").Eq(period)).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, err
	}

	payrolls := make([]model.Payroll, len(rows))
	for i, row := range rows {
		var details []model.PayDetail
		json.Unmarshal(row.Details, &details)
		payrolls[i] = model.Payroll{
			ID:           row.ID,
			TenantID:     row.TenantID,
			EmployeeID:   row.EmployeeID,
			EmployeeName: row.EmployeeName,
			EmployeeNIP:  row.EmployeeNIP,
			Period:       row.Period,
			BaseSalary:   row.BaseSalary,
			Allowances:   row.Allowances,
			Deductions:   row.Deductions,
			NetSalary:    row.NetSalary,
			Details:      details,
			Status:       row.Status,
			PaidAt:       row.PaidAt,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		}
	}
	return payrolls, nil
}

func (a *sdmAdapter) GetPayrollByEmployee(ctx context.Context, employeeID, period string) (*model.Payroll, error) {
	var row struct {
		ID         string     `db:"id"`
		TenantID   string     `db:"tenant_id"`
		EmployeeID string     `db:"employee_id"`
		Period     string     `db:"period"`
		BaseSalary int64      `db:"base_salary"`
		Allowances int64      `db:"allowances"`
		Deductions int64      `db:"deductions"`
		NetSalary  int64      `db:"net_salary"`
		Details    []byte     `db:"details"`
		Status     string     `db:"status"`
		PaidAt     *time.Time `db:"paid_at"`
	}

	found, err := a.db.From("payrolls").
		Where(goqu.C("employee_id").Eq(employeeID)).
		Where(goqu.C("period").Eq(period)).
		ScanStructContext(ctx, &row)
	if err != nil || !found {
		return nil, err
	}

	var details []model.PayDetail
	json.Unmarshal(row.Details, &details)

	return &model.Payroll{
		ID:         row.ID,
		TenantID:   row.TenantID,
		EmployeeID: row.EmployeeID,
		Period:     row.Period,
		BaseSalary: row.BaseSalary,
		Allowances: row.Allowances,
		Deductions: row.Deductions,
		NetSalary:  row.NetSalary,
		Details:    details,
		Status:     row.Status,
		PaidAt:     row.PaidAt,
	}, nil
}

func (a *sdmAdapter) UpdatePayrollStatus(ctx context.Context, id, status string) error {
	record := goqu.Record{"status": status, "updated_at": time.Now()}
	if status == model.PayrollStatusPaid {
		now := time.Now()
		record["paid_at"] = now
	}

	_, err := a.db.Update("payrolls").Set(record).
		Where(goqu.C("id").Eq(id)).
		Executor().ExecContext(ctx)
	return err
}

func (a *sdmAdapter) GetPayrollConfig(ctx context.Context, tenantID string) (*model.PayrollConfig, error) {
	var row struct {
		ID         string `db:"id"`
		TenantID   string `db:"tenant_id"`
		Components []byte `db:"components"`
	}

	found, err := a.db.From("payroll_configs").Where(goqu.C("tenant_id").Eq(tenantID)).ScanStructContext(ctx, &row)
	if err != nil || !found {
		// Return default config
		return &model.PayrollConfig{
			TenantID:   tenantID,
			Components: model.DefaultPayrollComponents(),
		}, nil
	}

	var components []model.PayComponent
	json.Unmarshal(row.Components, &components)

	return &model.PayrollConfig{
		ID:         row.ID,
		TenantID:   row.TenantID,
		Components: components,
	}, nil
}

func (a *sdmAdapter) SavePayrollConfig(ctx context.Context, config *model.PayrollConfig) error {
	componentsJSON, _ := json.Marshal(config.Components)

	if config.ID == "" {
		config.ID = uuid.New().String()
	}

	_, err := a.db.Insert("payroll_configs").Rows(
		goqu.Record{
			"id":         config.ID,
			"tenant_id":  config.TenantID,
			"components": componentsJSON,
		},
	).OnConflict(
		goqu.DoUpdate("tenant_id", goqu.Record{"components": componentsJSON}),
	).Executor().ExecContext(ctx)
	return err
}

// ==========================================
// ATTENDANCE OPERATIONS
// ==========================================

func (a *sdmAdapter) RecordAttendance(ctx context.Context, input *model.AttendanceInput) (*model.Attendance, error) {
	id := uuid.New().String()
	now := time.Now()
	date, _ := time.Parse("2006-01-02", input.Date)

	var checkIn, checkOut *time.Time
	if input.Status == model.AttendanceHadir {
		checkIn = &now
	}

	_, err := a.db.Insert("attendances").Rows(
		goqu.Record{
			"id":          id,
			"tenant_id":   input.TenantID,
			"employee_id": input.EmployeeID,
			"date":        date,
			"check_in":    checkIn,
			"check_out":   checkOut,
			"status":      input.Status,
			"notes":       input.Notes,
			"created_at":  now,
		},
	).OnConflict(
		goqu.DoUpdate("employee_id, date",
			goqu.Record{
				"status": input.Status,
				"notes":  input.Notes,
			}),
	).Executor().ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Attendance{
		ID:         id,
		TenantID:   input.TenantID,
		EmployeeID: input.EmployeeID,
		Date:       date,
		CheckIn:    checkIn,
		CheckOut:   checkOut,
		Status:     input.Status,
		Notes:      input.Notes,
		CreatedAt:  now,
	}, nil
}

func (a *sdmAdapter) GetAttendanceByDate(ctx context.Context, tenantID string, date string) ([]model.Attendance, error) {
	parsedDate, _ := time.Parse("2006-01-02", date)

	var rows []struct {
		ID           string     `db:"id"`
		TenantID     string     `db:"tenant_id"`
		EmployeeID   string     `db:"employee_id"`
		EmployeeName string     `db:"employee_name"`
		Date         time.Time  `db:"date"`
		CheckIn      *time.Time `db:"check_in"`
		CheckOut     *time.Time `db:"check_out"`
		Status       string     `db:"status"`
		Notes        string     `db:"notes"`
		CreatedAt    time.Time  `db:"created_at"`
	}

	err := a.db.From("attendances").
		Select("attendances.*", goqu.I("employees.name").As("employee_name")).
		Join(goqu.T("employees"), goqu.On(goqu.I("attendances.employee_id").Eq(goqu.I("employees.id")))).
		Where(goqu.I("attendances.tenant_id").Eq(tenantID)).
		Where(goqu.I("attendances.date").Eq(parsedDate)).
		ScanStructsContext(ctx, &rows)
	if err != nil {
		return nil, err
	}

	attendances := make([]model.Attendance, len(rows))
	for i, row := range rows {
		attendances[i] = model.Attendance{
			ID:           row.ID,
			TenantID:     row.TenantID,
			EmployeeID:   row.EmployeeID,
			EmployeeName: row.EmployeeName,
			Date:         row.Date,
			CheckIn:      row.CheckIn,
			CheckOut:     row.CheckOut,
			Status:       row.Status,
			Notes:        row.Notes,
			CreatedAt:    row.CreatedAt,
		}
	}
	return attendances, nil
}

func (a *sdmAdapter) GetAttendanceByEmployee(ctx context.Context, employeeID string, startDate, endDate string) ([]model.Attendance, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var attendances []model.Attendance
	err := a.db.From("attendances").
		Where(goqu.C("employee_id").Eq(employeeID)).
		Where(goqu.C("date").Gte(start)).
		Where(goqu.C("date").Lte(end)).
		Order(goqu.C("date").Desc()).
		ScanStructsContext(ctx, &attendances)
	return attendances, err
}

func (a *sdmAdapter) GetAttendanceSummary(ctx context.Context, tenantID, period string) (map[string]int, error) {
	// Parse period "2026-01" to get date range
	startDate, _ := time.Parse("2006-01", period)
	endDate := startDate.AddDate(0, 1, -1) // End of month

	type statusCount struct {
		Status string `db:"status"`
		Count  int    `db:"count"`
	}
	var counts []statusCount

	err := a.db.From("attendances").
		Select(goqu.C("status"), goqu.COUNT("*").As("count")).
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Where(goqu.C("date").Gte(startDate)).
		Where(goqu.C("date").Lte(endDate)).
		GroupBy("status").
		ScanStructsContext(ctx, &counts)
	if err != nil {
		return nil, err
	}

	summary := make(map[string]int)
	for _, c := range counts {
		summary[c.Status] = c.Count
	}
	return summary, nil
}
