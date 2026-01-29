package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCacheFields, downCacheFields)
}

// upCacheFields adds denormalized cache fields to reduce JOINs in common queries
// Based on Phase 3.3: Denormalized Cache roadmap
func upCacheFields(ctx context.Context, tx *sql.Tx) error {
	// Add cached fields to sekolah_siswa for faster lookups without JOINs
	queries := []struct {
		name  string
		query string
	}{
		// ==========================================
		// Cached class name in siswa table
		// ==========================================
		{
			name: "add cached_kelas_nama to sekolah_siswa",
			query: `ALTER TABLE sekolah_siswa 
				ADD COLUMN IF NOT EXISTS cached_kelas_nama VARCHAR(100)`,
		},
		{
			name: "add cached_wali_kelas to sekolah_siswa",
			query: `ALTER TABLE sekolah_siswa 
				ADD COLUMN IF NOT EXISTS cached_wali_kelas VARCHAR(100)`,
		},

		// ==========================================
		// Cached student info in student_grades
		// ==========================================
		{
			name: "add cached_student_nama to student_grades",
			query: `ALTER TABLE student_grades 
				ADD COLUMN IF NOT EXISTS cached_student_nama VARCHAR(255)`,
		},
		{
			name: "add cached_student_nis to student_grades",
			query: `ALTER TABLE student_grades 
				ADD COLUMN IF NOT EXISTS cached_student_nis VARCHAR(50)`,
		},
		{
			name: "add cached_subject_nama to student_grades",
			query: `ALTER TABLE student_grades 
				ADD COLUMN IF NOT EXISTS cached_subject_nama VARCHAR(255)`,
		},

		// ==========================================
		// Cached santri info in tahfidz_setoran
		// ==========================================
		{
			name: "add cached_santri_nama to tahfidz_setoran",
			query: `ALTER TABLE tahfidz_setoran 
				ADD COLUMN IF NOT EXISTS cached_santri_nama VARCHAR(255)`,
		},

		// ==========================================
		// Cached student info in spp_bills
		// ==========================================
		{
			name: "add cached_student_nama to spp_bills",
			query: `ALTER TABLE spp_bills 
				ADD COLUMN IF NOT EXISTS cached_student_nama VARCHAR(255)`,
		},
		{
			name: "add cached_student_nis to spp_bills",
			query: `ALTER TABLE spp_bills 
				ADD COLUMN IF NOT EXISTS cached_student_nis VARCHAR(50)`,
		},
		{
			name: "add cached_kelas_nama to spp_bills",
			query: `ALTER TABLE spp_bills 
				ADD COLUMN IF NOT EXISTS cached_kelas_nama VARCHAR(100)`,
		},

		// ==========================================
		// Cached info in sekolah_rapor
		// ==========================================
		{
			name: "add cached_santri_nama to sekolah_rapor",
			query: `ALTER TABLE sekolah_rapor 
				ADD COLUMN IF NOT EXISTS cached_santri_nama VARCHAR(255)`,
		},
		{
			name: "add cached_santri_nis to sekolah_rapor",
			query: `ALTER TABLE sekolah_rapor 
				ADD COLUMN IF NOT EXISTS cached_santri_nis VARCHAR(50)`,
		},
		{
			name: "add cached_kelas_nama to sekolah_rapor",
			query: `ALTER TABLE sekolah_rapor 
				ADD COLUMN IF NOT EXISTS cached_kelas_nama VARCHAR(100)`,
		},
	}

	for _, q := range queries {
		if _, err := tx.Exec(q.query); err != nil {
			// Log but don't fail - table might not exist yet
			fmt.Printf("Warning: cache field %s failed: %v\n", q.name, err)
		}
	}

	// Create function to sync cache on siswa update
	cacheFunction := `
		CREATE OR REPLACE FUNCTION sync_siswa_kelas_cache()
		RETURNS TRIGGER AS $$
		BEGIN
			IF NEW.kelas_id IS NOT NULL THEN
				SELECT nama INTO NEW.cached_kelas_nama 
				FROM sekolah_kelas WHERE id = NEW.kelas_id;
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	if _, err := tx.Exec(cacheFunction); err != nil {
		fmt.Printf("Warning: cache function failed: %v\n", err)
	}

	// Create trigger to auto-sync cache
	cacheTrigger := `
		DROP TRIGGER IF EXISTS sync_siswa_kelas_cache_trigger ON sekolah_siswa;
		CREATE TRIGGER sync_siswa_kelas_cache_trigger
		BEFORE INSERT OR UPDATE ON sekolah_siswa
		FOR EACH ROW
		EXECUTE FUNCTION sync_siswa_kelas_cache();
	`
	if _, err := tx.Exec(cacheTrigger); err != nil {
		fmt.Printf("Warning: cache trigger failed: %v\n", err)
	}

	return nil
}

func downCacheFields(ctx context.Context, tx *sql.Tx) error {
	// Drop trigger and function first
	tx.Exec("DROP TRIGGER IF EXISTS sync_siswa_kelas_cache_trigger ON sekolah_siswa")
	tx.Exec("DROP FUNCTION IF EXISTS sync_siswa_kelas_cache")

	// Drop cached columns
	columns := []struct {
		table  string
		column string
	}{
		{"sekolah_siswa", "cached_kelas_nama"},
		{"sekolah_siswa", "cached_wali_kelas"},
		{"student_grades", "cached_student_nama"},
		{"student_grades", "cached_student_nis"},
		{"student_grades", "cached_subject_nama"},
		{"tahfidz_setoran", "cached_santri_nama"},
		{"spp_bills", "cached_student_nama"},
		{"spp_bills", "cached_student_nis"},
		{"spp_bills", "cached_kelas_nama"},
		{"sekolah_rapor", "cached_santri_nama"},
		{"sekolah_rapor", "cached_santri_nis"},
		{"sekolah_rapor", "cached_kelas_nama"},
	}

	for _, c := range columns {
		query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", c.table, c.column)
		if _, err := tx.Exec(query); err != nil {
			fmt.Printf("Warning: dropping %s.%s failed: %v\n", c.table, c.column, err)
		}
	}

	return nil
}
