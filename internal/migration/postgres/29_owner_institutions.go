package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOwnerInstitutions, downOwnerInstitutions)
}

func upOwnerInstitutions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		-- Add is_owner column to users table
		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_owner BOOLEAN NOT NULL DEFAULT FALSE;
		
		-- Create owner_institutions table for owner's own institutions
		CREATE TABLE IF NOT EXISTS owner_institutions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			relationship_type VARCHAR(50) NOT NULL DEFAULT 'owned',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(owner_user_id, tenant_id)
		);

		CREATE INDEX IF NOT EXISTS idx_owner_institutions_owner ON owner_institutions(owner_user_id);
		CREATE INDEX IF NOT EXISTS idx_owner_institutions_tenant ON owner_institutions(tenant_id);
	`)
	return err
}

func downOwnerInstitutions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS owner_institutions;
		ALTER TABLE users DROP COLUMN IF EXISTS is_owner;
	`)
	return err
}
