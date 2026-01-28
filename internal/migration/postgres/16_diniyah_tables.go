package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDiniyah, downDiniyah)
}

func upDiniyah(ctx context.Context, tx *sql.Tx) error {
	// 1. Diniyah Kitab (Subject/Book)
	query := `
			CREATE TABLE IF NOT EXISTS sekolah_diniyah_kitab (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				tenant_id UUID NOT NULL,
				nama_kitab VARCHAR(255) NOT NULL,
				bidang_studi VARCHAR(100), -- Fiqih, Nahwu, Akhlak
				pengarang VARCHAR(255),
				keterangan TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			);
			CREATE TRIGGER update_sekolah_diniyah_kitab_updated_at
				BEFORE UPDATE ON sekolah_diniyah_kitab
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create sekolah_diniyah_kitab: %w", err)
	}

	return nil
}

func downDiniyah(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS sekolah_diniyah_kitab`); err != nil {
		return err
	}
	return nil
}
