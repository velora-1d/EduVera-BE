package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upEraporsTables, downEraporsTables)
}

func upEraporsTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS subjects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			code VARCHAR(20) NOT NULL,
			name VARCHAR(255) NOT NULL,
			grade_level INTEGER NOT NULL,
			curriculum VARCHAR(50) NOT NULL,
			is_muatan_lokal BOOLEAN DEFAULT FALSE,
			grading_config JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_subjects_tenant_id ON subjects(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_subjects_grade_level ON subjects(grade_level);
		CREATE INDEX IF NOT EXISTS idx_subjects_curriculum ON subjects(curriculum);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_subjects_code_tenant ON subjects(tenant_id, code);

		CREATE TABLE IF NOT EXISTS student_grades (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			student_id UUID NOT NULL,
			subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
			semester INTEGER NOT NULL CHECK (semester IN (1, 2)),
			academic_year VARCHAR(9) NOT NULL,
			
			nilai_tugas DECIMAL(5,2),
			nilai_ulangan_harian DECIMAL(5,2),
			nilai_uts DECIMAL(5,2),
			nilai_uas DECIMAL(5,2),
			nilai_pengetahuan DECIMAL(5,2),
			
			nilai_praktik DECIMAL(5,2),
			nilai_proyek DECIMAL(5,2),
			nilai_portofolio DECIMAL(5,2),
			nilai_keterampilan DECIMAL(5,2),
			
			nilai_akhir DECIMAL(5,2),
			predikat CHAR(1),
			
			deskripsi_pengetahuan TEXT,
			deskripsi_keterampilan TEXT,
			catatan_guru TEXT,
			
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_student_grades_tenant_id ON student_grades(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_student_grades_student_id ON student_grades(student_id);
		CREATE INDEX IF NOT EXISTS idx_student_grades_subject_id ON student_grades(subject_id);
		CREATE INDEX IF NOT EXISTS idx_student_grades_semester ON student_grades(semester);
		CREATE INDEX IF NOT EXISTS idx_student_grades_academic_year ON student_grades(academic_year);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_student_grades_unique ON student_grades(student_id, subject_id, semester, academic_year);

		DROP TRIGGER IF EXISTS update_subjects_updated_at ON subjects;
		CREATE TRIGGER update_subjects_updated_at BEFORE UPDATE ON subjects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

		DROP TRIGGER IF EXISTS update_student_grades_updated_at ON student_grades;
		CREATE TRIGGER update_student_grades_updated_at BEFORE UPDATE ON student_grades FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downEraporsTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE IF EXISTS student_grades;
		DROP TABLE IF EXISTS subjects;
	`)
	return err
}
