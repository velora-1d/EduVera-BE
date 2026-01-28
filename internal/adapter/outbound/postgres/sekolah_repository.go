package postgres_outbound_adapter

import (
	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
)

const (
	tableSiswa = "sekolah_siswa"
	tableGuru  = "sekolah_guru"
	tableMapel = "sekolah_mapel"
)

type sekolahAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewSekolahAdapter(
	db outbound_port.DatabaseExecutor,
) outbound_port.SekolahPort {
	return &sekolahAdapter{
		db: db,
	}
}

// ------ Siswa Implementation ------

func (a *sekolahAdapter) GetSiswaByTenant(tenantID string) ([]model.Siswa, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableSiswa).Where(goqu.Ex{"tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err // In production, handle table not found gracefully
	}
	defer rows.Close()

	var siswaList []model.Siswa
	for rows.Next() {
		var s model.Siswa
		// Assuming columns match struct order for simplicity in this sprint
		err := rows.Scan(&s.ID, &s.TenantID, &s.NIS, &s.Nama, &s.KelasID, &s.KelasNama, &s.Alamat, &s.Status)
		if err != nil {
			return nil, err
		}
		siswaList = append(siswaList, s)
	}

	return siswaList, nil
}

func (a *sekolahAdapter) CreateSiswa(siswa model.Siswa) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableSiswa).Rows(goqu.Record{
		"tenant_id":  siswa.TenantID,
		"nis":        siswa.NIS,
		"nama":       siswa.Nama,
		"kelas_id":   siswa.KelasID,
		"kelas_nama": siswa.KelasNama,
		"alamat":     siswa.Alamat,
		"status":     siswa.Status,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&siswa.ID)
}

// ------ Guru Implementation ------

func (a *sekolahAdapter) GetGuruByTenant(tenantID string) ([]model.Guru, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableGuru).Where(goqu.Ex{"tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guruList []model.Guru
	for rows.Next() {
		var g model.Guru
		err := rows.Scan(&g.ID, &g.TenantID, &g.NIP, &g.Nama, &g.Jenis, &g.Status)
		if err != nil {
			return nil, err
		}
		guruList = append(guruList, g)
	}
	return guruList, nil
}

func (a *sekolahAdapter) CreateGuru(guru model.Guru) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableGuru).Rows(goqu.Record{
		"tenant_id": guru.TenantID,
		"nip":       guru.NIP,
		"nama":      guru.Nama,
		"jenis":     guru.Jenis,
		"status":    guru.Status,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&guru.ID)
}

// ------ Mapel Implementation ------

func (a *sekolahAdapter) GetMapelByTenant(tenantID string) ([]model.Mapel, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableMapel).Where(goqu.Ex{"tenant_id": tenantID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mapelList []model.Mapel
	for rows.Next() {
		var m model.Mapel
		err := rows.Scan(&m.ID, &m.TenantID, &m.Kode, &m.Nama, &m.KKM)
		if err != nil {
			return nil, err
		}
		mapelList = append(mapelList, m)
	}
	return mapelList, nil
}
