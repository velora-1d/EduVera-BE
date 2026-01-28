package postgres_outbound_adapter

import (
	"database/sql"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
)

var (
	tableSiswa      = goqu.T("sekolah_siswa")
	tableAsrama     = goqu.T("sekolah_asrama")
	tableKamar      = goqu.T("sekolah_kamar")
	tablePenempatan = goqu.T("sekolah_penempatan")
	tableGuru       = goqu.T("sekolah_guru")
	tableMapel      = goqu.T("sekolah_mapel")

	tablePelanggaranAturan = goqu.T("sekolah_pelanggaran_aturan")
	tablePelanggaranSiswa  = goqu.T("sekolah_pelanggaran_siswa")
	tablePerizinan         = goqu.T("sekolah_perizinan")
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
	// Explicit select to ensure order and handle NULLs if needed via COALESCE or just matching types
	dataset := dialect.From(tableSiswa).Select(
		"id", "tenant_id", "nis", "nama", "kelas_id", "kelas_nama", "alamat", "nama_wali", "no_hp_wali", "status",
	).Where(goqu.Ex{"tenant_id": tenantID})

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
		// Scan into sql.NullString for nullable columns to be safe, then assign to struct
		var namaWali, noHpWali sql.NullString
		err := rows.Scan(&s.ID, &s.TenantID, &s.NIS, &s.Nama, &s.KelasID, &s.KelasNama, &s.Alamat, &namaWali, &noHpWali, &s.Status)
		if err != nil {
			return nil, err
		}
		s.NamaWali = namaWali.String
		s.NoHPWali = noHpWali.String
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
		"nama_wali":  siswa.NamaWali,
		"no_hp_wali": siswa.NoHPWali,
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

// ------ Kelas Implementation ------

func (a *sekolahAdapter) GetKelasByTenant(tenantID string) ([]model.Kelas, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From("sekolah_kelas").Where(goqu.Ex{"tenant_id": tenantID}).Order(goqu.I("urutan").Asc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kelasList []model.Kelas
	for rows.Next() {
		var k model.Kelas
		// Columns: id, tenant_id, nama, tingkat, urutan, status, created_at, updated_at
		// Need to scan properly. Assuming created_at/updated_at are ignored for now in struct,
		// but Scan must match table columns or we select specific columns.
		// Safer to assume SELECT * returns all, so we must scan all or change query to SELECT cols.
		// Let's change query to select explicit columns to be safe.
		// Re-writing query logic below.
		err := rows.Scan(&k.ID, &k.TenantID, &k.Nama, &k.Tingkat, &k.Urutan, &k.Status, &sql.NullTime{}, &sql.NullTime{})
		if err != nil {
			return nil, err
		}
		kelasList = append(kelasList, k)
	}
	return kelasList, nil
}

func (a *sekolahAdapter) CreateKelas(kelas *model.Kelas) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert("sekolah_kelas").Rows(goqu.Record{
		"tenant_id": kelas.TenantID,
		"nama":      kelas.Nama,
		"tingkat":   kelas.Tingkat,
		"urutan":    kelas.Urutan,
		"status":    kelas.Status,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&kelas.ID)
}

// ------ Asrama Implementation ------

func (a *sekolahAdapter) GetAsramaByTenant(tenantID string) ([]model.Asrama, error) {
	dialect := goqu.Dialect("postgres")
	// Join with guru for musyrif name, left join in case musyrif is null or deleted
	dataset := dialect.From(tableAsrama).
		LeftJoin(
			tableGuru.As("g"),
			goqu.On(tableAsrama.Col("musyrif_id").Eq(goqu.I("g.id"))),
		).
		Select(
			tableAsrama.Col("id"),
			tableAsrama.Col("tenant_id"),
			tableAsrama.Col("nama"),
			tableAsrama.Col("jenis"),
			tableAsrama.Col("musyrif_id"),
			goqu.COALESCE(goqu.I("g.nama"), "").As("musyrif_nama"),
			tableAsrama.Col("status"),
			tableAsrama.Col("created_at"),
			tableAsrama.Col("updated_at"),
		).Where(tableAsrama.Col("tenant_id").Eq(tenantID))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Asrama
	for rows.Next() {
		var m model.Asrama
		var musyrifID sql.NullString
		// Scan basic fields
		if err := rows.Scan(&m.ID, &m.TenantID, &m.Nama, &m.Jenis, &musyrifID, &m.Musyrif, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if musyrifID.Valid {
			id := musyrifID.String
			m.MusyrifID = &id
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateAsrama(m *model.Asrama) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tableAsrama).Rows(goqu.Record{
		"tenant_id":  m.TenantID,
		"nama":       m.Nama,
		"jenis":      m.Jenis,
		"musyrif_id": m.MusyrifID,
		"status":     m.Status,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (a *sekolahAdapter) GetKamarByAsrama(tenantID, asramaID string) ([]model.Kamar, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableKamar).
		Join(tableAsrama, goqu.On(tableKamar.Col("asrama_id").Eq(tableAsrama.Col("id")))).
		Select(
			tableKamar.Col("id"),
			tableKamar.Col("tenant_id"),
			tableKamar.Col("asrama_id"),
			tableAsrama.Col("nama").As("asrama_nama"),
			tableKamar.Col("nomor"),
			tableKamar.Col("kapasitas"),
			tableKamar.Col("status"),
			// Subquery for occupied count
			dialect.From(tablePenempatan).
				Select(goqu.COUNT("*")).
				Where(
					tablePenempatan.Col("kamar_id").Eq(tableKamar.Col("id")),
					tablePenempatan.Col("status").Eq("Aktif"),
				).As("terisi"),
			tableKamar.Col("created_at"),
			tableKamar.Col("updated_at"),
		).Where(
		tableKamar.Col("tenant_id").Eq(tenantID),
		goqu.Or(
			goqu.L("? = ''", asramaID), // If asramaID empty, return all
			tableKamar.Col("asrama_id").Eq(asramaID),
		),
	).Order(tableKamar.Col("nomor").Asc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Kamar
	for rows.Next() {
		var m model.Kamar
		if err := rows.Scan(&m.ID, &m.TenantID, &m.AsramaID, &m.AsramaNama, &m.Nomor, &m.Kapasitas, &m.Status, &m.Terisi, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreateKamar(m *model.Kamar) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tableKamar).Rows(goqu.Record{
		"tenant_id": m.TenantID,
		"asrama_id": m.AsramaID,
		"nomor":     m.Nomor,
		"kapasitas": m.Kapasitas,
		"status":    m.Status,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (a *sekolahAdapter) GetPenempatanByTenant(tenantID string) ([]model.Penempatan, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tablePenempatan).
		Join(tableSiswa, goqu.On(tablePenempatan.Col("santri_id").Eq(tableSiswa.Col("id")))).
		Join(tableKamar, goqu.On(tablePenempatan.Col("kamar_id").Eq(tableKamar.Col("id")))).
		Join(tableAsrama, goqu.On(tableKamar.Col("asrama_id").Eq(tableAsrama.Col("id")))).
		Select(
			tablePenempatan.Col("id"),
			tablePenempatan.Col("tenant_id"),
			tablePenempatan.Col("santri_id"),
			tableSiswa.Col("nama").As("santri_nama"),
			tablePenempatan.Col("kamar_id"),
			tableKamar.Col("nomor").As("kamar_nomor"),
			tableAsrama.Col("nama").As("asrama_nama"),
			tablePenempatan.Col("tanggal_masuk"),
			tablePenempatan.Col("status"),
			goqu.COALESCE(tablePenempatan.Col("keterangan"), "").As("keterangan"),
			tablePenempatan.Col("created_at"),
			tablePenempatan.Col("updated_at"),
		).Where(tablePenempatan.Col("tenant_id").Eq(tenantID)).
		Order(tablePenempatan.Col("created_at").Desc())

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Penempatan
	for rows.Next() {
		var m model.Penempatan
		if err := rows.Scan(&m.ID, &m.TenantID, &m.SantriID, &m.SantriNama, &m.KamarID, &m.KamarNomor, &m.AsramaNama, &m.TanggalMasuk, &m.Status, &m.Keterangan, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (a *sekolahAdapter) CreatePenempatan(m *model.Penempatan) error {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(tablePenempatan).Rows(goqu.Record{
		"tenant_id":     m.TenantID,
		"santri_id":     m.SantriID,
		"kamar_id":      m.KamarID,
		"tanggal_masuk": m.TanggalMasuk,
		"status":        m.Status,
		"keterangan":    m.Keterangan,
	}).Returning("id", "created_at", "updated_at")

	query, _, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return a.db.QueryRow(query).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}
