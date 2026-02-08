package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up36WhatsAppPhoneNumber, down36WhatsAppPhoneNumber)
}

func up36WhatsAppPhoneNumber(ctx context.Context, tx *sql.Tx) error {
	// Add phone_number column if not exists
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE tenant_whatsapp_sessions 
		ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20)
	`)
	if err != nil {
		// Column might already exist, continue
	}

	// Drop existing unique constraint on tenant_id if exists
	_, _ = tx.ExecContext(ctx, `
		ALTER TABLE tenant_whatsapp_sessions 
		DROP CONSTRAINT IF EXISTS tenant_whatsapp_sessions_tenant_id_key
	`)

	// Add unique constraint on instance_name (allows owner session with no tenant_id)
	_, _ = tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS tenant_whatsapp_sessions_instance_name_key 
		ON tenant_whatsapp_sessions (instance_name)
	`)

	// Make tenant_id nullable for owner sessions
	_, _ = tx.ExecContext(ctx, `
		ALTER TABLE tenant_whatsapp_sessions 
		ALTER COLUMN tenant_id DROP NOT NULL
	`)

	return nil
}

func down36WhatsAppPhoneNumber(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE tenant_whatsapp_sessions DROP COLUMN IF EXISTS phone_number`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DROP INDEX IF EXISTS tenant_whatsapp_sessions_instance_name_key`)
	return err
}
