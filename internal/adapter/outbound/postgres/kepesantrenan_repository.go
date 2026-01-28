package postgres_outbound_adapter

import (
	"database/sql"
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

// ------ Kepesantrenan Implementation ------

// Rules
func (a *sekolahAdapter) GetPelanggaranAturan(tenantID string) ([]model.PelanggaranAturan, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tablePelanggaranAturan).Where(goqu.Ex{"tenant_id": tenantID}).Order(goqu.I("poin").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PelanggaranAturan
	for rows.Next() {
		var m model.PelanggaranAturan
		if err := rows.Scan(&m.ID, &m.TenantID, &m.Judul, &m.Kategori, &m.Poin, &m.Level, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreatePelanggaranAturan(m *model.PelanggaranAturan) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tablePelanggaranAturan).Rows(goqu.Record{
		"tenant_id": m.TenantID,
		"judul":     m.Judul,
		"kategori":  m.Kategori,
		"poin":      m.Poin,
		"level":     m.Level,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

// Violations
func (a *sekolahAdapter) GetPelanggaranSiswa(tenantID string) ([]model.PelanggaranSiswa, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tablePelanggaranSiswa).
		Join(tableSiswa, goqu.On(tablePelanggaranSiswa.Col("santri_id").Eq(tableSiswa.Col("id")))).
		LeftJoin(tablePelanggaranAturan, goqu.On(tablePelanggaranSiswa.Col("aturan_id").Eq(tablePelanggaranAturan.Col("id")))). // Left join in case rule deleted
		Select(
			tablePelanggaranSiswa.Col("id"),
			tablePelanggaranSiswa.Col("tenant_id"),
			tablePelanggaranSiswa.Col("santri_id"),
			tableSiswa.Col("nama").As("santri_nama"),
			tablePelanggaranSiswa.Col("aturan_id"),
			goqu.COALESCE(tablePelanggaranAturan.Col("judul"), "").As("aturan_judul"),
			tablePelanggaranSiswa.Col("tanggal"),
			tablePelanggaranSiswa.Col("poin"),
			goqu.COALESCE(tablePelanggaranSiswa.Col("keterangan"), "").As("keterangan"),
			tablePelanggaranSiswa.Col("status"),
			goqu.COALESCE(tablePelanggaranSiswa.Col("sanksi"), "").As("sanksi"),
			tablePelanggaranSiswa.Col("created_at"),
			tablePelanggaranSiswa.Col("updated_at"),
		).Where(tablePelanggaranSiswa.Col("tenant_id").Eq(tenantID)).
		Order(tablePelanggaranSiswa.Col("tanggal").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PelanggaranSiswa
	for rows.Next() {
		var m model.PelanggaranSiswa
		var aturanID sql.NullString
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.SantriID, &m.SantriNama, &aturanID, &m.AturanJudul,
			&m.Tanggal, &m.Poin, &m.Keterangan, &m.Status, &m.Sanksi,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if aturanID.Valid {
			id := aturanID.String
			m.AturanID = &id
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreatePelanggaranSiswa(m *model.PelanggaranSiswa) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tablePelanggaranSiswa).Rows(goqu.Record{
		"tenant_id":  m.TenantID,
		"santri_id":  m.SantriID,
		"aturan_id":  m.AturanID,
		"tanggal":    m.Tanggal,
		"poin":       m.Poin,
		"keterangan": m.Keterangan,
		"status":     m.Status,
		"sanksi":     m.Sanksi,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

// Permissions
func (a *sekolahAdapter) GetPerizinan(tenantID string) ([]model.Perizinan, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tablePerizinan).
		Join(tableSiswa, goqu.On(tablePerizinan.Col("santri_id").Eq(tableSiswa.Col("id")))).
		LeftJoin(tableGuru, goqu.On(tablePerizinan.Col("penyetuju_id").Eq(tableGuru.Col("id")))).
		Select(
			tablePerizinan.Col("id"),
			tablePerizinan.Col("tenant_id"),
			tablePerizinan.Col("santri_id"),
			tableSiswa.Col("nama").As("santri_nama"),
			tablePerizinan.Col("tipe"),
			goqu.COALESCE(tablePerizinan.Col("alasan"), "").As("alasan"),
			tablePerizinan.Col("dari"),
			tablePerizinan.Col("sampai"),
			tablePerizinan.Col("status"),
			tablePerizinan.Col("penyetuju_id"),
			goqu.COALESCE(tableGuru.Col("nama"), "").As("penyetuju_nama"),
			tablePerizinan.Col("created_at"),
			tablePerizinan.Col("updated_at"),
		).Where(tablePerizinan.Col("tenant_id").Eq(tenantID)).
		Order(tablePerizinan.Col("created_at").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Perizinan
	for rows.Next() {
		var m model.Perizinan
		var penyetujuID sql.NullString
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.SantriID, &m.SantriNama, &m.Tipe, &m.Alasan, &m.Dari, &m.Sampai, &m.Status,
			&penyetujuID, &m.PenyetujuNama, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if penyetujuID.Valid {
			id := penyetujuID.String
			m.PenyetujuID = &id
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreatePerizinan(m *model.Perizinan) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tablePerizinan).Rows(goqu.Record{
		"tenant_id":    m.TenantID,
		"santri_id":    m.SantriID,
		"tipe":         m.Tipe,
		"alasan":       m.Alasan,
		"dari":         m.Dari,
		"sampai":       m.Sampai,
		"status":       m.Status,
		"penyetuju_id": m.PenyetujuID,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
