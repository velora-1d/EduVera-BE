package model

// Roles for Sekolah (10 roles)
const (
	RoleAdminSekolah  = "admin_sekolah"
	RoleKepalaSekolah = "kepala_sekolah"
	RoleWakilKepsek   = "wakil_kepsek"
	RoleWaliKelas     = "wali_kelas"
	RoleGuru          = "guru"
	RoleTataUsaha     = "tata_usaha"
	RoleBendahara     = "bendahara"
	RoleBK            = "bk"
	RolePerpustakaan  = "perpustakaan"
	RoleWaliSiswa     = "wali_siswa"
)

// Roles for Pesantren (6 roles)
const (
	RoleAdminPesantren = "admin_pesantren"
	RolePengasuh       = "pengasuh"
	RoleSekretaris     = "sekretaris"
	RoleBendaharaPes   = "bendahara_pesantren"
	RolePendidikan     = "pendidikan"
	RoleWaliSantri     = "wali_santri"
)

// Role metadata
type RoleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlanType    string `json:"plan_type"` // sekolah, pesantren
}

// All Sekolah roles
var SekolahRoles = []RoleInfo{
	{ID: RoleAdminSekolah, Name: "Admin Sekolah", Description: "Akses penuh ke semua fitur sekolah", PlanType: "sekolah"},
	{ID: RoleKepalaSekolah, Name: "Kepala Sekolah", Description: "Akses dashboard, laporan, dan approval", PlanType: "sekolah"},
	{ID: RoleWakilKepsek, Name: "Wakil Kepala Sekolah", Description: "Akses kurikulum, akademik, dan SDM", PlanType: "sekolah"},
	{ID: RoleWaliKelas, Name: "Wali Kelas", Description: "Akses kelas, siswa, nilai, dan absensi", PlanType: "sekolah"},
	{ID: RoleGuru, Name: "Guru", Description: "Akses jadwal, nilai, dan mapel", PlanType: "sekolah"},
	{ID: RoleTataUsaha, Name: "Tata Usaha", Description: "Akses data siswa, surat, dan arsip", PlanType: "sekolah"},
	{ID: RoleBendahara, Name: "Bendahara", Description: "Akses keuangan, SPP, dan laporan keuangan", PlanType: "sekolah"},
	{ID: RoleBK, Name: "Bimbingan Konseling", Description: "Akses data siswa dan konseling", PlanType: "sekolah"},
	{ID: RolePerpustakaan, Name: "Pustakawan", Description: "Akses perpustakaan dan peminjaman", PlanType: "sekolah"},
	{ID: RoleWaliSiswa, Name: "Wali Siswa", Description: "Akses nilai anak, SPP, dan pengumuman", PlanType: "sekolah"},
}

// All Pesantren roles
var PesantrenRoles = []RoleInfo{
	{ID: RoleAdminPesantren, Name: "Admin Pesantren", Description: "Akses penuh ke semua fitur pesantren", PlanType: "pesantren"},
	{ID: RolePengasuh, Name: "Pengasuh / Kyai", Description: "Akses dashboard, laporan, dan approval", PlanType: "pesantren"},
	{ID: RoleSekretaris, Name: "Sekretaris", Description: "Akses data santri, surat, dan arsip", PlanType: "pesantren"},
	{ID: RoleBendaharaPes, Name: "Bendahara", Description: "Akses keuangan, syahriah, dan laporan", PlanType: "pesantren"},
	{ID: RolePendidikan, Name: "Bagian Pendidikan", Description: "Akses tahfidz, diniyah, dan nilai santri", PlanType: "pesantren"},
	{ID: RoleWaliSantri, Name: "Wali Santri", Description: "Akses nilai anak, syahriah, dan pengumuman", PlanType: "pesantren"},
}

// GetRolesByPlanType returns all roles for a given plan type
func GetRolesByPlanType(planType string) []RoleInfo {
	switch planType {
	case "sekolah":
		return SekolahRoles
	case "pesantren":
		return PesantrenRoles
	case "hybrid":
		// Hybrid gets both
		all := make([]RoleInfo, 0, len(SekolahRoles)+len(PesantrenRoles))
		all = append(all, SekolahRoles...)
		all = append(all, PesantrenRoles...)
		return all
	default:
		return nil
	}
}

// IsValidRole checks if a role is valid for a given plan type
func IsValidRole(role string, planType string) bool {
	roles := GetRolesByPlanType(planType)
	for _, r := range roles {
		if r.ID == role {
			return true
		}
	}
	return false
}
