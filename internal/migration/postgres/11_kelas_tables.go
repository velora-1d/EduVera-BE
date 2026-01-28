package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upKelasTables, downKelasTables)
}

func upKelasTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sekolah_kelas (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			nama VARCHAR(100) NOT NULL,
			tingkat VARCHAR(50) NOT NULL,
			urutan INT DEFAULT 0,
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_sekolah_kelas_tenant ON sekolah_kelas(tenant_id);
	`)
	return err
}

func downKelasTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS sekolah_kelas;
	`)
	return err
}
