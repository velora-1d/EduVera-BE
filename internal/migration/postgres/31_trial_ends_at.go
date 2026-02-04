package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTrialEndsAt, downTrialEndsAt)
}

// upTrialEndsAt adds trial_ends_at column for trial period tracking
func upTrialEndsAt(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		-- Add trial_ends_at column to tenants table
		ALTER TABLE tenants ADD COLUMN IF NOT EXISTS trial_ends_at TIMESTAMP WITH TIME ZONE;

		-- Create index for efficient queries on trial expiry
		CREATE INDEX IF NOT EXISTS idx_tenants_trial_ends_at ON tenants(trial_ends_at) WHERE trial_ends_at IS NOT NULL;

		-- Add comment for documentation
		COMMENT ON COLUMN tenants.trial_ends_at IS 'Timestamp when trial period ends. NULL means no trial or trial cleared after upgrade.';
	`)
	return err
}

func downTrialEndsAt(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP INDEX IF EXISTS idx_tenants_trial_ends_at;
		ALTER TABLE tenants DROP COLUMN IF EXISTS trial_ends_at;
	`)
	return err
}
