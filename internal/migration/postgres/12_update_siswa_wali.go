package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateSiswaWali, downUpdateSiswaWali)
}

func upUpdateSiswaWali(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE sekolah_siswa 
		ADD COLUMN IF NOT EXISTS nama_wali VARCHAR(100),
		ADD COLUMN IF NOT EXISTS no_hp_wali VARCHAR(20);
	`)
	return err
}

func downUpdateSiswaWali(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE sekolah_siswa 
		DROP COLUMN IF EXISTS nama_wali,
		DROP COLUMN IF EXISTS no_hp_wali;
	`)
	return err
}
