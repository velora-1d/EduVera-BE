package model

type DashboardStats struct {
	TotalSantri   int64 `json:"total_santri"`
	TotalAsrama   int64 `json:"total_asrama"`
	TotalUstadz   int64 `json:"total_ustadz"`
	TotalPengurus int64 `json:"total_pengurus"`

	// Status Kepesantrenan
	AttendanceRate   float64 `json:"attendance_rate"`   // Percentage
	ActiveViolations int64   `json:"active_violations"` // Pelanggaran yang belum selesai

	// Financial
	CashBalance  float64 `json:"cash_balance"`
	IncomeMonth  float64 `json:"income_month"`  // Pemasukan bulan ini
	ExpenseMonth float64 `json:"expense_month"` // Pengeluaran bulan ini
}
