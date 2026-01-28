package postgres_outbound_adapter

import (
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

// ------ Laporan ------

func (a *sekolahAdapter) GetReportData(tenantID string, req model.ReportRequest) ([]model.ReportData, error) {
	// Determines which aggregation to run based on req.Type
	switch req.Type {
	case "kepesantrenan":
		return a.getKepesantrenanReport(tenantID)
	case "tahfidz":
		return a.getTahfidzReport(tenantID)
	case "diniyah":
		return a.getDiniyahReport(tenantID)
	case "sdm":
		return a.getSDMReport(tenantID) // Placeholder using SDM modules
	case "keuangan":
		return a.getKeuanganReport(tenantID)
	default:
		return []model.ReportData{}, nil
	}
}

func (a *sekolahAdapter) getKepesantrenanReport(tenantID string) ([]model.ReportData, error) {
	// Query: Join Siswa with Points/Exemptions (Mock logic for now as detailed tables might need complex joins)
	// Returning simplified mock data from DB context or basic fetch
	// Real implementation: SELECT s.nama, a.nama as asrama, count(p.id) as pelanggaran FROM sekolah_siswa s ...

	// For this Audit Phase, we'll return a sample aggregation based on real students if possible, or mock structure
	// Let's allow returning a mix for demonstration.

	// Simple Real Query: Get Students
	var list []model.ReportData

	ds := goqu.Dialect("postgres").From("sekolah_siswa").
		Select("nama_lengkap", "nis").
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Limit(10)

	query, _, _ := ds.ToSQL()
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nama, nis string
		rows.Scan(&nama, &nis)

		list = append(list, model.ReportData{
			"col1": nama,
			"col2": "Asrama 1", // Placeholder
			"col3": "Baik",     // Logic for Kedisiplinan
			"col4": "0",        // Logic for Pelanggaran Count
			"col5": "95%",      // Logic for Kehadiran
		})
	}

	return list, nil
}

func (a *sekolahAdapter) getTahfidzReport(tenantID string) ([]model.ReportData, error) {
	// Similar logic, fetching students
	var list []model.ReportData
	query := "SELECT nama_lengkap FROM sekolah_siswa WHERE tenant_id = $1 LIMIT 10"
	rows, err := a.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var nama string
		rows.Scan(&nama)
		list = append(list, model.ReportData{
			"col1": nama,
			"col2": "Juz 30",
			"col3": "Lancar",
			"col4": "Tercapai",
			"col5": "A",
		})
	}
	return list, nil
}

func (a *sekolahAdapter) getDiniyahReport(tenantID string) ([]model.ReportData, error) {
	var list []model.ReportData
	query := "SELECT nama_lengkap FROM sekolah_siswa WHERE tenant_id = $1 LIMIT 10"
	rows, err := a.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var nama string
		rows.Scan(&nama)
		list = append(list, model.ReportData{
			"col1": nama,
			"col2": "Ula A",
			"col3": "85",
			"col4": "Baik",
			"col5": "90%",
		})
	}
	return list, nil
}

func (a *sekolahAdapter) getSDMReport(tenantID string) ([]model.ReportData, error) {
	var list []model.ReportData
	// Fetch Employees
	query := "SELECT name, role FROM employees WHERE tenant_id = $1 LIMIT 10"
	rows, err := a.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name, role string
		rows.Scan(&name, &role)
		list = append(list, model.ReportData{
			"col1": name,
			"col2": role,
			"col3": "Baik",
			"col4": "Aktif",
			"col5": "100%",
		})
	}
	return list, nil
}

func (a *sekolahAdapter) getKeuanganReport(tenantID string) ([]model.ReportData, error) {
	// Mock Financial Summary (tenantID filtering to be added for real implementation)
	_ = tenantID // Placeholder usage until real query implemented
	return []model.ReportData{
		{"col1": "Pemasukan SPP", "col2": "Bulan Ini", "col3": "Rp 50.000.000", "col4": "Rp 45.000.000", "col5": "90%"},
		{"col1": "Pemasukan Tabungan", "col2": "Bulan Ini", "col3": "Rp 10.000.000", "col4": "Rp 12.000.000", "col5": "120%"},
		{"col1": "Pengeluaran Gaji", "col2": "Bulan Ini", "col3": "Rp 30.000.000", "col4": "Rp 30.000.000", "col5": "100%"},
	}, nil
}
