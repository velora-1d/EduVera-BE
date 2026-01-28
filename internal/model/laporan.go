package model

// ReportRequest defines filters for the report
type ReportRequest struct {
	Type   string `query:"type"`   // kepesantrenan, tahfidz, diniyah, sdm, keuangan
	Period string `query:"period"` // current_month, current_semester, current_year
}

// ReportData is a generic map for flexible report data
type ReportData map[string]interface{}

// Specific Report Structs (Optional, for better typing if needed)
type KepesantrenanReport struct {
	Santri       string `json:"santri"`
	Asrama       string `json:"asrama"`
	Kedisiplinan string `json:"kedisiplinan"` // Baik, Cukup, Kurang
	Pelanggaran  int    `json:"pelanggaran"`
	Kehadiran    string `json:"kehadiran"` // Percentage
}
