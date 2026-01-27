package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTenants, downTenants)
}

func upTenants(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(50) UNIQUE NOT NULL,
			plan_type VARCHAR(20) NOT NULL DEFAULT 'sekolah',
			institution_type VARCHAR(20),
			address TEXT,
			bank_name VARCHAR(100),
			account_number VARCHAR(50),
			account_holder VARCHAR(255),
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON tenants(subdomain);
		CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
	`)
	return err
}

func downTenants(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP INDEX IF EXISTS idx_tenants_status;
		DROP INDEX IF EXISTS idx_tenants_subdomain;
		DROP TABLE IF EXISTS tenants;
	`)
	return err
}
