package model

import "time"

// ==========================================
// SDM MODELS - Employee & Payroll System
// ==========================================

// Employee represents a school employee (Guru/Staf)
type Employee struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	NIP          string    `json:"nip"` // Nomor Induk Pegawai
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`          // Guru, Staf, Kepala Sekolah
	EmployeeType string    `json:"employee_type"` // PNS, Honorer, Kontrak
	Department   string    `json:"department"`
	JoinDate     time.Time `json:"join_date"`
	BaseSalary   int64     `json:"base_salary"` // Gaji pokok
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmployeeInput for creating/updating employees
type EmployeeInput struct {
	TenantID     string `json:"tenant_id"`
	NIP          string `json:"nip"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	EmployeeType string `json:"employee_type"`
	Department   string `json:"department"`
	JoinDate     string `json:"join_date"`
	BaseSalary   int64  `json:"base_salary"`
}

// Payroll represents monthly payroll record
type Payroll struct {
	ID           string      `json:"id"`
	TenantID     string      `json:"tenant_id"`
	EmployeeID   string      `json:"employee_id"`
	EmployeeName string      `json:"employee_name,omitempty"` // Joined
	EmployeeNIP  string      `json:"employee_nip,omitempty"`  // Joined
	Period       string      `json:"period"`                  // "2026-01"
	BaseSalary   int64       `json:"base_salary"`
	Allowances   int64       `json:"allowances"` // Total tunjangan
	Deductions   int64       `json:"deductions"` // Total potongan
	NetSalary    int64       `json:"net_salary"` // Total terima
	Status       string      `json:"status"`     // draft, pending, paid
	PaidAt       *time.Time  `json:"paid_at,omitempty"`
	Details      []PayDetail `json:"details"` // JSONB breakdown
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// PayDetail represents individual pay component
type PayDetail struct {
	Name   string `json:"name"` // "Tunjangan Transport", "BPJS"
	Type   string `json:"type"` // allowance, deduction
	Amount int64  `json:"amount"`
}

// PayrollInput for generating payroll
type PayrollInput struct {
	TenantID   string      `json:"tenant_id"`
	EmployeeID string      `json:"employee_id"`
	Period     string      `json:"period"`
	BaseSalary int64       `json:"base_salary"`
	Details    []PayDetail `json:"details"`
}

// GeneratePayrollInput for batch payroll generation
type GeneratePayrollInput struct {
	TenantID string `json:"tenant_id"`
	Period   string `json:"period"` // "2026-01"
}

// PayrollConfig tenant-level payroll configuration
type PayrollConfig struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenant_id"`
	Components []PayComponent `json:"components"` // JSONB
}

// PayComponent defines a payroll component template
type PayComponent struct {
	Name       string `json:"name"`       // "Tunjangan Transport"
	Type       string `json:"type"`       // allowance, deduction
	Amount     int64  `json:"amount"`     // Fixed amount
	Percentage int    `json:"percentage"` // Or % of base salary (use one)
	AppliesTo  string `json:"applies_to"` // all, guru, staf
}

// Attendance represents employee attendance record
type Attendance struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	EmployeeID   string     `json:"employee_id"`
	EmployeeName string     `json:"employee_name,omitempty"` // Joined
	Date         time.Time  `json:"date"`
	CheckIn      *time.Time `json:"check_in,omitempty"`
	CheckOut     *time.Time `json:"check_out,omitempty"`
	Status       string     `json:"status"` // hadir, sakit, izin, alpha
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

// AttendanceInput for recording attendance
type AttendanceInput struct {
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`
	Date       string `json:"date"`
	Status     string `json:"status"`
	Notes      string `json:"notes"`
}

// PaySlip represents printable pay slip data
type PaySlip struct {
	Employee      Employee  `json:"employee"`
	Payroll       Payroll   `json:"payroll"`
	SchoolName    string    `json:"school_name"`
	SchoolAddress string    `json:"school_address"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// Employee Types
const (
	EmployeeTypePNS          = "PNS"
	EmployeeTypeHonorer      = "Honorer"
	EmployeeTypeKontrak      = "Kontrak"
	EmployeeTypeTetapYayasan = "Tetap Yayasan"
)

// Employee Roles
const (
	EmployeeRoleGuru        = "Guru"
	EmployeeRoleStaf        = "Staf"
	EmployeeRoleKepsek      = "Kepala Sekolah"
	EmployeeRoleWakilKepsek = "Wakil Kepala Sekolah"
	EmployeeRoleTU          = "Tata Usaha"
)

// Payroll Status
const (
	PayrollStatusDraft   = "draft"
	PayrollStatusPending = "pending"
	PayrollStatusPaid    = "paid"
)

// Pay Component Types
const (
	PayComponentAllowance = "allowance"
	PayComponentDeduction = "deduction"
)

// Attendance Status
const (
	AttendanceHadir = "hadir"
	AttendanceSakit = "sakit"
	AttendanceIzin  = "izin"
	AttendanceAlpha = "alpha"
)

// Default Payroll Components (template)
func DefaultPayrollComponents() []PayComponent {
	return []PayComponent{
		{Name: "Tunjangan Transport", Type: PayComponentAllowance, Amount: 500000, AppliesTo: "all"},
		{Name: "Tunjangan Makan", Type: PayComponentAllowance, Amount: 300000, AppliesTo: "all"},
		{Name: "Tunjangan Jabatan", Type: PayComponentAllowance, Percentage: 10, AppliesTo: "guru"},
		{Name: "BPJS Kesehatan", Type: PayComponentDeduction, Percentage: 1, AppliesTo: "all"},
		{Name: "BPJS Ketenagakerjaan", Type: PayComponentDeduction, Percentage: 2, AppliesTo: "all"},
	}
}
