package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upIndexStrategy, downIndexStrategy)
}

// upIndexStrategy adds performance indexes for common query patterns
// Based on Phase 3: Database Optimization roadmap
// Using NoTx so individual index failures don't abort entire migration
func upIndexStrategy(ctx context.Context, db *sql.DB) error {
	queries := []struct {
		name  string
		query string
	}{
		// ==========================================
		// SPP & Keuangan - High traffic queries
		// ==========================================
		{
			name: "spp_bills tenant+status (filter tagihan)",
			query: `CREATE INDEX IF NOT EXISTS idx_spp_bills_tenant_status 
				ON spp_bills(tenant_id, status)`,
		},
		{
			name: "spp_bills tenant+period (filter bulanan)",
			query: `CREATE INDEX IF NOT EXISTS idx_spp_bills_tenant_period 
				ON spp_bills(tenant_id, bill_period)`,
		},
		{
			name: "spp_bills student (filter per siswa)",
			query: `CREATE INDEX IF NOT EXISTS idx_spp_bills_student 
				ON spp_bills(student_id)`,
		},

		// ==========================================
		// Siswa - Status filter & search
		// ==========================================
		{
			name: "sekolah_siswa tenant+status (filter status)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_siswa_tenant_status 
				ON sekolah_siswa(tenant_id, status)`,
		},
		{
			name: "sekolah_siswa nisn (lookup unique)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_siswa_nisn 
				ON sekolah_siswa(nisn) WHERE nisn IS NOT NULL AND nisn != ''`,
		},
		{
			name: "sekolah_siswa nama (text search)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_siswa_nama 
				ON sekolah_siswa USING gin(to_tsvector('simple', nama))`,
		},

		// ==========================================
		// Guru - Status & jenis filter
		// ==========================================
		{
			name: "sekolah_guru tenant+jenis (filter)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_guru_tenant_jenis 
				ON sekolah_guru(tenant_id, jenis)`,
		},

		// ==========================================
		// Student Grades - Composite for rapor
		// ==========================================
		{
			name: "student_grades composite (rapor query)",
			query: `CREATE INDEX IF NOT EXISTS idx_student_grades_rapor 
				ON student_grades(tenant_id, student_id, academic_year, semester)`,
		},

		// ==========================================
		// Tahfidz Setoran - Progress tracking
		// ==========================================
		{
			name: "tahfidz_setoran tenant+santri (progress)",
			query: `CREATE INDEX IF NOT EXISTS idx_tahfidz_setoran_tenant_santri 
				ON tahfidz_setoran(tenant_id, santri_id)`,
		},
		{
			name: "tahfidz_setoran tanggal (timeline)",
			query: `CREATE INDEX IF NOT EXISTS idx_tahfidz_setoran_tanggal 
				ON tahfidz_setoran(tanggal DESC)`,
		},

		// ==========================================
		// Rapor - Generated reports lookup
		// ==========================================
		{
			name: "sekolah_rapor tenant+santri (lookup)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_rapor_tenant_santri 
				ON sekolah_rapor(tenant_id, santri_id)`,
		},
		{
			name: "sekolah_rapor_nilai rapor_id (detail)",
			query: `CREATE INDEX IF NOT EXISTS idx_sekolah_rapor_nilai_rapor 
				ON sekolah_rapor_nilai(rapor_id)`,
		},

		// ==========================================
		// Kepesantrenan - Pelanggaran records
		// ==========================================
		{
			name: "pelanggaran_siswa tenant+santri",
			query: `CREATE INDEX IF NOT EXISTS idx_pelanggaran_siswa_tenant_santri 
				ON sekolah_pelanggaran_siswa(tenant_id, santri_id)`,
		},

		// ==========================================
		// Perizinan - Approval workflow
		// ==========================================
		{
			name: "perizinan tenant+status",
			query: `CREATE INDEX IF NOT EXISTS idx_perizinan_siswa_tenant_status 
				ON sekolah_perizinan_siswa(tenant_id, status)`,
		},

		// ==========================================
		// Asrama - Penempatan lookup
		// ==========================================
		{
			name: "asrama_penempatan santri",
			query: `CREATE INDEX IF NOT EXISTS idx_asrama_penempatan_santri 
				ON pesantren_asrama_penempatan(santri_id)`,
		},

		// ==========================================
		// SDM - Employee management
		// ==========================================
		{
			name: "sdm_employees tenant+status",
			query: `CREATE INDEX IF NOT EXISTS idx_sdm_employees_tenant_status 
				ON sdm_employees(tenant_id, status)`,
		},
		{
			name: "sdm_payroll tenant+period",
			query: `CREATE INDEX IF NOT EXISTS idx_sdm_payroll_tenant_period 
				ON sdm_payroll(tenant_id, period)`,
		},
	}

	for _, q := range queries {
		if _, err := db.Exec(q.query); err != nil {
			// Log but don't fail - table might not exist yet
			fmt.Printf("Warning: index %s failed: %v\n", q.name, err)
		}
	}

	return nil
}

func downIndexStrategy(ctx context.Context, db *sql.DB) error {
	// Drop indexes in reverse order
	indexes := []string{
		"idx_sdm_payroll_tenant_period",
		"idx_sdm_employees_tenant_status",
		"idx_asrama_penempatan_santri",
		"idx_perizinan_siswa_tenant_status",
		"idx_pelanggaran_siswa_tenant_santri",
		"idx_sekolah_rapor_nilai_rapor",
		"idx_sekolah_rapor_tenant_santri",
		"idx_tahfidz_setoran_tanggal",
		"idx_tahfidz_setoran_tenant_santri",
		"idx_student_grades_rapor",
		"idx_sekolah_guru_tenant_jenis",
		"idx_sekolah_siswa_nama",
		"idx_sekolah_siswa_nisn",
		"idx_sekolah_siswa_tenant_status",
		"idx_spp_bills_student",
		"idx_spp_bills_tenant_period",
		"idx_spp_bills_tenant_status",
	}

	for _, idx := range indexes {
		if _, err := db.Exec("DROP INDEX IF EXISTS " + idx); err != nil {
			fmt.Printf("Warning: dropping index %s failed: %v\n", idx, err)
		}
	}

	return nil
}
