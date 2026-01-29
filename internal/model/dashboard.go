package model

type DashboardStats struct {
	TotalSantri   int64 `json:"total_santri"`
	TotalAsrama   int64 `json:"total_asrama"`
	TotalUstadz   int64 `json:"total_ustadz"`
	TotalPengurus int64 `json:"total_pengurus"`

	// Status Kepesantrenan
	AttendanceRate   float64 `json:"attendance_rate"`   // Percentage
	ActiveViolations int64   `json:"active_violations"` // Pelanggaran yang belum selesai
	ActivePerizinan  int64   `json:"active_perizinan"`  // Perizinan yang sedang berjalan

	// Financial
	CashBalance  float64 `json:"cash_balance"`
	IncomeMonth  float64 `json:"income_month"`  // Pemasukan bulan ini
	ExpenseMonth float64 `json:"expense_month"` // Pengeluaran bulan ini
}

// SekolahDashboardStats holds statistics for Sekolah dashboard
type SekolahDashboardStats struct {
	TotalSiswa   int64   `json:"total_siswa"`
	TotalGuru    int64   `json:"total_guru"`
	TotalKelas   int64   `json:"total_kelas"`
	TagihanBulan float64 `json:"tagihan_bulan"` // Tagihan SPP bulan ini
	LunasCount   int64   `json:"lunas_count"`   // Jumlah yang sudah lunas
	BelumLunas   int64   `json:"belum_lunas"`   // Jumlah yang belum lunas
	TotalMapel   int64   `json:"total_mapel"`
	TotalJurusan int64   `json:"total_jurusan"`
}
