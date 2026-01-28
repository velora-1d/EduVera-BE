package model

// Sekolah Domain Models

type Siswa struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	NIS       string `json:"nis"`
	Nama      string `json:"nama"`
	KelasID   string `json:"kelas_id"`
	KelasNama string `json:"kelas_nama"` // Populated from join
	Alamat    string `json:"alamat"`
	NamaWali  string `json:"nama_wali"`
	NoHPWali  string `json:"no_hp_wali"`
	Status    string `json:"status"` // Aktif, Lulus, Pindah
}

type Guru struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	NIP      string `json:"nip"`
	Nama     string `json:"nama"`
	Jenis    string `json:"jenis"`  // Guru Mapel, Guru Kelas
	Status   string `json:"status"` // PNS, Honorer
}

type Mapel struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	KKM      int    `json:"kkm"`
}

type Kelas struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Nama     string `json:"nama"`
	Tingkat  string `json:"tingkat"` // Level like "10", "11", "12" or "Ula", "Wustha"
	Urutan   int    `json:"urutan"`
	Status   string `json:"status"`
}
