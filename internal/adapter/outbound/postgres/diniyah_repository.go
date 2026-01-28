package postgres_outbound_adapter

import (
	"prabogo/internal/model"

	"github.com/doug-martin/goqu/v9"
)

var tableDiniyahKitab = goqu.T("sekolah_diniyah_kitab")

func (a *sekolahAdapter) GetDiniyahKitab(tenantID string) ([]model.DiniyahKitab, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableDiniyahKitab).
		Select(
			tableDiniyahKitab.Col("id"),
			tableDiniyahKitab.Col("tenant_id"),
			tableDiniyahKitab.Col("nama_kitab"),
			goqu.COALESCE(tableDiniyahKitab.Col("bidang_studi"), "").As("bidang_studi"),
			goqu.COALESCE(tableDiniyahKitab.Col("pengarang"), "").As("pengarang"),
			goqu.COALESCE(tableDiniyahKitab.Col("keterangan"), "").As("keterangan"),
			tableDiniyahKitab.Col("created_at"),
			tableDiniyahKitab.Col("updated_at"),
		).Where(tableDiniyahKitab.Col("tenant_id").Eq(tenantID)).
		Order(tableDiniyahKitab.Col("nama_kitab").Asc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.DiniyahKitab
	for rows.Next() {
		var m model.DiniyahKitab
		if err := rows.Scan(
			&m.ID, &m.TenantID, &m.NamaKitab, &m.BidangStudi, &m.Pengarang, &m.Keterangan,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateDiniyahKitab(m *model.DiniyahKitab) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tableDiniyahKitab).Rows(goqu.Record{
		"tenant_id":    m.TenantID,
		"nama_kitab":   m.NamaKitab,
		"bidang_studi": m.BidangStudi,
		"pengarang":    m.Pengarang,
		"keterangan":   m.Keterangan,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
