package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upContent, downContent)
}

func upContent(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS contents (
			key VARCHAR(255) PRIMARY KEY,
			value TEXT NOT NULL,
			type VARCHAR(50) NOT NULL DEFAULT 'text',
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func downContent(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS contents")
	return err
}
