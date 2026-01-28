package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCurriculumTables, downCurriculumTables)
}

func upCurriculumTables(ctx context.Context, tx *sql.Tx) error {
	// 1. Curriculum References (Master Data)
	// Defines available curriculums in the system
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS curriculum_references (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(50) NOT NULL UNIQUE, -- e.g. "K13_REVISI", "MERDEKA", "PESANTREN_SALAF"
			name VARCHAR(100) NOT NULL,
			description TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		INSERT INTO curriculum_references (code, name, description) VALUES
		('K13_REVISI', 'Kurikulum 2013 Revisi', 'Menggunakan KKM dan penilaian Pengetahuan/Keterampilan'),
		('MERDEKA', 'Kurikulum Merdeka', 'Menggunakan Capaian Pembelajaran dan penilaian Formatif/Sumatif'),
		('PESANTREN_SALAF', 'Kurikulum Pesantren Salaf', 'Penilaian berbasis Kitab Kuning dan hafalan')
		ON CONFLICT (code) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to create curriculum_references: %w", err)
	}

	// 2. Grading Formulas (Config per Subject/Curriculum)
	// Replaces the JSONB inside subjects table with a more structured approach if needed,
	// but strictly we can keep using subjects.grading_config.
	// However, we need a way to version KKM history.

	// 3. KKM History (For K13 mainly) - Tracks KKM changes over semesters
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS kkm_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
			semester_id UUID NOT NULL, -- Logical link to a semester reference
			kkm_value INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(subject_id, semester_id)
		);
		CREATE INDEX IF NOT EXISTS idx_kkm_history_tenant ON kkm_history(tenant_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create kkm_history: %w", err)
	}

	// 4. Update student_grades table to support JSONB scoring
	// We keep old columns for backward compatibility but add flexible columns
	_, err = tx.Exec(`
		ALTER TABLE student_grades 
		ADD COLUMN IF NOT EXISTS component_scores JSONB DEFAULT '{}',
		ADD COLUMN IF NOT EXISTS grading_formula_id UUID,
		ADD COLUMN IF NOT EXISTS description_high TEXT, -- Capaian tertinggi (Merdeka)
		ADD COLUMN IF NOT EXISTS description_low TEXT;  -- Perlu bimbingan (Merdeka)
		
		CREATE INDEX IF NOT EXISTS idx_student_grades_component_scores ON student_grades USING gin(component_scores);
	`)
	if err != nil {
		return fmt.Errorf("failed to alter student_grades: %w", err)
	}

	return nil
}

func downCurriculumTables(ctx context.Context, tx *sql.Tx) error {
	// Reverting changes
	_, err := tx.Exec(`
		ALTER TABLE student_grades 
		DROP COLUMN IF EXISTS description_low,
		DROP COLUMN IF EXISTS description_high,
		DROP COLUMN IF EXISTS grading_formula_id,
		DROP COLUMN IF EXISTS component_scores;

		DROP TABLE IF EXISTS kkm_history;
		DROP TABLE IF EXISTS curriculum_references;
	`)
	return err
}
