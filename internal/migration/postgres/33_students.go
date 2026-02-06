package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUnifiedStudents, downUnifiedStudents)
}

func upUnifiedStudents(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		-- Unified Students Table (Siswa + Santri)
		CREATE TABLE IF NOT EXISTS students (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			
			-- Basic Info
			name VARCHAR(255) NOT NULL,
			nickname VARCHAR(100),
			gender VARCHAR(10), -- L/P or Laki-laki/Perempuan
			birth_place VARCHAR(100),
			birth_date DATE,
			photo_url TEXT,
			
			-- NIS (multiple types supported)
			nis VARCHAR(50),              -- NIS Sekolah
			nisn VARCHAR(20),             -- NISN Nasional
			nis_pesantren VARCHAR(50),    -- NIS Pondok/Pesantren
			
			-- Contact & Address
			address TEXT,
			phone VARCHAR(20),
			
			-- Parent/Guardian Info
			father_name VARCHAR(100),
			father_phone VARCHAR(20),
			father_occupation VARCHAR(100),
			mother_name VARCHAR(100),
			mother_phone VARCHAR(20),
			mother_occupation VARCHAR(100),
			guardian_name VARCHAR(100),
			guardian_phone VARCHAR(20),
			guardian_relation VARCHAR(50),
			
			-- Type & Classification
			type VARCHAR(20) NOT NULL DEFAULT 'siswa', -- siswa, santri, both
			jenjang VARCHAR(10),                        -- TK, SD, MI, SMP, MTs, SMA, MA, SMK
			
			-- Academic Relations
			kelas_id UUID REFERENCES sekolah_kelas(id),
			
			-- Pesantren Relations (for santri)
			kamar_id UUID REFERENCES sekolah_kamar(id),
			is_mukim BOOLEAN DEFAULT false,
			tahun_masuk INTEGER,
			
			-- Status
			status VARCHAR(20) DEFAULT 'active', -- active, graduated, dropped_out, transferred
			entry_date DATE DEFAULT CURRENT_DATE,
			exit_date DATE,
			exit_reason TEXT,
			
			-- Timestamps
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		-- Indexes for query optimization
		CREATE INDEX IF NOT EXISTS idx_students_tenant ON students(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_students_type ON students(tenant_id, type);
		CREATE INDEX IF NOT EXISTS idx_students_kelas ON students(kelas_id);
		CREATE INDEX IF NOT EXISTS idx_students_kamar ON students(kamar_id);
		CREATE INDEX IF NOT EXISTS idx_students_jenjang ON students(tenant_id, jenjang);
		CREATE INDEX IF NOT EXISTS idx_students_status ON students(tenant_id, status);
		
		-- Unique constraints for NIS (per tenant)
		CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nis_tenant 
			ON students(tenant_id, nis) WHERE nis IS NOT NULL AND nis != '';
		CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nisn_tenant 
			ON students(tenant_id, nisn) WHERE nisn IS NOT NULL AND nisn != '';
		CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nis_pesantren_tenant 
			ON students(tenant_id, nis_pesantren) WHERE nis_pesantren IS NOT NULL AND nis_pesantren != '';
		
		-- Trigger for updated_at
		CREATE TRIGGER update_students_updated_at
			BEFORE UPDATE ON students
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downUnifiedStudents(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TRIGGER IF EXISTS update_students_updated_at ON students;
		DROP TABLE IF EXISTS students;
	`)
	return err
}
