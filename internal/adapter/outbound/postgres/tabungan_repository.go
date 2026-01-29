package postgres_outbound_adapter

import (
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

var tableTabungan = goqu.T("sekolah_tabungan")
var tableTabunganMutasi = goqu.T("sekolah_tabungan_mutasi")

func (a *sekolahAdapter) GetTabunganList(tenantID string) ([]model.Tabungan, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableTabungan).
		Select(
			tableTabungan.Col("id"),
			tableTabungan.Col("tenant_id"),
			tableTabungan.Col("santri_id"),
			tableTabungan.Col("saldo"),
			tableTabungan.Col("status"),
			tableTabungan.Col("created_at"),
			tableTabungan.Col("updated_at"),
			goqu.COALESCE(tableSiswa.Col("nama_lengkap"), "").As("nama_santri"),
			goqu.COALESCE(tableSiswa.Col("nis"), "").As("nis"),
		).
		LeftJoin(tableSiswa, goqu.On(tableTabungan.Col("santri_id").Eq(tableSiswa.Col("id")))).
		Where(tableTabungan.Col("tenant_id").Eq(tenantID)).
		Order(tableTabungan.Col("updated_at").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Tabungan
	for rows.Next() {
		var m model.Tabungan
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.SantriID, &m.Saldo, &m.Status,
			&m.CreatedAt, &m.UpdatedAt, &m.NamaSantri, &m.NIS,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateTabunganMutasi(m *model.TabunganMutasi) error {
	dialect := goqu.Dialect("postgres")

	// Start Transaction (Logic handled in domain or repo, for simplicity doing repo side logic trigger here or just insert)
	// Ideally we update saldo here too.

	// Note: a.db interface might not support Begin() directly depending on implementation.
	// For MVP/Audit, executing sequentially.

	// 1. Insert Mutasi
	ds := dialect.Insert(tableTabunganMutasi).Rows(goqu.Record{
		"tabungan_id": m.TabunganID,
		"tenant_id":   m.TenantID,
		"tipe":        m.Tipe,
		"nominal":     m.Nominal,
		"keterangan":  m.Keterangan,
		"petugas":     m.Petugas,
	}).Returning("id", "created_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	if err := a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt); err != nil {
		return err
	}

	// 2. Update Saldo in Parent Table
	var updateOp interface{}
	if m.Tipe == "Debit" { // Masuk (Menambah Saldo Tabungan)
		updateOp = goqu.L("saldo + ?", m.Nominal)
	} else { // Kredit (Mengurangi Saldo Tabungan)
		updateOp = goqu.L("saldo - ?", m.Nominal)
	}

	dsUpdate := dialect.Update(tableTabungan).
		Set(goqu.Record{"saldo": updateOp}).
		Where(tableTabungan.Col("id").Eq(m.TabunganID))

	queryUpdate, _, err := dsUpdate.ToSQL()
	if err != nil {
		return err
	}
	if _, err := a.db.Exec(queryUpdate); err != nil {
		return err
	}

	return nil
}

// ------ Kalender ------

func (a *sekolahAdapter) GetKalenderEvents(tenantID string) ([]model.KalenderEvent, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("sekolah_kalender").
		Select(goqu.Star()).
		Where(goqu.C("tenant_id").Eq(tenantID)).
		Order(goqu.C("start_date").Asc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.KalenderEvent
	for rows.Next() {
		var m model.KalenderEvent
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.Title, &m.StartDate, &m.EndDate, &m.Category, &m.Description,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateKalenderEvent(m *model.KalenderEvent) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert("sekolah_kalender").Rows(goqu.Record{
		"tenant_id":   m.TenantID,
		"title":       m.Title,
		"start_date":  m.StartDate,
		"end_date":    m.EndDate,
		"category":    m.Category,
		"description": m.Description,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

// ------ Profil ------

func (a *sekolahAdapter) GetProfil(tenantID string) (*model.Profil, error) {
	// Delegating to safe implementation for readability
	return a.getProfilSafe(tenantID)
}

func (a *sekolahAdapter) getProfilSafe(tenantID string) (*model.Profil, error) {
	// 1. Get Tenant
	var tenantName, tenantAddress string
	q1, _, _ := goqu.Dialect("postgres").From("tenants").Select("name", "address").Where(goqu.C("id").Eq(tenantID)).ToSQL()
	if err := a.db.QueryRow(q1).Scan(&tenantName, &tenantAddress); err != nil {
		return nil, err
	}

	// 2. Get Profil
	var m model.Profil
	m.TenantID = tenantID
	m.NamaPesantren = tenantName
	m.Alamat = tenantAddress

	q2, _, _ := goqu.Dialect("postgres").From("sekolah_profil").Select(goqu.Star()).Where(goqu.C("tenant_id").Eq(tenantID)).ToSQL()
	err := a.db.QueryRow(q2).Scan(
		&m.ID, &m.TenantID, &m.JenisPesantren, &m.Deskripsi, &m.Website,
		&m.EmailKontak, &m.NoTelpKontak, &m.LogoURL, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		// If no rows, return basic info
		return &m, nil
	}
	return &m, nil
}

func (a *sekolahAdapter) UpdateProfil(tenantID string, m *model.ProfilUpdate) error {
	// Update Tenants Table (Name, Address) only if provided
	if m.NamaPesantren != "" || m.Alamat != "" {
		record := goqu.Record{}
		if m.NamaPesantren != "" {
			record["name"] = m.NamaPesantren
		}
		if m.Alamat != "" {
			record["address"] = m.Alamat
		}
		q1, _, _ := goqu.Dialect("postgres").Update("tenants").Set(record).Where(goqu.C("id").Eq(tenantID)).ToSQL()
		if _, err := a.db.Exec(q1); err != nil {
			return err
		}
	}

	// Update or Insert Sekolah Profil
	// Check if exists
	var exists bool
	a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM sekolah_profil WHERE tenant_id = $1)", tenantID).Scan(&exists)

	// Build record with non-empty fields only
	profilRecord := goqu.Record{}
	if m.JenisPesantren != "" {
		profilRecord["jenis_pesantren"] = m.JenisPesantren
	}
	if m.Deskripsi != "" {
		profilRecord["deskripsi"] = m.Deskripsi
	}
	if m.Curriculum != "" {
		profilRecord["curriculum"] = m.Curriculum
	}

	if len(profilRecord) == 0 {
		return nil // Nothing to update
	}

	if exists {
		q2, _, _ := goqu.Dialect("postgres").Update("sekolah_profil").Set(profilRecord).Where(goqu.C("tenant_id").Eq(tenantID)).ToSQL()
		if _, err := a.db.Exec(q2); err != nil {
			return err
		}
	} else {
		profilRecord["tenant_id"] = tenantID
		q2, _, _ := goqu.Dialect("postgres").Insert("sekolah_profil").Rows(profilRecord).ToSQL()
		if _, err := a.db.Exec(q2); err != nil {
			return err
		}
	}
	return nil
}

// ------ Dashboard Stats ------

func (a *sekolahAdapter) GetDashboardStats(tenantID string) (*model.SekolahDashboardStats, error) {
	stats := &model.SekolahDashboardStats{}

	// 1. Count Siswa
	q1 := "SELECT COUNT(*) FROM sekolah_siswa WHERE tenant_id = $1"
	a.db.QueryRow(q1, tenantID).Scan(&stats.TotalSiswa)

	// 2. Count Guru
	q2 := "SELECT COUNT(*) FROM sekolah_guru WHERE tenant_id = $1"
	a.db.QueryRow(q2, tenantID).Scan(&stats.TotalGuru)

	// 3. Count Kelas
	q3 := "SELECT COUNT(*) FROM sekolah_kelas WHERE tenant_id = $1"
	a.db.QueryRow(q3, tenantID).Scan(&stats.TotalKelas)

	// 4. Count Mapel
	q4 := "SELECT COUNT(*) FROM sekolah_mapel WHERE tenant_id = $1"
	a.db.QueryRow(q4, tenantID).Scan(&stats.TotalMapel)

	// 5. Tagihan bulan ini (from spp_bills with current month)
	q5 := `SELECT COALESCE(SUM(amount), 0) FROM spp_bills 
		   WHERE tenant_id = $1 
		   AND billing_month = EXTRACT(MONTH FROM CURRENT_DATE)
		   AND billing_year = EXTRACT(YEAR FROM CURRENT_DATE)`
	a.db.QueryRow(q5, tenantID).Scan(&stats.TagihanBulan)

	// 6. Count Lunas
	q6 := `SELECT COUNT(*) FROM spp_bills 
		   WHERE tenant_id = $1 AND status = 'paid' 
		   AND billing_month = EXTRACT(MONTH FROM CURRENT_DATE)
		   AND billing_year = EXTRACT(YEAR FROM CURRENT_DATE)`
	a.db.QueryRow(q6, tenantID).Scan(&stats.LunasCount)

	// 7. Count Belum Lunas
	q7 := `SELECT COUNT(*) FROM spp_bills 
		   WHERE tenant_id = $1 AND status = 'pending' 
		   AND billing_month = EXTRACT(MONTH FROM CURRENT_DATE)
		   AND billing_year = EXTRACT(YEAR FROM CURRENT_DATE)`
	a.db.QueryRow(q7, tenantID).Scan(&stats.BelumLunas)

	return stats, nil
}
