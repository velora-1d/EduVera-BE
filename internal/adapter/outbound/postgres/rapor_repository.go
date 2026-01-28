package postgres_outbound_adapter

import (
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

var tableRapor = goqu.T("sekolah_rapor")
var tableRaporNilai = goqu.T("sekolah_rapor_nilai")
var tableRaporPeriode = goqu.T("sekolah_rapor_periode")

func (a *sekolahAdapter) GetRaporList(tenantID string) ([]model.Rapor, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableRapor).
		Select(
			tableRapor.Col("id"),
			tableRapor.Col("tenant_id"),
			tableRapor.Col("periode_id"),
			tableRapor.Col("santri_id"),
			tableRapor.Col("status"),
			goqu.COALESCE(tableRapor.Col("catatan_wali_kelas"), "").As("catatan_wali_kelas"),
			tableRapor.Col("created_at"),
			tableRapor.Col("updated_at"),
			goqu.COALESCE(tableSiswa.Col("nama_lengkap"), "").As("nama_santri"),
			goqu.COALESCE(tableRaporPeriode.Col("nama"), "").As("nama_periode"),
		).
		LeftJoin(tableSiswa, goqu.On(tableRapor.Col("santri_id").Eq(tableSiswa.Col("id")))).
		LeftJoin(tableRaporPeriode, goqu.On(tableRapor.Col("periode_id").Eq(tableRaporPeriode.Col("id")))).
		Where(tableRapor.Col("tenant_id").Eq(tenantID)).
		Order(tableRapor.Col("created_at").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Rapor
	for rows.Next() {
		var m model.Rapor
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.PeriodeID, &m.SantriID, &m.Status, &m.CatatanWaliKelas,
			&m.CreatedAt, &m.UpdatedAt, &m.NamaSantri, &m.NamaPeriode,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	// TODO: Fetch NilaiList for each Rapor if needed, but list view usually light.
	return list, nil
}

func (a *sekolahAdapter) CreateRapor(m *model.Rapor) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tableRapor).Rows(goqu.Record{
		"tenant_id":          m.TenantID,
		"periode_id":         m.PeriodeID,
		"santri_id":          m.SantriID,
		"status":             m.Status,
		"catatan_wali_kelas": m.CatatanWaliKelas,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
