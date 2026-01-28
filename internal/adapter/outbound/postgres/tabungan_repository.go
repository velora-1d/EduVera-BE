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
	// Update Tenants Table (Name, Address)
	q1, _, _ := goqu.Dialect("postgres").Update("tenants").Set(goqu.Record{
		"name":    m.NamaPesantren,
		"address": m.Alamat,
	}).Where(goqu.C("id").Eq(tenantID)).ToSQL()

	if _, err := a.db.Exec(q1); err != nil {
		return err
	}

	// Update or Insert Sekolah Profil
	// Check if exists
	var exists bool
	a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM sekolah_profil WHERE tenant_id = $1)", tenantID).Scan(&exists)

	if exists {
		q2, _, _ := goqu.Dialect("postgres").Update("sekolah_profil").Set(goqu.Record{
			"jenis_pesantren": m.JenisPesantren,
			"deskripsi":       m.Deskripsi,
		}).Where(goqu.C("tenant_id").Eq(tenantID)).ToSQL()
		if _, err := a.db.Exec(q2); err != nil {
			return err
		}
	} else {
		q2, _, _ := goqu.Dialect("postgres").Insert("sekolah_profil").Rows(goqu.Record{
			"tenant_id":       tenantID,
			"jenis_pesantren": m.JenisPesantren,
			"deskripsi":       m.Deskripsi,
		}).ToSQL()
		if _, err := a.db.Exec(q2); err != nil {
			return err
		}
	}
	return nil
}
