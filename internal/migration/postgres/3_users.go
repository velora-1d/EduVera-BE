package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUsers, downUsers)
}

func upUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			whatsapp VARCHAR(20),
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) DEFAULT 'admin',
			is_active BOOLEAN DEFAULT false,
			email_verified_at TIMESTAMP,
			last_login_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
	`)
	return err
}

func downUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP INDEX IF EXISTS idx_users_role;
		DROP INDEX IF EXISTS idx_users_email;
		DROP INDEX IF EXISTS idx_users_tenant_id;
		DROP TABLE IF EXISTS users;
	`)
	return err
}
