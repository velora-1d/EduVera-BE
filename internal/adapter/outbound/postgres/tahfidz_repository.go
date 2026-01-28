package postgres_outbound_adapter

import (
	"database/sql"
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

var tableTahfidzSetoran = goqu.T("sekolah_tahfidz_setoran")

func (a *sekolahAdapter) GetTahfidzSetoran(tenantID string) ([]model.TahfidzSetoran, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableTahfidzSetoran).
		Join(tableSiswa, goqu.On(tableTahfidzSetoran.Col("santri_id").Eq(tableSiswa.Col("id")))).
		LeftJoin(tableGuru, goqu.On(tableTahfidzSetoran.Col("ustadz_id").Eq(tableGuru.Col("id")))).
		Select(
			tableTahfidzSetoran.Col("id"),
			tableTahfidzSetoran.Col("tenant_id"),
			tableTahfidzSetoran.Col("santri_id"),
			tableSiswa.Col("nama").As("santri_nama"),
			tableTahfidzSetoran.Col("ustadz_id"),
			goqu.COALESCE(tableGuru.Col("nama"), "").As("ustadz_nama"),
			tableTahfidzSetoran.Col("tanggal"),
			tableTahfidzSetoran.Col("juz"),
			goqu.COALESCE(tableTahfidzSetoran.Col("surah"), "").As("surah"),
			tableTahfidzSetoran.Col("ayat_awal"),
			tableTahfidzSetoran.Col("ayat_akhir"),
			tableTahfidzSetoran.Col("tipe"),
			goqu.COALESCE(tableTahfidzSetoran.Col("kualitas"), "").As("kualitas"),
			goqu.COALESCE(tableTahfidzSetoran.Col("catatan"), "").As("catatan"),
			tableTahfidzSetoran.Col("created_at"),
			tableTahfidzSetoran.Col("updated_at"),
		).Where(tableTahfidzSetoran.Col("tenant_id").Eq(tenantID)).
		Order(tableTahfidzSetoran.Col("tanggal").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TahfidzSetoran
	for rows.Next() {
		var m model.TahfidzSetoran
		var ustadzID sql.NullString
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.SantriID, &m.SantriNama, &ustadzID, &m.UstadzNama,
			&m.Tanggal, &m.Juz, &m.Surah, &m.AyatAwal, &m.AyatAkhir,
			&m.Tipe, &m.Kualitas, &m.Catatan, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if ustadzID.Valid {
			id := ustadzID.String
			m.UstadzID = &id
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateTahfidzSetoran(m *model.TahfidzSetoran) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tableTahfidzSetoran).Rows(goqu.Record{
		"tenant_id":  m.TenantID,
		"santri_id":  m.SantriID,
		"ustadz_id":  m.UstadzID,
		"tanggal":    m.Tanggal,
		"juz":        m.Juz,
		"surah":      m.Surah,
		"ayat_awal":  m.AyatAwal,
		"ayat_akhir": m.AyatAkhir,
		"tipe":       m.Tipe,
		"kualitas":   m.Kualitas,
		"catatan":    m.Catatan,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
