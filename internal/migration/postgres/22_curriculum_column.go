package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCurriculumColumn, downCurriculumColumn)
}

func upCurriculumColumn(ctx context.Context, tx *sql.Tx) error {
	query := `
		ALTER TABLE sekolah_profil
		ADD COLUMN IF NOT EXISTS curriculum VARCHAR(20) DEFAULT 'K13';
		
		COMMENT ON COLUMN sekolah_profil.curriculum IS 'Kurikulum: K13, MERDEKA, PESANTREN';
	`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to add curriculum column: %w", err)
	}
	return nil
}

func downCurriculumColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE sekolah_profil DROP COLUMN IF EXISTS curriculum`)
	return err
}
