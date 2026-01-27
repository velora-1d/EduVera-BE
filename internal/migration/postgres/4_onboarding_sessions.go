package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOnboardingSessions, downOnboardingSessions)
}

func upOnboardingSessions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS onboarding_sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			current_step INT DEFAULT 1,
			data JSONB DEFAULT '{}',
			expires_at TIMESTAMP,
			completed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_onboarding_tenant_id ON onboarding_sessions(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_onboarding_user_id ON onboarding_sessions(user_id);
	`)
	return err
}

func downOnboardingSessions(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP INDEX IF EXISTS idx_onboarding_user_id;
		DROP INDEX IF EXISTS idx_onboarding_tenant_id;
		DROP TABLE IF EXISTS onboarding_sessions;
	`)
	return err
}
